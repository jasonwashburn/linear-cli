package issue

import (
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var team, state, assignee string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Linear issues",
		Long:  `List Linear issues, optionally filtered by team, state, or assignee.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("Would list issues (team=%q, state=%q, assignee=%q)\n", team, state, assignee)
			return nil
		},
	}

	cmd.Flags().StringVar(&team, "team", "", "Filter by team name or key")
	cmd.Flags().StringVar(&state, "state", "", "Filter by issue state (e.g. Todo, In Progress, Done)")
	cmd.Flags().StringVar(&assignee, "assignee", "", "Filter by assignee username or email")

	return cmd
}
