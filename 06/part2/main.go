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

	time := util.ParseUintOrPanic(strings.ReplaceAll(strings.TrimPrefix(lines[0], "Time: "), " ", ""))
	fmt.Println("Time:", time)
	distance := util.ParseUintOrPanic(
		strings.ReplaceAll(strings.TrimPrefix(lines[1], "Distance: "), " ", ""))
	fmt.Println("Distance:", distance)

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
		panic(fmt.Errorf("Unable to win in the game. Mininal charging time to win is %d, while "+
			"longest disance may be covered after charging time %d", minChargingTimeToBeat,
			longestDistanceTime))
	}

	var winningWays uint
	if isLongestDistanceTimeInt {
		winningWays = (longestDistanceTime-minChargingTimeToBeat)*2 + 1
	} else {
		winningWays = (longestDistanceTime - minChargingTimeToBeat + 1) * 2
	}

	fmt.Println("Winning ways for game is", winningWays)
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
