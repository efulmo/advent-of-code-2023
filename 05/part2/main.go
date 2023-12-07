package main

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"unicode"

	"github.com/efulmo/advent-of-code-2023/util"
)

type Rule struct {
	destStart   uint
	sourceStart uint
	length      uint
}

type RuleSet struct {
	label string
	rules []Rule
}

type SeedRange struct {
	start  uint
	length uint
}

func (rs RuleSet) validate() {
	rulesLen := len(rs.rules)
	for rule1Idx, rule1 := range rs.rules {
		rule1SourceEnd := rule1.sourceStart + rule1.length
		rule1DestEnd := rule1.destStart + rule1.length

		for rule2Idx := rule1Idx + 1; rule2Idx < rulesLen; rule2Idx++ {
			rule2 := rs.rules[rule2Idx]
			rule2SourceEnd := rule2.sourceStart + rule2.length
			rule2DestEnd := rule2.destStart + rule2.length

			if (rule2.sourceStart >= rule1.sourceStart && rule2.sourceStart < rule1SourceEnd) ||
				(rule2SourceEnd > rule1.sourceStart && rule2SourceEnd <= rule1SourceEnd) {
				panic(fmt.Errorf("Rule %d%v conflicts with rule %d%v in ruleset %s in source ranges",
					rule2Idx, rule2,
					rule1Idx, rule1, rs.label))
			}

			if (rule2.destStart >= rule1.destStart && rule2.destStart < rule1DestEnd) ||
				(rule2DestEnd > rule1.destStart && rule2DestEnd <= rule1DestEnd) {
				panic(fmt.Errorf("Rule %d%v conflicts with rule %d%v in ruleset %s in dest ranges",
					rule2Idx, rule2,
					rule1Idx, rule1, rs.label))
			}
		}
	}
}

func (rs RuleSet) apply(seed uint) uint {
	for _, r := range rs.rules {
		sourceEnd := r.sourceStart + r.length
		if seed >= r.sourceStart && seed < sourceEnd {
			destShift := r.destStart - r.sourceStart
			return seed + destShift
		}
	}

	return seed
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	seedsStr := strings.TrimPrefix(lines[0], "seeds: ")
	seedsStrSlice := strings.Fields(seedsStr)

	var seedRanges []SeedRange
	seedStrSliceLen := uint(len(seedsStrSlice))
	for i := uint(0); i < seedStrSliceLen; i += 2 {
		seedRanges = append(seedRanges, SeedRange{
			start:  util.ParseUintOrPanic(seedsStrSlice[i]),
			length: util.ParseUintOrPanic(seedsStrSlice[i+1]),
		})
	}

	var ruleSets []RuleSet
	linesLen := uint(len(lines))
	for lineIdx := uint(2); lineIdx < linesLen; {
		ruleSet, nextRuleSetLineIdx := readRuleSet(lines, lineIdx)

		fmt.Printf("Parsed rule set at line %d: %v\n", lineIdx+1, ruleSet)

		ruleSets = append(ruleSets, ruleSet)
		lineIdx = nextRuleSetLineIdx
	}

	fmt.Println("Seed ranges:", seedRanges)

	var chans []chan uint
	for seedRangeIdx, seedRange := range seedRanges {
		fmt.Printf("Processing seed range %d/%d of size %d\n", seedRangeIdx+1, len(seedRanges),
			seedRange.length)

		ch := make(chan uint)
		chans = append(chans, ch)

		go func(sr SeedRange) {
			var minSeed uint = math.MaxUint
			for i := uint(0); i < sr.length; i++ {
				seed := sr.start + i
				for _, ruleSet := range ruleSets {
					seed = uint(ruleSet.apply(uint(seed)))
				}

				minSeed = min(minSeed, seed)
			}
			ch <- minSeed
		}(seedRange)
	}

	var finalSeeds []uint
	for _, ch := range chans {
		finalSeeds = append(finalSeeds, <-ch)
	}

	fmt.Println("Min seed:", slices.Min(finalSeeds))
}

func readRuleSet(lines []string, ruleSetStartIdx uint) (RuleSet, uint) {
	var rules []Rule
	var ruleSetLabel string

	linesLen := uint(len(lines))
	lineIdx := uint(ruleSetStartIdx)
	for ; lineIdx < linesLen; lineIdx++ {
		line := lines[lineIdx]

		// end of rule set; break
		if len(line) == 0 {
			break
		}

		// rule set label; skip
		if !unicode.IsDigit(rune(line[0])) {
			ruleSetLabel = strings.TrimSuffix(line, ":")
			continue
		}

		ruleNums := strings.Fields(line)
		ruleNumsLen := len(ruleNums)
		if ruleNumsLen != 3 {
			panic(fmt.Errorf("Invalid length of rule nums at line %d: %d", lineIdx+1, ruleNumsLen))
		}

		rules = append(rules, Rule{
			destStart:   util.ParseUintOrPanic(ruleNums[0]),
			sourceStart: util.ParseUintOrPanic(ruleNums[1]),
			length:      util.ParseUintOrPanic(ruleNums[2]),
		})
	}

	ruleSet := RuleSet{ruleSetLabel, rules}
	ruleSet.validate()

	return ruleSet, lineIdx + 1
}
