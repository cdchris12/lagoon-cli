package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
)

var rshCmd = &cobra.Command{
	Use:   "rsh",
	Short: "Use remote shell access",
	Long:  "Access the remote shell for a project's environment.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !ValidateToken() {
			fmt.Println("Need to run `lagoon login` first")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("You must provide the project name and environment name.")
			os.Exit(1)
		}
		//proc := exec.Command("ssh", fmt.Sprintf("-p %s -t %s-%s@%s", viper.GetString("lagoon_port"), args[0], args[1], viper.GetString("lagoon_hostname")))
		cmdArgs := []string{
			fmt.Sprintf("\\-p %s", viper.GetString("lagoon_port")),
			fmt.Sprintf("\\-t %s-%s@%s", args[0], args[1], viper.GetString("lagoon_hostname")),
		}
		proc := exec.Command("ssh", cmdArgs...)
		var out bytes.Buffer
		proc.Stdout = &out
		err := proc.Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(rshCmd)
}
