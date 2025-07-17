package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gmllt/clariti/models/event"
	"github.com/gmllt/clariti/server/drivers"
)

// IncidentHandler handles incident-related HTTP requests
type IncidentHandler struct {
	storage drivers.EventStorage
}

// NewIncidentHandler creates a new incident handler
func NewIncidentHandler(storage drivers.EventStorage) *IncidentHandler {
	return &IncidentHandler{storage: storage}
}

// writeJSON writes a JSON response
func (h *IncidentHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// writeError writes an error response
func (h *IncidentHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}

// HandleIncidents handles /incidents endpoint
func (h *IncidentHandler) HandleIncidents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllIncidents(w, r)
	case http.MethodPost:
		h.createIncident(w, r)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleIncidentByID handles /incidents/{id} endpoint
func (h *IncidentHandler) HandleIncidentByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id := strings.TrimPrefix(r.URL.Path, "/api/incidents/")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "Missing incident ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getIncident(w, r, id)
	case http.MethodPut:
		h.updateIncident(w, r, id)
	case http.MethodDelete:
		h.deleteIncident(w, r, id)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// getAllIncidents returns all incidents
func (h *IncidentHandler) getAllIncidents(w http.ResponseWriter, r *http.Request) {
	incidents, err := h.storage.GetAllIncidents()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to retrieve incidents")
		return
	}
	h.writeJSON(w, http.StatusOK, incidents)
}

// getIncident returns a specific incident
func (h *IncidentHandler) getIncident(w http.ResponseWriter, r *http.Request, id string) {
	incident, err := h.storage.GetIncident(id)
	if err != nil {
		if err == drivers.ErrNotFound {
			h.writeError(w, http.StatusNotFound, "Incident not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to retrieve incident")
		return
	}
	h.writeJSON(w, http.StatusOK, incident)
}

// createIncident creates a new incident
func (h *IncidentHandler) createIncident(w http.ResponseWriter, r *http.Request) {
	var incident event.Incident
	if err := json.NewDecoder(r.Body).Decode(&incident); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Ensure GUID is generated if not provided
	if incident.GUID == "" {
		// Create new incident with GUID
		newIncident := event.NewIncident(incident.Title, incident.Content, incident.Components, incident.IncidentCriticality, incident.Perpetual)
		incident = *newIncident
	}

	if err := h.storage.CreateIncident(&incident); err != nil {
		if err == drivers.ErrExists {
			h.writeError(w, http.StatusConflict, "Incident already exists")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create incident")
		return
	}

	h.writeJSON(w, http.StatusCreated, incident)
}

// updateIncident updates an existing incident
func (h *IncidentHandler) updateIncident(w http.ResponseWriter, r *http.Request, id string) {
	var incident event.Incident
	if err := json.NewDecoder(r.Body).Decode(&incident); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Ensure the ID in the URL matches the incident GUID
	incident.GUID = id

	if err := h.storage.UpdateIncident(&incident); err != nil {
		if err == drivers.ErrNotFound {
			h.writeError(w, http.StatusNotFound, "Incident not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to update incident")
		return
	}

	h.writeJSON(w, http.StatusOK, incident)
}

// deleteIncident deletes an incident
func (h *IncidentHandler) deleteIncident(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.storage.DeleteIncident(id); err != nil {
		if err == drivers.ErrNotFound {
			h.writeError(w, http.StatusNotFound, "Incident not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to delete incident")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
