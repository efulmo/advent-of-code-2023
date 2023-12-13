package main

import (
	"fmt"
	"math/bits"
	"slices"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	runeDamagedSpring     = '#'
	runeOperationalSpring = '.'
	runeUnknownSpring     = '?'
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	var damageVariantSum uint
	for lineIdx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			panic(fmt.Errorf("Unexpected fields count in the line - %d", len(fields)))
		}

		springMap := fields[0]
		damageCheckSum := fields[1]
		
		springMapUnfolded := springMap + string(runeUnknownSpring) +
			springMap + string(runeUnknownSpring) +
			springMap + string(runeUnknownSpring) +
			springMap + string(runeUnknownSpring) +
			springMap
		damageCheckSumUnfolded := damageCheckSum + "," +
			damageCheckSum + "," +
			damageCheckSum + "," +
			damageCheckSum + "," +
			damageCheckSum
		
		damagedSpringStrSl := strings.Split(damageCheckSumUnfolded, ",")
		damagedSpringsSequences := util.StringsToUints(damagedSpringStrSl)
		cnt := countDamageVariants(springMapUnfolded, damagedSpringsSequences)

		fmt.Printf("%d. Map %s %s has %d damage variants\n", lineIdx+1, springMap, damageCheckSum, cnt)

		damageVariantSum += cnt
	}

	fmt.Println("Damage variant sum:", damageVariantSum)
}

func countDamageVariants(springMap string, damagedSpringsCheckSum []uint) uint {
	var damagedSpringsTotal uint
	for _, cnt := range damagedSpringsCheckSum {
		damagedSpringsTotal += cnt
	}
	if damagedSpringsTotal == 0 {
		panic("Damaged springs total is 0")
	}

	var damagedSpringsMarked uint
	var unknownSpringsCount uint
	for _, r := range springMap {
		if r == runeDamagedSpring {
			damagedSpringsMarked++
		} else if r == runeUnknownSpring {
			unknownSpringsCount++
		}
	}

	if unknownSpringsCount == 0 {
		// panic("No unknown springs found!")
		fmt.Println("No unknown springs found!")
		return 1
	}

	if damagedSpringsMarked > damagedSpringsTotal {
		panic(fmt.Errorf("Marked damaged spring count %d is more than total damaged spring count %d", 
			damagedSpringsMarked, damagedSpringsTotal))
	}
	damagedSpringsToLocate := damagedSpringsTotal - damagedSpringsMarked

	if damagedSpringsToLocate == 0 {
		// panic("All damaged springs are marked!")
		fmt.Println("All damaged springs are marked!")
		return 1
	}
	if damagedSpringsToLocate == unknownSpringsCount {
		return 1
	}
	if damagedSpringsToLocate > unknownSpringsCount {
		panic("Not enough unknown springs to locate missing damaged springs")
	}

	variantMasks := createVariantMasks(unknownSpringsCount, damagedSpringsToLocate)
	// fmt.Printf("Variant masks generated: %d\n", len(variantMasks))
	// fmt.Println(variantMasks)

	var validMapCount uint
	for _, mask := range variantMasks {
		unknownSpringIdx := uint(0)
		// newSpringMaps := make(map[string]bool)

		var m strings.Builder
		for _, r := range springMap {
			if r == runeUnknownSpring {
				shouldPlaceDamagedSpring := mask[unknownSpringIdx]
				if shouldPlaceDamagedSpring {
					m.WriteRune(runeDamagedSpring)
				} else {
					m.WriteRune(runeOperationalSpring)
				}

				unknownSpringIdx++
			} else {
				m.WriteRune(r)
			}
		}

		newSpringMap := m.String()
		// if newSpringMaps[newSpringMap] {
		// 	panic("Duplicate spring map generated!")
		// } else {
		// 	newSpringMaps[newSpringMap] = true
		// }

		if isSpringMapValid(newSpringMap, damagedSpringsCheckSum) {
			validMapCount++

			// fmt.Printf("Map %d is valid: %s\n", validMapCount, newSpringMap)
		}
	}

	if validMapCount == 0 {
		panic(fmt.Errorf("No valid variants found: %s %v", springMap, damagedSpringsCheckSum))
	}

	return validMapCount
}

func createVariantMasks(unknownSpringsCount, damagedSpringsToLocate uint) [][]bool {
	var variants [][]bool

	combinationsCount := uint(1) << unknownSpringsCount
	for i := uint(0); i < combinationsCount; i++ {
		if uint(bits.OnesCount(i)) != damagedSpringsToLocate {
			continue
		}

		var variant []bool
		for j := uint(0); j < unknownSpringsCount; j++ {
			bitMask := uint(1) << j
			bitValue := i & bitMask > 0
			variant = append(variant, bitValue)
		}

		variants = append(variants, variant)
	}

	return variants
}

func isSpringMapValid(springMap string, damagedSpringSequences []uint) bool {
	damagedSpringSequenceLength := uint(0)
	var checkSum []uint
	
	for _, r := range springMap {
		if r == runeDamagedSpring {
			damagedSpringSequenceLength++
		} else if damagedSpringSequenceLength > 0 {
			checkSum = append(checkSum, damagedSpringSequenceLength)
			damagedSpringSequenceLength = 0
		}
	}

	if damagedSpringSequenceLength > 0 {
		checkSum = append(checkSum, damagedSpringSequenceLength)
	}

	return slices.Equal(checkSum, damagedSpringSequences)
}