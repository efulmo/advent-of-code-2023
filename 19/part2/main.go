package main

import (
	"fmt"
	"maps"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	kindCondition = 1
	kindRedirect  = 2

	decisionAccept = "A"
	decisionReject = "R"

	operatorMore = ">"
	operatorLess = "<"
)

type Rule struct {
	kind             uint8
	nextWorkflowName string

	// kindConditionOnly
	category, operator string
	value              uint16
}

type Range struct {
	moreThan, lessThan uint16
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	linesCount := uint(len(lines))
	var lineIdx uint
	workflows := make(map[string][]Rule)

	for lineIdx = uint(0); lineIdx < linesCount; lineIdx++ {
		line := lines[lineIdx]
		if len(line) == 0 {
			break
		}

		ruleParts := strings.Split(line, "{")
		name := ruleParts[0]

		rulesStr := strings.TrimRight(ruleParts[1], "}")
		rulesStrSl := strings.Split(rulesStr, ",")

		rules := make([]Rule, 0, 2)
		for _, rule := range rulesStrSl {
			if strings.ContainsAny(rule, "<>:") {
				category, operator, valueAndNextWorkflow := string(rule[0]), string(rule[1]), rule[2:]
				valueParts := strings.Split(valueAndNextWorkflow, ":")

				rules = append(rules, Rule{
					kind:             kindCondition,
					category:         category,
					operator:         operator,
					value:            uint16(util.ParseUintOrPanic(valueParts[0])),
					nextWorkflowName: valueParts[1],
				})

			} else {
				rules = append(rules, Rule{
					kind:             kindRedirect,
					nextWorkflowName: rule,
				})
			}
		}

		workflows[name] = rules
	}

	fmt.Printf("%d workflows are parsed\n", len(workflows))

	categoryCombination := map[string]Range{
		"a": {moreThan: 0, lessThan: 4001},
		"m": {moreThan: 0, lessThan: 4001},
		"s": {moreThan: 0, lessThan: 4001},
		"x": {moreThan: 0, lessThan: 4001},
	}

	acceptedCombos := getAcceptedCombos("in", workflows, categoryCombination)
	acceptedCombosLen := uint(len(acceptedCombos))
	fmt.Printf("%d derived combos found\n", acceptedCombosLen)

	uniqueAcceptedCombos := make(map[string]map[string]Range, len(acceptedCombos))
	for _, combo := range acceptedCombos {
		aRange := combo["a"]
		mRange := combo["m"]
		sRange := combo["s"]
		xRange := combo["x"]
		hash := fmt.Sprintf("%d<a<%d,%d<m<%d,%d<s<%d,%d<x<%d", aRange.lessThan, aRange.moreThan,
			mRange.lessThan, mRange.moreThan, sRange.lessThan, sRange.moreThan, xRange.lessThan,
			xRange.moreThan)
		uniqueAcceptedCombos[hash] = combo
	}
	fmt.Printf("%d unique combos left after filtering\n", len(uniqueAcceptedCombos))

	var totalCombos uint
	for _, combo := range uniqueAcceptedCombos {
		thisComboCount := uint(1)
		for _, aRange := range combo {
			diff := int(aRange.lessThan) - int(aRange.moreThan) - 1
			if diff > 0 {
				thisComboCount *= uint(diff)
			}
		}
		totalCombos += thisComboCount
	}
	fmt.Println("Total combo count:", totalCombos)
}

func getAcceptedCombos(
	workflowName string,
	workflows map[string][]Rule,
	prevCombo map[string]Range,
) []map[string]Range {
	if workflowName == decisionAccept {
		return []map[string]Range{prevCombo}
	}
	if workflowName == decisionReject {
		return []map[string]Range{}
	}

	workflow, found := workflows[workflowName]
	if !found {
		panic(fmt.Errorf("Invalid workflow name %s", workflowName))
	}

	curCombo := maps.Clone(prevCombo)
	var derivedCombos []map[string]Range

	for _, rule := range workflow {
		if rule.kind == kindRedirect {
			childCombos := getAcceptedCombos(rule.nextWorkflowName, workflows, curCombo)
			derivedCombos = append(derivedCombos, childCombos...)
		} else {
			switch rule.operator {
			case operatorMore:
				newRange := curCombo[rule.category]
				newRange.moreThan = max(rule.value, newRange.moreThan)
	
				matchingCombo := maps.Clone(curCombo)
				matchingCombo[rule.category] = newRange
	
				childCombos := getAcceptedCombos(rule.nextWorkflowName, workflows, matchingCombo)
				derivedCombos = append(derivedCombos, childCombos...)
	
				// change curCombo to fail the condition and proceed to the next rule
				failRange := curCombo[rule.category]
				failRange.lessThan = min(rule.value, failRange.lessThan) + 1
	
				curCombo[rule.category] = failRange
			case operatorLess:
				newRange := curCombo[rule.category]
				newRange.lessThan = min(rule.value, newRange.lessThan)
	
				matchingCombo := maps.Clone(curCombo)
				matchingCombo[rule.category] = newRange
	
				childCombos := getAcceptedCombos(rule.nextWorkflowName, workflows, matchingCombo)
				derivedCombos = append(derivedCombos, childCombos...)
	
				// change curCombo to fail the condition and proceed to the next rule
				failRange := curCombo[rule.category]
				failRange.moreThan = max(rule.value, failRange.moreThan) - 1
	
				curCombo[rule.category] = failRange
			default:
				panic(fmt.Errorf("Unexpected operator <%s>", rule.operator))
			}
		}
	}

	return derivedCombos
}
