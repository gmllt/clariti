package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gmllt/clariti/models/component"
	"github.com/gmllt/clariti/models/event"
	"github.com/gmllt/clariti/server/config"
	"github.com/gmllt/clariti/server/drivers"
)

// BenchmarkAPIHandler_HandleHealth tests the performance of the health endpoint
func BenchmarkAPIHandler_HandleHealth(b *testing.B) {
	storage := drivers.NewRAMStorage()
	cfg := &config.Config{}
	handler := NewAPIHandler(storage, cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		rr := httptest.NewRecorder()
		handler.HandleHealth(rr, req)
	}
}

// BenchmarkWeatherHandler_HandleWeather tests the performance of the weather endpoint
func BenchmarkWeatherHandler_HandleWeather(b *testing.B) {
	storage := drivers.NewRAMStorage()
	cfg := &config.Config{}
	handler := NewWeatherHandler(cfg, storage)

	// Setup test data
	platform := component.NewPlatform("AWS Production", "aws-prod")
	instance := component.NewInstance("EKS Cluster", "eks-cluster", platform)
	comp := component.NewComponent("API Gateway", "api-gateway", instance)

	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	incident := &event.Incident{
		BaseEvent: event.BaseEvent{
			GUID:           "incident-1",
			Title:          "API Issues",
			Content:        "API experiencing issues",
			Components:     []*component.Component{comp},
			StartEffective: &past,
			EndEffective:   &future,
		},
		IncidentCriticality: event.CriticalityMajorOutage,
	}

	storage.CreateIncident(incident)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/weather", nil)
		rr := httptest.NewRecorder()
		handler.HandleWeather(rr, req)
	}
}

// BenchmarkAPIHandler_JSON_Encoding tests the performance of JSON encoding in responses
func BenchmarkAPIHandler_JSON_Encoding(b *testing.B) {
	storage := drivers.NewRAMStorage()
	cfg := &config.Config{}
	handler := NewAPIHandler(storage, cfg)

	// Create test data
	testData := map[string]interface{}{
		"service": "clariti",
		"version": "1.0.0",
		"status":  "operational",
		"data": map[string]interface{}{
			"components": []string{"api", "database", "cache"},
			"metrics": map[string]int{
				"requests": 1000,
				"errors":   5,
				"latency":  250,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.writeJSON(rr, http.StatusOK, testData)
	}
}

// BenchmarkJSON_Marshaling tests raw JSON marshaling performance
func BenchmarkJSON_Marshaling(b *testing.B) {
	testData := map[string]interface{}{
		"service": "clariti",
		"version": "1.0.0",
		"status":  "operational",
		"data": map[string]interface{}{
			"components": []string{"api", "database", "cache"},
			"metrics": map[string]int{
				"requests": 1000,
				"errors":   5,
				"latency":  250,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(testData)
	}
}
