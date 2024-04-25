package cmd

import (
	"fmt"
	"strings"

	gpt4client "ephemyral/pkg"

	"github.com/spf13/cobra"
)

func generateTestCommand(directory string) (string, error) {
	filesList, err := getFileList(directory)
	if err != nil {
		return "", err
	}

	fullPrompt := TestCommandPrompt + strings.Join(filesList, "\n")

	gpt4client.SetDebug(false)
	return gpt4client.GetGPT4ResponseWithPrompt(fullPrompt)
}

var testCmd = &cobra.Command{
	Use:   "test [directory]",
	Short: "Generate and run a test command for the specified directory",
	Long:  "The 'test' command generates a testing command based on the structure of the project's files. It then updates the '.ephemyral' configuration file with the new testing command and executes it. Use this command to ensure your project adheres to testing standards and is free from test errors.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]

		if existingTestCommand, err := getExistingCommand(directory, "test"); err == nil && existingTestCommand != "" {
			fmt.Println("Running existing test command:", existingTestCommand)
			if err := executeCommand(directory, existingTestCommand); err != nil {
				fmt.Println("Error executing test command:", err)
			}
			return
		} else if err != nil {
			fmt.Println("Error reading .ephemyral file:", err)
			return
		}

		testCommand, err := generateTestCommand(directory)
		if err != nil {
			fmt.Println("Error generating test command:", err)
			return
		}

		refactoredTestCommand := filterOutCodeBlocks(testCommand)
		if err := updateEphemyralCommand(directory, "test", refactoredTestCommand); err != nil {
			fmt.Println("Error updating .ephemyral file:", err)
			return
		}

		fmt.Println("Successfully generated and updated test command:", refactoredTestCommand)
		if err := executeCommand(directory, refactoredTestCommand); err != nil {
			fmt.Println("Error executing test command:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
