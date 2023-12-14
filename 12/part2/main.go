package main

import (
	"fmt"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	runeDamagedSpring     = '#'
	runeOperationalSpring = '.'
	runeUnknownSpring     = '?'

	strDamagedSpring     = string(runeDamagedSpring)
	strOperationalSpring = string(runeOperationalSpring)
	strUnknownSpring     = string(runeUnknownSpring)
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	mapCache := make(map[string]uint)
	var damageVariantSum uint
	for lineIdx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			panic(fmt.Errorf("Unexpected fields count in the line - %d", len(fields)))
		}

		springMap := fields[0]
		damageCheckSum := fields[1]
		springMapUnfolded := springMap
		damageCheckSumUnfolded := damageCheckSum
		for i := uint(0); i < 4; i++ {
			springMapUnfolded += strUnknownSpring + springMap
			damageCheckSumUnfolded += "," + damageCheckSum
		}
		validateSpringMap(springMapUnfolded)

		checkSumStrSl := strings.Split(damageCheckSumUnfolded, ",")
		checkSum := util.StringsToUints(checkSumStrSl)

		variants := countDamageVariants(springMapUnfolded, checkSum, mapCache)

		fmt.Printf("%d. Map %s %s has %d damage variants\n", lineIdx+1, springMap, damageCheckSum,
			variants)

		damageVariantSum += variants
	}

	fmt.Println("Cache map size:", len(mapCache))
	fmt.Println("Damage variant sum:", damageVariantSum)
}

func validateSpringMap(m string) {
	for _, r := range m {
		if r != runeOperationalSpring && r != runeDamagedSpring && r != runeUnknownSpring {
			panic(fmt.Errorf("Unexpected run found in spring map: %c", r))
		}
	}
}

func countDamageVariants(
	springMap string,
	checkSum []uint,
	mapCache map[string]uint) uint {
	
	cacheKey := springMap + fmt.Sprintf("%v", checkSum)
	if vars, ok := mapCache[cacheKey]; ok {
		return vars
	}

	// spring map is empty
	if len(springMap) == 0 {
		// for good: all sequences are matched
		if (len(checkSum)) == 0 {
			util.DebugLog("%s %v: String map is empty as well as checksum - 1\n", springMap,
				checkSum)
			return 1
		}
		// for bad: some sequences left unmatched; invalid case
		util.DebugLog("%s %v: String map is empty but checksum isn't - 0\n", springMap,
			checkSum)
		return 0
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

	var checkSumSum uint
	for _, cnt := range checkSum {
		checkSumSum += cnt
	}

	// no damaged springs are expected anymore
	if checkSumSum == 0 {
		// but they exist; invalid case
		if damagedSpringsMarked > 0 {
			util.DebugLog("%s %v: Damaged springs are in the map, but they are not expected - 0\n",
				springMap, checkSum)
			mapCache[cacheKey] = 0
			return 0
		}
		// no damaged springs in the map; valid case
		util.DebugLog("%s %v: No damaged springs are in the map and no of them are expected - 1\n",
			springMap, checkSum)
		mapCache[cacheKey] = 1
		return 1
	}

	// no unknown springs
	if unknownSpringsCount == 0 && checkSumSum == 0 {
		util.DebugLog("%s %v: No unknown springs left - 1\n", springMap, checkSum)
		mapCache[cacheKey] = 1
		return 1
	}

	// checksum is too high; invalid case
	if checkSumSum > damagedSpringsMarked+unknownSpringsCount {
		util.DebugLog("%s %v: Checksum is too high - 0\n", springMap, checkSum)
		mapCache[cacheKey] = 0
		return 0
	}

	// to many damaged springs; invalid case
	if damagedSpringsMarked > checkSumSum {
		util.DebugLog("%s %v: Too many(%d) damaged springs are in the map for checksum - 0\n",
			springMap, checkSum, damagedSpringsMarked)
		mapCache[cacheKey] = 0
		return 0
	}

	damagedSpringsToLocate := checkSumSum - damagedSpringsMarked

	// to little unknown springs; invalid case
	if damagedSpringsToLocate > unknownSpringsCount {
		util.DebugLog("%s %v Too many damaged springs to locate(%d) for %d unknown springs - 0\n",
			springMap, checkSum, damagedSpringsToLocate, unknownSpringsCount)
		mapCache[cacheKey] = 0
		return 0
	}

	switch springMap[0] {
	case runeOperationalSpring:
		modifiedStringMap := strings.TrimLeft(springMap, strOperationalSpring)
		result := countDamageVariants(modifiedStringMap, checkSum, mapCache)

		mapCache[cacheKey] = result
		return result
	case runeDamagedSpring:
		firstSeq := checkSum[0]
		damagedSpringsAtBeginning := countStarting(springMap, runeDamagedSpring)
		if damagedSpringsAtBeginning == firstSeq {
			cutSpringMap := springMap[firstSeq:]
			cutCheckSum := checkSum[1:]

			// guarantee sequence end if next is unknown spring
			if uint(len(cutSpringMap)) > 0 {
				nextChar := cutSpringMap[0]
				if nextChar == runeUnknownSpring {
					cutSpringMap = strings.Replace(cutSpringMap, strUnknownSpring,
						strOperationalSpring, 1)
				}
			}

			result := countDamageVariants(cutSpringMap, cutCheckSum, mapCache)

			mapCache[cacheKey] = result
			return result
		} else if damagedSpringsAtBeginning > firstSeq {
			util.DebugLog("%s %v: Too long sequence of damaged springs at beginning - 0\n", springMap,
				checkSum)
			mapCache[cacheKey] = 0
			return 0
		} else if uint(len(springMap)) > damagedSpringsAtBeginning {
			nextChar := springMap[damagedSpringsAtBeginning]
			if nextChar == runeOperationalSpring {
				util.DebugLog("%s %v: Too short sequence of damaged springs at beginning - 0\n",
					springMap, checkSum)
				mapCache[cacheKey] = 0
				return 0
			}

			// next char is unknown spring
			modifiedSpringMap := strings.Replace(springMap, strUnknownSpring,
				strDamagedSpring, 1)
			result := countDamageVariants(modifiedSpringMap, checkSum, mapCache)
			
			mapCache[cacheKey] = result
			return result
		} else {
			util.DebugLog("%s %v: Unable to match first damaged springs sequence - 0\n", springMap,
				checkSum)
			mapCache[cacheKey] = 0
			return 0
		}
	case runeUnknownSpring:
		operationalCaseVariants := countDamageVariants(strOperationalSpring+springMap[1:],
			checkSum, mapCache)
		damagedCaseVariants := countDamageVariants(strDamagedSpring+springMap[1:],
			checkSum, mapCache)
		result := operationalCaseVariants + damagedCaseVariants

		mapCache[cacheKey] = result
		return result
	default:
		panic(fmt.Errorf("Unknown rune found in spring map: %c", springMap[0]))
	}
}

func countStarting(str string, run rune) uint {
	for i, r := range str {
		if r != run {
			return uint(i)
		}
	}
	return uint(len(str))
}
