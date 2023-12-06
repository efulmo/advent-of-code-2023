package util

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadInputFile() ([]string, error) {
	if len(os.Args) != 2 {
		return nil, fmt.Errorf("Usage: %s <input-file-path>", os.Args[0])
	}

	inputeFilePath := os.Args[1]
	fmt.Printf("Reading input from file <%s>\n", inputeFilePath)

	data, err := os.ReadFile(inputeFilePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading file <%s>: %s", inputeFilePath, err.Error())
	}

	fmt.Printf("Read %d bytes\n", len(data))

	lines := strings.Split(string(data), "\n")
	fmt.Printf("Read %d lines\n", len(lines))

	return lines, nil
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func AtoiOrPanic(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		PanicOnError(errors.Join(fmt.Errorf("Failed to parse <%s> to int", s), err))
	}

	return i
}