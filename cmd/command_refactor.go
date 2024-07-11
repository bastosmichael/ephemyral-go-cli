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
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	restoreContent := func() {
		err := os.WriteFile(filePath, fileContent, 0644)
		if err != nil {
			fmt.Println("Error restoring the original file content:", err)
		} else {
			fmt.Println("All retries failed, original content restored.")
		}
	}

	for retry := 0; retry <= retryCount; retry++ {
		if retry > 0 {
			fmt.Println("Retrying refactor... Attempt", retry)
			time.Sleep(retryDelay)
		}

		if refactorFile(filePath, string(fileContent), userPrompt, newFilePath, convID) {
			if (runBuild && !runCommand("build", filePath, convID, retryCount, retryDelay)) ||
				(runLint && !runCommand("lint", filePath, convID, retryCount, retryDelay)) ||
				(runTest && !runCommand("test", filePath, convID, retryCount, retryDelay)) ||
				(runDocs && !runCommand("docs", filePath, convID, retryCount, retryDelay)) {
				restoreContent()
				continue
			}
			return
		}
	}
	restoreContent()
}

func refactorFile(filePath, fileContent, userPrompt, newFilePath string, convID uuid.UUID) bool {
	fullPrompt := fmt.Sprintf(RefactorPromptPattern, userPrompt, fileContent)
	gpt4client.SetDebug(false)
	refactoredContent, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt, convID)
	if err != nil {
		fmt.Println("Error from LLM:", err)
		return false
	}

	filteredContent := filterOutCodeBlocks(refactoredContent)
	if strings.TrimSpace(filteredContent) == "" {
		fmt.Println("Insufficient content from LLM after filtering.")
		return false
	}

	targetFilePath := filePath
	if newFilePath != "" {
		targetFilePath = filepath.Join(newFilePath, filepath.Base(filePath))
	}

	if err := os.WriteFile(targetFilePath, []byte(filteredContent), 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return false
	}
	fmt.Println("File refactored successfully:", targetFilePath)
	return true
}

var refactorCmd = &cobra.Command{
	Use:   "refactor [file path] [prompt] [new file path]",
	Short: "Utilize an advanced LLM to refactor a give file or all files in a provided directory based on prompts, outputting, building and testing the improved code.",
	Long: `This command refactors a given file or all files in a directory by sending a prompt to an LLM 
and applying the suggested changes, replacing the file content or creating new files in the specified directory.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath, userPrompt, newFilePath := args[0], DefaultRefactorPrompt, ""
		if len(args) > 1 && strings.TrimSpace(args[1]) != "" {
			userPrompt = args[1]
		}
		if len(args) > 2 {
			newFilePath = args[2]
		}

		convID := uuid.New()
		fmt.Println(convID)

		retryCount, err := cmd.Flags().GetInt("retry")
		if err != nil {
			fmt.Println("Error reading retry count:", err)
			return
		}

		retryDelay := 2 * time.Second
		runBuild, _ := cmd.Flags().GetBool("build")
		runLint, _ := cmd.Flags().GetBool("lint")
		runTest, _ := cmd.Flags().GetBool("test")
		runDocs, _ := cmd.Flags().GetBool("docs")

		if fileInfo, err := os.Stat(filePath); err != nil {
			fmt.Println("Error accessing specified path:", err)
		} else if fileInfo.IsDir() {
			filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return err
				}
				executeRefactorWithRetries(path, userPrompt, newFilePath, convID, retryCount, retryDelay, runBuild, runLint, runTest, runDocs)
				return nil
			})
		} else {
			executeRefactorWithRetries(filePath, userPrompt, newFilePath, convID, retryCount, retryDelay, runBuild, runLint, runTest, runDocs)
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
