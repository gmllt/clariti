package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gmllt/clariti/server/config"
	"github.com/gmllt/clariti/server/drivers"
)

// APIHandler handles HTTP requests for the API
type APIHandler struct {
	storage drivers.EventStorage
	config  *config.Config
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(storage drivers.EventStorage, config *config.Config) *APIHandler {
	return &APIHandler{
		storage: storage,
		config:  config,
	}
}

// writeJSON writes a JSON response
func (h *APIHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// writeError writes an error response
func (h *APIHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}

// extractID extracts the ID from the URL path
func extractID(path, prefix string) string {
	if strings.HasPrefix(path, prefix) {
		return strings.TrimPrefix(path, prefix)
	}
	return ""
}

// Health check endpoint
func (h *APIHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "clariti-api",
	})
}

// Component info endpoints (read-only)
func (h *APIHandler) HandleComponents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "Only GET method allowed")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"platforms":  h.config.GetAllPlatforms(),
		"instances":  h.config.GetAllInstances(),
		"components": h.config.GetAllComponents(),
	})
}

func (h *APIHandler) HandlePlatforms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "Only GET method allowed")
		return
	}

	h.writeJSON(w, http.StatusOK, h.config.GetAllPlatforms())
}

func (h *APIHandler) HandleInstances(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "Only GET method allowed")
		return
	}

	h.writeJSON(w, http.StatusOK, h.config.GetAllInstances())
}

func (h *APIHandler) HandleComponentsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "Only GET method allowed")
		return
	}

	h.writeJSON(w, http.StatusOK, h.config.GetAllComponents())
}

// HandleComponentsHierarchy returns the full hierarchical structure
func (h *APIHandler) HandleComponentsHierarchy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "Only GET method allowed")
		return
	}

	h.writeJSON(w, http.StatusOK, h.config.Components)
}
