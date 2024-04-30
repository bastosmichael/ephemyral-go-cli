package cmd

import (
	"fmt"
	"strings"
	"time"

	gpt4client "ephemyral/pkg"

	"github.com/spf13/cobra"
)

func generateDocsCommand(directory string) (string, error) {
	filesList, err := getFileList(directory)
	if err != nil {
		return "", err
	}

	fullPrompt := DocsCommandPrompt + strings.Join(filesList, "\n")
	gpt4client.SetDebug(false)
	docsCommand, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(docsCommand) == "" {
		return "", fmt.Errorf("received empty docs command")
	}

	return docsCommand, nil
}

var docsCmd = &cobra.Command{
	Use:   "docs [directory]",
	Short: "Use a machine learning model to generate and execute a command for generating documentation, enhancing the clarity and maintainability of the codebase.",
	Long:  "The 'docs' command generates a command line that, when executed, produces documentation (such as a README or API documentation) for the project's files. It then updates the '.ephemyral' configuration file with the new command and executes it. This ensures that the project is well-documented and easier to maintain.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]

		// Get the retry count from flags, with a default of 3
		defaultRetryCount, err := cmd.Flags().GetInt("retry")
		if err != nil {
			fmt.Println("Error reading retry count:", err)
			return
		}

		// Try to get an existing docs command
		existingDocsCommand, err := getExistingCommand(directory, "docs")
		if err != nil {
			fmt.Println("Error reading .ephemyral file:", err)
			return
		}

		if existingDocsCommand != "" {
			// Retry execution with the existing command
			success := false
			for i := 0; i < defaultRetryCount; i++ {
				fmt.Println("Running existing docs command:", existingDocsCommand)
				if err := executeCommand(directory, existingDocsCommand); err != nil {
					fmt.Println("Error executing docs command:", err)
					time.Sleep(retryDelay) // wait before retrying
				} else {
					success = true
					fmt.Println("Successfully executed existing docs command:", existingDocsCommand) // Success message
					break
				}
			}

			if !success {
				fmt.Println("Failed to execute existing docs command after retries.")
				return
			}

			return
		}

		// Retry generating and executing the docs command
		var refactoredDocsCommand string
		success := false
		for i := 0; i < defaultRetryCount; i++ {
			docsCommand, err := generateDocsCommand(directory)
			if err != nil {
				fmt.Println("Error generating docs command:", err)
				time.Sleep(retryDelay) // wait before retrying
			} else {
				refactoredDocsCommand = filterOutCodeBlocks(docsCommand)
				fmt.Println("Successfully discovered docs command:", refactoredDocsCommand) // Success message
				if err := executeCommand(directory, refactoredDocsCommand); err != nil {
					fmt.Println("Error executing docs command:", err)
					time.Sleep(retryDelay) // wait before retrying
				} else {
					success = true
					fmt.Println("Successfully executed docs command:", refactoredDocsCommand) // Success message
					break
				}
			}
		}

		if !success {
			fmt.Println("Failed to generate or execute docs command after retries.")
			return
		}

		// Update the .ephemyral file with the successful docs command
		if err := updateEphemyralFile(directory, "docs", refactoredDocsCommand); err != nil {
			fmt.Println("Error updating .ephemyral file:", err)
			return
		}

		fmt.Println("Successfully updated .ephemyral with docs command:", refactoredDocsCommand) // Success message
	},
}

func init() {
	docsCmd.Flags().Int("retry", 3, "Number of retries for generating and executing docs command")
	rootCmd.AddCommand(docsCmd)
}
