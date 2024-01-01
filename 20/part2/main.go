package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	pulseKindLow  = "low"
	pulseKindHigh = "high"

	moduleKindButton      = "button"
	moduleKindBroadcaster = "broadcaster"
	moduleKindFlipFlop    = "%"
	moduleKindConjunction = "&"

	stateOff = false
)

type Module struct {
	kind, name        string
	outputModuleNames []string

	// for kindConjunction only
	inputModuleNames []string
}

type Pulse struct {
	sourceModuleName, kind, targetModuleName string
}

func main() {
	lines, err := util.ReadInputFile()
	util.PanicOnError(err)

	modules := make(map[string]Module, len(lines))
	var flipFlopNames, conjunctionNames []string

	for _, line := range lines {
		parts := strings.Split(line, "->")
		kindAndName, outputsStr := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		outputNames := strings.Split(outputsStr, ", ")

		kindAndNameFirstChar := kindAndName[:1]
		if kindAndNameFirstChar == moduleKindFlipFlop {
			name := kindAndName[1:]
			modules[name] = Module{
				kind:              moduleKindFlipFlop,
				name:              name,
				outputModuleNames: outputNames,
			}
			flipFlopNames = append(flipFlopNames, name)
		} else if kindAndNameFirstChar == moduleKindConjunction {
			name := kindAndName[1:]
			modules[name] = Module{
				kind:              moduleKindConjunction,
				name:              name,
				outputModuleNames: outputNames,
			}
			conjunctionNames = append(conjunctionNames, name)
		} else {
			// broadcaster
			modules[kindAndName] = Module{
				kind:              kindAndName,
				name:              kindAndName,
				outputModuleNames: outputNames,
			}
		}
	}
	fmt.Printf("%d modules are parsed\n", len(modules))

	for _, conjName := range conjunctionNames {
		var conjInputs []string
		for modName, module := range modules {
			if slices.Contains(module.outputModuleNames, conjName) {
				conjInputs = append(conjInputs, modName)
			}
		}

		conjModule := modules[conjName]
		conjModule.inputModuleNames = conjInputs
		modules[conjName] = conjModule
	}

	var rxModuleInputs []Module
	for _, module := range modules {
		if slices.Contains(module.outputModuleNames, "rx") {
			rxModuleInputs = append(rxModuleInputs, module)
		}
	}
	fmt.Printf("rx module has following input modules: %v\n", rxModuleInputs)

	rxModuleInput := rxModuleInputs[0]

	steps, err := detectPulse(rxModuleInput.name, pulseKindLow, "rx", 1, 1_000_000, flipFlopNames,
		conjunctionNames, modules)
	if err == nil {
		fmt.Printf("rx received a low pulse after %d button presses\n", steps[0])
		return
	}
	fmt.Println(err)

	var intervals []uint
	for _, inputName := range rxModuleInput.inputModuleNames {
		steps, err := detectPulse(inputName, pulseKindHigh, rxModuleInput.name, 5, 10_000_000,
			flipFlopNames, conjunctionNames, modules)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Module %s sent %d %s pulses on steps %v\n", inputName, len(steps),
				pulseKindHigh, steps)
			
			var stepIntervals []uint
			for i := uint(0); i < uint(len(steps))-1; i++ {
				stepIntervals = append(stepIntervals, steps[i+1] - steps[i])
			}

			allSame := true
			for i := uint(0); i < uint(len(stepIntervals))-1; i++ {
				if stepIntervals[i] != stepIntervals[i+1] {
					allSame = false
					break
				}
			}

			if allSame {
				fmt.Printf("All step intervals are the same - %d\n", stepIntervals[0])
				intervals = append(intervals, stepIntervals[0])
			}
		}
	}

	if len(intervals) == len(rxModuleInput.inputModuleNames) {
		fmt.Println("LCM for intervals:", lcm(intervals[0], intervals[1], intervals[2:]...))
	}
}

