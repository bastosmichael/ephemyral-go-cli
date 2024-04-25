// common.go
package cmd

import (
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

// EphemyralFile represents the structure of the .ephemyral YAML file.
type EphemyralFile struct {
    BuildCommand string `yaml:"build_command"`
    TestCommand  string `yaml:"test_command"`
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
    }

    return "", nil
}

// updateEphemyralCommand updates the specified key in the .ephemyral file.
func updateEphemyralCommand(directory, key, command string) error {
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
    return strings.Join(filteredLines, "\n")
}
