// +build !lint

package cmd

import (
	"fmt"
	"runtime"
	"strings"
	"time"
	"github.com/google/uuid"
	gpt4client "ephemyral/pkg"
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

// Generates and executes a new command of a given type.
func generateAndExecuteCommand(directory, commandType string, convID uuid.UUID, retryCount int, retryDelay time.Duration) error {
	var refactoredCommand string

	generationFailed := false
	executionFailed := false

	for i := 0; i < retryCount; i++ {
		// Generate a new command using the appropriate generator function.
		generator, found := commandGenerators[commandType]
		if !found {
			return fmt.Errorf("generator function not found for command type: %s", commandType)
		}

		command, err := generator(directory, convID)
		if err != nil {
			fmt.Println("Error generating command:", err)
			generationFailed = true
			time.Sleep(retryDelay)
			continue
		}

		refactoredCommand = filterOutCodeBlocks(command)
		fmt.Printf("Successfully generated %s command: %s\n", commandType, refactoredCommand)

		// Attempt to execute the command.
		if err := executeCommand(directory, refactoredCommand); err != nil {
			// Handle missing dependency error.
			dependencyCommand, depErr := generateDependencyCommand(refactoredCommand, err.Error(), convID)
			if depErr != nil {
				return fmt.Errorf("error generating dependency command: %v", depErr)
			}

			fmt.Printf("Running dependency installation command: %s\n", dependencyCommand)
			if depErr := executeCommand(directory, dependencyCommand); depErr != nil {
				fmt.Println("Error executing dependency command:", depErr)
				time.Sleep(retryDelay)
			} else {
				// Retry executing the original command.
				if err := executeCommand(directory, refactoredCommand); err != nil {
					fmt.Println("Error executing command:", err)
					executionFailed = true
					time.Sleep(retryDelay)
				} else {
					fmt.Printf("Successfully executed %s command after dependency installation: %s\n", commandType, refactoredCommand)
					// Update the .ephemyral file with the successful command.
					if err := updateEphemyralFile(directory, commandType, refactoredCommand); err != nil {
						fmt.Println("Error updating .ephemyral file:", err)
						return err
					}
					return nil
				}
			}
		} else {
			fmt.Printf("Successfully executed %s command: %s\n", commandType, refactoredCommand)
			// Update the .ephemyral file with the successful command.
			if err := updateEphemyralFile(directory, commandType, refactoredCommand); err != nil {
				fmt.Println("Error updating .ephemyral file:", err)
				return err
			}
			return nil
		}
	}

	if generationFailed && executionFailed {
		return fmt.Errorf("failed to generate or execute %s command after retries", commandType)
	}

	return fmt.Errorf("failed to execute %s command after retries", commandType)
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
