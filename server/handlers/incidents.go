package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gmllt/clariti/models/component"
	"github.com/gmllt/clariti/models/event"
	"github.com/gmllt/clariti/server/drivers"
)

// IncidentHandler handles incident-related HTTP requests
type IncidentHandler struct {
	storage drivers.EventStorage
}

// IncidentRequest represents the JSON structure for creating/updating incidents
type IncidentRequest struct {
	GUID        string   `json:"guid,omitempty"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Components  []string `json:"components"`            // Component codes as strings
	Criticality string   `json:"criticality,omitempty"` // Criticality as string
	Perpetual   bool     `json:"perpetual,omitempty"`
}

// parseCriticality converts string criticality to event.Criticality
func parseCriticality(criticalityStr string) (event.Criticality, error) {
	switch strings.ToLower(strings.TrimSpace(criticalityStr)) {
	case "operational":
		return event.CriticalityOperational, nil
	case "degraded":
		return event.CriticalityDegraded, nil
	case "partial outage", "partial_outage":
		return event.CriticalityPartialOutage, nil
	case "major outage", "major_outage":
		return event.CriticalityMajorOutage, nil
	case "under maintenance", "under_maintenance", "maintenance":
		return event.CriticalityUnderMaintenance, nil
	case "unknown", "":
		return event.CriticalityUnknown, nil
	default:
		return event.CriticalityUnknown, fmt.Errorf("invalid criticality '%s'. Valid values: operational, degraded, partial outage, major outage, under maintenance, unknown", criticalityStr)
	}
}

// ToIncident converts IncidentRequest to event.Incident by resolving component codes
func (req *IncidentRequest) ToIncident() (*event.Incident, error) {
	// Convert component codes to Component objects
	var components []*component.Component
	for _, code := range req.Components {
		// Create minimal component objects - in a real scenario, you'd want to
		// validate these against your component registry
		comp := &component.Component{
			BaseComponent: component.BaseComponent{
				Name: code, // For now, use code as name - could be improved
				Code: code,
			},
		}
		components = append(components, comp)
	}

	// Parse criticality
	criticality, err := parseCriticality(req.Criticality)
	if err != nil {
		return nil, err
	}

	return &event.Incident{
		BaseEvent: event.BaseEvent{
			GUID:       req.GUID,
			Title:      req.Title,
			Content:    req.Content,
			Components: components,
		},
		IncidentCriticality: criticality,
		Perpetual:           req.Perpetual,
	}, nil
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

// validateIncidentRequest validates incident request data and returns detailed error messages
func (h *IncidentHandler) validateIncidentRequest(req *IncidentRequest) []string {
	var errors []string

	// Validate required fields
	if strings.TrimSpace(req.Title) == "" {
		errors = append(errors, "field 'title' is required and cannot be empty")
	}

	if strings.TrimSpace(req.Content) == "" {
		errors = append(errors, "field 'content' is required and cannot be empty")
	}

	if len(req.Components) == 0 {
		errors = append(errors, "field 'components' is required and must contain at least one component")
	}

	// Validate criticality string
	if req.Criticality != "" {
		if _, err := parseCriticality(req.Criticality); err != nil {
			errors = append(errors, err.Error())
		}
	}

	// Validate component codes (basic validation)
	for i, code := range req.Components {
		if strings.TrimSpace(code) == "" {
			errors = append(errors, fmt.Sprintf("component at index %d has empty code", i))
		}
	}

	return errors
}

// validateIncident validates incident data and returns detailed error messages
func (h *IncidentHandler) validateIncident(incident *event.Incident) []string {
	var errors []string

	// Validate required fields
	if strings.TrimSpace(incident.Title) == "" {
		errors = append(errors, "field 'title' is required and cannot be empty")
	}

	if strings.TrimSpace(incident.Content) == "" {
		errors = append(errors, "field 'content' is required and cannot be empty")
	}

	if len(incident.Components) == 0 {
		errors = append(errors, "field 'components' is required and must contain at least one component")
	}

	// Validate criticality range
	if incident.IncidentCriticality < -1 || incident.IncidentCriticality > 4 {
		errors = append(errors, fmt.Sprintf("field 'criticality' must be between -1 and 4, got %d (0=operational, 1=degraded, 2=partial outage, 3=major outage, 4=maintenance, -1=unknown)", incident.IncidentCriticality))
	}

	// Validate component names (basic validation)
	for i, comp := range incident.Components {
		if comp == nil {
			errors = append(errors, fmt.Sprintf("component at index %d is null", i))
			continue
		}
		if strings.TrimSpace(comp.Code) == "" {
			errors = append(errors, fmt.Sprintf("component at index %d has empty code", i))
		}
	}

	return errors
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
	// Extract ID from path parameter
	id := r.PathValue("id")
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
	var req IncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	// Validate incident request data
	if validationErrors := h.validateIncidentRequest(&req); len(validationErrors) > 0 {
		h.writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":             "Validation failed",
			"validation_errors": validationErrors,
		})
		return
	}

	// Convert request to incident
	incident, err := req.ToIncident()
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid criticality value",
			"details": err.Error(),
		})
		return
	}

	// Ensure GUID is generated if not provided
	if incident.GUID == "" {
		// Create new incident with GUID
		newIncident := event.NewIncident(incident.Title, incident.Content, incident.Components, incident.IncidentCriticality, incident.Perpetual)
		incident = newIncident
	}

	if err := h.storage.CreateIncident(incident); err != nil {
		if err == drivers.ErrExists {
			h.writeError(w, http.StatusConflict, "Incident already exists")
			return
		}
		h.writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to create incident",
			"details": err.Error(),
		})
		return
	}

	h.writeJSON(w, http.StatusCreated, incident)
}

// updateIncident updates an existing incident
func (h *IncidentHandler) updateIncident(w http.ResponseWriter, r *http.Request, id string) {
	var req IncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	// Validate incident request data
	if validationErrors := h.validateIncidentRequest(&req); len(validationErrors) > 0 {
		h.writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":             "Validation failed",
			"validation_errors": validationErrors,
		})
		return
	}

	// Convert request to incident
	incident, err := req.ToIncident()
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid criticality value",
			"details": err.Error(),
		})
		return
	}

	// Ensure the ID in the URL matches the incident GUID
	incident.GUID = id

	if err := h.storage.UpdateIncident(incident); err != nil {
		if err == drivers.ErrNotFound {
			h.writeError(w, http.StatusNotFound, "Incident not found")
			return
		}
		h.writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to update incident",
			"details": err.Error(),
		})
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
