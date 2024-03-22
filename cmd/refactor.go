

// refactor.go
package cmd

import (
    gpt3client "ephemyral/pkg"
    "fmt"
    "io/ioutil"
    "strings"

    "github.com/spf13/cobra"
)

var refactorCmd = &cobra.Command{
    Use:   "refactor [file path] [prompt] [new file path]",
    Short: "Refactor a given file based on a prompt and output to a new file",
    Long: `This command refactors a given file by sending a prompt to an LLM 
and applying the suggested changes entirely, replacing the file content.`,
    Args: cobra.MinimumNArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        filePath := args[0]
        userPrompt := args[1]
        newFilePath := ""
        if len(args) > 2 {
            newFilePath = args[2]
        }

        // Read the file to be refactored.
        fileContent, err := ioutil.ReadFile(filePath)
        if err != nil {
            fmt.Println("Error reading file:", err)
            return
        }

        // Prepare the prompt by merging the user prompt with the file content.
        // Here you would 'vectorize' or otherwise prepare your file content if necessary.
        // For simplicity, we're directly using the content as part of the prompt.
        // Adjust this to meet the API's limitations as needed.
        fullPrompt := fmt.Sprintf("Refactor the following code based on this instruction: '%s'\n\n%s", userPrompt, string(fileContent))

        gpt3client.SetDebug(true)
        refactoredContent, err := gpt3client.GetLLMSuggestion(fullPrompt)
        if err != nil {
            fmt.Println("Error getting suggestion from LLM:", err)
            return
        }

        // Check if the refactored content is valid or as expected. This is a simplification.
        // In a real scenario, you'd need more robust handling here.
        if strings.TrimSpace(refactoredContent) == "" {
            fmt.Println("Received empty refactored content.")
            return
        }

        // If a new file path was provided, write the refactored content to that file.
        // Otherwise, overwrite the original file.
        if newFilePath != "" {
            err = ioutil.WriteFile(newFilePath, []byte(refactoredContent), 0644)
            if err != nil {
                fmt.Println("Error writing file:", err)
                return
            }
        } else {
            err = ioutil.WriteFile(filePath, []byte(refactoredContent), 0644)
            if err != nil {
                fmt.Println("Error writing file:", err)
                return
            }
        }

        fmt.Println("File refactored successfully.")
    },
}

func init() {
    rootCmd.AddCommand(refactorCmd)
}