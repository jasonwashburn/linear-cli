package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jasonwashburn/linear-cli/cmd/issue"
)

var rootCmd = &cobra.Command{
	Use:   "linear",
	Short: "A CLI for interacting with Linear",
	Long:  "linear-cli is a command-line tool for managing Linear issues from your terminal.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(issue.Cmd)
}

func initConfig() {
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
}
