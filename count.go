package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
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

	reader := bufio.NewReader(f)

	isInsideWord := false
	for {
		r, size, err := reader.ReadRune()
		if err != nil {
			break
		}

		res.Bytes += size
		if r == '\n' {
			res.Lines++
		}

		isSpace := unicode.IsSpace(r)

		if !isSpace && !isInsideWord {
			res.Words++
		}

		isInsideWord = !isSpace
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
