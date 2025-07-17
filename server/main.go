package main

import (
	"log"
	"os"

	"github.com/gmllt/clariti/server/core"
)

func main() {
	// Get configuration path
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	// Create server instance
	server, err := core.New(configPath)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Run server with graceful shutdown
	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
