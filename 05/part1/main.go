package main

import (
	"fmt"
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
	seedsStrSl := strings.Fields(seedsStr)
	var seeds []uint
	for _, seedStr := range seedsStrSl {
		seedInt := util.AtoiOrPanic(seedStr)
		seeds = append(seeds, uint(seedInt))
	}

	var ruleSets []RuleSet
	linesLen := uint(len(lines))
	for lineIdx := uint(2); lineIdx < linesLen; {
		ruleSet, nextRuleSetLineIdx := readRuleSet(lines, lineIdx)

		fmt.Printf("Parsed rule set at line %d: %v\n", lineIdx+1, ruleSet)

		ruleSets = append(ruleSets, ruleSet)
		lineIdx = nextRuleSetLineIdx
	}

	fmt.Println("Seeds:", seeds)
	for _, ruleSet := range ruleSets {
		for seedIdx, seed := range seeds {
			seeds[seedIdx] = uint(ruleSet.apply(uint(seed)))
		}

		fmt.Printf("%s applied: %v\n", ruleSet.label, seeds)
	}

	fmt.Println("Min seed:", slices.Min(seeds))
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
			destStart:   uint(util.AtoiOrPanic(ruleNums[0])),
			sourceStart: uint(util.AtoiOrPanic(ruleNums[1])),
			length:      uint(util.AtoiOrPanic(ruleNums[2])),
		})
	}

	ruleSet := RuleSet{ruleSetLabel, rules}
	ruleSet.validate()

	return ruleSet, lineIdx + 1
}
