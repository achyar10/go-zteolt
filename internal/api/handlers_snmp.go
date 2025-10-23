package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/achyar10/go-zteolt/internal/olt"
	"github.com/gofiber/fiber/v2"
)

// GetONUByBoardAndPON handles SNMP monitoring requests for ONU list
func (h *Handlers) GetONUByBoardAndPON(c *fiber.Ctx) error {
	var req SNMPMonitoringRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid request body"))
	}

	// Get board_id and pon_id from URL parameters
	boardID, err := strconv.Atoi(c.Params("board_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid board_id parameter"))
	}

	ponID, err := strconv.Atoi(c.Params("pon_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid pon_id parameter"))
	}

	// Validate parameters
	if boardID < 1 || boardID > 2 {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "board_id must be 1 or 2"))
	}

	if ponID < 1 || ponID > 16 {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "pon_id must be between 1 and 16"))
	}

	// Set default timeout if not specified
	timeout := 30 // seconds
	if req.Timeout > 0 {
		timeout = req.Timeout
	}

	// Initialize SNMP service with realistic data approach
	snmpService := olt.NewFinalSNMPService(time.Duration(timeout) * time.Second)

	// Convert API request to SNMP service request
	snmpReq := olt.SNMPRequest{
		Host:      req.Host,
		Port:      req.Port,
		Community: req.Community,
		BoardID:   boardID,
		PONID:     ponID,
		Timeout:   req.Timeout,
	}

	// Execute SNMP query
	ctx := c.Context()
	result, err := snmpService.GetONUByBoardAndPON(ctx, snmpReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			h.createAPIResponse(false, nil, fmt.Sprintf("SNMP query failed: %v", err)))
	}

	// Convert to API response format
	apiResponse := SNMPMonitoringResponse{
		Host:          result.Host,
		BoardID:       result.BoardID,
		PONID:         result.PONID,
		TotalONUs:     result.TotalONUs,
		ONUs:          convertToAPIONUInfo(result.ONUs),
		ExecutionTime: result.ExecutionTime,
		Timestamp:     result.Timestamp,
	}

	return c.JSON(h.createAPIResponse(true, apiResponse, ""))
}

// GetONUDetailsSNMP handles SNMP requests for specific ONU details
func (h *Handlers) GetONUDetailsSNMP(c *fiber.Ctx) error {
	var req SNMPONUDetailsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid request body"))
	}

	// Get parameters from URL
	boardID, err := strconv.Atoi(c.Params("board_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid board_id parameter"))
	}

	ponID, err := strconv.Atoi(c.Params("pon_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid pon_id parameter"))
	}

	onuID, err := strconv.Atoi(c.Params("onu_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid onu_id parameter"))
	}

	// Validate parameters
	if boardID < 1 || boardID > 2 {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "board_id must be 1 or 2"))
	}

	if ponID < 1 || ponID > 16 {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "pon_id must be between 1 and 16"))
	}

	if onuID < 1 || onuID > 128 {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "onu_id must be between 1 and 128"))
	}

	// Set default timeout if not specified
	timeout := 30 // seconds
	if req.Timeout > 0 {
		timeout = req.Timeout
	}

	// Initialize SNMP service with realistic data approach
	snmpService := olt.NewFinalSNMPService(time.Duration(timeout) * time.Second)

	// Convert API request to SNMP service request
	snmpReq := olt.SNMPRequest{
		Host:      req.Host,
		Port:      req.Port,
		Community: req.Community,
		BoardID:   boardID,
		PONID:     ponID,
		Timeout:   req.Timeout,
	}

	// Get specific ONU details directly (much faster - no walk needed)
	ctx := c.Context()
	targetONU, err := snmpService.GetONUDetails(ctx, snmpReq, onuID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			h.createAPIResponse(false, nil, fmt.Sprintf("SNMP query failed: %v", err)))
	}

	// Convert to API response format
	apiResponse := convertToAPIONUInfo([]olt.SNMPONUInfo{*targetONU})[0]

	return c.JSON(h.createAPIResponse(true, apiResponse, ""))
}

// GetEmptySlotsSNMP handles SNMP requests for empty ONU slots
func (h *Handlers) GetEmptySlotsSNMP(c *fiber.Ctx) error {
	var req SNMPEmptySlotsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid request body"))
	}

	// Get parameters from URL
	boardID, err := strconv.Atoi(c.Params("board_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid board_id parameter"))
	}

	ponID, err := strconv.Atoi(c.Params("pon_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid pon_id parameter"))
	}

	// Validate parameters
	if boardID < 1 || boardID > 2 {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "board_id must be 1 or 2"))
	}

	if ponID < 1 || ponID > 16 {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "pon_id must be between 1 and 16"))
	}

	// Set default timeout if not specified
	timeout := 30 // seconds
	if req.Timeout > 0 {
		timeout = req.Timeout
	}

	// Initialize SNMP service with realistic data approach
	snmpService := olt.NewFinalSNMPService(time.Duration(timeout) * time.Second)

	// Convert API request to SNMP service request
	snmpReq := olt.SNMPRequest{
		Host:      req.Host,
		Port:      req.Port,
		Community: req.Community,
		BoardID:   boardID,
		PONID:     ponID,
		Timeout:   req.Timeout,
	}

	// Get all ONUs for the board/PON
	ctx := c.Context()
	result, err := snmpService.GetONUByBoardAndPON(ctx, snmpReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			h.createAPIResponse(false, nil, fmt.Sprintf("SNMP query failed: %v", err)))
	}

	// Create a map of used ONU IDs
	usedONUIDs := make(map[int]bool)
	for _, onu := range result.ONUs {
		usedONUIDs[onu.ID] = true
	}

	// Find empty slots (1-128)
	var emptySlots []SNMPEmptySlot
	for i := 1; i <= 128; i++ {
		if !usedONUIDs[i] {
			emptySlots = append(emptySlots, SNMPEmptySlot{
				Board: boardID,
				PON:   ponID,
				ONUID: i,
			})
		}
	}

	// Convert to API response format
	apiResponse := SNMPEmptySlotsResponse{
		Host:          req.Host,
		BoardID:       boardID,
		PONID:         ponID,
		TotalEmpty:    len(emptySlots),
		EmptySlots:    emptySlots,
		ExecutionTime: result.ExecutionTime,
		Timestamp:     time.Now(),
	}

	return c.JSON(h.createAPIResponse(true, apiResponse, ""))
}

// Helper function to convert SNMP service ONU info to API model
func convertToAPIONUInfo(serviceONUs []olt.SNMPONUInfo) []SNMPONUInfo {
	var apiONUs []SNMPONUInfo
	for _, onu := range serviceONUs {
		apiONU := SNMPONUInfo{
			Board:        onu.Board,
			PON:          onu.PON,
			ID:           onu.ID,
			Name:         onu.Name,
			Description:  onu.Description,
			SerialNumber: onu.SerialNumber,
			OnuType:      onu.OnuType,
			RXPower:      onu.RXPower,
			TXPower:      onu.TXPower,
			Status:       onu.Status,
			IPAddress:    onu.IPAddress,
			LastOnline:   onu.LastOnline,
			Uptime:       onu.Uptime,
		}
		apiONUs = append(apiONUs, apiONU)
	}
	return apiONUs
}
