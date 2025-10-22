package api

import "time"

// Request/Response models for REST API

// AddONURequest represents request to add ONU
type AddONURequest struct {
	Host         string `json:"host" binding:"required"`
	Port         int    `json:"port" binding:"required"`
	User         string `json:"user" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Slot         int    `json:"slot" binding:"required"`
	OLTPort      int    `json:"olt_port" binding:"required"`
	ONU          int    `json:"onu" binding:"required"`
	SerialNumber string `json:"serial_number" binding:"required"`
	Code         string `json:"code" binding:"required"`
	RenderOnly   bool   `json:"render_only"`
}

// DeleteONURequest represents request to delete ONU
type DeleteONURequest struct {
	Host       string `json:"host" binding:"required"`
	Port       int    `json:"port" binding:"required"`
	User       string `json:"user" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Slot       int    `json:"slot" binding:"required"`
	OLTPort    int    `json:"olt_port" binding:"required"`
	ONU        int    `json:"onu" binding:"required"`
	RenderOnly bool   `json:"render_only"`
}

// CheckAttenuationRequest represents request to check attenuation
type CheckAttenuationRequest struct {
	Host       string `json:"host" binding:"required"`
	Port       int    `json:"port" binding:"required"`
	User       string `json:"user" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Slot       int    `json:"slot" binding:"required"`
	OLTPort    int    `json:"olt_port" binding:"required"`
	ONU        int    `json:"onu" binding:"required"`
	RenderOnly bool   `json:"render_only"`
}

// AttenuationData represents parsed attenuation information
type AttenuationData struct {
	Host        string        `json:"host"`
	Slot        int           `json:"slot"`
	Port        int           `json:"port"`
	ONU         int           `json:"onu"`
	Direction   string        `json:"direction"`
	OLTRxPower  float64       `json:"olt_rx_power_dbm"`
	OLTTxPower  float64       `json:"olt_tx_power_dbm"`
	ONURxPower  float64       `json:"onu_rx_power_dbm"`
	ONUTxPower  float64       `json:"onu_tx_power_dbm"`
	Attenuation float64       `json:"attenuation_db"`
	Status      string        `json:"status"`
	RawOutput   string        `json:"raw_output,omitempty"`
}

// CheckAttenuationResponse represents response for check attenuation
type CheckAttenuationResponse struct {
	Host       string              `json:"host"`
	Mode       string              `json:"mode"`
	Success    bool                `json:"success"`
	Error      string              `json:"error,omitempty"`
	Time       string              `json:"execution_time"`
	RenderOnly bool                `json:"render_only"`
	Data       *AttenuationDataDTO `json:"data,omitempty"`
}

// UnconfiguredONU represents parsed ONU unconfigured data
type UnconfiguredONUDTO struct {
	OLTIndex    string `json:"olt_index"`
	Model       string `json:"model"`
	SerialNumber string `json:"serial_number"`
	Slot        int    `json:"slot"`
	Port        int    `json:"port"`
}

// UnconfiguredONUListResponse represents response for check unconfigured ONUs
type UnconfiguredONUListResponse struct {
	Host          string                       `json:"host"`
	Mode          string                       `json:"mode"`
	Success       bool                         `json:"success"`
	Error         string                       `json:"error,omitempty"`
	Time          string                       `json:"execution_time"`
	RenderOnly    bool                         `json:"render_only"`
	Data          *UnconfiguredONUListDTO       `json:"data,omitempty"`
}

// UnconfiguredONUListDTO represents the data transfer object for unconfigured ONUs
type UnconfiguredONUListDTO struct {
	Host         string                     `json:"host"`
	TotalCount   int                        `json:"total_count"`
	ONUs         []UnconfiguredONUDTO       `json:"onus"`
	GroupedBySlot map[string][]UnconfiguredONUDTO `json:"grouped_by_slot,omitempty"`
	Status       string                     `json:"status"`
	RawOutput    string                     `json:"raw_output,omitempty"`
}

// AttenuationDataDTO represents the data transfer object for attenuation
type AttenuationDataDTO struct {
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

// CheckUnconfiguredRequest represents request to check unconfigured ONUs
type CheckUnconfiguredRequest struct {
	Host       string `json:"host" binding:"required"`
	Port       int    `json:"port" binding:"required"`
	User       string `json:"user" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RenderOnly bool   `json:"render_only"`
}

// BatchCommandsRequest represents request for batch commands
type BatchCommandsRequest struct {
	Host     string   `json:"host" binding:"required"`
	Port     int      `json:"port" binding:"required"`
	User     string   `json:"user" binding:"required"`
	Password string   `json:"password" binding:"required"`
	Commands []string `json:"commands" binding:"required"`
}

// APIResponse represents standard API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// ONUCommandResponse represents response for ONU operations
type ONUCommandResponse struct {
	Host       string   `json:"host"`
	Mode       string   `json:"mode"`
	Commands   []string `json:"commands"`
	Rendered   string   `json:"rendered,omitempty"`
	Output     string   `json:"output,omitempty"`
	Success    bool     `json:"success"`
	Error      string   `json:"error,omitempty"`
	Time       string   `json:"execution_time"`
	RenderOnly bool     `json:"render_only"`
}

// HealthCheckResponse represents health check response
type HealthCheckResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Services  map[string]string `json:"services"`
	Timestamp time.Time         `json:"timestamp"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}