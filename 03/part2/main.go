package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"

	"github.com/efulmo/advent-of-code-2023/util"
)

type PartID struct {
	lineIdx  uint
	startIdx uint // including
	endIdx   uint // excluding
	number   uint
}

type Point struct {
	rowIdx int
	colIdx int
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	gearConnections := make(map[Point][]PartID)
	for lineIdx, line := range lines {
		partIDs := parsePartIDs(line, uint(lineIdx))

		for _, partID := range partIDs {
			gearPoints := getGearPoints(lines, partID)

			if len(gearPoints) > 0 {
				fmt.Printf("%d. Part %d has %d gears around\n", lineIdx+1, partID.number, len(gearPoints))
			}

			for _, point := range gearPoints {
				connections := gearConnections[point]
				connections = append(connections, partID)

				gearConnections[point] = connections
			}
		}
	}

	fmt.Printf("%d gears detected around part IDs\n", len(gearConnections))

	var gearRatioSum uint
	for gearPoint, partIDs := range gearConnections {
		partsConnected := len(partIDs)
		if partsConnected == 2 {
			fmt.Printf("Gear %d:%d has %d and %d parts connected\n",
				gearPoint.rowIdx, gearPoint.colIdx,
				partIDs[0].number, partIDs[1].number)

			gearRatio := partIDs[0].number * partIDs[1].number
			gearRatioSum += gearRatio
		} else {
			fmt.Printf("Gear %d:%d has %d connections\n", gearPoint.rowIdx, gearPoint.colIdx, partsConnected)
		}
	}

	fmt.Println("Gear ratio sum:", gearRatioSum)
}

const NilIdx = -1

func parsePartIDs(line string, lineIdx uint) []PartID {
	var partIDs []PartID

	startRegionIdx := NilIdx
	for charIdx, r := range line {
		if unicode.IsDigit(r) {
			if startRegionIdx == NilIdx {
				startRegionIdx = charIdx
			}
		} else {
			if startRegionIdx != NilIdx {
				partIDs = append(partIDs, PartID{
					lineIdx:  lineIdx,
					startIdx: uint(startRegionIdx),
					endIdx:   uint(charIdx),
					number:   parsePartIDNumber(line, uint(startRegionIdx), uint(charIdx)),
				})

				// fmt.Printf("Parsed region %d:%d\n", startRegionIdx, charIdx)

				startRegionIdx = NilIdx
			}
		}
	}

	if startRegionIdx != NilIdx {
		partIDs = append(partIDs, PartID{
			lineIdx:  lineIdx,
			startIdx: uint(startRegionIdx),
			endIdx:   uint(len(line)),
			number:   parsePartIDNumber(line, uint(startRegionIdx), uint(len(line))),
		})
	}

	return partIDs
}

func parsePartIDNumber(line string, startIdx, endIdx uint) uint {
	partIDStr := line[startIdx:endIdx]

	// potential int truncation
	if len(partIDStr) > 6 {
		panic(fmt.Errorf("Potentionally too big part ID for int: %s", partIDStr))
	}

	return util.ParseUintOrPanic(partIDStr)
}

func partIDsToNumbers(partIDs []PartID) []uint {
	numbers := make([]uint, 0, len(partIDs))
	for _, partID := range partIDs {
		numbers = append(numbers, partID.number)
	}

	return numbers
}

const GearRune = '*'

func getGearPoints(lines []string, partID PartID) []Point {
	pointsAround := make([]Point, 0, 12)
	lineIdx := int(partID.lineIdx)
	startIdx := int(partID.startIdx)
	endIdx := int(partID.endIdx)

	// left & right
	pointsAround = append(pointsAround, Point{lineIdx, startIdx - 1})
	pointsAround = append(pointsAround, Point{lineIdx, endIdx})

	// top & bottom
	partIDLen := endIdx - startIdx
	for i := 0; i < partIDLen; i++ {
		pointsAround = append(pointsAround, Point{lineIdx - 1, startIdx + i})
		pointsAround = append(pointsAround, Point{lineIdx + 1, startIdx + i})
	}

	// top diagonals
	pointsAround = append(pointsAround, Point{lineIdx - 1, startIdx - 1})
	pointsAround = append(pointsAround, Point{lineIdx - 1, endIdx})

	// bottom diagonals
	pointsAround = append(pointsAround, Point{lineIdx + 1, startIdx - 1})
	pointsAround = append(pointsAround, Point{lineIdx + 1, endIdx})

	var gearPoints []Point
	for _, point := range pointsAround {
		r := getRuneAtOrDot(lines, point.rowIdx, point.colIdx)
		if r == GearRune {
			gearPoints = append(gearPoints, point)
		}
	}

	return gearPoints
}

const Dot = '.'

func getRuneAtOrDot(lines []string, rowIdx, colIdx int) rune {
	// negative indexes
	if rowIdx < 0 || colIdx < 0 {
		return Dot
	}

	// too big rowIdx
	if rowIdx >= len(lines) {
		return Dot
	}

	// too big colIdx
	line := lines[rowIdx]
	if colIdx >= len(line) {
		return Dot
	}

	return rune(line[colIdx])
}
