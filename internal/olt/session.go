package olt

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
	"time"
)

// Session represents a telnet session to OLT device
type Session struct {
	addr          string
	user          string
	pass          string
	promptPattern *regexp.Regexp
	conn          net.Conn
	timeout       time.Duration
	readBuf       bytes.Buffer
}

// NewSession creates a new OLT session
func NewSession(host string, port int, user, pass, promptRegex string, timeout time.Duration) (*Session, error) {
	if promptRegex == "" {
		promptRegex = `(?m)[>#]\s?$`
	}

	re, err := regexp.Compile(promptRegex)
	if err != nil {
		return nil, fmt.Errorf("invalid prompt regex: %w", err)
	}

	return &Session{
		addr:          fmt.Sprintf("%s:%d", host, port),
		user:          user,
		pass:          pass,
		promptPattern: re,
		timeout:       timeout,
	}, nil
}

// dial establishes connection to OLT
func (s *Session) dial() error {
	c, err := net.DialTimeout("tcp", s.addr, s.timeout)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}
	s.conn = c
	return nil
}

// Close closes the connection
func (s *Session) Close() {
	if s.conn != nil {
		_ = s.conn.Close()
	}
}

// stripIAC removes telnet negotiation bytes
func stripIAC(b []byte) []byte {
	out := make([]byte, 0, len(b))
	i := 0
	for i < len(b) {
		if b[i] == 255 { // IAC
			if i+2 < len(b) {
				i += 3
				continue
			}
			break
		}
		out = append(out, b[i])
		i++
	}
	return out
}

// readUntil reads until one of the patterns is matched
func (s *Session) readUntil(ctx context.Context, patterns ...*regexp.Regexp) (string, error) {
	if s.conn == nil {
		return "", errors.New("connection is nil")
	}

	buf := make([]byte, 4096)
	deadline := time.Now().Add(s.timeout)
	_ = s.conn.SetReadDeadline(deadline)

	for {
		select {
		case <-ctx.Done():
			return s.readBuf.String(), ctx.Err()
		default:
		}

		n, err := s.conn.Read(buf)
		if n > 0 {
			chunk := stripIAC(buf[:n])
			s.readBuf.Write(chunk)
			data := s.readBuf.String()
			for _, re := range patterns {
				if re.MatchString(data) {
					out := data
					s.readBuf.Reset()
					return out, nil
				}
			}
		}
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				return s.readBuf.String(), fmt.Errorf("read timeout")
			}
			if errors.Is(err, io.EOF) {
				return s.readBuf.String(), io.EOF
			}
			return s.readBuf.String(), err
		}
	}
}

// writeLine writes a line to the connection
func (s *Session) writeLine(line string) error {
	if s.conn == nil {
		return errors.New("connection is nil")
	}
	_ = s.conn.SetWriteDeadline(time.Now().Add(s.timeout))
	_, err := s.conn.Write([]byte(line + "\r\n"))
	return err
}

// Login authenticates to the OLT
func (s *Session) Login(ctx context.Context) (string, error) {
	if err := s.dial(); err != nil {
		return "", err
	}

	usernameRE := regexp.MustCompile(`(?i)(username|login)\s*:\s*$`)
	passwordRE := regexp.MustCompile(`(?i)password\s*:\s*$`)

	out1, _ := s.readUntil(ctx, usernameRE, passwordRE, s.promptPattern)

	if usernameRE.MatchString(out1) {
		if err := s.writeLine(s.user); err != nil {
			return out1, fmt.Errorf("write username: %w", err)
		}
		if _, err := s.readUntil(ctx, passwordRE); err != nil {
			return out1, fmt.Errorf("waiting password: %w", err)
		}
	}
	if passwordRE.MatchString(out1) || usernameRE.MatchString(out1) {
		if err := s.writeLine(s.pass); err != nil {
			return out1, fmt.Errorf("write password: %w", err)
		}
	}
	out2, err := s.readUntil(ctx, s.promptPattern)
	return out1 + out2, err
}

// Exec executes a single command
func (s *Session) Exec(ctx context.Context, cmd string) (string, error) {
	if strings.TrimSpace(cmd) == "" {
		return "", nil
	}
	if err := s.writeLine(cmd); err != nil {
		return "", err
	}
	out, err := s.readUntil(ctx, s.promptPattern)
	return out, err
}

// ExecBatch executes multiple commands
func (s *Session) ExecBatch(ctx context.Context, commands []string) (string, error) {
	var all strings.Builder
	for _, c := range commands {
		c = strings.TrimSpace(c)
		if c == "" || strings.HasPrefix(c, "#") {
			continue
		}

		out, err := s.Exec(ctx, c)
		all.WriteString(fmt.Sprintf(">>> %s\n%s\n", c, out))
		if err != nil {
			all.WriteString(fmt.Sprintf("ERR: %v\n", err))
		}
	}
	return all.String(), nil
}