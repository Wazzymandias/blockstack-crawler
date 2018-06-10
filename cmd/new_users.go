package cmd

import (
	"github.com/Wazzymandias/blockstack-profile-crawler/app"
	"github.com/spf13/cobra"
	"github.com/Wazzymandias/blockstack-profile-crawler/config"
	"time"
	"fmt"
)

var usersCmd = &cobra.Command{
	Use:   "new-users [OPTIONS]",
	Short: "display information related to new users for Blockstack apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO validate required flags and make sure values are correct

		var users map[string]map[string]bool
		var err error

		rh, err := app.NewRequestHandler()

		if err != nil {
			return err
		}

		if config.NewUsersSince != "" {
			var t time.Time
			layout := "2006-01-02"

			t, err = time.Parse(layout, config.NewUsersSince)

			if err != nil {
				return err
			}

			users, err = rh.RetrieveNewNames(t)
		} else {
			users, err = rh.RetrieveNames()
		}

		if err != nil {
			return err
		}

		printUsers(users)

		return rh.Shutdown()
	},
}

func init() {
	usersCmd.Flags().StringVarP(&config.NewUsersSince, "since", "s", "", "ISO 8601 formatted date [YYYY-MM-DD]")
}

func printUsers(u map[string]map[string]bool) {
	fmt.Println(u)
}

