package utils

import "os"

// FileExists returns true if a given file exists, and false if it doesn't.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
