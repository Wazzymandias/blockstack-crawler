package cmd

import "github.com/spf13/cobra"

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Fetch Blockstack app related information",
}

func init() {
	appCmd.AddCommand(newUsersCmd)
}
