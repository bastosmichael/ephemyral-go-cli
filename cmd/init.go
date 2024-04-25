package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new .ephemyral file for AI generated configurations, setting up the environment in a directory.",
	Long:  `The 'init' command is designed to help set up the necessary configurations for an Ephemyral-based project. When executed, the command checks if there is an existing '.ephemyral' file in the current directory. This file contains YAML-formatted information related to AI-generated build, test, and lint commands.
If the '.ephemyral' file does not exist, the command creates one with a basic template for build, test, and lint command configurations. This is useful for initializing a new Ephemyral task or project where AI-driven commands can be defined and customized later.
If a '.ephemyral' file is already present, the command confirms that the Ephemyral task has been initialized, allowing users to proceed with other tasks such as building, testing, or linting.
The newly created '.ephemyral' file has a default structure with placeholders for build, test, and lint commands, which can be edited as needed. The 'init' command provides a foundation for AI-based project management, ensuring that an essential configuration file is in place before additional tasks are performed.` ,
	Run: func(cmd *cobra.Command, args []string) {
		filename := ".ephemyral"
		if !fileExists(filename) {
			fmt.Println("No .ephemyral file found, creating one...")
			createEphemyralFile(filename)
		} else {
			fmt.Println("Ephemyral task initialized, .ephemyral file found")
		}
	},
}

func createEphemyralFile(filename string) {
	content := struct {
		BuildCommand string `yaml:"build_command"`
		TestCommand  string `yaml:"test_command"`
		LintCommand  string `yaml:"lint_command"`
	}{}
	data, err := yaml.Marshal(&content)
	if err != nil {
		fmt.Printf("Error creating YAML content: %v\n", err)
		return
	}
	if err = os.WriteFile(filename, data, 0644); err != nil {
		fmt.Printf("Error writing .ephemyral file: %v\n", err)
		return
	}
	fmt.Println(".ephemyral file created")
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
