package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEphemyralFile(t *testing.T) {
	testFile := EphemyralFile{
		BuildCommand: "go build",
		TestCommand:  "go test",
		LintCommand:  "golint",
	}

	require.Equal(t, "go build", testFile.BuildCommand)
	require.Equal(t, "go test", testFile.TestCommand)
	require.Equal(t, "golint", testFile.LintCommand)
}
