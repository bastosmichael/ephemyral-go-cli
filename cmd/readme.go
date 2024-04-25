package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	gpt4client "ephemyral/pkg"

	"github.com/spf13/cobra"
)

func findReadmeFile(directory string) (string, error) {
	filesList, err := getFileList(directory) // Assuming getFileList returns a slice of file names
	if err != nil {
		return "", err
	}

	// Prepare the full prompt for LLM
	fullPrompt := FindReadmeCommandPrompt + strings.Join(filesList, "\n")
	gpt4client.SetDebug(false)                                          // Disable debug mode if needed
	readmeFile, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt) // Get the LLM response
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(readmeFile) == "" {
		return "", fmt.Errorf("received an empty response from LLM while identifying README file")
	}

	return readmeFile, nil
}

func getReadmeContent(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading README file: %w", err)
	}
	return string(content), nil
}

func writeReadmeContent(filePath, content string) error {
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing README file: %w", err)
	}
	return nil
}

var readmeCmd = &cobra.Command{
	Use:   "readme [directory]",
	Short: "Identify the README file location and update its content with machine learning.",
	Long:  "The 'readme' command identifies the location of the README file in a project directory and updates it based on the content of the project files. It then adds this location to the '.ephemyral' configuration file for future reference.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]

		// Get the retry count from flags, with a default of 3
		defaultRetryCount, err := cmd.Flags().GetInt("retry")
		if err != nil {
			fmt.Println("Error reading retry count:", err)
			return
		}

		var readmeFilePath string
		success := false

		// Find the README file with retries
		for i := 0; i < defaultRetryCount; i++ {
			readmeFile, err := findReadmeFile(directory)
			if err != nil {
				fmt.Println("Error finding README file:", err)
				time.Sleep(retryDelay) // Wait before retrying
			} else {
				refactoredReadmeFile := filterOutCodeBlocks(readmeFile)

				// Join the directory and refactoredReadmeFile to get the complete path
				// readmeFilePath = filepath.Join(directory, refactoredReadmeFile)
				readmeFilePath = refactoredReadmeFile

				// Retry logic
				// if fileExists(refactoredReadmeFile) {
				success = true
				fmt.Println("Successfully found README file:", readmeFilePath) // Success message
				break
				// } else {
				// 	fmt.Println("README file not found:", readmeFilePath)
				// 	time.Sleep(retryDelay) // Wait before retrying
				// }
			}
		}

		if !success {
			fmt.Println("Failed to find README file after retries.")
			return
		}

		// Read the README content
		readmeContent, err := getReadmeContent(readmeFilePath)
		if err != nil {
			fmt.Println("Error reading README file:", err)
			return
		}

		// Prompt LLM to update the README content
		updatedReadme, err := gpt4client.GetGPT4ResponseWithPrompt(fmt.Sprintf("Please update this README:\n%s", readmeContent))
		if err != nil {
			fmt.Println("Error getting updated README from LLM:", err)
			return
		}

		// Write the updated README content
		if err := writeReadmeContent(readmeFilePath, updatedReadme); err != nil {
			fmt.Println("Error writing updated README:", err)
			return
		}

		// Update the .ephemyral file with the README location
		if err := updateEphemyralFile(directory, "readme", readmeFilePath); err != nil {
			fmt.Println("Error updating .ephemyral file:", err)
			return
		}

		fmt.Println("Successfully updated README and .ephemyral with its location.") // Success message
	},
}

func init() {
	readmeCmd.Flags().Int("retry", 3, "Number of retries for finding and updating README file")
	rootCmd.AddCommand(readmeCmd)
}
