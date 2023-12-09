package main

import (
	"fmt"
	"unicode"

	"github.com/efulmo/advent-of-code-2023/util"
)

type PartID struct {
	lineIdx  uint
	startIdx uint // including
	endIdx   uint // excluding
	number   uint
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	var partNumberSum uint

	for lineIdx, line := range lines {
		partIDs := parsePartIDs(line, uint(lineIdx))

		var adjacentPartIDs []PartID
		for _, partID := range partIDs {
			if isPartIDAdjacent(lines, partID) {
				adjacentPartIDs = append(adjacentPartIDs, partID)

				partNumberSum += partID.number
			}
		}

		// fmt.Printf("%d. %s. Parsed(%d): %v. Adj(%d): %v\n", lineIdx + 1, line,
		// 	len(partIDs), numberRegionsToNumbers(partIDs),
		// 	len(adjacentPartIDs), numberRegionsToNumbers(adjacentPartIDs))

		// if len(partIDs) != len(unadjacentPartIDs) {
		unadjacentPartIDs := getUnadjacentPartIDs(partIDs, adjacentPartIDs)
		fmt.Printf("%d.Parsed(%d): %v. Adj(%d): %v. Not(%d): %v\n", lineIdx+1,
			len(partIDs), partIDsToNumbers(partIDs),
			len(adjacentPartIDs), partIDsToNumbers(adjacentPartIDs),
			len(unadjacentPartIDs), partIDsToNumbers(unadjacentPartIDs))
		// }
	}

	fmt.Println("Part number sum:", partNumberSum)
}

const NilIdx = -1

func parsePartIDs(line string, lineIdx uint) []PartID {
	var partIDs []PartID

	startRegionIdx := NilIdx
	for charIdx, r := range line {
		if unicode.IsDigit(r) {
			if startRegionIdx == NilIdx {
				startRegionIdx = charIdx
			}
		} else {
			if startRegionIdx != NilIdx {
				partIDs = append(partIDs, PartID{
					lineIdx:  lineIdx,
					startIdx: uint(startRegionIdx),
					endIdx:   uint(charIdx),
					number:   parsePartIDNumber(line, uint(startRegionIdx), uint(charIdx)),
				})

				// fmt.Printf("Parsed region %d:%d\n", startRegionIdx, charIdx)

				startRegionIdx = NilIdx
			}
		}
	}

	if startRegionIdx != NilIdx {
		partIDs = append(partIDs, PartID{
			lineIdx:  lineIdx,
			startIdx: uint(startRegionIdx),
			endIdx:   uint(len(line)),
			number:   parsePartIDNumber(line, uint(startRegionIdx), uint(len(line))),
		})
	}

	return partIDs
}

func parsePartIDNumber(line string, startIdx, endIdx uint) uint {
	partIDStr := line[startIdx:endIdx]

	// potential int truncation
	if len(partIDStr) > 6 {
		panic(fmt.Errorf("Potentionally too big part ID for int: %s", partIDStr))
	}

	return util.ParseUintOrPanic(partIDStr)
}

func partIDsToNumbers(partIDs []PartID) []uint {
	numbers := make([]uint, 0, len(partIDs))
	for _, partID := range partIDs {
		numbers = append(numbers, partID.number)
	}

	return numbers
}

func isPartIDAdjacent(lines []string, part PartID) bool {
	runesAround := make([]rune, 0, 12)
	lineIdx := int(part.lineIdx)
	startIdx := int(part.startIdx)
	endIdx := int(part.endIdx)

	// left & right
	runesAround = append(runesAround, getRuneAtOrDot(lines, lineIdx, startIdx-1))
	runesAround = append(runesAround, getRuneAtOrDot(lines, lineIdx, endIdx))

	// top & bottom
	partIDLen := endIdx - startIdx
	for i := 0; i < partIDLen; i++ {
		runesAround = append(runesAround, getRuneAtOrDot(lines, lineIdx-1, startIdx+i))
		runesAround = append(runesAround, getRuneAtOrDot(lines, lineIdx+1, startIdx+i))
	}

	// top diagonals
	runesAround = append(runesAround, getRuneAtOrDot(lines, lineIdx-1, startIdx-1))
	runesAround = append(runesAround, getRuneAtOrDot(lines, lineIdx-1, endIdx))

	// bottom diagonals
	runesAround = append(runesAround, getRuneAtOrDot(lines, lineIdx+1, startIdx-1))
	runesAround = append(runesAround, getRuneAtOrDot(lines, lineIdx+1, endIdx))

	for _, r := range runesAround {
		if isSymbol(r) {
			return true
		}
	}

	return false
}

func getRuneAtOrDot(lines []string, rowIdx, colIdx int) rune {
	// negative indexes
	if rowIdx < 0 || colIdx < 0 {
		return Dot
	}

	// too big rowIdx
	if rowIdx >= len(lines) {
		return Dot
	}

	// too big colIdx
	line := lines[rowIdx]
	if colIdx >= len(line) {
		return Dot
	}

	return rune(line[colIdx])
}

const Dot = '.'

func isSymbol(ch rune) bool {
	return ch != Dot && !unicode.IsDigit(ch)
}

func getUnadjacentPartIDs(allPartIDs, adjacentPartIDs []PartID) []PartID {
	var unadjacentPartIDs []PartID

	m := make(map[PartID]bool)
	for _, region := range adjacentPartIDs {
		m[region] = true
	}

	for _, region := range allPartIDs {
		if !m[region] {
			unadjacentPartIDs = append(unadjacentPartIDs, region)
		}
	}

	return unadjacentPartIDs
}
