package main

import (
	"fmt"
	"strconv"
	"os"
	"encoding/csv"
	"time"
)

// Utility Functions ----------------------------------------------------------------------

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

// Bytecode Machine ----------------------------------------------------------------------

type machine struct {
	memory  []int
	input   []int
	output  []int
	program []int
	execIx  int
	status  string
	/* NEW: iniatilized
	   END: properly ended with opCode 99
	   EXH: exhausted - ended by running out of commands 
	   ERR: unknown command
	   HLT: halted due to output
	   WTI: waiting for input
	*/
}

// add one value to FIFO output
func (m *machine) addOutput(val int) {

	m.output = append(m.output, val)
}

// add one value to FIFO output
func (m *machine) addInput(val int) {

	m.input = append(m.input, val)
}

// consume output value 
func (m *machine) consumeOutput() (result int) {

	result = m.output[0]
	m.output = m.output[1:]
	return
}

// consume input value
func (m *machine) consumeInput() (result int, success bool) {

	if len(m.input) < 1 {
		return 0, false
	} 
	result   = m.input[0]
	m.input  = m.input[1:]
	success  = true
	return
}

// loadProgram (copies program into memory)
func (m *machine) loadProgram(prog []int) {

	m.program = prog
	m.memory  = make([]int, len(prog))
	m.execIx  = 0
	copy(m.memory, prog)
	m.status  = "NEW"
}

// reloads cached program into memory
func (m *machine) resetProgram() {
	copy(m.memory, m.program)
	m.execIx = 0
	m.status = "NEW"
}

// resets index
func (m *machine) resetIndex() {
	m.execIx = 0
}

// CPU
func (m *machine) execute(haltOnOutput bool) (code string) {

	log    := false // enable for logging
	cmd	   := 0
	modes  := 0
	res    := 0
	suc    := true
	// address mode for each parameter
	direct := []bool{false, false, false, false}
	// the parameters as they are
	param  := []int{0,0,0,0}
	// the parameters with addressing resolution
	parVal  := []int{0,0,0,0}
	// the total len of each command (other than 99)
	pLen   := []int{0,4,4,2,2,3,3,4,4}

	for m.execIx < len(m.memory) {

		cmd = m.memory[m.execIx] % 100

		// END
		if cmd == 99 {
			m.status = "END" 
			return m.status
		}

		if log { fmt.Print(m.memory[m.execIx:m.execIx+pLen[cmd]]) }

		// goes through all parameters and resolves
		// them immediately according to address mode
		modes = m.memory[m.execIx] / 100
		for i := 1; i < pLen[cmd]; i++ {
			direct[i] = ((modes % 10) == 1)
			modes     = modes / 10
			param[i]  = m.memory[m.execIx + i]
			if direct[i] {
				parVal[i] = param[i]
			} else {
				parVal[i] = m.memory[param[i]]
			}
		}		
		if log { fmt.Printf(" modes %v params %v", direct[1:pLen[cmd]], param[1:pLen[cmd]]) }

		// Core
		switch cmd {

		// ADD
		case 1:			
			m.memory[param[3]] = parVal[1] + parVal[2]
			m.execIx += pLen[cmd] 
			if log { fmt.Printf(" -> value %v to cell %v", m.memory[param[3]], param[3]) }

		// MUL
		case 2:
			m.memory[param[3]] = parVal[1] * parVal[2]
			m.execIx += pLen[cmd] 
			if log { fmt.Printf(" -> value %v to cell %v", m.memory[param[3]], param[3]) }

		// INP
		case 3:
			res, suc = m.consumeInput()
			if (!suc) {
				m.status = "WTI"
				return m.status
			}
			m.memory[param[1]] = res
			m.execIx += pLen[cmd] 
			if log { fmt.Printf(" -> input %v to cell %v", m.memory[param[1]], param[1]) }

		// OUT
		case 4:
			m.addOutput(parVal[1])
			m.execIx += pLen[cmd] 
			if log { fmt.Printf(" -> output value %v", parVal[1]) }
			if (haltOnOutput) {
				m.status = "HLT"
				return m.status
			}
		
		// NE0
		case 5:
			m.execIx += pLen[cmd] 
			if parVal[1] != 0 {
				m.execIx = parVal[2]
				if log { fmt.Printf(" exec index set to %v", parVal[2]) }
			}
			if log { fmt.Printf(" nothing happened") }
		
		// EQ0
		case 6:
			m.execIx += pLen[cmd] 
			if parVal[1] == 0 {
				m.execIx = parVal[2]
				if log { fmt.Printf(" exec index set to %v", parVal[2]) }
			}
			if log { fmt.Printf(" nothing happened") }
		
		// LTN
		case 7:
			if parVal[1] < parVal[2] {
				m.memory[param[3]] = 1
				if log { fmt.Printf(" cell %v set to 1", param[3]) }
			} else {
				m.memory[param[3]] = 0
				if log { fmt.Printf(" cell %v set to 0", param[3]) }
			}
			m.execIx += pLen[cmd]
		
		// EQL
		case 8:
			if parVal[1] == parVal[2] {
				m.memory[param[3]] = 1
				if log { fmt.Printf(" cell %v set to 1", param[3]) }
			} else {
				m.memory[param[3]] = 0
				if log { fmt.Printf(" cell %v set to 0", param[3]) }
			}
			m.execIx += pLen[cmd]

		default:
			fmt.Printf("\nUnknown Command %v !!!\n", cmd)
			m.status = "ERR"
			return m.status
		}
		
		if log { fmt.Printf(" / Ix now %v\n", m.execIx) }
	}

	m.status = "EXH"
	return m.status
}

