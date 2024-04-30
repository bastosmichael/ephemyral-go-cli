// common.go
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
	BuildCommand   string `yaml:"build_command"`
	TestCommand    string `yaml:"test_command"`
	LintCommand    string `yaml:"lint_command"`
	DocsCommand    string `yaml:"docs_command"`
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

	if key == "build" {
		ephemyral.BuildCommand = command
	} else if key == "test" {
		ephemyral.TestCommand = command
	} else if key == "lint" {
		ephemyral.LintCommand = command
	} else if key == "docs" {
		ephemyral.DocsCommand = command
	}

	data, err := yaml.Marshal(&ephemyral)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
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
