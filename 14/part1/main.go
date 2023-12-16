package main

import (
	"fmt"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	roundRock = byte('O')
	cubeRock  = byte('#')
	space     = byte('.')
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	var platform [][]byte
	for _, line := range lines {
		byteRow := make([]byte, 0, len(line))
		for _, r := range line {
			byteRow = append(byteRow, byte(r))
		}
		platform = append(platform, byteRow)
	}

	fmt.Println("Initial platform:")
	printBytes(platform)

	tiltNorth(platform)

	fmt.Println("Tilted platform:")
	printBytes(platform)

	fmt.Println("Load:", calculateNorhtBeamLoad(platform))
}

func tiltNorth(bytes [][]byte) {
	colCount := uint(len(bytes[0]))
	rowCount := uint(len(bytes))
	for colIdx := uint(0); colIdx < colCount; colIdx++ {
		for rowIdx := uint(0); rowIdx < rowCount; rowIdx++ {
			if bytes[rowIdx][colIdx] == roundRock {
				northernSpaceRowIdx := -1
				for i := int(rowIdx); i >= 0; i-- {
					b := bytes[i][colIdx]
					if b == space {
						northernSpaceRowIdx = i
					} else if b == cubeRock {
						break
					}
				}

				if northernSpaceRowIdx != -1 {
					swapBytes(bytes, rowIdx, colIdx, uint(northernSpaceRowIdx), colIdx)
				}
			}
		}
	}
}

func swapBytes(bytes [][]byte, row1, col1, row2, col2 uint) {
	bytes[row1][col1], bytes[row2][col2] = bytes[row2][col2], bytes[row1][col1]
}

func printBytes(bytes [][]byte) {
	for _, row := range bytes {
		for _, b := range row {
			fmt.Printf("%c", b)
		}
		fmt.Println()
	}
}

func calculateNorhtBeamLoad(bytes [][]byte) uint {
	colCount := uint(len(bytes[0]))
	rowCount := uint(len(bytes))

	var load uint
	for colIdx := uint(0); colIdx < colCount; colIdx++ {
		for rowIdx := uint(0); rowIdx < rowCount; rowIdx++ {
			if bytes[rowIdx][colIdx] == roundRock {
				rockLoad := rowCount - rowIdx
				
				util.DebugLog("Rock %d:%d has load %d\n", rowIdx, colIdx, rockLoad)
				load += rockLoad
			}
		}
	}

	return load
}
