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
	Short: "Initialize a new ephemyral task",
	Long:  `This command initializes a new ephemyral task and either finds or creates a .ephemyral file.`,
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
