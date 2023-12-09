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

func ParseUintOrPanic(s string) uint {
	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(errors.Join(fmt.Errorf("Failed to parse <%s> as uint", s), err))
	}
	return uint(u)
}

func StringsToUints(strs []string) []uint {
	var res []uint
	for _, s := range strs {
		res = append(res, ParseUintOrPanic(s))
	}
	return res
}

func StringsToInts(strs []string) []int {
	var res []int
	for _, s := range strs {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(errors.Join(fmt.Errorf("Failed to parse <%s> as int", s), err))
		}
		res = append(res, i)
	}
	return res
}

func SliceContainsSameValue[T comparable](sl []T, value T) bool {
	for _, u := range sl {
		if u != value {
			return false
		}
	}
	return true
}
