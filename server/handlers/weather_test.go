package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gmllt/clariti/server/config"
	"github.com/gmllt/clariti/server/drivers"
)

func TestNewWeatherHandler(t *testing.T) {
	cfg := &config.Config{}
	storage := drivers.NewRAMStorage()

	handler := NewWeatherHandler(cfg, storage)

	if handler == nil {
		t.Fatal("Expected weather handler to be created")
	}
	if handler.weatherService == nil {
		t.Error("Expected weather service to be initialized")
	}
}

func TestWeatherHandler_HandleWeather_GET(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Components: config.ComponentsConfig{
			Platforms: []config.PlatformConfig{
				{
					Name: "TestPlatform",
					Code: "TST",
					Instances: []config.InstanceConfig{
						{
							Name: "TestInstance",
							Code: "test",
							Components: []config.ComponentConfig{
								{Name: "TestComponent", Code: "comp"},
							},
						},
					},
				},
			},
		},
	}

	storage := drivers.NewRAMStorage()
	handler := NewWeatherHandler(cfg, storage)

	// Create request
	req, err := http.NewRequest("GET", "/api/v1/weather", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.HandleWeather(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	// Check content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}

	// Check response body structure
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	// Check required fields
	requiredFields := []string{"platforms", "instances", "components", "overall"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Expected field '%s' in response", field)
		}
	}
}

func TestWeatherHandler_HandleWeather_MethodNotAllowed(t *testing.T) {
	cfg := &config.Config{}
	storage := drivers.NewRAMStorage()
	handler := NewWeatherHandler(cfg, storage)

	// Test various non-GET methods
	methods := []string{"POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/api/v1/weather", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.HandleWeather(rr, req)

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("Expected status code %d for %s method, got %d", http.StatusMethodNotAllowed, method, status)
			}
		})
	}
}

func TestWeatherHandler_HandleWeather_WithData(t *testing.T) {
	// Setup with test data
	cfg := &config.Config{
		Components: config.ComponentsConfig{
			Platforms: []config.PlatformConfig{
				{
					Name: "TestPlatform",
					Code: "TST",
					Instances: []config.InstanceConfig{
						{
							Name: "TestInstance",
							Code: "test",
							Components: []config.ComponentConfig{
								{Name: "TestComponent", Code: "comp"},
							},
						},
					},
				},
			},
		},
	}

	storage := drivers.NewRAMStorage()
	handler := NewWeatherHandler(cfg, storage)

	req, err := http.NewRequest("GET", "/api/v1/weather", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HandleWeather(rr, req)

	// Parse response
	var response struct {
		Platforms  []interface{} `json:"platforms"`
		Instances  []interface{} `json:"instances"`
		Components []interface{} `json:"components"`
		Overall    struct {
			Status      int    `json:"status"`
			StatusLabel string `json:"status_label"`
		} `json:"overall"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}

	// Validate data
	if len(response.Platforms) != 1 {
		t.Errorf("Expected 1 platform, got %d", len(response.Platforms))
	}
	if len(response.Instances) != 1 {
		t.Errorf("Expected 1 instance, got %d", len(response.Instances))
	}
	if len(response.Components) != 1 {
		t.Errorf("Expected 1 component, got %d", len(response.Components))
	}

	// Should be operational by default
	if response.Overall.Status != 0 {
		t.Errorf("Expected operational status (0), got %d", response.Overall.Status)
	}
	if response.Overall.StatusLabel != "operational" {
		t.Errorf("Expected 'operational' status label, got '%s'", response.Overall.StatusLabel)
	}
}
