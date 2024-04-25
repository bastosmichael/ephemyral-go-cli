// lint.go
package cmd

import (
	"fmt"
	"strings"

	gpt4client "ephemyral/pkg"

	"github.com/spf13/cobra"
)

// generateLintCommand generates a linting command based on the file structure in the directory.
func generateLintCommand(directory string) (string, error) {
	filesList, err := getFileList(directory)
	if err != nil {
		return "", err
	}

	// Join all file names into a single prompt for linting.
	fullPrompt := "Based on the following file list, provide the simplest command line required to lint these files. The command must be in a single line and contain no extra text or commentary:\n" +
		strings.Join(filesList, "\n")

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
	Short: "Generate and run a linting command for the specified directory",
	Long:  "The 'lint' command generates a linting command based on the structure of the project's files. It then updates the '.ephemyral' configuration file with the new linting command and executes it. Use this command to ensure your project adheres to coding standards and is free from basic syntax errors.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]

		lintCommand, err := generateLintCommand(directory)

		if err != nil {
			fmt.Println("Error generating linting command:", err)
			return
		}

		refactoredLintCommand := filterOutCodeBlocks(lintCommand)

		if err := updateEphemyralCommand(directory, "lint", refactoredLintCommand); err != nil {
			fmt.Println("Error updating .ephemyral file:", err)
			return
		}

		fmt.Println("Successfully generated and updated lint command:", refactoredLintCommand)

		// Execute the linting command.
		if err := executeCommand(directory, refactoredLintCommand); err != nil {
			fmt.Println("Error executing linting command:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
