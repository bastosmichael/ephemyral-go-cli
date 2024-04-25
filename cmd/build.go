// build.go
package cmd

import (
    gpt4client "ephemyral/pkg"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"
)

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

var buildCmd = &cobra.Command{
    Use:   "build [directory]",
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        directory := args[0]

        existingBuildCommand, err := getExistingCommand(directory, "build")
        if err != nil {
            fmt.Println("Error reading .ephemyral file:", err)
            return
        }

        if existingBuildCommand != "" {
            fmt.Println("Running existing build command:", existingBuildCommand)
            if err := executeCommand(directory, existingBuildCommand); err != nil {
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

        if err := updateEphemyralCommand(directory, "build", refactoredBuildCommand); err != nil {
            fmt.Println("Error updating .ephemyral file:", err)
            return
        }

        fmt.Println("Successfully generated and updated build command:", refactoredBuildCommand)
        
        // Execute the new build command.
        if err := executeCommand(directory, buildCommand); err != nil {
            fmt.Println("Error executing new build command:", err)
        }
    },
}

func init() {
    rootCmd.AddCommand(buildCmd)
}
