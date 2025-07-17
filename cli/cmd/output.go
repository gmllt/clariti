package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// OutputItem represents a generic item for display
type OutputItem map[string]interface{}

// displayPretty shows data in a user-friendly format
func displayPretty(items []OutputItem, title string) {
	if len(items) == 0 {
		fmt.Printf("No %s found.\n", strings.ToLower(title))
		return
	}

	// Show summary first
	fmt.Printf("Getting %s...\n", strings.ToLower(title))
	fmt.Printf("OK\n")
	fmt.Println()

	// For single items, show details
	if len(items) == 1 {
		item := items[0]
		displayItemPretty(item)
		return
	}

	// For multiple items, show list format
	displayListPretty(items, title)
}

// displayItemPretty shows a single item in detail
func displayItemPretty(item OutputItem) {
	// Show key details first
	if title, ok := item["title"].(string); ok && title != "" {
		fmt.Printf("name:          %s\n", title)
	}
	if guid, ok := item["guid"].(string); ok && guid != "" {
		fmt.Printf("guid:          %s\n", guid)
	}
	if status := getStatus(item); status != "" {
		fmt.Printf("status:        %s\n", status)
	}
	if criticality := getCriticality(item); criticality != "" {
		fmt.Printf("criticality:   %s\n", criticality)
	}

	// Show content/description
	if content, ok := item["content"].(string); ok && content != "" {
		fmt.Printf("description:   %s\n", content)
	}

	// Show components
	if components := getComponents(item); len(components) > 0 {
		fmt.Printf("components:    %s\n", strings.Join(components, ", "))
	}

	// Show timing information
	if startPlanned := getTime(item, "start_planned"); startPlanned != "" {
		fmt.Printf("planned start: %s\n", startPlanned)
	}
	if endPlanned := getTime(item, "end_planned"); endPlanned != "" {
		fmt.Printf("planned end:   %s\n", endPlanned)
	}
	if startEffective := getTime(item, "start_effective"); startEffective != "" {
		fmt.Printf("actual start:  %s\n", startEffective)
	}
	if endEffective := getTime(item, "end_effective"); endEffective != "" {
		fmt.Printf("actual end:    %s\n", endEffective)
	}
}

// displayListPretty shows multiple items in a compact list
func displayListPretty(items []OutputItem, title string) {
	fmt.Printf("%s (%d):\n", title, len(items))

	// Find max widths for alignment
	maxNameWidth := 4         // minimum for "name"
	maxStatusWidth := 6       // minimum for "status"
	maxCriticalityWidth := 11 // minimum for "criticality"

	for _, item := range items {
		if title, ok := item["title"].(string); ok {
			if len(title) > maxNameWidth {
				maxNameWidth = len(title)
			}
		}
		if status := getStatus(item); len(status) > maxStatusWidth {
			maxStatusWidth = len(status)
		}
		if criticality := getCriticality(item); len(criticality) > maxCriticalityWidth {
			maxCriticalityWidth = len(criticality)
		}
	}

	// Header
	fmt.Printf("%-*s %-*s %-*s %s\n",
		maxNameWidth, "name",
		maxStatusWidth, "status",
		maxCriticalityWidth, "criticality",
		"guid")

	// Items
	for _, item := range items {
		name := getString(item, "title", "-")
		status := getStatus(item)
		if status == "" {
			status = "-"
		}
		criticality := getCriticality(item)
		if criticality == "" {
			criticality = "-"
		}
		guid := getString(item, "guid", "-")

		fmt.Printf("%-*s %-*s %-*s %s\n",
			maxNameWidth, name,
			maxStatusWidth, status,
			maxCriticalityWidth, criticality,
			guid)
	}
}

// Helper functions
func getString(item OutputItem, key, defaultValue string) string {
	if value, ok := item[key].(string); ok {
		return value
	}
	return defaultValue
}

func getStatus(item OutputItem) string {
	// Try different status fields
	if value, ok := item["status"].(string); ok && value != "" {
		return value
	}

	// For incidents/maintenance, calculate status from timing
	now := time.Now()

	startEffective := getTimeObj(item, "start_effective")
	endEffective := getTimeObj(item, "end_effective")
	startPlanned := getTimeObj(item, "start_planned")

	// Check if it's ongoing
	if !startEffective.IsZero() && endEffective.IsZero() {
		return "ongoing"
	}

	// Check if it's resolved
	if !endEffective.IsZero() {
		return "resolved"
	}

	// Check if it's planned
	if !startPlanned.IsZero() && startPlanned.After(now) {
		return "planned"
	}

	return "unknown"
}

