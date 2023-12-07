package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <input-file-path>\n", os.Args[0])
		return
	}

	inputeFilePath := os.Args[1]
	fmt.Printf("Reading input from file <%s>\n", inputeFilePath)

	data, err := os.ReadFile(inputeFilePath)
	if err != nil {
		fmt.Printf("Error reading file <%s>: %s\n", inputeFilePath, err.Error())
		os.Exit(1)
	}
	fmt.Printf("Read %d bytes\n", len(data))

	lines := strings.Split(string(data), "\n")
	fmt.Printf("Read %d lines\n", len(lines))

	var gamePowerSum uint
	for lineIdx, line := range lines {
		lineParts := strings.Split(line, ":")

		optimalCubeCount := map[string]uint{}

		rounds := strings.Split(strings.TrimSpace(lineParts[1]), ";")
		for roundIdx, round := range rounds {
			cubeSets := strings.Split(round, ", ")
			for cubeSetIdx, cubeSet := range cubeSets {
				setParts := strings.Split(strings.TrimSpace(cubeSet), " ")
				cubeCountStr := setParts[0]
				cubeCount, err := strconv.Atoi(cubeCountStr)
				if err != nil {
					fmt.Printf("Line %d, round %d, cubeset %d: Cannot parse cube count from <%s>: %s\n",
						lineIdx, roundIdx, cubeSetIdx, cubeCountStr, err.Error())
					os.Exit(1)
				}

				cubeColour := setParts[1]

				lowestCubeCount, ok := optimalCubeCount[cubeColour]
				if !ok || lowestCubeCount < uint(cubeCount) {
					optimalCubeCount[cubeColour] = uint(cubeCount)
				}
			}
		}

		gamePower := 1
		for _, cubeCount := range optimalCubeCount {
			gamePower *= int(cubeCount)
		}

		fmt.Printf("%d. %s\nOptimal count: %v. Game power: %d\n",
			lineIdx, line, optimalCubeCount, gamePower)
		gamePowerSum += uint(gamePower)
	}

	fmt.Println("Sum of possible games:", gamePowerSum)
}
