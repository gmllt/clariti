package routes

import (
	"encoding/json"
	"net/http"
)

// APIDoc represents the API documentation structure
type APIDoc struct {
	Service   string             `json:"service"`
	Version   string             `json:"version"`
	Endpoints map[string][]Route `json:"endpoints"`
}

// Route represents a single API route
type Route struct {
	Path         string   `json:"path"`
	Methods      []string `json:"methods"`
	Description  string   `json:"description"`
	AuthRequired bool     `json:"auth_required"`
}

// setupDocumentationRoutes adds API documentation endpoint
func setupDocumentationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		docs := APIDoc{
			Service: "Clariti API",
			Version: "1.0.0",
			Endpoints: map[string][]Route{
				"health": {
					{
						Path:         "/health",
						Methods:      []string{"GET"},
						Description:  "Server health check",
						AuthRequired: false,
					},
				},
				"components": {
					{
						Path:         "/api/components",
						Methods:      []string{"GET"},
						Description:  "Get all component information (flattened)",
						AuthRequired: false,
					},
					{
						Path:         "/api/components/hierarchy",
						Methods:      []string{"GET"},
						Description:  "Get hierarchical component structure",
						AuthRequired: false,
					},
					{
						Path:         "/api/platforms",
						Methods:      []string{"GET"},
						Description:  "Get all platforms",
						AuthRequired: false,
					},
					{
						Path:         "/api/instances",
						Methods:      []string{"GET"},
						Description:  "Get all instances",
						AuthRequired: false,
					},
					{
						Path:         "/api/components/list",
						Methods:      []string{"GET"},
						Description:  "Get all components with relationships",
						AuthRequired: false,
					},
				},
				"incidents": {
					{
						Path:         "/api/incidents",
						Methods:      []string{"GET", "POST"},
						Description:  "List all incidents or create new incident",
						AuthRequired: false, // GET is free, POST requires auth
					},
					{
						Path:         "/api/incidents/{id}",
						Methods:      []string{"GET", "PUT", "DELETE"},
						Description:  "Get, update or delete specific incident",
						AuthRequired: false, // GET is free, PUT/DELETE require auth
					},
				},
				"planned-maintenances": {
					{
						Path:         "/api/planned-maintenances",
						Methods:      []string{"GET", "POST"},
						Description:  "List all planned maintenances or create new one",
						AuthRequired: false, // GET is free, POST requires auth
					},
					{
						Path:         "/api/planned-maintenances/{id}",
						Methods:      []string{"GET", "PUT", "DELETE"},
						Description:  "Get, update or delete specific planned maintenance",
						AuthRequired: false, // GET is free, PUT/DELETE require auth
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(docs)
	})
}
