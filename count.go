package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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

func PrintHeader(w io.Writer, opts DisplayOptions) {
	labels := []string{}

	if opts.ShouldShowLines() {
		labels = append(labels, "lines")
	}

	if opts.ShouldShowWords() {
		labels = append(labels, "words")
	}

	if opts.ShouldShowBytes() {
		labels = append(labels, "bytes")
	}

	lines := strings.Join(labels, "\t") + "\t"

	fmt.Fprint(w, lines+"\n")
}

func (c Counts) Print(w io.Writer, opts DisplayOptions, suffixes ...string) {
	stats := []string{}

	if opts.ShouldShowLines() {
		stats = append(stats, strconv.Itoa(c.Lines))
	}

	if opts.ShouldShowWords() {
		stats = append(stats, strconv.Itoa(c.Words))
	}

	if opts.ShouldShowBytes() {
		stats = append(stats, strconv.Itoa(c.Bytes))
	}

	line := strings.Join(stats, "\t") + "\t"

	fmt.Fprint(w, line)

	suffixStr := strings.Join(suffixes, " ")
	if suffixStr != "" {
		fmt.Fprintf(w, " %s", suffixStr)
	}

	fmt.Fprint(w, "\n")
}

func GetCounts(f io.Reader) Counts {
	res := Counts{}

	buf := make([]byte, 32*1024)
	isInsideWord := false

	for {
		n, err := f.Read(buf)
		if n > 0 {
			res.Bytes += n

			for i := 0; i < n; i++ {
				c := buf[i]

				if c == '\n' {
					res.Lines++
				}

				isSpace := c <= ' ' && (c == ' ' || c == '\n' || c == '\t' || c == '\r' || c == '\v' || c == '\f')

				if !isSpace && !isInsideWord {
					res.Words++
				}

				isInsideWord = !isSpace
			}
		}

		if err != nil {
			break
		}
	}

	return res
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
