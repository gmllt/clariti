package routes

import (
	"net/http"

	"github.com/gmllt/clariti/server/config"
	"github.com/gmllt/clariti/server/handlers"
	v1 "github.com/gmllt/clariti/server/routes/v1"
)

// Setup configures all the routes for the application
func Setup(mux *http.ServeMux, h *handlers.Handlers, cfg *config.Config) {
	// Health check routes
	setupHealthRoutes(mux, h)

	// API index and documentation
	setupAPIIndexRoutes(mux)
	setupDocumentationRoutes(mux)

	// API v1 routes
	v1.SetupV1Routes(mux, h)
	v1.SetupV1DocumentationRoutes(mux)
}

// setupHealthRoutes configures health check endpoints
func setupHealthRoutes(mux *http.ServeMux, h *handlers.Handlers) {
	mux.HandleFunc("/health", h.API.HandleHealth)
}
