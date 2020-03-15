package ebook

import (
	"errors"
	"path/filepath"

	epub "github.com/vitpelekhaty/fbimgextract/formats/epub"
	fb2 "github.com/vitpelekhaty/fbimgextract/formats/fictionbook/fictionbook2"
)

// IEbookImageReader image extractor interface
type IEbookImageReader interface {
	Extract(name string) ([]byte, error)
	List() (map[string]string, error)
	Close() error
}

// ErrorUnsupportedFormat
var ErrorUnsupportedFormat = errors.New("unsupported format")

// NewEBookImageReader returns FictionBookImageReader object
func NewEBookImageReader(path string) (IEbookImageReader, error) {
	ext := filepath.Ext(path)

	switch ext {
	case ".fb2":
		return fb2.NewImageReader(path)
	case ".epub":
		return epub.NewImageReader(path)
	default:
		return nil, ErrorUnsupportedFormat
	}
}
