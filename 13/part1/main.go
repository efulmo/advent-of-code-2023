package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	patterns := groupPatterns(lines)
	fmt.Printf("Detected %d patterns\n", len(patterns))

	var pointsSum uint
	for pIdx, pattern := range patterns {
		util.DebugLog("%d.\n", pIdx+1)
		pointsSum += calculatePatternPoints(pattern)
	}
	fmt.Println("Pattern points sum:", pointsSum)
}

func groupPatterns(lines []string) [][]string {
	var patterns [][]string
	var pattern []string

	for _, line := range lines {
		if len(line) > 0 {
			pattern = append(pattern, line)
		} else {
			patterns = append(patterns, pattern)
			pattern = make([]string, 0)
		}
	}

	if len(pattern) != 0 {
		patterns = append(patterns, pattern)
	}

	return patterns
}

func calculatePatternPoints(pattern []string) uint {
	foundMirrowBelowRowIdx := findHorizontalMirrow(pattern)
	if foundMirrowBelowRowIdx != -1 {
		util.DebugLog("Found mirror at row %d\n", foundMirrowBelowRowIdx+1)
		return uint(foundMirrowBelowRowIdx+1) * 100
	}

	foundMirrorRightToColIdx := findVerticalMirror(pattern)
	if foundMirrorRightToColIdx != -1 {
		util.DebugLog("Found mirror at col %d\n", foundMirrorRightToColIdx+1)
		return uint(foundMirrorRightToColIdx + 1)
	}

	panic(fmt.Errorf("Unable to find mirror"))
}

func findHorizontalMirrow(pattern []string) int {
	rowsTotal := uint(len(pattern))

	sameRowsMap := make(map[uint][]uint)
	for i := uint(0); i < rowsTotal-1; i++ {
		for j := i + 1; j < rowsTotal; j++ {
			if pattern[i] == pattern[j] {
				sameRowsMap[i] = append(sameRowsMap[i], j)
				sameRowsMap[j] = append(sameRowsMap[j], i)
			}
		}
	}

	util.DebugLog("Same rows: %v\n", sameRowsMap)

	foundMirrowBelowRowIdx := -1
findMirrorLoop:
	for rowIdx := uint(0); rowIdx < rowsTotal; rowIdx++ {
		sameRows, exist := sameRowsMap[rowIdx]
		if !exist {
			continue
		}

	verifyMirrowPositionLoop:
		for _, sameRowIdx := range sameRows {
			minRowIdx, maxRowIdx := min(rowIdx, sameRowIdx), max(rowIdx, sameRowIdx)
			if minRowIdx != 0 && maxRowIdx != rowsTotal-1 {
				continue
			}

			mirroredRowsCount := maxRowIdx - minRowIdx + 1
			mirrorBelowRowIdx := minRowIdx + mirroredRowsCount/2 - 1

			for i := uint(0); i < mirroredRowsCount/2; i++ {
				if !slices.Contains(sameRowsMap[mirrorBelowRowIdx-i], mirrorBelowRowIdx+1+i) {
					continue verifyMirrowPositionLoop
				}
			}

			foundMirrowBelowRowIdx = int(mirrorBelowRowIdx)
			break findMirrorLoop
		}
	}

	return foundMirrowBelowRowIdx
}

func findVerticalMirror(pattern []string) int {
	transposedPattern := transposePattern(pattern)
	return findHorizontalMirrow(transposedPattern)
}

func transposePattern(pattern []string) []string {
	colCount := uint(len(pattern[0]))
	rowCount := uint(len(pattern))

	var newPattern []string
	for colIdx := uint(0); colIdx < colCount; colIdx++ {
		var newPatternRow strings.Builder
		for rowIdx := uint(0); rowIdx < rowCount; rowIdx++ {
			newPatternRow.WriteByte(pattern[rowIdx][colIdx])
		}
		newPattern = append(newPattern, newPatternRow.String())
	}

	return newPattern
}
