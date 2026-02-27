package issue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jasonwashburn/linear-cli/internal/linear"
)

var getCmd = &cobra.Command{
	Use:   "get <issue-id>",
	Short: "Show details for an issue",
	Long:  "Show full details for a single Linear issue by its identifier (e.g. ENG-123).",
	Args:  cobra.ExactArgs(1),
	RunE:  runGet,
}

var getJSON bool

func init() {
	getCmd.Flags().BoolVar(&getJSON, "json", false, "Output as JSON")
}

func runGet(cmd *cobra.Command, args []string) error {
	apiKey := viper.GetString("LINEAR_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: LINEAR_API_KEY environment variable is not set")
		os.Exit(1)
	}
	client := linear.NewClient(apiKey)

	issue, err := client.GetIssue(context.Background(), args[0])
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	if getJSON {
		return json.NewEncoder(os.Stdout).Encode(issue)
	}

	assignee := "Unassigned"
	if issue.Assignee != nil {
		assignee = issue.Assignee.DisplayName
	}

	fmt.Printf("Identifier: %s\n", issue.Identifier)
	fmt.Printf("Title:      %s\n", issue.Title)
	fmt.Printf("State:      %s\n", issue.State.Name)
	fmt.Printf("Priority:   %s\n", linear.PriorityLabel(issue.Priority))
	fmt.Printf("Assignee:   %s\n", assignee)
	fmt.Printf("Team:       %s (%s)\n", issue.Team.Name, issue.Team.Key)
	fmt.Printf("Created:    %s\n", issue.CreatedAt)
	fmt.Printf("Updated:    %s\n", issue.UpdatedAt)
	fmt.Printf("URL:        %s\n", issue.URL)
	if issue.Description != "" {
		fmt.Printf("\nDescription:\n%s\n", issue.Description)
	}
	return nil
}
