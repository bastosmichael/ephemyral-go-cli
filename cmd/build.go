// build.go
package cmd

import (
	gpt4client "ephemyral/pkg"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// generateBuildCommand generates a build command by listing all files in the root directory and subdirectories.
func generateBuildCommand(directory string) (string, error) {
	filesList, err := getFileList(directory)
	if err != nil {
		return "", err
	}

	// Join all file names into a single prompt.
	fullPrompt := "Based on the following file list, provide the simplest command line required to build these files. The command must be in a single line and contain no extra text or commentary:\n" +
		strings.Join(filesList, "\n")

	gpt4client.SetDebug(false)
	buildCommand, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(buildCommand) == "" {
		return "", fmt.Errorf("received empty build command")
	}

	return buildCommand, nil
}

var buildCmd = &cobra.Command{
	Use:   "build [directory]",
	Short: "Generate and run a build command for the specified directory",
	Long:  "The 'build' command retrieves the build command specified in the '.ephemyral' configuration file within the given directory. If no build command is specified, the command generates a build command based on the structure of the directory given. It then updates the '.ephemyral' file with the new build command and executes it. Use this command to ensure your project runs the appropriate build before other tasks.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]

		existingBuildCommand, err := getExistingCommand(directory, "build")
		if err != nil {
			fmt.Println("Error reading .ephemyral file:", err)
			return
		}

		if existingBuildCommand != "" {
			fmt.Println("Running existing build command:", existingBuildCommand)
			if err := executeCommand(directory, existingBuildCommand); err != nil {
				fmt.Println("Error executing new build command:", err)
				return
			}
			return
		}

		buildCommand, err := generateBuildCommand(directory)

		if err != nil {
			fmt.Println("Error generating build command:", err)
			return
		}

		refactoredBuildCommand := filterOutCodeBlocks(buildCommand)

		if err := updateEphemyralCommand(directory, "build", refactoredBuildCommand); err != nil {
			fmt.Println("Error updating .ephemyral file:", err)
			return
		}

		fmt.Println("Successfully generated and updated build command:", refactoredBuildCommand)

		// Execute the new build command.
		if err := executeCommand(directory, buildCommand); err != nil {
			fmt.Println("Error executing new build command:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
