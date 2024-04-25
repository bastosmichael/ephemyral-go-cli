// build.go
package cmd

import (
    gpt4client "ephemyral/pkg"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"
    "gopkg.in/yaml.v2"
)

// EphemyralFile represents the structure of the .ephemyral YAML file.
type EphemyralFile struct {
    BuildCommand string `yaml:"build_command"`
    TestCommand  string `yaml:"test_command"`
}

// getExistingBuildCommand reads the existing build command from the .ephemyral file.
func getExistingBuildCommand(directory string) (string, error) {
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

    return ephemyral.BuildCommand, nil
}

// generateBuildCommand generates a build command by listing all files in the directory and subdirectories.
func generateBuildCommand(directory string) (string, error) {
    var filesList []string

    // Walk through the directory and its subdirectories to get all file names.
    err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() { // Only add files, not directories.
            relativePath, err := filepath.Rel(directory, path) // Get relative paths.
            if err != nil {
                return err
            }
            filesList = append(filesList, relativePath)
        }

        return nil
    }) // Correctly close this anonymous function

    if err != nil {
        return "", err
    }

    // Join all file names into a single prompt.
    fullPrompt := strings.Join(filesList, "\n")

    gpt4client.SetDebug(true)
    buildCommand, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt)
    if err != nil {
        return "", err
    }

    if strings.TrimSpace(buildCommand) == "" {
        return "", fmt.Errorf("received empty build command")
    }

    return buildCommand, nil
}

// updateEphemyralBuildCommand updates the .ephemyral file with a new build command.
func updateEphemyralBuildCommand(directory, buildCommand string) error {
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

    ephemyral.BuildCommand = buildCommand

    data, err := yaml.Marshal(&ephemyral)
    if err != nil {
        return err
    }

    return os.WriteFile(filename, data, 0644) // Correct the WriteFile function call.
}

var buildCmd = &cobra.Command{
    Use:   "build [directory]",
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        directory := args[0]

        existingBuildCommand, err := getExistingBuildCommand(directory)
        if err != nil {
            fmt.Println("Error reading .ephemyral file:", err)
            return
        }

        if existingBuildCommand != "" {
            fmt.Println("Running existing build command:", existingBuildCommand)
            return
        }

        buildCommand, err := generateBuildCommand(directory)
        if err != nil {
            fmt.Println("Error generating build command:", err)
            return
        }

        if err := updateEphemyralBuildCommand(directory, buildCommand); err != nil {
            fmt.Println("Error updating .ephemyral file:", err)
            return
        }

        fmt.Println("Successfully generated and updated build command:", buildCommand)
        // Implement build command execution logic here.
    },
}

func init() {
    rootCmd.AddCommand(buildCmd)
}
