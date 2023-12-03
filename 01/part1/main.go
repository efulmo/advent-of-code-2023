package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <input-file-path>\n", os.Args[0])
		return
	}

	inputeFilePath := os.Args[1]
	fmt.Printf("Reading input from file <%s>\n", inputeFilePath)

	data, err := os.ReadFile(inputeFilePath)
	if err != nil {
		fmt.Printf("Error reading file <%s>: %s\n", inputeFilePath, err.Error())
		os.Exit(1)
	}
	fmt.Printf("Read %d bytes\n", len(data))

	lines := strings.Split(string(data), "\n")
	fmt.Printf("Read %d lines\n", len(lines))

	var sum uint
	for lineIdx, line := range lines {
		leftIdx, rightIdx := -1, -1
		
		for charIdx, r := range line {
			if unicode.IsDigit(r) {
				rightIdx = charIdx

				if leftIdx == -1 {
					leftIdx = charIdx
				}
			}
		}

		if leftIdx == -1 || rightIdx == -1 {
			fmt.Printf("Line %d. No numbers found\n", lineIdx)
			os.Exit(1)
		}

		leftDigit := string(line[leftIdx])
		rightDigit := string(line[rightIdx])
		fmt.Printf("%d. %s. Detected %s, %s\n", lineIdx, line, leftDigit, rightDigit)
		
		lineValueStr := leftDigit + rightDigit
		lineValue, err := strconv.Atoi(lineValueStr)
		if err != nil {
			fmt.Printf("Line %d. Unable to parse line value %s\n", lineIdx, lineValueStr)
			os.Exit(1)
		}

		sum += uint(lineValue)
	}

	fmt.Println("Total sum:", sum)
}