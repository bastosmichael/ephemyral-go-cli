//go:build !lint
// +build !lint

package cmd

import (
	gpt4client "ephemyral/pkg"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

		convID := new(uuid.UUID)

		var userPrompt string
		if len(args) > 1 && strings.TrimSpace(args[1]) != "" {
			userPrompt = args[1] // Use the provided prompt if available
		} else {
			userPrompt = "Create a new code file with meaningful content." // Default prompt
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("Error accessing specified path:", err)
			return
		}

		if fileInfo.IsDir() {
			// Handle directory - generate multiple files
			filesList, err := getFileList(filePath)
			if err != nil {
				fmt.Println("Error accessing files in directory:", err)
				return
			}
			for _, name := range filesList {
				generateNewFile(filepath.Join(filePath, name), userPrompt, convID)
			}
		} else {
			// Generate single file
			generateNewFile(filePath, userPrompt, convID)
		}
	},
}

func generateNewFile(filePath string, userPrompt string, convID *uuid.UUID) {
	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("File already exists:", filePath)
		return
	}

	fullPrompt := fmt.Sprintf("Create a new code file based on this prompt: %s.", userPrompt)

	newFileContent, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt, convID)
	if err != nil {
		fmt.Println("Error generating new file content:", err)
		return
	}

	err = os.WriteFile(filePath, []byte(newFileContent), 0644)
	if err != nil {
		fmt.Println("Error writing new file:", err)
		return
	}

	fmt.Println("New file created successfully:", filePath)
}

func init() {
	rootCmd.AddCommand(createCmd)
}
