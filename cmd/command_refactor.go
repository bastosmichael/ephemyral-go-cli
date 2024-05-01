// +build !lint

package cmd

import (
	gpt4client "ephemyral/pkg"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func executeRefactorWithRetries(filePath, userPrompt, newFilePath string, retryCount int, retryDelay time.Duration) {
	for retry := 0; retry <= retryCount; retry++ {
		if retry > 0 {
			fmt.Println("Retrying refactor... Attempt", retry)
			time.Sleep(retryDelay)
		}

		success := refactorFile(filePath, userPrompt, newFilePath)
		if success {
			break
		}
	}
}

func refactorFile(filePath, userPrompt, newFilePath string) bool {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return false
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
		return false
	}

	refactoredContent = filterOutCodeBlocks(refactoredContent)

	if strings.TrimSpace(refactoredContent) == "" {
		fmt.Println("Invalid or insufficient content received. Expected specific code changes only.")
		return false
	}

	targetFilePath := filePath
	if newFilePath != "" {
		if fileInfo, _ := os.Stat(newFilePath); fileInfo != nil && fileInfo.IsDir() {
			targetFilePath = filepath.Join(newFilePath, filepath.Base(filePath))
		} else {
			targetFilePath = newFilePath
		}
	}

	err = os.WriteFile(targetFilePath, []byte(refactoredContent), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return false
	}

	fmt.Println("File refactored successfully:", targetFilePath)
	return true
}

var refactorCmd = &cobra.Command{
	Use:   "refactor [file path] [prompt] [new file path]",
	Short: "Utilize an advanced LLM to refactor given files or all files in a directory based on a prompt, outputting the improved code to a new location.",
	Long: `This command refactors a given file or all files in a directory by sending a prompt to an LLM 
and applying the suggested changes, replacing the file content or creating new files in the specified directory.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		userPrompt := DefaultRefactorPrompt
		if len(args) > 1 && strings.TrimSpace(args[1]) != "" {
			userPrompt = args[1]
		}

		newFilePath := ""
		if len(args) > 2 {
			newFilePath = args[2]
		}

		defaultRetryCount, err := cmd.Flags().GetInt("retry")
		if err != nil {
			fmt.Println("Error reading retry count:", err)
			return
		}

		retryDelay := 2 * time.Second // Example retry delay

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("Error accessing specified path:", err)
			return
		}

		if fileInfo.IsDir() {
			err := filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					fmt.Println("Error accessing path:", path, err)
					return err
				}
				if !info.IsDir() {
					executeRefactorWithRetries(path, userPrompt, newFilePath, defaultRetryCount, retryDelay)
				}
				return nil
			})

			if err != nil {
				fmt.Println("Error during directory walk:", err)
				return
			}
		} else {
			executeRefactorWithRetries(filePath, userPrompt, newFilePath, defaultRetryCount, retryDelay)
		}
	},
}

func init() {
	refactorCmd.Flags().Int("retry", 3, "Number of retries for refactoring files")
	rootCmd.AddCommand(refactorCmd)
}
