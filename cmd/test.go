package cmd

import (
	"fmt"
	"strings"
	"time"

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
	Short: "Deploy AI models to generate and run optimized test commands for the specified directory, enhancing test accuracy and efficiency.",
	Long:  "The 'test' command generates a testing command based on the structure of the project's files. It then updates the '.ephemyral' configuration file with the new testing command and executes it. Use this command to ensure your project adheres to testing standards and is free from test errors.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]

		// Get the retry count from flags, with a default of 3
		defaultRetryCount, err := cmd.Flags().GetInt("retry")
		if err != nil {
			fmt.Println("Error reading retry count:", err)
			return
		}

		// Try to get existing test command
		existingTestCommand, err := getExistingCommand(directory, "test")
		if err != nil {
			fmt.Println("Error reading .ephemyral file:", err)
			return
		}

		if existingTestCommand != "" {
			// Retry execution with the existing command
			success := false
			for i := 0; i < defaultRetryCount; i++ {
				fmt.Println("Running existing test command:", existingTestCommand)
				if err := executeCommand(directory, existingTestCommand); err != nil {
					fmt.Println("Error executing test command:", err)
					time.Sleep(retryDelay) // wait before retrying
				} else {
					success = true
					fmt.Println("Successfully executed existing test command:", existingTestCommand) // Success message
					break
				}
			}

			if !success {
				fmt.Println("Failed to execute existing test command after retries.")
				return
			}

			return
		}

		// Retry generating and executing the test command
		var refactoredTestCommand string
		success := false
		for i := 0; i < defaultRetryCount; i++ {
			testCommand, err := generateTestCommand(directory)
			if err != nil {
				fmt.Println("Error generating test command:", err)
				time.Sleep(retryDelay) // wait before retrying
			} else {
				refactoredTestCommand = filterOutCodeBlocks(testCommand)
				fmt.Println("Successfully generated test command:", refactoredTestCommand) // Success message
				if err := executeCommand(directory, refactoredTestCommand); err != nil {
					fmt.Println("Error executing test command:", err)
					time.Sleep(retryDelay) // wait before retrying
				} else {
					success = true
					fmt.Println("Successfully executed test command:", refactoredTestCommand) // Success message
					break
				}
			}
		}

		if !success {
			fmt.Println("Failed to generate or execute test command after retries.")
			return
		}

		// Update the .ephemyral file with the successful test command
		if err := updateEphemyralCommand(directory, "test", refactoredTestCommand); err != nil {
			fmt.Println("Error updating .ephemyral file:", err)
			return
		}

		fmt.Println("Successfully updated .ephemyral with test command:", refactoredTestCommand) // Success message
	},
}

func init() {
	testCmd.Flags().Int("retry", 3, "Number of retries for generating and executing test command")
	rootCmd.AddCommand(testCmd)
}
