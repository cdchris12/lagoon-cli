package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var projectInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Details about a project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// @todo list current projects and allow choosing?
			fmt.Println("You must provide a project name.")
			os.Exit(1)
		}
		if len(args) > 1 {
			fmt.Println("Too many arguments.")
			os.Exit(1)
		}
		projectName := args[0]
		var responseData ProjectByName
		err := GraphQLRequest(fmt.Sprintf(`query {
  projectByName(name: "%s") {
    id,
    name,
    gitUrl,
    subfolder,
    branches,
    pullrequests,
    productionEnvironment,
    environments {
      name
      environmentType
      deployType
      route
    }
    autoIdle,
    storageCalc,
    developmentEnvironmentsLimit,
  }
}`, projectName), &responseData)

		if err != nil {
			panic(err)
		}
		project := responseData.ProjectByName
		var currentDevEnvironments int = 0
		for _, environment := range project.Environments {
			if environment.EnvironmentType == "development" {
				currentDevEnvironments++
			}
		}

		fmt.Println(projectName)
		fmt.Println()
		fmt.Println(fmt.Sprintf("Git URL: %s", project.GitURL))
		fmt.Println(fmt.Sprintf("Branches Pattern: %s", project.Branches))
		fmt.Println(fmt.Sprintf("Pull Requests: %s", project.Pullrequests))
		fmt.Println(fmt.Sprintf("Production Environment: %s", project.ProductionEnvironment))
		fmt.Println(fmt.Sprintf("Development Environments: %d / %d", currentDevEnvironments, project.DevelopmentEnvironmentsLimit))
		fmt.Println()
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoWrapText(true)
		table.SetHeader([]string{"Name", "Deploy Type", "Environment Type", "Route"})
		for _, environment := range project.Environments {
			table.Append([]string{
				environment.Name,
				environment.DeployType,
				environment.EnvironmentType,
				environment.Route,
			})
		}
		table.Render()
	},
}

func init() {
	projectCmd.AddCommand(projectInfoCmd)
}