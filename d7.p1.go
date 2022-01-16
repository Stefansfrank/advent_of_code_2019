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
func setOutput(val int) {
	output = append(output, val)
	return
}

// input
func getInput() (result int) {

	result = input[0]
	input  = input[1:]
	return
}

// loadProgram (copies program into memory)
func loadProgram(prog []int) {

	memory  = make([]int, len(prog))
	copy(memory, prog)
	return
}

// CPU
func execute() {

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
		if cmd == 99 { return }
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

		// ADD
		case 1:			
			memory[param[3]] = parVal[1] + parVal[2]
			execIx += pLen[cmd] 
			if log { fmt.Printf(" -> value %v to cell %v", result, param[3]) }

		// MUL
		case 2:
			memory[param[3]] = parVal[1] * parVal[2]
			execIx += pLen[cmd] 
			if log { fmt.Printf(" -> value %v to cell %v", result, param[3]) }

		// INP
		case 3:
			memory[param[1]] = getInput()
			execIx += pLen[cmd] 
			if log { fmt.Printf(" -> input %v to cell %v", result, param[1]) }

		// OUT
		case 4:
			setOutput(parVal[1])
			execIx += pLen[cmd] 
			if log { fmt.Printf(" -> output value %v", parVal[1]) }
		
		// NE0
		case 5:
			execIx += pLen[cmd] 
			if parVal[1] != 0 {
				execIx = parVal[2]
				if log { fmt.Printf(" exec index set to %v", parVal[2]) }
			}
			if log { fmt.Printf(" nothing happened") }
		
		// EQ0
		case 6:
			execIx += pLen[cmd] 
			if parVal[1] == 0 {
				execIx = parVal[2]
				if log { fmt.Printf(" exec index set to %v", parVal[2]) }
			}
			if log { fmt.Printf(" nothing happened") }
		
		// LTN
		case 7:
			if parVal[1] < parVal[2] {
				memory[param[3]] = 1
				if log { fmt.Printf(" cell %v set to 1", param[3]) }
			} else {
				memory[param[3]] = 0
				if log { fmt.Printf(" cell %v set to 0", param[3]) }
			}
			execIx += pLen[cmd]
		
		// EQL
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

	return
}

var input   []int
var output  []int
var memory  []int
var program []int

// this is the function to be executed for each permutation
type execFun func([]int) int
func execPerm(sequence []int) int {

	output = []int{0}

	for _ ,phase := range sequence {

		input    = []int{phase, output[0]}
		output   = []int{}
		loadProgram(program)
		execute()
	}
	return output[0]

}

// that's the permutation algorithm
func heapPermutation(a []int, size int, exec execFun) (result int) {

	// end of recursion
	if size == 1 {
		result = exec(a)
	}

	// recursive pair switching
	for i := 0; i < size; i++ {
		tmp := heapPermutation(a, size-1, exec)
		if tmp > result { result = tmp }

		if size%2 == 1 {
			a[0], a[size-1] = a[size-1], a[0]
		} else {
			a[i], a[size-1] = a[size-1], a[i]
		}
	}

	return
}


// MAIN ----
func main () {

	start := time.Now()

	program  = readCsvFile2Int("d7.input.txt")

	fmt.Printf("Maximum thrust: %v\n", heapPermutation([]int{0,1,2,3,4}, 5, execPerm))

	fmt.Printf("Execution time: %v\n", time.Since(start))
}