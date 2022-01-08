package main

import (
	"fmt"
	"time"
	"ac19/ac19cpu"
)

// runs empty input once on all machines 
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

// executes one network tic by reading the outputs of all units
// feeding it to the correct inputs and execute 
// if it encounters target 255 it returns the second value of the packet
func networkStep(cpus []ac19cpu.Machine) (bool, int) {

	for i := 0; i < 50; i++ {

		for len(cpus[i].Output) > 0 {
			tgt := cpus[i].ConsumeOutput()
			if tgt == 255 {
				cpus[i].ConsumeOutput()
				return true, cpus[i].ConsumeOutput()
			}
			cpus[tgt].AddInput(cpus[i].ConsumeOutput(), cpus[i].ConsumeOutput())
		}
	}

	for i := 0; i < 50; i++ {
		if len(cpus[i].Input) == 0 {
			cpus[i].Input = []int{-1}
		}
		cpus[i].Execute(false)		
	}

	return false, 0
} 

// -------------------------------------- Main ------------------------------------------------------------
func main() {

	start := time.Now()

	cpus := networkInit("d23.input.txt")

	end := false
	res := 0
	for !end {
		fmt.Print(".")
		end, res = networkStep(cpus)
	}
	fmt.Printf("\nY value of first packet to 255: %v\n", res)

	fmt.Printf("Execution time: %v\n", time.Since(start))
}