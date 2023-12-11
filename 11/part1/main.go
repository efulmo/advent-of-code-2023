package main

import (
	"fmt"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

type Galaxy struct {
	ID     uint
	rowIdx uint
	colIdx uint
}

type GalaxyPair struct {
	a uint
	b uint
}

const (
	runeGalaxy = '#'
	runeDot    = '.'
)
const charDot = string(runeDot)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	galaxies := findGalaxies(lines)
	galaxiesLen := uint(len(galaxies))
	fmt.Printf("%d galaxies parsed\n", galaxiesLen)
	// fmt.Println(galaxies)

	var expandingRows []uint
	for rowIdx, line := range lines {
		if strings.Count(line, charDot) == len(line) {
			expandingRows = append(expandingRows, uint(rowIdx))
		}
	}
	fmt.Printf("Found %d expanding rows\n", len(expandingRows))
	// fmt.Println(expandingRows)

	colCount := uint(len(lines[0]))
	var expandingCols []uint
byCols:
	for colIdx := uint(0); colIdx < colCount; colIdx++ {
		for _, line := range lines {
			if line[colIdx] != runeDot {
				continue byCols
			}
		}
		expandingCols = append(expandingCols, colIdx)
	}
	fmt.Printf("Found %d expanding cols\n", len(expandingCols))
	// fmt.Println(expandingCols)

	var pairs []GalaxyPair
	for i := uint(1); i <= galaxiesLen; i++ {
		for j := uint(i + 1); j <= galaxiesLen; j++ {
			pairs = append(pairs, GalaxyPair{i, j})
		}
	}
	fmt.Printf("%d galaxy pairs are built\n", len(pairs))
	// fmt.Println(pairs)

	pathLengthSum := uint(0)
	for _, pair := range pairs {
		g1 := galaxies[pair.a]
		g2 := galaxies[pair.b]
		minRowIdx := min(g1.rowIdx, g2.rowIdx)
		maxRowIdx := max(g1.rowIdx, g2.rowIdx)
		minColIdx := min(g1.colIdx, g2.colIdx)
		maxColIdx := max(g1.colIdx, g2.colIdx)
		pathLength := (maxRowIdx - minRowIdx) + (maxColIdx - minColIdx) + 
			countBetween(expandingRows, minRowIdx, maxRowIdx) + 
			countBetween(expandingCols, minColIdx, maxColIdx)

		// fmt.Printf("Path between G%d(%d:%d) and G%d(%d:%d) is %d\n", g1.ID, g1.rowIdx, g1.colIdx, 
		// 	g2.ID, g2.rowIdx, g2.colIdx, pathLength)
		pathLengthSum += pathLength
	}
	fmt.Println("Paths length sum:", pathLengthSum)
}

func findGalaxies(lines []string) map[uint]Galaxy {
	galaxies := make(map[uint]Galaxy)
	for rowIdx, line := range lines {
		for colIdx, r := range line {
			if r == runeGalaxy {
				ID := uint(len(galaxies) + 1)
				galaxies[ID] = Galaxy{
					ID:     ID,
					rowIdx: uint(rowIdx),
					colIdx: uint(colIdx),
				}
			}
		}
	}

	return galaxies
}

func countBetween(expandingLines []uint, fromLineIdx, toLineIdx uint) uint {
	cnt := uint(0)
	for _, lineIdx := range expandingLines {
		if lineIdx > fromLineIdx && lineIdx < toLineIdx {
			cnt++
		}
	}

	return cnt
}