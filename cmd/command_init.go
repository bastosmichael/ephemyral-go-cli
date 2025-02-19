// +build !lint

package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/sha3"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new .ephemyral file for AI-generated configurations, setting up the environment in a directory.",
	Long: `The 'init' command is designed to help set up the necessary configurations for an Ephemyral-based project. When executed, the command checks if there is an existing '.ephemyral' file in the current directory. This file contains YAML-formatted information related to AI-generated build, test, and lint commands.
If the '.ephemyral' file does not exist, the command creates one with a basic template for build, test, and lint command configurations. This is useful for initializing a new Ephemyral task or project where AI-driven commands can be defined and customized later.
If a '.ephemyral' file is already present, the command confirms that the Ephemyral task has been initialized, allowing users to proceed with other tasks such as building, testing, or linting.
The newly created '.ephemyral' file has a default structure with placeholders for build, test, and lint commands, which can be edited as needed. The 'init' command provides a foundation for AI-based project management, ensuring that an essential configuration file is in place before additional tasks are performed.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename := ".ephemyral"
		if !fileExists(filename) {
			fmt.Println("No .ephemyral file found, creating one...")
			createEphemyralFile(filename)
		} else {
			fmt.Println("Ephemyral task initialized, .ephemyral file found")
			checkAndDecryptAPIKey(filename)
		}
	},
}

func createEphemyralFile(filename string) {
	reader := bufio.NewReader(os.Stdin)

	// Ask for the OpenAI API key
	fmt.Print("Enter your OpenAI API key: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	// Ask if the user wants to encrypt the API key
	fmt.Print("Would you like to encrypt your API key with a passphrase? (yes/no): ")
	encryptOption, _ := reader.ReadString('\n')
	encryptOption = strings.TrimSpace(strings.ToLower(encryptOption))

	if encryptOption == "yes" {
		fmt.Print("Enter a passphrase for encryption: ")
		passphrase, _ := reader.ReadString('\n')
		passphrase = strings.TrimSpace(passphrase)
		encryptedAPIKey, err := encrypt(apiKey, passphrase)
		if err != nil {
			fmt.Printf("Error encrypting API key: %v\n", err)
			return
		}
		apiKey = encryptedAPIKey
	}

	// Use the EphemyralFile struct to create the default content
	content := EphemyralFile{
		OpenAIAPIKey: apiKey,
		BuildCommand: "",
		TestCommand:  "",
		LintCommand:  "",
		DocsCommand:  "",
	}

	// Marshal the content into YAML format
	data, err := yaml.Marshal(&content)
	if err != nil {
		fmt.Printf("Error creating YAML content: %v\n", err)
		return
	}

	// Write the YAML data to the .ephemyral file
	if err = os.WriteFile(filename, data, 0644); err != nil {
		fmt.Printf("Error writing .ephemyral file: %v\n", err)
		return
	}
	fmt.Println(".ephemyral file created")
}

func checkAndDecryptAPIKey(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading .ephemyral file: %v\n", err)
		return
	}

	var content EphemyralFile
	if err := yaml.Unmarshal(data, &content); err != nil {
		fmt.Printf("Error parsing .ephemyral file: %v\n", err)
		return
	}

	apiKey := content.OpenAIAPIKey
	reader := bufio.NewReader(os.Stdin)

	// Check if the API key looks encrypted (base64 encoded string)
	if looksLikeEncrypted(apiKey) {
		fmt.Print("The API key appears to be encrypted. Would you like to decrypt it and display it? (yes/no): ")
		decryptOption, _ := reader.ReadString('\n')
		decryptOption = strings.TrimSpace(strings.ToLower(decryptOption))

		if decryptOption == "yes" {
			fmt.Print("Enter the passphrase for decryption: ")
			passphrase, _ := reader.ReadString('\n')
			passphrase = strings.TrimSpace(passphrase)
			decryptedAPIKey, err := decrypt(apiKey, passphrase)
			if err != nil {
				fmt.Printf("Error decrypting API key: %v\n", err)
				return
			}
			fmt.Printf("Decrypted API key: %s\n", decryptedAPIKey)
		}
	}
}

func createHash(key string) string {
	hash := sha3.New256()
	hash.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))[:32]
}

func looksLikeEncrypted(text string) bool {
	_, err := base64.StdEncoding.DecodeString(text)
	return err == nil
}
