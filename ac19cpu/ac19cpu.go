package ac19cpu

import (
	"fmt"
	"os"
	"encoding/csv"
	"strconv"
)

type Machine struct {
	memory  []int
	Input   []int
	Output  []int
	Program []int
	execIx  int
	base    int
	Status  string
	/* NEW: iniatilized
	   END: properly ended with opCode 99
	   EXH: exhausted - ended by running out of commands 
	   ERR: unknown command
	   HLT: halted due to Output
	   WTI: waiting for Input
	*/
}

// direct manipulation of memory
func (m *Machine) ManMemory(adr, val int) {
	m.setMem(adr, val)
}

// add one value to FIFO Output
func (m *Machine) AddOutput(val int) {

	m.Output = append(m.Output, val)
}

// add multiple values to FIFO Input
func (m *Machine) AddInput(vals ...int) {

	m.Input = append(m.Input, vals...)
}

// consume Output value 
func (m *Machine) ConsumeOutput() (result int) {

	result = m.Output[0]
	m.Output = m.Output[1:]
	return
}

// consume Input value
func (m *Machine) ConsumeInput() (result int, success bool) {

	if len(m.Input) < 1 {
		return 0, false
	} 
	result   = m.Input[0]
	m.Input  = m.Input[1:]
	success  = true
	return
}

// loadProgram (copies Program into memory)
func (m *Machine) LoadProgram(prog []int) {

	m.Program = prog
	m.memory  = make([]int, len(prog))
	m.execIx  = 0
	copy(m.memory, prog)
	m.Status  = "NEW"
	m.base    = 0
}

// loadss program from csv file
func (m *Machine) LoadProgramFromCsv(name string) {

	file, _ := os.Open(name)
	defer file.Close()

	nums := []int{}
	numStrs, _ := csv.NewReader(file).ReadAll()

	for _, numStr := range numStrs[0] {
		nums = append(nums, atoi(numStr))
	}	
	m.LoadProgram(nums)
}

// reloads cached Program into memory
func (m *Machine) ResetProgram() {
	copy(m.memory, m.Program)
	m.execIx = 0
	m.Status = "NEW"
	m.base   = 0
}

// resets index
func (m *Machine) ResetIndex() {
	m.execIx = 0
}

// expands memory
func (m *Machine) extendMemory(to int) {
	newMem := make([]int, to+5) // to make room for a command ...
	m.memory = append(m.memory, newMem...)
}

// 
func (m *Machine) getMem(addr int) int {
	if addr >= len(m.memory) {
		m.extendMemory(addr + 5)
	}
	return m.memory[addr]
}

// 
func (m *Machine) setMem(addr, val int) {
	if addr >= len(m.memory) {
		m.extendMemory(addr + 5)
	}
	m.memory[addr] = val
}

// CPU
func (m *Machine) Execute(haltOnOutput bool) (code string) {

	log    := false // enable for logging
	cmd	   := 0
	modes  := 0
	res    := 0
	suc    := true
	// address mode for each parameter
	mode := []int{0,0,0,0}
	// address mode resolution by caculating the index of the relevant cell
	parIx  := []int{0,0,0,0}
	// the total len of each command (other than 99)
	pLen   := []int{0,4,4,2,2,3,3,4,4,2}

	for m.execIx < len(m.memory) {

		cmd = m.memory[m.execIx] % 100

		// END
		if cmd == 99 {
			m.Status = "END" 
			return m.Status
		}

		if log { fmt.Print(m.memory[m.execIx:m.execIx+pLen[cmd]]) }

		// goes through all parameters and resolves
		// them immediately according to address mode
		modes = m.getMem(m.execIx) / 100
		for i := 1; i < pLen[cmd]; i++ {
			mode[i]   = (modes % 10)
			modes     = modes / 10
			// param[i]  = m.getMem(m.execIx + i)
			switch mode[i] {
			case 1: // direct addressing
				parIx[i] = m.execIx + i
			case 0: // indirect addressing
				parIx[i] = m.getMem(m.execIx + i)
			case 2: // realtive indirect adressing
				parIx[i] = m.getMem(m.execIx + i) + m.base

			}
		}		
		if log { fmt.Printf(" modes %v", mode[1:pLen[cmd]]) } 

		// Core
		switch cmd {

		// ADD
		case 1:			
			m.setMem(parIx[3], m.getMem(parIx[1]) + m.getMem(parIx[2]))
			m.execIx += pLen[cmd] 
			if log { fmt.Printf(" -> value %v to cell %v", m.memory[parIx[3]], parIx[3]) }

		// MUL
		case 2:
			m.setMem(parIx[3], m.getMem(parIx[1]) * m.getMem(parIx[2]))
			m.execIx += pLen[cmd] 
			if log { fmt.Printf(" -> value %v to cell %v", m.memory[parIx[3]], parIx[3]) }

		// INP
		case 3:
			res, suc = m.ConsumeInput()
			if (!suc) {
				m.Status = "WTI"
				return m.Status
			}
			m.setMem(parIx[1], res)
			m.execIx += pLen[cmd] 
			if log { fmt.Printf(" -> Input %v to cell %v", res, parIx[1]) }
			//fmt.Printf(" exIx: %v base: %v", m.execIx, m.base)

		// OUT
		case 4:
			m.AddOutput(m.getMem(parIx[1]))
			m.execIx += pLen[cmd] 
			if log { fmt.Printf(" -> Output value %v", m.memory[parIx[1]]) }
			if (haltOnOutput) {
				m.Status = "HLT"
				return m.Status
			}
		
		// NE0
		case 5:
			m.execIx += pLen[cmd] 
			if m.getMem(parIx[1]) != 0 {
				m.execIx = m.getMem(parIx[2])
				if log { fmt.Printf(" exec index set to %v", m.memory[parIx[2]]) }
			} else { if log { fmt.Printf(" nothing happened") } }
		
		// EQ0
		case 6:
			m.execIx += pLen[cmd] 
			if m.getMem(parIx[1]) == 0 {
				m.execIx = m.getMem(parIx[2])
				if log { fmt.Printf(" exec index set to %v", m.memory[parIx[2]]) }
			} else { if log { fmt.Printf(" nothing happened") } }
		
		// LTN
		case 7:
			if m.getMem(parIx[1]) < m.getMem(parIx[2]) {
				m.setMem(parIx[3], 1)
				if log { fmt.Printf(" cell %v set to 1", parIx[3]) }
			} else {
				m.setMem(parIx[3], 0)
				if log { fmt.Printf(" cell %v set to 0", parIx[3]) }
			}
			m.execIx += pLen[cmd]
		
		// EQL
		case 8:
			if m.getMem(parIx[1]) == m.getMem(parIx[2]) {
				m.setMem(parIx[3], 1)
				if log { fmt.Printf(" cell %v set to 1", parIx[3]) }
			} else {
				m.setMem(parIx[3], 0)
				if log { fmt.Printf(" cell %v set to 0", parIx[3]) }
			}
			m.execIx += pLen[cmd]

		// SBS
		case 9:
			m.base += m.getMem(parIx[1])
			m.execIx += pLen[cmd]

		default:
			fmt.Printf("\nUnknown Command %v !!!\n", cmd)
			m.Status = "ERR"
			return m.Status
		}
		
		if log { fmt.Printf(" / Ix now %v\n", m.execIx) }
	}

	m.Status = "EXH"
	return m.Status
}

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}
