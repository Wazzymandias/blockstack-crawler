package cmd

import (
	"fmt"
	"github.com/Wazzymandias/blockstack-profile-crawler/app"
	"github.com/spf13/cobra"
)

var newUsersCmd = &cobra.Command{
	Use:   "new-users [OPTIONS]",
	Short: "display information related to new users for Blockstack apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO validate required flags and make sure values are correct

		rh, err := app.NewRequestHandler()

		if err != nil {
			return err
		}

		//t, err := time.Parse(time.RFC3339, "")
		//
		//if err != nil {
		//	return err
		//}
		//newUsers, err := rh.RetrieveNewNames(t)

		newUsers, err := rh.RetrieveNames()

		if err != nil {
			return err
		}

		printUsers(newUsers)

		return rh.Shutdown()
	},
}

func printUsers(u map[string]map[string]bool) {
	fmt.Println(u)
}
