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
		// get a new token if the current one is invalid
		valid := graphql.VerifyTokenExpiry()
		if valid == false {
			loginErr := loginToken()
			if loginErr != nil {
				panic(loginErr)
			}
		}

		var responseData WhatIsThere
		err := graphql.GraphQLRequest(`
query whatIsThere {
	allProjects {
		id
		gitUrl
		name,
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
		table.SetHeader([]string{"ID", "Project Name", "Git URL"})
		for _, project := range responseData.AllProjects {
			table.Append([]string{
				fmt.Sprintf("%d", project.ID),
				project.Name,
				project.GitURL,
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
