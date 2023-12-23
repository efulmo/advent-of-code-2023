package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	runeEmpty              = '.'
	runeMirrorForward      = '/'
	runeMirrorBackward     = '\\'
	runeSplitterHorizontal = '-'
	runeSplitterVertical   = '|'

	runeTileEnergized = '#'
	runeTileRegular   = '.'

	directionRight = 1
	directionDown  = 2
	directionLeft  = 3
	directionUp    = 4
)

type Coord struct {
	rowIdx, colIdx int
}

type StartPosition struct {
	coord     Coord
	direction uint8
}

func main() {
	contraption, err := util.ReadInputFile()
	util.PanicOnError(err)

	rowsTotal := uint(len(contraption))
	columnsTotal := uint(len(contraption[0]))
	fmt.Printf("A contraption %dx%d is read\n", rowsTotal, columnsTotal)

	var startPositions []StartPosition
	for colIdx := uint(0); colIdx < columnsTotal; colIdx++ {
		startPositions = append(startPositions, StartPosition{
			coord: Coord{
				rowIdx: 0,
				colIdx: int(colIdx),
			},
			direction: directionDown,
		})
		startPositions = append(startPositions, StartPosition{
			coord: Coord{
				rowIdx: int(rowsTotal - 1),
				colIdx: int(colIdx),
			},
			direction: directionUp,
		})
	}
	for rowIdx := uint(0); rowIdx < rowsTotal; rowIdx++ {
		startPositions = append(startPositions, StartPosition{
			coord: Coord{
				rowIdx: int(rowIdx),
				colIdx: 0,
			},
			direction: directionRight,
		})
		startPositions = append(startPositions, StartPosition{
			coord: Coord{
				rowIdx: int(rowIdx),
				colIdx: int(columnsTotal - 1),
			},
			direction: directionLeft,
		})
	}

	positionsCount := uint(len(startPositions))
	var maxVisitedTiles uint
	for i, position := range startPositions {
		visitedTiles := make(map[Coord][]uint8)

		simulateBeam(contraption, visitedTiles, position.coord, position.direction)
		
		visitedTilesCount := uint(len(visitedTiles))
		maxVisitedTiles = max(maxVisitedTiles, visitedTilesCount)

		fmt.Printf("%d/%d: %d\n", i+1, positionsCount, maxVisitedTiles)
	}

	fmt.Println("Max visited tiles:", maxVisitedTiles)
}

func simulateBeam(
	contraction []string,
	visitedTiles map[Coord][]uint8,
	startCoord Coord,
	startDirection uint8,
) {
	coord := startCoord
	direction := startDirection

simulationLoop:
	for {
		rowIdx, colIdx := coord.rowIdx, coord.colIdx

		visitedFromDirections := visitedTiles[coord]
		if slices.Contains(visitedFromDirections, direction) {
			util.DebugLog("Tile %d:%d was already visited from direction %d. Stopping beam simulation\n",
				rowIdx+1, colIdx+1, direction)
			break simulationLoop
		} else {
			visitedTiles[coord] = append(visitedFromDirections, direction)
		}

		run := rune(contraction[rowIdx][colIdx])

		switch run {
		case runeEmpty:
			nextCoord := getNextCoord(contraction, coord, direction)
			if nextCoord == nil {
				break simulationLoop
			}
			coord = *nextCoord
		case runeMirrorForward, runeMirrorBackward:
			direction = reflectBeam(run, direction)
			nextCoord := getNextCoord(contraction, coord, direction)
			if nextCoord == nil {
				break simulationLoop
			}
			coord = *nextCoord
		case runeSplitterHorizontal:
			if direction == directionRight || direction == directionLeft {
				nextCoord := getNextCoord(contraction, coord, direction)
				if nextCoord == nil {
					break simulationLoop
				}
				coord = *nextCoord
			} else {
				nextCoord := getNextCoord(contraction, coord, directionRight)
				if nextCoord != nil {
					simulateBeam(contraction, visitedTiles, *nextCoord, directionRight)
				}

				nextCoord = getNextCoord(contraction, coord, directionLeft)
				if nextCoord != nil {
					simulateBeam(contraction, visitedTiles, *nextCoord, directionLeft)
				}

				break simulationLoop
			}
		case runeSplitterVertical:
			if direction == directionDown || direction == directionUp {
				nextCoord := getNextCoord(contraction, coord, direction)
				if nextCoord == nil {
					break simulationLoop
				}
				coord = *nextCoord
			} else {
				nextCoord := getNextCoord(contraction, coord, directionDown)
				if nextCoord != nil {
					simulateBeam(contraction, visitedTiles, *nextCoord, directionDown)
				}

				nextCoord = getNextCoord(contraction, coord, directionUp)
				if nextCoord != nil {
					simulateBeam(contraction, visitedTiles, *nextCoord, directionUp)
				}

				break simulationLoop
			}
		}
	}
}

func getNextCoord(contraction []string, coord Coord, direction uint8) *Coord {
	var nextCoord Coord

	switch direction {
	case directionRight:
		nextCoord = Coord{
			rowIdx: coord.rowIdx,
			colIdx: coord.colIdx + 1,
		}
	case directionDown:
		nextCoord = Coord{
			rowIdx: coord.rowIdx + 1,
			colIdx: coord.colIdx,
		}
	case directionLeft:
		nextCoord = Coord{
			rowIdx: coord.rowIdx,
			colIdx: coord.colIdx - 1,
		}
	case directionUp:
		nextCoord = Coord{
			rowIdx: coord.rowIdx - 1,
			colIdx: coord.colIdx,
		}
	}

	if nextCoord.rowIdx >= 0 && nextCoord.rowIdx < len(contraction) {
		nextRow := contraction[nextCoord.rowIdx]
		if nextCoord.colIdx >= 0 && nextCoord.colIdx < len(nextRow) {
			util.DebugLog("Next coord %d:%d is found from %d:%d in %d direction\n",
				nextCoord.rowIdx+1, nextCoord.colIdx+1, coord.rowIdx+1, coord.colIdx+1, direction)
			return &nextCoord
		}
	}

	return nil
}

func reflectBeam(mirror rune, direction uint8) uint8 {
	switch mirror {
	case runeMirrorForward:
		switch direction {
		case directionRight:
			return directionUp
		case directionDown:
			return directionLeft
		case directionLeft:
			return directionDown
		case directionUp:
			return directionRight
		}
	case runeMirrorBackward:
		switch direction {
		case directionRight:
			return directionDown
		case directionDown:
			return directionRight
		case directionLeft:
			return directionUp
		case directionUp:
			return directionLeft
		}
	}

	panic(fmt.Errorf("Unexpected mirror %c or direction %d", mirror, direction))
}

func formatVisitedTiles(rowsTotal, columnsTotal uint, tiles map[Coord][]uint8) string {
	var rows []string
	for rowIdx := uint(0); rowIdx < rowsTotal; rowIdx++ {
		var rowBld strings.Builder
		for colIdx := uint(0); colIdx < columnsTotal; colIdx++ {
			_, wasTileVisited := tiles[Coord{int(rowIdx), int(colIdx)}]
			if wasTileVisited {
				rowBld.WriteRune(runeTileEnergized)
			} else {
				rowBld.WriteRune(runeTileRegular)
			}
		}
		rows = append(rows, fmt.Sprintf("%3d. %v", rowIdx+1, rowBld.String()))
	}

	return strings.Join(rows, "\n")
}
