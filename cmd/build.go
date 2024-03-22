// build.go
package cmd

import (
    gpt4client "ephemyral/pkg"
    "fmt"
    "io/ioutil"
    "strings"

    "github.com/spf13/cobra"
)

// Extracting the core functionality into a separate function for testing.
func GenerateBuildCommand(directory string) (string, error) {
    files, err := ioutil.ReadDir(directory)
    if err != nil {
        return "", err // Return error to be handled/tested.
    }

    var fullPrompt strings.Builder
    for _, file := range files {
        fileContent, err := ioutil.ReadFile(directory + "/" + file.Name())
        if err != nil {
            return "", err // Return error to be handled/tested.
        }
        fullPrompt.WriteString(string(fileContent) + "\n\n")
    }

    gpt4client.SetDebug(true)
    buildCommand, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt.String())
    if err != nil {
        return "", err // Return error to be handled/tested.
    }

    if strings.TrimSpace(buildCommand) == "" {
        return "", fmt.Errorf("received empty build command") // Using error to signify empty command.
    }

    // Typically, you'd run the build command here, but for testing, we'll just return it.
    return buildCommand, nil
}

var buildCmd = &cobra.Command{
    Use:   "build [directory]",
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        buildCommand, err := GenerateBuildCommand(args[0])
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        
        fmt.Println("Successfully generated build command:", buildCommand)
        // Implement build command execution logic.
    },
}

func init() {
    rootCmd.AddCommand(buildCmd)
}
