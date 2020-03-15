package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// cmdRoot
var cmdRoot = &cobra.Command{
	Use:  "fbimgextract",
	Long: "Tool for listing images in ebooks and extracting them from ebooks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.UsageString())
	},
}

// cmdVersion show version
var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version")
	},
}

// ErrorUndefinedImageID
var ErrorUndefinedImageID = errors.New("undefined image ID")

// ErrorOutputIsNotDir
var ErrorOutputIsNotDir = errors.New("output is not a dir")

// cmdExtract extract images
var cmdExtract = &cobra.Command{
	Use:   "extract",
	Short: "extract image(-s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := os.Stat(Input)

		if err != nil {
			return err
		}

		if All {
			if !IsDir(Output) {
				return ErrorOutputIsNotDir
			}

			_, err = ExtractAll(Input, Output)

			return err
		}

		if strings.Trim(ImageName, " ") == "" {
			return ErrorUndefinedImageID
		}

		var (
			filename string
			bookPath string
		)

		if strings.Trim(Output, " ") == "" {
			bookPath, _ = filepath.Split(Input)
			_, filename = filepath.Split(ImageName)

			filename = filepath.Join(bookPath, filename)
		} else {
			if IsDir(Output) {
				_, filename = filepath.Split(ImageName)
				filename = filepath.Join(Output, filename)
			} else {
				filename = Output
			}
		}

		return Extract(Input, ImageName, filename)
	},
}

// cmdList list images
var cmdList = &cobra.Command{
	Use:   "list",
	Short: "list images",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := os.Stat(Input)

		if err != nil {
			return err
		}

		return List(Input)
	},
}

// Execute run app
func Execute() error {
	return cmdRoot.Execute()
}

func init() {
	cmdList.Flags().StringVarP(&Input, "input", "i", "", "ebook filename")

	cmdExtract.Flags().StringVarP(&Input, "input", "i", "", "ebook filename")
	cmdExtract.Flags().StringVarP(&Output, "output", "o", "", "output folder")
	cmdExtract.Flags().StringVarP(&ImageName, "image", "I", "", "image to extract")
	cmdExtract.Flags().BoolVarP(&All, "all", "a", false, "extract all images")

	cmdRoot.AddCommand(cmdVersion, cmdExtract, cmdList)
}
