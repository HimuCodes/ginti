package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	filename := "utf8.txt"

	log.SetFlags(0)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("failed to read file:", err)
	}

	wordCount := CountWords(file)
	fmt.Println(wordCount)
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
