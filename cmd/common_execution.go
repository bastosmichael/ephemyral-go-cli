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

const (
	BashCmd = "bash"
	BashOpt = "-c"
)

func runCommand(cmdType, filePath string, convID uuid.UUID, retryCount int, retryDelay time.Duration) bool {
	directory, err := findEphemyralDirectory(filePath)
	if err != nil {
		return logError("Error:", err)
	}

	return executeCommandOrLog(directory, cmdType, convID, retryCount, retryDelay)
}

func logError(msg string, err error) bool {
	fmt.Println(msg, err)
	return false
}

func executeCommandOrLog(directory, cmdType string, convID uuid.UUID, retryCount int, retryDelay time.Duration) bool {
	cmd, err := getExistingCommandOrError(directory, cmdType)
	if err != nil || cmd == "" {
		return logError("Error reading existing "+cmdType+" command:", err)
	}

	return logExecutionResult(directory, cmd, cmdType, convID, retryCount, retryDelay)
}

func logExecutionResult(directory, cmd, cmdType string, convID uuid.UUID, retryCount int, retryDelay time.Duration) bool {
	if err := executeWithRetries(directory, cmd, cmdType, convID, retryCount, retryDelay); err != nil {
		return logError("Failed to execute "+cmdType+" command:", err)
	}
	fmt.Println(cmdType, "command executed successfully.")
	return true
}

func executeCommand(directory, command string) error {
	cmd := createCommand(directory, command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createCommand(directory, command string) *exec.Cmd {
	cmd := exec.Command(BashCmd, BashOpt, command)
	cmd.Dir = directory
	return cmd
}

func getExistingCommandOrError(directory, commandType string) (string, error) {
	cmd, err := getExistingCommand(directory, commandType)
	return cmd, wrapError(err, "reading existing "+commandType+" command")
}

func wrapError(err error, msg string) error {
	if err != nil {
		return fmt.Errorf("error %s: %v", msg, err)
	}
	return nil
}

func executeCommandOfType(directory, commandType string, convID uuid.UUID, retryCount int, retryDelay time.Duration) error {
	cmd, err := getExistingCommandOrError(directory, commandType)
	if err != nil {
		return err
	}

	return executeOrGenerateCommand(directory, cmd, commandType, convID, retryCount, retryDelay)
}

func executeOrGenerateCommand(directory, cmd, cmdType string, convID uuid.UUID, retryCount int, retryDelay time.Duration) error {
	if cmd != "" {
		return executeWithRetryHandling(directory, cmd, cmdType, convID, retryCount, retryDelay)
	}

	return generateAndExecuteCommand(directory, cmdType, convID, retryCount, retryDelay)
}

func executeWithRetryHandling(directory, cmd, cmdType string, convID uuid.UUID, retryCount int, retryDelay time.Duration) error {
	if err := executeWithRetries(directory, cmd, cmdType, convID, retryCount, retryDelay); err != nil {
		return fmt.Errorf("failed to execute %s command after retries: %v", cmdType, err)
	}
	return nil
}

func executeWithRetries(directory, command, commandType string, convID uuid.UUID, retryCount int, retryDelay time.Duration) error {
	for i := 0; i < retryCount; i++ {
		if err := tryExecuteCommand(directory, command, commandType, convID, retryDelay); err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to execute %s command after retries", commandType)
}

func tryExecuteCommand(directory, command, commandType string, convID uuid.UUID, retryDelay time.Duration) error {
	fmt.Printf("Running %s command: %s\n", commandType, command)
	if err := executeCommand(directory, command); err != nil {
		return handleExecutionError(directory, command, commandType, convID, err, retryDelay)
	}
	fmt.Printf("Successfully executed %s command: %s\n", commandType, command)
	return nil
}

func handleExecutionError(directory, command, commandType string, convID uuid.UUID, err error, retryDelay time.Duration) error {
	fmt.Println("Error executing command:", err)
	dependencyCommand, depErr := generateDependencyCommand(command, err.Error(), convID)
	if depErr != nil {
		return fmt.Errorf("error generating dependency command: %v", depErr)
	}

	return tryDependencyCommand(directory, command, commandType, dependencyCommand, retryDelay)
}

func tryDependencyCommand(directory, command, commandType, dependencyCommand string, retryDelay time.Duration) error {
	fmt.Printf("Running dependency installation command: %s\n", dependencyCommand)
	if depErr := executeCommand(directory, dependencyCommand); depErr != nil {
		fmt.Println("Error executing dependency command:", depErr)
		time.Sleep(retryDelay)
	} else {
		return reattemptOriginalCommand(directory, command, commandType, retryDelay)
	}
	return nil
}

func reattemptOriginalCommand(directory, command, commandType string, retryDelay time.Duration) error {
	if err := executeCommand(directory, command); err != nil {
		fmt.Println("Error executing command:", err)
		time.Sleep(retryDelay)
	} else {
		fmt.Printf("Successfully executed %s command after dependency installation: %s\n", commandType, command)
		return nil
	}
	return nil
}
