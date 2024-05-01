// +build !lint

package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// getFileList retrieves a list of all non-directory file names in the specified directory and its subdirectories,
// skipping specified directories like ".git".
func getFileList(directory string) ([]string, error) {
	var filesList []string

	// Read the files in the root directory.
	rootFiles, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	// Add non-directory files from the root directory to the list.
	for _, file := range rootFiles {
		if !file.IsDir() {
			filesList = append(filesList, file.Name())
		}
	}

	// Walk through the directory and its subdirectories.
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.HasSuffix(info.Name(), ".git") {
				return filepath.SkipDir
			}
		} else {
			relativePath, err := filepath.Rel(directory, path)
			if err != nil {
				return err
			}
			filesList = append(filesList, relativePath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filesList, nil
}
