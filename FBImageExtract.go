package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	ebook "./formats"
	epub "./formats/epub"
	fb2 "./formats/fictionbook/fictionbook2"
)

var input = flag.String("i", "", "input file [required]")
var output = flag.String("o", "", "output folder [optional]")
var cover = flag.Bool("c", false, "extract cover only [optional]")
var logfilename = flag.String("l", "", "log filename [optional]")

func logFatal(logger *log.Logger, err error) {

	if logger != nil {
		fmt.Printf("%v\n", err)
		logger.Fatal(err)
	} else {
		log.Fatal(err)
	}

}

func getImageExtractor(filename string) (ebook.IImageExtractor, error) {

	var fext string

	fext = filepath.Ext(filename)

	switch fext {
	case ".fb2":
		imageExtractor := new(fb2.ImageExtractor)
		return imageExtractor, nil
	case ".epub":
		imageExtractor := new(epub.ImageExtractor)
		return imageExtractor, nil
	}

	return nil, errors.New("unsupported file format")
}

func main() {

	var (
		logger    *log.Logger
		extractor ebook.IImageExtractor
		count     int
		duration  time.Duration
		err       error
	)

	start := time.Now()

	flag.Parse()

	if strings.Trim(*logfilename, " ") != "" {

		f, ferr := os.Create(*logfilename)
		defer f.Close()

		if ferr == nil {
			logger = log.New(f, "fbimageextract: ", log.Lshortfile|log.Ldate|log.Ltime)
		} else {
			log.Fatal(ferr)
		}

	}

	if _, err = os.Stat(*input); os.IsNotExist(err) {
		logFatal(logger, errors.New("input file: no such file"))
	}

	if *output == "" {

		if *output, err = filepath.Abs(*input); err != nil {
			logFatal(logger, err)
		}

		*output = filepath.Dir(*output)

	}

	if extractor, err = getImageExtractor(*input); err != nil {
		logFatal(logger, err)
	}

	count, err = extractor.Extract(*input, *output, *cover)

	end := time.Now()
	duration = end.Sub(start)

	fmt.Printf("%d image(s) extracted in %.3fs\n", count, duration.Seconds())

	if err != nil {
		logFatal(logger, err)
	}
}
