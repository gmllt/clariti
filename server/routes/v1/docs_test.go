package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupV1DocumentationRoutes(t *testing.T) {
	mux := http.NewServeMux()
	SetupV1DocumentationRoutes(mux)

	// Test documentation endpoint
	req, err := http.NewRequest("GET", "/api/v1/docs", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check status
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}

	// Parse response
	var docs V1APIDoc
	if err := json.Unmarshal(rr.Body.Bytes(), &docs); err != nil {
		t.Fatalf("Error unmarshaling docs: %v", err)
	}

	// Validate structure
	if docs.Service != "Clariti API" {
		t.Errorf("Expected service 'Clariti API', got '%s'", docs.Service)
	}
	if docs.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", docs.Version)
	}
	if docs.BaseURL != "/api/v1" {
		t.Errorf("Expected base URL '/api/v1', got '%s'", docs.BaseURL)
	}

	// Check endpoints exist
	expectedEndpointGroups := []string{"components", "incidents", "planned-maintenances", "weather"}
	for _, group := range expectedEndpointGroups {
		if _, exists := docs.Endpoints[group]; !exists {
			t.Errorf("Expected endpoint group '%s' in documentation", group)
		}
	}
}

func TestV1DocumentationRoute_MethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	SetupV1DocumentationRoutes(mux)

	// Test non-GET methods
	methods := []string{"POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/api/v1/docs", nil)
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

func TestV1APIDoc_Structure(t *testing.T) {
	// Test the documentation structure matches expected format
	mux := http.NewServeMux()
	SetupV1DocumentationRoutes(mux)

	req, err := http.NewRequest("GET", "/api/v1/docs", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	var docs V1APIDoc
	if err := json.Unmarshal(rr.Body.Bytes(), &docs); err != nil {
		t.Fatalf("Error unmarshaling docs: %v", err)
	}

	// Check weather endpoint specifically
	weatherEndpoints, exists := docs.Endpoints["weather"]
	if !exists {
		t.Fatal("Expected 'weather' endpoint group")
	}

	if len(weatherEndpoints) != 1 {
		t.Errorf("Expected 1 weather endpoint, got %d", len(weatherEndpoints))
	}

	weatherEndpoint := weatherEndpoints[0]
	if weatherEndpoint.Path != "/api/v1/weather" {
		t.Errorf("Expected weather path '/api/v1/weather', got '%s'", weatherEndpoint.Path)
	}

	if len(weatherEndpoint.Methods) != 1 || weatherEndpoint.Methods[0] != "GET" {
		t.Errorf("Expected weather endpoint to only support GET, got %v", weatherEndpoint.Methods)
	}

	if weatherEndpoint.AuthRequired {
		t.Error("Expected weather endpoint to not require auth")
	}

	if weatherEndpoint.Description == "" {
		t.Error("Expected weather endpoint to have a description")
	}
}

func TestV1Route_Structure(t *testing.T) {
	// Test V1Route struct
	route := V1Route{
		Path:         "/api/v1/test",
		Methods:      []string{"GET", "POST"},
		Description:  "Test endpoint",
		AuthRequired: true,
	}

	if route.Path != "/api/v1/test" {
		t.Errorf("Expected path '/api/v1/test', got '%s'", route.Path)
	}

	if len(route.Methods) != 2 {
		t.Errorf("Expected 2 methods, got %d", len(route.Methods))
	}

	expectedMethods := []string{"GET", "POST"}
	for i, method := range expectedMethods {
		if route.Methods[i] != method {
			t.Errorf("Expected method '%s' at index %d, got '%s'", method, i, route.Methods[i])
		}
	}

	if !route.AuthRequired {
		t.Error("Expected auth to be required")
	}

	if route.Description != "Test endpoint" {
		t.Errorf("Expected description 'Test endpoint', got '%s'", route.Description)
	}
}

func TestV1Documentation_AllEndpoints(t *testing.T) {
	mux := http.NewServeMux()
	SetupV1DocumentationRoutes(mux)

	req, err := http.NewRequest("GET", "/api/v1/docs", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	var docs V1APIDoc
	json.Unmarshal(rr.Body.Bytes(), &docs)

	// Count total endpoints
	totalEndpoints := 0
	for _, routes := range docs.Endpoints {
		totalEndpoints += len(routes)
	}

	// We expect at least 10 endpoints (5 component + 2 incident + 2 maintenance + 1 weather)
	if totalEndpoints < 10 {
		t.Errorf("Expected at least 10 endpoints, got %d", totalEndpoints)
	}

	// Check specific endpoint paths exist
	expectedPaths := []string{
		"/api/v1/components",
		"/api/v1/platforms",
		"/api/v1/incidents",
		"/api/v1/planned-maintenances",
		"/api/v1/weather",
	}

	foundPaths := make(map[string]bool)
	for _, routes := range docs.Endpoints {
		for _, route := range routes {
			foundPaths[route.Path] = true
		}
	}

	for _, expectedPath := range expectedPaths {
		if !foundPaths[expectedPath] {
			t.Errorf("Expected to find path '%s' in documentation", expectedPath)
		}
	}
}
