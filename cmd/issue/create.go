package issue

import (
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var title, team, description, assignee string
	var priority int

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Linear issue",
		Long:  `Create a new Linear issue with the given title and optional fields.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf(
				"Would create issue (title=%q, team=%q, description=%q, priority=%d, assignee=%q)\n",
				title, team, description, priority, assignee,
			)
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Issue title (required)")
	cmd.Flags().StringVar(&team, "team", "", "Team name or key")
	cmd.Flags().StringVar(&description, "description", "", "Issue description")
	cmd.Flags().IntVar(&priority, "priority", 0, "Priority level: 0 (none), 1 (urgent), 2 (high), 3 (medium), 4 (low)")
	cmd.Flags().StringVar(&assignee, "assignee", "", "Assignee username or email")

	if err := cmd.MarkFlagRequired("title"); err != nil {
		panic(err)
	}

	return cmd
}
