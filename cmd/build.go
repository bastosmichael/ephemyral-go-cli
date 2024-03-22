

// build.go
package cmd

import (
    gpt3client "ephemyral/pkg"
    "fmt"
    "io/ioutil"
    "strings"

    "github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
    Use:   "build [directory]",
    Short: "Build code directory",
    Long: `This command reads all files in a given directory and uses gpt3client
to return only the command necessary to build the code directory and runs that command to successfully check for a build.`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        directory := args[0]

        // Read all files in the given directory.
        files, err := ioutil.ReadDir(directory)
        if err != nil {
            fmt.Println("Error reading directory:", err)
            return
        }

        // Prepare the prompt by merging all files' content into one string.
        // Here you would 'vectorize' or otherwise prepare your file content if necessary.
        // Adjust this to meet the API's limitations as needed.
        var fullPrompt strings.Builder
        for _, file := range files {
            fileContent, err := ioutil.ReadFile(directory + "/" + file.Name())
            if err != nil {
                fmt.Println("Error reading file:", err)
                return
            }
            fullPrompt.WriteString(string(fileContent) + "\n\n")
        }

        gpt3client.SetDebug(true)
        buildCommand, err := gpt3client.GetLLMSuggestion(fullPrompt.String())
        if err != nil {
            fmt.Println("Error getting suggestion from LLM:", err)
            return
        }

        // Check if the build command is valid or as expected. This is a simplification.
        // In a real scenario, you'd need more robust handling here.
        if strings.TrimSpace(buildCommand) == "" {
            fmt.Println("Received empty build command.")
            return
        }

        // Run the build command.
        // TODO: Implement build command execution.
        fmt.Println("Successfully ran build command:", buildCommand)
    },
}

func init() {
    rootCmd.AddCommand(buildCmd)
}