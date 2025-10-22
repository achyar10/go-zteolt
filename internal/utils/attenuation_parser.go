package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// AttenuationData represents parsed attenuation information
type AttenuationData struct {
	Host        string  `json:"host"`
	Slot        int     `json:"slot"`
	Port        int     `json:"port"`
	ONU         int     `json:"onu"`
	Direction   string  `json:"direction"`
	OLTRxPower  float64 `json:"olt_rx_power_dbm"`
	OLTTxPower  float64 `json:"olt_tx_power_dbm"`
	ONURxPower  float64 `json:"onu_rx_power_dbm"`
	ONUTxPower  float64 `json:"onu_tx_power_dbm"`
	Attenuation float64 `json:"attenuation_db"`
	Status      string  `json:"status"`
	RawOutput   string  `json:"raw_output,omitempty"`
}

// ParseAttenuationOutput parses the raw output from show pon power attenuation command
func ParseAttenuationOutput(host string, slot, port, onu int, rawOutput string) *AttenuationData {
	data := &AttenuationData{
		Host:      host,
		Slot:      slot,
		Port:      port,
		ONU:       onu,
		RawOutput: rawOutput,
		Status:    "unknown",
	}

	// Remove ANSI escape sequences and clean the output
	cleanOutput := cleanOutput(rawOutput)

	// Parse up direction (upstream)
	upData := parseDirection(cleanOutput, "up")
	if upData != nil {
		data.Direction = "up"
		data.OLTRxPower = upData.OLTRxPower
		data.ONUTxPower = upData.ONUTxPower
		data.Attenuation = upData.Attenuation
		data.Status = determineStatus(upData.Attenuation)
	}

	// Parse down direction (downstream)
	downData := parseDirection(cleanOutput, "down")
	if downData != nil {
		// If we already have up data, we need to merge or create separate entries
		if data.Direction == "up" {
			// Create combined response with both directions
			data.OLTTxPower = downData.OLTTxPower
			data.ONURxPower = downData.ONURxPower
			data.Direction = "both"
		} else {
			data.Direction = "down"
			data.OLTTxPower = downData.OLTTxPower
			data.ONURxPower = downData.ONURxPower
			data.Attenuation = downData.Attenuation
			data.Status = determineStatus(downData.Attenuation)
		}
	}

	return data
}

// DirectionData represents parsed data for one direction
type DirectionData struct {
	OLTRxPower  float64
	ONUTxPower  float64
	OLTTxPower  float64
	ONURxPower  float64
	Attenuation float64
}

// parseDirection parses data for specific direction (up/down)
func parseDirection(output, direction string) *DirectionData {
	// Pattern to match the attenuation output lines
	// Example: "up      Rx :-28.827(dbm)      Tx:2.200(dbm)        31.027(dB)"
	pattern := fmt.Sprintf(`(?i)%s\s+Rx\s*:(?P<rx>[-\d.]+)\s*\(dbm\)\s*Tx\s*:(?P<tx>[-\d.]+)\s*\(dbm\)\s+(?P<att>[-\d.]+)\s*\(db\)`, direction)

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(output)
	if matches == nil {
		return nil
	}

	// Get named subexpression indices
	names := re.SubexpNames()
	result := make(map[string]string)
	for i, match := range matches {
		if i > 0 && names[i] != "" {
			result[names[i]] = match
		}
	}

	data := &DirectionData{}

	// Parse Rx power
	if rxStr, ok := result["rx"]; ok {
		if rx, err := strconv.ParseFloat(rxStr, 64); err == nil {
			if direction == "up" {
				data.OLTRxPower = rx
			} else {
				data.ONURxPower = rx
			}
		}
	}

	// Parse Tx power
	if txStr, ok := result["tx"]; ok {
		if tx, err := strconv.ParseFloat(txStr, 64); err == nil {
			if direction == "up" {
				data.ONUTxPower = tx
			} else {
				data.OLTTxPower = tx
			}
		}
	}

	// Parse attenuation
	if attStr, ok := result["att"]; ok {
		if att, err := strconv.ParseFloat(attStr, 64); err == nil {
			data.Attenuation = att
		}
	}

	return data
}

// cleanOutput removes ANSI escape sequences and extra whitespace
func cleanOutput(output string) string {
	// Remove ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[mGKHJABCD]`)
	cleaned := ansiRegex.ReplaceAllString(output, "")

	// Remove \r characters
	cleaned = strings.ReplaceAll(cleaned, "\r", "")

	// Clean up extra whitespace
	lines := strings.Split(cleaned, "\n")
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// determineStatus determines the status based on attenuation value
func determineStatus(attenuation float64) string {
	if attenuation < 0 {
		return "error"
	} else if attenuation <= 10 {
		return "excellent"
	} else if attenuation <= 15 {
		return "good"
	} else if attenuation <= 25 {
		return "normal"
	} else if attenuation <= 30 {
		return "warning"
	} else {
		return "critical"
	}
}

// GetStatusDescription returns human readable status description
func GetStatusDescription(status string) string {
	switch status {
	case "excellent":
		return "Excellent signal quality (< 10 dB)"
	case "good":
		return "Good signal quality (10-15 dB)"
	case "normal":
		return "Normal signal quality (15-25 dB)"
	case "warning":
		return "Warning: High attenuation (25-30 dB)"
	case "critical":
		return "Critical: Very high attenuation (> 30 dB)"
	case "error":
		return "Error: Invalid data"
	default:
		return "Unknown status"
	}
}