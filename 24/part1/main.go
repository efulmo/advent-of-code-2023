package main

import (
	"fmt"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

type Hailstone struct {
	lineIdx                         uint16
	startX, startY, startZ          uint
	velocityX, velocityY, velocityZ int
}

type TestArea struct {
	fromX, toX, fromY, toY float64
}

// for sample
// var testArea = TestArea{7, 27, 7, 27}
// for input
var testArea = TestArea{200000000000000, 400000000000000, 200000000000000, 400000000000000}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	hailstones := []Hailstone{}
	replacer := strings.NewReplacer("@", "", ",", "")
	for lineIdx, line := range lines {
		fields := strings.Fields(replacer.Replace(line))
		hailstones = append(hailstones, Hailstone{
			lineIdx:   uint16(lineIdx),
			startX:    util.ParseUintOrPanic(fields[0]),
			startY:    util.ParseUintOrPanic(fields[1]),
			startZ:    util.ParseUintOrPanic(fields[2]),
			velocityX: util.ParseIntOrPanic(fields[3]),
			velocityY: util.ParseIntOrPanic(fields[4]),
			velocityZ: util.ParseIntOrPanic(fields[5]),
		})
	}

	fmt.Printf("%d hailstones are parsed\n", len(hailstones))
	// fmt.Println(hailstones)

	crossesInTestArea := uint16(0)
	hailstonesCount := uint(len(hailstones))
	for i := uint(0); i < hailstonesCount-1; i++ {
		for j := i+1; j < hailstonesCount; j++ {
			stone1 := hailstones[i]
			stone2 := hailstones[j]

			if pathsCrossInTestArea(stone1, stone2) {
				crossesInTestArea++
			}
		}
	}
	fmt.Printf("%d hailstone paths crossed in the test area\n", crossesInTestArea)
}

func pathsCrossInTestArea(stone1, stone2 Hailstone) bool {
	time2divider := stone2.velocityX * stone1.velocityY - stone1.velocityX * stone2.velocityY
	if time2divider == 0 {
		fmt.Printf("Stones %d and %d never cross\n", stone1.lineIdx+1, stone2.lineIdx+1)
		return false
	}

	time2 := float64(stone1.velocityX * (int(stone2.startY) - int(stone1.startY)) -
		stone1.velocityY * (int(stone2.startX) - int(stone1.startX))) /
		float64(time2divider)
	if time2 < 0 {
		fmt.Printf("Stones %d and %d crossed in point past of the second stone\n", stone1.lineIdx+1, 
			stone2.lineIdx+1)
		return false
	}
	
	time1 := (float64(int(stone2.startX) - int(stone1.startX)) + float64(stone2.velocityX) * time2) /
		float64(stone1.velocityX)
	if time1 < 0 {
		fmt.Printf("Stones %d and %d crossed in point past of the first stone\n", stone1.lineIdx+1, 
			stone2.lineIdx+1)
		return false
	}
	
	x0 := float64(stone2.startX) + float64(stone2.velocityX) * time2
	y0 := float64(stone2.startY) + float64(stone2.velocityY) * time2
	
	withingTestArea := x0 >= testArea.fromX && x0 <= testArea.toX &&
		y0 >= testArea.fromY && y0 <= testArea.toY
	
	insideTestAreaStr := " inside test area"
	if !withingTestArea {
		insideTestAreaStr = " outside of test area"
	}

	fmt.Printf("Stones %d and %d cross in point %3.1f:%3.1f%s on time %.1f and %.1f\n", stone1.lineIdx+1, 
		stone2.lineIdx+1, x0, y0, insideTestAreaStr, time1, time2)

	return withingTestArea
}