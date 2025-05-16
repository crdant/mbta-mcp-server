package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestMockServer(t *testing.T) {
	// Create a mock server that returns a simple response
	server := MockServer(MockAPIResponse(http.StatusOK, `{"status":"ok"}`))
	defer server.Close()

	// Send a request to the mock server
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to send request to mock server: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf(`Expected status to be "ok", got "%s"`, result["status"])
	}
}

func TestSetupTestEnv(t *testing.T) {
	// Save original environment variables
	originalValue := os.Getenv("TEST_ENV_VAR")

	// Set up test environment
	cleanup := SetupTestEnv(map[string]string{
		"TEST_ENV_VAR": "test_value",
	})
	defer cleanup()

	// Check if environment variable was set
	if value := os.Getenv("TEST_ENV_VAR"); value != "test_value" {
		t.Errorf(`Expected TEST_ENV_VAR to be "test_value", got "%s"`, value)
	}

	// Run cleanup
	cleanup()

	// Check if environment variable was restored
	if value := os.Getenv("TEST_ENV_VAR"); value != originalValue {
		t.Errorf(`Expected TEST_ENV_VAR to be restored to "%s", got "%s"`, originalValue, value)
	}
}