package main

import (
	"fmt"
	"time"
	"ac19/ac19cpu"
)

// the NAT memory
type nat struct {
	x, y int
}

// creates all units and executes the initial run
func networkInit(file string) (cpus []ac19cpu.Machine) {
	cpus = make([]ac19cpu.Machine, 50)
	for i := 0; i < 50; i++ {
		cpus[i] = ac19cpu.Machine{}
		cpus[i].LoadProgramFromCsv(file)
		cpus[i].Input = []int{i}
		cpus[i].Execute(false)
	}
	return
}

// executes one minute, detects idle and sends NAT memory to 0
// detects if there are two idles in a row with identical Y values
func networkStep(cpus []ac19cpu.Machine, cnat *nat, lastIdle bool, lastY0 int) (end, idle bool, Y0 int) {

	// collect output
	for i := 0; i < 50; i++ {
		for len(cpus[i].Output) > 0 {
			tgt := cpus[i].ConsumeOutput()
			if tgt == 255 {
				// fill NAT
				cnat.x = cpus[i].ConsumeOutput()
				cnat.y = cpus[i].ConsumeOutput()
			} else {
				cpus[tgt].AddInput(cpus[i].ConsumeOutput(), cpus[i].ConsumeOutput())
			}
		}
	}

	// determine idle state
	idlecs := 0
	for i := 0; i < 50; i++ {
		if len(cpus[i].Input) == 0 {
			cpus[i].Input = []int{-1}
			idlecs += 1
		}
	}

	// deal with idle state
	if idlecs == 50 {
		if lastIdle && cnat.y == lastY0 {
			fmt.Printf("\nSent %v to cpu 0 twice in a row\n", cnat.y)
			return true, true, cnat.y
		}
		cpus[0].AddInput(cnat.x, cnat.y)
		Y0 = cnat.y
		idle = true
	} else {
		idle = false
	}

	// executes units
	for i := 0; i < 50; i++ {
		cpus[i].Execute(false)		
	}

	return 
} 



// -------------------------------------- Main ------------------------------------------------------------
func main() {

	start := time.Now()

	cpus := networkInit("d23.input.txt")
	cnat := &nat{}

	end  := false
	idle := false
	Y0   := 0
	for !end {
		fmt.Print(".")
		end, idle, Y0 = networkStep(cpus, cnat, idle, Y0)
	}

	fmt.Printf("Execution time: %v\n", time.Since(start))
}