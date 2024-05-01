// +build !lint

package cmd

import (
	"fmt"
	"strings"

	gpt4client "ephemyral/pkg"

	"github.com/spf13/cobra"
)

func generateTestCommand(directory string) (string, error) {
	filesList, err := getFileList(directory)
	if err != nil {
		return "", err
	}

	fullPrompt := TestCommandPrompt + strings.Join(filesList, "\n")

	gpt4client.SetDebug(false)
	return gpt4client.GetGPT4ResponseWithPrompt(fullPrompt)
}

var testCmd = &cobra.Command{
	Use:   "test [directory]",
	Short: "Deploy AI models to generate and run optimized test commands for the specified directory, enhancing test accuracy and efficiency.",
	Long:  "The 'test' command generates a testing command based on the structure of the project's files. It then updates the '.ephemyral' configuration file with the new testing command and executes it. Use this command to ensure your project adheres to testing standards and is free from test errors.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]

		// Get the retry count from flags, with a default of 3
		defaultRetryCount, err := cmd.Flags().GetInt("retry")
		if err != nil {
			fmt.Println("Error reading retry count:", err)
			return
		}
		
		if err := executeCommandOfType(directory, "test", defaultRetryCount, retryDelay); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	testCmd.Flags().Int("retry", 3, "Number of retries for generating and executing test command")
	rootCmd.AddCommand(testCmd)
}
