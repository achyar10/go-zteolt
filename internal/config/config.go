package config

import (
	"time"
)

// Config holds application configuration
type Config struct {
	Server struct {
		Host         string        `json:"host"`
		Port         int           `json:"port"`
		ReadTimeout  time.Duration `json:"read_timeout"`
		WriteTimeout time.Duration `json:"write_timeout"`
		IdleTimeout  time.Duration `json:"idle_timeout"`
	} `json:"server"`

	OLT struct {
		DefaultTimeout  time.Duration `json:"default_timeout"`
		WriteTimeout    time.Duration `json:"write_timeout"`
		MaxRetries      int           `json:"max_retries"`
		ParallelWorkers int           `json:"parallel_workers"`
	} `json:"olt"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	cfg := &Config{}

	// Server defaults
	cfg.Server.Host = "0.0.0.0"
	cfg.Server.Port = 8080
	cfg.Server.ReadTimeout = 30 * time.Second
	cfg.Server.WriteTimeout = 30 * time.Second
	cfg.Server.IdleTimeout = 60 * time.Second

	// OLT defaults
	cfg.OLT.DefaultTimeout = 8 * time.Second
	cfg.OLT.WriteTimeout = 24 * time.Second
	cfg.OLT.MaxRetries = 2
	cfg.OLT.ParallelWorkers = 8

	return cfg
}