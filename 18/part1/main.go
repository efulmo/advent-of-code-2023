package main

import (
	"fmt"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	directionUp    = "U"
	directionDown  = "D"
	directionLeft  = "L"
	directionRight = "R"
)

type Coord struct {
	x, y int
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	startCoord := Coord{
		x: 0,
		y: 0,
	}
	coords := []Coord{startCoord}
	var perimiter uint

	for _, line := range lines {
		fields := strings.Fields(line)
		direction, lengthStr := fields[0], fields[1]
		length := util.ParseUintOrPanic(lengthStr)

		prevCoord := coords[len(coords)-1]
		var newCoord Coord

		switch direction {
		case directionUp:
			newCoord = Coord{
				x: prevCoord.x,
				y: prevCoord.y + int(length),
			}
		case directionDown:
			newCoord = Coord{
				x: prevCoord.x,
				y: prevCoord.y - int(length),
			}
		case directionRight:
			newCoord = Coord{
				x: prevCoord.x + int(length),
				y: prevCoord.y,
			}
		case directionLeft:
			newCoord = Coord{
				x: prevCoord.x - int(length),
				y: prevCoord.y,
			}
		default:
			panic(fmt.Errorf("Unknown direction: %s", direction))
		}

		coords = append(coords, newCoord)
		perimiter += length

		// fmt.Printf("%s %d: %d:%d -> %d:%d\n", direction, length, prevCoord.x, prevCoord.y,
		// 	newCoord.x, newCoord.y)
	}

	coordsLen := uint(len(coords))
	if coords[0] != coords[coordsLen-1] {
		panic(fmt.Errorf("Lava pool isn't closed"))
	}

	var sum int
	for i := uint(0); i < coordsLen-1; i++ {
		coord, nextCoord := coords[i], coords[i+1]

		toAdd := coord.x * nextCoord.y
		toSubtract := nextCoord.x * coord.y
		sum += toAdd - toSubtract

		// fmt.Printf("calculations: add %d, subtract %d\n", toAdd, toSubtract)
	}

	internalArea := abs(sum) / 2
	edgeArea := int(perimiter)/2 + 1
	fmt.Println("Area:", internalArea+edgeArea)
}

func abs(i int) int {
	if i >= 0 {
		return i
	}
	return -i
}
