package main

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

type Hailstone struct {
	line                            uint16
	xStart, yStart, zStart          int
	xVelocity, yVelocity, zVelocity int
}

const (
	simulationDuration = 1000 // time range to brute force
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	hailstones := []Hailstone{}
	stonesByXVelocity := make(map[int][]Hailstone, len(lines))
	stonesByYVelocity := make(map[int][]Hailstone, len(lines))
	stonesByZVelocity := make(map[int][]Hailstone, len(lines))

	replacer := strings.NewReplacer("@", "", ",", "")
	for lineIdx, line := range lines {
		fields := strings.Fields(replacer.Replace(line))
		hailstone := Hailstone{
			line:      uint16(lineIdx),
			xStart:    util.ParseIntOrPanic(fields[0]),
			yStart:    util.ParseIntOrPanic(fields[1]),
			zStart:    util.ParseIntOrPanic(fields[2]),
			xVelocity: util.ParseIntOrPanic(fields[3]),
			yVelocity: util.ParseIntOrPanic(fields[4]),
			zVelocity: util.ParseIntOrPanic(fields[5]),
		}
		hailstones = append(hailstones, hailstone)

		vx := hailstone.xVelocity
		stonesByXVelocity[vx] = append(stonesByXVelocity[vx], hailstone)

		vy := hailstone.yVelocity
		stonesByYVelocity[vy] = append(stonesByYVelocity[vy], hailstone)

		vz := hailstone.zVelocity
		stonesByZVelocity[vz] = append(stonesByZVelocity[vz], hailstone)
	}

	fmt.Printf("%d hailstones are parsed\n", len(hailstones))
	// fmt.Println(hailstones)

	rockXVelocities := detectRockVelocities(stonesByXVelocity, func(h Hailstone) int {
		return h.xStart
	})
	if len(rockXVelocities) == 0 {
		panic("No common rock X velocities detected")
	}
	fmt.Println("Rock X velocities detected:", util.MapKeysToSortedSlice(rockXVelocities))

	rockYVelocities := detectRockVelocities(stonesByYVelocity, func(h Hailstone) int {
		return h.yStart
	})
	if len(rockYVelocities) == 0 {
		panic("No common rock Y velocities detected")
	}
	fmt.Println("Rock Y velocities detected:", util.MapKeysToSortedSlice(rockYVelocities))

	rockZVelocities := detectRockVelocities(stonesByZVelocity, func(h Hailstone) int {
		return h.zStart
	})
	if len(rockZVelocities) == 0 {
		panic("No common rock Z velocities detected")
	}
	fmt.Println("Rock Z velocities detected:", util.MapKeysToSortedSlice(rockZVelocities))

	var xStart, yStart, zStart int
velocityLoop:
	for xVelocity := range rockXVelocities {
		for yVelocity := range rockYVelocities {
			for zVelocity := range rockZVelocities {
				for _, hailstone := range hailstones {
					// fmt.Printf("Checking hailstone #%d\n", hailstoneIdx+1)
					xStart, yStart, zStart, err = findFirstMatchingStartCoordinates(hailstones,
						hailstone, xVelocity, yVelocity, zVelocity)
					if err == nil {
						break velocityLoop
					}
				}
			}
		}
	}

	if err != nil {
		panic(fmt.Errorf("Failed to find matching rock start coordinates: %s", err.Error()))
	}
	fmt.Printf("Rock start coordinates: %d,%d,%d. Sum: %d\n", xStart, yStart, zStart,
		xStart+yStart+zStart)
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func primeFactors(n uint) []uint {
	if n <= 0 {
		return []uint{}
	}

	factors := []uint{}
	toFactor := n
	root := uint(math.Sqrt(float64(toFactor)))

	for i := uint(2); i <= root; i++ {
		if toFactor%i == 0 {
			factors = append(factors, i)

			toFactor /= i

			root = uint(math.Sqrt(float64(toFactor)))
			i = 1 // will be increased to 2 before next iteration
		}
	}

	// 1 isn't prime
	if toFactor != 1 {
		factors = append(factors, toFactor)
	}

	return factors
}

func getProducts(nums []uint) map[uint]bool {
	numsLen := uint(len(nums))
	if numsLen == 0 {
		return map[uint]bool{}
	}
	if numsLen == 1 {
		return map[uint]bool{nums[0]: true}
	}
	if numsLen == 2 {
		return map[uint]bool{nums[0] * nums[1]: true}
	}
	if numsLen > 64 {
		panic(fmt.Errorf("Too high number of nums: %d. Max: 64", numsLen))
	}

	products := make(map[uint]bool, numsLen)
	maxBitMaskValue := uint(math.Pow(2, float64(numsLen)) - 1)
	for bitMask := uint(0); bitMask <= maxBitMaskValue; bitMask++ {
		if bits.OnesCount(bitMask) > 1 {
			product := uint(1)
			for i := uint(0); i < numsLen; i++ {
				bit := bitMask & (1 << i)
				if bit != 0 {
					product *= nums[i]
				}
			}

			products[product] = true
		}
	}

	return products
}

func getFactors(n uint) map[uint]bool {
	primeFactors := primeFactors(n)

	primeFactorsExtended := make([]uint, 0, len(primeFactors)+1)
	primeFactorsExtended = append(primeFactorsExtended, 1)
	primeFactorsExtended = append(primeFactorsExtended, primeFactors...)

	products := getProducts(primeFactorsExtended)
	products[n] = true

	return products
}

func getPossibleVelocities(hailstoneVelocity int, diffFactors map[uint]bool) map[int]bool {
	possibleVelocities := make(map[int]bool, len(diffFactors))
	for factor := range diffFactors {
		iFactor := int(factor)
		possibleVelocities[hailstoneVelocity+iFactor] = true
		possibleVelocities[hailstoneVelocity-iFactor] = true
	}

	return possibleVelocities
}

func detectRockVelocities(
	stonesByVelocity map[int][]Hailstone,
	stoneToStart func(Hailstone) int,
) map[int]bool {
	pairCountByPossibleVelocity := make(map[int]uint)
	stonePairWithSameVelocityCount := uint(0)
	for velocity, stones := range stonesByVelocity {
		if len(stones) != 2 {
			continue
		}
		stonePairWithSameVelocityCount++

		startDiff := uint(abs(int(stoneToStart(stones[0])) - int(stoneToStart(stones[1]))))
		// fmt.Printf("Stones %d and %d have same velocity %d. Start diff %d\n",
		// 	stones[0].lineIdx, stones[1].lineIdx, velocity, startDiff)

		diffFactors := getFactors(startDiff)
		// fmt.Printf("Diff factors[%d]:\n", len(diffFactors))
		// fmt.Println(util.SetToSlice(diffFactors))

		possibleRockVelocities := getPossibleVelocities(velocity, diffFactors)
		// fmt.Printf("Possible rock velocities[%d]:\n", len(possibleRockVelocities))
		// fmt.Println(util.SetToSlice(possibleRockVelocities))

		for velocity := range possibleRockVelocities {
			pairCountByPossibleVelocity[velocity]++
		}
	}

	commonVelocities := map[int]bool{}
	for velocity, pairCount := range pairCountByPossibleVelocity {
		if pairCount == stonePairWithSameVelocityCount {
			commonVelocities[velocity] = true
		}
	}

	return commonVelocities
}

func findFirstMatchingStartCoordinates(
	hailstones []Hailstone,
	baseHailstone Hailstone,
	rockXVelocity, rockYVelocity, rockZVelocity int,
) (xStart, yStart, zStart int, err error) {
	hailstonesCount := uint(len(hailstones))
	stoneLineByCollisionTime := make(map[float64]uint16)

	for t := uint(0); t < simulationDuration; t++ {
		rockXStart := baseHailstone.xStart + int(t)*(baseHailstone.xVelocity-rockXVelocity)
		rockYStart := baseHailstone.yStart + int(t)*(baseHailstone.yVelocity-rockYVelocity)
		rockZStart := baseHailstone.zStart + int(t)*(baseHailstone.zVelocity-rockZVelocity)

		clear(stoneLineByCollisionTime)

		// check if rock can hit all other hailstones from these start coordinates on different moments
	hailstonesVerification:
		for stoneIdx := uint(1); stoneIdx < hailstonesCount; stoneIdx++ {
			stone := hailstones[stoneIdx]

			// rock and hailstone having same velocities will never collide if they don't start
			// movement from same coordinates; collision with all hailstones is guaranteed by the
			// puzzle description, so current hailstone start coordinates are rock start coordinates
			// as well!
			if stone.xVelocity == rockXVelocity && stone.yVelocity == rockYVelocity &&
				stone.zVelocity == rockZVelocity {
				return stone.xStart, stone.yStart, stone.zStart, nil
			}

			// fmt.Printf("%d. xDiff=%d yDiff=%d zDiff=%d\n", t, stone.xStart-rockXStart, stone.yStart-rockYStart,
			// 	stone.zStart-rockZStart)

			xCollisionTime := float64(stone.xStart-rockXStart) / float64(rockXVelocity-stone.xVelocity)
			yCollisionTime := float64(stone.yStart-rockYStart) / float64(rockYVelocity-stone.yVelocity)
			zCollisionTime := float64(stone.zStart-rockZStart) / float64(rockZVelocity-stone.zVelocity)

			collisionTimes := []float64{xCollisionTime, yCollisionTime, zCollisionTime}
			// fmt.Printf("t: %v times: %v\n", t, collisionTimes)

			var finiteCollisionTime float64
			finiteCollisionTimeFound := false
			for _, time := range collisionTimes {
				if !math.IsInf(time, 0) {
					finiteCollisionTime = time
					finiteCollisionTimeFound = true
					break
				}
			}
			if !finiteCollisionTimeFound {
				panic(fmt.Errorf("No finite collision time found in %v", collisionTimes))
			}

			for _, time := range collisionTimes {
				// rock and the hailstone have the same axis velocity; they collide all the time,
				// but not on all axes, so the search must continue
				if math.IsInf(time, 0) {
					continue
				}

				// collision time must be
				// 1. same for all axis
				// 2. in future
				// 3. integer
				if time != finiteCollisionTime || time < 0 || math.Mod(time, 1.0) != 0 {
					break hailstonesVerification
				}
			}

			// collision with current hailstone is valid here
			stoneLineByCollisionTime[finiteCollisionTime] = stone.line
		}

		if uint(len(stoneLineByCollisionTime)) == hailstonesCount-1 {
			return rockXStart, rockYStart, rockZStart, nil
		}
	}

	return 0, 0, 0, errors.New("No matching start and velocity found")
}
