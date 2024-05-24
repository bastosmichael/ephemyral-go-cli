//go:build !lint
// +build !lint

package cmd

import (
	gpt4client "ephemyral/pkg"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

var retryDelay = 2 * time.Second

// Function type for generating commands.
type commandGenerator func(directory string, convID uuid.UUID) (string, error)

// Map associating command types with their respective generation functions.
var commandGenerators = map[string]commandGenerator{
	"test":  generateTestCommand,
	"lint":  generateLintCommand,
	"build": generateBuildCommand,
	"docs":  generateDocsCommand,
}

// Generates and executes a new command of a given type
func generateAndExecuteCommand(directory, commandType string, convID uuid.UUID, retryCount int, retryDelay time.Duration) error {
	generator, found := commandGenerators[commandType]
	if !found {
		return fmt.Errorf("generator function not found for command type: %s", commandType)
	}

	for i := 0; i < retryCount; i++ {
		// Generate a new command using the appropriate generator function
		command, err := generator(directory, convID)
		if err != nil {
			fmt.Println("Error generating command:", err)
			time.Sleep(retryDelay)
			continue
		}

		refactoredCommand := filterOutCodeBlocks(command)
		fmt.Printf("Successfully generated %s command: %s\n", commandType, refactoredCommand)

		// Execute the generated command with retries
		if err := executeWithRetries(directory, refactoredCommand, commandType, convID, retryCount, retryDelay); err != nil {
			fmt.Println(err)
			time.Sleep(retryDelay)
		} else {
			// Update the .ephemyral file with the successful command
			if err := updateEphemyralFile(directory, commandType, refactoredCommand); err != nil {
				fmt.Println("Error updating .ephemyral file:", err)
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("failed to generate or execute %s command after retries", commandType)
}

func generateDependencyCommand(failedCommand, errorMessage string, convID uuid.UUID) (string, error) {
	// Determine the operating system
	osType := runtime.GOOS

	// Construct a prompt to handle missing dependencies
	prompt := fmt.Sprintf("The following command '%s' failed with the error '%s'. Based on this error and the current operating system '%s', provide the simplest single-line command to install all necessary dependencies. The response should contain no comments, explanations, or code blocks, and if multiple commands are needed, they should be separated by '&&'. Include necessary flags like '-y' for automatic confirmation:\n", failedCommand, errorMessage, osType)
	dependencyCommand, err := gpt4client.GetGPT4ResponseWithPrompt(prompt, convID)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(dependencyCommand) == "" {
		return "", fmt.Errorf("received empty dependency command")
	}

	return dependencyCommand, nil
}
