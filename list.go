package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	ebook "github.com/vitpelekhaty/fbimgextract/formats"
)

// List
func List(path string) error {
	reader, err := ebook.NewEBookImageReader(path)

	if err != nil {
		return err
	}

	images, err := reader.List()

	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 20, 0, '\t', 0)

	for ID, mediaType := range images {
		fmt.Fprintf(tw, "%s\t%s\t\n", ID, mediaType)
	}

	tw.Flush()

	return nil
}
