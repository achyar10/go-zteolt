package api

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/achyar10/go-zteolt/internal/config"
	"github.com/achyar10/go-zteolt/internal/olt"
	"github.com/achyar10/go-zteolt/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// Handlers holds API handlers
type Handlers struct {
	oltService   *olt.Service
	templateMgr  *config.TemplateManager
	requestIDGen func() string
}

// NewHandlers creates new API handlers
func NewHandlers(oltService *olt.Service, templateMgr *config.TemplateManager) *Handlers {
	return &Handlers{
		oltService:  oltService,
		templateMgr: templateMgr,
		requestIDGen: func() string {
			return fmt.Sprintf("%d", time.Now().UnixNano())
		},
	}
}

// createAPIResponse creates standard API response
func (h *Handlers) createAPIResponse(success bool, data interface{}, errorMsg string) APIResponse {
	response := APIResponse{
		Success:   success,
		Data:      data,
		Error:     errorMsg,
		Timestamp: time.Now(),
		RequestID: h.requestIDGen(),
	}
	return response
}

// AddONU handles add ONU requests
func (h *Handlers) AddONU(c *fiber.Ctx) error {
	var req AddONURequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid request body"))
	}

	// Render commands using template
	commands, _, err := h.templateMgr.RenderTemplate("add-onu", map[string]interface{}{
		"Slot":         req.Slot,
		"Port":         req.OLTPort,
		"Onu":          req.ONU,
		"SerialNumber": req.SerialNumber,
		"Code":         req.Code,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			h.createAPIResponse(false, nil, "Template rendering failed"))
	}

	if req.RenderOnly {
		return c.JSON(h.createAPIResponse(true, ONUCommandResponse{
			Host:       req.Host,
			Mode:       "add-onu",
			Commands:   commands,
			RenderOnly: true,
			Success:    true,
		}, ""))
	}

	// Execute commands on OLT
	oltReq := olt.OLTRequest{
		Host:     req.Host,
		Port:     req.Port,
		User:     req.User,
		Password: req.Password,
		Commands: commands,
	}

	ctx := c.Context()
	result, err := h.oltService.ExecuteCommands(ctx, oltReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			h.createAPIResponse(false, nil, err.Error()))
	}

	response := ONUCommandResponse{
		Host:       result.Host,
		Mode:       "add-onu",
		Commands:   commands,
		Output:     result.Output,
		Success:    result.Success,
		Error:      result.Error,
		Time:       result.Time,
		RenderOnly: false,
	}

	return c.JSON(h.createAPIResponse(true, response, ""))
}

// DeleteONU handles delete ONU requests
func (h *Handlers) DeleteONU(c *fiber.Ctx) error {
	var req DeleteONURequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid request body"))
	}

	// Render commands using template
	commands, _, err := h.templateMgr.RenderTemplate("delete-onu", map[string]interface{}{
		"Slot": req.Slot,
		"Port": req.OLTPort,
		"Onu":  req.ONU,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			h.createAPIResponse(false, nil, "Template rendering failed"))
	}

	if req.RenderOnly {
		return c.JSON(h.createAPIResponse(true, ONUCommandResponse{
			Host:       req.Host,
			Mode:       "delete-onu",
			Commands:   commands,
			RenderOnly: true,
			Success:    true,
		}, ""))
	}

	// Execute commands on OLT
	oltReq := olt.OLTRequest{
		Host:     req.Host,
		Port:     req.Port,
		User:     req.User,
		Password: req.Password,
		Commands: commands,
	}

	ctx := c.Context()
	result, err := h.oltService.ExecuteCommands(ctx, oltReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			h.createAPIResponse(false, nil, err.Error()))
	}

	response := ONUCommandResponse{
		Host:       result.Host,
		Mode:       "delete-onu",
		Commands:   commands,
		Output:     result.Output,
		Success:    result.Success,
		Error:      result.Error,
		Time:       result.Time,
		RenderOnly: false,
	}

	return c.JSON(h.createAPIResponse(true, response, ""))
}

