package cmd

import (
	"fmt"
	"github.com/Wazzymandias/blockstack-crawler/config"
	"github.com/Wazzymandias/blockstack-crawler/worker"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var namesCmd = &cobra.Command{
	Use:   "names [OPTIONS]",
	Short: "display information related to users for Blockstack apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		var names map[string]map[string]bool
		var err error

		log.Println("[names] creating new name worker")
		nw, err := worker.NewNameWorker()

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

			log.Println("[names] attempting to retrieve new names since ", t.String())
			names, err = nw.RetrieveNewNames(t)
		} else {
			log.Println("[names] fetching latest")
			names, err = nw.RetrieveNames()
		}

		if err != nil {
			return err
		}

		fmt.Println(names)
		prettyPrintNames(names)

		return nw.Shutdown()
	},
}

func init() {
	namesCmd.Flags().StringVarP(&config.NewUsersSince, "since", "s", "",
		"ISO 8601 formatted date [YYYY-MM-DD]")

	namesCmd.Flags().StringVarP(&config.OutputFormat, "format", "t",
		config.DefaultOutputFormat, "output names in prettified json or text format")

	namesCmd.Flags().StringVarP(&config.OutputFile, "outfile", "o", "",
		"write results to file rather than printing to standard output")
}

func prettyPrintNames(n map[string]map[string]bool) error {
	switch config.OutputFormat {
	case "json", "JSON":
		prettyPrintJSON(n)
	case "txt", "text":
		prettyPrintTxt(n)
	default:
		return fmt.Errorf("unsupported format specified: %s", config.OutputFormat)
	}

	return nil
}

func prettyPrintJSON(n map[string]map[string]bool) {
	if config.OutputFile != "" {
		// write to file
	}
}

func prettyPrintTxt(n map[string]map[string]bool) {
	if config.OutputFile != "" {
		// write to file
	}
}
