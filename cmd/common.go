// common.go
package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// EphemyralFile represents the structure of the .ephemyral YAML file.
type EphemyralFile struct {
    BuildCommand string `yaml:"build_command"`
    TestCommand  string `yaml:"test_command"`
    LintCommand  string `yaml:"lint_command"`
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
    } else if key == "lint" {
        ephemyral.LintCommand = command
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
