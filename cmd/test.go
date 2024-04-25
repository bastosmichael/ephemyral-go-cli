// test.go
package cmd

import (
    "fmt"
    "strings"

    "github.com/spf13/cobra"
    gpt4client "ephemyral/pkg"
)

// generateTestCommand generates a test command based on the file structure in the directory.
func generateTestCommand(directory string) (string, error) {
    filesList, err := getFileList(directory)
    if err != nil {
        return "", err
    }

    // Join all file names into a single prompt.
    fullPrompt := "Based on the following file list, provide the simplest command line required to test these files. The command must be in a single line and contain no extra text or commentary:\n" +
                  strings.Join(filesList, "\n")

    gpt4client.SetDebug(true)
    testCommand, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt)
    if err != nil {
        return "", err
    }

    if strings.TrimSpace(testCommand) == "" {
        return "", fmt.Errorf("received empty test command")
    }

    return testCommand, nil
}

var testCmd = &cobra.Command{
    Use:   "test [directory]",
    Short: "Generate and run a test command for the specified directory",
    Long:  "The 'test' command retrieves the test command specified in the '.ephemyral' configuration file within the given directory. If no test command is specified, the command generates a test command based on the structure of the project's files. It then updates the '.ephemyral' file with the new test command and executes it. Use this command to ensure your project runs the appropriate tests before deployment or other build-related tasks.",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        directory := args[0]

        existingTestCommand, err := getExistingCommand(directory, "test")
        if err != nil {
            fmt.Println("Error reading .ephemyral file:", err)
            return
        }

        if existingTestCommand != "" {
            fmt.Println("Running existing test command:", existingTestCommand)
            if err := executeCommand(directory, existingTestCommand); err != nil {
                fmt.Println("Error executing new test command:", err)
                return
            }
            return
        }

        testCommand, err := generateTestCommand(directory)

        if err != nil {
            fmt.Println("Error generating test command:", err)
            return
        }

        refactoredTestCommand := filterOutCodeBlocks(testCommand)

        if err := updateEphemyralCommand(directory, "test", refactoredTestCommand); err != nil {
            fmt.Println("Error updating .ephemyral file:", err)
            return
        }

        fmt.Println("Successfully generated and updated test command:", refactoredTestCommand)
        
        // Execute the new test command.
        if err := executeCommand(directory, testCommand); err != nil {
            fmt.Println("Error executing new test command:", err)
        }
    },
}

func init() {
    rootCmd.AddCommand(testCmd)
}
