// refactor.go
package cmd

import (
	gpt3client "ephemyral/pkg"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var refactorCmd = &cobra.Command{
	Use:   "refactor [file path] [prompt]",
	Short: "Refactor a given file based on a prompt",
	Long:  `This command refactors a given file by sending a prompt to an LLM and applying the suggested changes.`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		prompt := args[1]

		gpt3client.SetDebug(true)
		suggestion, err := gpt3client.GetLLMSuggestion(prompt)
		if err != nil {
			fmt.Println("Error getting suggestion from LLM:", err)
			return
		}

		// Read the file to be refactored.
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		// Apply the suggestion to the file content.
		// This is a simplification. You'll need to parse the suggestion and apply it accordingly.
		refactoredContent := string(fileContent) + "\n\n// Suggestion from LLM:\n" + suggestion

		// Write the refactored content back to the file.
		err = ioutil.WriteFile(filePath, []byte(refactoredContent), 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return
		}

		fmt.Println("File refactored successfully.")
	},
}

func init() {
	rootCmd.AddCommand(refactorCmd)
}
