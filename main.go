package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	filename := "words.txt"

	log.SetFlags(0)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("failed to read file:", err)
	}

	wordCount := CountWordsInFile(file)
	fmt.Println(wordCount)
}

func CountWordsInFile(file *os.File) int {
	wordCount := 0
	isInsideWord := false

	const bufferSize = 1
	buffer := make([]byte, bufferSize)

	leftover := []byte{}

	for {
		size, err := file.Read(buffer)
		if err != nil {
			break
		}

		subBuffer := append(leftover, buffer[:size]...)

		for len(subBuffer) > 0 {

			r, rsize := utf8.DecodeRune(subBuffer)

			if r == utf8.RuneError {
				break
			}
			subBuffer = subBuffer[rsize:]

			if !unicode.IsSpace(r) && !isInsideWord {
				wordCount++
			}

			isInsideWord = !unicode.IsSpace(r)
		}

		leftover = nil
		leftover = append(leftover, subBuffer...)
	}

	return wordCount
}

func CountWords(data []byte) int {
	words := strings.Fields(string(data))
	return len(words)
}
