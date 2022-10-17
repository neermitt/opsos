package utils

import (
	"io"
	"os"
)

func FileExists(filename string) bool {
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) || err != nil {
		return false
	}
	return !fileInfo.IsDir()
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

	return Get(format)(w, data)
}
