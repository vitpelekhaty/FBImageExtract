package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// HasExtension checks whether a filename has an extension
func HasExtension(filename string) bool {

	var ext string

	ext = path.Ext(filename)
	return strings.Trim(ext, " ") != ""

}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {

	srcinfo, err := os.Stat(src)

	if err != nil {
		return err
	}

	if !srcinfo.Mode().IsRegular() {
		return fmt.Errorf("CopyFile: non regular source file %s (%q)", srcinfo.Name(), srcinfo.Mode().String())
	}

	dstinfo, err := os.Stat(dst)

	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !dstinfo.Mode().IsRegular() {
			return fmt.Errorf("CopyFile: non regular destination file %s (%q)", dstinfo.Name(), dstinfo.Mode().String())
		}

		if os.SameFile(srcinfo, dstinfo) {
			return errors.New("CopyFile: cannot copy a file itself")
		}
	}

	err = copyFileContents(src, dst)
	return err
}

func copyFileContents(src, dst string) error {

	in, err := os.Open(src)

	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer func() {
		copyerr := out.Close()
		if err == nil {
			err = copyerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	err = out.Sync()

	return err
}
