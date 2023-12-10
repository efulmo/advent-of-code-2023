package main

import (
	"fmt"

	"github.com/efulmo/advent-of-code-2023/util"
)

type Coord struct {
	rowIdx int
	colIdx int
}

type Tile struct {
	coord Coord
	char  string
}

type Step struct {
	direction uint8
	fromTile  Tile
	toTile    Tile
}

const (
	directionUp    = 1
	directionRight = 2
	directionDown  = 3
	directionLeft  = 4
)

const (
	charStart  = "S"
	charGround = "."

	charUpDown  = "|"
	charUpRight = "F"
	charUpLeft  = "7"

	charRightLeft = "-"
	charRightUp   = "J"
	charRightDown = charUpLeft

	charDownUp    = charUpDown
	charDownRight = "L"
	charDownLeft  = charRightUp

	charLeftRight = charRightLeft
	charLeftUp    = charDownRight
	charLeftDown  = charUpRight
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	var startCoord Coord
	for lineIdx, line := range lines {
		for colIdx, r := range line {
			if string(r) == charStart {
				startCoord = Coord{
					rowIdx: lineIdx,
					colIdx: colIdx,
				}
				break
			}
		}
	}

	var currentTile = getTileAt(lines, startCoord)
	pathLength := uint(1)
	var previousTile Tile

	for {
		nextStep := getNextStep(lines, currentTile, previousTile)

		// start is found again; the loop has closed
		if nextStep.toTile.char == charStart {
			break
		}

		fmt.Printf("%d. Taking a step in direction %d to tile %d:%d with rune %s\n", pathLength,
			nextStep.direction, nextStep.toTile.coord.rowIdx+1, nextStep.toTile.coord.colIdx+1,
			nextStep.toTile.char)

		previousTile = currentTile
		currentTile = nextStep.toTile
		pathLength++
	}

	fmt.Println("Path length:", pathLength)
	fmt.Println("Farthest tile:", pathLength/2)
}

func getNextStep(lines []string, currentTile Tile, previousTile Tile) Step {
	availableDirs := getAvailableDirectionsFromChar(currentTile.char)

	currentCoord := currentTile.coord
	for _, dir := range availableDirs {
		var nextCoord Coord
		switch dir {
		case directionUp:
			nextCoord = Coord{currentCoord.rowIdx - 1, currentCoord.colIdx}
		case directionRight:
			nextCoord = Coord{currentCoord.rowIdx, currentCoord.colIdx + 1}
		case directionDown:
			nextCoord = Coord{currentCoord.rowIdx + 1, currentCoord.colIdx}
		case directionLeft:
			nextCoord = Coord{currentCoord.rowIdx, currentCoord.colIdx - 1}
		}

		// don't go back
		if previousTile.coord == nextCoord {
			continue
		}

		nextTile := getTileAt(lines, nextCoord)

		// check if getting to that tile allowed from this direction
		if !isStepDestinationValid(dir, nextTile) {
			continue
		}

		return Step{
			direction: dir,
			fromTile:  previousTile,
			toTile:    nextTile,
		}
	}

	panic(fmt.Errorf("Failed to find next tile from %d:%d(%s)", currentCoord.rowIdx,
		currentCoord.rowIdx, currentTile.char))
}

func getAvailableDirectionsFromChar(c string) []uint8 {
	var dirs []uint8

	if c == charStart || c == charDownUp || c == charRightUp || c == charLeftUp {
		dirs = append(dirs, directionUp)
	}
	if c == charStart || c == charLeftRight || c == charUpRight || c == charDownRight {
		dirs = append(dirs, directionRight)
	}
	if c == charStart || c == charUpDown || c == charRightDown || c == charLeftDown {
		dirs = append(dirs, directionDown)
	}
	if c == charStart || c == charRightLeft || c == charUpLeft || c == charDownLeft {
		dirs = append(dirs, directionLeft)
	}

	dirsLen := uint(len(dirs))
	if dirsLen < 1 || dirsLen > 4 {
		panic(fmt.Errorf("Unexpected number of directions for char %s: %v", c, dirs))
	}

	return dirs
}

func getTileAt(lines []string, coord Coord) Tile {
	return Tile{
		coord: coord,
		char:  getCharAt(lines, coord),
	}
}

func getCharAt(lines []string, coord Coord) string {
	row, col := coord.rowIdx, coord.colIdx
	if row < 0 || col < 0 || row >= len(lines) {
		return charGround
	}

	line := lines[row]
	if col >= len(line) {
		return charGround
	}

	return string(line[col])
}

func isStepDestinationValid(direction uint8, to Tile) bool {
	r := to.char
	if r == charStart {
		return true
	}

	switch direction {
	case directionUp:
		return r == charUpDown || r == charUpRight || r == charUpLeft
	case directionRight:
		return r == charRightLeft || r == charRightUp || r == charRightDown
	case directionDown:
		return r == charDownUp || r == charDownRight || r == charDownLeft
	case directionLeft:
		return r == charLeftRight || r == charLeftUp || r == charLeftDown
	default:
		panic(fmt.Errorf("Unexpected direction %d", direction))
	}
}
