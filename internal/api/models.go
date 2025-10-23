package api

import "time"

// Request/Response models for REST API

// AddONURequest represents request to add ONU
type AddONURequest struct {
	Host           string `json:"host" binding:"required"`
	Port           int    `json:"port" binding:"required"`
	User           string `json:"user" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Board          int    `json:"board" binding:"required"`
	PON            int    `json:"pon" binding:"required"`
	ONU            int    `json:"onu" binding:"required"`
	SerialNumber   string `json:"serial_number" binding:"required"`
	Name           string `json:"name" binding:"required"`
	SecretPassword string `json:"secret_password" binding:"required"`
	Description    string `json:"description" binding:"required"`
	VlanID         int    `json:"vlan_id" binding:"required"`
	TcontProfile   string `json:"tcont_profile" binding:"required"`
	TrafficLimit   string `json:"traffic_limit" binding:"required"`
	RenderOnly     bool   `json:"render_only"`
}

// DeleteONURequest represents request to delete ONU
type DeleteONURequest struct {
	Host       string `json:"host" binding:"required"`
	Port       int    `json:"port" binding:"required"`
	User       string `json:"user" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Board      int    `json:"board" binding:"required"`
	PON        int    `json:"pon" binding:"required"`
	ONU        int    `json:"onu" binding:"required"`
	RenderOnly bool   `json:"render_only"`
}

// CheckAttenuationRequest represents request to check attenuation
type CheckAttenuationRequest struct {
	Host       string `json:"host" binding:"required"`
	Port       int    `json:"port" binding:"required"`
	User       string `json:"user" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Board      int    `json:"board" binding:"required"`
	PON        int    `json:"pon" binding:"required"`
	ONU        int    `json:"onu" binding:"required"`
	RenderOnly bool   `json:"render_only"`
}

// AttenuationData represents parsed attenuation information
type AttenuationData struct {
	Host        string  `json:"host"`
	Board       int     `json:"board"`
	PON         int     `json:"pon"`
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
	OLTIndex     string `json:"olt_index"`
	Model        string `json:"model"`
	SerialNumber string `json:"serial_number"`
	Board        int    `json:"board"`
	PON          int    `json:"pon"`
}

// UnconfiguredONUListResponse represents response for check unconfigured ONUs
type UnconfiguredONUListResponse struct {
	Host       string                  `json:"host"`
	Mode       string                  `json:"mode"`
	Success    bool                    `json:"success"`
	Error      string                  `json:"error,omitempty"`
	Time       string                  `json:"execution_time"`
	RenderOnly bool                    `json:"render_only"`
	Data       *UnconfiguredONUListDTO `json:"data,omitempty"`
}

// UnconfiguredONUListDTO represents the data transfer object for unconfigured ONUs
type UnconfiguredONUListDTO struct {
	Host           string                          `json:"host"`
	TotalCount     int                             `json:"total_count"`
	ONUs           []UnconfiguredONUDTO            `json:"onus"`
	GroupedByBoard map[string][]UnconfiguredONUDTO `json:"grouped_by_board,omitempty"`
	Status         string                          `json:"status"`
	RawOutput      string                          `json:"raw_output,omitempty"`
}

// AttenuationDataDTO represents the data transfer object for attenuation
type AttenuationDataDTO struct {
	Host        string  `json:"host"`
	Board       int     `json:"board"`
	PON         int     `json:"pon"`
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

// RebootONURequest represents request to reboot ONU
type RebootONURequest struct {
	Host       string `json:"host" binding:"required"`
	Port       int    `json:"port" binding:"required"`
	User       string `json:"user" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Board      int    `json:"board" binding:"required"`
	PON        int    `json:"pon" binding:"required"`
	ONU        int    `json:"onu" binding:"required"`
	RenderOnly bool   `json:"render_only"`
}

// RebootONUResponse represents response for ONU reboot
type RebootONUResponse struct {
	Host       string `json:"host"`
	Mode       string `json:"mode"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	Time       string `json:"execution_time"`
	RenderOnly bool   `json:"render_only"`
}

// SaveConfigurationRequest represents request to save configuration
type SaveConfigurationRequest struct {
	Host       string `json:"host" binding:"required"`
	Port       int    `json:"port" binding:"required"`
	User       string `json:"user" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Timeout    int    `json:"timeout,omitempty"` // custom timeout in seconds (default: 300s for save operations)
	RenderOnly bool   `json:"render_only"`
}

// SaveConfigurationResponse represents response for save configuration
type SaveConfigurationResponse struct {
	Host       string `json:"host"`
	Mode       string `json:"mode"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	Time       string `json:"execution_time"`
	RenderOnly bool   `json:"render_only"`
	Output     string `json:"output,omitempty"`
	Status     string `json:"status,omitempty"` // success, in_progress, failed, timeout
	TimeoutUsed int   `json:"timeout_used,omitempty"` // timeout used in seconds (for debugging)
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
	Data      any        `json:"data,omitempty"`
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

// SNMP Monitoring Request/Response models

// SNMPMonitoringRequest represents SNMP monitoring request
type SNMPMonitoringRequest struct {
	Host      string `json:"host" binding:"required"`
	Port      int    `json:"port" binding:"required"`
	Community string `json:"community" binding:"required"`
	Timeout   int    `json:"timeout,omitempty"` // optional timeout in seconds
}

// SNMPONUInfo represents ONU information from SNMP
type SNMPONUInfo struct {
	Board        int    `json:"board"`
	PON          int    `json:"pon"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	SerialNumber string `json:"serial_number"`
	OnuType      string `json:"onu_type"`
	RXPower      string `json:"rx_power"`
	TXPower      string `json:"tx_power"`
	Status       string `json:"status"`
	IPAddress    string `json:"ip_address,omitempty"`
	LastOnline   string `json:"last_online,omitempty"`
	Uptime       string `json:"uptime,omitempty"`
}

// SNMPMonitoringResponse represents SNMP monitoring response
type SNMPMonitoringResponse struct {
	Host          string        `json:"host"`
	BoardID       int           `json:"board_id"`
	PONID         int           `json:"pon_id"`
	TotalONUs     int           `json:"total_onus"`
	ONUs          []SNMPONUInfo `json:"onus"`
	ExecutionTime string        `json:"execution_time"`
	Timestamp     time.Time     `json:"timestamp"`
}

// SNMPONUDetailsRequest represents request for specific ONU details
type SNMPONUDetailsRequest struct {
	Host      string `json:"host" binding:"required"`
	Port      int    `json:"port" binding:"required"`
	Community string `json:"community" binding:"required"`
	Timeout   int    `json:"timeout,omitempty"`
}

// SNMPEmptySlotsRequest represents request for empty ONU boards
type SNMPEmptySlotsRequest struct {
	Host      string `json:"host" binding:"required"`
	Port      int    `json:"port" binding:"required"`
	Community string `json:"community" binding:"required"`
	Timeout   int    `json:"timeout,omitempty"`
}

// SNMPEmptySlot represents an empty ONU board
type SNMPEmptySlot struct {
	Board int `json:"board"`
	PON   int `json:"pon"`
	ONUID int `json:"onu_id"`
}

// SNMPEmptySlotsResponse represents response for empty ONU boards
type SNMPEmptySlotsResponse struct {
	Host          string          `json:"host"`
	BoardID       int             `json:"board_id"`
	PONID         int             `json:"pon_id"`
	TotalEmpty    int             `json:"total_empty"`
	EmptySlots    []SNMPEmptySlot `json:"empty_boards"`
	ExecutionTime string          `json:"execution_time"`
	Timestamp     time.Time       `json:"timestamp"`
}
