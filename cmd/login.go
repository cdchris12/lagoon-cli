package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"os"
	"strings"
	"net"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log into a Lagoon instance",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, _ := os.UserHomeDir()
		config := &ssh.ClientConfig{
			User: "lagoon",
			Auth: []ssh.AuthMethod{
				publicKey(fmt.Sprintf("%s/.ssh/id_rsa.pub", homeDir)),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		var err error

		conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", viper.GetString("lagoon_hostname"), viper.GetString("lagoon_port")), config)
		if err != nil {
			panic(err)
		}
		session, err := conn.NewSession()
		if err != nil {
			_ = conn.Close()
			panic(err)
		}

		out, err := session.CombinedOutput("token")
		if err != nil {
			panic(err)
		}
		err = conn.Close()
		viper.Set("lagoon_token", strings.TrimSpace(string(out)))
		err = viper.WriteConfig()
		if err != nil {
			panic(err)
		}
		fmt.Println("Token fetched and saved.")
	},
}

func publicKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// First, try to look for an unencrypted private key
	signer, err := ssh.ParsePrivateKey(key)
	if err.Error() != "ssh: cannot decode encrypted private keys" {
		panic(err)
	} else if err == nil {
		// return unencrypted private key
		return ssh.PublicKeys(signer)
	}

	//#TODO
	// Connect to SSH agent to ask for unencrypted private keys
	if sshAgentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		sshAgent := agent.NewClient(sshAgentConn)

		keys, _ := sshAgent.List()
		if len(keys) > 0 {
			// There are key(s) in the agent
			defer sshAgentConn.Close
			return ssh.PublicKeysCallback(sshAgent.Signers)
		}
	}

	// Handle encrypted private keys
	fmt.Println("Found an encrypted private key!")
	fmt.Printf("Enter passphrase for '%s': ", path)
	bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	signer, err = ssh.ParsePrivateKeyWithPassphrase(key, bytePassword)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}
