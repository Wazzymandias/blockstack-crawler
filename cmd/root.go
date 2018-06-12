// Package cmd implements CLI commands that perform Blockstack related operations
package cmd

import (
	"fmt"
	"github.com/Wazzymandias/blockstack-crawler/config"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   config.ProgramName,
	Short: fmt.Sprintf("%s scrapes Blockstack API and outputs in text or json format", config.ProgramName),
}

func init() {
	rootCmd.PersistentFlags().StringVar(&config.DatabaseType, "db", config.DefaultDBType,
		"type of database to store results in")

	rootCmd.PersistentFlags().StringVar(&config.StorageType, "store", config.DefaultStorageType,
		"type of storage to use for persisting file data")

	rootCmd.AddCommand(namesCmd)
}

// Execute starts CLI tool. If exits with
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
