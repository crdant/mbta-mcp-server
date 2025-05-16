package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

// MockServer creates a test HTTP server that returns pre-configured responses
func MockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// MockAPIResponse is a helper for creating handler functions that return mock API responses
func MockAPIResponse(statusCode int, responseBody string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}
}

// LoadFixture loads a test fixture from the test/fixtures directory
func LoadFixture(name string) ([]byte, error) {
	path := filepath.Join("../../test/fixtures", name)
	return os.ReadFile(path)
}

// LoadJSONFixture loads and unmarshals a JSON test fixture
func LoadJSONFixture(name string, v interface{}) error {
	data, err := LoadFixture(name)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// SetupTestEnv sets environment variables for testing and returns a cleanup function
func SetupTestEnv(envMap map[string]string) func() {
	originalValues := make(map[string]string)

	// Save original values and set new ones
	for key, value := range envMap {
		originalValues[key] = os.Getenv(key)
		_ = os.Setenv(key, value)
	}

	// Return cleanup function
	return func() {
		for key, value := range originalValues {
			if value == "" {
				_ = os.Unsetenv(key)
			} else {
				_ = os.Setenv(key, value)
			}
		}
	}
}