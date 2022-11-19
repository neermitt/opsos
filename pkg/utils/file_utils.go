package utils

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// FileExists checks if file exists and is not a directory
func FileExists(filename string) bool {
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) || err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

// IsDir Checks if the filename exists and is a directory
func IsDir(filename string) bool {
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) || err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// PrintOrWriteToFile converts the provided value to given format and writes it to the specified file
func PrintOrWriteToFile(format string, filePath string, data any, fileMode os.FileMode) error {
	var w io.Writer = os.Stdout
	if filePath != "" {
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	return GetFormatter(format)(w, data)
}

// TrimBasePathFromPath trims the base path prefix from the path
func TrimBasePathFromPath(basePath string, path string) string {
	return strings.TrimPrefix(path, basePath)
}

// JoinAbsolutePathWithPath checks if the provided path is absolute. If the provided path is relative, it joins the base path with the path and returns the absolute path
func JoinAbsolutePathWithPath(basePath string, providedPath string) (string, error) {
	// If the provided path is an absolute path, return it
	if filepath.IsAbs(providedPath) {
		return providedPath, nil
	}

	// Join the base path with the provided path
	joinedPath := path.Join(basePath, providedPath)

	// If the joined path is an absolute path, return it
	if filepath.IsAbs(joinedPath) {
		return joinedPath, nil
	}

	// Convert the joined path to an absolute path
	absPath, err := filepath.Abs(joinedPath)
	if err != nil {
		return "", err
	}

	// Check if the final absolute path exists in the file system
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		return "", err
	}

	return absPath, nil
}
