package main

import (
	"fmt"
	"strconv"
	"os"
	"encoding/csv"
	"time"
)

// no error handling ...
func readCsvFile2Int (name string) (nums []int) {
	
	file, _ := os.Open(name)
	defer file.Close()

	numStrs, _ := csv.NewReader(file).ReadAll()

	for _, numStr := range numStrs[0] {
		nums = append(nums, atoi(numStr))
	}	

	return

}

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

// output
func setOutput(val int, output []int) []int {
	return append(output, val)
}

// input
func getInput(input []int) (int, []int) {
	return input[0], input[1:]
}

//main processor
func execute(memory, input, output []int) []int {

	log    := false // enable for logging
	execIx := 0
	cmd	   := 0
	modes  := 0
	result := 0
	// address mode for each parameter
	direct := []bool{false, false, false, false}
	// the parameters as they are
	param  := []int{0,0,0,0}
	// the parameters with addressing resolution
	parVal  := []int{0,0,0,0}
	// the total len of each command (other than 99)
	pLen   := []int{0,4,4,2,2,3,3,4,4}	

	// for i := 0; i < 20; i++ { // limited loop for logging
	for execIx < len(memory) {

		cmd = memory[execIx] % 100
		if cmd == 99 { return output }
		if log { fmt.Print(memory[execIx:execIx+pLen[cmd]]) }
		modes = memory[execIx] / 100

		// goes through all parameters and resolves
		// them immediately according to address mode
		for i := 1; i < pLen[cmd]; i++ {
			direct[i] = ((modes % 10) == 1)
			modes     = modes / 10
			param[i]  = memory[execIx + i]
			if direct[i] {
				parVal[i] = param[i]
			} else {
				parVal[i] = memory[param[i]]
			}
		}		
		if log { fmt.Printf(" modes %v params %v", direct[1:pLen[cmd]], param[1:pLen[cmd]]) }

		switch cmd {
		case 1:			
			memory[param[3]] = parVal[1] + parVal[2]
			execIx += pLen[cmd] 
			if log { fmt.Printf(" -> value %v to cell %v", result, param[3]) }

		case 2:
			memory[param[3]] = parVal[1] * parVal[2]
			execIx += pLen[cmd] 
			if log { fmt.Printf(" -> value %v to cell %v", result, param[3]) }

		case 3:
			result, input = getInput(input)
			memory[param[1]] = result
			execIx += pLen[cmd] 
			if log { fmt.Printf(" -> input %v to cell %v", result, param[1]) }

		case 4:
			output = setOutput(parVal[1], output)
			execIx += pLen[cmd] 
			if log { fmt.Printf(" -> output value %v", parVal[1]) }
		
		case 5:
			execIx += pLen[cmd] 
			if parVal[1] != 0 {
				execIx = parVal[2]
				if log { fmt.Printf(" exec index set to %v", parVal[2]) }
			}
			if log { fmt.Printf(" nothing happened") }
		
		case 6:
			execIx += pLen[cmd] 
			if parVal[1] == 0 {
				execIx = parVal[2]
				if log { fmt.Printf(" exec index set to %v", parVal[2]) }
			}
			if log { fmt.Printf(" nothing happened") }
		
		case 7:
			if parVal[1] < parVal[2] {
				memory[param[3]] = 1
				if log { fmt.Printf(" cell %v set to 1", param[3]) }
			} else {
				memory[param[3]] = 0
				if log { fmt.Printf(" cell %v set to 0", param[3]) }
			}
			execIx += pLen[cmd]
		
		case 8:
			if parVal[1] == parVal[2] {
				memory[param[3]] = 1
				if log { fmt.Printf(" cell %v set to 1", param[3]) }
			} else {
				memory[param[3]] = 0
				if log { fmt.Printf(" cell %v set to 0", param[3]) }
			}
			execIx += pLen[cmd]

		default:
			fmt.Printf("\nUnknown Command %v !!!\n", cmd)
		}
		
		if log { fmt.Printf(" / Ix now %v\n", execIx) }
	}

	return output
}

// MAIN ----
func main () {

	start := time.Now()

	memory := readCsvFile2Int("d5.input.txt")

	input  := []int{5}
	output := []int{}
	output = execute(memory, input, output)

	fmt.Printf("Output: %v\n", output)
	fmt.Printf("Execution time: %v\n", time.Since(start))
}