// CheckAttenuation handles check attenuation requests
func (h *Handlers) CheckAttenuation(c *fiber.Ctx) error {
	var req CheckAttenuationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid request body"))
	}

	if req.RenderOnly {
		return c.JSON(h.createAPIResponse(true, CheckAttenuationResponse{
			Host:       req.Host,
			Mode:       "check-attenuation",
			RenderOnly: true,
			Success:    true,
		}, ""))
	}

	// Render commands using template
	commands, _, err := h.templateMgr.RenderTemplate("check-attenuation", map[string]interface{}{
		"Slot": req.Slot,
		"Port": req.OLTPort,
		"Onu":  req.ONU,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			h.createAPIResponse(false, nil, "Template rendering failed"))
	}

	// Execute commands on OLT
	oltReq := olt.OLTRequest{
		Host:     req.Host,
		Port:     req.Port,
		User:     req.User,
		Password: req.Password,
		Commands: commands,
	}

	ctx := c.Context()
	result, err := h.oltService.ExecuteCommands(ctx, oltReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			h.createAPIResponse(false, nil, err.Error()))
	}

	// Parse the output to extract structured attenuation data
	var attenuationData *AttenuationDataDTO
	if result.Success && result.Output != "" {
		// Extract the actual output (remove headers and prompts)
		actualOutput := extractAttenuationOutput(result.Output)
		parsedData := utils.ParseAttenuationOutput(req.Host, req.Slot, req.OLTPort, req.ONU, actualOutput)

		// Convert to DTO
		if parsedData != nil {
			attenuationData = &AttenuationDataDTO{
				Host:        parsedData.Host,
				Slot:        parsedData.Slot,
				Port:        parsedData.Port,
				ONU:         parsedData.ONU,
				Direction:   parsedData.Direction,
				OLTRxPower:  parsedData.OLTRxPower,
				OLTTxPower:  parsedData.OLTTxPower,
				ONURxPower:  parsedData.ONURxPower,
				ONUTxPower:  parsedData.ONUTxPower,
				Attenuation: parsedData.Attenuation,
				Status:      parsedData.Status,
				RawOutput:   parsedData.RawOutput,
			}
		}
	}

	response := CheckAttenuationResponse{
		Host:       result.Host,
		Mode:       "check-attenuation",
		Success:    result.Success,
		Error:      result.Error,
		Time:       result.Time,
		RenderOnly: false,
		Data:       attenuationData,
	}

	return c.JSON(h.createAPIResponse(true, response, ""))
}

// extractAttenuationOutput extracts the actual attenuation data from the full command output
func extractAttenuationOutput(fullOutput string) string {
	lines := strings.Split(fullOutput, "\n")
	var outputLines []string
	inData := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip headers and prompts
		if strings.Contains(line, "===") || strings.HasPrefix(line, ">>>") ||
			strings.HasPrefix(line, "ZXAN") || line == "" {
			continue
		}

		// Start collecting data when we see relevant content
		if strings.Contains(line, "OLT") || strings.Contains(line, "ONU") ||
			strings.Contains(line, "Attenuation") || strings.Contains(line, "up") ||
			strings.Contains(line, "down") {
			inData = true
		}

		if inData {
			outputLines = append(outputLines, line)
		}
	}

	return strings.Join(outputLines, "\n")
}

// CheckUnconfigured handles check unconfigured ONU requests
func (h *Handlers) CheckUnconfigured(c *fiber.Ctx) error {
	var req CheckUnconfiguredRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid request body"))
	}

	if req.RenderOnly {
		return c.JSON(h.createAPIResponse(true, UnconfiguredONUListResponse{
			Host:       req.Host,
			Mode:       "check-unconfigured",
			RenderOnly: true,
			Success:    true,
		}, ""))
	}

	// Render commands using template
	commands, _, err := h.templateMgr.RenderTemplate("check-unconfigured", nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			h.createAPIResponse(false, nil, "Template rendering failed"))
	}

	// Execute commands on OLT
	oltReq := olt.OLTRequest{
		Host:     req.Host,
		Port:     req.Port,
		User:     req.User,
		Password: req.Password,
		Commands: commands,
	}

	ctx := c.Context()
	result, err := h.oltService.ExecuteCommands(ctx, oltReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			h.createAPIResponse(false, nil, err.Error()))
	}

	// Parse the output to extract structured unconfigured ONU data
	var unconfiguredData *UnconfiguredONUListDTO
	if result.Success && result.Output != "" {
		// Extract the actual output (remove headers and prompts)
		actualOutput := extractUnconfiguredOutput(result.Output)
		parsedData := utils.ParseUnconfiguredONUOutput(req.Host, actualOutput)

		// Convert to DTO
		if parsedData != nil {
			onus := make([]UnconfiguredONUDTO, len(parsedData.ONUs))
			for i, onu := range parsedData.ONUs {
				onus[i] = UnconfiguredONUDTO{
					OLTIndex:     onu.OLTIndex,
					Model:        onu.Model,
					SerialNumber: onu.SerialNumber,
					Slot:         onu.Slot,
					Port:         onu.Port,
				}
			}

			// Convert grouped by slot
			grouped := make(map[string][]UnconfiguredONUDTO)
			for slot, slotONUs := range parsedData.GroupedBySlot {
				converted := make([]UnconfiguredONUDTO, len(slotONUs))
				for i, onu := range slotONUs {
					converted[i] = UnconfiguredONUDTO{
						OLTIndex:     onu.OLTIndex,
						Model:        onu.Model,
						SerialNumber: onu.SerialNumber,
						Slot:         onu.Slot,
						Port:         onu.Port,
					}
				}
				grouped[slot] = converted
			}

			unconfiguredData = &UnconfiguredONUListDTO{
				Host:          parsedData.Host,
				TotalCount:    parsedData.TotalCount,
				ONUs:          onus,
				GroupedBySlot: grouped,
				Status:        utils.GetUnconfiguredStatus(parsedData.TotalCount),
				RawOutput:     parsedData.RawOutput,
			}
		}
	}

	response := UnconfiguredONUListResponse{
		Host:       result.Host,
		Mode:       "check-unconfigured",
		Success:    result.Success,
		Error:      result.Error,
		Time:       result.Time,
		RenderOnly: false,
		Data:       unconfiguredData,
	}

	return c.JSON(h.createAPIResponse(true, response, ""))
}

