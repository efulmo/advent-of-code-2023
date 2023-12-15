package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

type Coord struct {
	rowIdx, colIdx uint
}

const (
	horizontal = 1
	vertical   = 2
)

type Mirror struct {
	idx         uint
	orientation uint8
}

const (
	runeAsh  = '.'
	runeRock = '#'
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
	oldMirror := findHorizontalMirrow(pattern, -1)
	if oldMirror == nil {
		oldMirror = findVerticalMirror(pattern, -1)
	}

	if oldMirror == nil {
		panic(fmt.Errorf("Unable to find old mirror"))
	}

	util.DebugLog("Found mirror in %d orientation at idx %d\n", oldMirror.orientation, oldMirror.idx+1)

	smudgesMap := findPotentialSmudges(pattern)
	transposedSmudges := findPotentialSmudges(transposePattern(pattern))
	for tSmudge := range transposedSmudges {
		correctedSmudge := Coord{
			rowIdx: tSmudge.colIdx,
			colIdx: tSmudge.rowIdx,
		}
		smudgesMap[correctedSmudge] = true
	}

	var smudges []Coord
	for smudge := range smudgesMap {
		smudges = append(smudges, smudge)
	}
	slices.SortFunc(smudges, func(c1, c2 Coord) int {
		res := 0
		if c1.rowIdx != c2.rowIdx {
			res = int(c1.rowIdx) - int(c2.rowIdx)
		} else if c1.colIdx != c2.colIdx {
			res = int(c1.colIdx) - int(c2.colIdx)
		}

		// util.DebugLog("Comparing %v and %v: %d\n", c1, c2, res)
		return res
	})

	util.DebugLog("%d potential smudges are found: %v\n", len(smudges), smudges)

	var newMirror *Mirror
	for _, smudge := range smudges {
		util.DebugLog("Testing smudge %v\n", smudge)

		// fix the smudge
		updatedRowBytes := []byte(pattern[smudge.rowIdx])
		oldSmudgeChar := updatedRowBytes[smudge.colIdx]
		newSmudgeChar := runeAsh
		if oldSmudgeChar == runeAsh {
			newSmudgeChar = runeRock
		}
		updatedRowBytes[smudge.colIdx] = byte(newSmudgeChar)

		patternCopy := make([]string, len(pattern))
		copy(patternCopy, pattern)
		patternCopy[smudge.rowIdx] = string(updatedRowBytes)

		ignoreIdx := -1
		if oldMirror.orientation == horizontal {
			ignoreIdx = int(oldMirror.idx)
		}
		newMirror = findHorizontalMirrow(patternCopy, ignoreIdx)

		if newMirror == nil || *newMirror == *oldMirror {
			ignoreIdx := -1
			if oldMirror.orientation == vertical {
				ignoreIdx = int(oldMirror.idx)
			}
			newMirror = findVerticalMirror(patternCopy, ignoreIdx)
		}

		if newMirror != nil && *newMirror != *oldMirror {
			break
		}
	}

	if newMirror == nil || *newMirror == *oldMirror {
		panic(fmt.Errorf("Unable to find mirror after trying all smudges"))
	}

	util.DebugLog("Found new mirror in %d orientation at idx %d\n", newMirror.orientation,
		newMirror.idx+1)

	switch newMirror.orientation {
	case horizontal:
		return uint(newMirror.idx+1) * 100
	case vertical:
		return uint(newMirror.idx) + 1
	default:
		panic(fmt.Errorf("Unknown mirror orientation: %d", newMirror.orientation))
	}
}

func findHorizontalMirrow(pattern []string, ignoreIdx int) *Mirror {
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

			// ignore old mirror if it's passed
			if int(mirrorBelowRowIdx) == ignoreIdx {
				continue
			}

			foundMirrowBelowRowIdx = int(mirrorBelowRowIdx)
			break findMirrorLoop
		}
	}

	if foundMirrowBelowRowIdx != -1 {
		return &Mirror{
			idx:         uint(foundMirrowBelowRowIdx),
			orientation: horizontal,
		}
	}
	return nil
}

func findVerticalMirror(pattern []string, ignoreIdx int) *Mirror {
	transposedPattern := transposePattern(pattern)
	mirror := findHorizontalMirrow(transposedPattern, ignoreIdx)
	if mirror != nil {
		mirror.orientation = vertical
	}
	return mirror
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

func findPotentialSmudges(pattern []string) map[Coord]bool {
	rowsTotal := uint(len(pattern))

	smudges := make(map[Coord]bool)
	for i := uint(0); i < rowsTotal-1; i++ {
		for j := i + 1; j < rowsTotal; j++ {
			colIndexes := findDifferentIndexes(pattern[i], pattern[j])
			if len(colIndexes) != 1 {
				continue
			}

			for _, colIdx := range colIndexes {
				smudges[Coord{
					rowIdx: i,
					colIdx: colIdx,
				}] = true
				smudges[Coord{
					rowIdx: j,
					colIdx: colIdx,
				}] = true
			}
		}
	}

	return smudges
}

func findDifferentIndexes(s1, s2 string) []uint {
	s1Len := uint(len(s1))
	if s1Len != uint(len(s2)) {
		panic(fmt.Errorf("Strings of different lengths: %d, %d", s1Len, len(s2)))
	}

	var differentIndexes []uint
	for i := uint(0); i < s1Len; i++ {
		if s1[i] != s2[i] {
			differentIndexes = append(differentIndexes, i)
		}
	}

	return differentIndexes
}
