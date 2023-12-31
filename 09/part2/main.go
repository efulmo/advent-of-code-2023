package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	var nextValueSum int
	for lineIdx, line := range lines {
		vals := util.StringsToInts(strings.Fields(line))
		slices.Reverse(vals)
		nextValue := predictNextValue(uint(lineIdx), vals)

		fmt.Printf("%d. %v... %d\n", lineIdx+1, vals, nextValue)

		nextValueSum += nextValue
	}

	fmt.Println("Next values sum:", nextValueSum)
}

func predictNextValue(lineIdx uint, vals []int) int {
	var diffs [][]int
	diffs = append(diffs, vals)

	for {
		row := diffs[len(diffs)-1]
		if util.SliceContainsSameValue(row, 0) {
			break
		}
		rowLen := uint(len(row))

		newRow := make([]int, 0, rowLen-1)
		for i := uint(0); i < rowLen-1; i++ {
			diff := row[i+1] - row[i]
			newRow = append(newRow, diff)
		}
		diffs = append(diffs, newRow)
	}

	fmt.Printf("%d. Calculated all diffs:\n", lineIdx+1)
	for _, row := range diffs {
		fmt.Println(row)
	}

	for i := uint(len(diffs)) - 1; i > 0; i-- {
		row := diffs[i]
		upperRow := diffs[i-1]
		upperRowNextVal := upperRow[len(upperRow)-1] + row[len(row)-1]
		diffs[i-1] = append(upperRow, upperRowNextVal)
	}

	fmt.Printf("%d. After next values are added:\n", lineIdx+1)
	for _, row := range diffs {
		fmt.Println(row)
	}

	return diffs[0][len(diffs[0])-1]
}
