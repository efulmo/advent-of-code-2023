package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

type End struct {
	x, y, z uint16
}

type Brick struct {
	id                  string
	lowerEnd, higherEnd End
}

func (b Brick) getProjection() Projection {
	minX, maxX := min(b.lowerEnd.x, b.higherEnd.x), max(b.lowerEnd.x, b.higherEnd.x)
	minY, maxY := min(b.lowerEnd.y, b.higherEnd.y), max(b.lowerEnd.y, b.higherEnd.y)

	return Projection{
		fromX: minX,
		toX:   maxX,
		fromY: minY,
		toY:   maxY,
	}
}

type Projection struct {
	fromX, toX, fromY, toY uint16
}

func (p Projection) overlapsWith(other Projection) bool {
	noOverlap := p.fromX > other.toX || // this is to the right of other
		p.toX < other.fromX || // this is to the left of other
		p.fromY > other.toY || // this is above of other
		p.toY < other.fromY // this is below of other

	return !noOverlap
}

func newBrick(lineIdx int, end1, end2 End) Brick {
	lowerEnd, higherEnd := end1, end2
	if lowerEnd.z > higherEnd.z {
		lowerEnd, higherEnd = higherEnd, lowerEnd
	}

	brick := Brick{
		id:        toStringId(uint16(lineIdx)),
		lowerEnd:  lowerEnd,
		higherEnd: higherEnd,
	}
	if brick.lowerEnd.z == 0 || brick.higherEnd.z == 0 {
		panic(fmt.Errorf("Block on line %d in on ground level", lineIdx+1))
	}

	axesState := []bool{
		brick.lowerEnd.x != brick.higherEnd.x,
		brick.lowerEnd.y != brick.higherEnd.y,
		brick.lowerEnd.z != brick.higherEnd.z,
	}

	var differentAxesCount uint8
	for _, axisIsDifferent := range axesState {
		if axisIsDifferent {
			differentAxesCount++
		}
	}
	if differentAxesCount > 1 {
		panic(fmt.Errorf("Block on line %d has more than 1 coordinate different: %v", lineIdx+1,
			brick))
	}

	return brick
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	brickById := make(map[string]Brick)
	brickIdsByLowerEndZ := make(map[uint16]map[string]bool)
	brickIdsByHigherEndZ := make(map[uint16]map[string]bool)
	var maxZ uint16

	for lineIdx, line := range lines {
		endStrs := strings.Split(line, "~")
		if len(endStrs) != 2 {
			panic(fmt.Errorf("Line %d has %d ends", lineIdx+1, len(endStrs)))
		}
		end1CoordStrs := strings.Split(endStrs[0], ",")
		end2CoordStrs := strings.Split(endStrs[1], ",")

		brick := newBrick(lineIdx, End{
			x: uint16(util.ParseUintOrPanic(end1CoordStrs[0])),
			y: uint16(util.ParseUintOrPanic(end1CoordStrs[1])),
			z: uint16(util.ParseUintOrPanic(end1CoordStrs[2])),
		},
			End{
				x: uint16(util.ParseUintOrPanic(end2CoordStrs[0])),
				y: uint16(util.ParseUintOrPanic(end2CoordStrs[1])),
				z: uint16(util.ParseUintOrPanic(end2CoordStrs[2])),
			},
		)

		brickById[brick.id] = brick

		addBrickIdToLevelMap(brickIdsByLowerEndZ, brick.lowerEnd.z, brick.id)
		addBrickIdToLevelMap(brickIdsByHigherEndZ, brick.higherEnd.z, brick.id)

		maxZ = max(maxZ, brick.higherEnd.z)
	}

	fmt.Printf("Parsed %d bricks:\n", len(brickById))
	fmt.Println(formatBricksMap(brickById, func(b1, b2 Brick) int {
		return strings.Compare(b1.id, b2.id)
	}, nil))

	supportingBricksById := applyGravity(brickById, brickIdsByLowerEndZ, brickIdsByHigherEndZ, maxZ)

	fmt.Println("After gravity is applied:")
	fmt.Println(formatBricksMap(brickById, func(b1, b2 Brick) int {
		return int(b1.lowerEnd.z) - int(b2.lowerEnd.z)
	}, supportingBricksById))

	theOnlySupportingBrickIds := make(map[string]bool)
	for _, brickIds := range supportingBricksById {
		if len(brickIds) == 1 {
			maps.Copy(theOnlySupportingBrickIds, brickIds)
		}
	}
	fmt.Println("Bricks safe to disintegrate:", len(brickById)-len(theOnlySupportingBrickIds))
}

/*
0  -> A
1  -> B
25 -> Z
26 -> BA
27 -> BB
*/
func toStringId(toConvert uint16) string {
	num := toConvert
	s := ""

	for {
		reminder := num % 26
		s = string(rune('A'+reminder)) + s

		num = (num - reminder) / 26
		if num == 0 {
			break
		}
	}

	return s
}

