package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "bfc",
	Short: "bfc scrapes Blockstack API and indexes app and profile information",
}

func init() {
	rootCmd.AddCommand(appCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
