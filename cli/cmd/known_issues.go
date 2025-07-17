package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
)

// KnownIssue represents a known issue (different from incidents)
type KnownIssue struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
	ComponentID string `json:"component_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// knownIssueCmd represents the known-issue command
var knownIssueCmd = &cobra.Command{
	Use:     "known-issue",
	Short:   "Manage known issues",
	Long:    "Create, list, update and delete known issues (non-critical problems)",
	Aliases: []string{"ki", "issue"},
}

// knownIssueListCmd lists all known issues
var knownIssueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all known issues",
	Long:  "List all known issues from the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		trace("Executing known-issue list command")

		client := getAPIClient()
		resp, err := client.makeRequest("GET", "/api/v1/known-issues", nil)
		if err != nil {
			return fmt.Errorf("failed to list known issues: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("list known issues failed: %s", string(body))
		}

		return outputData(body, "Known Issues")
	},
}

// knownIssueGetCmd gets a specific known issue
var knownIssueGetCmd = &cobra.Command{
	Use:   "get [issue-id]",
	Short: "Get a specific known issue",
	Long:  "Get details of a specific known issue by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueID := args[0]
		trace("Executing known-issue get command for ID: %s", issueID)

		client := getAPIClient()
		resp, err := client.makeRequest("GET", "/api/v1/known-issues/"+issueID, nil)
		if err != nil {
			return fmt.Errorf("failed to get known issue: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("get known issue failed: %s", string(body))
		}

		return outputData(body, "Known Issue Details")
	},
}

// knownIssueCreateCmd creates a new known issue
var knownIssueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new known issue",
	Long:  "Create a new known issue with title, description, severity and component",
	RunE: func(cmd *cobra.Command, args []string) error {
		trace("Executing known-issue create command")

		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		severity, _ := cmd.Flags().GetString("severity")
		componentID, _ := cmd.Flags().GetString("component")

		if title == "" {
			return fmt.Errorf("title is required")
		}

		issue := map[string]interface{}{
			"title":           title,
			"content":         description,
			"start_effective": time.Now().Format(time.RFC3339), // Default to now for known issues
			"perpetual":       true,                            // Known issues are typically ongoing
			"criticality":     severity,                        // Use severity directly as criticality
		}

		// Add component if specified (server expects array of strings)
		if componentID != "" {
			issue["components"] = []string{componentID}
		} else {
			// Server validation requires at least one component
			issue["components"] = []string{"general"}
		}

		client := getAPIClient()
		resp, err := client.makeRequest("POST", "/api/v1/known-issues", issue)
		if err != nil {
			return fmt.Errorf("failed to create known issue: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 201 {
			return fmt.Errorf("create known issue failed: %s", string(body))
		}

		return outputData(body, "Known Issue Created")
	},
}

// knownIssueStopCmd stops/resolves a known issue (sets effective end time)
var knownIssueStopCmd = &cobra.Command{
	Use:   "stop [issue-id]",
	Short: "Stop/resolve a known issue",
	Long:  "Stop/resolve a known issue by setting the effective end time to now",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueID := args[0]
		trace("Executing known-issue stop command for ID: %s", issueID)

		// Update with current time as effective end
		update := map[string]interface{}{
			"end_effective": time.Now().Format(time.RFC3339),
		}

		client := getAPIClient()
		resp, err := client.makeRequest("PUT", "/api/v1/known-issues/"+issueID, update)
		if err != nil {
			return fmt.Errorf("failed to stop known issue: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("stop known issue failed: %s", string(body))
		}

		fmt.Printf("Known issue %s resolved successfully\n", issueID)
		return outputData(body, "Known Issue Resolved")
	},
}

func init() {
	// Add known-issue command to root
	rootCmd.AddCommand(knownIssueCmd)

	// Add subcommands to known-issue
	knownIssueCmd.AddCommand(knownIssueListCmd)
	knownIssueCmd.AddCommand(knownIssueGetCmd)
	knownIssueCmd.AddCommand(knownIssueCreateCmd)
	knownIssueCmd.AddCommand(knownIssueStopCmd)

	// Flags for create command
	knownIssueCreateCmd.Flags().String("title", "", "Known issue title (required)")
	knownIssueCreateCmd.Flags().String("description", "", "Known issue description")
	knownIssueCreateCmd.Flags().String("severity", "degraded", "Known issue criticality (operational|degraded|partial outage|major outage|under maintenance)")
	knownIssueCreateCmd.Flags().String("component", "", "Component ID")
	knownIssueCreateCmd.MarkFlagRequired("title")
}
