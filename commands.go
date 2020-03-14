package main

import (
	"fmt"

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

// cmdExtract extract images
var cmdExtract = &cobra.Command{
	Use:   "extract",
	Short: "extract image(-s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Extract(Input, ImageName, Output)
	},
}

// cmdList list images
var cmdList = &cobra.Command{
	Use:   "list",
	Short: "list images",
	RunE: func(cmd *cobra.Command, args []string) error {
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

	cmdRoot.AddCommand(cmdVersion, cmdExtract, cmdList)
}
