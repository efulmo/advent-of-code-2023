package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	runeOperationAdd    = '='
	runeOperationRemove = '-'
)

type Instruction struct {
	lensLabel   string
	operation   rune
	focalLength uint8
}

type Lens struct {
	label  string
	length uint8
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	instructionsStr := strings.Split(lines[0], ",")
	fmt.Printf("Found %d instructions\n", len(instructionsStr))

	boxes := make(map[uint8][]Lens)
	for instrIdx, instrStr := range instructionsStr {
		instr := parseInstruction(instrStr)
		util.DebugLog("Parsed instruction: %v\n", instr)

		boxIdx := computeHash(instr.lensLabel)
		lenses := boxes[boxIdx]
		findLensByLabel := func(lens Lens) bool {
			return lens.label == instr.lensLabel
		}

		switch instr.operation {
		case runeOperationAdd:
			newLens := Lens{
				label:  instr.lensLabel,
				length: instr.focalLength,
			}
			lensIdx := slices.IndexFunc(lenses, findLensByLabel)

			if lensIdx == -1 {
				boxes[boxIdx] = append(lenses, newLens)
			} else {
				lenses[lensIdx] = newLens
				boxes[boxIdx] = lenses
			}
		case runeOperationRemove:
			boxes[boxIdx] = slices.DeleteFunc(lenses, findLensByLabel)
		default:
			panic(fmt.Errorf("Unknown command at idx: %d", instrIdx))
		}
	}

	fmt.Println("Boxes:")
	fmt.Println(boxes)

	var totalFocusingPower uint
	for boxIdx, lenses := range boxes {
		for lensIdx, lens := range lenses {
			lensPower := (uint(boxIdx) + 1) * uint(lensIdx+1) * uint(lens.length)
			totalFocusingPower += lensPower
		}
	}
	fmt.Println("Total focusing power:", totalFocusingPower)
}

func parseInstruction(s string) Instruction {
	var labelRunes []rune
	for _, r := range s {
		if unicode.IsLetter(r) {
			labelRunes = append(labelRunes, r)
		} else {
			break
		}
	}

	operation := s[len(labelRunes)]
	var lensLength uint8
	if operation == runeOperationAdd {
		lensLengthStr := string(s[len(labelRunes)+1])

		parsedLensLength, err := strconv.ParseUint(lensLengthStr, 10, 4)
		if err != nil {
			panic(errors.Join(fmt.Errorf("Failed to parse lens length from %s", lensLengthStr)))
		}
		lensLength = uint8(parsedLensLength)
	}

	return Instruction{
		lensLabel:   string(labelRunes),
		operation:   rune(operation),
		focalLength: uint8(lensLength),
	}
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
