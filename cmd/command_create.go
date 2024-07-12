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
	"github.com/muesli/termenv"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

var (
	automode      bool
	maxIterations int
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
		printPanel(convID.String(), "Conversation ID", "cyan")

		userPrompt := "Create a new code file with meaningful content."
		if len(args) > 1 && strings.TrimSpace(args[1]) != "" {
			userPrompt = args[1] // Use the provided prompt if available
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			printPanel(fmt.Sprintf("Error accessing specified path: %s", err), "Error", "red")
			return
		}

		retryCount, err := cmd.Flags().GetInt("retry")
		if err != nil {
			printPanel(fmt.Sprintf("Error reading retry count: %s", err), "Error", "red")
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
				printPanel(fmt.Sprintf("Error accessing files in directory: %s", err), "Error", "red")
				return
			}
			for _, name := range filesList {
				generateNewFile(filepath.Join(filePath, name), userPrompt, convID, retryCount, retryDelay, runBuild, runLint, runTest, runDocs)
			}
		} else {
			generateNewFile(filePath, userPrompt, convID, retryCount, retryDelay, runBuild, runLint, runTest, runDocs)
		}

		if automode {
			runAutomode(maxIterations)
		}
	},
}

func printPanel(content, title, color string) {
	p := termenv.ColorProfile()
	panel := termenv.String(content).Foreground(p.Color(color))
	fmt.Printf("[%s] %s: %s\n", title, panel.String(), content)
}

func runAutomode(iterations int) {
	for i := 0; i < iterations; i++ {
		fmt.Println("Automode iteration:", i+1)
		// Perform actions for each iteration
		time.Sleep(1 * time.Second) // Placeholder for actual work
	}
	printPanel("Automode completed.", "Automode", "green")
}

func generateNewFile(filePath string, userPrompt string, convID uuid.UUID, retryCount int, retryDelay time.Duration, runBuild, runLint, runTest, runDocs bool) {
	existingContent := ""
	if _, err := os.Stat(filePath); err == nil {
		content, err := os.ReadFile(filePath)
		if err != nil {
			printPanel(fmt.Sprintf("Error reading existing file content: %s", err), "Error", "red")
			return
		}
		existingContent = string(content)
	}

	fullPrompt := fmt.Sprintf("Create a new code file based on this prompt: %s.", userPrompt)

	retryErr := retryWithDelay(retryCount, retryDelay, func() error {
		newFileContent, err := gpt4client.GetGPT4ResponseWithPrompt(fullPrompt, convID)
		if err != nil {
			return fmt.Errorf("error generating new file content: %w", err)
		}

		if strings.TrimSpace(newFileContent) == "" {
			return fmt.Errorf("invalid or insufficient content received")
		}

		filteredContent := filterOutCodeBlocks(newFileContent)
		if strings.TrimSpace(filteredContent) == "" {
			return fmt.Errorf("filtered content is empty")
		}

		if existingContent != "" {
			diff, err := generateAndApplyDiff(existingContent, filteredContent, filePath)
			if err != nil {
				return fmt.Errorf("error applying diff: %w", err)
			}
			printPanel(fmt.Sprintf("Changes applied:\n%s", diff), "Diff", "green")
		} else {
			if err := os.WriteFile(filePath, []byte(filteredContent), 0644); err != nil {
				return fmt.Errorf("error writing new file: %w", err)
			}
			printPanel("New file created successfully: "+filePath, "Success", "green")
		}

		if (runBuild && !runCommand("build", filePath, convID, retryCount, retryDelay)) ||
			(runLint && !runCommand("lint", filePath, convID, retryCount, retryDelay)) ||
			(runTest && !runCommand("test", filePath, convID, retryCount, retryDelay)) ||
			(runDocs && !runCommand("docs", filePath, convID, retryCount, retryDelay)) {
			os.Remove(filePath)
			return fmt.Errorf("command execution failed, file creation was unsuccessful")
		}

		return nil
	})

	if retryErr != nil {
		printPanel("All retries failed. File creation was unsuccessful.", "Error", "red")
	}
}

func retryWithDelay(attempts int, delay time.Duration, function func() error) error {
	var err error
	for i := 0; i <= attempts; i++ {
		if i > 0 {
			printPanel(fmt.Sprintf("Retrying... Attempt %d", i), "Retry", "yellow")
			time.Sleep(delay)
		}
		if err = function(); err == nil {
			return nil
		}
	}
	return err
}

func generateAndApplyDiff(originalContent, newContent, path string) (string, error) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(originalContent, newContent, false)
	unifiedDiff := dmp.DiffPrettyText(diffs)

	if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
		return "", err
	}

	return unifiedDiff, nil
}

func init() {
	createCmd.Flags().Int("retry", 3, "Number of retries for creating files")
	createCmd.Flags().Bool("build", false, "Run build command after creating files")
	createCmd.Flags().Bool("lint", false, "Run lint command after creating files")
	createCmd.Flags().Bool("test", false, "Run test command after creating files")
	createCmd.Flags().Bool("docs", false, "Run docs command after creating files")
	createCmd.Flags().BoolVar(&automode, "automode", false, "Run in automode")
	createCmd.Flags().IntVar(&maxIterations, "max-iterations", 25, "Maximum iterations for automode")
	rootCmd.AddCommand(createCmd)
}
