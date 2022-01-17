package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
)

// no error handling ...
func readTxtFile(name string) (lines []string) {
	
	file, _ := os.Open(name)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {		
		lines = append(lines, scanner.Text())
	}

	return

}

// simple integer absolute function
func abs(v int) int {
	if v >= 0 { return v }
	return -v
}

// executes the repetition of the input code pattern
// but at the same time throws away everything below 
// offset as there are only zeros in the matrix below the diagonal
func truncRptCode(code []int, offset, nmRept int) []int {
	max   := len(code)*nmRept
	tCode := make([]int, max - offset)

	for i := range tCode {
		tCode[i] = code[(offset+i) % len(code)]
	}
	return tCode
}

// executes phases under the assumption that
// the offset is in the second half of the new total code 
// Reason: beyond half, the transformation matrix has a very simple design:
// 1) the diagonal is all 1
// 2) everything below the diagonal is 0
// 3) everything above the diagonal is 1
// Due to 2) => all digits below offset can be ignored
// Due to 3) => no patterns have to be computed and the sum can be simplified 
func findMessage(code []int, offset, nmRept, numPhs, numDgt int) []int {

	if offset < len(code)*nmRept/2 {
		fmt.Println("Warning: Result might be wrong since offset too low!")
	} 

	// throw away everything below offset
	tCode := truncRptCode(code, offset, nmRept)
	tcLen := len(tCode)

	// phase loop
	for ph := 0; ph < numPhs; ph++ {
		
		// due to the diagonal nature of the matrix I can reduce 
		// this to one loop	going from the last line back to the first
		// letting the sum build up (all multipliers are 1)
		sum := 0
		for i := tcLen-1; i > -1; i-- {
			sum += tCode[i]
			tCode[i] = abs(sum) % 10
		}
	}

	return tCode[:numDgt]
}

// parse string
func parse(sCode string) []int {
	nCode := make([]int, len(sCode))
	for i, s := range sCode {
		nCode[i] = int(s - '0')
	}
	return nCode
}

// parse back from []int into a string
func parseBack(code []int) string {
	sCode := ""
	for _, v := range code {
		sCode += fmt.Sprintf("%v", v)
	}
	return sCode
}

func main() {
	start := time.Now()
	code := parse(readTxtFile("d16.input.txt")[0])
	offset := 0
	for i := 0; i < 7; i++ {
		offset = offset * 10 + code[i]
	}
	fmt.Printf("\nPart 2 Result: %v\n\n", parseBack(findMessage(code, offset, 10000, 100, 8)))

	fmt.Printf("Execution time: %v\n", time.Since(start))
} 