package cmd

import (
	"fmt"
	"os"
	"strings"

	//"github.com/cdchris12/lagoon-cli/graphql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configDefaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Set the default lagoon to use",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Not enough arguments. Requires: lagoon name")
			os.Exit(1)
		}
		lagoonName := args[0]
		viper.Set("default", strings.TrimSpace(string(lagoonName)))
		err := viper.WriteConfig()
		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Sprintf("Updating default lagoon to %s", lagoonName))
	},
}

func init() {
	configCmd.AddCommand(configDefaultCmd)
}
