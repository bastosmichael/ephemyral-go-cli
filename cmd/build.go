package cmd

import (
	gpt4client "ephemyral/pkg"
	"fmt"
	"strings"

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

		existingBuildCommand, err := getExistingCommand(directory, "build")
		if err != nil {
			fmt.Printf("Error reading .ephemyral file: %v\n", err)
			return
		}

		if existingBuildCommand != "" {
			fmt.Println("Running existing build command:", existingBuildCommand)
			if err := executeCommand(directory, existingBuildCommand); err != nil {
				fmt.Printf("Error executing build command: %v\n", err)
			}
			return
		}

		buildCommand, err := generateBuildCommand(directory)
		if err != nil {
			fmt.Printf("Error generating build command: %v\n", err)
			return
		}

		refactoredBuildCommand := filterOutCodeBlocks(buildCommand)

		if err := updateEphemyralCommand(directory, "build", refactoredBuildCommand); err != nil {
			fmt.Printf("Error updating .ephemyral file: %v\n", err)
			return
		}

		fmt.Println("Build command generated and updated:", refactoredBuildCommand)

		if err := executeCommand(directory, refactoredBuildCommand); err != nil {
			fmt.Printf("Error executing build command: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
