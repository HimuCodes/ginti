package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
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

func isSpace(c byte) bool {
	return c <= ' ' && (c == ' ' || c == '\n' || c == '\t' || c == '\r' || c == '\v' || c == '\f')
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

				space := isSpace(c)

				if !space && !isInsideWord {
					res.Words++
				}

				isInsideWord = !space
			}
		}

		if err != nil {
			break
		}
	}

	return res
}

func countChunk(f *os.File, start, end int64) Counts {
	res := Counts{}
	res.Bytes = int(end - start)

	buf := make([]byte, 32*1024)
	curr := start
	isInsideWord := false

	// Boundary check: Peek at previous byte to see if we're in the middle of a word
	if start != 0 {
		peek := make([]byte, 1)
		_, err := f.ReadAt(peek, start-1)
		if err == nil && !isSpace(peek[0]) {
			isInsideWord = true
		}
	}

	for curr < end {
		toRead := int64(len(buf))
		if curr+toRead > end {
			toRead = end - curr
		}

		n, err := f.ReadAt(buf[:toRead], curr)
		if n > 0 {
			for i := 0; i < n; i++ {
				c := buf[i]
				if c == '\n' {
					res.Lines++
				}

				space := isSpace(c)
				if !space && !isInsideWord {
					res.Words++
				}
				isInsideWord = !space
			}
			curr += int64(n)
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

	info, err := file.Stat()
	if err != nil {
		return Counts{}, err
	}

	fileSize := info.Size()

	// If the file is small, don't bother with goroutines
	if fileSize < 1024*1024 {
		return GetCounts(file), nil
	}

	numWorkers := runtime.NumCPU()
	chunkSize := fileSize / int64(numWorkers)

	var wg sync.WaitGroup
	results := make(chan Counts, numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = fileSize
		}

		wg.Add(1)
		go func(s, e int64) {
			defer wg.Done()
			results <- countChunk(file, s, e)
		}(start, end)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	total := Counts{}
	for res := range results {
		total = total.Add(res)
	}

	return total, nil
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
