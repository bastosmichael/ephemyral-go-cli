// init.go
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new ephemyral task",
	Long:  `This command initializes a new ephemyral task and either finds or creates a .ephemyral file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Define the file path
		filename := ".ephemyral"

		// Check if the file exists
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			// File does not exist, create it
			fmt.Println("No .ephemyral file found, creating one...")

			// Define the default content as a YAML structure
			content := struct {
				BuildCommand string `yaml:"build_command"`
				TestCommand  string `yaml:"test_command"`
			}{
				BuildCommand: "", // Placeholder for build command
				TestCommand:  "", // Placeholder for test command
			}

			// Marshal the content into YAML format
			data, err := yaml.Marshal(&content)
			if err != nil {
				fmt.Println("Error creating YAML content:", err)
				return
			}

			// Write the YAML data to the file
			if err := ioutil.WriteFile(filename, data, 0644); err != nil {
				fmt.Println("Error writing .ephemyral file:", err)
				return
			}

			fmt.Println(".ephemyral file created")
		} else {
			// File exists
			fmt.Println("Ephemyral task initialized, .ephemyral file found")
		}
	},
}
