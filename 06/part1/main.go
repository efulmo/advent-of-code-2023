package main

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	times := util.StringsToUint(strings.Fields(lines[0])[1:])
	fmt.Println("Time:", times)
	distances := util.StringsToUint(strings.Fields(lines[1])[1:])
	fmt.Println("Distances:", distances)

	timesLen := uint(len(times))
	if len(times) != len(distances) {
		panic(fmt.Errorf("Number of times %d is different to number of distances - %d", timesLen,
			len(distances)))
	}

	winningWaysCountProd := uint(1)
	for i := uint(0); i < timesLen; i++ {
		time, distance := times[i], distances[i]

		boatSpeed, err := calculateBoatSpeed(time, distance)
		if err != nil {
			panic(errors.Join(fmt.Errorf("Failed to calculate boat speed for time %d, distance %d",
				time, distance)))
		}
		isGameSpeedInt := math.Round(boatSpeed) == boatSpeed

		var minChargingTimeToBeat uint
		if isGameSpeedInt {
			minChargingTimeToBeat = uint(boatSpeed) + 1
		} else {
			minChargingTimeToBeat = uint(math.Ceil(boatSpeed))
		}

		longestDistanceTimeF := float64(time) / 2
		longestDistanceTime := uint(longestDistanceTimeF)
		isLongestDistanceTimeInt := math.Round(longestDistanceTimeF) == longestDistanceTimeF

		if minChargingTimeToBeat >= longestDistanceTime {
			fmt.Printf("Unable to win in game %d. Mininal charging time to win is %d, while longest "+
				"disance may be covered after charging time %d", i+1, minChargingTimeToBeat,
				longestDistanceTime)
			continue
		}

		var winningWays uint
		if isLongestDistanceTimeInt {
			winningWays = (longestDistanceTime - minChargingTimeToBeat) * 2 + 1
		} else {
			winningWays = (longestDistanceTime - minChargingTimeToBeat + 1) * 2
		}
		fmt.Printf("Winning ways for game %d is %d\n", i+1, winningWays)

		winningWaysCountProd *= winningWays
	}

	fmt.Println("Winning ways prod:", winningWaysCountProd)
}

func calculateBoatSpeed(time, distance uint) (float64, error) {
	timeF, distanceF := float64(time), float64(distance)

	d := timeF*timeF - 4*distanceF
	if d < 0 {
		return 0.0, fmt.Errorf("Negative discriminant %f", d)
	}
	sqrtD := math.Sqrt(d)
	speed1 := (-timeF + sqrtD) / -2
	speed2 := (-timeF - sqrtD) / -2

	if speed1 < 0 && speed2 < 0 {
		return 0.0, errors.New("Both roots are negative")
	}

	if speed1 < 0 {
		return speed2, nil
	} else if speed2 < 0 {
		return speed1, nil
	}

	return min(speed1, speed2), nil
}
