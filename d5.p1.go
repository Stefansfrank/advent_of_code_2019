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
	return input[0], input//[1:]
}

//main processor
func execute(memory, input, output []int) []int {

	log    := false
	execIx := 0
	cmd	   := 0
	modes  := 0
	result := 0
	// address mode for each parameter
	direct := []bool{false, false, false, false}
	param  := []int{0,0,0,0}
	// the total len of each command (other than 99)
	pLen   := []int{0,4,4,2,2}	

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
		}

		if log { fmt.Printf(" modes %v params %v", direct[1:pLen[cmd]], param[1:pLen[cmd]]) }
		switch cmd {
		case 1:
			if direct[1] {
				result = param[1]
			} else {
				result = memory[param[1]]
			}

			if direct[2] {
				result += param[2]
			} else {
				result += memory[param[2]]
			}

			memory[param[3]] = result
			if log { fmt.Printf(" -> value %v to cell %v", result, param[3]) }
		case 2:
			if direct[1] {
				result = param[1]
			} else {
				result = memory[param[1]]
			}

			if direct[2] {
				result *= param[2]
			} else {
				result *= memory[param[2]]
			}

			memory[param[3]] = result
			if log { fmt.Printf(" -> value %v to cell %v", result, param[3]) }
		case 3:
			result, input = getInput(input)
			memory[param[1]] = result
			if log { fmt.Printf(" -> input %v to cell %v", result, param[1]) }
		case 4:
			if direct[1] {
				output = setOutput(param[1], output)
				if log { fmt.Printf(" -> output value %v", param[1]) }
			} else {
				output = setOutput(memory[param[1]], output)
				if log { fmt.Printf(" -> output %v from cell %v", memory[param[1]], param[1]) }
			}
		}
		execIx += pLen[cmd] 
		if log { fmt.Printf(" / Ix now %v\n", execIx) }
	}

	return output
}

// MAIN ----
func main () {

	start := time.Now()

	memory := readCsvFile2Int("d5.input.txt")

	input  := []int{1}
	output := []int{}
	output = execute(memory, input, output)

	fmt.Printf("Output: %v\n", output)
	fmt.Printf("Execution time: %v\n", time.Since(start))
}