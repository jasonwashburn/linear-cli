package issue

import (
	"github.com/spf13/cobra"
)

// Cmd is the parent "issue" subcommand.
var Cmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage Linear issues",
	Long:  "Create, read, update, and delete Linear issues.",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)
}
