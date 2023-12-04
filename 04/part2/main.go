package main

import (
	"container/list"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	guessCountByCardId := make(map[uint8]uint8, len(lines))
	cardsQueue := list.New()

	for lineIdx, line := range lines {
		lineParts := strings.Split(line, ":")
		cardIdStr := strings.TrimPrefix(lineParts[0], "Card")
		cardIdInt, err := strconv.Atoi(strings.TrimSpace(cardIdStr))
		if err != nil {
			util.PanicOnError(errors.Join(fmt.Errorf("Line %d. Error parsing <%s> as int", lineIdx+1,
			cardIdStr), err))
		}
		cardId := uint8(cardIdInt)

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

		var guessedNumbersCount uint8
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

		fmt.Printf("Card %d. Guessed numbers: %d\n", cardId, guessedNumbersCount)

		guessCountByCardId[cardId] = guessedNumbersCount
		cardsQueue.PushBack(cardId)
	}

	var cardsProcessed, maxQueueLen uint
	for cardsQueue.Len() > 0 {
		maxQueueLen = max(maxQueueLen, uint(cardsQueue.Len()))

		queueEl := cardsQueue.Front()
		cardId := queueEl.Value.(uint8)
		for i := uint8(1); i <= guessCountByCardId[cardId]; i++ {
			cardsQueue.PushFront(cardId + i)
		}

		cardsQueue.Remove(queueEl)
		cardsProcessed++
	}

	fmt.Println("Max queue len:", maxQueueLen)
	fmt.Println("Cards total:", cardsProcessed)
}
