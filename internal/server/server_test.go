// ABOUTME: This file contains tests for the MCP server initialization.
// ABOUTME: It verifies server creation, configuration, and error handling.

package server

import (
	"testing"

	"github.com/crdant/mbta-mcp-server/internal/config"
)

func TestNewServer(t *testing.T) {
	// Test with valid configuration
	t.Run("Valid configuration", func(t *testing.T) {
		cfg := &config.Config{
			APIKey:      "test-api-key",
			Debug:       true,
			LogLevel:    "debug",
			Timeout:     30,
			APIBaseURL:  "https://api-test.mbta.com",
			Environment: "test",
		}

		server, err := New(cfg)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if server == nil {
			t.Error("Expected server instance, got nil")
		}
	})

	// Test with nil configuration
	t.Run("Nil configuration", func(t *testing.T) {
		server, err := New(nil)
		if err == nil {
			t.Error("Expected error for nil config, got nil")
		}
		if server != nil {
			t.Errorf("Expected nil server for nil config, got %v", server)
		}
	})
}

func TestServerMetadata(t *testing.T) {
	// Create a test configuration
	cfg := &config.Config{
		APIKey:      "test-api-key",
		Debug:       false,
		LogLevel:    "info",
		Timeout:     30,
		APIBaseURL:  "https://api-test.mbta.com",
		Environment: "test",
	}

	// Create a server instance
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test SetMetadata functionality
	t.Run("Set metadata on server", func(t *testing.T) {
		// This is just testing the interface exists
		// The actual functionality will be tested after implementation
		server.SetMetadata("test_key", "test_value")
		// If we reach here without panicking, the method exists
	})
}

func TestStartServer(t *testing.T) {
	// Create a test configuration
	cfg := &config.Config{
		APIKey:      "test-api-key",
		Debug:       false,
		LogLevel:    "info",
		Timeout:     30,
		APIBaseURL:  "https://api-test.mbta.com",
		Environment: "test",
	}

	// Create a server instance
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test Start method exists
	t.Run("Start method exists", func(t *testing.T) {
		// Don't actually start the server as it would block the test
		// Just check the method exists on the interface
		startMethod := func() {
			_ = server.Start()
		}

		// If we can compile this, the method exists
		// This is just a compile-time check
		_ = startMethod
	})
}
