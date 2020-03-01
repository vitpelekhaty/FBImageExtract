package ebook

import (
	"errors"
	"path/filepath"

	epub "github.com/vitpelekhaty/fbimgextract/formats/epub"
	fb2 "github.com/vitpelekhaty/fbimgextract/formats/fictionbook/fictionbook2"
)

// ImageList list of images
type ImageList map[string]string

// IEbookImageReader image extractor interface
type IEbookImageReader interface {
	Extract(name string) ([]byte, error)
	List() (ImageList, error)
	Close() error
}

// NewEBookImageReader returns FictionBookImageReader object
func NewEBookImageReader(path string) (IEbookImageReader, error) {
	ext := filepath.Ext(path)

	switch ext {
	case ".fb2":
		return fb2.NewImageReader(path)
	case ".epub":
		return epub.NewImageReader(path)
	default:
		return nil, errors.New("unsupported format")
	}
}

// ExtractResult result of image extraction
type ExtractResult struct {
	Ok    bool
	Error error
}
