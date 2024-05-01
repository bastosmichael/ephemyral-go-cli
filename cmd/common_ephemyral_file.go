//go:build !lint
// +build !lint

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// EphemyralFile represents the structure of the .ephemyral YAML file.
type EphemyralFile struct {
	BuildCommand string `yaml:"build-command"`
	TestCommand  string `yaml:"test-command"`
	LintCommand  string `yaml:"lint-command"`
	DocsCommand  string `yaml:"docs-command"`
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

	switch key {
	case "build":
		return ephemyral.BuildCommand, nil
	case "test":
		return ephemyral.TestCommand, nil
	case "lint":
		return ephemyral.LintCommand, nil
	case "docs":
		return ephemyral.DocsCommand, nil
	default:
		return "", fmt.Errorf("unknown key: %s", key)
	}
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

func findEphemyralDirectory(filePath string) (string, error) {
	dir := filepath.Dir(filePath)

	for {
		// Check if ".ephemyral" file exists in the current directory
		ephemyralPath := filepath.Join(dir, ".ephemyral")
		if _, err := os.Stat(ephemyralPath); err == nil {
			return dir, nil
		}

		// Move up to the parent directory
		parentDir := filepath.Dir(dir)
		if parentDir == dir { // If we've reached the root directory
			break
		}
		dir = parentDir
	}

	return "", fmt.Errorf(".ephemyral file not found in any directory upwards from %s", filePath)
}
