package main

import (
	"fmt"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	var prevValueSum int
	for lineIdx, line := range lines {
		vals := util.StringsToInts(strings.Fields(line))
		prevValue := predictPrevValue(uint(lineIdx), vals)

		fmt.Printf("%d. %d... %v\n", lineIdx+1, prevValue, vals)

		prevValueSum += prevValue
	}

	fmt.Println("Next values sum:", prevValueSum)
}

func predictPrevValue(lineIdx uint, vals []int) int {
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
		upperRowPrevVal := upperRow[0] - row[0]

		updatedUpperRow := make([]int, 0, len(upperRow)+1)
		updatedUpperRow = append(updatedUpperRow, upperRowPrevVal)
		updatedUpperRow = append(updatedUpperRow, upperRow...)
		diffs[i-1] = updatedUpperRow
	}

	fmt.Printf("%d. After prev values are added:\n", lineIdx+1)
	for _, row := range diffs {
		fmt.Println(row)
	}

	return diffs[0][0]
}
