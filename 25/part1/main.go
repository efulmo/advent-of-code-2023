package main

import (
	"fmt"
	"maps"
	"math"
	"math/big"
	"slices"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	edgeCountToDelete = 3

	maxPathLength = math.MaxUint
)

type Edge struct {
	smallerPart, biggerPart string
}

func (e Edge) String() string {
	return fmt.Sprintf("%s/%s", e.smallerPart, e.biggerPart)
}

func newEdge(part1, part2 string) Edge {
	smallerPart := min(part1, part2)
	biggerPart := max(part1, part2)

	return Edge{
		smallerPart: smallerPart,
		biggerPart:  biggerPart,
	}
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	// build edges map
	connectedPartsByName := make(map[string]map[string]bool, len(lines))
	for _, line := range lines {
		lineParts := strings.Split(line, ":")
		partName := lineParts[0]
		connectedPartNames := strings.Fields(lineParts[1])
		connectedPartsMap := make(map[string]bool, len(connectedPartNames))
		for _, partName := range connectedPartNames {
			connectedPartsMap[partName] = true
		}

		existingConnectedParts := connectedPartsByName[partName]
		if existingConnectedParts == nil {
			connectedPartsByName[partName] = connectedPartsMap
		} else {
			for partToConnect := range connectedPartsMap {
				existingConnectedParts[partToConnect] = true
			}
		}

		for connectedPartName := range connectedPartsMap {
			reverseConnectedParts := connectedPartsByName[connectedPartName]
			if reverseConnectedParts == nil {
				reverseConnectedParts = map[string]bool{
					partName: true,
				}
				connectedPartsByName[connectedPartName] = reverseConnectedParts
			} else {
				reverseConnectedParts[partName] = true
			}
		}
	}

	fmt.Println("Parts parsed:", len(connectedPartsByName))
	// for partName, connectedParts := range connectedPartsByName {
	// 	fmt.Printf("%s: %s\n", partName, strings.Join(util.MapKeysToSortedSlice(connectedParts), ", "))
	// }

	// collect stats on shortest paths
	partNames := util.MapKeysToSortedSlice(connectedPartsByName)
	partsCount := uint(len(partNames))
	edgeCountByEdge := make(map[Edge]uint)

	// fmt.Println(findShortestPathBfs(connectedPartsByName, "rzs", "xhk"))

	allCombinationsNumber := uint(partsCount * partsCount)
	currentCombinationNumber := uint(0)

collectEdgeUsageStatsLoop:
	for i := uint(0); i < partsCount; i++ {
		for j := uint(0); j < partsCount; j++ {
			currentCombinationNumber++
			startPart, endPart := partNames[i], partNames[j]

			// fmt.Printf("Finding path from %s to %s\n", part1, part2)
			path := findShortestPathDijkstra(startPart, endPart, connectedPartsByName)

			if currentCombinationNumber%1000 == 0 {
				fmt.Printf("%d/%d. Path from %s to %s is %v\n", currentCombinationNumber, allCombinationsNumber,
					startPart, endPart, path)
			}

			if currentCombinationNumber == 50_000 {
				break collectEdgeUsageStatsLoop
			}

			pathLen := uint(len(path))
			for pathPartIdx := uint(0); pathPartIdx < pathLen-1; pathPartIdx++ {
				part1 := path[pathPartIdx]
				part2 := path[pathPartIdx+1]
				edgeCountByEdge[newEdge(part1, part2)]++
			}
		}
	}

	// build edge stats
	edgesByUsageCount := make(map[uint]map[Edge]bool)
	for edge, count := range edgeCountByEdge {
		edges := edgesByUsageCount[count]
		if edges == nil {
			edges = make(map[Edge]bool)
			edgesByUsageCount[count] = edges
		}
		edges[edge] = true
	}

	edgeCounts := util.MapKeysToSortedSlice(edgesByUsageCount)
	fmt.Printf("Counts: %v\n", edgeCounts)

	fmt.Println("Edges by usage:")
	for i := len(edgeCounts) - 1; i >= 0; i-- {
		count := edgeCounts[i]
		edges := util.MapKeysToSlice(edgesByUsageCount[count])
		slices.SortFunc(edges, compareEdgesByStringForm)

		fmt.Printf("%d: %v\n", count, edges)
	}

	// find edges to delete
	edgesToDelete := make(map[Edge]bool)

selectingMiddleEdges:
	for i := len(edgeCounts) - 1; i >= 0; i-- {
		count := edgeCounts[i]
		edges := edgesByUsageCount[count]

		for edge := range edges {
			edgesToDelete[edge] = true
			if len(edgesToDelete) == edgeCountToDelete {
				break selectingMiddleEdges
			}
		}
	}
	edgesToDeleteSl := util.MapKeysToSlice(edgesToDelete)
	slices.SortFunc(edgesToDeleteSl, compareEdgesByStringForm)
	fmt.Println("Detected middle edges to delete:", edgesToDeleteSl)

	// delete middle edges
	for edge := range edgesToDelete {
		connectedParts := connectedPartsByName[edge.smallerPart]
		delete(connectedParts, edge.biggerPart)

		reverseConnectedParts := connectedPartsByName[edge.biggerPart]
		delete(reverseConnectedParts, edge.smallerPart)
	}

	// detect clusters
	partsByClusterId := detectClusters(connectedPartsByName)
	if len(partsByClusterId) != 2 {
		panic(fmt.Sprintf("%d clusters detected while 2 are expected!", len(partsByClusterId)))
	}

	clusterSizeProduct := uint(1)
	for clusterId, parts := range partsByClusterId {
		partsCount := uint(len(parts))
		fmt.Printf("Cluster %d has %d parts\n", clusterId, partsCount)
		fmt.Println(util.MapKeysToSortedSlice(parts))

		clusterSizeProduct *= partsCount
	}
	fmt.Println("Cluster size product:", clusterSizeProduct)
}

