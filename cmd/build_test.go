package cmd_test

import (
    "ephemyral/cmd"
    "testing"
)

func TestGenerateBuildCommand(t *testing.T) {
    // This is a basic example. Your test setup might need adjustments based on your environment.
    dir := "testData" // Assuming testData contains mock files for testing.
    expectedCommand := "echo 'build command'" // This is a mock expected output.

    command, err := cmd.GenerateBuildCommand(dir)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    if command != expectedCommand {
        t.Errorf("Expected %s, got %s", expectedCommand, command)
    }
}
