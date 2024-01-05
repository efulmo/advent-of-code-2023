package main

import (
	"fmt"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	charStart  = "S"
	charGarden = "."

	steps = 64
)

type Coord struct {
	rowIdx, colIdx uint8
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	var startCoord Coord
	for rowIdx, line := range lines {
		if startIdx := strings.Index(line, charStart); startIdx != -1 {
			startCoord = Coord{
				rowIdx: uint8(rowIdx),
				colIdx: uint8(startIdx),
			}
			break
		}
	}

	fmt.Printf("Start is detected at coord %d:%d\n", startCoord.rowIdx+1, startCoord.colIdx+1)

	prevCoords := map[Coord]bool{
		startCoord: true,
	}
	
	for i := uint(0); i < steps; i++ {
		curCoords := make(map[Coord]bool)
		for prevCoord := range prevCoords {
			// up
			coord, valid := getGardenCoordIfValid(lines, int(prevCoord.rowIdx)-1, int(prevCoord.colIdx))
			if valid {
				curCoords[coord] = true
			}

			// down
			coord, valid = getGardenCoordIfValid(lines, int(prevCoord.rowIdx)+1, int(prevCoord.colIdx))
			if valid {
				curCoords[coord] = true
			}

			// right
			coord, valid = getGardenCoordIfValid(lines, int(prevCoord.rowIdx), int(prevCoord.colIdx)+1)
			if valid {
				curCoords[coord] = true
			}

			// left
			coord, valid = getGardenCoordIfValid(lines, int(prevCoord.rowIdx), int(prevCoord.colIdx)-1)
			if valid {
				curCoords[coord] = true
			}
		}
		prevCoords = curCoords
	}

	fmt.Printf("Reachable coords in %d steps: %d\n", steps, len(prevCoords))
}

func getGardenCoordIfValid(lines []string, rowIdx, colIdx int) (Coord, bool) {
	if rowIdx < 0 || colIdx < 0 || rowIdx >= len(lines) {
		return Coord{}, false
	}

	line := lines[rowIdx]
	if colIdx >= len(line) {
		return Coord{}, false
	}

	charAtCoord := line[colIdx : colIdx+1]
	if charAtCoord != charStart && charAtCoord != charGarden {
		return Coord{}, false
	}

	return Coord{
		rowIdx: uint8(rowIdx),
		colIdx: uint8(colIdx),
	}, true
}
