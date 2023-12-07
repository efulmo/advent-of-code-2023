package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Number struct {
	digit uint8
	chars string
}

var numbers = []Number{
	{1, "one"},
	{2, "two"},
	{3, "three"},
	{4, "four"},
	{5, "five"},
	{6, "six"},
	{7, "seven"},
	{8, "eight"},
	{9, "nine"},
}

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
		var leftDigit, rightDigit string

		lineLen := len(line)
		for charIdx := 0; charIdx < lineLen; charIdx++ {
			r := rune(line[charIdx])
			if unicode.IsDigit(r) {
				rightDigit = string(r)

				if leftDigit == "" {
					leftDigit = rightDigit
				}

				continue
			}

			charDigit := parseCharDigit(line, charIdx)
			if charDigit != nil {
				rightDigit = strconv.Itoa(int(charDigit.digit))

				if leftDigit == "" {
					leftDigit = rightDigit
				}
			}
		}

		if leftDigit == "" || rightDigit == "" {
			fmt.Printf("Line %d. No numbers found: %s\n", lineIdx, line)
			os.Exit(1)
		}

		lineValueStr := leftDigit + rightDigit
		lineValue, err := strconv.Atoi(lineValueStr)
		if err != nil {
			fmt.Printf("Line %d. Unable to parse line value %s\n", lineIdx, lineValueStr)
			os.Exit(1)
		}

		fmt.Printf("%d. %s. Detected %s, %s => %d\n", lineIdx, line, leftDigit, rightDigit, lineValue)

		sum += uint(lineValue)
	}

	fmt.Println("Total sum:", sum)
}

func parseCharDigit(line string, idx int) *Number {
	lineLen := len(line)
	for _, num := range numbers {
		if idx+len(num.chars) > lineLen {
			continue
		}

		if stringIsAtIdx(line, num.chars, idx) {
			return &num
		}
	}

	return nil
}

func stringIsAtIdx(input, sample string, idx int) bool {
	if idx < 0 {
		fmt.Printf("stringIsAtIdx: input=%s, sample=%s, idx=%d. idx < 0\n", input, sample, idx)
		os.Exit(1)
	}

	inputLen := len(input)
	sampleLen := len(sample)

	// idx is out of input length
	if idx >= inputLen {
		return false
	}

	// sample is too long to match
	if idx+sampleLen > inputLen {
		return false
	}

	for i := 0; i < sampleLen; i++ {
		if input[idx+i] != sample[i] {
			return false
		}
	}

	return true
}
