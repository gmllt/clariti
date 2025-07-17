package core

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gmllt/clariti/logger"
	"github.com/gmllt/clariti/server/config"
	"github.com/gmllt/clariti/server/drivers"
	"github.com/gmllt/clariti/server/handlers"
	"github.com/gmllt/clariti/server/metrics"
	"github.com/gmllt/clariti/server/middleware"
	"github.com/gmllt/clariti/server/routes"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Server represents the HTTP server with all its dependencies
type Server struct {
	config     *config.Config
	storage    drivers.EventStorage
	httpServer *http.Server
	handlers   *handlers.Handlers
}

// New creates a new server instance
func New(configPath string) (*Server, error) {
	log := logger.GetDefault().WithComponent("Server")
	log.WithField("config_path", configPath).Info("Creating new server instance")

	// Load configuration
	log.Debug("Loading configuration")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.WithError(err).Error("Failed to load configuration")
		return nil, err
	}
	log.Info("Configuration loaded successfully")

	// Initialize storage driver
	log.Info("Initializing storage driver")
	storage, err := drivers.NewStorage(cfg)
	if err != nil {
		log.WithError(err).Error("Failed to initialize storage driver")
		return nil, err
	}
	log.Info("Storage driver initialized successfully")

	// Initialize handlers
	log.Debug("Initializing request handlers")
	handlers := handlers.New(storage, cfg)

	// Initialize metrics
	log.Debug("Initializing application metrics")
	metrics.InitializeApplicationInfo()
	log.Info("Application metrics initialized")

	// Create HTTP server
	log.Debug("Setting up HTTP routes")
	mux := http.NewServeMux()

	// Setup metrics endpoint - control tower instrumentation
	mux.Handle("/metrics", promhttp.Handler())
	log.Debug("Metrics endpoint configured at /metrics")

	// Setup routes
	routes.Setup(mux, handlers, cfg)
	log.Info("HTTP routes configured")

	// Configure server with timeouts
	log.Debug("Configuring HTTP server with timeouts")
	httpServer := &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      middleware.CORS(middleware.BasicAuth(cfg)(middleware.MetricsMiddleware(logger.GetDefault())(mux))),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.WithField("address", cfg.GetAddress()).Info("Server instance created successfully")
	return &Server{
		config:     cfg,
		storage:    storage,
		httpServer: httpServer,
		handlers:   handlers,
	}, nil
}

// NewWithConfig creates a new server instance with provided config and storage (useful for testing)
func NewWithConfig(cfg *config.Config, storage drivers.EventStorage) *Server {
	// Initialize handlers
	handlers := handlers.New(storage, cfg)

	// Initialize metrics
	metrics.InitializeApplicationInfo()

	// Create HTTP server
	mux := http.NewServeMux()

	// Setup metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// Setup routes
	routes.Setup(mux, handlers, cfg)

	// Configure server with timeouts
	httpServer := &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      middleware.CORS(middleware.BasicAuth(cfg)(middleware.MetricsMiddleware(logger.GetDefault())(mux))),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &Server{
		config:     cfg,
		storage:    storage,
		httpServer: httpServer,
		handlers:   handlers,
	}
}

// Handler returns the HTTP handler for testing purposes
func (s *Server) Handler() http.Handler {
	return s.httpServer.Handler
}

// Run starts the server and handles graceful shutdown
func (s *Server) Run() error {
	log := logger.GetDefault().WithComponent("Server")
	log.Info("Starting server run sequence")

	// Channel to listen for interrupt signal to terminate
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	log.Debug("Signal handlers configured for graceful shutdown")

	// Start server in a goroutine
	go func() {
		scheme := "HTTP"
		if s.config.IsHTTPSEnabled() {
			scheme = "HTTPS"
			log.WithField("cert_file", s.config.Server.CertFile).Info("HTTPS mode enabled")
		}

		address := s.config.GetAddress()
		log.WithField("scheme", scheme).WithField("address", address).Info("Starting Clariti server")

		var err error
		if s.config.IsHTTPSEnabled() {
			log.Debug("Starting HTTPS server")
			err = s.httpServer.ListenAndServeTLS(s.config.Server.CertFile, s.config.Server.KeyFile)
		} else {
			log.Debug("Starting HTTP server")
			err = s.httpServer.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("Server failed to start")
			panic(err) // Fatal error
		}
	}()

	// Start metrics updater goroutine - flight operations monitor
	go s.updateMetricsPeriodically(log)

	// Wait for interrupt signal
	log.Info("Server ready, waiting for requests")
	<-stop
	log.Warn("Shutdown signal received, stopping server")

	// Create a context with timeout for graceful shutdown
	log.Info("Starting graceful shutdown with 30s timeout")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server forced to shutdown")
		return err
	}

	log.Info("Server stopped gracefully")
	return nil
}

