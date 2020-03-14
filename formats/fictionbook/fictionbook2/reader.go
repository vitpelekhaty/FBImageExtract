package fictionbook2

import (
	"encoding/base64"
)

// FictionBookImageReader image extractor
type FictionBookImageReader struct {
	book *book
}

// NewImageReader
func NewImageReader(path string) (*FictionBookImageReader, error) {
	book, err := openBook(path)

	if err != nil {
		return nil, err
	}

	return &FictionBookImageReader{book: book}, nil
}

// Extract
func (self *FictionBookImageReader) Extract(name string) ([]byte, error) {
	for _, image := range self.book.Binaries {
		if image.ID == name {
			return base64.StdEncoding.DecodeString(image.Raw)
		}
	}

	return make([]byte, 0), nil
}

// List
func (self *FictionBookImageReader) List() (map[string]string, error) {
	images := make(map[string]string)

	for _, image := range self.book.Binaries {
		images[image.ID] = image.ContentType
	}

	return images, nil
}

// Close
func (self *FictionBookImageReader) Close() error {
	return nil
}
