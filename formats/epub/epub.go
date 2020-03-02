package epub

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
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

// EpubImageReader epub image reader
type EpubImageReader struct {
	reader *zip.ReadCloser
}

func NewImageReader(path string) (*EpubImageReader, error) {
	reader, err := zip.OpenReader(path)

	if err != nil {
		return nil, err
	}

	return &EpubImageReader{reader: reader}, nil
}

func (self *EpubImageReader) Extract(name string) ([]byte, error) {
	return make([]byte, 0), nil
}

func (self *EpubImageReader) List() (map[string]string, error) {
	images := make(map[string]string)

	container, err := self.file("META-INF/container.xml")

	if err != nil {
		return images, err
	}

	_, err = self.rootFiles(container)

	if err != nil {
		return images, err
	}

	return images, nil
}

func (self *EpubImageReader) Close() error {
	return self.reader.Close()
}

var ErrorNoSuchFile = errors.New("no such file")

func (self *EpubImageReader) file(path string) (*zip.File, error) {
	for _, f := range self.reader.File {
		if f.Name == path {
			return f, nil
		}
	}

	return nil, ErrorNoSuchFile
}

func (self *EpubImageReader) rootFiles(container *zip.File) ([]*zip.File, error) {
	rootFiles := make([]*zip.File, 0)

	_, err := self.data(container)

	if err != nil {
		return rootFiles, err
	}

	return rootFiles, nil
}

func (self *EpubImageReader) data(file *zip.File) ([]byte, error) {
	var buf bytes.Buffer

	writer := bufio.NewWriter(&buf)

	reader, err := file.Open()

	if err != nil {
		return buf.Bytes(), err
	}

	defer reader.Close()

	_, err = io.Copy(writer, reader)

	return buf.Bytes(), err
}

func (self *EpubImageReader) rootFilePaths(containerData []byte) ([]string, error) {
	paths := make([]string, 0)
	return paths, nil
}

/*
// Extract extracts images from epub file filename into folder dir and
// returns number of exported images and error if exists
func (extractor *EpubImageReader) Extract(filename string, dir string, coveronly bool) (int, error) {
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

	return 0, nil
}


func (extractor *EpubImageReader) extractEpub(filename string, dest string) error {

	err := futil.ZipExtract(filename, dest)
	return err

}

func (extractor *EpubImageReader) getRootFiles(tempdir string) ([]string, error) {

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

func (extractor *EpubImageReader) readContainer(containerPath string) ([]string, error) {

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

func (extractor *EpubImageReader) readManifest(filename string, items *[]string, coveronly bool) error {

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
*/