func detectPulse(
	sourceModuleName, pulseKind, targetModuleName string,
	times, iterations uint,
	flipFlopNames, conjunctionNames []string,
	modules map[string]Module,
) ([]uint, error) {
	flipFlopStates := make(map[string]bool, len(flipFlopNames))
	for _, name := range flipFlopNames {
		flipFlopStates[name] = stateOff
	}

	conjunctionStates := make(map[string]map[string]string, len(conjunctionNames))
	for _, name := range conjunctionNames {
		inputNames := modules[name].inputModuleNames
		inputStates := make(map[string]string, len(inputNames))
		for _, inputName := range inputNames {
			inputStates[inputName] = pulseKindLow
		}
		conjunctionStates[name] = inputStates
	}

	var pulseHappenedOnSteps []uint
	for i := uint(1); i <= iterations; i++ {
		pulseHappened := isPulseSentOnButtonPress(sourceModuleName, pulseKind, targetModuleName,
			flipFlopStates, conjunctionStates, modules)
		if pulseHappened {
			pulseHappenedOnSteps = append(pulseHappenedOnSteps, i)

			if uint(len(pulseHappenedOnSteps)) == times {
				return pulseHappenedOnSteps, nil
			}
		}
	}
	return pulseHappenedOnSteps, fmt.Errorf("Failed to detect %s pulse %s->%s %d times after %d iterations",
		pulseKind, sourceModuleName, targetModuleName, times, iterations)
}

func isPulseSentOnButtonPress(
	wantedSourceModuleName, wantedPulseKind, wantedTargetModuleName string,
	flipFlopStates map[string]bool,
	conjunctionStates map[string]map[string]string,
	modules map[string]Module,
) bool {
	pulsesToProcess := []Pulse{{
		sourceModuleName: moduleKindButton,
		kind:             pulseKindLow,
		targetModuleName: moduleKindBroadcaster,
	}}
	var isPulseSent bool

	for len(pulsesToProcess) > 0 {
		pulse := pulsesToProcess[0]
		pulsesToProcess = pulsesToProcess[1:]

		if pulse.sourceModuleName == wantedSourceModuleName && pulse.kind == wantedPulseKind &&
			pulse.targetModuleName == wantedTargetModuleName {
			isPulseSent = true
		}

		targetModule, found := modules[pulse.targetModuleName]
		if !found {
			// fmt.Printf("WARN: module %s isn't found\n", pulse.targetModuleName)
			continue
		}

		switch targetModule.kind {
		case moduleKindFlipFlop:
			if pulse.kind == pulseKindLow {
				prevState := flipFlopStates[pulse.targetModuleName]
				flipFlopStates[pulse.targetModuleName] = !prevState

				outputPulseKind := pulseKindHigh
				if prevState {
					outputPulseKind = pulseKindLow
				}

				for _, outputName := range targetModule.outputModuleNames {
					pulsesToProcess = append(pulsesToProcess, Pulse{
						sourceModuleName: pulse.targetModuleName,
						kind:             outputPulseKind,
						targetModuleName: outputName,
					})
				}
			}
		case moduleKindConjunction:
			state := conjunctionStates[pulse.targetModuleName]
			state[pulse.sourceModuleName] = pulse.kind

			allInputsSentHigh := true
			for _, prevPulse := range state {
				if prevPulse == pulseKindLow {
					allInputsSentHigh = false
					break
				}
			}

			var outputPulseKind = pulseKindHigh
			if allInputsSentHigh {
				outputPulseKind = pulseKindLow
			}

			for _, outputName := range targetModule.outputModuleNames {
				pulsesToProcess = append(pulsesToProcess, Pulse{
					sourceModuleName: pulse.targetModuleName,
					kind:             outputPulseKind,
					targetModuleName: outputName,
				})
			}
		case moduleKindBroadcaster:
			for _, outputName := range targetModule.outputModuleNames {
				pulsesToProcess = append(pulsesToProcess, Pulse{
					sourceModuleName: pulse.targetModuleName,
					kind:             pulse.kind,
					targetModuleName: outputName,
				})
			}
		default:
			panic(fmt.Errorf("Module %s of kind %s cannot be a pulse target", targetModule.name,
				targetModule.kind))
		}
	}

	return isPulseSent
}

func gcd(a, b uint) uint {
	for b != 0 {
		temp := b
		b = a % b
		a = temp
	}
	return a
}

func lcm(a, b uint, integers ...uint) uint {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}
