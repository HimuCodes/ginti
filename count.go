package main

import (
	"bufio"
	"io"
	"os"
)

type Counts struct {
	Bytes int
	Words int
	Lines int
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
