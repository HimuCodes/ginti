package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"text/tabwriter"
)

type DisplayOptions struct {
	ShowBytes  bool
	ShowWords  bool
	ShowLines  bool
	ShowHeader bool
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

	flag.BoolVar(
		&opts.ShowHeader,
		"header",
		false,
		"Used to toggle whether or not to show the headers for all the data being parsed",
	)

	flag.Parse()

	log.SetFlags(0)

	wr := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', tabwriter.AlignRight)

	totals := Counts{}

	filenames := flag.Args()
	didError := false

	wg := sync.WaitGroup{}
	wg.Add(len(filenames))

	l := sync.Mutex{}

	if opts.ShowHeader {
		PrintHeader(wr, opts)
	}

	for _, filename := range filenames {
		go func() {
			defer wg.Done()

			counts, err := CountFile(filename)
			if err != nil {
				didError = true
				fmt.Fprintln(os.Stderr, "counter:", err)

				return
			}

			l.Lock()
			defer l.Unlock()
			// state mutated below
			totals = totals.Add(counts)

			counts.Print(wr, opts, filename)
		}()
	}
	wg.Wait()

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