func compareEdgesByStringForm(e1, e2 Edge) int {
	return strings.Compare(e1.String(), e2.String())
}

func findShortestPathDijkstra(
	startPart, endPart string,
	connectedPartsByName map[string]map[string]bool,
) []string {
	if startPart == endPart {
		return []string{startPart}
	}

	partsToVisit := []string{startPart}
	visitedParts := map[string]bool{}
	pathFromStartToPart := map[string][]string{}

	for len(partsToVisit) > 0 {
		currentPart := partsToVisit[0]
		partsToVisit = partsToVisit[1:]

		currentPartPath := pathFromStartToPart[currentPart]

		if currentPart == endPart {
			return append(currentPartPath, endPart)
		}

		connectedParts := connectedPartsByName[currentPart]
		for connectedPart := range connectedParts {
			if visitedParts[connectedPart] {
				continue
			}

			connectedPartPath, pathExists := pathFromStartToPart[connectedPart]
			if !pathExists || len(connectedPartPath) > len(currentPartPath)+1 {
				newConnectedPartPath := slices.Clone(currentPartPath)
				newConnectedPartPath = append(newConnectedPartPath, currentPart)
				pathFromStartToPart[connectedPart] = newConnectedPartPath
			}

			partsToVisit = append(partsToVisit, connectedPart)
		}

		visitedParts[currentPart] = true
	}

	return append(pathFromStartToPart[endPart], endPart)
}

func detectClusters(connectedPartsByName map[string]map[string]bool) map[uint]map[string]bool {
	if len(connectedPartsByName) == 0 {
		return map[uint]map[string]bool{}
	}

	partsOutsideClusters := make(map[string]bool)
	for part := range connectedPartsByName {
		partsOutsideClusters[part] = true
	}

	clusterIdByPart := make(map[string]uint)
	currentClusterId := uint(1)

	for len(partsOutsideClusters) > 0 {
		firstClusterPart, _ := util.GetAnyMapKey(partsOutsideClusters)

		fmt.Printf("New cluster ID=%d is detected\n", currentClusterId)
		clusterIdByPart[firstClusterPart] = currentClusterId
		delete(partsOutsideClusters, firstClusterPart)

		partsToVisit := maps.Clone(connectedPartsByName[firstClusterPart])
		for len(partsToVisit) > 0 {
			part, _ := util.GetAnyMapKey(partsToVisit)
			delete(partsToVisit, part)

			partClusterId := clusterIdByPart[part]
			if partClusterId != 0 {
				continue
			}

			clusterIdByPart[part] = currentClusterId
			delete(partsOutsideClusters, part)

			maps.Copy(partsToVisit, connectedPartsByName[part])
		}

		currentClusterId++
	}

	partsByClusterId := make(map[uint]map[string]bool)
	for part, clusterId := range clusterIdByPart {
		parts := partsByClusterId[clusterId]
		if parts == nil {
			parts = make(map[string]bool)
		}

		parts[part] = true
		partsByClusterId[clusterId] = parts
	}

	return partsByClusterId
}