// GetConfig returns the server configuration
func (s *Server) GetConfig() *config.Config {
	return s.config
}

// GetStorage returns the storage driver
func (s *Server) GetStorage() drivers.EventStorage {
	return s.storage
}

// updateMetricsPeriodically updates business metrics every 30 seconds - operational status monitoring
func (s *Server) updateMetricsPeriodically(log *logrus.Entry) {
	log.Debug("Starting periodic metrics updater - flight operations monitor")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.updateBusinessMetrics(log)
		}
	}
}

// updateBusinessMetrics collects and updates business-related metrics - operational status collection
func (s *Server) updateBusinessMetrics(log *logrus.Entry) {
	log.Debug("Updating business metrics - collecting operational status")

	// Update uptime
	metrics.UpdateUptimeMetrics()

	// Count components
	componentCount := 0
	for _, platform := range s.config.Components.Platforms {
		for _, instance := range platform.Instances {
			componentCount += len(instance.Components)
		}
	}
	metrics.UpdateComponentsCount(float64(componentCount))

	// Get real data from storage - flight operations data collection

	// Count incidents by severity and status
	incidents, err := s.storage.GetAllIncidents()
	if err != nil {
		log.WithError(err).Warn("Failed to retrieve incidents for metrics - using default values")
		// Set default values on error
		metrics.UpdateIncidentMetrics("critical", "active", 0)
		metrics.UpdateIncidentMetrics("major", "active", 0)
		metrics.UpdateIncidentMetrics("minor", "active", 0)
	} else {
		// Count incidents by severity and status
		incidentCounts := make(map[string]map[string]int)
		for _, incident := range incidents {
			severity := incident.Criticality().String()
			status := string(incident.Status())

			if incidentCounts[severity] == nil {
				incidentCounts[severity] = make(map[string]int)
			}
			incidentCounts[severity][status]++
		}

		// Update metrics for all combinations
		severities := []string{"operational", "degraded", "partial outage", "major outage", "under maintenance", "unknown"}
		statuses := []string{"ongoing", "resolved", "unknown"}

		for _, severity := range severities {
			for _, status := range statuses {
				count := 0
				if incidentCounts[severity] != nil {
					count = incidentCounts[severity][status]
				}
				metrics.UpdateIncidentMetrics(severity, status, float64(count))
			}
		}

		log.WithField("total_incidents", len(incidents)).Debug("Incident metrics updated from storage")
	}

	// Count planned maintenances by status
	maintenances, err := s.storage.GetAllPlannedMaintenances()
	if err != nil {
		log.WithError(err).Warn("Failed to retrieve planned maintenances for metrics - using default values")
		// Set default values on error
		metrics.UpdatePlannedMaintenanceMetrics("scheduled", 0)
		metrics.UpdatePlannedMaintenanceMetrics("active", 0)
	} else {
		// Count maintenances by status
		maintenanceCounts := make(map[string]int)
		for _, maintenance := range maintenances {
			status := string(maintenance.Status())
			maintenanceCounts[status]++
		}

		// Update metrics for common statuses
		statuses := []string{"scheduled", "ongoing", "completed", "unknown"}
		for _, status := range statuses {
			count := maintenanceCounts[status]
			metrics.UpdatePlannedMaintenanceMetrics(status, float64(count))
		}

		log.WithField("total_maintenances", len(maintenances)).Debug("Planned maintenance metrics updated from storage")
	}

	log.WithField("component_count", componentCount).Debug("Business metrics updated - operational status collected")
}
