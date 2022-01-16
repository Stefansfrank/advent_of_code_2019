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

//main processor
func execute(memory []int) []int {

	execIx := 0
	cmd	   := 0

	for execIx < len(memory) {
		switch cmd = memory[execIx]; cmd {
		case 1:
			memory[memory[execIx+3]] = memory[memory[execIx+1]] + memory[memory[execIx + 2]]
		case 2:
			memory[memory[execIx+3]] = memory[memory[execIx+1]] * memory[memory[execIx + 2]]
		case 99:
			return memory
		}
		execIx += 4
	}

	return memory
}

// prime the memory
func prime(memory []int, noun, verb int) {
	memory[1] = noun
	memory[2] = verb
}

// return the seed
func seed(noun, verb int) int {
	return noun*100 + verb
}

// loop through seeds
func loop(memory []int) int {

	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			newMemory := make([]int, len(memory))
			copy(newMemory, memory)
			prime(newMemory, i, j)
			if execute(newMemory)[0] == 19690720 {
				return seed(i, j)
			}
			
		}
	}
	return 0
}

// MAIN ----
func main () {

	start := time.Now()

	memory := readCsvFile2Int("d2.input.txt")

	// make a copy for part 1 so the processor can change the memory
	cpmem  := make([]int, len(memory))
	copy(cpmem, memory)
	
	prime(cpmem, 12, 2)
	fmt.Printf("\nFirst cell: %v\n", execute(cpmem)[0])

	fmt.Printf("Seed: %v\n\n", loop(memory))
	fmt.Printf("Execution time: %v\n", time.Since(start))
}