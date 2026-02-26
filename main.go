package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	totalLines := 0
	totalWords := 0
	totalBytes := 0

	filenames := os.Args[1:]

	didError := false

	for _, filename := range filenames {

		counts, err := CountFile(filename)
		if err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "counter:", err)
			continue
		}

		totalLines += counts.Lines
		totalWords += counts.Words
		totalBytes += counts.Bytes

		fmt.Printf("%8d %8d %8d %s\n", counts.Lines, counts.Words, counts.Bytes, filename)
	}

	if len(filenames) == 0 {
		counts := GetCounts(os.Stdin)
		fmt.Printf("%8d %8d %8d\n", counts.Lines, counts.Words, counts.Bytes)
	}

	if len(filenames) > 1 {
		fmt.Printf("%8d %8d %8d total\n", totalLines, totalWords, totalBytes)
	}

	if didError {
		os.Exit(1)
	}
}
