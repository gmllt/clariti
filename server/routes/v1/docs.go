package v1

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/common/version"
)

// V1APIDoc represents the API v1 documentation structure
type V1APIDoc struct {
	Service   string               `json:"service"`
	Version   string               `json:"version"`
	BaseURL   string               `json:"base_url"`
	Endpoints map[string][]V1Route `json:"endpoints"`
}

// V1Route represents a single API v1 route
type V1Route struct {
	Path         string   `json:"path"`
	Methods      []string `json:"methods"`
	Description  string   `json:"description"`
	AuthRequired bool     `json:"auth_required"`
}

// SetupV1DocumentationRoutes adds API v1 documentation endpoint
func SetupV1DocumentationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/docs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		docs := V1APIDoc{
			Service: "Clariti API",
			Version: version.Version,
			BaseURL: "/api/v1",
			Endpoints: map[string][]V1Route{
				"components": {
					{
						Path:         "/api/v1/components",
						Methods:      []string{"GET"},
						Description:  "Get all component information (flattened)",
						AuthRequired: false,
					},
					{
						Path:         "/api/v1/components/hierarchy",
						Methods:      []string{"GET"},
						Description:  "Get hierarchical component structure",
						AuthRequired: false,
					},
					{
						Path:         "/api/v1/platforms",
						Methods:      []string{"GET"},
						Description:  "Get all platforms",
						AuthRequired: false,
					},
					{
						Path:         "/api/v1/instances",
						Methods:      []string{"GET"},
						Description:  "Get all instances",
						AuthRequired: false,
					},
					{
						Path:         "/api/v1/components/list",
						Methods:      []string{"GET"},
						Description:  "Get all components with relationships",
						AuthRequired: false,
					},
				},
				"incidents": {
					{
						Path:         "/api/v1/incidents",
						Methods:      []string{"GET", "POST"},
						Description:  "List all incidents or create new incident (POST requires auth)",
						AuthRequired: false,
					},
					{
						Path:         "/api/v1/incidents/{id}",
						Methods:      []string{"GET", "PUT", "DELETE"},
						Description:  "Get, update or delete specific incident (PUT/DELETE require auth)",
						AuthRequired: false,
					},
				},
				"planned-maintenances": {
					{
						Path:         "/api/v1/planned-maintenances",
						Methods:      []string{"GET", "POST"},
						Description:  "List all planned maintenances or create new one (POST requires auth)",
						AuthRequired: false,
					},
					{
						Path:         "/api/v1/planned-maintenances/{id}",
						Methods:      []string{"GET", "PUT", "DELETE"},
						Description:  "Get, update or delete specific planned maintenance (PUT/DELETE require auth)",
						AuthRequired: false,
					},
				},
				"weather": {
					{
						Path:         "/api/v1/weather",
						Methods:      []string{"GET"},
						Description:  "Get current service weather (status overview based on active incidents and maintenances)",
						AuthRequired: false,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(docs); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}
