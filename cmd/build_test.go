// build_test.go
package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

// executeCommand executes a cobra.Command with given arguments, and returns the output
func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err := root.Execute()
	return buf.String(), err
}

// TestBuildCommand tests the build command
func TestBuildCommand(t *testing.T) {
	output, err := executeCommand(rootCmd, "build")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := "Building the project...\n" // Update this based on your actual command output
	if output != expected {
		t.Errorf("Expected output %q but got %q", expected, output)
	}
}
