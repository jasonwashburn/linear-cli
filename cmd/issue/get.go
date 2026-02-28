package issue

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <issue-id>",
		Short: "Get details for a Linear issue",
		Long:  `Display the full details of a single Linear issue by its ID.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			issueID := args[0]
			fmt.Fprintf(cmd.OutOrStdout(), "Would get issue %q\n", issueID)
			return nil
		},
	}

	return cmd
}
