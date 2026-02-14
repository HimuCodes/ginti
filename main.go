package main

import (
	"fmt"
	"os"
)

func main() {
	data, _ := os.ReadFile("./words.txt")

	wordCount := CountWords(data)
	fmt.Println(wordCount)
}

func CountWords(data []byte) int {
	wordCount := 0
	if len(data) == 0 {
		return 0
	}
	for _, x := range data {
		if x == ' ' {
			fmt.Println("space detected")
			wordCount++
		}
	}
	wordCount++

	return wordCount
}
