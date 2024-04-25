// common.go
package cmd

import "strings"

func filterOutCodeBlocks(content string) string {
    lines := strings.Split(content, "\n")
    filteredLines := make([]string, 0)
    for _, line := range lines {
        if !strings.Contains(line, "```") {
            filteredLines = append(filteredLines, line)
        }
    }
    return strings.Join(filteredLines, "\n")
}
