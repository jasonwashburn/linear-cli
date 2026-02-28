package cmd

import (
	"os"

	"github.com/jasonwashburn/linear-cli/cmd/issue"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "linear",
	Short: "A CLI for managing Linear issues and resources",
	Long:  `linear-cli lets you manage Linear issues and resources directly from the terminal.`,
}

// Execute runs the root command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")

	rootCmd.AddCommand(issue.NewIssueCmd())
}

func initConfig() {
	viper.AutomaticEnv()
}
