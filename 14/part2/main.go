package main

import (
	"crypto/sha256"
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

	platformByHash := make(map[string][][]byte)
	nextPlatformHashByHash := make(map[string]string)
	var previousPlatformHash string
	var platformHash string

	for i := uint(0); i < 1000000000; i++ {
		if platformHash == "" {
			platformHash = computePlatformHash(platform)
		}

		nextPlatformHash, found := nextPlatformHashByHash[platformHash]

		if (i+1)%1000000 == 0 {
			util.DebugLog("%dM. Next platform found - %t. Cache size: %d\n", (i+1)/1000000, found,
				len(nextPlatformHashByHash))
		}

		if found {
			// get next platform from cache
			platform = platformByHash[nextPlatformHash]

			// save connection with previous platform to cache
			nextPlatformHashByHash[previousPlatformHash] = platformHash
			previousPlatformHash = platformHash
			platformHash = nextPlatformHash

			continue
		}

		nextPlatformHashByHash[previousPlatformHash] = platformHash
		platformByHash[platformHash] = deepCopy(platform)
		previousPlatformHash = platformHash
		platformHash = ""

		doTitlCycle(platform)
	}

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

func tiltSouth(bytes [][]byte) {
	colCount := uint(len(bytes[0]))
	rowCount := uint(len(bytes))
	for colIdx := uint(0); colIdx < colCount; colIdx++ {
		for rowIdx := int(rowCount - 1); rowIdx >= 0; rowIdx-- {
			if bytes[rowIdx][colIdx] == roundRock {
				southernSpaceRowIdx := -1
				for i := rowIdx; i < int(rowCount); i++ {
					b := bytes[i][colIdx]
					if b == space {
						southernSpaceRowIdx = i
					} else if b == cubeRock {
						break
					}
				}

				if southernSpaceRowIdx != -1 {
					swapBytes(bytes, uint(rowIdx), colIdx, uint(southernSpaceRowIdx), colIdx)
				}
			}
		}
	}
}

func tiltEast(bytes [][]byte) {
	colCount := uint(len(bytes[0]))
	rowCount := uint(len(bytes))
	for rowIdx := uint(0); rowIdx < rowCount; rowIdx++ {
		for colIdx := int(colCount - 1); colIdx >= 0; colIdx-- {
			if bytes[rowIdx][colIdx] == roundRock {
				easternSpaceColIdx := -1
				for i := uint(colIdx); i < colCount; i++ {
					b := bytes[rowIdx][i]
					if b == space {
						easternSpaceColIdx = int(i)
					} else if b == cubeRock {
						break
					}
				}

				if easternSpaceColIdx != -1 {
					swapBytes(bytes, rowIdx, uint(colIdx), rowIdx, uint(easternSpaceColIdx))
				}
			}
		}
	}
}

func tiltWest(bytes [][]byte) {
	colCount := uint(len(bytes[0]))
	rowCount := uint(len(bytes))
	for rowIdx := uint(0); rowIdx < rowCount; rowIdx++ {
		for colIdx := uint(0); colIdx < colCount; colIdx++ {
			if bytes[rowIdx][colIdx] == roundRock {
				westernSpaceColIdx := -1
				for i := int(colIdx); i >= 0; i-- {
					b := bytes[rowIdx][i]
					if b == space {
						westernSpaceColIdx = i
					} else if b == cubeRock {
						break
					}
				}

				if westernSpaceColIdx != -1 {
					swapBytes(bytes, rowIdx, colIdx, rowIdx, uint(westernSpaceColIdx))
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

func doTitlCycle(bytes [][]byte) {
	tiltNorth(bytes)
	tiltWest(bytes)
	tiltSouth(bytes)
	tiltEast(bytes)
}

func calculateNorhtBeamLoad(bytes [][]byte) uint {
	colCount := uint(len(bytes[0]))
	rowCount := uint(len(bytes))

	var load uint
	for colIdx := uint(0); colIdx < colCount; colIdx++ {
		for rowIdx := uint(0); rowIdx < rowCount; rowIdx++ {
			if bytes[rowIdx][colIdx] == roundRock {
				rockLoad := rowCount - rowIdx

				// util.DebugLog("Rock %d:%d has load %d\n", rowIdx, colIdx, rockLoad)
				load += rockLoad
			}
		}
	}

	return load
}

func computePlatformHash(bytes [][]byte) string {
	hasher := sha256.New()
	for _, row := range bytes {
		hasher.Write(row)
	}
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func deepCopy(platform [][]byte) [][]byte {
	platformCopy := make([][]byte, 0, len(platform))
	for _, row := range platform {
		rowCopy := make([]byte, len(row))
		copy(rowCopy, row)

		platformCopy = append(platformCopy, rowCopy)
	}

	return platformCopy
}
