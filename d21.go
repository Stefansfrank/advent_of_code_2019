package main

import (
	"fmt"
	"time"
	"ac19/ac19cpu"
)


// Helpers to handle input and output for the CPU
func cmdInp(s string) (out []int) {
	out = []int{}
	for _, c := range s {
		out = append(out, int(c))
	}
	out = append(out, int('\n'))
	return 
}

func cmdOut(out []int) (s string) {
	s = ""
	for _, c := range out {
		s = s + fmt.Sprintf("%c",c)
	}
	return 
}

// Program Part 1 
// with the concept of jumping early if there is a landing spot 
func prgm1() (p []string) {
	p = []string{}
	p = append(p, "NOT C J")
	p = append(p, "AND D J")
	p = append(p, "NOT A T")
	p = append(p, "OR T J")
	p = append(p, "WALK")
	return
}

// Program Part 2
// this needs an additional check before jumping early in order
// to ensure a subsequeent jump will suceed (checking H)
func prgm2() (p []string) {
	p = []string{}
	p = append(p, "NOT C J")
	p = append(p, "NOT B T")
	p = append(p, "OR T J")
	p = append(p, "AND D J")
	p = append(p, "AND H J")
	p = append(p, "NOT A T")
	p = append(p, "OR T J")
	p = append(p, "RUN")
	return
}


// -------------------------------------- Main ------------------------------------------------------------
func main() {

	start := time.Now()

	cpu := ac19cpu.Machine{}
	cpu.LoadProgramFromCsv("d21.input.txt")

	prog := prgm1()
	fmt.Println("In:","")
	cpu.Execute(false)
	fmt.Println("Out:",cmdOut(cpu.Output))

	for _, code := range prog {
		cpu.Output = []int{}
		cpu.Input = cmdInp(code)
		fmt.Println("In:", code)
		cpu.Execute(false)
		fmt.Println("Out:", cmdOut(cpu.Output))
	}
	fmt.Println("Walk result: ", cpu.Output[len(cpu.Output)-1])

	cpu.ResetProgram()

	prog = prgm2()
	fmt.Println("In:","")
	cpu.Execute(false)
	fmt.Println("Out:",cmdOut(cpu.Output))

	for _, code := range prog {
		cpu.Output = []int{}
		cpu.Input = cmdInp(code)
		fmt.Println("In:", code)
		cpu.Execute(false)
		fmt.Println("Out:", cmdOut(cpu.Output))
	}
	fmt.Println("Run result: ", cpu.Output[len(cpu.Output)-1])

	fmt.Printf("Execution time: %v\n", time.Since(start))
}