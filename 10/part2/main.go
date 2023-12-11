package main

import (
	"container/list"
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

		pathCluster[nextStep.toTile] = true
		previousTile = currentTile
		currentTile = nextStep.toTile
		step++
	}

	fmt.Printf("Path cluster contains %d tiles: %v\n", len(pathCluster), printCluster(pathCluster))

	clusters := map[uint]map[Tile]bool{
		clusterPath: pathCluster,
	}

	for rowIdx, line := range lines {
		for colIdx := range line {
			t := Tile{
				rowIdx: rowIdx,
				colIdx: colIdx,
			}

			if !inAnyCluster(clusters, t) {
				newCluster := buildCluster(lines, clusters, t)
				newClusterID := uint(len(clusters))
				clusters[newClusterID] = newCluster

				fmt.Printf("Cluster #%d is detected with %d tiles: %v\n", newClusterID,
					len(newCluster), printCluster(newCluster))
			}
		}
	}

	lastRowIdx := len(lines) - 1
	lastColIdx := len(lines[0]) - 1
	clustersLen := uint(len(clusters))
	var enclosedClusters []map[Tile]bool
	var enclosedClustersIDs []uint

	for i := uint(1); i < clustersLen; i++ {
		cluster := clusters[i]
		isClusterAdjacentToBorder := false
		
		for tile := range cluster {
			if tile.rowIdx == 0 || tile.rowIdx == lastRowIdx ||
				tile.colIdx == 0 || tile.colIdx == lastColIdx {
				isClusterAdjacentToBorder = true
				break
			}
		}

		if !isClusterAdjacentToBorder {
			enclosedClusters = append(enclosedClusters, cluster)
			enclosedClustersIDs = append(enclosedClustersIDs, i)
		}
	}

	fmt.Printf("%d enclosed clusters found out of %d: %v\n", len(enclosedClusters), clustersLen-1,
		enclosedClustersIDs)

	enclosedTileCount := uint(0)
	for _, cluster := range enclosedClusters {
		enclosedTileCount += uint(len(cluster))
	}

	fmt.Println("Enclosed tiles found:", enclosedTileCount)
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

func inAnyCluster(clusters map[uint]map[Tile]bool, tile Tile) bool {
	for _, cluster := range clusters {
		if cluster[tile] {
			return true
		}
	}

	return false
}

func buildCluster(lines []string, clusters map[uint]map[Tile]bool, startTile Tile) map[Tile]bool {
	cluster := map[Tile]bool{
		startTile: true,
	}

	tilesToProcess := list.New()
	tilesToProcess.PushBack(startTile)

	for tilesToProcess.Len() > 0 {
		listEl := tilesToProcess.Front()
		currentTile := listEl.Value.(Tile)

		for _, dir := range allDirections {
			potentialClusterTile := getTileInDirection(currentTile, dir)
			if !isTileValid(lines, potentialClusterTile) {
				continue
			}

			if !cluster[potentialClusterTile] && !inAnyCluster(clusters, potentialClusterTile) {
				cluster[potentialClusterTile] = true
				tilesToProcess.PushBack(potentialClusterTile)
			}
		}

		tilesToProcess.Remove(listEl)
	}

	return cluster
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

	var tilesFormatted []string
	for _, tile := range tilesSl {
		tilesFormatted = append(tilesFormatted, fmt.Sprintf("%d:%d", tile.rowIdx+1, tile.colIdx+1))
	}

	return strings.Join(tilesFormatted, ", ")
}
