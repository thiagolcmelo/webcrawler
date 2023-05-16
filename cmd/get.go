// Package cmd contains the COBRA CLI app used as entry for the webcrawler tool
package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/thiagolcmelo/webcrawler/src"
	"github.com/thiagolcmelo/webcrawler/src/memory"
)

var (
	backoff           time.Duration
	backoffMultiplier int
	format            string
	output            string
	retries           int
	timeout           time.Duration
	verbose           bool
	workers           int
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [flags] domain",
	Short: "It triggers the webcrawler to explore a domain",
	Long: `It triggers the webcrawler to explore a domain
The domain must be provided as a position argument.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		seed := args[0]
		if seed == "" {
			fmt.Println("expected a domain to explore")
			return
		}

		if format != "json" && format != "json-formatted" && format != "raw" {
			fmt.Println("output format can be json, json-formatted or raw")
			return
		}

		if !verbose {
			log.SetOutput(io.Discard)
		}

		ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
		defer cancel()

		frontier := memory.NewFrontier()
		storage := memory.NewStorage()
		events := memory.NewEvents()

		orchestrator := src.NewOrchestrator(
			ctx,
			workers,
			frontier,
			storage,
			events,
			retries,
			backoff,
			backoffMultiplier,
		)
		orchestrator.Start(seed)

		if output != "" {
			// open output file
			outputWriter, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}
			// close fo on exit and check for its returned error
			defer func() {
				if err := outputWriter.Close(); err != nil {
					panic(err)
				}
			}()

			switch format {
			case "raw":
				orchestrator.PrintReport(outputWriter, false, false)
			case "json":
				orchestrator.PrintReport(outputWriter, true, false)
			case "json-formatted":
				orchestrator.PrintReport(outputWriter, true, true)
			}
		} else {
			switch format {
			case "raw":
				orchestrator.PrintReport(os.Stdout, false, false)
			case "json":
				orchestrator.PrintReport(os.Stdout, true, false)
			case "json-formatted":
				orchestrator.PrintReport(os.Stdout, true, true)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().DurationVarP(&backoff, "backoff", "b", 500*time.Millisecond, "how long the client should wait before attempting a retry after a failed request")
	getCmd.Flags().IntVarP(&backoffMultiplier, "backoff-multiplier", "m", 2, "how much the backoff duration should increase between each retry attempt")
	getCmd.Flags().StringVarP(&format, "format", "f", "json", "output format can be json, json-formatted or raw (dummy tree structure)")
	getCmd.Flags().StringVarP(&output, "output", "o", "", "filename to write output to, if empty, it will print to stdout")
	getCmd.Flags().IntVarP(&retries, "retries", "r", 1, "how many times the client should attempt to retry a failed request per individual download")
	getCmd.Flags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "for how long the webcrawler will explore the domain")
	getCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "use it to print logs")
	getCmd.Flags().IntVarP(&workers, "workers", "w", 3, "number of concurrent workers")
}
