package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	totals := Counts{}

	filenames := os.Args[1:]

	didError := false

	for _, filename := range filenames {

		counts, err := CountFile(filename)
		if err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "counter:", err)
			continue
		}

		totals = totals.Add(counts)

		counts.Print(os.Stdout, filename)
	}

	if len(filenames) == 0 {
		GetCounts(os.Stdin).Print(os.Stdout)
	}

	if len(filenames) > 1 {
		totals.Print(os.Stdout, "total")
	}

	if didError {
		os.Exit(1)
	}
}
