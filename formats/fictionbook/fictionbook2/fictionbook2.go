package fictionbook2

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"

	ebook "../.."
	futil "../../utils"
)

var imageExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png"}

type fictionBookBinary struct {
	ID          string `xml:"id,attr"`
	ContentType string `xml:"content-type,attr"`
	Raw         string `xml:",innerxml"`
}

type fictionBookCoverImageProps struct {
	XMLName xml.Name `xml:"image"`
	Href    string   `xml:"href,attr"`
}

type fictionBookCoverPage struct {
	XMLName xml.Name                   `xml:"coverpage"`
	Image   fictionBookCoverImageProps `xml:"image"`
}

type fictionBookTitleInfo struct {
	XMLName   xml.Name             `xml:"title-info"`
	CoverPage fictionBookCoverPage `xml:"coverpage"`
}

type fictionBookDescription struct {
	XMLName   xml.Name             `xml:"description"`
	TitleInfo fictionBookTitleInfo `xml:"title-info"`
}

type fictionBook struct {
	XMLName     xml.Name               `xml:"FictionBook"`
	Description fictionBookDescription `xml:"description"`
	Binaries    []fictionBookBinary    `xml:"binary"`
}

// ImageExtractor - image extractor from fictionbooks
type ImageExtractor struct {
	ebook.IImageExtractor
}

// Extract extracts images from fictionbook file filename into folder dir and
// returns number of exported images and error if exists
func (extractor *ImageExtractor) Extract(filename string, dir string, coveronly bool) (int, error) {

	var (
		count      int
		coverimage string
		cansave    bool
		fb         fictionBook
		cimage     fictionBookCoverImageProps
		wg         sync.WaitGroup
	)

	reschan := make(chan ebook.ExtractResult)

	xmlFile, err := os.Open(filename)
	defer xmlFile.Close()

	if err != nil {
		return count, err
	}

	xmlData, err := ioutil.ReadAll(xmlFile)

	if err != nil {
		return count, err
	}

	if err := xml.Unmarshal(xmlData, &fb); err != nil {
		return count, err
	}

	cimage = fb.Description.TitleInfo.CoverPage.Image
	coverimage = strings.Trim(cimage.Href, "#")

	for _, fbBinary := range fb.Binaries {

		cansave = true

		if coveronly {
			cansave = coverimage == fbBinary.ID
		}

		if cansave {

			wg.Add(1)

			go func(bin fictionBookBinary, dst string) {

				defer wg.Done()

				err := extractor.saveImage(bin, dst)
				reschan <- ebook.ExtractResult{Ok: err == nil, Error: err}

			}(fbBinary, dir)
		}
	}

	go func() {

		wg.Wait()
		close(reschan)

	}()

	for result := range reschan {
		if result.Ok {
			count++
		} else {
			return count, result.Error
		}
	}

	return count, nil
}

func (extractor *ImageExtractor) saveImage(binary fictionBookBinary, dir string) error {

	var (
		filename    string
		contentType string
	)

	filename = path.Join(dir, binary.ID)
	contentType = strings.ToLower(binary.ContentType)

	if !futil.HasExtension(filename) {
		filename = strings.Join([]string{filename, imageExtensions[contentType]}, "")
	}

	imageData, err := DecodeBase64(binary.Raw)

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, imageData, 0100664); err != nil {
		return err
	}

	return nil
}
