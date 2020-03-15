package main

import (
	"os"
	"path/filepath"

	ebook "github.com/vitpelekhaty/fbimgextract/formats"
)

// Extract
func Extract(path, imageName, output string) error {
	reader, err := ebook.NewEBookImageReader(path)

	if err != nil {
		return err
	}

	defer reader.Close()

	b, err := reader.Extract(imageName)

	if err != nil {
		return err
	}

	return SaveToFile(b, output)
}

// ExtractAll
func ExtractAll(path, outputDir string) (int, error) {
	reader, err := ebook.NewEBookImageReader(path)

	if err != nil {
		return 0, err
	}

	defer reader.Close()

	images, err := reader.List()

	if err != nil {
		return 0, err
	}

	var n = 0
	var filename string

	for imageID := range images {
		b, err := reader.Extract(imageID)

		if err != nil {
			return n, err
		}

		filename = filepath.Join(outputDir, imageID)

		err = os.MkdirAll(filename, 0)

		if err != nil {
			return n, err
		}

		err = SaveToFile(b, filename)

		if err != nil {
			return n, err
		}

		n++
	}

	return n, nil
}
