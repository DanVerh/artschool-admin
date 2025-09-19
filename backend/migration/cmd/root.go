package cmd

import (
    "github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Migration CLI for Artschool Admin",
    Long:  "CLI tool to run and rollback MongoDB migrations for Artschool Admin",
}

// Execute runs the root command
func Execute() error {
    return RootCmd.Execute()
}