func getCriticality(item OutputItem) string {
	if value, ok := item["criticality"].(string); ok {
		return value
	}
	if value, ok := item["criticality"].(float64); ok {
		// Convert numeric criticality to string
		criticalityMap := map[int]string{
			0: "operational",
			1: "degraded",
			2: "partial outage",
			3: "major outage",
			4: "maintenance",
		}
		if str, exists := criticalityMap[int(value)]; exists {
			return str
		}
	}
	return ""
}

func getComponents(item OutputItem) []string {
	if components, ok := item["components"].([]interface{}); ok {
		var result []string
		for _, comp := range components {
			if compMap, ok := comp.(map[string]interface{}); ok {
				if code, ok := compMap["code"].(string); ok {
					result = append(result, code)
				} else if name, ok := compMap["name"].(string); ok {
					result = append(result, name)
				}
			}
		}
		return result
	}
	return nil
}

func getTime(item OutputItem, key string) string {
	timeObj := getTimeObj(item, key)
	if timeObj.IsZero() {
		return ""
	}
	return timeObj.Format("2006-01-02 15:04:05 MST")
}

func getTimeObj(item OutputItem, key string) time.Time {
	if timeStr, ok := item[key].(string); ok && timeStr != "" {
		if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
			return t
		}
	}
	return time.Time{}
}

// outputDataPretty handles user-friendly output
func outputDataPretty(data []byte, title string) error {
	// Try to parse as single object first
	var item OutputItem
	if err := json.Unmarshal(data, &item); err == nil {
		// Special handling for weather data
		if _, hasOverall := item["overall"]; hasOverall {
			displayWeatherPretty(item, title)
			return nil
		}

		displayPretty([]OutputItem{item}, title)
		return nil
	}

	// Try to parse as array
	var items []OutputItem
	if err := json.Unmarshal(data, &items); err == nil {
		displayPretty(items, title)
		return nil
	}

	// Fallback to JSON output if parsing fails
	fmt.Printf("Response: %s\n", string(data))
	return nil
} // displayWeatherPretty shows weather data in a comprehensive format
func displayWeatherPretty(item OutputItem, title string) {
	fmt.Printf("Getting %s...\n", strings.ToLower(title))
	fmt.Printf("OK\n")
	fmt.Println()

	// Show overall status first
	if overall, ok := item["overall"].(map[string]interface{}); ok {
		status := getString(OutputItem(overall), "status_label", "unknown")
		lastUpdated := getTime(OutputItem(overall), "last_updated")

		fmt.Printf("Overall Status: %s\n", status)
		if lastUpdated != "" {
			fmt.Printf("Last Updated:   %s\n", lastUpdated)
		}
		fmt.Println()
	}

	// Show platforms
	if platforms, ok := item["platforms"].([]interface{}); ok && len(platforms) > 0 {
		fmt.Printf("Platforms:\n")
		for _, p := range platforms {
			if platform, ok := p.(map[string]interface{}); ok {
				name := getString(OutputItem(platform), "platform", "")
				status := getString(OutputItem(platform), "status_label", "unknown")
				fmt.Printf("  %-20s %s\n", name, status)
			}
		}
		fmt.Println()
	}

	// Show instances
	if instances, ok := item["instances"].([]interface{}); ok && len(instances) > 0 {
		fmt.Printf("Instances:\n")
		for _, i := range instances {
			if instance, ok := i.(map[string]interface{}); ok {
				name := getString(OutputItem(instance), "instance", "")
				platform := getString(OutputItem(instance), "platform", "")
				status := getString(OutputItem(instance), "status_label", "unknown")
				fmt.Printf("  %-20s %-20s %s\n", name, platform, status)
			}
		}
		fmt.Println()
	}

	// Show components
	if components, ok := item["components"].([]interface{}); ok && len(components) > 0 {
		fmt.Printf("Components:\n")
		for _, c := range components {
			if component, ok := c.(map[string]interface{}); ok {
				name := getString(OutputItem(component), "component", "")
				instance := getString(OutputItem(component), "instance", "")
				status := getString(OutputItem(component), "status_label", "unknown")
				fmt.Printf("  %-25s %-20s %s\n", name, instance, status)
			}
		}
	}
}
