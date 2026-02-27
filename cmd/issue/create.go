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

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	Long:  "Create a new Linear issue with the provided title and optional fields.",
	RunE:  runCreate,
}

var (
	createTitle       string
	createTeam        string
	createDescription string
	createPriority    int
	createAssignee    string
	createParent      string
	createJSON        bool
)

func init() {
	createCmd.Flags().StringVar(&createTitle, "title", "", "Issue title (required)")
	createCmd.Flags().StringVar(&createTeam, "team", "", "Team ID (required)")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Issue description")
	createCmd.Flags().IntVar(&createPriority, "priority", -1, "Priority: 0=None, 1=Urgent, 2=High, 3=Medium, 4=Low")
	createCmd.Flags().StringVar(&createAssignee, "assignee", "", "Assignee user ID")
	createCmd.Flags().StringVar(&createParent, "parent", "", "Parent issue ID (creates a sub-issue)")
	createCmd.Flags().BoolVar(&createJSON, "json", false, "Output as JSON")
	_ = createCmd.MarkFlagRequired("title")
	_ = createCmd.MarkFlagRequired("team")
}

func runCreate(cmd *cobra.Command, args []string) error {
	apiKey := viper.GetString("LINEAR_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: LINEAR_API_KEY environment variable is not set")
		os.Exit(1)
	}
	client := linear.NewClient(apiKey)

	input := linear.IssueCreateInput{
		Title:       createTitle,
		TeamID:      createTeam,
		Description: createDescription,
		AssigneeID:  createAssignee,
		ParentID:    createParent,
	}
	if createPriority >= 0 {
		p := createPriority
		input.Priority = &p
	}

	issue, err := client.CreateIssue(context.Background(), input)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	if createJSON {
		return json.NewEncoder(os.Stdout).Encode(issue)
	}

	fmt.Printf("Created issue %s: %s\n", issue.Identifier, issue.Title)
	fmt.Printf("URL: %s\n", issue.URL)
	return nil
}
