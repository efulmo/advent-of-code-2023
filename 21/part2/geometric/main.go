package main

import (
	"fmt"
	"math"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	charStart  = "S"
	charGarden = "."
	charRock   = "#"

	steps = 6
	// steps = 64
)

type Coord struct {
	rowIdx, colIdx uint8
}

type Diff struct {
	row, col int8
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	accesibleTilesCount := uint(0)
	var startCoord *Coord

	// find start coord and count non-rock tiles
	for rowIdx, line := range lines {
		for colIdx, r := range line {
			b := byte(r)
			if b == charStart[0] {
				startCoord = &Coord{
					rowIdx: uint8(rowIdx),
					colIdx: uint8(colIdx),
				}
				accesibleTilesCount++
			} else if b == charGarden[0] {
				accesibleTilesCount++
			}
		}
	}

	if startCoord == nil {
		panic(fmt.Errorf("Starting coord isn't found"))
	}
	fmt.Printf("Start is detected at coord %d:%d\n", startCoord.rowIdx+1, startCoord.colIdx+1)

	directionDiffs := []Diff{
		{0, 1},  // right
		{1, 0},  // down
		{0, -1}, // left
		{-1, 0}, // up
	}
	minDistanceByCoord := map[Coord]uint{
		*startCoord: 0,
	}
	prevCoords := map[Coord]bool{
		*startCoord: true,
	}

	step := uint(1)
	for {
		distancesMeasured := len(minDistanceByCoord)
		curCoords := make(map[Coord]bool)

		for prevCoord := range prevCoords {
			for _, diff := range directionDiffs {
				coord, valid := getGardenCoordIfValid(lines, int(prevCoord.rowIdx)+int(diff.row),
					int(prevCoord.colIdx)+int(diff.col))
				if valid {
					curCoords[coord] = true

					curDistance, alreadyVisited := minDistanceByCoord[coord]
					minDistance := step
					if alreadyVisited {
						minDistance = min(curDistance, step)
					}
					minDistanceByCoord[coord] = minDistance
				}
			}
		}
		prevCoords = curCoords

		// stop when no new garden tile is added during the step
		if len(minDistanceByCoord) == distancesMeasured {
			break
		}

		step++
	}

	fmt.Printf("Min distances are calculated for %d garden tiles out of %d in %d steps\n",
		len(minDistanceByCoord), accesibleTilesCount, step)

	for rowIdx, row := range lines {
		for colIdx := range row {
			coord := Coord{uint8(rowIdx), uint8(colIdx)}
			distance := minDistanceByCoord[coord]
			if distance == 0 && row[colIdx] != charRock[0] {
				fmt.Printf("%d:%d: %d\n", rowIdx+1, colIdx+1, distance)
			}
		}
	}

	var evenTilesCount, evenFarTilesCount uint
	for _, distance := range minDistanceByCoord {
		if distance%2 == 0 {
			evenTilesCount++

			if distance > 65 {
				evenFarTilesCount++
			}
		}
	}
	oddTilesCount := uint(len(minDistanceByCoord)) - evenTilesCount
	oddFarTilesCount := uint(len(minDistanceByCoord)) - evenFarTilesCount

	const gardenRepeatedCount = 202300
	visitedTiles := math.Pow(gardenRepeatedCount+1, 2)*float64(oddTilesCount) +
		math.Pow(gardenRepeatedCount, 2)*float64(evenTilesCount) -
		(gardenRepeatedCount+1)*float64(oddFarTilesCount) +
		gardenRepeatedCount*float64(evenFarTilesCount)

	fmt.Printf("Total visited tiles: %d\n", uint(visitedTiles))
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
