package cmd

import (
	"fmt"
	"os"
	//"encoding/json"
	//"bytes"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a lagoon configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Not enough arguments. Requires: lagoon name")
			os.Exit(1)
		}
		lagoonName := args[0]
		fmt.Println(fmt.Sprintf("Deleting config for lagoon: %s", lagoonName))
		if yesNo() {
			err := Unset(lagoonName)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	configCmd.AddCommand(configDeleteCmd)
}

func Unset(key string) error {
	delete(viper.Get("lagoons").(map[string]interface{}), key)
	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}
