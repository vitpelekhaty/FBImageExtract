package epub

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"sync"

	ebook ".."
	futil "../utils"
)

type rootFile struct {
	FullPath  string `xml:"full-path,attr"`
	MediaType string `xml:"media-type,attr"`
}

type rootFiles struct {
	RootFiles []rootFile `xml:"rootfile"`
}

type metaContainer struct {
	Roots []rootFiles `xml:"rootfiles"`
}

type manifestitem struct {
	XMLName   xml.Name `xml:"item"`
	ID        string   `xml:"id,attr"`
	Href      string   `xml:"href,attr"`
	MediaType string   `xml:"media-type,attr"`
}

type manifest struct {
	XMLName xml.Name       `xml:"manifest"`
	Items   []manifestitem `xml:"item"`
}

type meta struct {
	XMLName xml.Name `xml:"meta"`
	Name    string   `xml:"name,attr"`
	Content string   `xml:"content,attr"`
}

type metadata struct {
	XMLName xml.Name `xml:"metadata"`
	Meta    meta     `xml:"meta"`
}

type pckg struct {
	XMLName  xml.Name `xml:"package"`
	Metadata metadata `xml:"metadata"`
	Manifest manifest `xml:"manifest"`
}

var coreImageTypes = []string{
	"image/gif",
	"image/jpeg",
	"image/png",
	"image/svg+xml"}

// ImageExtractor - image extractor from epub books
type ImageExtractor struct {
	ebook.IImageExtractor
}

// Extract extracts images from epub file filename into folder dir and
// returns number of exported images and error if exists
func (extractor *ImageExtractor) Extract(filename string, dir string, coveronly bool) (int, error) {

	var (
		count       int
		roots       []string
		contentPath string
		items       []string
		wg          sync.WaitGroup
	)

	reschan := make(chan ebook.ExtractResult)

	tempdir, err := ioutil.TempDir("", "fbie")

	if err != nil {
		return count, err
	}

	defer os.RemoveAll(tempdir)

	if err = extractor.extractEpub(filename, tempdir); err != nil {
		return count, err
	}

	if roots, err = extractor.getRootFiles(tempdir); err != nil {
		return count, err
	}

	for _, root := range roots {

		contentPath = path.Join(tempdir, root)

		err := extractor.readManifest(contentPath, &items, coveronly)

		if err != nil {
			return count, err
		}

		for _, imagepath := range items {

			src := path.Join(tempdir, imagepath)
			dest := path.Join(dir, imagepath)

			wg.Add(1)

			go func(s, d string) {

				defer wg.Done()

				err := futil.CopyFile(s, d)

				reschan <- ebook.ExtractResult{Ok: err == nil, Error: err}

			}(src, dest)
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

func (extractor *ImageExtractor) extractEpub(filename string, dest string) error {

	err := futil.ZipExtract(filename, dest)
	return err

}

func (extractor *ImageExtractor) getRootFiles(tempdir string) ([]string, error) {

	var (
		roots         []string
		containerPath string
		err           error
	)

	containerPath = path.Join(tempdir, "META-INF/container.xml")

	if _, err = os.Stat(containerPath); os.IsNotExist(err) {
		return roots, errors.New("invalid epub book format")
	}

	roots, err = extractor.readContainer(containerPath)

	if err != nil {
		return roots, err
	}

	return roots, nil

}

func (extractor *ImageExtractor) readContainer(containerPath string) ([]string, error) {

	var (
		roots     []string
		container metaContainer
	)

	xmlFile, err := os.Open(containerPath)
	defer xmlFile.Close()

	if err != nil {
		return roots, err
	}

	xmlData, err := ioutil.ReadAll(xmlFile)

	if err != nil {
		return roots, err
	}

	if err := xml.Unmarshal(xmlData, &container); err != nil {
		return roots, err
	}

	for _, r := range container.Roots {
		for _, rf := range r.RootFiles {
			roots = append(roots, rf.FullPath)
		}
	}

	return roots, nil

}

func (extractor *ImageExtractor) readManifest(filename string, items *[]string, coveronly bool) error {

	var (
		pack        pckg
		metaContent string
		metaName    string
		canAppend   bool
	)

	xmlFile, err := os.Open(filename)
	defer xmlFile.Close()

	if err != nil {
		return err
	}

	xmlData, err := ioutil.ReadAll(xmlFile)

	if err != nil {
		return err
	}

	if err := xml.Unmarshal(xmlData, &pack); err != nil {
		return err
	}

	metaName = pack.Metadata.Meta.Name
	metaContent = pack.Metadata.Meta.Content

	for _, item := range pack.Manifest.Items {

		canAppend = true

		if coveronly {
			canAppend = metaName == "cover" && metaContent == item.ID
		}

		if ok := listContainsString(item.Href, *items); !ok && canAppend {
			if ok := listContainsString(item.MediaType, coreImageTypes); ok {
				*items = append(*items, item.Href)
			}
		}

	}

	return nil
}
