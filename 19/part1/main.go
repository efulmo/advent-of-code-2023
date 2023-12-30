package main

import (
	"fmt"
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
	category string
	operator string
	value    uint16
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

	fmt.Printf("%d workflows are parsed:\n", len(workflows))
	for _, w := range workflows {
		fmt.Println(w)
	}

	var parts []map[string]uint16
	for lineIdx++; lineIdx < linesCount; lineIdx++ {
		line := strings.Trim(lines[lineIdx], "{}")
		categoryStrs := strings.Split(line, ",")

		part := make(map[string]uint16, 4)
		for _, catStr := range categoryStrs {
			name, valueStr := catStr[:1], catStr[2:]
			part[name] = uint16(util.ParseUintOrPanic(valueStr))
		}
		parts = append(parts, part)
	}

	fmt.Printf("%d parts are parsed:\n", len(parts))
	for _, p := range parts {
		fmt.Println(p)
	}

	var acceptedPartsSum uint
	for _, part := range parts {
		decision := analyzePart(part, "in", workflows)
		if decision == decisionAccept {
			for _, catVal := range part {
				acceptedPartsSum += uint(catVal)
			}
		}
	}
	fmt.Println("Accepted parts sum:", acceptedPartsSum)
}

func analyzePart(part map[string]uint16, workflowName string, workflows map[string][]Rule) string {
	if workflowName == decisionAccept || workflowName == decisionReject {
		return workflowName
	}

	workflow, found := workflows[workflowName]
	if !found {
		panic(fmt.Errorf("Invalid workflow name %s", workflowName))
	}

	for _, rule := range workflow {
		if rule.kind == kindRedirect {
			return analyzePart(part, rule.nextWorkflowName, workflows)
		}

		var partMatches bool
		switch rule.operator {
		case operatorMore:
			partMatches = part[rule.category] > rule.value
		case operatorLess:
			partMatches = part[rule.category] < rule.value
		default:
			panic(fmt.Errorf("Unexpected operator %s", rule.operator))
		}

		if partMatches {
			return analyzePart(part, rule.nextWorkflowName, workflows)
		}
	}

	panic(fmt.Errorf("Workflow %s doesn't end with redirect rule", workflowName))
}
