package cmd

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v2"
)

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
