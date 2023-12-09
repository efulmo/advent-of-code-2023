package main

import (
	"fmt"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

type Node struct {
	name          string
	leftNodeName  string
	rightNodeName string
}

const (
	startNodeName  = "AAA"
	finishNodeName = "ZZZ"

	commandRight = 'R'
	commandLeft  = 'L'
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	commands := lines[0]
	nodeByName := parseNodes(lines[2:])

	commandIdx := uint(0)
	nodeName := startNodeName
	commandCount := uint(len(commands))
	stepsMade := uint(0)
	
	fmt.Println("At node AAA")

	for {
		node := nodeByName[nodeName]
		command := rune(commands[commandIdx])
		
		prevNodeName := nodeName
		if command == commandLeft {
			nodeName = node.leftNodeName
		} else {
			nodeName = node.rightNodeName
		}

		commandIdx++
		if commandIdx >= commandCount {
			commandIdx = 0
		}
		
		stepsMade++
		fmt.Printf("Step %d. %s -> %s\n", stepsMade, prevNodeName, nodeName)

		if nodeName == finishNodeName {
			break
		}
	}
	fmt.Printf("The way took %d steps\n", stepsMade)
}

func parseNodes(lines []string) map[string]Node {
	nodeByName := make(map[string]Node, len(lines))
	r := strings.NewReplacer("=", "", "(", "", ")", "", ",", "")

	for _, line := range lines {
		fields := strings.Fields(r.Replace(line))
		name := fields[0]
		nodeByName[name] = Node{
			name:          name,
			leftNodeName:  fields[1],
			rightNodeName: fields[2],
		}
	}

	return nodeByName
}
