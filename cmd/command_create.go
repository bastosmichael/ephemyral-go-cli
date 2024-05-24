//go:build !lint
// +build !lint

package cmd

import (
	gpt4client "ephemyral/pkg"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [file path] [prompt]",
	Short: "Employ a language model to generate new code files based on a natural language prompt. If the file path is a directory, it generates multiple AI-crafted files.",
	Long: `This command generates a new code file based on a given prompt. 
If the file path is a directory, it uses a query to determine the file names and creates new files based on the provided prompt.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		convID := uuid.New()
		fmt.Println(convID)

		userPrompt := "Create a new code file with meaningful content."
		if len(args) > 1 && strings.TrimSpace(args[1]) != "" {
			userPrompt = args[1] // Use the provided prompt if available
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("Error accessing specified path:", err)
			return
		}

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

		if fileInfo.IsDir() {
			filesList, err := getFileList(filePath)
			if err != nil {
				fmt.Println("Error accessing files in directory:", err)
				return
			}
			for _, name := range filesList {
				generateNewFile(filepath.Join(filePath, name), userPrompt, convID, retryCount, retryDelay, runBuild, runLint, runTest, runDocs)
			}
		} else {
			generateNewFile(filePath, userPrompt, convID, retryCount, retryDelay, runBuild, runLint, runTest, runDocs)
		}
	},
}

func generateNewFile(filePath string, userPrompt string, convID uuid.UUID, retryCount int, retryDelay time.Duration, runBuild, runLint, runTest, runDocs bool) {
	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("File already exists:", filePath)
		return
	}

	fullPrompt := fmt.Sprintf("Create a new code file based on this prompt: %s.", userPrompt)

	for retry := 0; retry <= retryCount; retry++ {
		if retry > 0 {
			fmt.Println("Retrying creation... Attempt", retry)
			time.Sleep(retryDelay)
		}

		newFileContent, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt, convID)
		if err != nil {
			fmt.Println("Error generating new file content:", err)
			continue
		}

		if strings.TrimSpace(newFileContent) == "" {
			fmt.Println("Invalid or insufficient content received.")
			continue
		}

		if err := os.WriteFile(filePath, []byte(newFileContent), 0644); err != nil {
			fmt.Println("Error writing new file:", err)
			continue
		}

		fmt.Println("New file created successfully:", filePath)

		if (runBuild && !runCommand("build", filePath, convID, retryCount, retryDelay)) ||
			(runLint && !runCommand("lint", filePath, convID, retryCount, retryDelay)) ||
			(runTest && !runCommand("test", filePath, convID, retryCount, retryDelay)) ||
			(runDocs && !runCommand("docs", filePath, convID, retryCount, retryDelay)) {
			os.Remove(filePath)
			continue
		}

		return
	}

	fmt.Println("All retries failed. File creation was unsuccessful.")
}

func init() {
	createCmd.Flags().Int("retry", 3, "Number of retries for creating files")
	createCmd.Flags().Bool("build", false, "Run build command after creating files")
	createCmd.Flags().Bool("lint", false, "Run lint command after creating files")
	createCmd.Flags().Bool("test", false, "Run test command after creating files")
	createCmd.Flags().Bool("docs", false, "Run docs command after creating files")
	rootCmd.AddCommand(createCmd)
}
