package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// targetCmd represents the target command
var targetCmd = &cobra.Command{
	Use:   "target",
	Short: "Manage API targets (saved configurations)",
	Long:  "Manage API targets to save server URL and authentication credentials",
}

// targetSetCmd sets a new target
var targetSetCmd = &cobra.Command{
	Use:   "set [target-name]",
	Short: "Set a new target configuration",
	Long:  "Set a new target with server URL and authentication credentials",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetName := args[0]

		url, _ := cmd.Flags().GetString("url")
		user, _ := cmd.Flags().GetString("username")
		pass, _ := cmd.Flags().GetString("password")

		// Use current flags or environment if not provided
		if url == "" {
			url = serverURL
		}
		if user == "" {
			user = username
		}
		if pass == "" {
			pass = password
		}

		if url == "" {
			return fmt.Errorf("URL is required (use --url flag or CLARITI_SERVER_URL env var)")
		}

		err := setTarget(targetName, url, user, pass)
		if err != nil {
			return fmt.Errorf("failed to set target: %w", err)
		}

		fmt.Printf("Target '%s' set successfully\n", targetName)
		fmt.Printf("URL: %s\n", url)
		fmt.Printf("Username: %s\n", user)

		return nil
	},
}

// targetUseCmd switches to an existing target
var targetUseCmd = &cobra.Command{
	Use:   "use [target-name]",
	Short: "Switch to an existing target",
	Long:  "Switch to an existing target configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetName := args[0]

		err := switchTarget(targetName)
		if err != nil {
			return fmt.Errorf("failed to switch target: %w", err)
		}

		fmt.Printf("Switched to target '%s'\n", targetName)
		return nil
	},
}

// targetListCmd lists all targets
var targetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured targets",
	Long:  "List all configured targets with their details",
	RunE: func(cmd *cobra.Command, args []string) error {
		targets, current, err := listTargets()
		if err != nil {
			return fmt.Errorf("failed to list targets: %w", err)
		}

		if len(targets) == 0 {
			fmt.Println("No targets configured. Use 'clariti-cli target set' to create one.")
			return nil
		}

		fmt.Println("Configured targets:")
		fmt.Println()

		for name, target := range targets {
			marker := "  "
			if name == current {
				marker = "* "
			}

			fmt.Printf("%s%s\n", marker, name)
			fmt.Printf("    URL: %s\n", target.URL)
			fmt.Printf("    Username: %s\n", target.Username)
			fmt.Println()
		}

		if current != "" {
			fmt.Printf("Current target: %s\n", current)
		}

		return nil
	},
}

// targetDeleteCmd deletes a target
var targetDeleteCmd = &cobra.Command{
	Use:   "delete [target-name]",
	Short: "Delete a target configuration",
	Long:  "Delete a target configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetName := args[0]

		err := deleteTarget(targetName)
		if err != nil {
			return fmt.Errorf("failed to delete target: %w", err)
		}

		fmt.Printf("Target '%s' deleted successfully\n", targetName)
		return nil
	},
}

// targetCurrentCmd shows current target
var targetCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current target",
	Long:  "Show the currently active target configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		target, err := getCurrentTarget()
		if err != nil {
			fmt.Printf("No current target: %s\n", err.Error())
			return nil
		}

		fmt.Printf("Current target: %s\n", target.Name)
		fmt.Printf("URL: %s\n", target.URL)
		fmt.Printf("Username: %s\n", target.Username)

		return nil
	},
}

func init() {
	// Add target command to root
	rootCmd.AddCommand(targetCmd)

	// Add subcommands
	targetCmd.AddCommand(targetSetCmd)
	targetCmd.AddCommand(targetUseCmd)
	targetCmd.AddCommand(targetListCmd)
	targetCmd.AddCommand(targetDeleteCmd)
	targetCmd.AddCommand(targetCurrentCmd)

	// Flags for set command
	targetSetCmd.Flags().String("url", "", "Server URL")
	targetSetCmd.Flags().String("username", "", "Username for authentication")
	targetSetCmd.Flags().String("password", "", "Password for authentication")
}
