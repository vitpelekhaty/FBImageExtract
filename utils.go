package main

import (
	"os"
)

// IsDir
func IsDir(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		return fi.IsDir()
	}

	return false
}

// SaveToFile
func SaveToFile(b []byte, path string) error {
	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(b)

	return err
}
