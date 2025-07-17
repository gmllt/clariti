package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gmllt/clariti/models/component"
	"github.com/gmllt/clariti/models/event"
	"github.com/gmllt/clariti/server/drivers"
)

// PlannedMaintenanceHandler handles planned maintenance HTTP requests
type PlannedMaintenanceHandler struct {
	storage drivers.EventStorage
}

// PlannedMaintenanceRequest represents the JSON structure for creating/updating planned maintenances
type PlannedMaintenanceRequest struct {
	GUID         string    `json:"guid,omitempty"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Components   []string  `json:"components"` // Component codes as strings
	StartPlanned time.Time `json:"start_planned"`
	EndPlanned   time.Time `json:"end_planned"`
	Cancelled    bool      `json:"cancelled,omitempty"`
}

// ToPlannedMaintenance converts PlannedMaintenanceRequest to event.PlannedMaintenance by resolving component codes
func (req *PlannedMaintenanceRequest) ToPlannedMaintenance() *event.PlannedMaintenance {
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

	return &event.PlannedMaintenance{
		BaseEvent: event.BaseEvent{
			GUID:       req.GUID,
			Title:      req.Title,
			Content:    req.Content,
			Components: components,
		},
		StartPlanned: req.StartPlanned,
		EndPlanned:   req.EndPlanned,
		Cancelled:    req.Cancelled,
	}
}

// NewPlannedMaintenanceHandler creates a new planned maintenance handler
func NewPlannedMaintenanceHandler(storage drivers.EventStorage) *PlannedMaintenanceHandler {
	return &PlannedMaintenanceHandler{storage: storage}
}

// writeJSON writes a JSON response
func (h *PlannedMaintenanceHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// writeError writes an error response
func (h *PlannedMaintenanceHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}

// validatePlannedMaintenanceRequest validates planned maintenance request data and returns detailed error messages
func (h *PlannedMaintenanceHandler) validatePlannedMaintenanceRequest(req *PlannedMaintenanceRequest) []string {
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

	// Validate time fields
	if req.StartPlanned.IsZero() {
		errors = append(errors, "field 'start_planned' is required and must be a valid ISO 8601 timestamp")
	}

	if req.EndPlanned.IsZero() {
		errors = append(errors, "field 'end_planned' is required and must be a valid ISO 8601 timestamp")
	}

	// Validate time logic
	if !req.StartPlanned.IsZero() && !req.EndPlanned.IsZero() {
		if req.EndPlanned.Before(req.StartPlanned) {
			errors = append(errors, "field 'end_planned' must be after 'start_planned'")
		}
		if req.StartPlanned.Before(time.Now().Add(-24 * time.Hour)) {
			errors = append(errors, "field 'start_planned' cannot be more than 24 hours in the past")
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

// validatePlannedMaintenance validates planned maintenance data and returns detailed error messages
func (h *PlannedMaintenanceHandler) validatePlannedMaintenance(maintenance *event.PlannedMaintenance) []string {
	var errors []string

	// Validate required fields
	if strings.TrimSpace(maintenance.Title) == "" {
		errors = append(errors, "field 'title' is required and cannot be empty")
	}

	if strings.TrimSpace(maintenance.Content) == "" {
		errors = append(errors, "field 'content' is required and cannot be empty")
	}

	if len(maintenance.Components) == 0 {
		errors = append(errors, "field 'components' is required and must contain at least one component")
	}

	// Validate time fields
	if maintenance.StartPlanned.IsZero() {
		errors = append(errors, "field 'start_planned' is required and must be a valid ISO 8601 timestamp")
	}

	if maintenance.EndPlanned.IsZero() {
		errors = append(errors, "field 'end_planned' is required and must be a valid ISO 8601 timestamp")
	}

	// Validate time logic
	if !maintenance.StartPlanned.IsZero() && !maintenance.EndPlanned.IsZero() {
		if maintenance.EndPlanned.Before(maintenance.StartPlanned) {
			errors = append(errors, "field 'end_planned' must be after 'start_planned'")
		}
		if maintenance.StartPlanned.Before(time.Now().Add(-24 * time.Hour)) {
			errors = append(errors, "field 'start_planned' cannot be more than 24 hours in the past")
		}
	}

	// Validate component names (basic validation)
	for i, comp := range maintenance.Components {
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

// HandlePlannedMaintenances handles /planned-maintenances endpoint
func (h *PlannedMaintenanceHandler) HandlePlannedMaintenances(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllPlannedMaintenances(w, r)
	case http.MethodPost:
		h.createPlannedMaintenance(w, r)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandlePlannedMaintenanceByID handles /planned-maintenances/{id} endpoint
func (h *PlannedMaintenanceHandler) HandlePlannedMaintenanceByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path parameter
	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "Missing planned maintenance ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getPlannedMaintenance(w, r, id)
	case http.MethodPut:
		h.updatePlannedMaintenance(w, r, id)
	case http.MethodDelete:
		h.deletePlannedMaintenance(w, r, id)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// getAllPlannedMaintenances returns all planned maintenances
func (h *PlannedMaintenanceHandler) getAllPlannedMaintenances(w http.ResponseWriter, r *http.Request) {
	maintenances, err := h.storage.GetAllPlannedMaintenances()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to retrieve planned maintenances")
		return
	}
	h.writeJSON(w, http.StatusOK, maintenances)
}

// getPlannedMaintenance returns a specific planned maintenance
func (h *PlannedMaintenanceHandler) getPlannedMaintenance(w http.ResponseWriter, r *http.Request, id string) {
	maintenance, err := h.storage.GetPlannedMaintenance(id)
	if err != nil {
		if err == drivers.ErrNotFound {
			h.writeError(w, http.StatusNotFound, "Planned maintenance not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to retrieve planned maintenance")
		return
	}
	h.writeJSON(w, http.StatusOK, maintenance)
}

// createPlannedMaintenance creates a new planned maintenance
func (h *PlannedMaintenanceHandler) createPlannedMaintenance(w http.ResponseWriter, r *http.Request) {
	var req PlannedMaintenanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	// Validate planned maintenance request data
	if validationErrors := h.validatePlannedMaintenanceRequest(&req); len(validationErrors) > 0 {
		h.writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":             "Validation failed",
			"validation_errors": validationErrors,
		})
		return
	}

	// Convert request to planned maintenance
	maintenance := req.ToPlannedMaintenance()

	// Ensure GUID is generated if not provided
	if maintenance.GUID == "" {
		// Create new planned maintenance with GUID - use provided times or current time
		newMaintenance := event.NewPlannedMaintenance(maintenance.Title, maintenance.Content, maintenance.Components, maintenance.StartPlanned, maintenance.EndPlanned)
		maintenance = newMaintenance
	}

	if err := h.storage.CreatePlannedMaintenance(maintenance); err != nil {
		if err == drivers.ErrExists {
			h.writeError(w, http.StatusConflict, "Planned maintenance already exists")
			return
		}
		h.writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to create planned maintenance",
			"details": err.Error(),
		})
		return
	}

	h.writeJSON(w, http.StatusCreated, maintenance)
}

// updatePlannedMaintenance updates an existing planned maintenance
func (h *PlannedMaintenanceHandler) updatePlannedMaintenance(w http.ResponseWriter, r *http.Request, id string) {
	var req PlannedMaintenanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	// Validate planned maintenance request data
	if validationErrors := h.validatePlannedMaintenanceRequest(&req); len(validationErrors) > 0 {
		h.writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":             "Validation failed",
			"validation_errors": validationErrors,
		})
		return
	}

	// Convert request to planned maintenance
	maintenance := req.ToPlannedMaintenance()

	// Ensure the ID in the URL matches the maintenance GUID
	maintenance.GUID = id

	if err := h.storage.UpdatePlannedMaintenance(maintenance); err != nil {
		if err == drivers.ErrNotFound {
			h.writeError(w, http.StatusNotFound, "Planned maintenance not found")
			return
		}
		h.writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to update planned maintenance",
			"details": err.Error(),
		})
		return
	}

	h.writeJSON(w, http.StatusOK, maintenance)
}

// deletePlannedMaintenance deletes a planned maintenance
func (h *PlannedMaintenanceHandler) deletePlannedMaintenance(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.storage.DeletePlannedMaintenance(id); err != nil {
		if err == drivers.ErrNotFound {
			h.writeError(w, http.StatusNotFound, "Planned maintenance not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to delete planned maintenance")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
