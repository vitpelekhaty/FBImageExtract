package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ZipExtract extracts content from zip archive file filename into directory dest
func ZipExtract(filename string, dest string) error {

	zr, err := zip.OpenReader(filename)

	if err != nil {
		return err
	}

	defer zr.Close()

	for _, file := range zr.File {

		rc, err := file.Open()

		if err != nil {
			return err
		}

		defer rc.Close()

		fpath := filepath.Join(dest, file.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {

			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			out, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())

			if err != nil {
				return err
			}

			_, err = io.Copy(out, rc)

			out.Close()

			if err != nil {
				return err
			}
		}

	}

	return nil
}
