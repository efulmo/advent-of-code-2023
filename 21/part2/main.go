package main

import (
	"fmt"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	charStart  = "S"
	charGarden = "."
	charRock   = "#"

	targetStep = 26_501_365
)

type Coord struct {
	rowIdx, colIdx int
}

type Diff struct {
	row, col int8
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	var startCoord *Coord
	for rowIdx, line := range lines {
		if startIdx := strings.Index(line, charStart); startIdx != -1 {
			startCoord = &Coord{
				rowIdx: rowIdx,
				colIdx: startIdx,
			}
			break
		}
	}
	if startCoord == nil {
		panic(fmt.Errorf("Starting coord isn't found"))
	}
	fmt.Printf("Start is detected at coord %d:%d\n", startCoord.rowIdx+1, startCoord.colIdx+1)

	prevCoords := map[Coord]bool{
		*startCoord: true,
	}

	directionDiffs := []Diff{
		{0, 1},  // right
		{1, 0},  // down
		{0, -1}, // left
		{-1, 0}, // up
	}

	fieldSize := uint(len(lines))
	initialFieldSteps := fieldSize / 2
	firstFieldSteps := initialFieldSteps + fieldSize
	secondFieldSteps := initialFieldSteps + 2*fieldSize
	thirdFieldSteps := initialFieldSteps + 3*fieldSize

	reachedTilesByStep := map[uint]uint{
		// 64:         0,
		initialFieldSteps: 0,
		firstFieldSteps:   0,
		secondFieldSteps:  0,
		thirdFieldSteps:   0,
	}

	for step := uint(1); step <= thirdFieldSteps; step++ {
		curCoords := make(map[Coord]bool)
		for prevCoord := range prevCoords {
			for _, diff := range directionDiffs {
				coord, valid := getGardenCoordIfValid(lines, prevCoord.rowIdx+int(diff.row),
					prevCoord.colIdx+int(diff.col))
				if valid {
					curCoords[coord] = true
				}
			}
		}
		// fmt.Printf("%d tiles reached after %d steps\n", len(curCoords), step)

		if _, tracked := reachedTilesByStep[step]; tracked {
			reachedTilesByStep[step] = uint(len(curCoords))
		}

		prevCoords = curCoords
	}

	fmt.Printf("Simulation results: %v\n", reachedTilesByStep)

	vals := []uint{
		reachedTilesByStep[initialFieldSteps],
		reachedTilesByStep[firstFieldSteps],
		reachedTilesByStep[secondFieldSteps],
		reachedTilesByStep[thirdFieldSteps],
	}
	fmt.Printf("Reached tiles on %d step: %d\n", targetStep, predictNthValue(vals,
		(targetStep-initialFieldSteps)/fieldSize+1))
}

func getGardenCoordIfValid(lines []string, rowIdx, colIdx int) (Coord, bool) {
	linesLen := uint(len(lines))

	normalizedRowIdx := normalizeIdx(rowIdx, linesLen)

	line := lines[normalizedRowIdx]
	colsLen := uint(len(line))

	normalizedColIdx := normalizeIdx(colIdx, colsLen)

	charAtCoord := line[normalizedColIdx : normalizedColIdx+1]
	if charAtCoord != charStart && charAtCoord != charGarden {
		return Coord{}, false
	}

	return Coord{
		rowIdx: rowIdx,
		colIdx: colIdx,
	}, true
}

func normalizeIdx(idx int, idxTotal uint) int {
	normalizedIdx := idx
	if abs(normalizedIdx) >= int(idxTotal) {
		normalizedIdx %= int(idxTotal)
	}
	if normalizedIdx < 0 {
		normalizedIdx += int(idxTotal)
	}

	return normalizedIdx
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func predictNthValue(vals []uint, targetValIdx uint) uint {
	var diffs [][]uint
	diffs = append(diffs, vals)

	for {
		row := diffs[len(diffs)-1]
		if util.SliceContainsSameValue(row, 0) {
			break
		}
		rowLen := uint(len(row))

		newRow := make([]uint, 0, rowLen-1)
		for i := uint(0); i < rowLen-1; i++ {
			diff := row[i+1] - row[i]
			newRow = append(newRow, diff)
		}
		diffs = append(diffs, newRow)
	}

	for _, row := range diffs {
		fmt.Println(row)
	}

	fmt.Printf("Extrapolating to %d step\n", targetValIdx)

	for uint(len(diffs[0])) < targetValIdx {
		for i := uint(len(diffs)) - 1; i > 0; i-- {
			row := diffs[i]
			upperRow := diffs[i-1]
			upperRowNextVal := upperRow[len(upperRow)-1] + row[len(row)-1]
			diffs[i-1] = append(upperRow, upperRowNextVal)
		}
	}

	return diffs[0][len(diffs[0])-1]
}