// Main control programs ---------------------------------------------------------------

// array of machine instances (global)
var mms []machine

// this is the function to be executed for each permutation
type execFun func([]int) int
func execPerm(sequence []int) int {

	// log the cycles for each sequence
	log    := false
	numM   := len(mms)
	prev   := 0

	// simulating the last machine output to be 0
	// in order to prime the first machine
	mms[numM-1].addOutput(0)
	
	// adding the phase input only once and 
	// reset all machines to the original program
	for i ,phase := range sequence {
		mms[i].input = []int{phase}
		mms[i].resetProgram()
	}
	
	lastStatus := "NEW"
	lastResult := 0

	if log { fmt.Printf("S%v\n", sequence) }
	for true {
		for i := 0; i < numM; i++ {

			if log { fmt.Printf("M%v",i) }

			// cyclic indexing and copying previous output
			if i == 0 {
				prev = numM - 1
			} else {
				prev = i - 1
			}
			mms[i].addInput(mms[prev].consumeOutput())
			if log { fmt.Printf("|I%v",mms[i].input) }
	
			lastStatus = mms[i].execute(true)
			if log { fmt.Printf("|ST:%v", lastStatus) }
			if log { fmt.Printf("|O%v(I%v)\n",mms[i].output,mms[i].input) }

			// if last machine, cache the output in case END is detected subsequently
			if i == numM - 1 {
				lastResult = mms[i].output[0]
			}

			// one of the machines properly stops
			if lastStatus == "END" {
				if log { fmt.Printf("End Result: %v\n", lastResult)}
				return lastResult
			}

		}
	}

	// this can never happen since the loop is endless unless END is encountered
	return 0

}

// The recursive Heap permutation algorithm
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

	program  := readCsvFile2Int("d7.input.txt")

	numM    := 5
	seq     := []int{9,7,8,5,6}

	mms = make([]machine, numM)	
	for i := 0; i < numM; i++ {
		mms[i].loadProgram(program)
	}

	fmt.Printf("Maximum thrust: %v\n", heapPermutation(seq, numM, execPerm))

	fmt.Printf("Execution time: %v\n", time.Since(start))
}