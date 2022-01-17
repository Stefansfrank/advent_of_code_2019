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
	base    int
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
	m.base    = 0
}

// reloads cached program into memory
func (m *machine) resetProgram() {
	copy(m.memory, m.program)
	m.execIx = 0
	m.status = "NEW"
	m.base   = 0
}

// resets index
func (m *machine) resetIndex() {
	m.execIx = 0
}

// expands memory
func (m *machine) extendMemory(to int) {
	newMem := make([]int, to+5) // to make room for a command ...
	m.memory = append(m.memory, newMem...)
}

// 
func (m *machine) getMem(addr int) int {
	if addr > len(m.memory) {
		m.extendMemory(addr + 5)
	}
	return m.memory[addr]
}

// 
func (m *machine) setMem(addr, val int) {
	if addr > len(m.memory) {
		m.extendMemory(addr + 5)
	}
	m.memory[addr] = val
}

// CPU
func (m *machine) execute(haltOnOutput bool) (code string) {

	log    := false // enable for logging
	cmd	   := 0
	modes  := 0
	res    := 0
	suc    := true
	// address mode for each parameter
	mode := []int{0,0,0,0}
	// the parameters as they are
	param  := []int{0,0,0,0}
	// the parameters with addressing resolution
	parVal  := []int{0,0,0,0} // if parameters indicate a value to be read
	parRef  := []int{0,0,0,0} // if parameters indicate a cell to be written to
	// the total len of each command (other than 99)
	pLen   := []int{0,4,4,2,2,3,3,4,4,2}

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
		modes = m.getMem(m.execIx) / 100
		for i := 1; i < pLen[cmd]; i++ {
			mode[i]   = (modes % 10)
			modes     = modes / 10
			param[i]  = m.getMem(m.execIx + i)
			switch mode[i] {
			case 1:
				parVal[i] = param[i]
				parRef[i] = m.execIx + i
			case 0:
				parVal[i] = m.getMem(param[i])
				parRef[i] = param[i]
			case 2:
				parVal[i] = m.getMem(m.base+param[i])
				parRef[i] = m.base+param[i]

			}
		}		
		if log { fmt.Printf(" modes %v params %v", mode[1:pLen[cmd]], param[1:pLen[cmd]]) }

		// Core
		switch cmd {

		// ADD
		case 1:			
			m.setMem(parRef[3], parVal[1] + parVal[2])
			m.execIx += pLen[cmd] 
			if log { fmt.Printf(" -> value %v to cell %v", m.memory[parRef[3]], parRef[3]) }

		// MUL
		case 2:
			m.setMem(parRef[3], parVal[1] * parVal[2])
			m.execIx += pLen[cmd] 
			if log { fmt.Printf(" -> value %v to cell %v", m.memory[parRef[3]], parRef[3]) }

		// INP
		case 3:
			res, suc = m.consumeInput()
			if (!suc) {
				m.status = "WTI"
				return m.status
			}
			m.setMem(parRef[1], res)
			m.execIx += pLen[cmd] 
			if log { fmt.Printf(" -> input %v to cell %v", res, parRef[1]) }
			//fmt.Printf(" exIx: %v base: %v", m.execIx, m.base)

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
				m.setMem(parRef[3], 1)
				if log { fmt.Printf(" cell %v set to 1", parRef[3]) }
			} else {
				m.setMem(parRef[3], 0)
				if log { fmt.Printf(" cell %v set to 0", parRef[3]) }
			}
			m.execIx += pLen[cmd]
		
		// EQL
		case 8:
			if parVal[1] == parVal[2] {
				m.setMem(parRef[3], 1)
				if log { fmt.Printf(" cell %v set to 1", parRef[3]) }
			} else {
				m.setMem(parRef[3], 0)
				if log { fmt.Printf(" cell %v set to 0", parRef[3]) }
			}
			m.execIx += pLen[cmd]

		// SBS
		case 9:
			m.base += parVal[1]
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

// MAIN ----
func main () {

	start := time.Now()

	program  := readCsvFile2Int("d9.input.txt")
	//program  := []int{109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99}

	mms := machine{}	
	mms.loadProgram(program)

	// 
	mms.addInput(1)
	mms.execute(false)
	fmt.Printf("Status: %v\n", mms.status)
	fmt.Printf("Result: %v\n", mms.consumeOutput())

	//
	mms.resetProgram() 
	mms.addInput(2)
	mms.execute(false)
	fmt.Printf("Status: %v\n", mms.status)
	fmt.Printf("Result: %v\n", mms.consumeOutput())

	fmt.Printf("Execution time: %v\n", time.Since(start))
}