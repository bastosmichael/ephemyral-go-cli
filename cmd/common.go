// common.go
//go:build !lint
// +build !lint

package cmd

import (
	gpt4client "ephemyral/pkg"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// EphemyralFile represents the structure of the .ephemyral YAML file.
type EphemyralFile struct {
	BuildCommand string `yaml:"build-command"`
	TestCommand  string `yaml:"test-command"`
	LintCommand  string `yaml:"lint-command"`
	DocsCommand  string `yaml:"docs-command"`
}

var retryDelay = 2 * time.Second

func generateDependencyCommand(failedCommand, errorMessage string) (string, error) {
	// Determine the operating system
	osType := runtime.GOOS

	// Construct a prompt to handle missing dependencies
	prompt := fmt.Sprintf("The following command '%s' failed with the error '%s'. Based on this error and the current operating system '%s', provide the simplest single-line command to install all necessary dependencies. The response should contain no comments, explanations, or code blocks, and if multiple commands are needed, they should be separated by '&&'. Include necessary flags like '-y' for automatic confirmation:\n", failedCommand, errorMessage, osType)

	dependencyCommand, err := gpt4client.GetGPT4ResponseWithPrompt(prompt)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(dependencyCommand) == "" {
		return "", fmt.Errorf("received empty dependency command")
	}

	return dependencyCommand, nil
}

// getFileList retrieves a list of all non-directory file names in the specified directory and its subdirectories,
// skipping specified directories like ".git".
func getFileList(directory string) ([]string, error) {
	var filesList []string

	// Read the files in the root directory.
	rootFiles, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	// Add non-directory files from the root directory to the list.
	for _, file := range rootFiles {
		if !file.IsDir() {
			filesList = append(filesList, file.Name())
		}
	}

	// Walk through the directory and its subdirectories.
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.HasSuffix(info.Name(), ".git") {
				return filepath.SkipDir
			}
		} else {
			relativePath, err := filepath.Rel(directory, path)
			if err != nil {
				return err
			}
			filesList = append(filesList, relativePath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filesList, nil
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

// getExistingCommand reads the existing command from the .ephemyral file based on the key.
func getExistingCommand(directory, key string) (string, error) {
	filename := directory + "/.ephemyral"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "", nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	var ephemyral EphemyralFile
	if err := yaml.Unmarshal(data, &ephemyral); err != nil {
		return "", err
	}

	if key == "build" {
		return ephemyral.BuildCommand, nil
	} else if key == "test" {
		return ephemyral.TestCommand, nil
	} else if key == "lint" {
		return ephemyral.LintCommand, nil
	} else if key == "docs" {
		return ephemyral.DocsCommand, nil
	}

	return "", nil
}

// updateEphemyralFile updates the specified key in the .ephemyral file.
func updateEphemyralFile(directory, key, command string) error {
	filename := directory + "/.ephemyral"

	var ephemyral EphemyralFile
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		data, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(data, &ephemyral); err != nil {
			return err
		}
	}

	switch key {
	case "build":
		ephemyral.BuildCommand = command
	case "test":
		ephemyral.TestCommand = command
	case "lint":
		ephemyral.LintCommand = command
	case "docs":
		ephemyral.DocsCommand = command
	default:
		return fmt.Errorf("unknown key: %s", key)
	}

	data, err := yaml.Marshal(&ephemyral)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		fmt.Println("Error updating .ephemyral file:", err)
		return err
	}

	fmt.Printf("Successfully updated .ephemyral with %s command: %s\n", key, command)
	return nil
}

func filterOutCodeBlocks(content string) string {
	lines := strings.Split(content, "\n")
	filteredLines := make([]string, 0)
	for _, line := range lines {
		if !strings.Contains(line, "```") {
			filteredLines = append(filteredLines, line)
		}
	}
	// Join the filtered lines and ensure it ends with a newline
	result := strings.Join(filteredLines, "\n")
	if !strings.HasSuffix(result, "\n") {
		result += "\n" // Add a newline if it doesn't exist
	}
	return result
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// Function type for generating commands.
type commandGenerator func(directory string) (string, error)

// Map associating command types with their respective generation functions.
var commandGenerators = map[string]commandGenerator{
	"test":  generateTestCommand,
	"lint":  generateLintCommand,
	"build": generateBuildCommand,
	"docs":  generateDocsCommand,
}

// Executes a command of a given type (e.g., test, lint, build, docs) in the specified directory.
func executeCommandOfType(directory, commandType string, retryCount int, retryDelay time.Duration) error {
	// Try to get an existing command of the given type.
	existingCommand, err := getExistingCommand(directory, commandType)
	if err != nil {
		return fmt.Errorf("error reading existing %s command: %v", commandType, err)
	}

	if existingCommand != "" {
		// Retry execution with the existing command.
		if err := retryExecution(directory, existingCommand, commandType, retryCount, retryDelay); err != nil {
			return fmt.Errorf("failed to execute %s command after retries: %v", commandType, err)
		}
		return nil
	}

	// Generate and execute a new command.
	return generateAndExecuteCommand(directory, commandType, retryCount, retryDelay)
}

// Retries execution of a given command a specified number of times.
func retryExecution(directory, command, commandType string, retryCount int, retryDelay time.Duration) error {
	generationFailed := false
	executionFailed := false

	for i := 0; i < retryCount; i++ {
		fmt.Printf("Running %s command: %s\n", commandType, command)
		if err := executeCommand(directory, command); err != nil {
			fmt.Println("Error executing command:", err)
			executionFailed = true
			time.Sleep(retryDelay)
		} else {
			fmt.Printf("Successfully executed %s command: %s\n", commandType, command)
			return nil
		}
	}

	if generationFailed && executionFailed {
		return fmt.Errorf("failed to generate or execute %s command after retries", commandType)
	}

	return fmt.Errorf("failed to execute %s command after retries", commandType)
}

// Generates and executes a new command of a given type.
func generateAndExecuteCommand(directory, commandType string, retryCount int, retryDelay time.Duration) error {
	var refactoredCommand string

	generationFailed := false
	executionFailed := false

	for i := 0; i < retryCount; i++ {
		// Generate a new command using the appropriate generator function.
		generator, found := commandGenerators[commandType]
		if !found {
			return fmt.Errorf("generator function not found for command type: %s", commandType)
		}

		command, err := generator(directory)
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
			dependencyCommand, depErr := generateDependencyCommand(refactoredCommand, err.Error())
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
					// Update the .ephemyral file with the successful command
					if err := updateEphemyralFile(directory, commandType, refactoredCommand); err != nil {
						fmt.Println("Error updating .ephemyral file:", err)
						return err
					}
					return nil
				}
			}
		} else {
			fmt.Printf("Successfully executed %s command: %s\n", commandType, refactoredCommand)
			// Update the .ephemyral file with the successful command
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
