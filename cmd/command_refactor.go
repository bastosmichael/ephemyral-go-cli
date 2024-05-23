//go:build !lint
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

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func executeRefactorWithRetries(filePath, userPrompt, newFilePath string, convID uuid.UUID, retryCount int, retryDelay time.Duration, runBuild, runLint, runTest, runDocs bool) {
	// Reading the original file content
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Attempting the retries
	for retry := 0; retry <= retryCount; retry++ {
		if retry > 0 {
			fmt.Println("Retrying refactor... Attempt", retry)
			time.Sleep(retryDelay)
		}

		success := refactorFile(filePath, string(fileContent), userPrompt, newFilePath, convID)
		if success {
			// Run the additional commands if specified
			if runBuild && !runExistingBuildCommand(filePath, convID, retryCount, retryDelay) {
				restoreOriginalContent(filePath, fileContent)
				continue
			}

			if runLint && !runExistingLintCommand(filePath, convID, retryCount, retryDelay) {
				restoreOriginalContent(filePath, fileContent)
				continue
			}

			if runTest && !runExistingTestCommand(filePath, convID, retryCount, retryDelay) {
				restoreOriginalContent(filePath, fileContent)
				continue
			}

			if runDocs && !runExistingDocsCommand(filePath, convID, retryCount, retryDelay) {
				restoreOriginalContent(filePath, fileContent)
				continue
			}

			return
		}
	}

	// If all retries fail, restore the original content
	restoreOriginalContent(filePath, fileContent)
}

func refactorFile(filePath, fileContent, userPrompt, newFilePath string, convID uuid.UUID) bool {
	fullPrompt := fmt.Sprintf(
		RefactorPromptPattern,
		userPrompt,
		fileContent,
	)

	gpt4client.SetDebug(false)
	refactoredContent, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt, convID)
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

func restoreOriginalContent(filePath string, fileContent []byte) {
	err := os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		fmt.Println("Error restoring the original file content:", err)
	} else {
		fmt.Println("All retries failed, original content restored.")
	}
}

func runExistingBuildCommand(filePath string, convID uuid.UUID, retryCount int, retryDelay time.Duration) bool {
	directory, err := findEphemyralDirectory(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	existingCommand, err := getExistingCommandOrError(directory, "build")
	if err != nil {
		fmt.Println("Error reading existing build command:", err)
		return false
	}

	if existingCommand != "" {
		if err := retryExecution(directory, existingCommand, "build", convID, retryCount, retryDelay); err != nil {
			fmt.Println("Failed to execute build command:", err)
			return false
		} else {
			fmt.Println("Build command executed successfully.")
			return true
		}
	}
	return true
}

func runExistingLintCommand(filePath string, convID uuid.UUID, retryCount int, retryDelay time.Duration) bool {
	directory, err := findEphemyralDirectory(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	existingCommand, err := getExistingCommandOrError(directory, "lint")
	if err != nil {
		fmt.Println("Error reading existing lint command:", err)
		return false
	}

	if existingCommand != "" {
		if err := retryExecution(directory, existingCommand, "lint", convID, retryCount, retryDelay); err != nil {
			fmt.Println("Failed to execute lint command:", err)
			return false
		} else {
			fmt.Println("Lint command executed successfully.")
			return true
		}
	}
	return true
}

func runExistingTestCommand(filePath string, convID uuid.UUID, retryCount int, retryDelay time.Duration) bool {
	directory, err := findEphemyralDirectory(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	existingCommand, err := getExistingCommandOrError(directory, "test")
	if err != nil {
		fmt.Println("Error reading existing test command:", err)
		return false
	}

	if existingCommand != "" {
		if err := retryExecution(directory, existingCommand, "test", convID, retryCount, retryDelay); err != nil {
			fmt.Println("Failed to execute test command:", err)
			return false
		} else {
			fmt.Println("Test command executed successfully.")
			return true
		}
	}
	return true
}

func runExistingDocsCommand(filePath string, convID uuid.UUID, retryCount int, retryDelay time.Duration) bool {
	directory, err := findEphemyralDirectory(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	existingCommand, err := getExistingCommandOrError(directory, "docs")
	if err != nil {
		fmt.Println("Error reading existing docs command:", err)
		return false
	}

	if existingCommand != "" {
		if err := retryExecution(directory, existingCommand, "docs", convID, retryCount, retryDelay); err != nil {
			fmt.Println("Failed to execute docs command:", err)
			return false
		} else {
			fmt.Println("Docs command executed successfully.")
			return true
		}
	}
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

		convID := uuid.New()
		fmt.Println(convID)

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

		retryDelay := 2 * time.Second

		runBuild, _ := cmd.Flags().GetBool("build")
		runLint, _ := cmd.Flags().GetBool("lint")
		runTest, _ := cmd.Flags().GetBool("test")
		runDocs, _ := cmd.Flags().GetBool("docs")

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
					executeRefactorWithRetries(path, userPrompt, newFilePath, convID, defaultRetryCount, retryDelay, runBuild, runLint, runTest, runDocs)
				}
				return nil
			})

			if err != nil {
				fmt.Println("Error during directory walk:", err)
				return
			}
		} else {
			executeRefactorWithRetries(filePath, userPrompt, newFilePath, convID, defaultRetryCount, retryDelay, runBuild, runLint, runTest, runDocs)
		}
	},
}

func init() {
	refactorCmd.Flags().Int("retry", 3, "Number of retries for refactoring files")
	refactorCmd.Flags().Bool("build", false, "Run build command after refactoring")
	refactorCmd.Flags().Bool("lint", false, "Run lint command after refactoring")
	refactorCmd.Flags().Bool("test", false, "Run test command after refactoring")
	refactorCmd.Flags().Bool("docs", false, "Run docs command after refactoring")
	rootCmd.AddCommand(refactorCmd)
}
