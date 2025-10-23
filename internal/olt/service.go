package olt

import (
	"context"
	"fmt"
	"time"
)

// Service provides OLT operations
type Service struct {
	timeout time.Duration
}

// NewService creates a new OLT service
func NewService(timeout time.Duration) *Service {
	return &Service{
		timeout: timeout,
	}
}

// OLTRequest represents a request to OLT device
type OLTRequest struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Prompt   string `json:"prompt"`
	Commands []string `json:"commands"`
}

// OLTResponse represents response from OLT device
type OLTResponse struct {
	Host    string `json:"host"`
	Output  string `json:"output"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Time    string `json:"execution_time"`
}

// ExecuteCommands executes commands on OLT device
func (s *Service) ExecuteCommands(ctx context.Context, req OLTRequest) (*OLTResponse, error) {
	start := time.Now()

	// Set total timeout
	totalCtx, cancel := context.WithTimeout(ctx, s.timeout*3)
	defer cancel()

	// Create session
	sess, err := NewSession(req.Host, req.Port, req.User, req.Password, req.Prompt, s.timeout)
	if err != nil {
		return &OLTResponse{
			Host:    req.Host,
			Success: false,
			Error:   fmt.Sprintf("session creation failed: %v", err),
			Time:    time.Since(start).String(),
		}, nil
	}
	defer sess.Close()

	// Login
	header := fmt.Sprintf("== %s:%d ==\n", req.Host, req.Port)
	_, err = sess.Login(totalCtx)
	if err != nil {
		return &OLTResponse{
			Host:    req.Host,
			Output:  header,
			Success: false,
			Error:   fmt.Sprintf("login failed: %v", err),
			Time:    time.Since(start).String(),
		}, nil
	}

	// Disable paging
	_, _ = sess.Exec(totalCtx, "terminal length 0")
	_, _ = sess.Exec(totalCtx, "screen-length 0 temporary")
	_, _ = sess.Exec(totalCtx, "disable clipaging")

	// Execute commands
	batchOut, err := sess.ExecBatch(totalCtx, req.Commands)
	output := header + batchOut

	success := err == nil
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	return &OLTResponse{
		Host:    req.Host,
		Output:  output,
		Success: success,
		Error:   errorMsg,
		Time:    time.Since(start).String(),
	}, nil
}

// ExecuteCommandsWithCustomTimeout executes commands with custom timeout
func (s *Service) ExecuteCommandsWithCustomTimeout(ctx context.Context, req OLTRequest, customTimeout time.Duration) (*OLTResponse, error) {
	start := time.Now()

	// Use custom timeout or default
	timeout := customTimeout
	if timeout == 0 {
		timeout = s.timeout
	}

	// Set total timeout with extended buffer for save operations
	totalCtx, cancel := context.WithTimeout(ctx, timeout*5) // 5x buffer for save operations
	defer cancel()

	// Create session with custom timeout
	sess, err := NewSession(req.Host, req.Port, req.User, req.Password, req.Prompt, timeout)
	if err != nil {
		return &OLTResponse{
			Host:    req.Host,
			Success: false,
			Error:   fmt.Sprintf("session creation failed: %v", err),
			Time:    time.Since(start).String(),
		}, nil
	}
	defer sess.Close()

	// Login
	header := fmt.Sprintf("== %s:%d ==\n", req.Host, req.Port)
	_, err = sess.Login(totalCtx)
	if err != nil {
		return &OLTResponse{
			Host:    req.Host,
			Output:  header,
			Success: false,
			Error:   fmt.Sprintf("login failed: %v", err),
			Time:    time.Since(start).String(),
		}, nil
	}

	// Disable paging
	_, _ = sess.Exec(totalCtx, "terminal length 0")
	_, _ = sess.Exec(totalCtx, "screen-length 0 temporary")
	_, _ = sess.Exec(totalCtx, "disable clipaging")

	// Execute commands
	batchOut, err := sess.ExecBatch(totalCtx, req.Commands)
	output := header + batchOut

	success := err == nil
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
		// Check for timeout specifically
		if totalCtx.Err() == context.DeadlineExceeded {
			errorMsg = "Operation timed out. The save configuration process may take several minutes on busy OLTs. Consider increasing the timeout parameter."
		}
	}

	return &OLTResponse{
		Host:    req.Host,
		Output:  output,
		Success: success,
		Error:   errorMsg,
		Time:    time.Since(start).String(),
	}, nil
}

// RenderCommand renders a single command for testing
func (s *Service) RenderCommand(req OLTRequest) *OLTResponse {
	return &OLTResponse{
		Host:    req.Host,
		Output:  fmt.Sprintf("Would execute: %v", req.Commands),
		Success: true,
		Time:    "0s",
	}
}