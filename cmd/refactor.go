// refactor.go
package cmd

import (
    gpt4client "ephemyral/pkg"
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

        // Prepare a precise prompt to the LLM.
        fullPrompt := fmt.Sprintf("Analyze the following code and return only the completely refactored or optimized code based on this instruction: '%s' Provide the refactored or optimized version only. Do not include any additional text or unchanged code.\n\n%s", userPrompt, string(fileContent))

        gpt4client.SetDebug(true)
        refactoredContent, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt)
        if err != nil {
            fmt.Println("Error getting suggestion from LLM:", err)
            return
        }

        // Filter out any lines containing triple backticks
        refactoredContent = filterOutCodeBlocks(refactoredContent)

        // Validation: Ensure the returned content is strictly code and it's only the changed parts.
        if !isCode(refactoredContent) || strings.TrimSpace(refactoredContent) == "" {
            fmt.Println("Invalid or insufficient content received. Expected specific code changes only.")
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

func isCode(content string) bool {
    // Simple check for code structure; adjust according to your needs.
    return strings.Contains(content, "func") || strings.Contains(content, "import")
}

func filterOutCodeBlocks(content string) string {
    lines := strings.Split(content, "\n")
    filteredLines := []string{}
    for _, line := range lines {
        if !strings.Contains(line, "```") {
            filteredLines = append(filteredLines, line)
        }
    }
    return strings.Join(filteredLines, "\n")
}

func init() {
    rootCmd.AddCommand(refactorCmd)
}
