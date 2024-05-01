// +build !lint

package cmd

import (
	"strings"
)

func filterOutCodeBlocks(content string) string {
	lines := strings.Split(content, "\n")
	filteredLines := make([]string, 0)
	for _, line := range lines {
		if !strings.Contains(line, "```") {
			filteredLines = append(filteredLines, line)
		}
	}
	// Join the filtered lines and ensure it ends with a newline
	result := strings.Join(filteredLines, "\n")
	if !strings.HasSuffix(result, "\n") {
		result += "\n" // Add a newline if it doesn't exist
	}
	return result
}
