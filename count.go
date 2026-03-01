package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Counts struct {
	Bytes int
	Words int
	Lines int
}

// Add will modify the values of the count by
// adding the values from the other.
func (c Counts) Add(other Counts) Counts {
	c.Bytes += other.Bytes
	c.Words += other.Words
	c.Lines += other.Lines

	return c
}

func (c Counts) Print(w io.Writer, filenames ...string) {
	fmt.Fprintf(w, "%8d %8d %8d", c.Lines, c.Words, c.Bytes)

	for _, filename := range filenames {
		fmt.Fprintf(w, " %s", filename)
	}
	fmt.Fprintln(w)
}

func GetCounts(f io.ReadSeeker) Counts {
	const offsetStart = 0

	byteCount := CountBytes(f)
	f.Seek(offsetStart, io.SeekStart)

	wordCount := CountWords(f)
	f.Seek(offsetStart, io.SeekStart)

	lineCount := CountLines(f)
	f.Seek(offsetStart, io.SeekStart)

	return Counts{
		Bytes: byteCount,
		Words: wordCount,
		Lines: lineCount,
	}
}

func CountFile(filename string) (Counts, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Counts{}, err
	}

	defer file.Close()

	counts := GetCounts(file)

	return counts, nil
}

func CountWords(file io.Reader) int {
	wordCount := 0

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		wordCount++
	}
	return wordCount
}

func CountLines(r io.Reader) int {
	linesCount := 0

	reader := bufio.NewReader(r)

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}

		if r == '\n' {
			linesCount++
		}

	}

	return linesCount
}

func CountBytes(r io.Reader) int {
	byteCount, _ := io.Copy(io.Discard, r)
	return int(byteCount)
}
