package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gmllt/clariti/server/config"
	"github.com/gmllt/clariti/server/drivers"
)

// WeatherHandler handles weather-related requests
type WeatherHandler struct {
	weatherService *drivers.WeatherService
}

// NewWeatherHandler creates a new weather handler
func NewWeatherHandler(cfg *config.Config, storage *drivers.RAMStorage) *WeatherHandler {
	return &WeatherHandler{
		weatherService: drivers.NewWeatherService(cfg, storage),
	}
}

// HandleWeather returns the current "weather" status of all services
func (wh *WeatherHandler) HandleWeather(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	summary, err := wh.weatherService.GetWeatherSummary()
	if err != nil {
		http.Error(w, "Failed to get weather summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(summary); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
