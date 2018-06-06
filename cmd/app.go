package cmd

import "github.com/spf13/cobra"

var appCmd = &cobra.Command{}

func init() {
	appCmd.AddCommand(appUsersCmd)
}
