package epub

import (
	"path/filepath"
	"strings"
)

// EpubImageReader
type EpubImageReader struct {
	file *File
}

// NewReader
func NewImageReader(path string) (*EpubImageReader, error) {
	file, err := Open(path)

	if err != nil {
		return nil, err
	}

	return &EpubImageReader{file: file}, nil
}

// Extract
func (self *EpubImageReader) Extract(name string) ([]byte, error) {
	return self.file.Read(name)
}

// List
func (self *EpubImageReader) List() (map[string]string, error) {
	images := make(map[string]string)

	roots, err := self.file.RootFiles()

	if err != nil {
		return images, err
	}

	filter := func(path, mediaType string) bool {
		var res = false

		for _, coreImageType := range CoreImageTypes {
			res = res || strings.EqualFold(mediaType, coreImageType)

			if res {
				return res
			}
		}

		return res
	}

	var dir, href string

	for _, root := range roots {
		pack, err := self.file.Package(root.FullPath)

		if err != nil {
			return images, err
		}

		dir = filepath.Dir(root.FullPath)

		for _, item := range pack.Manifest.Items {
			href = filepath.Join(dir, item.Href)

			if filter(href, item.MediaType) {
				images[href] = item.MediaType
			}
		}
	}

	return images, nil
}

// Close
func (self *EpubImageReader) Close() error {
	return self.file.Close()
}
