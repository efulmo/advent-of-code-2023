package util

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"slices"
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

func ParseIntOrPanic(s string) int {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(errors.Join(fmt.Errorf("Failed to parse <%s> as int", s), err))
	}
	return int(i)
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

func DebugLog(str string, params ...any) {
	fmt.Printf(str, params...)
}

func MapKeysToSortedSlice[K cmp.Ordered, V any](m map[K]V) []K {
	sl := MapKeysToSlice(m)
	slices.Sort(sl)
	return sl
}

func MapKeysToSlice[K comparable, V any](m map[K]V) []K {
	sl := make([]K, 0, len(m))
	for key := range m {
		sl = append(sl, key)
	}
	return sl
}

func GetAnyMapKey[K comparable, V any](m map[K]V) (K, error) {
	for k := range m {
		return k, nil
	}

	var zeroValueKey K
	return zeroValueKey, errors.New("Unable to get key from empty map")
}