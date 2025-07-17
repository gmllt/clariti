package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
)

// Incident represents an incident (using the same structure as the server)
type Incident struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
	ComponentID string `json:"component_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// incidentCmd represents the incident command
var incidentCmd = &cobra.Command{
	Use:   "incident",
	Short: "Manage critical incidents",
	Long:  "Create, list, update and delete critical incidents (active outages and major problems)",
}

// incidentListCmd lists all incidents
var incidentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all incidents",
	Long:  "List all incidents from the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		trace("Executing incident list command")

		client := getAPIClient()
		resp, err := client.makeRequest("GET", "/api/v1/incidents", nil)
		if err != nil {
			return fmt.Errorf("failed to list incidents: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("list incidents failed: %s", string(body))
		}

		return outputData(body, "Incidents")
	},
}

// incidentGetCmd gets a specific incident
var incidentGetCmd = &cobra.Command{
	Use:   "get [incident-id]",
	Short: "Get a specific incident",
	Long:  "Get details of a specific incident by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		incidentID := args[0]
		trace("Executing incident get command for ID: %s", incidentID)

		client := getAPIClient()
		resp, err := client.makeRequest("GET", "/api/v1/incidents/"+incidentID, nil)
		if err != nil {
			return fmt.Errorf("failed to get incident: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("get incident failed: %s", string(body))
		}

		return outputData(body, "Incident Details")
	},
}

// incidentCreateCmd creates a new incident
var incidentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new critical incident",
	Long:  "Create a new critical incident for active outages and major problems",
	RunE: func(cmd *cobra.Command, args []string) error {
		trace("Executing incident create command")

		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		severity, _ := cmd.Flags().GetString("severity")
		componentID, _ := cmd.Flags().GetString("component")

		if title == "" {
			return fmt.Errorf("title is required")
		}

		incident := map[string]interface{}{
			"title":           title,
			"content":         description,
			"start_effective": time.Now().Format(time.RFC3339), // Default to now for incidents
			"criticality":     severity,                        // Use severity directly as criticality
		}

		// Add component if specified (server expects array of strings)
		if componentID != "" {
			incident["components"] = []string{componentID}
		} else {
			// Server validation requires at least one component
			incident["components"] = []string{"general"}
		}

		client := getAPIClient()
		resp, err := client.makeRequest("POST", "/api/v1/incidents", incident)
		if err != nil {
			return fmt.Errorf("failed to create incident: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 201 {
			return fmt.Errorf("create incident failed: %s", string(body))
		}

		return outputData(body, "Incident Created")
	},
}

// incidentUpdateCmd updates an incident
var incidentUpdateCmd = &cobra.Command{
	Use:   "update [incident-id]",
	Short: "Update an incident",
	Long:  "Update an existing incident",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		incidentID := args[0]
		trace("Executing incident update command for ID: %s", incidentID)

		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		severity, _ := cmd.Flags().GetString("severity")
		status, _ := cmd.Flags().GetString("status")
		componentID, _ := cmd.Flags().GetString("component")

		// Build update object with only non-empty fields
		update := make(map[string]interface{})
		if title != "" {
			update["title"] = title
		}
		if description != "" {
			update["description"] = description
		}
		if severity != "" {
			update["severity"] = severity
		}
		if status != "" {
			update["status"] = status
		}
		if componentID != "" {
			update["component_id"] = componentID
		}

		if len(update) == 0 {
			return fmt.Errorf("no fields to update")
		}

		client := getAPIClient()
		resp, err := client.makeRequest("PUT", "/api/v1/incidents/"+incidentID, update)
		if err != nil {
			return fmt.Errorf("failed to update incident: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("update incident failed: %s", string(body))
		}

		return outputData(body, "Incident Updated")
	},
}

// incidentDeleteCmd deletes an incident
var incidentDeleteCmd = &cobra.Command{
	Use:   "delete [incident-id]",
	Short: "Delete an incident",
	Long:  "Delete an existing incident by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		incidentID := args[0]
		trace("Executing incident delete command for ID: %s", incidentID)

		client := getAPIClient()
		resp, err := client.makeRequest("DELETE", "/api/v1/incidents/"+incidentID, nil)
		if err != nil {
			return fmt.Errorf("failed to delete incident: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 204 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("delete incident failed: %s", string(body))
		}

		fmt.Printf("Incident %s deleted successfully\n", incidentID)
		return nil
	},
}

// incidentStopCmd stops/resolves an incident (sets effective end time)
var incidentStopCmd = &cobra.Command{
	Use:   "stop [incident-id]",
	Short: "Stop/resolve a critical incident",
	Long:  "Stop/resolve a critical incident by setting the effective end time to now",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		incidentID := args[0]
		trace("Executing incident stop command for ID: %s", incidentID)

		client := getAPIClient()

		// First, get the current incident data
		resp, err := client.makeRequest("GET", "/api/v1/incidents/"+incidentID, nil)
		if err != nil {
			return fmt.Errorf("failed to get incident: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("get incident failed: %s", string(body))
		}

		// Parse the existing incident data
		var incident map[string]interface{}
		if err := json.Unmarshal(body, &incident); err != nil {
			return fmt.Errorf("failed to parse incident data: %w", err)
		}

		// Convert components from objects to string array (server format difference)
		if components, ok := incident["components"].([]interface{}); ok {
			var componentCodes []string
			for _, comp := range components {
				if compMap, ok := comp.(map[string]interface{}); ok {
					if code, ok := compMap["code"].(string); ok {
						componentCodes = append(componentCodes, code)
					}
				}
			}
			incident["components"] = componentCodes
		}

		// Convert criticality from number to string (server format difference)
		if criticality, ok := incident["criticality"].(float64); ok {
			criticalityMap := map[int]string{
				0: "operational",
				1: "degraded",
				2: "partial outage",
				3: "major outage",
				4: "under maintenance",
			}
			if critStr, exists := criticalityMap[int(criticality)]; exists {
				incident["criticality"] = critStr
			}
		}

		// Add end_effective to the existing data
		incident["end_effective"] = time.Now().Format(time.RFC3339)

		// Update the incident with all required fields
		resp, err = client.makeRequest("PUT", "/api/v1/incidents/"+incidentID, incident)
		if err != nil {
			return fmt.Errorf("failed to stop incident: %w", err)
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("stop incident failed: %s", string(body))
		}

		fmt.Printf("Incident %s resolved successfully\n", incidentID)
		return outputData(body, "Incident Resolved")
	},
}

func init() {
	// Add incident command to root
	rootCmd.AddCommand(incidentCmd)

	// Add subcommands to incident
	incidentCmd.AddCommand(incidentListCmd)
	incidentCmd.AddCommand(incidentGetCmd)
	incidentCmd.AddCommand(incidentCreateCmd)
	incidentCmd.AddCommand(incidentUpdateCmd)
	incidentCmd.AddCommand(incidentStopCmd)
	incidentCmd.AddCommand(incidentDeleteCmd)

	// Flags for create command
	incidentCreateCmd.Flags().String("title", "", "Incident title (required)")
	incidentCreateCmd.Flags().String("description", "", "Incident description")
	incidentCreateCmd.Flags().String("severity", "partial outage", "Incident criticality (operational|degraded|partial outage|major outage|under maintenance)")
	incidentCreateCmd.Flags().String("component", "", "Component ID")
	incidentCreateCmd.MarkFlagRequired("title")

	// Flags for update command
	incidentUpdateCmd.Flags().String("title", "", "Incident title")
	incidentUpdateCmd.Flags().String("description", "", "Incident description")
	incidentUpdateCmd.Flags().String("severity", "", "Incident criticality (operational|degraded|partial outage|major outage|under maintenance)")
	incidentUpdateCmd.Flags().String("status", "", "Incident status (investigating|identified|monitoring|resolved)")
	incidentUpdateCmd.Flags().String("component", "", "Component ID")
}
