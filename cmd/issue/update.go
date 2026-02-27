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

var updateCmd = &cobra.Command{
	Use:   "update <issue-id>",
	Short: "Update an existing issue",
	Long:  "Update writable fields on an existing Linear issue by its identifier.",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdate,
}

var (
	updateTitle       string
	updateDescription string
	updateState       string
	updatePriority    int
	updateAssignee    string
	updateJSON        bool
)

func init() {
	updateCmd.Flags().StringVar(&updateTitle, "title", "", "New title")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "New description")
	updateCmd.Flags().StringVar(&updateState, "state", "", "New state ID")
	updateCmd.Flags().IntVar(&updatePriority, "priority", -1, "New priority: 0=None, 1=Urgent, 2=High, 3=Medium, 4=Low")
	updateCmd.Flags().StringVar(&updateAssignee, "assignee", "", "New assignee user ID")
	updateCmd.Flags().BoolVar(&updateJSON, "json", false, "Output as JSON")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	apiKey := viper.GetString("LINEAR_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: LINEAR_API_KEY environment variable is not set")
		os.Exit(1)
	}
	client := linear.NewClient(apiKey)

	input := linear.IssueUpdateInput{}
	if cmd.Flags().Changed("title") {
		input.Title = &updateTitle
	}
	if cmd.Flags().Changed("description") {
		input.Description = &updateDescription
	}
	if cmd.Flags().Changed("state") {
		input.StateID = &updateState
	}
	if cmd.Flags().Changed("priority") {
		input.Priority = &updatePriority
	}
	if cmd.Flags().Changed("assignee") {
		input.AssigneeID = &updateAssignee
	}

	issue, err := client.UpdateIssue(context.Background(), args[0], input)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	if updateJSON {
		return json.NewEncoder(os.Stdout).Encode(issue)
	}

	fmt.Printf("Updated issue %s: %s\n", issue.Identifier, issue.Title)
	fmt.Printf("State: %s | Priority: %s\n", issue.State.Name, linear.PriorityLabel(issue.Priority))
	return nil
}
