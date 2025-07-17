package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	serverURL    string
	username     string
	password     string
	outputFormat string
	traceFile    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clariti-cli",
	Short: "Clariti CLI - Manage incidents and maintenance",
	Long: `Clariti CLI is a command line tool to interact with Clariti API.
	
Environment Variables:
  CLARITI_SERVER_URL    - Server URL (default: http://localhost:8080)
  CLARITI_USERNAME      - Basic auth username
  CLARITI_PASSWORD      - Basic auth password
  CLARITI_OUTPUT_FORMAT - Output format: json, yaml, table (default: pretty)
  CLARITI_TRACE_FILE    - Trace output file (default: stdout)
  CLARITI_TRACE_ENABLED - Enable trace output (default: false)
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", getEnvWithDefault("CLARITI_SERVER_URL", "http://localhost:8080"), "Clariti server URL")
	rootCmd.PersistentFlags().StringVar(&username, "username", os.Getenv("CLARITI_USERNAME"), "Basic auth username")
	rootCmd.PersistentFlags().StringVar(&password, "password", os.Getenv("CLARITI_PASSWORD"), "Basic auth password")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "output", getEnvWithDefault("CLARITI_OUTPUT_FORMAT", ""), "Output format (json|yaml|table)")
	rootCmd.PersistentFlags().StringVar(&traceFile, "trace-file", os.Getenv("CLARITI_TRACE_FILE"), "Trace output file (empty for stdout)")
}

func initConfig() {
	// Initialize tracing if enabled
	if os.Getenv("CLARITI_TRACE_ENABLED") == "true" {
		initTracing()
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func initTracing() {
	// Setup tracing output
	traceOutput := os.Stdout
	if traceFile != "" {
		file, err := os.OpenFile(traceFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open trace file: %v\n", err)
			return
		}
		traceOutput = file
	}

	fmt.Fprintf(traceOutput, "[TRACE] Clariti CLI started\n")
	fmt.Fprintf(traceOutput, "[TRACE] Server URL: %s\n", serverURL)
	fmt.Fprintf(traceOutput, "[TRACE] Output format: %s\n", outputFormat)
}
