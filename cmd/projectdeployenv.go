package cmd

import (
	"fmt"
	"os"

	"github.com/cdchris12/lagoon-cli/graphql"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var projectDeployEnvCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an environment",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			fmt.Println("Not enough arguments. Requires: project name and environment.")
			os.Exit(1)
		}
		projectName := args[0]
		projectEnvironment := args[1]

		// get a new token if the current one is invalid
		valid := graphql.VerifyTokenExpiry()
		if valid == false {
			loginErr := loginToken()
			if loginErr != nil {
				panic(loginErr)
			}
		}

		fmt.Println(fmt.Sprintf("Deploying %s-%s", projectName, projectEnvironment))

		if yesNo() {
			var responseData DeployResult
			err := graphql.GraphQLRequest(fmt.Sprintf(`mutation {
    deployEnvironmentBranch(
      input: {
        project:{name:"%s"}
        branchName:"%s"
      }
    )
  }`, projectName, projectEnvironment), &responseData)
			if err != nil {
				panic(err)
			}
			if responseData.DeployEnvironmentBranch == "success" {
				fmt.Println(fmt.Sprintf("Result: %s", aurora.Green(responseData.DeployEnvironmentBranch)))
			} else {
				fmt.Println(fmt.Sprintf("Result: %s", aurora.Yellow(responseData.DeployEnvironmentBranch)))
			}
		}

	},
}

func init() {
	projectCmd.AddCommand(projectDeployEnvCmd)
}
