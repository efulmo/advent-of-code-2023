package main

import (
	"fmt"

	"github.com/efulmo/advent-of-code-2023/util"
)

type Tile struct {
	rowIdx int
	colIdx int
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
	charStart         = "S"
	charGround        = "."
	charBeyondBorders = "?"

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

	var startTile Tile
	for lineIdx, line := range lines {
		for colIdx, r := range line {
			if string(r) == charStart {
				startTile = Tile{
					rowIdx: lineIdx,
					colIdx: colIdx,
				}
				break
			}
		}
	}

	var currentTile = startTile
	pathLength := uint(1)
	var previousTile Tile

	for {
		nextStep := getNextStep(lines, currentTile, previousTile)

		// start is found again; the loop has closed
		nextTileChar := getCharAt(lines, nextStep.toTile)
		if nextTileChar == charStart {
			break
		}

		fmt.Printf("%d. Taking a step in direction %d to tile %d:%d with char %s\n", pathLength,
			nextStep.direction, nextStep.toTile.rowIdx+1, nextStep.toTile.colIdx+1,
			nextTileChar)

		previousTile = currentTile
		currentTile = nextStep.toTile
		pathLength++
	}

	fmt.Println("Path length:", pathLength)
	fmt.Println("Farthest tile:", pathLength/2)
}

func getNextStep(lines []string, currentTile Tile, previousTile Tile) Step {
	availableDirs := getAvailableDirectionsFromChar(getCharAt(lines, currentTile))

	for _, dir := range availableDirs {
		var nextTile Tile
		switch dir {
		case directionUp:
			nextTile = Tile{currentTile.rowIdx - 1, currentTile.colIdx}
		case directionRight:
			nextTile = Tile{currentTile.rowIdx, currentTile.colIdx + 1}
		case directionDown:
			nextTile = Tile{currentTile.rowIdx + 1, currentTile.colIdx}
		case directionLeft:
			nextTile = Tile{currentTile.rowIdx, currentTile.colIdx - 1}
		}

		// don't go back
		if previousTile == nextTile {
			continue
		}

		// check if getting to that tile allowed from this direction
		if !isStepDestinationValid(lines, dir, nextTile) {
			continue
		}

		return Step{
			direction: dir,
			fromTile:  previousTile,
			toTile:    nextTile,
		}
	}

	panic(fmt.Errorf("Failed to find next tile from %d:%d(%s)", currentTile.rowIdx,
		currentTile.rowIdx, getCharAt(lines, currentTile)))
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

func getCharAt(lines []string, tile Tile) string {
	row, col := tile.rowIdx, tile.colIdx
	if row < 0 || col < 0 || row >= len(lines) {
		return charBeyondBorders
	}

	line := lines[row]
	if col >= len(line) {
		return charBeyondBorders
	}

	return string(line[col])
}

func isStepDestinationValid(lines []string, direction uint8, to Tile) bool {
	r := getCharAt(lines, to)
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
