package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gmllt/clariti/models/component"
	"github.com/gmllt/clariti/models/event"
	"github.com/gmllt/clariti/server/config"
	"github.com/gmllt/clariti/server/core"
	"github.com/gmllt/clariti/server/drivers"
)

func TestServer_Integration(t *testing.T) {
	// Create test config
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: "0", // Let OS choose port
		},
		Auth: config.AuthConfig{
			AdminUsername: "admin",
			AdminPassword: "test123",
		},
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
								{Name: "API", Code: "api"},
								{Name: "Database", Code: "db"},
							},
						},
					},
				},
			},
		},
	}

	// Create server
	storage := drivers.NewRAMStorage()
	server := core.NewWithConfig(cfg, storage)

	// Create test server
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	// Test API flow
	t.Run("API Index", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var apiIndex map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&apiIndex)

		if service, exists := apiIndex["service"]; !exists || service != "Clariti API" {
			t.Error("Expected service to be 'Clariti API'")
		}
	})

	t.Run("Weather Before Incidents", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/weather")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		var weather map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&weather)

		overall := weather["overall"].(map[string]interface{})
		if status := overall["status"].(float64); status != 0 {
			t.Errorf("Expected operational status (0), got %v", status)
		}
	})

	t.Run("Create Incident", func(t *testing.T) {
		incidentData := map[string]interface{}{
			"title":   "API Degradation",
			"content": "API response times are slow",
			"components": []map[string]interface{}{
				{
					"name": "API",
					"code": "api",
					"instance": map[string]interface{}{
						"name": "TestInstance",
						"code": "test",
						"platform": map[string]interface{}{
							"name": "TestPlatform",
							"code": "TST",
						},
					},
				},
			},
			"criticality": event.CriticalityDegraded,
			"perpetual":   false,
		}

		jsonData, _ := json.Marshal(incidentData)

		req, _ := http.NewRequest("POST", ts.URL+"/api/v1/incidents", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth("admin", "test123")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}

		var incident map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&incident)

		if guid, exists := incident["guid"]; !exists || guid == "" {
			t.Error("Expected incident to have GUID")
		}
	})

	t.Run("Weather After Incident", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/weather")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		var weather map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&weather)

		overall := weather["overall"].(map[string]interface{})
		if status := overall["status"].(float64); status != float64(event.CriticalityDegraded) {
			t.Errorf("Expected degraded status (%d), got %v", event.CriticalityDegraded, status)
		}

		// Check active events
		activeEvents := overall["active_events"].([]interface{})
		if len(activeEvents) != 1 {
			t.Errorf("Expected 1 active event, got %d", len(activeEvents))
		}
	})

	t.Run("Create Planned Maintenance", func(t *testing.T) {
		now := time.Now()
		maintenanceData := map[string]interface{}{
			"title":   "Database Upgrade",
			"content": "Upgrading database to latest version",
			"components": []map[string]interface{}{
				{
					"name": "Database",
					"code": "db",
					"instance": map[string]interface{}{
						"name": "TestInstance",
						"code": "test",
						"platform": map[string]interface{}{
							"name": "TestPlatform",
							"code": "TST",
						},
					},
				},
			},
			"start_planned": now.Format(time.RFC3339),
			"end_planned":   now.Add(2 * time.Hour).Format(time.RFC3339),
		}

		jsonData, _ := json.Marshal(maintenanceData)

		req, _ := http.NewRequest("POST", ts.URL+"/api/v1/planned-maintenances", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth("admin", "test123")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})

	t.Run("Weather With Multiple Events", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/weather")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		var weather map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&weather)

		overall := weather["overall"].(map[string]interface{})

		// Should show highest criticality (maintenance = 4 > degraded = 1)
		if status := overall["status"].(float64); status != float64(event.CriticalityUnderMaintenance) {
			t.Errorf("Expected under maintenance status (%d), got %v", event.CriticalityUnderMaintenance, status)
		}

		// Should have both events
		activeEvents := overall["active_events"].([]interface{})
		if len(activeEvents) != 2 {
			t.Errorf("Expected 2 active events, got %d", len(activeEvents))
		}
	})
}

func TestServer_ErrorHandling(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Host: "localhost", Port: "0"},
		Auth:   config.AuthConfig{AdminUsername: "admin", AdminPassword: "test123"},
		Components: config.ComponentsConfig{
			Platforms: []config.PlatformConfig{
				{Name: "Test", Code: "TST", Instances: []config.InstanceConfig{}},
			},
		},
	}

	storage := drivers.NewRAMStorage()
	server := core.NewWithConfig(cfg, storage)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	t.Run("Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", ts.URL+"/api/v1/incidents", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth("admin", "test123")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Unauthorized Request", func(t *testing.T) {
		req, _ := http.NewRequest("POST", ts.URL+"/api/v1/incidents", bytes.NewBufferString("{}"))
		req.Header.Set("Content-Type", "application/json")
		// No auth header

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", ts.URL+"/api/v1/weather", nil)
		req.SetBasicAuth("admin", "test123") // Add auth for middleware

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", resp.StatusCode)
		}
	})
}

func TestServer_Documentation(t *testing.T) {
	cfg := &config.Config{
		Components: config.ComponentsConfig{
			Platforms: []config.PlatformConfig{
				{Name: "Test", Code: "TST", Instances: []config.InstanceConfig{}},
			},
		},
	}

	storage := drivers.NewRAMStorage()
	server := core.NewWithConfig(cfg, storage)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	t.Run("API Documentation", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/docs")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var docs map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&docs)

		if endpoints, exists := docs["endpoints"]; !exists {
			t.Error("Expected endpoints in documentation")
		} else {
			endpointMap := endpoints.(map[string]interface{})

			// Check weather endpoint is documented
			if _, exists := endpointMap["weather"]; !exists {
				t.Error("Expected weather endpoints in documentation")
			}
		}
	})
}

func BenchmarkServer_Weather(b *testing.B) {
	cfg := &config.Config{
		Components: config.ComponentsConfig{
			Platforms: []config.PlatformConfig{
				{
					Name: "BenchPlatform",
					Code: "BNC",
					Instances: []config.InstanceConfig{
						{
							Name: "BenchInstance",
							Code: "bench",
							Components: []config.ComponentConfig{
								{Name: "API", Code: "api"},
								{Name: "DB", Code: "db"},
								{Name: "Cache", Code: "cache"},
							},
						},
					},
				},
			},
		},
	}

	storage := drivers.NewRAMStorage()
	server := core.NewWithConfig(cfg, storage)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	// Add some test incidents for more realistic benchmark
	platform := component.NewPlatform("BenchPlatform", "BNC")
	instance := component.NewInstance("BenchInstance", "bench", platform)
	comp := component.NewComponent("API", "api", instance)

	for i := 0; i < 10; i++ {
		incident := event.NewFiringIncident(
			fmt.Sprintf("Test Incident %d", i),
			"Benchmark incident",
			[]*component.Component{comp},
			event.CriticalityDegraded,
		)
		storage.CreateIncident(incident)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := http.Get(ts.URL + "/api/v1/weather")
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}
