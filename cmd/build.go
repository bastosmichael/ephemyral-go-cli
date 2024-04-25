// build.go
package cmd

import (
    gpt4client "ephemyral/pkg"
    "fmt"
    "os"
    "os/exec"
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

// generateBuildCommand generates a build command by listing all files in the root directory and subdirectories.
func generateBuildCommand(directory string) (string, error) {
    var filesList []string

    // First, get the files in the root directory.
    rootFiles, err := os.ReadDir(directory) // Read the content of the root directory.
    if err != nil {
        return "", err
    }

    for _, file := range rootFiles {
        if !file.IsDir() { // Only add files, not directories.
            filesList = append(filesList, file.Name()) // Add the file names to the list.
        }
    }

    // Walk through the directory and its subdirectories to get all other file names.
    err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Ignore the .git directories.
        if info.IsDir() && strings.HasSuffix(info.Name(), ".git") {
            return filepath.SkipDir // Skip .git directories.
        }

        if info.IsDir() && filepath.Base(path) == directory {
            return nil // Skip the root directory as we've already processed its files.
        }

        if !info.IsDir() { // Only add files, not directories.
            relativePath, err := filepath.Rel(directory, path) // Get relative paths.
            if err != nil {
                return err
            }
            filesList = append(filesList, relativePath)
        }

        return nil
    }) // Correctly close this anonymous function.

    if err != nil {
        return "", err
    }

    // Join all file names into a single prompt.
    fullPrompt := "Based on the following file list, provide the simplest command line required to build these files. The command must be in a single line and contain no extra text or commentary:\n" +
                  strings.Join(filesList, "\n")

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

// executeBuildCommand executes the build command using os/exec.
func executeBuildCommand(directory, buildCommand string) error {
    // Create a new command from the buildCommand string.
    cmd := exec.Command("bash", "-c", buildCommand)

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
            if err := executeBuildCommand(directory, existingBuildCommand); err != nil {
                fmt.Println("Error executing new build command:", err)
                return
            }
            return
        }

        buildCommand, err := generateBuildCommand(directory)

        if err != nil {
            fmt.Println("Error generating build command:", err)
            return
        }

        refactoredBuildCommand := filterOutCodeBlocks(buildCommand)

        if err := updateEphemyralBuildCommand(directory, refactoredBuildCommand); err != nil {
            fmt.Println("Error updating .ephemyral file:", err)
            return
        }

        fmt.Println("Successfully generated and updated build command:", refactoredBuildCommand)
        
        // Execute the new build command.
        if err := executeBuildCommand(directory, buildCommand); err != nil {
            fmt.Println("Error executing new build command:", err)
        }
    },
}

func init() {
    rootCmd.AddCommand(buildCmd)
}
