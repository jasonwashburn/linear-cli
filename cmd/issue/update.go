package issue

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	var title, description, state, assignee string
	var priority int

	cmd := &cobra.Command{
		Use:   "update <issue-id>",
		Short: "Update a Linear issue",
		Long:  `Update one or more fields on an existing Linear issue.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			issueID := args[0]
			fmt.Fprintf(cmd.OutOrStdout(),
				"Would update issue %q (title=%q, description=%q, state=%q, priority=%d, assignee=%q)\n",
				issueID, title, description, state, priority, assignee,
			)
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "New issue title")
	cmd.Flags().StringVar(&description, "description", "", "New issue description")
	cmd.Flags().StringVar(&state, "state", "", "New issue state (e.g. Todo, In Progress, Done)")
	cmd.Flags().IntVar(&priority, "priority", 0, "Priority level: 0 (none), 1 (urgent), 2 (high), 3 (medium), 4 (low)")
	cmd.Flags().StringVar(&assignee, "assignee", "", "New assignee username or email")

	return cmd
}
