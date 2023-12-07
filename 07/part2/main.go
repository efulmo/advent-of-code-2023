package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	comboNotSet    = 0
	combo5ofAKind  = 1
	combo4ofAKind  = 2
	comboFullHouse = 3
	combo3ofAKind  = 4
	combo2Pair     = 5
	combo1Pair     = 6
	comboHighCard  = 7
)

const cardChars = "AKQT98765432J"
const jokerChar = "J"

var cardRanks map[rune]uint8

type Hand struct {
	cards string
	bid   uint
	combo uint8
}

func (h *Hand) check() {
	if len(h.cards) != 5 {
		panic(fmt.Errorf("Invalid number of cards: %d", len(h.cards)))
	}

	for _, card := range h.cards {
		if !strings.ContainsRune(cardChars, card) {
			panic(fmt.Errorf("Invalid card found: %c", card))
		}
	}
}

func (h *Hand) getCombination() uint8 {
	if h.combo != comboNotSet {
		return h.combo
	}

	countByCard := map[string]uint{}
	for _, card := range h.cards {
		countByCard[string(card)]++
	}

	jokersCount := countByCard[jokerChar]
	if jokersCount > 0 {
		var maxNonJokerCard string
		maxNonJokerCardCount := uint(0)
		for card, count := range countByCard {
			if card != jokerChar && count > maxNonJokerCardCount {
				maxNonJokerCardCount = count
				maxNonJokerCard = card
			}
		}

		if maxNonJokerCard != "" {
			delete(countByCard, jokerChar)
			countByCard[maxNonJokerCard] += jokersCount
		}
	}

	var cardCounts []uint
	for _, count := range countByCard {
		cardCounts = append(cardCounts, count)
	}
	cardCountsLen := len(cardCounts)

	var combo uint8 = comboHighCard
	switch cardCountsLen {
	case 1:
		combo = combo5ofAKind
	case 2:
		if slices.Contains(cardCounts, 4) {
			combo = combo4ofAKind
		} else {
			combo = comboFullHouse
		}
	case 3:
		if slices.Contains(cardCounts, 3) {
			combo = combo3ofAKind
		} else {
			combo = combo2Pair
		}
	case 4:
		combo = combo1Pair
	}

	h.combo = combo

	return combo
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	var hands []*Hand
	for lineIdx, line := range lines {
		parts := strings.Fields(line)
		hand := &Hand{
			cards: parts[0],
			bid:   util.ParseUintOrPanic(parts[1]),
		}
		hand.check()

		fmt.Printf("%d. Hand %v is read\n", lineIdx+1, hand)

		hands = append(hands, hand)
	}

	fmt.Println("Hands before sorting:")
	fmt.Println(ptrsToHands(hands))

	slices.SortFunc(hands, compareCards)

	fmt.Println("Hands after sorting:")
	fmt.Println(ptrsToHands(hands))

	handsLen := len(hands)
	var winnings uint
	for i := 0; i < handsLen; i++ {
		rank := handsLen - i
		winnings += uint(rank) * hands[i].bid
	}
	fmt.Println("Winnings:", winnings)
}

func compareCards(h1, h2 *Hand) int {
	h1Combo, h2Combo := h1.getCombination(), h2.getCombination()

	// fmt.Printf("Comparing %v and %v\n", *h1, *h2)

	if h1Combo != h2Combo {
		return int(h1Combo) - int(h2Combo)
	}

	ranks := getCardRanks()
	for idx, h1r := range h1.cards {
		rank1, rank2 := ranks[h1r], ranks[rune(h2.cards[idx])]
		if rank1 == rank2 {
			continue
		}
		return int(rank1) - int(rank2)
	}

	fmt.Printf("Equal hands detected! %v and %v\n", *h1, *h2)
	return 0
}

func ptrsToHands(ptrs []*Hand) []Hand {
	var hands []Hand
	for _, ptr := range ptrs {
		hands = append(hands, *ptr)
	}
	return hands
}

func getCardRanks() map[rune]uint8 {
	if cardRanks == nil {
		cardRanks = make(map[rune]uint8, len(cardChars))
		for idx, r := range cardChars {
			cardRanks[r] = uint8(idx)
		}
	}

	return cardRanks
}
