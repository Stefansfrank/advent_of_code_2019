package main

import (
	"fmt"
	"time"
	"os"
	"bufio"
)

// simple text read - no error handling ...
func readTxtFile (name string) (lines []string) {	
	file, _ := os.Open(name)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {		
		lines = append(lines, scanner.Text())
	}
	return
}

// initializes the map with the input lines
func initial(inp []string) [][]byte {
	mp := make([][]byte, 7)
	for i := 0; i<7; i++ {
		mp[i] = make([]byte, 7)
		if i > 0 && i < 6 {
			for j, c := range inp[i-1] {
				mp[i][j+1] = ('.' - byte(c)) / ('.' - '#')
			}
		}
	}
	return mp
}

// executes one minute in the simulation
func step(mp [][]byte) (nmp [][]byte) {
	nmp = make([][]byte, len(mp))
	nmp[0] = make([]byte, len(mp[0]))
	nmp[6] = make([]byte, len(mp[0]))
	for y := 1; y < 6; y++ {
		nmp[y] = make([]byte, len(mp[0]))
		for x := 1; x < 6; x ++ {
			sm := (mp[y+1][x] + mp[y-1][x] + mp[y][x+1] + mp[y][x-1])
			if sm == 1 || sm == 2 && mp[y][x] == 0 {
				nmp[y][x] = 1
			}
		}
	}
	return
}

// creates a linear string from the grid used as a hash 
// in order to easily check for repetition
func hash(mp [][]byte) (h string) {
	h = ""
	for y := 1; y < 6; y++ {
		for x := 1; x < 6; x ++ {
			h += fmt.Sprintf("%v",mp[y][x])
		}
	}
	return
}

// calculates the biodiversity rating
func biod(mp [][]byte) (res int) {
	for y := 5; y > 0; y -- {
		for x := 5; x > 0; x -- {
			res <<= 1
			res += int(mp[y][x])
		}
	}
	return
}

func main () {
	start := time.Now()

	inp := readTxtFile("d24.input.txt")
	mp  := initial(inp)
	hit := map[string]bool{ hash(mp):true }

	end := false
	for !end {
		mp  = step(mp)
		hs := hash(mp)
		if hit[hs] {
			end = true
			fmt.Println("Biodiversity upon repetition:", biod(mp))
		} else {
			hit[hs] = true
		}
	}

	fmt.Println("Execution time: ", time.Since(start))
}
