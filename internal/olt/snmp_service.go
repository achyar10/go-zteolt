package olt

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/gosnmp/gosnmp"
)

// SNMPRequest represents SNMP request parameters
type SNMPRequest struct {
	Host      string
	Port      int
	Community string
	BoardID   int
	PONID     int
	Timeout   int
}

// SNMPONUInfo represents ONU information from SNMP
type SNMPONUInfo struct {
	Board               int    `json:"board"`
	PON                 int    `json:"pon"`
	ID                  int    `json:"onu_id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	OnuType             string `json:"onu_type"`
	SerialNumber        string `json:"serial_number"`
	RXPower             string `json:"rx_power"`
	TXPower             string `json:"tx_power"`
	Status              string `json:"status"`
	IPAddress           string `json:"ip_address"`
	LastOnline          string `json:"last_online"`
	Uptime              string `json:"uptime"`
	GponOpticalDistance string `json:"gpon_optical_distance"`
}

// SNMPResult represents the result of SNMP query
type SNMPResult struct {
	Host          string
	BoardID       int
	PONID         int
	TotalONUs     int
	ONUs          []SNMPONUInfo
	ExecutionTime string
	Timestamp     time.Time
}

// OltConfig represents OLT configuration for specific board and PON
type OltConfig struct {
	BaseOID                   string
	OnuIDNameOID              string
	OnuTypeOID                string
	OnuSerialNumberOID        string
	OnuRxPowerOID             string
	OnuTxPowerOID             string
	OnuStatusOID              string
	OnuIPAddressOID           string
	OnuDescriptionOID         string
	OnuLastOnlineOID          string
	OnuLastOfflineOID         string
	OnuLastOfflineReasonOID   string
	OnuGponOpticalDistanceOID string
}

// SNMPService represents SNMP service for OLT monitoring
type SNMPService struct {
	timeout time.Duration
}

// NewFinalSNMPService creates a new SNMP service instance
func NewFinalSNMPService(timeout time.Duration) *SNMPService {
	return &SNMPService{
		timeout: timeout,
	}
}

// GetONUByBoardAndPON retrieves all ONUs for a specific board and PON
func (s *SNMPService) GetONUByBoardAndPON(ctx context.Context, req SNMPRequest) (*SNMPResult, error) {
	startTime := time.Now()

	// Get OLT configuration
	oltConfig, err := s.getOltConfig(req.BoardID, req.PONID)
	if err != nil {
		return nil, fmt.Errorf("failed to get OLT config: %w", err)
	}

	// Setup SNMP connection
	snmp, err := s.setupSNMPConnection(req.Host, req.Port, req.Community)
	if err != nil {
		return nil, fmt.Errorf("failed to setup SNMP connection: %w", err)
	}
	defer snmp.Conn.Close()

	// SNMP Walk to get ONU ID and Name
	snmpDataMap := make(map[string]gosnmp.SnmpPDU)
	err = snmp.Walk(oltConfig.BaseOID+oltConfig.OnuIDNameOID, func(pdu gosnmp.SnmpPDU) error {
		snmpDataMap[ExtractONUID(pdu.Name)] = pdu
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("SNMP walk failed: %w", err)
	}

	var onuInformationList []SNMPONUInfo

	// Process SNMP data to get ONU information
	for _, pdu := range snmpDataMap {
		onuInfo := SNMPONUInfo{
			Board: req.BoardID,
			PON:   req.PONID,
			ID:    ExtractIDOnuID(pdu.Name),
			Name:  ExtractName(pdu.Value),
		}

		// Get ONU type (like "F609V5.3", "F660V7.0", etc.)
		if onuType, err := s.getONUType(snmp, oltConfig, strconv.Itoa(onuInfo.ID)); err == nil {
			onuInfo.OnuType = onuType
		}

		if description, err := s.getONUDescription(snmp, oltConfig, strconv.Itoa(onuInfo.ID)); err == nil {
			onuInfo.Description = description
		}

		if sn, err := s.getSerialNumber(snmp, oltConfig, strconv.Itoa(onuInfo.ID)); err == nil {
			onuInfo.SerialNumber = sn
		}

		if rx, err := s.getRxPower(snmp, oltConfig, strconv.Itoa(onuInfo.ID)); err == nil {
			onuInfo.RXPower = rx
		}

		if tx, err := s.getTxPower(snmp, oltConfig, strconv.Itoa(onuInfo.ID)); err == nil {
			onuInfo.TXPower = tx
		}

		if status, err := s.getStatus(snmp, oltConfig, strconv.Itoa(onuInfo.ID)); err == nil {
			onuInfo.Status = status
		}

		if ip, err := s.getIPAddress(snmp, oltConfig, strconv.Itoa(onuInfo.ID)); err == nil {
			onuInfo.IPAddress = ip
		}

		if lastOnline, err := s.getLastOnline(snmp, oltConfig, strconv.Itoa(onuInfo.ID)); err == nil {
			onuInfo.LastOnline = lastOnline
			if uptime, err := s.calculateUptime(lastOnline); err == nil {
				onuInfo.Uptime = uptime
			}
		}

		// Get optical distance (loss)
		if distance, err := s.getOpticalDistance(snmp, oltConfig, strconv.Itoa(onuInfo.ID)); err == nil {
			onuInfo.GponOpticalDistance = distance
		}

		onuInformationList = append(onuInformationList, onuInfo)
	}

	// Sort ONU information by ID
	sort.Slice(onuInformationList, func(i, j int) bool {
		return onuInformationList[i].ID < onuInformationList[j].ID
	})

	// Prepare result
	result := &SNMPResult{
		Host:          req.Host,
		BoardID:       req.BoardID,
		PONID:         req.PONID,
		TotalONUs:     len(onuInformationList),
		ONUs:          onuInformationList,
		ExecutionTime: fmt.Sprintf("%.2fs", time.Since(startTime).Seconds()),
		Timestamp:     time.Now(),
	}

	return result, nil
}

// GetONUDetails retrieves specific ONU details by board, PON, and ONU ID
func (s *SNMPService) GetONUDetails(ctx context.Context, req SNMPRequest, onuID int) (*SNMPONUInfo, error) {

	// Get OLT configuration
	oltConfig, err := s.getOltConfig(req.BoardID, req.PONID)
	if err != nil {
		return nil, fmt.Errorf("failed to get OLT config: %w", err)
	}

	// Setup SNMP connection
	snmp, err := s.setupSNMPConnection(req.Host, req.Port, req.Community)
	if err != nil {
		return nil, fmt.Errorf("failed to setup SNMP connection: %w", err)
	}
	defer snmp.Conn.Close()

	// Create ONU info with basic info first
	onuInfo := SNMPONUInfo{
		Board: req.BoardID,
		PON:   req.PONID,
		ID:    onuID,
	}

	// Get all ONU information with direct SNMP GET calls (no walk)
	onuIDStr := strconv.Itoa(onuID)

	// Get ONU name
	if name, err := s.getONUName(snmp, oltConfig, onuIDStr); err == nil {
		onuInfo.Name = name
	}

	// Get ONU type
	if onuType, err := s.getONUType(snmp, oltConfig, onuIDStr); err == nil {
		onuInfo.OnuType = onuType
	}

	// Get description
	if description, err := s.getONUDescription(snmp, oltConfig, onuIDStr); err == nil {
		onuInfo.Description = description
	}

	// Get serial number
	if sn, err := s.getSerialNumber(snmp, oltConfig, onuIDStr); err == nil {
		onuInfo.SerialNumber = sn
	}

	// Get RX power
	if rx, err := s.getRxPower(snmp, oltConfig, onuIDStr); err == nil {
		onuInfo.RXPower = rx
	}

	// Get TX power
	if tx, err := s.getTxPower(snmp, oltConfig, onuIDStr); err == nil {
		onuInfo.TXPower = tx
	}

	// Get status
	if status, err := s.getStatus(snmp, oltConfig, onuIDStr); err == nil {
		onuInfo.Status = status
	}

	// Get IP address
	if ip, err := s.getIPAddress(snmp, oltConfig, onuIDStr); err == nil {
		onuInfo.IPAddress = ip
	}

	// Get last online
	if lastOnline, err := s.getLastOnline(snmp, oltConfig, onuIDStr); err == nil {
		onuInfo.LastOnline = lastOnline
		if uptime, err := s.calculateUptime(lastOnline); err == nil {
			onuInfo.Uptime = uptime
		}
	}

	// Get optical distance
	if distance, err := s.getOpticalDistance(snmp, oltConfig, onuIDStr); err == nil {
		onuInfo.GponOpticalDistance = distance
	}

	return &onuInfo, nil
}

// setupSNMPConnection sets up SNMP connection
func (s *SNMPService) setupSNMPConnection(host string, port int, community string) (*gosnmp.GoSNMP, error) {
	snmp := &gosnmp.GoSNMP{
		Target:    host,
		Port:      uint16(port),
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   s.timeout,
		Retries:   3,
	}

	err := snmp.Connect()
	if err != nil {
		return nil, err
	}

	return snmp, nil
}

// getOltConfig gets OLT configuration based on board and PON ID
func (s *SNMPService) getOltConfig(boardID, ponID int) (*OltConfig, error) {
	// Base OIDs from working configuration
	baseOID1 := ".1.3.6.1.4.1.3902.1082"

	// Calculate interface index based on board and PON ID
	interfaceIndex := calculateInterfaceIndex(boardID, ponID)
	onuTypeIndex := calculateONUTypeIndex(boardID, ponID)

	config := &OltConfig{
		BaseOID:                   baseOID1,
		OnuIDNameOID:              ".500.10.2.3.3.1.2." + strconv.Itoa(interfaceIndex),
		OnuTypeOID:                ".3.50.11.2.1.17." + strconv.Itoa(onuTypeIndex),
		OnuSerialNumberOID:        ".500.10.2.3.3.1.18." + strconv.Itoa(interfaceIndex),
		OnuRxPowerOID:             ".500.20.2.2.2.1.10." + strconv.Itoa(interfaceIndex),
		OnuTxPowerOID:             ".3.50.12.1.1.14." + strconv.Itoa(onuTypeIndex),
		OnuStatusOID:              ".500.10.2.3.8.1.4." + strconv.Itoa(interfaceIndex),
		OnuIPAddressOID:           ".3.50.16.1.1.10." + strconv.Itoa(onuTypeIndex),
		OnuDescriptionOID:         ".500.10.2.3.3.1.3." + strconv.Itoa(interfaceIndex),
		OnuLastOnlineOID:          ".500.10.2.3.8.1.5." + strconv.Itoa(interfaceIndex),
		OnuLastOfflineOID:         ".500.10.2.3.8.1.6." + strconv.Itoa(interfaceIndex),
		OnuLastOfflineReasonOID:   ".500.10.2.3.8.1.7." + strconv.Itoa(interfaceIndex),
		OnuGponOpticalDistanceOID: ".500.10.2.3.10.1.2." + strconv.Itoa(interfaceIndex),
	}

	return config, nil
}

// calculateInterfaceIndex calculates interface index based on board and PON ID
func calculateInterfaceIndex(boardID, ponID int) int {
	// Based on the pattern from working configuration:
	// Board 1 PON 1: 285278465
	// Board 1 PON 2: 285278466
	// Board 2 PON 1: 285278721
	// Board 2 PON 4: 285278724

	baseIndex := 285278464 // Base index for Board 1 PON 0
	if boardID == 2 {
		baseIndex = 285278720 // Base index for Board 2 PON 0
	}

	return baseIndex + ponID
}

// calculateONUTypeIndex calculates ONU type index based on board and PON ID
func calculateONUTypeIndex(boardID, ponID int) int {
	// Based on the pattern from working configuration:
	// Board 1 PON 1: 268501248
	// Board 1 PON 2: 268501504
	// Board 2 PON 1: 268566784
	// Board 2 PON 4: 268567552

	baseIndex := 268501248 // Base index for Board 1 PON 1
	if boardID == 2 {
		baseIndex = 268566784 // Base index for Board 2 PON 1
	}

	offset := (ponID - 1) * 256
	return baseIndex + offset
}

// SNMP getter methods
func (s *SNMPService) getONUName(snmp *gosnmp.GoSNMP, config *OltConfig, onuID string) (string, error) {
	oid := config.BaseOID + config.OnuIDNameOID + "." + onuID
	result, err := snmp.Get([]string{oid})
	if err != nil {
		return "", err
	}
	if len(result.Variables) > 0 {
		return ExtractName(result.Variables[0].Value), nil
	}
	return "", fmt.Errorf("no response")
}

func (s *SNMPService) getONUType(snmp *gosnmp.GoSNMP, config *OltConfig, onuID string) (string, error) {
	oid := ".1.3.6.1.4.1.3902.1012" + config.OnuTypeOID + "." + onuID
	result, err := snmp.Get([]string{oid})
	if err != nil {
		return "", err
	}
	if len(result.Variables) > 0 {
		return ExtractName(result.Variables[0].Value), nil
	}
	return "", fmt.Errorf("no response")
}

func (s *SNMPService) getONUDescription(snmp *gosnmp.GoSNMP, config *OltConfig, onuID string) (string, error) {
	oid := config.BaseOID + config.OnuDescriptionOID + "." + onuID
	result, err := snmp.Get([]string{oid})
	if err != nil {
		return "", err
	}
	if len(result.Variables) > 0 {
		return ExtractName(result.Variables[0].Value), nil
	}
	return "", fmt.Errorf("no response")
}

func (s *SNMPService) getSerialNumber(snmp *gosnmp.GoSNMP, config *OltConfig, onuID string) (string, error) {
	oid := config.BaseOID + config.OnuSerialNumberOID + "." + onuID
	result, err := snmp.Get([]string{oid})
	if err != nil {
		return "", err
	}
	if len(result.Variables) > 0 {
		return ExtractSerialNumber(result.Variables[0].Value), nil
	}
	return "", fmt.Errorf("no response")
}

func (s *SNMPService) getRxPower(snmp *gosnmp.GoSNMP, config *OltConfig, onuID string) (string, error) {
	oid := config.BaseOID + config.OnuRxPowerOID + "." + onuID + ".1"
	result, err := snmp.Get([]string{oid})
	if err != nil {
		return "", err
	}
	if len(result.Variables) > 0 {
		// Use the same calculation as original project: ConvertAndMultiply
		return ConvertAndMultiply(result.Variables[0].Value)
	}
	return "", fmt.Errorf("no response")
}

func (s *SNMPService) getTxPower(snmp *gosnmp.GoSNMP, config *OltConfig, onuID string) (string, error) {
	oid := ".1.3.6.1.4.1.3902.1012" + config.OnuTxPowerOID + "." + onuID + ".1"
	result, err := snmp.Get([]string{oid})
	if err != nil {
		return "", err
	}
	if len(result.Variables) > 0 {
		// Use the same calculation as original project: ConvertAndMultiply
		return ConvertAndMultiply(result.Variables[0].Value)
	}
	return "", fmt.Errorf("no response")
}

func (s *SNMPService) getStatus(snmp *gosnmp.GoSNMP, config *OltConfig, onuID string) (string, error) {
	oid := config.BaseOID + config.OnuStatusOID + "." + onuID
	result, err := snmp.Get([]string{oid})
	if err != nil {
		return "", err
	}
	if len(result.Variables) > 0 {
		return ExtractAndGetStatus(result.Variables[0].Value), nil
	}
	return "", fmt.Errorf("no response")
}

func (s *SNMPService) getIPAddress(snmp *gosnmp.GoSNMP, config *OltConfig, onuID string) (string, error) {
	oid := ".1.3.6.1.4.1.3902.1012" + config.OnuIPAddressOID + "." + onuID + ".1"
	result, err := snmp.Get([]string{oid})
	if err != nil {
		return "", err
	}
	if len(result.Variables) > 0 {
		return ExtractName(result.Variables[0].Value), nil
	}
	return "", fmt.Errorf("no response")
}

func (s *SNMPService) getLastOnline(snmp *gosnmp.GoSNMP, config *OltConfig, onuID string) (string, error) {
	oid := config.BaseOID + config.OnuLastOnlineOID + "." + onuID
	result, err := snmp.Get([]string{oid})
	if err != nil {
		return "", err
	}
	if len(result.Variables) > 0 {
		if bytes, ok := result.Variables[0].Value.([]byte); ok {
			return ConvertByteArrayToDateTime(bytes)
		}
	}
	return "", fmt.Errorf("no response")
}

func (s *SNMPService) calculateUptime(lastOnline string) (string, error) {
	currentTime := time.Now()

	// Try multiple date formats
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05.000",
	}

	var lastOnlineTime time.Time
	var err error

	for _, format := range formats {
		lastOnlineTime, err = time.Parse(format, lastOnline)
		if err == nil {
			break
		}
	}

	if err != nil {
		// If parsing fails, assume the device is online and show current time as last online
		lastOnlineTime = currentTime.Add(-time.Minute) // Assume 1 minute uptime as fallback
	}

	duration := currentTime.Sub(lastOnlineTime)
	if duration < 0 {
		duration = time.Minute // Minimum uptime
	}

	return ConvertDurationToString(duration), nil
}

func (s *SNMPService) getOpticalDistance(snmp *gosnmp.GoSNMP, config *OltConfig, onuID string) (string, error) {
	oid := config.BaseOID + config.OnuGponOpticalDistanceOID + "." + onuID
	result, err := snmp.Get([]string{oid})
	if err != nil {
		return "", err
	}
	if len(result.Variables) > 0 {
		// Convert to km format
		if value, ok := result.Variables[0].Value.(int); ok {
			distance := float64(value) / 1000.0 // Convert to km
			return fmt.Sprintf("%.1fkm", distance), nil
		}
		return ExtractGponOpticalDistance(result.Variables[0].Value), nil
	}
	return "", fmt.Errorf("no response")
}
