package issue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jasonwashburn/linear-cli/internal/linear"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	Long:  "List Linear issues, optionally filtered by team, state, or assignee.",
	RunE:  runList,
}

var (
	listTeam     string
	listState    string
	listAssignee string
	listJSON     bool
)

func init() {
	listCmd.Flags().StringVar(&listTeam, "team", "", "Filter by team ID or key")
	listCmd.Flags().StringVar(&listState, "state", "", "Filter by state name (e.g. In Progress)")
	listCmd.Flags().StringVar(&listAssignee, "assignee", "", "Filter by assignee ID")
	listCmd.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")
}

func runList(cmd *cobra.Command, args []string) error {
	apiKey := viper.GetString("LINEAR_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: LINEAR_API_KEY environment variable is not set")
		os.Exit(1)
	}
	client := linear.NewClient(apiKey)

	filter := linear.IssueFilter{
		TeamID:     listTeam,
		State:      listState,
		AssigneeID: listAssignee,
	}

	issues, err := client.ListIssues(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	if listJSON {
		return json.NewEncoder(os.Stdout).Encode(issues)
	}

	if len(issues) == 0 {
		fmt.Println("No issues found.")
		return nil
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Title", "State", "Priority", "Assignee", "Team"})
	for _, issue := range issues {
		assignee := ""
		if issue.Assignee != nil {
			assignee = issue.Assignee.DisplayName
		}
		t.AppendRow(table.Row{
			issue.Identifier,
			truncate(issue.Title, 50),
			issue.State.Name,
			linear.PriorityLabel(issue.Priority),
			assignee,
			issue.Team.Key,
		})
	}
	t.Render()
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
