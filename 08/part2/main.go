package main

import (
	"fmt"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

type Node struct {
	num           uint
	name          string
	leftNodeName  string
	rightNodeName string
}

const (
	startNodeNameSuffix  = "A"
	finishNodeNameSuffix = "Z"

	commandRight = 'R'
	commandLeft  = 'L'

	cycleDetectionRounds = 5
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	commands := lines[0]
	nodeByName := parseNodes(lines[2:])

	var startNodeNames []string
	for name := range nodeByName {
		if strings.HasSuffix(name, startNodeNameSuffix) {
			startNodeNames = append(startNodeNames, name)
		}
	}
	startNodeCount := uint(len(startNodeNames))

	fmt.Printf("%d nodes parsed, %d of them are starting nodes: %v\n", len(nodeByName),
		startNodeCount, startNodeNames)

	commandCount := uint(len(commands))
	var ghostPathLengths []uint

	for ghostIdx, startNodeName := range startNodeNames {
		var pathLengths []uint
		var pathHashes []uint
		nodeName := startNodeName

		for finishNodesMetCount := uint(0); finishNodesMetCount < cycleDetectionRounds; finishNodesMetCount++ {
			commandIdx := uint(0)
			stepsMade := uint(0)
			pathHash := uint(0)

			fmt.Printf("Ghost %d starts with node %s\n", ghostIdx, startNodeName)

			for {
				command := rune(commands[commandIdx])
				node := nodeByName[nodeName]

				var newNodeName string
				if command == commandLeft {
					newNodeName = node.leftNodeName
				} else {
					newNodeName = node.rightNodeName
				}

				// oldNodeName := nodeName
				nodeName = newNodeName
				stepsMade++
				pathHash += node.num

				// fmt.Printf("Step %d. %s -> %s\n", stepsMade, oldNodeName, newNodeName)

				// prepare to select next command
				commandIdx++
				if commandIdx >= commandCount {
					commandIdx = 0
				}

				if strings.HasSuffix(newNodeName, finishNodeNameSuffix) {
					break
				}
			}

			pathLengths = append(pathLengths, stepsMade)
			pathHashes = append(pathHashes, pathHash)
			fmt.Printf("Ghost %d reached finish node %s in %d steps. Path hash %d\n", ghostIdx,
				nodeName, stepsMade, pathHash)
		}

		if util.SliceContainsSameValue(pathLengths, pathLengths[0]) &&
			util.SliceContainsSameValue(pathHashes[1:], pathHashes[1]) {
			ghostPathLengths = append(ghostPathLengths, pathLengths[0])
		} else {
			panic(fmt.Errorf("Ghost %d path isn't cycled", ghostIdx))
		}
	}

	fmt.Println("Ghost path lengths:", ghostPathLengths)
	fmt.Println("Shortest common path length:", lcm(ghostPathLengths[0], ghostPathLengths[1],
		ghostPathLengths...))
}

func parseNodes(lines []string) map[string]Node {
	nodeByName := make(map[string]Node, len(lines))
	r := strings.NewReplacer("=", "", "(", "", ")", "", ",", "")

	for lineIdx, line := range lines {
		fields := strings.Fields(r.Replace(line))
		name := fields[0]
		nodeByName[name] = Node{
			num:           uint(lineIdx) + 1,
			name:          name,
			leftNodeName:  fields[1],
			rightNodeName: fields[2],
		}
	}

	return nodeByName
}

func gcd(a, b uint) uint {
	for b != 0 {
		temp := b
		b = a % b
		a = temp
	}
	return a
}

func lcm(a, b uint, integers ...uint) uint {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}
