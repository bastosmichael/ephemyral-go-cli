// +build !lint

package cmd

import (
	gpt4client "ephemyral/pkg"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// Generates a build command.
func generateBuildCommand(directory string, convID uuid.UUID) (string, error) {
	filesList, err := getFileList(directory)
	if err != nil {
		return "", err
	}

	fullPrompt := BuildCommandPrompt + strings.Join(filesList, "\n")
	gpt4client.SetDebug(true)
	buildCommand, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt, convID)
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

		convID := uuid.New()
		fmt.Println(convID)

		if err := executeCommandOfType(directory, "build", convID, defaultRetryCount, retryDelay); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	buildCmd.Flags().Int("retry", 3, "Number of retries for generating and executing build command")
	rootCmd.AddCommand(buildCmd)
}
