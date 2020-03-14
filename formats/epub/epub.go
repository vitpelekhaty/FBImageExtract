package epub

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

type RootFile struct {
	FullPath  string `xml:"full-path,attr"`
	MediaType string `xml:"media-type,attr"`
}

type RootFiles struct {
	RootFiles []RootFile `xml:"rootfile"`
}

type MetaContainer struct {
	Roots []RootFiles `xml:"rootfiles"`
}

type ManifestItem struct {
	XMLName   xml.Name `xml:"item"`
	ID        string   `xml:"id,attr"`
	Href      string   `xml:"href,attr"`
	MediaType string   `xml:"media-type,attr"`
}

type Manifest struct {
	XMLName xml.Name       `xml:"manifest"`
	Items   []ManifestItem `xml:"item"`
}

type Meta struct {
	XMLName xml.Name `xml:"meta"`
	Name    string   `xml:"name,attr"`
	Content string   `xml:"content,attr"`
}

type Metadata struct {
	XMLName xml.Name `xml:"metadata"`
	Meta    Meta     `xml:"meta"`
}

type Package struct {
	XMLName  xml.Name `xml:"package"`
	Metadata Metadata `xml:"metadata"`
	Manifest Manifest `xml:"manifest"`
}

var CoreImageTypes = []string{
	"image/gif",
	"image/jpeg",
	"image/png",
	"image/svg+xml"}

// File struct reader
type File struct {
	reader *zip.ReadCloser
}

// Open
func Open(path string) (*File, error) {
	reader, err := zip.OpenReader(path)

	if err != nil {
		return nil, err
	}

	return &File{reader: reader}, nil
}

// Close
func (self *File) Close() error {
	return self.reader.Close()
}

// Read
func (self *File) Read(section string) ([]byte, error) {
	s, err := self.section(section)

	if err != nil {
		return make([]byte, 0), err
	}

	var buf bytes.Buffer

	writer := bufio.NewWriter(&buf)

	reader, err := s.Open()

	if err != nil {
		return buf.Bytes(), err
	}

	defer reader.Close()

	_, err = io.Copy(writer, reader)

	return buf.Bytes(), err
}

// RootFiles
func (self *File) RootFiles() ([]RootFile, error) {
	d, err := self.Read("META-INF/container.xml")

	if err != nil {
		return make([]RootFile, 0), err
	}

	var container MetaContainer

	err = xml.Unmarshal(d, &container)

	if err != nil {
		return make([]RootFile, 0), err
	}

	roots := make([]RootFile, 0)

	for _, r := range container.Roots {
		roots = append(roots, r.RootFiles...)
	}

	return roots, nil
}

// Package
func (self *File) Package(rootPath string) (*Package, error) {
	d, err := self.Read(rootPath)

	if err != nil {
		return nil, err
	}

	var pack Package

	err = xml.Unmarshal(d, &pack)

	if err != nil {
		return nil, err
	}

	return &pack, nil
}

// ErrorNoCover
var ErrorNoCover = errors.New("no cover image")

// CoverImage
func (self *File) CoverImage(rootPath string) (*ManifestItem, error) {
	pack, err := self.Package(rootPath)

	if err != nil {
		return nil, err
	}

	meta := pack.Metadata.Meta

	if !strings.EqualFold(meta.Name, "cover") {
		return nil, ErrorNoCover
	}

	for _, item := range pack.Manifest.Items {
		if item.ID == meta.Content {
			return &item, nil
		}
	}

	return nil, nil
}

// section
func (self *File) section(path string) (*zip.File, error) {
	for _, f := range self.reader.File {
		if f.Name == path {
			return f, nil
		}
	}

	return nil, fmt.Errorf("no section %s", path)
}
