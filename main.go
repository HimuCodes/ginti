package main

import (
	"fmt"
	"os"
)

func main() {
	data, _ := os.ReadFile("./words.txt")

	_ = data

	fmt.Println("data:", string(data))
}