// extractUnconfiguredOutput extracts the actual unconfigured ONU data from the full command output
func extractUnconfiguredOutput(fullOutput string) string {
	lines := strings.Split(fullOutput, "\n")
	var outputLines []string
	inData := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip headers and prompts
		if strings.Contains(line, "===") || strings.HasPrefix(line, ">>>") ||
			strings.HasPrefix(line, "ZXAN") || line == "" {
			continue
		}

		// Start collecting data when we see relevant content
		if strings.Contains(line, "OltIndex") || strings.Contains(line, "Model") ||
			strings.Contains(line, "SN") || strings.Contains(line, "gpon-olt_") {
			inData = true
		}

		if inData {
			outputLines = append(outputLines, line)
		}
	}

	return strings.Join(outputLines, "\n")
}

// BatchCommands handles batch command requests
func (h *Handlers) BatchCommands(c *fiber.Ctx) error {
	var req BatchCommandsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			h.createAPIResponse(false, nil, "Invalid request body"))
	}

	// Execute commands on OLT
	oltReq := olt.OLTRequest{
		Host:     req.Host,
		Port:     req.Port,
		User:     req.User,
		Password: req.Password,
		Commands: req.Commands,
	}

	ctx := c.Context()
	result, err := h.oltService.ExecuteCommands(ctx, oltReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			h.createAPIResponse(false, nil, err.Error()))
	}

	response := ONUCommandResponse{
		Host:     result.Host,
		Mode:     "batch",
		Commands: req.Commands,
		Output:   result.Output,
		Success:  result.Success,
		Error:    result.Error,
		Time:     result.Time,
	}

	return c.JSON(h.createAPIResponse(true, response, ""))
}

// HealthCheck handles health check requests
func (h *Handlers) HealthCheck(c *fiber.Ctx) error {
	response := HealthCheckResponse{
		Status:  "healthy",
		Version: "1.0.0",
		Uptime:  "0h 0m 0s", // TODO: Calculate actual uptime
		Services: map[string]string{
			"olt":       "ok",
			"templates": "ok",
			"fiber":     "ok",
		},
		Timestamp: time.Now(),
	}

	return c.JSON(h.createAPIResponse(true, response, ""))
}

// ListTemplates handles template listing requests
func (h *Handlers) ListTemplates(c *fiber.Ctx) error {
	templates := h.templateMgr.GetAvailableTemplates()

	data := map[string]interface{}{
		"templates": templates,
		"count":     len(templates),
	}

	return c.JSON(h.createAPIResponse(true, data, ""))
}

// APIInfo handles root path requests
func (h *Handlers) APIInfo(c *fiber.Ctx) error {
	data := map[string]interface{}{
		"name":      "ZTE OLT Management API",
		"version":   "1.0.0",
		"status":    "running",
		"framework": "Fiber v2",
		"endpoints": map[string]string{
			"health":             "/api/v1/health",
			"templates":          "/api/v1/templates",
			"add_onu":            "/api/v1/onu/add",
			"delete_onu":         "/api/v1/onu/delete",
			"check_attenuation":  "/api/v1/onu/check-attenuation",
			"check_unconfigured": "/api/v1/onu/check-unconfigured",
			"batch_commands":     "/api/v1/batch/commands",
		},
	}

	return c.JSON(h.createAPIResponse(true, data, ""))
}

// LoggingMiddleware logs HTTP requests
func (h *Handlers) LoggingMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	// Log request
	log.Printf("%s %s %s", c.Method(), c.Path(), c.IP())

	// Continue to next handler
	err := c.Next()

	// Log response time
	log.Printf("%s %s completed in %v", c.Method(), c.Path(), time.Since(start))

	return err
}
