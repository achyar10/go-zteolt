package api

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all API routes using Fiber
func SetupRoutes(handlers *Handlers) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "ZTE OLT Management API",
		ServerHeader: "ZTE-OLT-API/1.0",
	})

	// Apply middleware
	app.Use(handlers.LoggingMiddleware)
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusOK)
		}

		return c.Next()
	})

	// Root path - API info
	app.Get("/", handlers.APIInfo)

	// API version 1
	v1 := app.Group("/api/v1")

	// Health check
	v1.Get("/health", handlers.HealthCheck)

	// Template management
	v1.Get("/templates", handlers.ListTemplates)

	// ONU operations
	v1.Post("/onu/add", handlers.AddONU)
	v1.Post("/onu/delete", handlers.DeleteONU)
	v1.Post("/onu/check-attenuation", handlers.CheckAttenuation)
	v1.Post("/onu/check-unconfigured", handlers.CheckUnconfigured)

	// Batch operations
	v1.Post("/batch/commands", handlers.BatchCommands)

	return app
}