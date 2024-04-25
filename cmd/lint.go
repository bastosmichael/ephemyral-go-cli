package cmd

import (
	"fmt"
	"strings"

	gpt4client "ephemyral/pkg"
	"github.com/spf13/cobra"
)

func generateLintCommand(directory string) (string, error) {
	filesList, err := getFileList(directory)
	if err != nil {
		return "", err
	}

	fullPrompt := LintCommandPrompt + strings.Join(filesList, "\n")
	gpt4client.SetDebug(false)
	lintCommand, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(lintCommand) == "" {
		return "", fmt.Errorf("received empty lint command")
	}

	return lintCommand, nil
}

var lintCmd = &cobra.Command{
	Use:   "lint [directory]",
	Short: "Use a machine learning model to generate and execute a linting command, improving code quality by identifying patterns and anomalies.",
	Long:  "The 'lint' command generates a linting command based on the structure of the project's files. It then updates the '.ephemyral' configuration file with the new linting command and executes it. Use this command to ensure your project adheres to coding standards and is free from basic syntax errors.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]
		existingLintCommand, err := getExistingCommand(directory, "lint")
		if err != nil {
			fmt.Printf("Error reading .ephemyral file: %v\n", err)
			return
		}
		if existingLintCommand != "" {
			fmt.Printf("Running existing lint command: %s\n", existingLintCommand)
			if err := executeCommand(directory, existingLintCommand); err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}
			return
		}

		lintCommand, err := generateLintCommand(directory)
		if err != nil {
			fmt.Printf("Error generating lint command: %v\n", err)
			return
		}

		refactoredLintCommand := filterOutCodeBlocks(lintCommand)
		if err := updateEphemyralCommand(directory, "lint", refactoredLintCommand); err != nil {
			fmt.Printf("Error updating .ephemyral file: %v\n", err)
			return
		}

		fmt.Printf("Successfully generated and updated lint command: %s\n", refactoredLintCommand)
		if err := executeCommand(directory, refactoredLintCommand); err != nil {
			fmt.Printf("Error executing new lint command: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
