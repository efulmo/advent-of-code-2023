package main

import (
	"fmt"
	"slices"
	"strings"

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

var allDirections = []uint8{
	directionUp,
	directionRight,
	directionDown,
	directionLeft,
}

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

const clusterPath = 0

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
	var previousTile Tile
	step := uint(1)
	path := []Tile{startTile}
	pathCluster := map[Tile]bool{
		startTile: true,
	}

	for {
		nextStep := getNextStep(lines, currentTile, previousTile)

		// start is found again; the loop has closed
		nextTileChar := getCharAt(lines, nextStep.toTile)
		if nextTileChar == charStart {
			break
		}

		fmt.Printf("%d. Taking a step in direction %d to tile %d:%d with char %s\n", step,
			nextStep.direction, nextStep.toTile.rowIdx+1, nextStep.toTile.colIdx+1,
			nextTileChar)

		path = append(path, nextStep.toTile)
		pathCluster[nextStep.toTile] = true
		previousTile = currentTile
		currentTile = nextStep.toTile
		step++
	}

	fmt.Printf("Path cluster contains %d tiles: %v\n", len(pathCluster), printCluster(pathCluster))

	replacementChar := getStartTileReplacement(lines, path)
	lines[startTile.rowIdx] = strings.ReplaceAll(lines[startTile.rowIdx], charStart, replacementChar)

	var enclosedTiles []Tile

	for rowIdx, line := range lines {
		isWhithinPath := false
		for colIdx := range line {
			t := Tile{
				rowIdx: rowIdx,
				colIdx: colIdx,
			}
			c := getCharAt(lines, t)

			if pathCluster[t] {
				// north -> inverse
				if c == charDownUp || c == charLeftUp || c == charRightUp {
					isWhithinPath = !isWhithinPath
				}
			} else if isWhithinPath {
				enclosedTiles = append(enclosedTiles, t)
			}
		}
	}

	fmt.Printf("%d enclosed tiles found: %v\n", len(enclosedTiles), printTiles(enclosedTiles))
}

func getNextStep(lines []string, currentTile Tile, previousTile Tile) Step {
	availableDirs := getAvailableDirectionsFromChar(getCharAt(lines, currentTile))

	for _, dir := range availableDirs {
		nextTile := getTileInDirection(currentTile, dir)

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

func getTileInDirection(tile Tile, direction uint8) Tile {
	switch direction {
	case directionUp:
		return Tile{tile.rowIdx - 1, tile.colIdx}
	case directionRight:
		return Tile{tile.rowIdx, tile.colIdx + 1}
	case directionDown:
		return Tile{tile.rowIdx + 1, tile.colIdx}
	case directionLeft:
		return Tile{tile.rowIdx, tile.colIdx - 1}
	default:
		panic(fmt.Errorf("Invalid direction passed: %d", direction))
	}
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
	if isTileValid(lines, tile) {
		return string(lines[tile.rowIdx][tile.colIdx])
	}

	return charBeyondBorders
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

func isTileValid(lines []string, tile Tile) bool {
	row, col := tile.rowIdx, tile.colIdx
	if row < 0 || col < 0 || row >= len(lines) {
		return false
	}

	line := lines[row]
	if col >= len(line) {
		return false
	}

	return true
}

func printCluster(cluster map[Tile]bool) string {
	var tilesSl []Tile
	for tile := range cluster {
		tilesSl = append(tilesSl, tile)
	}
	slices.SortFunc(tilesSl, func(t1, t2 Tile) int {
		if t1.rowIdx != t2.rowIdx {
			return t1.rowIdx - t2.rowIdx
		}
		if t1.colIdx != t2.colIdx {
			return t1.colIdx - t2.colIdx
		}
		return 0
	})

	return printTiles(tilesSl)
}

func printTiles(tiles []Tile) string {
	var tilesFormatted []string
	for _, tile := range tiles {
		tilesFormatted = append(tilesFormatted, fmt.Sprintf("%d:%d", tile.rowIdx+1, tile.colIdx+1))
	}

	return strings.Join(tilesFormatted, ", ")
}

func getStartTileReplacement(lines []string, path []Tile) string {
	startTile := path[0]
	lastTile := path[len(path)-1]
	preLastTile := path[len(path)-2]

	stepFromLast := getNextStep(lines, lastTile, preLastTile)
	stepFromStart := getNextStep(lines, startTile, lastTile)

	// -
	if (stepFromLast.direction == directionRight || stepFromStart.direction == directionRight) &&
		(stepFromLast.direction == directionLeft || stepFromStart.direction == directionLeft) {
		return charLeftRight
	}
	// |
	if (stepFromLast.direction == directionUp || stepFromStart.direction == directionUp) &&
		(stepFromLast.direction == directionDown || stepFromStart.direction == directionDown) {
		return charDownUp
	}
	// F
	if (stepFromLast.direction == directionUp && stepFromStart.direction == directionRight) ||
		(stepFromLast.direction == directionLeft && stepFromStart.direction == directionDown) {
		return charUpRight
	}
	// 7
	if (stepFromLast.direction == directionUp && stepFromStart.direction == directionLeft) ||
		(stepFromLast.direction == directionRight && stepFromStart.direction == directionDown) {
		return charUpLeft
	}
	// L
	if (stepFromLast.direction == directionDown && stepFromStart.direction == directionRight) ||
		(stepFromLast.direction == directionLeft && stepFromStart.direction == directionUp) {
		return charDownRight
	}
	// J
	if (stepFromLast.direction == directionDown && stepFromStart.direction == directionLeft) ||
		(stepFromLast.direction == directionRight && stepFromStart.direction == directionUp) {
		return charDownLeft
	}

	panic(fmt.Errorf("Failed to detect start tile replacement. From last direction %d. " +
		"From start direction %d", stepFromLast.direction, stepFromStart.direction))
}
