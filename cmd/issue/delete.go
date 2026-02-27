package issue

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jasonwashburn/linear-cli/internal/linear"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <issue-id>",
	Short: "Delete an issue",
	Long:  "Delete a Linear issue by its identifier. Prompts for confirmation unless --yes is passed.",
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

var deleteYes bool

func init() {
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation prompt")
}

func runDelete(cmd *cobra.Command, args []string) error {
	apiKey := viper.GetString("LINEAR_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: LINEAR_API_KEY environment variable is not set")
		os.Exit(1)
	}
	client := linear.NewClient(apiKey)

	issueID := args[0]

	if !deleteYes {
		fmt.Printf("Are you sure you want to delete issue %q? [y/N] ", issueID)
		reader := bufio.NewReader(os.Stdin)
		answer, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading confirmation: %w", err)
		}
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	if err := client.DeleteIssue(context.Background(), issueID); err != nil {
		return fmt.Errorf("error: %w", err)
	}

	fmt.Printf("Deleted issue %s.\n", issueID)
	return nil
}
