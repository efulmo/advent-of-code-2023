package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	var pointsTotal uint
	for lineIdx, line := range lines {
		lineParts := strings.Split(line, ":")
		numbersStrTrimmed := strings.TrimSpace(lineParts[1])
		numbers := strings.Split(numbersStrTrimmed, " | ")
		luckyNumbers := strings.Split(numbers[0], " ")

		luckyNumbersMap := make(map[string]bool)
		for _, num := range luckyNumbers {
			numTrimmed := strings.TrimSpace(num)
			if len(numTrimmed) == 0 {
				continue
			}

			luckyNumbersMap[numTrimmed] = true
		}

		var guessedNumbersCount uint
		cardNumbers := strings.Split(numbers[1], " ")
		for _, num := range cardNumbers {
			numTrimmed := strings.TrimSpace(num)
			if len(numTrimmed) == 0 {
				continue
			}

			_, exists := luckyNumbersMap[numTrimmed]
			if exists {
				guessedNumbersCount++
			}
		}

		if guessedNumbersCount > 0 {
			cardValue := uint(math.Pow(2, float64(guessedNumbersCount-1)))

			fmt.Printf("%d.%v. Lucky found: %d. Card value: %d\n", lineIdx+1,
				numbersStrTrimmed, guessedNumbersCount, cardValue)

			pointsTotal += uint(cardValue)
		} else {
			fmt.Printf("%d.%v. No numbers guessed\n", lineIdx+1, numbersStrTrimmed)
		}
	}

	fmt.Println("Points total:", pointsTotal)
}
