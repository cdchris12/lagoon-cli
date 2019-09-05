package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configLagoonCmd = &cobra.Command{
	Use:   "lagoon",
	Short: "Add a lagoon to use",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Not enough arguments. Requires: lagoon name")
			os.Exit(1)
		}
		lagoonName := args[0]

		lagoonHostname := Prompt(fmt.Sprintf("Lagoon Hostname (%s)", viper.GetString("lagoons."+lagoonName+".hostname")))
		lagoonPort := Prompt(fmt.Sprintf("Lagoon Port (%d)", viper.GetInt("lagoons."+lagoonName+".port")))
		lagoonGraphQL := Prompt(fmt.Sprintf("Lagoon GraphQL endpoint (%s)", viper.GetString("lagoons."+lagoonName+".graphql")))

		viper.Set("lagoons."+lagoonName+".hostname", lagoonHostname)
		viper.Set("lagoons."+lagoonName+".port", lagoonPort)
		viper.Set("lagoons."+lagoonName+".graphql", lagoonGraphQL)
		err := viper.WriteConfig()
		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Sprintf("\nAdded a new lagoon named: %s", lagoonName))
	},
}

func init() {
	configCmd.AddCommand(configLagoonCmd)
}

var inputScanner = bufio.NewScanner(os.Stdin)

// GetInput reads input from an input buffer and returns the result as a string.
func GetInput() string {
	inputScanner.Scan()
	return strings.TrimSpace(inputScanner.Text())
}

// Prompt gets input with a prompt and returns the input
func Prompt(prompt string) string {
	fullPrompt := fmt.Sprintf("%s", prompt)
	fmt.Print(fullPrompt + ": ")
	return GetInput()
}
