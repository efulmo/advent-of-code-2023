package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var cubeLimits = map[string]uint8{
	"red": 12,
	"green": 13,
	"blue": 14,
}

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

	var possibleGameSum uint
	for lineIdx, line := range lines {
		lineParts := strings.Split(line, ":")
		gameParts := strings.Split(lineParts[0], " ")
		gameIdStr := gameParts[1]
		gameId, err := strconv.Atoi(gameIdStr)
		if err != nil {
			fmt.Printf("Line %d: Cannot parse game ID from <%s>: %s\n", lineIdx, gameIdStr, err.Error())
			os.Exit(1)
		}

		isGamePossible := true

		rounds := strings.Split(strings.TrimSpace(lineParts[1]), ";")
		roundsLoop:
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
				cubeLimit, ok := cubeLimits[cubeColour]
				if !ok {
					fmt.Printf("Line %d, round %d, cubeset %d: Unknown cube colour %s\n", 
						lineIdx, roundIdx, cubeSetIdx, cubeColour)
					os.Exit(1)
				}

				if cubeCount > int(cubeLimit) {
					isGamePossible = false
					break roundsLoop
				}
			}
		}

		if isGamePossible {
			fmt.Printf("Line %d. Game %d is possible\n", lineIdx, gameId)
			possibleGameSum += uint(gameId)
		}
	}

	fmt.Println("Sum of possible games:", possibleGameSum)
}
