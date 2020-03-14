package epub

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var cases = []struct {
	path  string
	cover string
}{
	{
		path:  "testdata/test.epub",
		cover: "cover.jpg",
	},
}

func listImages(path string) (map[string]string, error) {
	reader, err := NewImageReader(path)

	if err != nil {
		return make(map[string]string), err
	}

	defer reader.Close()

	return reader.List()
}

func extractCover(path, cover string) ([]byte, error) {
	reader, err := NewImageReader(path)

	if err != nil {
		return make([]byte, 0), err
	}

	defer reader.Close()

	return reader.Extract(cover)
}

func TestList(t *testing.T) {
	_, testFilename, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	for _, test := range cases {
		path := filepath.Join(filepath.Dir(testFilename), test.path)

		images, err := listImages(path)

		if err != nil {
			t.Fatal(err)
		}

		if len(images) > 0 {
			t.Fatal(fmt.Errorf("listing images error in %s", test.path))
		}
	}
}

func TestExtract(t *testing.T) {
	_, testFilename, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	for _, test := range cases {
		path := filepath.Join(filepath.Dir(testFilename), test.path)

		if strings.Trim(test.cover, " ") != "" {
			b, err := extractCover(path, test.cover)

			if err != nil {
				t.Fatal(err)
			}

			if len(b) == 0 {
				t.Fatal(fmt.Errorf("reading cover error in %s", test.path))
			}
		}
	}
}
