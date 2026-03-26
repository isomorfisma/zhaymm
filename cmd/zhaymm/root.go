package main

import (
	"github.com/spf13/cobra"
)

// rootCmd is the base command
var rootCmd = &cobra.Command{
	Use:   "zhaymm",
	Short: "zhaymm is a blazing-fast database seeder",
	Long:  `A zero-memory-bloat database seeder & anonymization pipeline powered by Go.`,
	// This below is what gets executed when user passes no sub-command
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to root command and set flags
func Execute() error {
	return rootCmd.Execute()
}