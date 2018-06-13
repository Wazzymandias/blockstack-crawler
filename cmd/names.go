package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Wazzymandias/blockstack-crawler/config"
	"github.com/Wazzymandias/blockstack-crawler/names"
	"github.com/Wazzymandias/blockstack-crawler/worker"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var namesCmd = &cobra.Command{
	Use:   "names [OPTIONS]",
	Short: "display information related to users for Blockstack apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		// n is set of names map[namespace][SetOfNames]
		var n map[string]map[string]bool
		var err error

		log.Println("creating new name worker")
		nw, err := worker.NewNameWorker()

		if err != nil {
			return fmt.Errorf("error creating new name worker: %+v", err)
		}

		log.Println("attempting to retrieve names")
		if config.NewUsersSince != "" {
			var t time.Time
			layout := "2006-01-02"

			t, err = time.Parse(layout, config.NewUsersSince)

			if err != nil {
				return fmt.Errorf("error parsing time (YYYY-MM-DD required): %+v", err)
			}

			n, err = nw.RetrieveNewNames(t)
		} else {
			n, err = nw.RetrieveNames()
		}

		if err != nil {
			return fmt.Errorf("error retrieving names: %+v", err)
		}

		err = prettyPrintNames(names.MapToSlice(n))

		if err != nil {
			return fmt.Errorf("error printing names: %+v", err)
		}

		return nw.Shutdown()
	},
}

func init() {
	namesCmd.Flags().StringVarP(&config.NewUsersSince, "since", "s", "",
		"ISO 8601 formatted date [YYYY-MM-DD]")

	namesCmd.Flags().StringVarP(&config.OutputFormat, "format", "f",
		config.DefaultOutputFormat, "output names in prettified json or text format")

	namesCmd.Flags().StringVarP(&config.OutputFile, "outfile", "o", "",
		"write results to file rather than printing to standard output")

	namesCmd.Flags().DurationVarP(&config.Timeout, "timeout", "t", config.DefaultTimeout,
		"timeout for API requests")

	namesCmd.Flags().Uint64VarP(&config.BatchSize, "batch", "b", config.DefaultBatchSize,
		"number of concurrent requests to make to API")
}

func prettyPrintNames(n map[string][]string) error {
	switch config.OutputFormat {
	case "json", "JSON":
		return prettyPrintJSON(n)
	case "txt", "text":
		return prettyPrintTxt(n)
	default:
		return fmt.Errorf("unsupported format specified: %s", config.OutputFormat)
	}
}

func prettyPrintJSON(n map[string][]string) error {
	nb, err := json.MarshalIndent(&n, "", "    ")

	if err != nil {
		return fmt.Errorf("error marhsalling json: %+v", err)
	}

	if config.OutputFile != "" {
		out := config.OutputFile

		if !strings.HasSuffix(out, ".json") {
			out = out + ".json"
		}

		return ioutil.WriteFile(out, nb, 0644)
	}

	fmt.Println(string(nb))
	return nil
}

func prettyPrintTxt(n map[string][]string) error {
	var buf bytes.Buffer

	for k, v := range n {
		buf.WriteString(fmt.Sprintf("%s:\n", k))
		for _, name := range v {
			buf.WriteString(fmt.Sprintf("\t%s\n", name))
		}
	}
	if config.OutputFile != "" {
		return ioutil.WriteFile(config.OutputFile, buf.Bytes(), 0644)
	}

	fmt.Println(buf.String())
	return nil
}
