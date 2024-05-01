// +build !lint

package cmd

import (
	"fmt"
	"strings"

	gpt4client "ephemyral/pkg"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func generateLintCommand(directory string, convID uuid.UUID) (string, error) {
	filesList, err := getFileList(directory)
	if err != nil {
		return "", err
	}

	fullPrompt := LintCommandPrompt + strings.Join(filesList, "\n")
	gpt4client.SetDebug(false)
	lintCommand, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt, convID)
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

		// Get the retry count from flags, with a default of 3
		defaultRetryCount, err := cmd.Flags().GetInt("retry")
		if err != nil {
			fmt.Println("Error reading retry count:", err)
			return
		}

		convID := uuid.New()
		fmt.Println(convID)

		if err := executeCommandOfType(directory, "lint", convID, defaultRetryCount, retryDelay); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	lintCmd.Flags().Int("retry", 3, "Number of retries for generating and executing lint command")
	rootCmd.AddCommand(lintCmd)
}
