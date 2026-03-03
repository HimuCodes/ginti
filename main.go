package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

type DisplayOptions struct {
	ShowBytes bool
	ShowWords bool
	ShowLines bool
}

func (d DisplayOptions) ShouldShowBytes() bool {
	if !d.ShowBytes && !d.ShowWords && !d.ShowLines {
		return true
	}

	return d.ShowBytes
}

func (d DisplayOptions) ShouldShowWords() bool {
	if !d.ShowBytes && !d.ShowWords && !d.ShowLines {
		return true
	}

	return d.ShowWords
}
func (d DisplayOptions) ShouldShowLines() bool {
	if !d.ShowBytes && !d.ShowWords && !d.ShowLines {
		return true
	}

	return d.ShowLines
}
func main() {
	opts := DisplayOptions{}

	flag.BoolVar(
		&opts.ShowLines,
		"l",
		false,
		"Used to toggle whether or not to show the line count",
	)

	flag.BoolVar(
		&opts.ShowBytes,
		"c",
		false,
		"Used to toggle whether or not to show the byte count",
	)

	flag.BoolVar(
		&opts.ShowWords,
		"w",
		false,
		"Used to toggle whether or not to show the word count",
	)

	flag.Parse()

	log.SetFlags(0)

	wr := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', tabwriter.AlignRight)

	totals := Counts{}

	filenames := flag.Args()
	didError := false

	for _, filename := range filenames {

		counts, err := CountFile(filename)
		if err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "counter:", err)
			continue
		}

		totals = totals.Add(counts)

		counts.Print(wr, opts, filename)
	}

	if len(filenames) == 0 {
		GetCounts(os.Stdin).Print(wr, opts)
	}

	if len(filenames) > 1 {
		totals.Print(wr, opts, "total")
	}

	wr.Flush()
	if didError {
		os.Exit(1)
	}
}
