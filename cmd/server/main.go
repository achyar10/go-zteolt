package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/achyar10/go-zteolt/internal/api"
	"github.com/achyar10/go-zteolt/internal/config"
	"github.com/achyar10/go-zteolt/internal/olt"
)

func main() {
	// Parse command line flags
	var (
		host = flag.String("host", "0.0.0.0", "Server host")
		port = flag.Int("port", 8080, "Server port")
		dev  = flag.Bool("dev", false, "Development mode")
	)
	flag.Parse()

	// Load configuration
	cfg := config.DefaultConfig()
	cfg.Server.Host = *host
	cfg.Server.Port = *port

	// Initialize services
	log.Println("🚀 Initializing ZTE OLT Management API...")

	// Initialize template manager
	templateMgr, err := config.NewTemplateManager()
	if err != nil {
		log.Fatalf("❌ Failed to initialize template manager: %v", err)
	}
	log.Printf("✅ Loaded %d templates", len(templateMgr.GetAvailableTemplates()))

	// Initialize OLT service
	oltService := olt.NewService(cfg.OLT.DefaultTimeout)
	log.Printf("✅ OLT service initialized with timeout: %v", cfg.OLT.DefaultTimeout)

	// Initialize API handlers
	handlers := api.NewHandlers(oltService, templateMgr)

	// Setup Fiber routes
	app := api.SetupRoutes(handlers)

	// Configure Fiber server
	app.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))

	// Start server in goroutine
	go func() {
		log.Printf("🌐 Starting Fiber server on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if *dev {
			log.Printf("🔧 Development mode enabled")
			log.Printf("📖 API documentation: http://%s:%d", cfg.Server.Host, cfg.Server.Port)
			log.Printf("❤️  Health check: http://%s:%d/api/v1/health", cfg.Server.Host, cfg.Server.Port)
		}

		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("❌ Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown Fiber server
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("⚠️  Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
