package cmd

import (
	"fmt"
	"os"

	"github.com/cdchris12/lagoon-cli/graphql"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var projectDeleteEnvCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an environment",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			fmt.Println("Not enough arguments. Requires: project name and environment.")
			os.Exit(1)
		}
		projectName := args[0]
		projectEnvironment := args[1]

		fmt.Println(fmt.Sprintf("Deleting %s-%s", projectName, projectEnvironment))

		if yesNo() {
			var responseData DeleteResult
			err := graphql.GraphQLRequest(fmt.Sprintf(`mutation {
    deleteEnvironment(
      input: {
        project:"%s"
        name:"%s"
        execute:true
      }
    )
  }`, projectName, projectEnvironment), &responseData)
			if err != nil {
				panic(err)
			}
			if responseData.DeleteEnvironment == "success" {
				fmt.Println(fmt.Sprintf("Result: %s", aurora.Green(responseData.DeleteEnvironment)))
			} else {
				fmt.Println(fmt.Sprintf("Result: %s", aurora.Yellow(responseData.DeleteEnvironment)))
			}
		}

	},
}

func init() {
	projectCmd.AddCommand(projectDeleteEnvCmd)
}
