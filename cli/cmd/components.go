package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// componentsCmd represents the components command
var componentsCmd = &cobra.Command{
	Use:   "components",
	Short: "Manage components",
	Long:  "List and view component hierarchy",
}

// componentsListCmd lists all components
var componentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all components",
	Long:  "List all components from the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		trace("Executing components list command")

		client := getAPIClient()
		resp, err := client.makeRequest("GET", "/api/v1/components", nil)
		if err != nil {
			return fmt.Errorf("failed to list components: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("list components failed: %s", string(body))
		}

		return outputData(body, "Components")
	},
}

// componentsTreeCmd shows component hierarchy
var componentsTreeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Show component hierarchy",
	Long:  "Show the component hierarchy tree",
	RunE: func(cmd *cobra.Command, args []string) error {
		trace("Executing components tree command")

		client := getAPIClient()
		resp, err := client.makeRequest("GET", "/api/v1/components/hierarchy", nil)
		if err != nil {
			return fmt.Errorf("failed to get component tree: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("get component tree failed: %s", string(body))
		}

		return outputData(body, "Component Hierarchy")
	},
}

// platformsCmd lists platforms
var platformsCmd = &cobra.Command{
	Use:   "platforms",
	Short: "List platforms",
	Long:  "List all platforms",
	RunE: func(cmd *cobra.Command, args []string) error {
		trace("Executing platforms command")

		client := getAPIClient()
		resp, err := client.makeRequest("GET", "/api/v1/platforms", nil)
		if err != nil {
			return fmt.Errorf("failed to list platforms: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("list platforms failed: %s", string(body))
		}

		return outputData(body, "Platforms")
	},
}

// instancesCmd lists instances
var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "List instances",
	Long:  "List all instances",
	RunE: func(cmd *cobra.Command, args []string) error {
		trace("Executing instances command")

		client := getAPIClient()
		resp, err := client.makeRequest("GET", "/api/v1/instances", nil)
		if err != nil {
			return fmt.Errorf("failed to list instances: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("list instances failed: %s", string(body))
		}

		return outputData(body, "Instances")
	},
}

func init() {
	// Add components command to root
	rootCmd.AddCommand(componentsCmd)
	rootCmd.AddCommand(platformsCmd)
	rootCmd.AddCommand(instancesCmd)

	// Add subcommands to components
	componentsCmd.AddCommand(componentsListCmd)
	componentsCmd.AddCommand(componentsTreeCmd)
}
