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
	"github.com/gmllt/clariti/server/middleware"
	"github.com/gmllt/clariti/server/routes"
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

	// Create HTTP server
	log.Debug("Setting up HTTP routes")
	mux := http.NewServeMux()

	// Setup routes
	routes.Setup(mux, handlers, cfg)
	log.Info("HTTP routes configured")

	// Configure server with timeouts
	log.Debug("Configuring HTTP server with timeouts")
	httpServer := &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      middleware.CORS(middleware.BasicAuth(cfg)(mux)),
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

	// Create HTTP server
	mux := http.NewServeMux()

	// Setup routes
	routes.Setup(mux, handlers, cfg)

	// Configure server with timeouts
	httpServer := &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      middleware.CORS(middleware.BasicAuth(cfg)(mux)),
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
