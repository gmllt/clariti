package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gmllt/clariti/server/config"
	"github.com/gmllt/clariti/server/drivers"
	"github.com/gmllt/clariti/server/handlers"
)

func TestSetup(t *testing.T) {
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
	Setup(mux, h, cfg)

	// Test that main routes are registered
	testCases := []struct {
		method string
		path   string
		status int
	}{
		{"GET", "/health", http.StatusOK},
		{"GET", "/api", http.StatusOK},
		{"GET", "/api/v1/weather", http.StatusOK},
		{"GET", "/api/v1/platforms", http.StatusOK},
		{"GET", "/api/v1/docs", http.StatusOK},
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

func TestSetupHealthRoutes(t *testing.T) {
	mux := http.NewServeMux()
	storage := drivers.NewRAMStorage()
	cfg := &config.Config{}
	h := handlers.New(storage, cfg)

	setupHealthRoutes(mux, h)

	// Test health endpoint
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check that response contains status
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Error unmarshaling health response: %v", err)
	}

	if status, exists := response["status"]; !exists {
		t.Error("Expected 'status' field in health response")
	} else if status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%v'", status)
	}
}

func TestSetupAPIIndexRoutes(t *testing.T) {
	mux := http.NewServeMux()
	setupAPIIndexRoutes(mux)

	// Test API index endpoint
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}

	// Parse response
	var apiIndex APIIndex
	if err := json.Unmarshal(rr.Body.Bytes(), &apiIndex); err != nil {
		t.Fatalf("Error unmarshaling API index: %v", err)
	}

	// Validate structure
	if apiIndex.Service != "Clariti API" {
		t.Errorf("Expected service 'Clariti API', got '%s'", apiIndex.Service)
	}

	if len(apiIndex.Versions) == 0 {
		t.Error("Expected at least one version in API index")
	}

	// Check v1 exists
	v1Found := false
	for _, version := range apiIndex.Versions {
		if version.Version == "v1" {
			v1Found = true
			if version.Status != "stable" {
				t.Errorf("Expected v1 status 'stable', got '%s'", version.Status)
			}
			if version.DocsURL != "/api/v1/docs" {
				t.Errorf("Expected v1 docs URL '/api/v1/docs', got '%s'", version.DocsURL)
			}
		}
	}

	if !v1Found {
		t.Error("Expected to find v1 in API versions")
	}
}

func TestAPIIndex_MethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	setupAPIIndexRoutes(mux)

	// Test non-GET methods
	methods := []string{"POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/api", nil)
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

func TestRoutes_Integration(t *testing.T) {
	// Test that all routes work together
	mux := http.NewServeMux()
	storage := drivers.NewRAMStorage()
	cfg := &config.Config{
		Components: config.ComponentsConfig{
			Platforms: []config.PlatformConfig{
				{
					Name: "IntegrationTest",
					Code: "INT",
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

	Setup(mux, h, cfg)

	// Test API flow: index -> docs -> weather

	// 1. Get API index
	req, _ := http.NewRequest("GET", "/api", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	var apiIndex APIIndex
	json.Unmarshal(rr.Body.Bytes(), &apiIndex)

	if len(apiIndex.Versions) == 0 {
		t.Fatal("No versions found in API index")
	}

	// 2. Get documentation
	req, _ = http.NewRequest("GET", "/api/v1/docs", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected docs to be available, got status %d", rr.Code)
	}

	// 3. Get weather
	req, _ = http.NewRequest("GET", "/api/v1/weather", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected weather to be available, got status %d", rr.Code)
	}

	var weatherResponse map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &weatherResponse)

	// Should have our test platform
	platforms, exists := weatherResponse["platforms"]
	if !exists {
		t.Error("Expected platforms in weather response")
	}

	platformList, ok := platforms.([]interface{})
	if !ok || len(platformList) == 0 {
		t.Error("Expected at least one platform in weather response")
	}
}
