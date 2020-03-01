package fictionbook2

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
)

// ImageExtensions supported image formats
var ImageExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

// binarySection binary section
type binarySection struct {
	ID          string `xml:"id,attr"`
	ContentType string `xml:"content-type,attr"`
	Raw         string `xml:",innerxml"`
}

// coverProperties cover image properties
type coverProperties struct {
	XMLName xml.Name `xml:"image"`
	Href    string   `xml:"href,attr"`
}

// cover image section
type cover struct {
	XMLName xml.Name        `xml:"coverpage"`
	Image   coverProperties `xml:"image"`
}

// titleInfo title info section
type titleInfo struct {
	XMLName   xml.Name `xml:"title-info"`
	CoverPage cover    `xml:"coverpage"`
}

// description section
type description struct {
	XMLName   xml.Name  `xml:"description"`
	TitleInfo titleInfo `xml:"title-info"`
}

// book file structure
type book struct {
	XMLName     xml.Name        `xml:"FictionBook"`
	Description description     `xml:"description"`
	Binaries    []binarySection `xml:"binary"`
}

func openBook(path string) (*book, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	data, err := ioutil.ReadAll(f)

	if err != nil {
		return nil, err
	}

	var book book

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = identCharsetReader

	err = decoder.Decode(&book)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

func identCharsetReader(encoding string, input io.Reader) (io.Reader, error) {
	return input, nil
}
