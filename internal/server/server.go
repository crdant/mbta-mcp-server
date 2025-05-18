// ABOUTME: This file implements the MCP server for the MBTA application.
// ABOUTME: It provides server initialization and configuration handling.

package server

import (
	"fmt"
	"log"
	"os"

	"github.com/crdant/mbta-mcp-server/internal/config"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// Server represents the MBTA MCP server.
type Server struct {
	mcpServer *mcpserver.MCPServer
	config    *config.Config
}

// New creates a new MBTA MCP server with the provided configuration.
// It initializes the MCP server and sets up basic configuration.
func New(cfg *config.Config) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	// Create options for MCP server
	var serverOpts []mcpserver.ServerOption

	// Enable logging if debug is enabled
	if cfg.Debug {
		serverOpts = append(serverOpts, mcpserver.WithLogging())
	}

	// Create MCP server with name and version
	mcpServer := mcpserver.NewMCPServer(
		"MBTA MCP Server",
		"0.1.0", // This should eventually be dynamic
		serverOpts...,
	)

	// Create and return server
	server := &Server{
		mcpServer: mcpServer,
		config:    cfg,
	}

	return server, nil
}

// SetMetadata sets additional metadata for the MCP server.
// This is a placeholder as the current API doesn't directly support metadata.
func (s *Server) SetMetadata(key string, value interface{}) {
	// Log the attempt for now
	log.Printf("Setting metadata %s = %v (note: not supported in current MCP API)", key, value)
}

// Start starts the MCP server using stdio protocol.
// It blocks until the server is stopped or encounters an error.
func (s *Server) Start() error {
	log.Println("Starting MBTA MCP Server with stdio protocol")

	// Apply middleware before starting
	s.ApplyMiddleware()

	// Configure stdio options
	var stdioOpts []mcpserver.StdioOption

	// Add custom error logger
	errorLogger := log.New(os.Stderr, "MBTA-MCP-ERROR: ", log.LstdFlags)
	stdioOpts = append(stdioOpts, mcpserver.WithErrorLogger(errorLogger))

	// Start the server with stdio transport and options
	return mcpserver.ServeStdio(s.mcpServer, stdioOpts...)
}

// This implementation is in middleware.go

// This implementation is in handlers.go
