package config

import (
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// Save current environment to restore later
	originalAPIKey := os.Getenv("MBTA_API_KEY")
	originalDebug := os.Getenv("DEBUG")
	originalLogLevel := os.Getenv("LOG_LEVEL")
	originalTimeout := os.Getenv("TIMEOUT_SECONDS")
	originalAPIURL := os.Getenv("MBTA_API_URL")
	originalEnv := os.Getenv("ENVIRONMENT")

	// Ensure we restore environment after test
	defer func() {
		os.Setenv("MBTA_API_KEY", originalAPIKey)
		os.Setenv("DEBUG", originalDebug)
		os.Setenv("LOG_LEVEL", originalLogLevel)
		os.Setenv("TIMEOUT_SECONDS", originalTimeout)
		os.Setenv("MBTA_API_URL", originalAPIURL)
		os.Setenv("ENVIRONMENT", originalEnv)
	}()

	// Test with default values
	t.Run("Default values", func(t *testing.T) {
		os.Unsetenv("MBTA_API_KEY")
		os.Unsetenv("DEBUG")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("TIMEOUT_SECONDS")
		os.Unsetenv("MBTA_API_URL")
		os.Unsetenv("ENVIRONMENT")

		config := New()

		if config.APIKey != "" {
			t.Errorf("Expected empty APIKey, got %s", config.APIKey)
		}
		if config.Debug != false {
			t.Errorf("Expected Debug to be false, got %v", config.Debug)
		}
		if config.LogLevel != "info" {
			t.Errorf("Expected LogLevel to be info, got %s", config.LogLevel)
		}
		if config.Timeout != 30*time.Second {
			t.Errorf("Expected Timeout to be 30s, got %v", config.Timeout)
		}
		if config.APIBaseURL != "https://api-v3.mbta.com" {
			t.Errorf("Expected APIBaseURL to be https://api-v3.mbta.com, got %s", config.APIBaseURL)
		}
		if config.Environment != "development" {
			t.Errorf("Expected Environment to be development, got %s", config.Environment)
		}
	})

	// Test with custom values
	t.Run("Custom values", func(t *testing.T) {
		os.Setenv("MBTA_API_KEY", "test-api-key")
		os.Setenv("DEBUG", "true")
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("TIMEOUT_SECONDS", "60")
		os.Setenv("MBTA_API_URL", "http://test-api.example.com")
		os.Setenv("ENVIRONMENT", "test")

		config := New()

		if config.APIKey != "test-api-key" {
			t.Errorf("Expected APIKey to be test-api-key, got %s", config.APIKey)
		}
		if config.Debug != true {
			t.Errorf("Expected Debug to be true, got %v", config.Debug)
		}
		if config.LogLevel != "debug" {
			t.Errorf("Expected LogLevel to be debug, got %s", config.LogLevel)
		}
		if config.Timeout != 60*time.Second {
			t.Errorf("Expected Timeout to be 60s, got %v", config.Timeout)
		}
		if config.APIBaseURL != "http://test-api.example.com" {
			t.Errorf("Expected APIBaseURL to be http://test-api.example.com, got %s", config.APIBaseURL)
		}
		if config.Environment != "test" {
			t.Errorf("Expected Environment to be test, got %s", config.Environment)
		}
	})

	// Test with invalid boolean value
	t.Run("Invalid boolean", func(t *testing.T) {
		os.Setenv("DEBUG", "not-a-bool")

		config := New()

		if config.Debug != false {
			t.Errorf("Expected Debug to be false for invalid input, got %v", config.Debug)
		}
	})

	// Test with invalid integer value
	t.Run("Invalid integer", func(t *testing.T) {
		os.Setenv("TIMEOUT_SECONDS", "not-an-int")

		config := New()

		if config.Timeout != 30*time.Second {
			t.Errorf("Expected Timeout to be 30s for invalid input, got %v", config.Timeout)
		}
	})
}