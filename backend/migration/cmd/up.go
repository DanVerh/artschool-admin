package cmd

import (
    "log"
    "github.com/DanVerh/artschool-admin/backend/migration/internal/application"
    "github.com/spf13/cobra"
)

func UpCommand(app *application.App) *cobra.Command {
    return &cobra.Command{
        Use:   "up",
        Short: "Apply all up migrations",
        Run: func(cmd *cobra.Command, args []string) {
            if err := app.Up(); err != nil {
                log.Fatal(err)
            }
        },
    }
}
