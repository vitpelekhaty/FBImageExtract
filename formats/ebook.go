package ebook

// IImageExtractor - image extractor interface
type IImageExtractor interface {
	Extract(filename string, dir string, coveronly bool) (int, error)
}

// ExtractResult - result of image extraction
type ExtractResult struct {
	Ok    bool
	Error error
}
