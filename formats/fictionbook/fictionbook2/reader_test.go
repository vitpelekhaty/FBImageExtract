package fictionbook2

import (
	"path/filepath"
	"runtime"
	"testing"
)

var cases = []struct {
	path   string
	images map[string]string
}{
	{
		path: "testdata/example1.fb2",
		images: map[string]string{
			"fish1.png": "image/png",
			"fish.jpg":  "image/jpeg",
			"free.png":  "image/png",
			"subs.png":  "image/png",
			"paid.png":  "image/png",
		},
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

func TestList(t *testing.T) {
	var done bool

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

		done = len(images) == len(test.images)

		if !done {
			t.Fail()
		}

		for tk, tv := range test.images {
			v, ok := images[tk]

			done = (done && ok) && (v == tv)
		}

		if !done {
			t.Fail()
		}
	}
}

func extractImage(path, imageID string) ([]byte, error) {
	reader, err := NewImageReader(path)

	if err != nil {
		return make([]byte, 0), err
	}

	defer reader.Close()

	return reader.Extract(imageID)
}

func TestExtract(t *testing.T) {
	_, testFilename, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	for _, test := range cases {
		path := filepath.Join(filepath.Dir(testFilename), test.path)

		for imageID := range test.images {
			data, err := extractImage(path, imageID)

			if err != nil {
				t.Fatal(err)
			}

			if len(data) == 0 {
				t.Fail()
			}
		}
	}
}