func addBrickIdToLevelMap(bricksByLevelMap map[uint16]map[string]bool, level uint16, brickId string) {
	levelBricks := bricksByLevelMap[level]
	if levelBricks == nil {
		levelBricks = make(map[string]bool)
	}
	levelBricks[brickId] = true
	bricksByLevelMap[level] = levelBricks
}

func removeBrickIdFromLevelMap(bricksByLevelMap map[uint16]map[string]bool, level uint16, brickId string) {
	levelBricks := bricksByLevelMap[level]
	if levelBricks == nil {
		return
	}
	delete(levelBricks, brickId)

	if len(levelBricks) == 0 {
		levelBricks = nil
	}
	bricksByLevelMap[level] = levelBricks
}

func applyGravity(
	brickById map[string]Brick,
	brickIdsByLowerEndZ, brickIdsByHigherEndZ map[uint16]map[string]bool,
	maxZ uint16,
) map[string]map[string]bool {
	supportingBricksById := make(map[string]map[string]bool)

	// over all levels
	for z := uint16(2); z <= maxZ; z++ {
		brickIds := brickIdsByLowerEndZ[z]

		// over all bricks on the current level
		for brickId := range brickIds {
			brick := brickById[brickId]
			prj := brick.getProjection()

			// over lower levels
			for zDiff := uint16(1); brick.lowerEnd.z-zDiff > 0; zDiff++ {
				targetZ := brick.lowerEnd.z - zDiff
				brickIdsOnThisZ := brickIdsByHigherEndZ[targetZ]

				supportingBrickIds := make(map[string]bool)
				for lowerBrickId := range brickIdsOnThisZ {
					lowerBrick := brickById[lowerBrickId]
					if prj.overlapsWith(lowerBrick.getProjection()) {
						supportingBrickIds[lowerBrickId] = true
					}
				}

				// some bricks found that would support current brick
				if len(supportingBrickIds) > 0 {
					// supporting brick is more than 1 level lower; apply gravity!
					if zDiff > 1 {
						removeBrickIdFromLevelMap(brickIdsByLowerEndZ, brick.lowerEnd.z, brickId)
						removeBrickIdFromLevelMap(brickIdsByHigherEndZ, brick.higherEnd.z, brickId)

						brick.lowerEnd.z = brick.lowerEnd.z - (zDiff - 1)
						brick.higherEnd.z = brick.higherEnd.z - (zDiff - 1)

						addBrickIdToLevelMap(brickIdsByLowerEndZ, brick.lowerEnd.z, brickId)
						addBrickIdToLevelMap(brickIdsByHigherEndZ, brick.higherEnd.z, brickId)

						fmt.Printf("Gravity moved brick %s %d levels below\n", brick.id, zDiff-1)
					}

					supportingBricksById[brickId] = supportingBrickIds

					break
				}
			}

			// lowest brick is still in the air; apply gravity!
			if len(supportingBricksById[brickId]) == 0 && brick.lowerEnd.z > 1 {
				levelDiff := brick.lowerEnd.z - 1

				removeBrickIdFromLevelMap(brickIdsByLowerEndZ, brick.lowerEnd.z, brickId)
				removeBrickIdFromLevelMap(brickIdsByHigherEndZ, brick.higherEnd.z, brickId)
				
				brick.lowerEnd.z = brick.lowerEnd.z - levelDiff
				brick.higherEnd.z = brick.higherEnd.z - levelDiff

				addBrickIdToLevelMap(brickIdsByLowerEndZ, brick.lowerEnd.z, brickId)
				addBrickIdToLevelMap(brickIdsByHigherEndZ, brick.higherEnd.z, brickId)

				fmt.Printf("Gravity moved brick %s %d level below - to ground\n", brick.id, levelDiff)
			}

			brickById[brickId] = brick
		}
	}

	return supportingBricksById
}

func formatBricksMap(
	brickById map[string]Brick,
	comparator func(b1, b2 Brick) int,
	supportingBricksById map[string]map[string]bool,
) string {
	var bricks []Brick
	for _, brick := range brickById {
		bricks = append(bricks, brick)
	}
	slices.SortFunc(bricks, comparator)
	
	return formatBricks(bricks, supportingBricksById)
}

func formatBricks(bricks []Brick, supportingBricksById map[string]map[string]bool) string {
	var formatted []string
	for _, b := range bricks {
		formatted = append(formatted, formatBrick(b, supportingBricksById[b.id]))
	}
	return strings.Join(formatted, "\n")
}

func formatBrick(b Brick, supportingBrickIds map[string]bool) string {
	var supportingBrickIdsSl []string
	for brickId := range supportingBrickIds {
		supportingBrickIdsSl = append(supportingBrickIdsSl, brickId)
	}
	slices.Sort(supportingBrickIdsSl)
	supportingBrickIdsStr := strings.Join(supportingBrickIdsSl, ",")

	return fmt.Sprintf("[%s]%d:%d:%d~%d:%d:%d^%s", b.id, b.lowerEnd.x, b.lowerEnd.y, b.lowerEnd.z,
		b.higherEnd.x, b.higherEnd.y, b.higherEnd.z, supportingBrickIdsStr)
}
