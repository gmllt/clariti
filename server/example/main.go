package main

import (
	"log"

	"github.com/gmllt/clariti/server/core"
)

// Example of how to use the server programmatically
func main() {
	// Create server instance with custom config
	server, err := core.New("config.yaml")
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Access server components if needed
	config := server.GetConfig()
	storage := server.GetStorage()

	log.Printf("Server configuration loaded for %s:%s", config.Server.Host, config.Server.Port)
	log.Printf("Storage driver: %T", storage)

	// Run the server (this will block until interrupted)
	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
