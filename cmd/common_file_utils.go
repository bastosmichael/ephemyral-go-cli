//go:build !lint
// +build !lint

package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

const gitDirSuffix = ".git"

// fileExists checks if a file exists at the specified path.
func fileExists(filename string) bool {
	return !os.IsNotExist(checkFileStat(filename))
}

// checkFileStat returns the os.Stat result for a given filename.
func checkFileStat(filename string) error {
	_, err := os.Stat(filename)
	return err
}

// getFileList retrieves a list of all non-directory file names in the specified directory and its subdirectories,
// skipping specified directories like ".git".
func getFileList(directory string) ([]string, error) {
	filesList, err := readRootFiles(directory)
	if err != nil {
		return nil, err
	}

	subFilesList, err := walkDirectory(directory)
	if err != nil {
		return nil, err
	}

	return append(filesList, subFilesList...), nil
}

// readRootFiles reads the files in the root directory and adds non-directory files to the list.
func readRootFiles(directory string) ([]string, error) {
	var filesList []string
	rootFiles, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range rootFiles {
		if !file.IsDir() {
			filesList = append(filesList, file.Name())
		}
	}
	return filesList, nil
}

// walkDirectory walks through the directory and its subdirectories to collect file paths.
func walkDirectory(directory string) ([]string, error) {
	var filesList []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		return walkFunc(path, info, err, directory, &filesList)
	})
	if err != nil {
		return nil, err
	}
	return filesList, nil
}

// walkFunc processes each file or directory encountered during the walk.
func walkFunc(path string, info os.FileInfo, err error, directory string, filesList *[]string) error {
	if err != nil {
		return err
	}

	if shouldSkipDir(info) {
		return filepath.SkipDir
	}

	if !info.IsDir() {
		relativePath, err := getRelativePath(directory, path)
		if err != nil {
			return err
		}
		*filesList = append(*filesList, relativePath)
	}
	return nil
}

// shouldSkipDir determines if a directory should be skipped.
func shouldSkipDir(info os.FileInfo) bool {
	return info.IsDir() && strings.HasSuffix(info.Name(), gitDirSuffix)
}

// getRelativePath calculates the relative path of a file from the base directory.
func getRelativePath(base, target string) (string, error) {
	return filepath.Rel(base, target)
}
