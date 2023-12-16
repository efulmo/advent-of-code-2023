package main

import (
	"fmt"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	instructions := strings.Split(lines[0], ",")
	fmt.Printf("Found %d instructions\n", len(instructions))

	var hashSum uint
	for _, instr := range instructions {
		hash := computeHash(instr)
		util.DebugLog("%s has %d hash\n", instr, hash)

		hashSum += uint(hash)
	}
	fmt.Println("Hash sum:", hashSum)
}

func computeHash(s string) uint8 {
	var hash uint
	for _, r := range s {
		hash += uint(r)
		hash *= 17
		hash %= 256
	}

	return uint8(hash)
}