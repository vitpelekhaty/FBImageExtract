package main

import (
	"os"

	ebook "github.com/vitpelekhaty/fbimgextract/formats"
)

// Extract
func Extract(path, imageName, output string) error {
	reader, err := ebook.NewEBookImageReader(path)

	if err != nil {
		return err
	}

	b, err := reader.Extract(imageName)

	if err != nil {
		return err
	}

	return SaveToFile(b, output)
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
