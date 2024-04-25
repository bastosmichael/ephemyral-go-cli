package cmd

import (
	gpt4client "ephemyral/pkg"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func generateBuildCommand(directory string) (string, error) {
	filesList, err := getFileList(directory)
	if err != nil {
		return "", err
	}

	fullPrompt := BuildCommandPrompt + strings.Join(filesList, "\n")
	gpt4client.SetDebug(false)
	buildCommand, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt)
	if err != nil || strings.TrimSpace(buildCommand) == "" {
		return "", fmt.Errorf("error generating or empty build command")
	}

	return buildCommand, nil
}

var buildCmd = &cobra.Command{
	Use:   "build [directory]",
	Short: "Use AI to intelligently generate and execute a build command for the specified directory, optimizing for performance and efficiency.",
	Long:  "The 'build' command generates a building command based on the structure of the project's files. It then updates the '.ephemyral' configuration file with the new build command and executes it. Use this command to ensure your project builds correctly and is free from errors.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]

		// Get the retry count from flags, with a default of 3
		defaultRetryCount, err := cmd.Flags().GetInt("retry")
		if err != nil {
			fmt.Println("Error reading retry count:", err)
			return
		}

		// Try to get existing build command
		existingBuildCommand, err := getExistingCommand(directory, "build")
		if err != nil {
			fmt.Println("Error reading .ephemyral file:", err)
			return
		}

		if existingBuildCommand != "" {
			// Retry execution with the existing command
			success := false
			for i := 0; i < defaultRetryCount; i++ {
				fmt.Println("Running existing build command:", existingBuildCommand)
				if err := executeCommand(directory, existingBuildCommand); err != nil {
					fmt.Println("Error executing build command:", err)
					time.Sleep(retryDelay) // wait before retrying
				} else {
					success = true
					fmt.Println("Successfully executed existing build command:", existingBuildCommand) // Success message
					break
				}
			}

			if !success {
				fmt.Println("Failed to execute existing build command after retries.")
				return
			}

			return
		}

		// Retry generating and executing the build command
		var refactoredBuildCommand string
		success := false
		for i := 0; i < defaultRetryCount; i++ {
			buildCommand, err := generateBuildCommand(directory)
			if err != nil {
				fmt.Println("Error generating build command:", err)
				time.Sleep(retryDelay) // wait before retrying
			} else {
				refactoredBuildCommand = filterOutCodeBlocks(buildCommand)
				fmt.Println("Successfully generated build command:", refactoredBuildCommand) // Success message
				if err := executeCommand(directory, refactoredBuildCommand); err != nil {
					fmt.Println("Error executing build command:", err)
					time.Sleep(retryDelay) // wait before retrying
				} else {
					success = true
					fmt.Println("Successfully executed build command:", refactoredBuildCommand) // Success message
					break
				}
			}
		}

		if !success {
			fmt.Println("Failed to generate or execute build command after retries.")
			return
		}

		// Update the .ephemyral file with the successful build command
		if err := updateEphemyralFile(directory, "build", refactoredBuildCommand); err != nil {
			fmt.Println("Error updating .ephemyral file:", err)
			return
		}

		fmt.Println("Successfully updated .ephemyral with build command:", refactoredBuildCommand) // Success message
	},
}

func init() {
	buildCmd.Flags().Int("retry", 3, "Number of retries for generating and executing build command")
	rootCmd.AddCommand(buildCmd)
}
