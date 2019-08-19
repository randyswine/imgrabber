package helper

import (
	"fmt"
	"os"
	"strings"
)

// MakeDestinationDir creates a directory to load.
func MakeDestinationDir(dirpath string) error {
	if IsPathExists(dirpath) == false {
		err := os.Mkdir(dirpath, os.ModeDir)
		if err != nil {
			return fmt.Errorf("error of download images (make dirpath): %v", err)
		}
	}
	return nil
}

/*
	IsPathExists returns a boolean indicating whether the error is
	known to report that a file or directory does not exist.
*/
func IsPathExists(path string) bool {
	result := true
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		result = false
	}
	return result
}

func IsCorrectExtension(url string, extensions ...string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(url, ext) {
			return true
		}
	}
	return false
}
