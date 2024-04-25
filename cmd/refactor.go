// refactor.go
package cmd

import (
    gpt4client "ephemyral/pkg"
    "fmt"
    "io/fs"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"
)

var refactorCmd = &cobra.Command{
    Use:   "refactor [file path] [prompt] [new file path]",
    Short: "Utilize an advanced LLM to refactor given files or all files in a directory based on a prompt, outputting the improved code to a new location.",
    Long: `This command refactors a given file or all files in a directory by sending a prompt to an LLM 
and applying the suggested changes, replacing the file content or creating new files in the specified directory.`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        filePath := args[0]

        // Default prompt if userPrompt is not provided
        userPrompt := DefaultRefactorPrompt
        
        if len(args) > 1 && strings.TrimSpace(args[1]) != "" {
            userPrompt = args[1] // Use the provided prompt if available
        }

        newFilePath := ""
        if len(args) > 2 {
            newFilePath = args[2]
        }

        // Check if the provided filePath is a directory
        fileInfo, err := os.Stat(filePath)
        if err != nil {
            fmt.Println("Error accessing specified path:", err)
            return
        }

        if fileInfo.IsDir() {
            // Handle directory
            filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
                if err != nil {
                    fmt.Println("Error accessing path during walk:", path, err)
                    return err
                }
                if !info.IsDir() {
                    refactorFile(path, userPrompt, newFilePath)
                }
                return nil
            })
        } else {
            // Handle single file
            refactorFile(filePath, userPrompt, newFilePath)
        }
    },
}

func refactorFile(filePath, userPrompt, newFilePath string) {
    fileContent, err := ioutil.ReadFile(filePath)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

    fullPrompt := fmt.Sprintf(
        RefactorPromptPattern,
        userPrompt,
        string(fileContent),
    )

    gpt4client.SetDebug(false)
    refactoredContent, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt)
    if err != nil {
        fmt.Println("Error getting suggestion from LLM:", err)
        return
    }

    // Filter out lines with triple backticks
    refactoredContent = filterOutCodeBlocks(refactoredContent)

    if strings.TrimSpace(refactoredContent) == "" {
        fmt.Println("Invalid or insufficient content received. Expected specific code changes only.")
        return
    }

    targetFilePath := filePath
    if newFilePath != "" {
        if fileInfo, _ := os.Stat(newFilePath); fileInfo != nil && fileInfo.IsDir() {
            targetFilePath = filepath.Join(newFilePath, filepath.Base(filePath))
        } else {
            targetFilePath = newFilePath
        }
    }

    err = ioutil.WriteFile(targetFilePath, []byte(refactoredContent), 0644)
    if err != nil {
        fmt.Println("Error writing file:", err)
        return
    }

    fmt.Println("File refactored successfully:", targetFilePath)
}

func init() {
    rootCmd.AddCommand(refactorCmd)
}
