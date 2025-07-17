package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gmllt/clariti/server/config"
	"github.com/gmllt/clariti/server/drivers"
	"github.com/gmllt/clariti/server/handlers"
)

func TestSetupV1Routes(t *testing.T) {
	// Setup
	mux := http.NewServeMux()
	storage := drivers.NewRAMStorage()
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
	h := handlers.New(storage, cfg)

	// Setup routes
	SetupV1Routes(mux, h)

	// Test that routes are registered by making requests
	testCases := []struct {
		method string
		path   string
		status int
	}{
		{"GET", "/api/v1/platforms", http.StatusOK},
		{"GET", "/api/v1/instances", http.StatusOK},
		{"GET", "/api/v1/components", http.StatusOK},
		{"GET", "/api/v1/components/hierarchy", http.StatusOK},
		{"GET", "/api/v1/components/list", http.StatusOK},
		{"GET", "/api/v1/incidents", http.StatusOK},
		{"GET", "/api/v1/planned-maintenances", http.StatusOK},
		{"GET", "/api/v1/weather", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != tc.status {
				t.Errorf("Expected status %d for %s %s, got %d", tc.status, tc.method, tc.path, rr.Code)
			}
		})
	}
}

func TestSetupV1WeatherRoutes(t *testing.T) {
	// Setup
	mux := http.NewServeMux()
	storage := drivers.NewRAMStorage()
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
	h := handlers.New(storage, cfg)

	// Setup only weather routes
	setupV1WeatherRoutes(mux, h)

	// Test weather endpoint
	req, err := http.NewRequest("GET", "/api/v1/weather", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check response structure
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	requiredFields := []string{"platforms", "instances", "components", "overall"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Expected field '%s' in weather response", field)
		}
	}
}

func TestV1WeatherRoute_MethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	storage := drivers.NewRAMStorage()
	cfg := &config.Config{}
	h := handlers.New(storage, cfg)

	setupV1WeatherRoutes(mux, h)

	// Test non-GET methods
	methods := []string{"POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/api/v1/weather", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d for %s method, got %d", http.StatusMethodNotAllowed, method, rr.Code)
			}
		})
	}
}

func TestSetupV1ComponentRoutes(t *testing.T) {
	mux := http.NewServeMux()
	storage := drivers.NewRAMStorage()
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
	h := handlers.New(storage, cfg)

	setupV1ComponentRoutes(mux, h)

	// Test component endpoints
	componentEndpoints := []string{
		"/api/v1/components",
		"/api/v1/components/hierarchy",
		"/api/v1/platforms",
		"/api/v1/instances",
		"/api/v1/components/list",
	}

	for _, endpoint := range componentEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Expected status %d for %s, got %d", http.StatusOK, endpoint, rr.Code)
			}

			// Check that response is valid JSON
			var response interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("Invalid JSON response for %s: %v", endpoint, err)
			}
		})
	}
}

func TestSetupV1EventRoutes(t *testing.T) {
	mux := http.NewServeMux()
	storage := drivers.NewRAMStorage()
	cfg := &config.Config{}
	h := handlers.New(storage, cfg)

	setupV1EventRoutes(mux, h)

	// Test event endpoints (GET should work)
	eventEndpoints := []string{
		"/api/v1/incidents",
		"/api/v1/planned-maintenances",
	}

	for _, endpoint := range eventEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Expected status %d for %s, got %d", http.StatusOK, endpoint, rr.Code)
			}

			// Check that response is valid JSON
			var response interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("Invalid JSON response for %s: %v", endpoint, err)
			}
		})
	}
}
