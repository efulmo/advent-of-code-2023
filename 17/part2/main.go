package main

import (
	"container/list"
	"fmt"
	"slices"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	directionRight = "r"
	directionDown  = "d"
	directionLeft  = "l"
	directionUp    = "u"
	directionNone  = "n"

	infinityHealLoss = 999_999_999

	minStepsInSameDirection = 4
	maxStepsInSameDirection = 10
)

type Coord struct {
	rowIdx, colIdx uint8
}

type Node struct {
	coord                Coord
	inDirection          string
	stepsMadeInDirection uint8
}

type NodeInfo struct {
	totalHeatLoss uint
	prevNode      Node
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	rowsTotal := uint8(len(lines))
	colsTotal := uint8(len(lines[0]))

	heatLossByCoord := make(map[Coord]uint8, rowsTotal*colsTotal)
	for rowIdx, row := range lines {
		for colIdx := range row {
			coord := Coord{
				rowIdx: uint8(rowIdx),
				colIdx: uint8(colIdx),
			}
			heatLossByCoord[coord] = uint8(util.ParseUintOrPanic(row[colIdx : colIdx+1]))
		}
	}

	startNode := Node{
		coord:                Coord{0, 0},
		inDirection:          directionNone,
		stepsMadeInDirection: 0,
	}
	nodeInfos := make(map[Node]NodeInfo, rowsTotal*colsTotal*4)
	nodeInfos[startNode] = NodeInfo{
		totalHeatLoss: 0,
		prevNode:      startNode,
	}

	analyzedNodes := make(map[Node]bool)

	nodesToAnalyze := []Node{startNode}
	nodeComparator := byTotalHeatLossComparator(nodeInfos)

	for len(nodesToAnalyze) > 0 {
		// heatLosses := nodesToTotalHeatLoss(nodesToAnalyze, nodeInfos)
		// fmt.Println("Selecting 1st node from nodes to analyze:", heatLosses)
		// if !slices.IsSortedFunc(nodesToAnalyze, nodeComparator) || !slices.IsSorted(heatLosses) {
		// 	panic(fmt.Errorf("Nodes to analyzed aren't sorted: %v", heatLosses)))
		// }

		// find next node with minimal total heat loss
		currentNode := nodesToAnalyze[0]
		nodesToAnalyze = nodesToAnalyze[1:]
		currentNodeInfo := nodeInfos[currentNode]

		neighbourNodes := getNeighbourNodes(currentNode, analyzedNodes, rowsTotal, colsTotal)

		// coord := currentNode.coord
		// util.DebugLog("Analyzing node %d:%d:%s:%d with total hit loss %d and %d neighbours. Queue: %d\n",
		// 	coord.rowIdx+1, coord.colIdx+1, currentNode.inDirection, currentNode.stepsMadeInDirection,
		// 	currentNodeInfo.totalHeatLoss, len(neighbourNodes),
		// 	len(nodesToAnalyze),
		// )

		for _, neighbourNode := range neighbourNodes {
			neighbourNodeCoord := neighbourNode.coord
			neighbourNodeInfo := nodeInfos[neighbourNode]
			neighbourNodeTotalHeatLoss := neighbourNodeInfo.totalHeatLoss
			if neighbourNodeTotalHeatLoss == 0 {
				neighbourNodeTotalHeatLoss = infinityHealLoss
			}

			newTotalHeatLoss := currentNodeInfo.totalHeatLoss + uint(heatLossByCoord[neighbourNodeCoord])

			// util.DebugLog("New total heat loss for %d:%d:%s:%d node: %d. Old: %d\n",
			// 	neighbourNodeCoord.rowIdx+1, neighbourNodeCoord.colIdx+1, neighbourNode.inDirection,
			// 	neighbourNode.stepsMadeInDirection, newTotalHeatLoss, neighbourNodeTotalHeatLoss)

			if newTotalHeatLoss < neighbourNodeTotalHeatLoss {
				// util.DebugLog("New total heat loss is set for node %s - %d. Old: %d\n",
				// 	formatNode(neighbourNode, heatLossByCoord),
				// 	newTotalHeatLoss, neighbourNodeTotalHeatLoss)

				nodeInfos[neighbourNode] = NodeInfo{
					totalHeatLoss: newTotalHeatLoss,
					prevNode:      currentNode,
				}

				nodesToAnalyze = insertSortingAwareIfAbsent(nodesToAnalyze, neighbourNode, nodeComparator)
			}
		}

		analyzedNodes[currentNode] = true
	}

	finishCoord := Coord{
		rowIdx: rowsTotal - 1,
		colIdx: colsTotal - 1,
	}
	var finishNode Node
	finishNodeFound := false
	for node, nodeInfo := range nodeInfos {
		if node.coord == finishCoord {
			fmt.Printf("Finish node found %s with total heat loss %d\n", formatNode(node,
				heatLossByCoord), nodeInfo.totalHeatLoss)

			if node.stepsMadeInDirection >= minStepsInSameDirection &&
				node.stepsMadeInDirection <= maxStepsInSameDirection && (
					!finishNodeFound || nodeInfo.totalHeatLoss < nodeInfos[finishNode].totalHeatLoss) {
				finishNode = node
				finishNodeFound = true
			}
		}
	}
	finishNodeInfo := nodeInfos[finishNode]

	fmt.Println("Unique analyzed coords:", len(analyzedNodes))
	if !finishNodeFound {
		fmt.Println("No path to finish node found")
	} else {
		fmt.Printf("Finish node: %s. Heat loss: %d\n", formatNode(finishNode, heatLossByCoord),
			finishNodeInfo.totalHeatLoss)
		// fmt.Println("Optimal path:")
		// formattedPathNodes := formatNodePath(buildNodePath(finishNode, nodeInfos), heatLossByCoord)
		// for step, formattedMode := range formattedPathNodes {
		// 	fmt.Printf("%d. %s\n", step+1, formattedMode)
		// }
	}
}

