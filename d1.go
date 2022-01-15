package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
)

// no error handling ...
func readTxtFile2Int (name string) (nums []int) {
	
	file, _ := os.Open(name)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {		
		nums = append(nums, atoi(scanner.Text()))
	}

	return

}

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

// fuel computation
func calcFuel (module int, rec bool) (fuel int) {
	
	fuel = (module / 3) - 2

	if !rec { return }

	if fuel < 1 { return 0 }

	fuel += calcFuel (fuel, true)

	return 
}

// MAIN ----
func main () {

	start := time.Now()

	modules := readTxtFile2Int("d1.input.txt")
	//modules := []int{12, 14, 1969, 100756}
	
	totalFuel := 0
	for _, module := range modules {
		totalFuel += calcFuel(module, false)
	}
	fmt.Printf("\nTotal fuel need (Part 1): %v\n", totalFuel)

	totalFuel = 0
	for _, module := range modules {
		totalFuel += calcFuel(module, true)
	}
	fmt.Printf("Total fuel need (Part 2): %v\n\n", totalFuel)

	fmt.Printf("Execution time: %v\n", time.Since(start))
}