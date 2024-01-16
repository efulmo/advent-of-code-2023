package main

import (
	"errors"
	"fmt"
	"maps"
	"slices"
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

	lastRowIdx := uint8(rowsTotal - 1)
	endCoordColIdx := strings.Index(lines[lastRowIdx], charPath)
	if endCoordColIdx == -1 {
		panic(errors.New("Start corrd isn't found"))
	}
	endCoord := Coord{lastRowIdx, uint8(endCoordColIdx)}

	crossings := getCrossings(lines, startCoord)
	fmt.Println("Detected crossings:", formatCoordsMap(crossings))

	graph := buildGraph(lines, crossings, startCoord, endCoord)
	fmt.Println("Graph:")
	printGraph(graph)

	path, err := getLongestPathToEnd(graph, startCoord, endCoord, []Coord{})
	if err != nil {
		fmt.Println("Path to end isn't found:", err.Error())
	} else {
		fmt.Printf("Longest path contains %d nodes: %v\n", len(path), formatCoords(path))
		fmt.Println("Length:", computePathLength(graph, path))
	}
}

func getCrossings(lines []string, startCoord Coord) map[Coord]bool {
	crossings := map[Coord]bool{}
	coordsToVisit := []Coord{startCoord}
	visitedCoords := map[Coord]bool{}
	emptyCoordsSet := map[Coord]bool{}

	for len(coordsToVisit) > 0 {
		coord := coordsToVisit[0]
		coordsToVisit = coordsToVisit[1:]

		nextAvailableSteps := getNextValidSteps(lines, coord, emptyCoordsSet)
		if len(nextAvailableSteps) > 2 {
			crossings[coord] = true
		}

		nextValidSteps := getNextValidSteps(lines, coord, visitedCoords)
		coordsToVisit = append(coordsToVisit, nextValidSteps...)
		visitedCoords[coord] = true
	}

	return crossings
}

func buildGraph(lines []string, crossings map[Coord]bool, startCoord, endCoord Coord) map[Coord]map[Coord]uint {
	graph := map[Coord]map[Coord]uint{}

	nodes := maps.Clone(crossings)
	nodes[startCoord] = true
	nodes[endCoord] = true

	for node := range nodes {
		graph[node] = map[Coord]uint{}

		stepsFromCrossing := getNextValidSteps(lines, node, map[Coord]bool{})

		for _, step := range stepsFromCrossing {
			stepsToClosestCrossing := uint(1)
			currentCoord := step
			visitedCoords := map[Coord]bool{node: true}
			nextSteps := getNextValidSteps(lines, currentCoord, visitedCoords)

			for len(nextSteps) == 1 {
				stepsToClosestCrossing++

				visitedCoords[currentCoord] = true

				currentCoord = nextSteps[0]
				nextSteps = getNextValidSteps(lines, currentCoord, visitedCoords)
			}

			if len(nextSteps) > 1 || currentCoord == startCoord || currentCoord == endCoord {
				graph[node][currentCoord] = stepsToClosestCrossing
			} else {
				fmt.Printf("Reached dead for node %d:%d end at %d:%d. Next steps: %d\n", node.rowIdx+1,
					node.colIdx+1, currentCoord.rowIdx+1, currentCoord.colIdx+1, len(nextSteps))
			}
		}
	}

	return graph
}

func formatCoordsMap(coords map[Coord]bool) []string {
	coordsSl := make([]Coord, 0, len(coords))
	for coord := range coords {
		coordsSl = append(coordsSl, coord)
	}
	sortCoords(coordsSl)

	return formatCoords(coordsSl)
}

func formatCoords(coords []Coord) []string {
	formattedCoords := make([]string, 0, len(coords))
	for _, coord := range coords {
		formattedCoords = append(formattedCoords, fmt.Sprintf("%d:%d", coord.rowIdx+1, coord.colIdx+1))
	}

	return formattedCoords
}

func sortCoords(coords []Coord) {
	slices.SortFunc(coords, func(c1, c2 Coord) int {
		if c1.rowIdx != c2.rowIdx {
			return int(c1.rowIdx) - int(c2.rowIdx)
		}
		return int(c1.colIdx) - int(c2.colIdx)
	})
}

func printGraph(graph map[Coord]map[Coord]uint) {
	fromNodes := make([]Coord, 0, len(graph))
	for node := range graph {
		fromNodes = append(fromNodes, node)
	}
	sortCoords(fromNodes)

	for _, fromNode := range fromNodes {
		toNodes := graph[fromNode]
		toNodeStrs := make([]string, 0, len(toNodes))
		for toNode, distance := range toNodes {
			toNodeStrs = append(toNodeStrs, fmt.Sprintf("[%d]%d:%d->%d:%d", distance, fromNode.rowIdx+1,
				fromNode.colIdx+1, toNode.rowIdx+1, toNode.colIdx+1))
		}
		fmt.Printf("%d:%d: %s\n", fromNode.rowIdx+1, fromNode.colIdx+1, strings.Join(toNodeStrs, ", "))
	}
}

func getNextValidSteps(lines []string, coord Coord, visitedCoords map[Coord]bool) []Coord {
	currentCoordChar := lines[coord.rowIdx][coord.colIdx : coord.colIdx+1]

	switch {
	case charForest == currentCoordChar:
		panic(fmt.Errorf("I am in the forest at %d:%d", coord.rowIdx+1, coord.colIdx+1))
	case charPath == currentCoordChar, slopeChars[currentCoordChar]:
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

func getLongestPathToEnd(
	graph map[Coord]map[Coord]uint,
	currentCoord, endCoord Coord,
	path []Coord,
) ([]Coord, error) {
	var longestPath []Coord
	var longestPathLegth uint
	for nextCoord := range graph[currentCoord] {
		
		// recursion exit condition
		if nextCoord == endCoord {
			finishedPath := slices.Clone(path)
			finishedPath = append(finishedPath, currentCoord, endCoord)

			pathLength := computePathLength(graph, finishedPath)
			if pathLength > longestPathLegth {
				longestPath = finishedPath
				longestPathLegth = pathLength
			}

		// general case
		} else if !slices.Contains(path, nextCoord) {
			furtherPath := slices.Clone(path)
			furtherPath = append(furtherPath, currentCoord)
			finishedPath, err := getLongestPathToEnd(graph, nextCoord, endCoord, furtherPath)
			if err == nil {
				pathLength := computePathLength(graph, finishedPath)
				if pathLength > longestPathLegth {
					longestPath = finishedPath
					longestPathLegth = pathLength
				}
			}
		}
	}

	if longestPath == nil {
		return nil, fmt.Errorf("No path to end found from %d:%d", currentCoord.rowIdx+1, currentCoord.colIdx+1)
	}

	return longestPath, nil
}

func computePathLength(graph map[Coord]map[Coord]uint, path []Coord) uint {
	if len(path) == 0 {
		panic(errors.New("Empty path is passed for length calculation"))
	}

	var length uint
	for i := uint(0); i < uint(len(path)-1); i++ {
		currentCoord, nextCoord := path[i], path[i+1]
		length += graph[currentCoord][nextCoord]
	}

	return length
}
