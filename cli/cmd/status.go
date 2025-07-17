package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check server health",
	Long:  "Check if the Clariti server is running and healthy",
	RunE: func(cmd *cobra.Command, args []string) error {
		trace("Executing health command")

		client := getAPIClient()
		resp, err := client.makeRequest("GET", "/health", nil)
		if err != nil {
			return fmt.Errorf("health check failed: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("server unhealthy: %s", string(body))
		}

		return outputData(body, "Health check successful")
	},
}

// weatherCmd represents the weather command
var weatherCmd = &cobra.Command{
	Use:   "weather",
	Short: "Get service status overview",
	Long:  "Get the overall weather status of all services",
	RunE: func(cmd *cobra.Command, args []string) error {
		trace("Executing weather command")

		client := getAPIClient()
		resp, err := client.makeRequest("GET", "/api/v1/weather", nil)
		if err != nil {
			return fmt.Errorf("weather check failed: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("weather request failed: %s", string(body))
		}

		return outputData(body, "Service weather")
	},
}

// outputData formats and outputs data based on the selected format
func outputData(data []byte, title string) error {
	trace("Outputting data in format: %s", outputFormat)

	switch outputFormat {
	case "json":
		// Pretty print JSON
		var jsonData interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			return fmt.Errorf("invalid JSON response: %w", err)
		}
		formatted, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format JSON: %w", err)
		}
		fmt.Println(string(formatted))

	case "yaml":
		// Convert JSON to YAML (simple implementation)
		var jsonData interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			return fmt.Errorf("invalid JSON response: %w", err)
		}
		// For now, just pretty print JSON (YAML support can be added later)
		formatted, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format data: %w", err)
		}
		fmt.Println("# YAML output (formatted as JSON for now)")
		fmt.Println(string(formatted))

	case "table":
		// Legacy table format for compatibility
		fmt.Printf("=== %s ===\n", title)
		var jsonData map[string]interface{}
		if err := json.Unmarshal(data, &jsonData); err == nil {
			// For objects, display key-value pairs
			for key, value := range jsonData {
				switch v := value.(type) {
				case string:
					fmt.Printf("%-20s: %s\n", key, v)
				case float64:
					fmt.Printf("%-20s: %.0f\n", key, v)
				case bool:
					fmt.Printf("%-20s: %t\n", key, v)
				default:
					fmt.Printf("%-20s: %v\n", key, v)
				}
			}
		} else {
			// For arrays or complex data, try to display nicely
			var jsonArray []interface{}
			if err := json.Unmarshal(data, &jsonArray); err == nil {
				fmt.Printf("Found %d items:\n", len(jsonArray))
				for i, item := range jsonArray {
					fmt.Printf("--- Item %d ---\n", i+1)
					if itemMap, ok := item.(map[string]interface{}); ok {
						for key, value := range itemMap {
							fmt.Printf("%-20s: %v\n", key, value)
						}
					} else {
						fmt.Printf("%v\n", item)
					}
				}
			} else {
				// Fallback to raw output
				fmt.Println(string(data))
			}
		}

	default:
		// Default to pretty style for better UX
		return outputDataPretty(data, title)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(weatherCmd)
}
