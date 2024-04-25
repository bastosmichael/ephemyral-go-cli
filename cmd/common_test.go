// common_test.go
package cmd

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestEphemyralFile(t *testing.T) {
    testFile := EphemyralFile{
        BuildCommand: "go build",
        TestCommand:  "go test",
        LintCommand:  "golint",
    }

    assert.Equal(t, testFile.BuildCommand, "go build")
    assert.Equal(t, testFile.TestCommand, "go test")
    assert.Equal(t, testFile.LintCommand, "golint")
}
