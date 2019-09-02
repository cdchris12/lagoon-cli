package cmd

import (
	"fmt"
	"os"

	"github.com/cdchris12/lagoon-cli/app"
	"github.com/cdchris12/lagoon-cli/graphql"

	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Show your projects, or details about a project",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !graphql.HasValidToken() {
			fmt.Println("Need to run `lagoon login` first")
			os.Exit(1)
		}
		cmdProject, _ = app.GetLocalProject()
	},
}
