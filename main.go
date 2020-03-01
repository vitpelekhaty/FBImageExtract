package main

import (
	"log"
)

var (
	// Input file
	Input string
	// OutputDir output directory
	OutputDir string
	// CoverOnly extract cover only
	CoverOnly bool
)

func main() {
	if err := Execute(); err != nil {
		log.Fatal(err)
	}
}
