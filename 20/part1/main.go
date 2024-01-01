package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/efulmo/advent-of-code-2023/util"
)

const (
	pulseLow  = "low"
	pulseHigh = "high"

	kindButton      = "button"
	kindBroadcaster = "broadcaster"
	kindFlipFlop    = "%"
	kindConjunction = "&"

	stateOn  = true
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
		if kindAndNameFirstChar == kindFlipFlop {
			name := kindAndName[1:]
			modules[name] = Module{
				kind:              kindFlipFlop,
				name:              name,
				outputModuleNames: outputNames,
			}
			flipFlopNames = append(flipFlopNames, name)
		} else if kindAndNameFirstChar == kindConjunction {
			name := kindAndName[1:]
			modules[name] = Module{
				kind:              kindConjunction,
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

	flipFlopStates := make(map[string]bool, len(flipFlopNames))
	for _, name := range flipFlopNames {
		flipFlopStates[name] = stateOff
	}

	conjunctionStates := make(map[string]map[string]string, len(conjunctionNames))
	for _, name := range conjunctionNames {
		inputNames := modules[name].inputModuleNames
		inputStates := make(map[string]string, len(inputNames))
		for _, inputName := range inputNames {
			inputStates[inputName] = pulseLow
		}
		conjunctionStates[name] = inputStates
	}

	var lowPulseCount, highPulseCount uint
	for i := uint(0); i < 1000; i++ {
		lowPulses, highPulses := pressButton(flipFlopStates, conjunctionStates, modules)
		lowPulseCount += lowPulses
		highPulseCount += highPulses
	}
	fmt.Printf("%d low * %d high = %d\n", lowPulseCount, highPulseCount, lowPulseCount*highPulseCount)
}

func pressButton(
	flipFlopStates map[string]bool,
	conjunctionStates map[string]map[string]string,
	modules map[string]Module,
) (uint, uint) {
	pulsesToProcess := []Pulse{{
		sourceModuleName: kindButton,
		kind:             pulseLow,
		targetModuleName: kindBroadcaster,
	}}

	var lowPulses, highPulses uint
	for len(pulsesToProcess) > 0 {
		pulse := pulsesToProcess[0]
		pulsesToProcess = pulsesToProcess[1:]

		switch pulse.kind {
		case pulseLow:
			lowPulses++
		case pulseHigh:
			highPulses++
		default:
			panic(fmt.Errorf("Unknown pulse kind %s if received", pulse.kind))
		}

		targetModule, found := modules[pulse.targetModuleName]
		if !found {
			// fmt.Printf("WARN: module %s isn't found\n", pulse.targetModuleName)
			continue
		}

		switch targetModule.kind {
		case kindFlipFlop:
			if pulse.kind == pulseLow {
				prevState := flipFlopStates[pulse.targetModuleName]
				flipFlopStates[pulse.targetModuleName] = !prevState

				outputPulseKind := pulseHigh
				if prevState {
					outputPulseKind = pulseLow
				}

				for _, outputName := range targetModule.outputModuleNames {
					pulsesToProcess = append(pulsesToProcess, Pulse{
						sourceModuleName: pulse.targetModuleName,
						kind:             outputPulseKind,
						targetModuleName: outputName,
					})
				}
			}
		case kindConjunction:
			state := conjunctionStates[pulse.targetModuleName]
			state[pulse.sourceModuleName] = pulse.kind

			allInputsSentHigh := true
			for _, prevPulse := range state {
				if prevPulse == pulseLow {
					allInputsSentHigh = false
					break
				}
			}

			var outputPulseKind = pulseHigh
			if allInputsSentHigh {
				outputPulseKind = pulseLow
			}

			for _, outputName := range targetModule.outputModuleNames {
				pulsesToProcess = append(pulsesToProcess, Pulse{
					sourceModuleName: pulse.targetModuleName,
					kind:             outputPulseKind,
					targetModuleName: outputName,
				})
			}
		case kindBroadcaster:
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

	return lowPulses, highPulses
}
