package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	APIKey      string
	Debug       bool
	LogLevel    string
	Timeout     time.Duration
	APIBaseURL  string
	Environment string
}

// New creates a new configuration from environment variables
func New() *Config {
	return &Config{
		APIKey:      getEnv("MBTA_API_KEY", ""),
		Debug:       getEnvBool("DEBUG", false),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Timeout:     time.Duration(getEnvInt("TIMEOUT_SECONDS", 30)) * time.Second,
		APIBaseURL:  getEnv("MBTA_API_URL", "https://api-v3.mbta.com"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

// getEnv retrieves an environment variable with a fallback value
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

// getEnvBool retrieves a boolean environment variable with a fallback
func getEnvBool(key string, fallback bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return boolValue
}

// getEnvInt retrieves an integer environment variable with a fallback
func getEnvInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}