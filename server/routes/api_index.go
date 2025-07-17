package routes

import (
	"encoding/json"
	"net/http"
)

// APIVersionInfo represents information about an API version
type APIVersionInfo struct {
	Version     string `json:"version"`
	Status      string `json:"status"`
	DocsURL     string `json:"docs_url"`
	Description string `json:"description"`
}

// APIIndex represents the main API index with available versions
type APIIndex struct {
	Service     string           `json:"service"`
	Description string           `json:"description"`
	Versions    []APIVersionInfo `json:"versions"`
}

// setupAPIIndexRoutes adds the main API index endpoint
func setupAPIIndexRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		index := APIIndex{
			Service:     "Clariti API",
			Description: "REST API for managing incidents and planned maintenances with hierarchical component structure",
			Versions: []APIVersionInfo{
				{
					Version:     "v1",
					Status:      "stable",
					DocsURL:     "/api/v1/docs",
					Description: "Current stable version with full incident and planned maintenance management",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(index)
	})
}
