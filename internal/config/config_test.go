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
		_ = os.Setenv("MBTA_API_KEY", originalAPIKey)
		_ = os.Setenv("DEBUG", originalDebug)
		_ = os.Setenv("LOG_LEVEL", originalLogLevel)
		_ = os.Setenv("TIMEOUT_SECONDS", originalTimeout)
		_ = os.Setenv("MBTA_API_URL", originalAPIURL)
		_ = os.Setenv("ENVIRONMENT", originalEnv)
	}()

	// Test with default values
	t.Run("Default values", func(t *testing.T) {
		_ = os.Unsetenv("MBTA_API_KEY")
		_ = os.Unsetenv("DEBUG")
		_ = os.Unsetenv("LOG_LEVEL")
		_ = os.Unsetenv("TIMEOUT_SECONDS")
		_ = os.Unsetenv("MBTA_API_URL")
		_ = os.Unsetenv("ENVIRONMENT")

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
		_ = os.Setenv("MBTA_API_KEY", "test-api-key")
		_ = os.Setenv("DEBUG", "true")
		_ = os.Setenv("LOG_LEVEL", "debug")
		_ = os.Setenv("TIMEOUT_SECONDS", "60")
		_ = os.Setenv("MBTA_API_URL", "http://test-api.example.com")
		_ = os.Setenv("ENVIRONMENT", "test")

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
		_ = os.Setenv("DEBUG", "not-a-bool")

		config := New()

		if config.Debug != false {
			t.Errorf("Expected Debug to be false for invalid input, got %v", config.Debug)
		}
	})

	// Test with invalid integer value
	t.Run("Invalid integer", func(t *testing.T) {
		_ = os.Setenv("TIMEOUT_SECONDS", "not-an-int")

		config := New()

		if config.Timeout != 30*time.Second {
			t.Errorf("Expected Timeout to be 30s for invalid input, got %v", config.Timeout)
		}
	})
}