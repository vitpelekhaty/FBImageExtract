package epub

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestOpenClose(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	path := filepath.Join(filepath.Dir(file), "testdata/test.epub")

	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}

	book, err := Open(path)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = book.Close(); err != nil {
			t.Log(err)
		}
	}()
}

func TestRead(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	path := filepath.Join(filepath.Dir(file), "testdata/test.epub")

	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}

	book, err := Open(path)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = book.Close(); err != nil {
			t.Log(err)
		}
	}()

	b, err := book.Read("META-INF/container.xml")

	if err != nil {
		t.Fatal(err)
	}

	if len(b) == 0 {
		t.FailNow()
	}
}

func TestRootFiles(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	path := filepath.Join(filepath.Dir(file), "testdata/test.epub")

	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}

	book, err := Open(path)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = book.Close(); err != nil {
			t.Log(err)
		}
	}()

	roots, err := book.RootFiles()

	if err != nil {
		t.Fatal(err)
	}

	if len(roots) == 0 {
		t.FailNow()
	}
}

func TestPackage(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	path := filepath.Join(filepath.Dir(file), "testdata/test.epub")

	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}

	book, err := Open(path)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = book.Close(); err != nil {
			t.Log(err)
		}
	}()

	roots, err := book.RootFiles()

	if err != nil {
		t.Fatal(err)
	}

	if len(roots) == 0 {
		t.FailNow()
	}

	pack, err := book.Package(roots[0].FullPath)

	if err != nil {
		t.Fatal(err)
	}

	if pack == nil {
		t.FailNow()
	}
}

func TestCoverImage(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	path := filepath.Join(filepath.Dir(file), "testdata/test.epub")

	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}

	book, err := Open(path)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = book.Close(); err != nil {
			t.Log(err)
		}
	}()

	roots, err := book.RootFiles()

	if err != nil {
		t.Fatal(err)
	}

	if len(roots) == 0 {
		t.FailNow()
	}

	cover, err := book.CoverImage(roots[0].FullPath)

	if err != nil {
		t.Fatal(err)
	}

	if cover == nil {
		t.FailNow()
	}
}
