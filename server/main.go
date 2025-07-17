package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/gmllt/clariti/server/core"
	"github.com/prometheus/common/version"
)

var (
	app        = kingpin.New("clariti-server", "Clariti status page server - A modern status page system")
	configPath = app.Flag("config", "Path to configuration file").Short('c').Envar("CONFIG_PATH").Default("config.yaml").String()

	// Commands
	versionInfo = app.Command("version", "Show detailed version information")
	serverCmd   = app.Command("serve", "Start the Clariti server").Default()
)

func main() {
	// Configure application metadata
	app.Version(version.Version)
	app.Author("GMLLT")
	app.HelpFlag.Short('h')

	// Parse command line arguments
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	// Handle commands
	switch command {
	case versionInfo.FullCommand():
		fmt.Println(version.Print("clariti-server"))
		return
	case serverCmd.FullCommand():
		// Default: start server
	}

	// Create server instance
	server, err := core.New(*configPath)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Run server with graceful shutdown
	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
