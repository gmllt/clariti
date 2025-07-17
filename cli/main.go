package main

import (
	"fmt"
	"os"

	"github.com/gmllt/clariti/cli/cmd"
)

func main() {
	// Execute CLI
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
