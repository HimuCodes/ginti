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

type FileCountsResult struct {
	counts   Counts
	filename string
	err      error
	idx      int
}

func (d DisplayOptions) AnyTrue() bool {
	return d.ShowBytes || d.ShowWords || d.ShowLines
}

func (d DisplayOptions) ShouldShowBytes() bool {
	return !d.AnyTrue() || d.ShowBytes
}

func (d DisplayOptions) ShouldShowWords() bool {
	return !d.AnyTrue() || d.ShowWords
}

func (d DisplayOptions) ShouldShowLines() bool {
	return !d.AnyTrue() || d.ShowLines
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

	if opts.ShowHeader {
		PrintHeader(wr, opts)
	}

	ch := CountFiles(filenames)

	results := make([]FileCountsResult, len(filenames))

	for res := range ch {
		results[res.idx] = res
	}

	for _, res := range results {
		if res.err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "counter:", res.err)
			continue
		}

		totals = totals.Add(res.counts)
		res.counts.Print(wr, opts, res.filename)
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

func CountFiles(filenames []string) <-chan FileCountsResult {
	ch := make(chan FileCountsResult)
	wg := sync.WaitGroup{}

	for i, filename := range filenames {
		wg.Add(1)
		go func(fname string) {
			defer wg.Done()

			counts, err := CountFile(fname)
			ch <- FileCountsResult{
				counts:   counts,
				filename: fname,
				err:      err,
				idx:      i,
			}
		}(filename)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
