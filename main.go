package main

import (
	"log"
)

var (
	// Input file
	Input string
	// Output output filename or directory
	Output string
	// ImageName image to extract
	ImageName string
	// All
	All bool
)

func main() {
	if err := Execute(); err != nil {
		log.Fatal(err)
	}
}