func byTotalHeatLossComparator(nodeInfos map[Node]NodeInfo) func(Node, Node) int {
	return func(n1, n2 Node) int {
		info1, found1 := nodeInfos[n1]
		if !found1 {
			panic(fmt.Errorf("Cannot find node info for %v", n1))
		}

		info2, found2 := nodeInfos[n2]
		if !found2 {
			panic(fmt.Errorf("Cannot find node info for %v", n2))
		}

		return int(info1.totalHeatLoss) - int(info2.totalHeatLoss)
	}
}

func insertSortingAwareIfAbsent(nodes []Node, nodeToInsert Node, comparator func(Node, Node) int) []Node {
	newNodeIdx, _ := slices.BinarySearchFunc(nodes, nodeToInsert, comparator)
	isNodeAlreadyPresent := slices.Contains(nodes, nodeToInsert)

	if !isNodeAlreadyPresent {
		nodes = slices.Insert(nodes, newNodeIdx, nodeToInsert)
	}

	return nodes
}

func getNeighbourNodes(node Node, analyzedNodes map[Node]bool, rowsTotal, colsTotal uint8) []Node {
	var allowedNextMoveDirections []string

	if node.inDirection != directionNone && node.stepsMadeInDirection < minStepsInSameDirection {
		allowedNextMoveDirections = []string{node.inDirection}
	} else {
		switch node.inDirection {
		// starting node only; may go to any direction
		case directionNone:
			allowedNextMoveDirections = []string{directionRight, directionDown}
		case directionRight:
			allowedNextMoveDirections = []string{directionRight, directionDown, directionUp}
		case directionDown:
			allowedNextMoveDirections = []string{directionRight, directionDown, directionLeft}
		case directionLeft:
			allowedNextMoveDirections = []string{directionDown, directionLeft, directionUp}
		case directionUp:
			allowedNextMoveDirections = []string{directionRight, directionLeft, directionUp}
		default:
			panic(fmt.Errorf("Unexpected direction %s", node.inDirection))
		}
	}

	coord := node.coord

	var neighboursNodes []Node
	for _, direction := range allowedNextMoveDirections {
		stepsMadeInThatDirection := uint8(0)
		if node.inDirection == direction {
			stepsMadeInThatDirection = node.stepsMadeInDirection
		}

		// all steps in the direction are made; proceed with other directions
		if stepsMadeInThatDirection >= maxStepsInSameDirection {
			continue
		}

		var newCoord Coord
		var isNewCoordValid bool
		switch direction {
		case directionRight:
			newCoord, isNewCoordValid = newCoordIfValid(int(coord.rowIdx), int(coord.colIdx+1),
				rowsTotal, colsTotal)
		case directionDown:
			newCoord, isNewCoordValid = newCoordIfValid(int(coord.rowIdx+1), int(coord.colIdx),
				rowsTotal, colsTotal)
		case directionLeft:
			newCoord, isNewCoordValid = newCoordIfValid(int(coord.rowIdx), int(coord.colIdx)-1,
				rowsTotal, colsTotal)
		case directionUp:
			newCoord, isNewCoordValid = newCoordIfValid(int(coord.rowIdx)-1, int(coord.colIdx),
				rowsTotal, colsTotal)
		}

		if isNewCoordValid {
			newNode := Node{
				coord:                newCoord,
				inDirection:          direction,
				stepsMadeInDirection: stepsMadeInThatDirection + 1,
			}
			if !analyzedNodes[newNode] {
				neighboursNodes = append(neighboursNodes, newNode)
			}
		}
	}

	return neighboursNodes
}

func newCoordIfValid(rowIdx, colIdx int, rowsTotal, colsTotal uint8) (Coord, bool) {
	if rowIdx < 0 || rowIdx >= int(rowsTotal) ||
		colIdx < 0 || colIdx >= int(colsTotal) {
		return Coord{}, false
	}
	return Coord{uint8(rowIdx), uint8(colIdx)}, true
}

func nodesToTotalHeatLoss(nodes []Node, nodeInfos map[Node]NodeInfo) []uint {
	var heatLosses []uint
	for _, node := range nodes {
		heatLosses = append(heatLosses, nodeInfos[node].totalHeatLoss)
	}

	return heatLosses
}

func buildNodePath(finishNode Node, nodeInfos map[Node]NodeInfo) []Node {
	nodes := list.New()

	var startCoord Coord
	for node := finishNode; node.coord != startCoord; node = nodeInfos[node].prevNode {
		nodes.PushFront(node)
	}
	nodes.PushFront(nodeInfos[nodes.Front().Value.(Node)].prevNode)

	nodesSl := make([]Node, 0, nodes.Len())
	for e := nodes.Front(); e != nil; e = e.Next() {
		nodesSl = append(nodesSl, e.Value.(Node))
	}

	return nodesSl
}

func formatNodePath(nodePath []Node, heatLossByCoord map[Coord]uint8) []string {
	var coordsStr []string
	for _, node := range nodePath {
		coordsStr = append(coordsStr, formatNode(node, heatLossByCoord))
	}

	return coordsStr
}

func formatNode(node Node, heatLossByCoord map[Coord]uint8) string {
	coord := node.coord
	return fmt.Sprintf("%d:%d:%s:%d[%d]", coord.rowIdx+1, coord.colIdx+1, node.inDirection,
		node.stepsMadeInDirection, heatLossByCoord[coord])
}
