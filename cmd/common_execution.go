//go:build !lint
// +build !lint

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
)

func runCommand(cmdType, filePath string, convID uuid.UUID, retryCount int, retryDelay time.Duration) bool {
	directory, err := findEphemyralDirectory(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	existingCommand, err := getExistingCommandOrError(directory, cmdType)
	if err != nil || existingCommand == "" {
		fmt.Println("Error reading existing", cmdType, "command:", err)
		return false
	}

	if err := executeWithRetries(directory, existingCommand, cmdType, convID, retryCount, retryDelay); err != nil {
		fmt.Println("Failed to execute", cmdType, "command:", err)
		return false
	}
	fmt.Println(cmdType, "command executed successfully.")
	return true
}

// executeBuildCommand executes the build command using os/exec.
func executeCommand(directory, command string) error {
	// Create a new command from the buildCommand string.
	cmd := exec.Command("bash", "-c", command)

	// Set the command's working directory to the specified one.
	cmd.Dir = directory

	// Redirect output to the console (or you could handle it in other ways).
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command and check for errors.
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func getExistingCommandOrError(directory, commandType string) (string, error) {
	existingCommand, err := getExistingCommand(directory, commandType)
	if err != nil {
		return "", fmt.Errorf("error reading existing %s command: %v", commandType, err)
	}
	return existingCommand, nil
}

// Executes a command of a given type (e.g., test, lint, build, docs) in the specified directory.
func executeCommandOfType(directory, commandType string, convID uuid.UUID, retryCount int, retryDelay time.Duration) error {
	// Try to get an existing command of the given type.
	existingCommand, err := getExistingCommandOrError(directory, commandType)
	if err != nil {
		return err
	}

	if existingCommand != "" {
		// Retry execution with the existing command.
		if err := executeWithRetries(directory, existingCommand, commandType, convID, retryCount, retryDelay); err != nil {
			return fmt.Errorf("failed to execute %s command after retries: %v", commandType, err)
		}
		return nil
	}

	// Generate and execute a new command.
	return generateAndExecuteCommand(directory, commandType, convID, retryCount, retryDelay)
}

// Shared function to handle command execution with retries and dependency management
func executeWithRetries(directory, command, commandType string, convID uuid.UUID, retryCount int, retryDelay time.Duration) error {
	for i := 0; i < retryCount; i++ {
		fmt.Printf("Running %s command: %s\n", commandType, command)
		if err := executeCommand(directory, command); err != nil {
			fmt.Println("Error executing command:", err)

			// Handle missing dependency error
			dependencyCommand, depErr := generateDependencyCommand(command, err.Error(), convID)
			if depErr != nil {
				return fmt.Errorf("error generating dependency command: %v", depErr)
			}

			fmt.Printf("Running dependency installation command: %s\n", dependencyCommand)
			if depErr := executeCommand(directory, dependencyCommand); depErr != nil {
				fmt.Println("Error executing dependency command:", depErr)
				time.Sleep(retryDelay)
			} else {
				// Retry executing the original command
				if err := executeCommand(directory, command); err != nil {
					fmt.Println("Error executing command:", err)
					time.Sleep(retryDelay)
				} else {
					fmt.Printf("Successfully executed %s command after dependency installation: %s\n", commandType, command)
					return nil
				}
			}
		} else {
			fmt.Printf("Successfully executed %s command: %s\n", commandType, command)
			return nil
		}
	}
	return fmt.Errorf("failed to execute %s command after retries", commandType)
}
