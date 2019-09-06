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
		loginErr := loginToken()
		if loginErr != nil {
			panic(loginErr)
		}
		fmt.Println("Token fetched and saved.")
	},
}

func publicKey(path string) (ssh.AuthMethod, func() error) {
	noopCloseFunc := func() error { return nil }

	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// Connect to SSH agent to ask for unencrypted private keys
	if sshAgentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		sshAgent := agent.NewClient(sshAgentConn)

		keys, _ := sshAgent.List()
		if len(keys) > 0 {
			// There are key(s) in the agent
			//defer sshAgentConn.Close()
			return ssh.PublicKeysCallback(sshAgent.Signers), sshAgentConn.Close
		}
	}

	// Try to look for an unencrypted private key
	signer, err := ssh.ParsePrivateKey(key)
	if err.Error() != "ssh: cannot decode encrypted private keys" {
		panic(err)
	} else if err == nil {
	// return unencrypted private key
	return ssh.PublicKeys(signer), noopCloseFunc
	}

	// Handle encrypted private keys
	fmt.Println("Found an encrypted private key!")
	fmt.Printf("Enter passphrase for '%s': ", path)
	bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	signer, err = ssh.ParsePrivateKeyWithPassphrase(key, bytePassword)
	return ssh.PublicKeys(signer), noopCloseFunc
}

func loginToken() error {
	homeDir, _ := os.UserHomeDir()
	authMethod, closeSSHAgent := publicKey(fmt.Sprintf("%s/.ssh/id_rsa", homeDir))
	config := &ssh.ClientConfig{
		User: "lagoon",
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	defer closeSSHAgent()

	var err error

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", viper.GetString("lagoons."+cmdLagoon+".hostname"), viper.GetString("lagoons."+cmdLagoon+".port")), config)
	if err != nil {
		//panic(err)
		return err
	}
	session, err := conn.NewSession()
	if err != nil {
		_ = conn.Close()
		//panic(err)
		return err
	}

	out, err := session.CombinedOutput("token")
	if err != nil {
		//panic(err)
		return err
	}
	err = conn.Close()
	viper.Set("lagoons."+cmdLagoon+".token", strings.TrimSpace(string(out)))
	err = viper.WriteConfig()
	if err != nil {
		//panic(err)
		return err
	}
	return nil
}
