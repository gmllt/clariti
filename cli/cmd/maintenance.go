package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
)

// Maintenance represents a planned maintenance (using the same structure as the server)
type Maintenance struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ComponentID string `json:"component_id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// maintenanceCmd represents the maintenance command
var maintenanceCmd = &cobra.Command{
	Use:   "maintenance",
	Short: "Manage planned maintenances",
	Long:  "Create, list, update and delete planned maintenances",
}

// maintenanceListCmd lists all maintenances
var maintenanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all planned maintenances",
	Long:  "List all planned maintenances from the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		trace("Executing maintenance list command")

		client := getAPIClient()
		resp, err := client.makeRequest("GET", "/api/v1/planned-maintenances", nil)
		if err != nil {
			return fmt.Errorf("failed to list maintenances: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("list maintenances failed: %s", string(body))
		}

		return outputData(body, "Planned Maintenances")
	},
}

// maintenanceGetCmd gets a specific maintenance
var maintenanceGetCmd = &cobra.Command{
	Use:   "get [maintenance-id]",
	Short: "Get a specific planned maintenance",
	Long:  "Get details of a specific planned maintenance by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		maintenanceID := args[0]
		trace("Executing maintenance get command for ID: %s", maintenanceID)

		client := getAPIClient()
		resp, err := client.makeRequest("GET", "/api/v1/planned-maintenances/"+maintenanceID, nil)
		if err != nil {
			return fmt.Errorf("failed to get maintenance: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("get maintenance failed: %s", string(body))
		}

		return outputData(body, "Maintenance Details")
	},
}

// maintenanceCreateCmd creates a new maintenance
var maintenanceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new planned maintenance",
	Long:  "Create a new planned maintenance with title, description, and schedule",
	RunE: func(cmd *cobra.Command, args []string) error {
		trace("Executing maintenance create command")

		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		componentID, _ := cmd.Flags().GetString("component")
		startTime, _ := cmd.Flags().GetString("start-time")
		endTime, _ := cmd.Flags().GetString("end-time")

		if title == "" {
			return fmt.Errorf("title is required")
		}

		// Use default times if not provided
		now := time.Now()
		if startTime == "" {
			startTime = now.Add(1 * time.Hour).Format(time.RFC3339) // Default to 1 hour from now
		}
		if endTime == "" {
			startTimeObj, err := time.Parse(time.RFC3339, startTime)
			if err != nil {
				return fmt.Errorf("invalid start time format: %w", err)
			}
			endTime = startTimeObj.Add(2 * time.Hour).Format(time.RFC3339) // Default to 2 hours after start
		}

		// Create maintenance request with proper format for server
		maintenance := map[string]interface{}{
			"title":         title,
			"content":       description, // Server expects 'content', not 'description'
			"start_planned": startTime,
			"end_planned":   endTime,
		}

		// Add component if specified (server expects array of strings)
		if componentID != "" {
			maintenance["components"] = []string{componentID}
		} else {
			// Server validation requires at least one component
			maintenance["components"] = []string{"general"}
		}

		client := getAPIClient()
		resp, err := client.makeRequest("POST", "/api/v1/planned-maintenances", maintenance)
		if err != nil {
			return fmt.Errorf("failed to create maintenance: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 201 {
			return fmt.Errorf("create maintenance failed: %s", string(body))
		}

		return outputData(body, "Maintenance Created")
	},
}

// maintenanceUpdateCmd updates a maintenance
var maintenanceUpdateCmd = &cobra.Command{
	Use:   "update [maintenance-id]",
	Short: "Update a planned maintenance",
	Long:  "Update an existing planned maintenance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		maintenanceID := args[0]
		trace("Executing maintenance update command for ID: %s", maintenanceID)

		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")
		componentID, _ := cmd.Flags().GetString("component")
		startTime, _ := cmd.Flags().GetString("start-time")
		endTime, _ := cmd.Flags().GetString("end-time")

		// Build update object with only non-empty fields
		update := make(map[string]interface{})
		if title != "" {
			update["title"] = title
		}
		if description != "" {
			update["description"] = description
		}
		if status != "" {
			update["status"] = status
		}
		if componentID != "" {
			update["component_id"] = componentID
		}
		if startTime != "" {
			update["start_time"] = startTime
		}
		if endTime != "" {
			update["end_time"] = endTime
		}

		if len(update) == 0 {
			return fmt.Errorf("no fields to update")
		}

		client := getAPIClient()
		resp, err := client.makeRequest("PUT", "/api/v1/planned-maintenances/"+maintenanceID, update)
		if err != nil {
			return fmt.Errorf("failed to update maintenance: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("update maintenance failed: %s", string(body))
		}

		return outputData(body, "Maintenance Updated")
	},
}

// maintenanceStartCmd starts a planned maintenance (sets effective start time)
var maintenanceStartCmd = &cobra.Command{
	Use:   "start [maintenance-id]",
	Short: "Start a planned maintenance",
	Long:  "Start a planned maintenance by setting the effective start time to now",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		maintenanceID := args[0]
		trace("Executing maintenance start command for ID: %s", maintenanceID)

		client := getAPIClient()

		// First, get the existing maintenance data
		resp, err := client.makeRequest("GET", "/api/v1/planned-maintenances/"+maintenanceID, nil)
		if err != nil {
			return fmt.Errorf("failed to get maintenance: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("get maintenance failed: %s", string(body))
		}

		// Parse the existing maintenance data
		var maintenance map[string]interface{}
		if err := json.Unmarshal(body, &maintenance); err != nil {
			return fmt.Errorf("failed to parse maintenance data: %w", err)
		}

		// Convert components from objects to string array (server format difference)
		if components, ok := maintenance["components"].([]interface{}); ok {
			var componentCodes []string
			for _, comp := range components {
				if compMap, ok := comp.(map[string]interface{}); ok {
					if code, ok := compMap["code"].(string); ok {
						componentCodes = append(componentCodes, code)
					}
				}
			}
			maintenance["components"] = componentCodes
		}

		// Add start_effective to the existing data
		maintenance["start_effective"] = time.Now().Format(time.RFC3339)

		// Update the maintenance with all required fields
		resp, err = client.makeRequest("PUT", "/api/v1/planned-maintenances/"+maintenanceID, maintenance)
		if err != nil {
			return fmt.Errorf("failed to start maintenance: %w", err)
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("start maintenance failed: %s", string(body))
		}

		fmt.Printf("Maintenance %s started successfully\n", maintenanceID)
		return outputData(body, "Maintenance Started")
	},
}

// maintenanceStopCmd stops a planned maintenance (sets effective end time)
var maintenanceStopCmd = &cobra.Command{
	Use:   "stop [maintenance-id]",
	Short: "Stop a planned maintenance",
	Long:  "Stop a planned maintenance by setting the effective end time to now",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		maintenanceID := args[0]
		trace("Executing maintenance stop command for ID: %s", maintenanceID)

		// Update with current time as effective end
		update := map[string]interface{}{
			"end_effective": time.Now().Format(time.RFC3339),
		}

		client := getAPIClient()
		resp, err := client.makeRequest("PUT", "/api/v1/planned-maintenances/"+maintenanceID, update)
		if err != nil {
			return fmt.Errorf("failed to stop maintenance: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("stop maintenance failed: %s", string(body))
		}

		fmt.Printf("Maintenance %s stopped successfully\n", maintenanceID)
		return outputData(body, "Maintenance Stopped")
	},
}

// maintenanceDeleteCmd deletes a maintenance
var maintenanceDeleteCmd = &cobra.Command{
	Use:   "delete [maintenance-id]",
	Short: "Delete a planned maintenance",
	Long:  "Delete an existing planned maintenance by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		maintenanceID := args[0]
		trace("Executing maintenance delete command for ID: %s", maintenanceID)

		client := getAPIClient()
		resp, err := client.makeRequest("DELETE", "/api/v1/planned-maintenances/"+maintenanceID, nil)
		if err != nil {
			return fmt.Errorf("failed to delete maintenance: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 204 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("delete maintenance failed: %s", string(body))
		}

		fmt.Printf("Maintenance %s deleted successfully\n", maintenanceID)
		return nil
	},
}

func init() {
	// Add maintenance command to root
	rootCmd.AddCommand(maintenanceCmd)

	// Add subcommands to maintenance
	maintenanceCmd.AddCommand(maintenanceListCmd)
	maintenanceCmd.AddCommand(maintenanceGetCmd)
	maintenanceCmd.AddCommand(maintenanceCreateCmd)
	maintenanceCmd.AddCommand(maintenanceUpdateCmd)
	maintenanceCmd.AddCommand(maintenanceStartCmd)
	maintenanceCmd.AddCommand(maintenanceStopCmd)
	maintenanceCmd.AddCommand(maintenanceDeleteCmd)

	// Flags for create command
	maintenanceCreateCmd.Flags().String("title", "", "Maintenance title (required)")
	maintenanceCreateCmd.Flags().String("description", "", "Maintenance description")
	maintenanceCreateCmd.Flags().String("component", "", "Component ID")
	maintenanceCreateCmd.Flags().String("start-time", "", "Start time (RFC3339 format)")
	maintenanceCreateCmd.Flags().String("end-time", "", "End time (RFC3339 format)")
	maintenanceCreateCmd.MarkFlagRequired("title")

	// Flags for update command
	maintenanceUpdateCmd.Flags().String("title", "", "Maintenance title")
	maintenanceUpdateCmd.Flags().String("description", "", "Maintenance description")
	maintenanceUpdateCmd.Flags().String("status", "", "Maintenance status (scheduled|in-progress|completed|cancelled)")
	maintenanceUpdateCmd.Flags().String("component", "", "Component ID")
	maintenanceUpdateCmd.Flags().String("start-time", "", "Start time (RFC3339 format)")
	maintenanceUpdateCmd.Flags().String("end-time", "", "End time (RFC3339 format)")
}
