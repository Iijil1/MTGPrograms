package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type transition struct {
	symbol    int
	direction string
	state     string
}

func exit(c int) {
	switch c {
	case 1:
		fmt.Fprintln(os.Stderr, "TM in standard format required.")
		os.Exit(1)
	case 2:
		fmt.Fprintln(os.Stderr, "Initial Tape is invalid")
		os.Exit(2)
	}
}

func main() {

	//get input
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	tmString := scanner.Text()
	if len(tmString) == 0 {
		exit(1)
	}
	scanner.Scan()
	initialTape := scanner.Text()

	//interpret input
	stateStrings := strings.Split(tmString, "_")
	numSymbols := len(stateStrings[0]) / 3
	symbols := []int{}
	for i := 0; i < numSymbols; i++ {
		symbols = append(symbols, i)
	}
	states := []string{}
	tm := map[string][]*transition{}
	for i, stateString := range stateStrings {
		if len(stateString) != 3*numSymbols {
			exit(1)
		}
		state := string(rune('A' + i))
		states = append(states, state)
		tm[state] = []*transition{}

		for i := 0; i < numSymbols; i++ {
			newSymbol := int(stateString[3*i] - '0')
			newDirection := string(stateString[3*i+1])
			newState := string(stateString[3*i+2])
			if newState < string('A') || newState >= string(rune('A'+len(stateStrings))) {
				tm[state] = append(tm[state], nil)
				continue
			}
			if newSymbol < 0 || newSymbol >= numSymbols {
				exit(1)
			}
			switch newDirection {
			case "L":
				tm[state] = append(tm[state], &transition{newSymbol, "Left", newState})
			case "R":
				tm[state] = append(tm[state], &transition{newSymbol, "Right", newState})
			default:
				exit(1)
			}
		}
	}
	startState := states[0]
	startSymbol := symbols[0]
	tapeInt := 0
	if len(initialTape) > 0 {
		startSymbol = int(initialTape[0] - '0')
		if startSymbol < 0 || startSymbol >= numSymbols {
			exit(2)
		}
		initialTape = initialTape[1:]
		for len(initialTape) > 0 {
			nextSymbol := int(initialTape[len(initialTape)-1] - '0')
			initialTape = initialTape[:len(initialTape)-1]
			if nextSymbol < 0 || nextSymbol >= numSymbols {
				exit(2)
			}
			tapeInt *= numSymbols
			tapeInt += nextSymbol
		}
	}

	//build TWM program
	dirs := []string{"Left", "Right"}

	//--name the clocks
	// clocks := []string{"Input"}
	// for _, dir := range dirs {
	// 	clocks = append(clocks, dir+"Tape")
	// 	clocks = append(clocks, dir+"Temp")
	// 	clocks = append(clocks, dir+"Trans")
	// 	clocks = append(clocks, dir+"Mult")
	// 	for _, sy := range symbols {
	// 		clocks = append(clocks, dir+"Div"+fmt.Sprint(sy))
	// 		clocks = append(clocks, dir+"Write"+fmt.Sprint(sy))
	// 	}
	// }
	// for _, st := range states {
	// 	for _, sy := range symbols {
	// 		clocks = append(clocks, "Trans"+st+fmt.Sprint(sy))
	// 	}
	// }
	clocks := []string{"Input", "LeftTape", "LeftTemp", "LeftTrans"}
	for _, sy := range symbols {
		clocks = append(clocks, "LeftDiv"+fmt.Sprint(sy))
	}
	clocks = append(clocks, "LeftMult")
	for _, sy := range symbols {
		clocks = append(clocks, "LeftWrite"+fmt.Sprint(sy))
	}
	for _, st := range states {
		for _, sy := range symbols {
			clocks = append(clocks, "Trans"+st+fmt.Sprint(sy))
		}
	}
	for _, sy := range symbols {
		clocks = append(clocks, "RightWrite"+fmt.Sprint(sy))
	}
	clocks = append(clocks, "RightMult")
	for _, sy := range symbols {
		clocks = append(clocks, "RightDiv"+fmt.Sprint(sy))
	}
	clocks = append(clocks, "RightTrans", "RightTemp", "RightTape")

	program := map[string]map[string]int{}

	//--fill Header
	program["Input"] = map[string]int{"Input": len(clocks)}
	for _, clock := range clocks[1:] {
		program["Input"][clock] = len(clocks) - 1
	}
	for i, dir := range dirs {
		//--fill Tape
		program[dir+"Tape"] = map[string]int{"Input": 2}
		for _, clock := range clocks[1:] {
			program[dir+"Tape"][clock] = 0
		}
		program[dir+"Tape"][dir+"Tape"] = 10
		program[dir+"Tape"][dir+"Mult"] = 9
		for _, sy := range symbols {
			program[dir+"Tape"][dir+"Div"+fmt.Sprint(sy)] = 9
		}
		//--fill Temp
		program[dir+"Temp"] = map[string]int{"Input": 2}
		for _, clock := range clocks[1:] {
			program[dir+"Temp"][clock] = 0
		}
		program[dir+"Temp"][dir+"Temp"] = 8
		program[dir+"Temp"][dir+"Trans"] = 7
		//--fill Trans
		program[dir+"Trans"] = map[string]int{"Input": 2}
		for _, clock := range clocks[1:] {
			program[dir+"Trans"][clock] = 2
		}
		program[dir+"Trans"][dir+"Tape"] = 4
		program[dir+"Trans"][dir+"Temp"] = 0
		//--fill Mult
		program[dir+"Mult"] = map[string]int{"Input": 2}
		for _, clock := range clocks[1:] {
			program[dir+"Mult"][clock] = 2
		}
		program[dir+"Mult"][dir+"Tape"] = 0
		program[dir+"Mult"][dir+"Temp"] = 2 + 2*numSymbols
		//--fill Divs
		for _, rowSy := range symbols {
			inputAdjustment := 0
			if rowSy != startSymbol {
				inputAdjustment = 1
			}
			program[dir+"Div"+fmt.Sprint(rowSy)] = map[string]int{"Input": 2 + inputAdjustment}
			for _, clock := range clocks[1:] {
				program[dir+"Div"+fmt.Sprint(rowSy)][clock] = 2
			}
			program[dir+"Div"+fmt.Sprint(rowSy)][dir+"Tape"] = 0
			program[dir+"Div"+fmt.Sprint(rowSy)][dir+"Temp"] = 2
			if rowSy == numSymbols-1 {
				program[dir+"Div"+fmt.Sprint(rowSy)][dir+"Temp"] += 2
			}
			for _, colDir := range dirs {
				program[dir+"Div"+fmt.Sprint(rowSy)][colDir+"Div"+fmt.Sprint(rowSy)] = 3
			}
			for _, state := range states {
				program[dir+"Div"+fmt.Sprint(rowSy)]["Trans"+state+fmt.Sprint(rowSy)] = 3
			}
			for _, colDir := range dirs {
				program[dir+"Div"+fmt.Sprint(rowSy)][colDir+"Div"+fmt.Sprint((rowSy+1)%numSymbols)] = 1
			}
			for _, state := range states {
				program[dir+"Div"+fmt.Sprint(rowSy)]["Trans"+state+fmt.Sprint((rowSy+1)%numSymbols)] = 1
			}
		}
		//--fill Writes
		for _, rowSy := range symbols {
			program[dir+"Write"+fmt.Sprint(rowSy)] = map[string]int{"Input": 2}
			for _, clock := range clocks[1:] {
				program[dir+"Write"+fmt.Sprint(rowSy)][clock] = 9
			}
			program[dir+"Write"+fmt.Sprint(rowSy)][dir+"Write"+fmt.Sprint(rowSy)] = 11
			program[dir+"Write"+fmt.Sprint(rowSy)][dirs[1-i]+"Mult"] = 0
			program[dir+"Write"+fmt.Sprint(rowSy)][dir+"Tape"] = 5 + 2*rowSy
			program[dir+"Write"+fmt.Sprint(rowSy)][dir+"Temp"] = 5
			program[dir+"Write"+fmt.Sprint(rowSy)][dir+"Trans"] = 5
			program[dir+"Write"+fmt.Sprint(rowSy)][dir+"Mult"] = 5
			for _, colSy := range symbols {
				program[dir+"Write"+fmt.Sprint(rowSy)][dir+"Div"+fmt.Sprint(colSy)] = 0
			}
		}
	}
	//--fill Transitions
	for _, rowState := range states {
		for _, rowSymbol := range symbols {
			inputAdjustment := 0
			if rowState != startState {
				inputAdjustment += 1
			}
			if rowSymbol != startSymbol {
				inputAdjustment += 1
			}
			program["Trans"+rowState+fmt.Sprint(rowSymbol)] = map[string]int{"Input": 1 + inputAdjustment}
			if tm[rowState][rowSymbol] == nil {
				for _, clock := range clocks[1:] {
					program["Trans"+rowState+fmt.Sprint(rowSymbol)][clock] = 0
				}
				continue
			}
			targetState := tm[rowState][rowSymbol].state
			targetSymbol := tm[rowState][rowSymbol].symbol
			dir := tm[rowState][rowSymbol].direction
			antidir := "Left"
			if dir == "Left" {
				antidir = "Right"
			}
			program["Trans"+rowState+fmt.Sprint(rowSymbol)][dir+"Tape"] = 0
			for _, colSymbol := range symbols {
				symbolAdjustment := 0
				if colSymbol == rowSymbol {
					symbolAdjustment += 1
				}
				if colSymbol == startSymbol {
					symbolAdjustment -= 1
				}
				program["Trans"+rowState+fmt.Sprint(rowSymbol)][dir+"Div"+fmt.Sprint(colSymbol)] =
					1 + symbolAdjustment
			}
			program["Trans"+rowState+fmt.Sprint(rowSymbol)][dir+"Temp"] = 2
			program["Trans"+rowState+fmt.Sprint(rowSymbol)][dir+"Trans"] = 3
			program["Trans"+rowState+fmt.Sprint(rowSymbol)][antidir+"Tape"] = 4
			program["Trans"+rowState+fmt.Sprint(rowSymbol)][antidir+"Mult"] = 5
			program["Trans"+rowState+fmt.Sprint(rowSymbol)][antidir+"Temp"] = 6
			program["Trans"+rowState+fmt.Sprint(rowSymbol)][antidir+"Trans"] = 7
			for _, colSymbol := range symbols {
				program["Trans"+rowState+fmt.Sprint(rowSymbol)][antidir+"Write"+fmt.Sprint(colSymbol)] = 10
				program["Trans"+rowState+fmt.Sprint(rowSymbol)][dir+"Write"+fmt.Sprint(colSymbol)] = 10
			}
			program["Trans"+rowState+fmt.Sprint(rowSymbol)][antidir+"Write"+fmt.Sprint(targetSymbol)] = 8
			correction := 0
			for _, colState := range states {
				stateAdj := 0
				if colState == targetState {
					stateAdj -= 1
				}
				if colState == rowState {
					stateAdj += 1
				}
				for _, colSymbol := range symbols {
					symbolAdjustment := 0
					if colSymbol == rowSymbol {
						symbolAdjustment += 1
					}
					if colSymbol == startSymbol {
						symbolAdjustment -= 1
					}
					program["Trans"+rowState+fmt.Sprint(rowSymbol)]["Trans"+colState+fmt.Sprint(colSymbol)] =
						stateAdj + symbolAdjustment
					if colState == targetState && colSymbol == symbols[0] {
						currentValue := -1
						if colState != rowState {
							currentValue += 1
						}
						if colSymbol != rowSymbol {
							currentValue += 1
						}
						correction = 9 - currentValue -
							program["Trans"+rowState+fmt.Sprint(rowSymbol)]["Trans"+colState+fmt.Sprint(colSymbol)]
					}
				}
			}
			//adjusting so the correct state triggers after the write
			for _, colState := range states {
				for _, colSymbol := range symbols {
					program["Trans"+rowState+fmt.Sprint(rowSymbol)]["Trans"+colState+fmt.Sprint(colSymbol)] += correction
				}
			}
			program["Trans"+rowState+fmt.Sprint(rowSymbol)][dir+"Mult"] = 10
			for _, colSymbol := range symbols {
				symbolAdjustment := 0
				if colSymbol == rowSymbol {
					symbolAdjustment += 1
				}
				if colSymbol == startSymbol {
					symbolAdjustment -= 1
				}
				program["Trans"+rowState+fmt.Sprint(rowSymbol)][antidir+"Div"+fmt.Sprint(colSymbol)] =
					10 + symbolAdjustment
			}
		}
	}
	//--add initial tape
	program["RightTape"]["Input"] += 2 * tapeInt

	//print TWM program
	programString := "["
	for _, row := range clocks {
		programString += "["
		for _, column := range clocks {
			programString += fmt.Sprint(program[row][column], ",")
		}
		programString = programString[:len(programString)-1]
		programString += fmt.Sprintf("],\n")
	}
	programString = programString[:len(programString)-2]
	programString += "]"
	fmt.Println(programString)
}
