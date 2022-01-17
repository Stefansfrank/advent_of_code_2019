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

// coefficient
func mult(row, col int) int {
	switch (col / row) % 4 {
	case 0: { return 0 }
	case 1: { return 1 }
	case 2: { return 0 }
	case 3: { return -1 }
	}
	return 0
}

// executes one transofrmation of the input code
func trans(code []int) []int {

	nCode := make([]int, len(code))
	for ix := range code {
		for i := range code {
			nCode[ix] += code[i] * mult(ix+1, i+1)
		} 
		nCode[ix] = abs(nCode[ix]) % 10
	}
	return nCode
}

// parse string
func parse(sCode string) []int {
	nCode := make([]int, len(sCode))
	for i, s := range sCode {
		nCode[i] = int(s - '0')
	}
	return nCode
}

// parse back
func parseBack(code []int, f,t int) string {
	sCode := ""
	for i, v := range code {
		if i == f && f != 0 { sCode += fmt.Sprintf("|")}
		if v > -1 {
			sCode += fmt.Sprintf("%v", v)
		} else {
			sCode += fmt.Sprintf("-")
		}
		if i == t && t != 0 { sCode += fmt.Sprintf("|")}
	}
	return sCode
}

func main() {
	start := time.Now()
	code := parse(readTxtFile("d16.input.txt")[0])

	for i := 0; i < 100; i++ {
		code = trans(code)
	}
	fmt.Printf("\nPart 1 Result: %v\n\n", code[0:8])
	fmt.Printf("Execution time: %v\n", time.Since(start))
} 