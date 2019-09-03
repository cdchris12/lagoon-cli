package cmd

import (
	"fmt"
	"os"

	"github.com/cdchris12/lagoon-cli/graphql"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show your projects",
	Run: func(cmd *cobra.Command, args []string) {
		var responseData WhatIsThere
		err := graphql.GraphQLRequest(`
query whatIsThere {
	allProjects {
		id
		gitUrl
		name,
		customer {
		  id,
		  name
		}
		environments {
		  environmentType,
		  route
		}
	  }
}
`, &responseData)
		if err != nil {
			panic(err)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoWrapText(true)
		table.SetHeader([]string{"ID", "Name", "Customer", "Git URL", "URL"})
		for _, project := range responseData.AllProjects {
			productionEnvironmentRoute, err := getProductionEnvironment(project.Environments, project.Name)
			if err != nil {
				panic(err)
			}
			table.Append([]string{
				fmt.Sprintf("%d", project.ID),
				project.Name,
				project.Customer.Name,
				project.GitURL,
				productionEnvironmentRoute,
			})
		}
		table.Render()
		fmt.Println()
		fmt.Println("To view a project's details, run `lagoon project info {name}`.")
	},
}

func init() {
	projectCmd.AddCommand(projectListCmd)
}

func getProductionEnvironment(environments []Environments, projectName string) (string, error) {
	for _, environment := range environments {
		if environment.EnvironmentType == "production" {
			return &environment.Route, nil
		}
	}
	//#TODO
	// Make this print in red, if possible
	if len(environments) == 0 {
		// No environments in this project
		fmt.Printf("The project %s has no environments! Skipping...", projectName)
		return "", nil
	} else {
		// Project has environments, but none are set as production
		fmt.Printf("No production environment could be found for the %s project! Defaulting to the first environment...", projectName)
		return &environments[0].Route, nil
	}
}
