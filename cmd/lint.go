package cmd

import (
	"fmt"
	"strings"
	"time"

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

		// Get the retry count from flags, with a default of 3
		defaultRetryCount, err := cmd.Flags().GetInt("retry")
		if err != nil {
			fmt.Println("Error reading retry count:", err)
			return
		}

		// Try to get existing lint command
		existingLintCommand, err := getExistingCommand(directory, "lint")
		if err != nil {
			fmt.Println("Error reading .ephemyral file:", err)
			return
		}

		if existingLintCommand != "" {
			// Retry execution with the existing command
			success := false
			var lastError error
			for i := 0; i < defaultRetryCount; i++ {
				fmt.Println("Running existing lint command:", existingLintCommand)
				lastError = executeCommand(directory, existingLintCommand)
				if lastError != nil {
					fmt.Println("Error executing lint command:", lastError)
					time.Sleep(retryDelay) // wait before retrying
				} else {
					success = true
					fmt.Println("Successfully executed existing lint command:", existingLintCommand) // Success message
					break
				}
			}

			if !success {
				fmt.Println("Failed to execute existing lint command after retries:", lastError)
				return
			}

			return
		}

		// Generate and execute a new lint command
		var refactoredLintCommand string
		success := false
		var lastError error
		for i := 0; i < defaultRetryCount; i++ {
			lintCommand, err := generateLintCommand(directory)
			if err != nil {
				lastError = err
				fmt.Println("Error generating lint command:", err)
				time.Sleep(retryDelay)
				continue
			}

			refactoredLintCommand = filterOutCodeBlocks(lintCommand)
			fmt.Println("Successfully generated lint command:", refactoredLintCommand)

			lastError = executeCommand(directory, refactoredLintCommand)
			if lastError != nil {
				// Handle missing dependency error
				dependencyCommand, depErr := generateDependencyCommand(refactoredLintCommand, lastError.Error())
				if depErr != nil {
					fmt.Println("Error generating dependency command:", depErr)
					return
				}

				fmt.Println("Running dependency installation command:", dependencyCommand)
				lastError = executeCommand(directory, dependencyCommand)
				if lastError != nil {
					fmt.Println("Error executing dependency command:", lastError)
					time.Sleep(retryDelay)
				} else {
					// Retry executing lint command
					lastError = executeCommand(directory, refactoredLintCommand)
					if lastError != nil {
						fmt.Println("Error executing lint command:", lastError)
						time.Sleep(retryDelay)
					} else {
						success = true
						fmt.Println("Successfully executed lint command after dependency installation:", refactoredLintCommand)
						break
					}
				}
			} else {
				success = true
				fmt.Println("Successfully executed lint command:", refactoredLintCommand)
				break
			}
		}

		if !success {
			fmt.Println("Failed to generate or execute lint command after retries:", lastError)
			return
		}

		// Update the .ephemyral file with the successful lint command
		if err := updateEphemyralFile(directory, "lint", refactoredLintCommand); err != nil {
			fmt.Println("Error updating .ephemyral file:", err)
			return
		}

		fmt.Println("Successfully updated .ephemyral with lint command:", refactoredLintCommand) // Success message
	},
}

func init() {
	lintCmd.Flags().Int("retry", 3, "Number of retries for generating and executing lint command")
	rootCmd.AddCommand(lintCmd)
}
