// +build !lint

package cmd

import (
	"fmt"
	"strings"

	gpt4client "ephemyral/pkg"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func generateDocsCommand(directory string, convID uuid.UUID) (string, error) {
	filesList, err := getFileList(directory)
	if err != nil {
		return "", err
	}

	fullPrompt := DocsCommandPrompt + strings.Join(filesList, "\n")
	gpt4client.SetDebug(false)
	docsCommand, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt, convID)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(docsCommand) == "" {
		return "", fmt.Errorf("received empty docs command")
	}

	return docsCommand, nil
}

var docsCmd = &cobra.Command{
	Use:   "docs [directory]",
	Short: "Generate and execute a command to generate documentation, enhancing the codebase's maintainability.",
	Long:  "The 'docs' command creates a command to produce documentation (like a README or API documentation) for the project's files. It then updates the '.ephemyral' configuration file with the new command and executes it.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]

		// Get the retry count from flags, with a default of 3
		defaultRetryCount, err := cmd.Flags().GetInt("retry")
		if err != nil {
			fmt.Println("Error reading retry count:", err)
			return
		}

		convID := uuid.New()
		fmt.Println(convID)

		if err := executeCommandOfType(directory, "docs", convID, defaultRetryCount, retryDelay); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	docsCmd.Flags().Int("retry", 3, "Number of retries for generating and executing docs command")
	rootCmd.AddCommand(docsCmd)
}
