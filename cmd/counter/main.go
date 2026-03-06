package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/HimuCodes/ginti/counter"
	"github.com/HimuCodes/ginti/display"
)

func main() {
	showLines := flag.Bool(
		"l",
		false,
		"Used to toggle whether or not to show the line count",
	)

	showBytes := flag.Bool(
		"c",
		false,
		"Used to toggle whether or not to show the byte count",
	)

	showWords := flag.Bool(
		"w",
		false,
		"Used to toggle whether or not to show the word count",
	)

	showHeader := flag.Bool(
		"header",
		false,
		"Used to toggle whether or not to show the headers for all the data being parsed",
	)

	flag.Parse()

	opts := display.NewOptions(*showBytes, *showWords, *showLines, *showHeader)

	log.SetFlags(0)

	wr := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', tabwriter.AlignRight)

	totals := counter.Counts{}

	filenames := flag.Args()
	didError := false

	if opts.ShouldShowHeader() {
		counter.PrintHeader(wr, opts)
	}

	for _, filename := range filenames {

		counts, err := counter.CountFile(filename)
		if err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "counter:", err)
			continue
		}

		totals = totals.Add(counts)

		counts.Print(wr, opts, filename)
	}

	if len(filenames) == 0 {
		counter.GetCounts(os.Stdin).Print(wr, opts)
	}

	if len(filenames) > 1 {
		totals.Print(wr, opts, "total")
	}

	wr.Flush()
	if didError {
		os.Exit(1)
	}
}
