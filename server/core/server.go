package core

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// Initialize storage driver
	storage := drivers.NewRAMStorage()

	// Initialize handlers
	handlers := handlers.New(storage, cfg)

	// Create HTTP server
	mux := http.NewServeMux()

	// Setup routes
	routes.Setup(mux, handlers, cfg)

	// Apply middleware
	handler := middleware.CORS(middleware.BasicAuth(cfg)(mux))

	httpServer := &http.Server{
		Addr:    cfg.GetAddress(),
		Handler: handler,
	}

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

// Run starts the server and handles graceful shutdown
func (s *Server) Run() error {
	// Channel to listen for interrupt signal to terminate
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting Clariti server on %s", s.config.GetAddress())
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	log.Println("Server stopped gracefully")
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

// Handler returns the HTTP handler for testing purposes
func (s *Server) Handler() http.Handler {
	return s.httpServer.Handler
}
