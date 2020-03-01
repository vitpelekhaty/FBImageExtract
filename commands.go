package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cmdRoot
var cmdRoot = &cobra.Command{
	Short: "fbimgextract",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.UsageString()
	},
}

// cmdVersion show version
var cmdVersion = &cobra.Command{
	Short: "version",
	Long:  "show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version")
	},
}

// cmdExtract extract images
var cmdExtract = &cobra.Command{
	Short: "extract",
	Long:  "extract images",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

// cmdList list images
var cmdList = &cobra.Command{
	Short: "list",
	Long:  "list images",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

// Execute run app
func Execute() error {
	return cmdRoot.Execute()
}
