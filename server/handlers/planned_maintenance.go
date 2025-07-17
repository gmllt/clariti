package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gmllt/clariti/models/event"
	"github.com/gmllt/clariti/server/drivers"
)

// PlannedMaintenanceHandler handles planned maintenance HTTP requests
type PlannedMaintenanceHandler struct {
	storage drivers.EventStorage
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
	var maintenance event.PlannedMaintenance
	if err := json.NewDecoder(r.Body).Decode(&maintenance); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Ensure GUID is generated if not provided
	if maintenance.GUID == "" {
		// Create new planned maintenance with GUID - use provided times or current time
		newMaintenance := event.NewPlannedMaintenance(maintenance.Title, maintenance.Content, maintenance.Components, maintenance.StartPlanned, maintenance.EndPlanned)
		maintenance = *newMaintenance
	}

	if err := h.storage.CreatePlannedMaintenance(&maintenance); err != nil {
		if err == drivers.ErrExists {
			h.writeError(w, http.StatusConflict, "Planned maintenance already exists")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create planned maintenance")
		return
	}

	h.writeJSON(w, http.StatusCreated, maintenance)
}

// updatePlannedMaintenance updates an existing planned maintenance
func (h *PlannedMaintenanceHandler) updatePlannedMaintenance(w http.ResponseWriter, r *http.Request, id string) {
	var maintenance event.PlannedMaintenance
	if err := json.NewDecoder(r.Body).Decode(&maintenance); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Ensure the ID in the URL matches the maintenance GUID
	maintenance.GUID = id

	if err := h.storage.UpdatePlannedMaintenance(&maintenance); err != nil {
		if err == drivers.ErrNotFound {
			h.writeError(w, http.StatusNotFound, "Planned maintenance not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to update planned maintenance")
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
