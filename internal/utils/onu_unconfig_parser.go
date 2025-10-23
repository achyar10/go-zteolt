package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// UnconfiguredONU represents parsed ONU unconfigured data
type UnconfiguredONU struct {
	OLTIndex     string `json:"olt_index"`
	Model        string `json:"model"`
	SerialNumber string `json:"serial_number"`
	Board        int    `json:"board"`
	PON          int    `json:"pon"`
}

// UnconfiguredONUList represents list of unconfigured ONUs
type UnconfiguredONUList struct {
	Host          string                       `json:"host"`
	TotalCount    int                          `json:"total_count"`
	ONUs          []UnconfiguredONU            `json:"onus"`
	GroupedBySlot map[string][]UnconfiguredONU `json:"grouped_by_slot,omitempty"`
	RawOutput     string                       `json:"raw_output,omitempty"`
}

// ParseUnconfiguredONUOutput parses the raw output from show pon onu uncfg command
func ParseUnconfiguredONUOutput(host string, rawOutput string) *UnconfiguredONUList {
	data := &UnconfiguredONUList{
		Host:      host,
		ONUs:      []UnconfiguredONU{},
		RawOutput: rawOutput,
	}

	// Remove ANSI escape sequences and clean the output
	cleanOutput := cleanONUOutput(rawOutput)

	// Parse ONU entries
	onus := parseONUEntries(cleanOutput)
	data.ONUs = onus
	data.TotalCount = len(onus)

	// Group by slot for easier analysis
	data.GroupedBySlot = groupONUsBySlot(onus)

	return data
}

// parseONUEntries parses individual ONU entries from the output
func parseONUEntries(output string) []UnconfiguredONU {
	var onus []UnconfiguredONU

	// Pattern to match ONU lines
	// Example: "gpon-olt_1/1/14     F660V8.0                 RTEGC6A1BF4D"
	pattern := `^(gpon-olt_\d+/\d+/\d+)\s+(\w+\S*)\s+(\w+)$`

	re := regexp.MustCompile(pattern)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "OltIndex") || strings.Contains(line, "----") {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		// Extract slot and port from OLT index
		board, pon := parseOLTSlotPort(matches[1])

		onu := UnconfiguredONU{
			OLTIndex:     matches[1],
			Model:        matches[2],
			SerialNumber: matches[3],
			Board:        board,
			PON:          pon,
		}

		onus = append(onus, onu)
	}

	return onus
}

// parseOLTSlotPort extracts slot and port from OLT index string
func parseOLTSlotPort(oltIndex string) (int, int) {
	// Pattern: gpon-olt_1/1/14
	pattern := `gpon-olt_\d+/(\d+)/(\d+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(oltIndex)

	if len(matches) < 3 {
		return 0, 0
	}

	slot, _ := strconv.Atoi(matches[1])
	port, _ := strconv.Atoi(matches[2])

	return slot, port
}

// groupONUsBySlot groups ONUs by their slot
func groupONUsBySlot(onus []UnconfiguredONU) map[string][]UnconfiguredONU {
	grouped := make(map[string][]UnconfiguredONU)

	for _, onu := range onus {
		slotKey := strconv.Itoa(onu.Board)
		grouped[slotKey] = append(grouped[slotKey], onu)
	}

	return grouped
}

// cleanONUOutput removes ANSI escape sequences and extra whitespace
func cleanONUOutput(output string) string {
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

// GetONUModelDescription returns human readable model description
func GetONUModelDescription(model string) string {
	modelUpper := strings.ToUpper(model)

	switch {
	case strings.Contains(modelUpper, "F660"):
		return "ZTE F660 Series GPON ONU"
	case strings.Contains(modelUpper, "F620"):
		return "ZTE F620 Series GPON ONU"
	case strings.Contains(modelUpper, "F660V8"):
		return "ZTE F660V8 GPON ONU (Latest)"
	case strings.Contains(modelUpper, "F660V5"):
		return "ZTE F660V5 GPON ONU (Mid-gen)"
	case strings.Contains(modelUpper, "F660V3"):
		return "ZTE F660V3 GPON ONU (Older)"
	case strings.Contains(modelUpper, "F601"):
		return "ZTE F601 GPON ONU"
	case strings.Contains(modelUpper, "AN5506"):
		return "Fiberhome AN5506 Series ONU"
	case strings.Contains(modelUpper, "HG8245"):
		return "Huawei HG8245 Series ONU"
	default:
		return "Unknown ONU Model: " + model
	}
}

// GetUnconfiguredStatus returns status based on count
func GetUnconfiguredStatus(count int) string {
	switch {
	case count == 0:
		return "all_configured"
	case count <= 5:
		return "few_unconfigured"
	case count <= 15:
		return "some_unconfigured"
	default:
		return "many_unconfigured"
	}
}

// GetStatusDescription returns human readable status description
func GetUnconfiguredStatusDescription(status string) string {
	switch status {
	case "all_configured":
		return "All ONUs are properly configured"
	case "few_unconfigured":
		return "Few ONUs need configuration (1-5 units)"
	case "some_unconfigured":
		return "Several ONUs need configuration (6-15 units)"
	case "many_unconfigured":
		return "Many ONUs need configuration (>15 units)"
	default:
		return "Unknown status"
	}
}
