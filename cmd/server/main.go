package main

import (
	"fmt"
	"log"
	"os"

	"github.com/crdant/mbta-mcp-server/internal/config"
	"github.com/crdant/mbta-mcp-server/internal/server"
)

// Version will be set at build time
// Format: semver+build.COMMIT_SHA (e.g., 1.2.3+build.a1b2c3d)
var Version = "0.1.0+build.dev"

func main() {
	log.Printf("Starting MBTA MCP Server version %s", Version)

	// Load configuration from environment variables
	cfg := config.New()

	// Check for required API key
	if cfg.APIKey == "" {
		log.Println("Warning: No MBTA_API_KEY environment variable found. API functionality will be limited.")
	}

	// Configure logging
	if cfg.Debug {
		log.Println("Debug mode enabled")
	}
	log.Printf("Log level set to: %s", cfg.LogLevel)
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("MBTA API URL: %s", cfg.APIBaseURL)
	log.Printf("Request timeout: %v", cfg.Timeout)

	// Initialize MCP server
	mcpServer, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize MCP server: %v", err)
		os.Exit(1)
	}

	// Register default handlers for transit information
	mcpServer.RegisterDefaultHandlers()

	// Start the server with stdio protocol
	fmt.Println("MBTA MCP Server started successfully")
	fmt.Println("Using stdio protocol - input and output are on stdin/stdout")
	if err := mcpServer.Start(); err != nil {
		log.Fatalf("Error running MCP server: %v", err)
		os.Exit(1)
	}
}