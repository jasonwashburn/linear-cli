package issue

import (
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <issue-id>",
		Short: "Delete a Linear issue",
		Long:  `Delete a Linear issue by its ID.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			issueID := args[0]
			cmd.Printf("Would delete issue %q\n", issueID)
			return nil
		},
	}

	return cmd
}
