package main

import (
	"fmt"
	"os"
)

func main() {
	data, _ := os.ReadFile("./words.txt")

	_ = data

	wordCount := 0

	for _, x := range data {
		if x == ' ' {
			fmt.Println("space detected")
			wordCount++
		}
	}

	wordCount++

	fmt.Println(wordCount)

	fmt.Println("data:", string(data))
}
