package main

import (
	"errors"
	"fmt"
	"maps"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	charPath       = "."
	charForest     = "#"
	charSlopeRight = ">"
	charSlopeLeft  = "<"
	charSlopeDown  = "v"
	charSlopeUp    = "^"
)

type Coord struct {
	rowIdx, colIdx uint8
}

type DirectionDiff struct {
	rowiDff, colDiff int8
}

var (
	directionDown  = DirectionDiff{1, 0}
	directionUp    = DirectionDiff{-1, 0}
	directionRight = DirectionDiff{0, 1}
	directionLeft  = DirectionDiff{0, -1}
	allDirections  = []DirectionDiff{directionUp, directionRight, directionDown, directionLeft}

	slopeChars = map[string]bool{
		charSlopeUp:    true,
		charSlopeDown:  true,
		charSlopeRight: true,
		charSlopeLeft:  true,
	}
	directionBySlopeChar = map[string]DirectionDiff{
		charSlopeUp:    directionUp,
		charSlopeDown:  directionDown,
		charSlopeRight: directionRight,
		charSlopeLeft:  directionLeft,
	}
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	rowsTotal := len(lines)

	if rowsTotal < 2 {
		panic(fmt.Errorf("Input has too little lines: %d", rowsTotal))
	}

	startCoordColIdx := strings.Index(lines[0], charPath)
	if startCoordColIdx == -1 {
		panic(errors.New("Start corrd isn't found"))
	}
	startCoord := Coord{0, uint8(startCoordColIdx)}

	lastRowIdx := uint8(rowsTotal-1)
	endCoordColIdx := strings.Index(lines[lastRowIdx], charPath)
	if endCoordColIdx == -1 {
		panic(errors.New("Start corrd isn't found"))
	}
	endCoord := Coord{lastRowIdx, uint8(endCoordColIdx)}

	path, err := getLongestPathToEnd(lines, startCoord, endCoord, map[Coord]bool{})
	if err != nil {
		fmt.Println("Path to end isn't found:", err.Error())
	} else {
		fmt.Printf("Longest path is %d steps long\n", len(path)-1)
	}
}

func getLongestPathToEnd(
	lines []string,
	startCoord, endCoord Coord,
	visitedCoords map[Coord]bool,
) (map[Coord]bool, error) {
	currentCoord := startCoord
	nextCoords := getNextValidSteps(lines, currentCoord, visitedCoords)

	// no crossing; keep walking till there are available steps
	for len(nextCoords) == 1 {
		nextCoord := nextCoords[0]
		// fmt.Printf("Step %d:%d -> %d:%d\n", currentCoord.rowIdx+1, currentCoord.colIdx+1,
		// 	nextCoord.rowIdx+1, nextCoord.colIdx+1)

		visitedCoords[currentCoord] = true
		currentCoord = nextCoord

		nextCoords = getNextValidSteps(lines, currentCoord, visitedCoords)
	}

	// either expected path end or dead end
	if len(nextCoords) == 0 {
		if currentCoord == endCoord {
			visitedCoords[currentCoord] = true
			return visitedCoords, nil
		}
		return nil, fmt.Errorf("Dead end at %d:%d", currentCoord.rowIdx+1, currentCoord.colIdx+1)
	}

	// crossing met
	var longestPath map[Coord]bool
	for _, nextStep := range nextCoords {
		visitedCoordsCopy := maps.Clone(visitedCoords)
		path, err := getLongestPathToEnd(lines, nextStep, endCoord, visitedCoordsCopy)
		if err == nil {
			// fmt.Printf("End coord reached. Path length - %d, longest so far - %d\n", len(path), 
			// 	len(longestPath))

			if len(path) > len(longestPath) {
				longestPath = path
			}
		}
	}

	if longestPath != nil {
		return longestPath, nil
	}
	return nil, fmt.Errorf("No ways from crossing %d:%d", currentCoord.rowIdx+1, currentCoord.colIdx+1)
}

func getNextValidSteps(lines []string, coord Coord, visitedCoords map[Coord]bool) []Coord {
	currentCoordChar := lines[coord.rowIdx][coord.colIdx : coord.colIdx+1]

	switch {
	case charForest == currentCoordChar:
		panic(fmt.Errorf("I am in the forest at %d:%d", coord.rowIdx+1, coord.colIdx+1))
	case slopeChars[currentCoordChar]:
		direction := directionBySlopeChar[currentCoordChar]
		nextCoord, isValid := getCoordIfValid(lines, int16(coord.rowIdx)+int16(direction.rowiDff),
			int16(coord.colIdx)+int16(direction.colDiff), visitedCoords)
		if !isValid {
			panic(fmt.Errorf("No valid next step from %d:%d", coord.rowIdx+1, coord.colIdx+1))
			// return []Coord{}
		}

		return []Coord{nextCoord}
	case charPath == currentCoordChar:
		var nextValidSteps []Coord
		for _, direction := range allDirections {
			nextCoord, isValid := getCoordIfValid(lines, int16(coord.rowIdx)+int16(direction.rowiDff),
				int16(coord.colIdx)+int16(direction.colDiff), visitedCoords)
			if isValid {
				nextValidSteps = append(nextValidSteps, nextCoord)
			}
		}
		return nextValidSteps
	default:
		panic(fmt.Errorf("Unexpected char at %d:%d", coord.rowIdx+1, coord.colIdx+1))
	}
}

func getCoordIfValid(
	lines []string,
	rowIdx, colIdx int16,
	visitedCoords map[Coord]bool,
) (coord Coord, isValid bool) {
	if rowIdx < 0 || colIdx < 0 {
		return Coord{}, false
	}
	if rowIdx >= int16(len(lines)) {
		return Coord{}, false
	}
	row := lines[rowIdx]
	if colIdx >= int16(len(row)) {
		return Coord{}, false
	}
	char := lines[rowIdx][colIdx : colIdx+1]
	if char == charForest {
		return Coord{}, false
	}
	coord = Coord{uint8(rowIdx), uint8(colIdx)}

	return coord, !visitedCoords[coord]
}
