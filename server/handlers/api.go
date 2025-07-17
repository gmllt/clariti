package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gmllt/clariti/logger"
	"github.com/gmllt/clariti/server/config"
	"github.com/gmllt/clariti/server/drivers"
	"github.com/gmllt/clariti/utils"
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

// writeJSON writes a JSON response using object pools for better performance
func (h *APIHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// writeError writes an error response using pooled objects
func (h *APIHandler) writeError(w http.ResponseWriter, status int, message string) {
	errorMap := utils.GetStringMap()
	errorMap["error"] = message
	defer utils.PutStringMap(errorMap)

	h.writeJSON(w, status, errorMap)
}

// Health check endpoint with optimized response
func (h *APIHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	log := logger.GetDefault().WithComponent("APIHandler")
	log.Debug("Health check requested")

	healthMap := utils.GetStringMap()
	healthMap["status"] = "healthy"
	healthMap["service"] = "clariti-api"
	defer utils.PutStringMap(healthMap)

	log.Info("Health check completed successfully")
	h.writeJSON(w, http.StatusOK, healthMap)
}

// Component info endpoints (read-only) with optimized allocations
func (h *APIHandler) HandleComponents(w http.ResponseWriter, r *http.Request) {
	log := logger.GetDefault().WithComponent("APIHandler")
	log.WithField("method", r.Method).Debug("Components request received")

	if r.Method != http.MethodGet {
		log.WithField("method", r.Method).Warn("Method not allowed for components endpoint")
		h.writeError(w, http.StatusMethodNotAllowed, "Only GET method allowed")
		return
	}

	response := utils.GetJSONResponse()
	response["platforms"] = h.config.GetAllPlatforms()
	response["instances"] = h.config.GetAllInstances()
	response["components"] = h.config.GetAllComponents()
	defer utils.PutJSONResponse(response)

	log.Info("Components data retrieved successfully")
	h.writeJSON(w, http.StatusOK, response)
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
	log := logger.GetDefault().WithComponent("APIHandler")
	log.Debug("Components list request received")

	if r.Method != http.MethodGet {
		log.WithField("method", r.Method).Warn("Method not allowed for components list endpoint")
		h.writeError(w, http.StatusMethodNotAllowed, "Only GET method allowed")
		return
	}

	components := h.config.GetAllComponents()
	log.WithField("component_count", len(components)).Info("Components list retrieved successfully")
	h.writeJSON(w, http.StatusOK, components)
}

// HandleComponentsHierarchy returns the full hierarchical structure
func (h *APIHandler) HandleComponentsHierarchy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "Only GET method allowed")
		return
	}

	h.writeJSON(w, http.StatusOK, h.config.Components)
}
