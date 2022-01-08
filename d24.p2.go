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

// a new 5x5 map
func newMp() (mp [][]int) {
	mp  = make([][]int, 5)
	for y := 0; y < 5; y++ {
		mp[y] = make([]int, 5)
	}	
	return 
}

// create a new stack of maps with 'numLr' layers
func newStck(numLr int) (stck map[int][][]int) {
	stck = map[int][][]int{0:newMp()}
	for i := 1; i <= numLr; i++ {
		stck[i]  = newMp()
		stck[-i] = newMp()
	}
	return
}

// creates the initial map from input and adds two buffer layers in both directions
func initial(inp []string) (stck map[int][][]int) {

	mp := newMp()
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			mp[y][x] = int(('.' - byte(inp[y][x])) / ('.' - '#'))
		}
	}
	stck = map[int][][]int{0:mp, 1:newMp(), 2:newMp(), -1:newMp(), -2:newMp()}
	return 
}

// this is somewhat ugly and calculates the sum of neightbors
// including the neighbors from bordering layers.
func nSum(lr, x, y int, stck map[int][][]int) (sm int) {

	switch x {
	case 0: 
		sm += stck[lr-1][2][1]
		sm += stck[lr][y][1]
	case 1:
		if y == 2 {
			sm +=  stck[lr+1][0][0] + stck[lr+1][1][0] + stck[lr+1][2][0] + stck[lr+1][3][0] + stck[lr+1][4][0]
		} else {
			sm += stck[lr][y][2]
		}
		sm += stck[lr][y][0]
	case 2: 
		if y == 2 {
			return 0 
		} else {
			sm += stck[lr][y][1]
			sm += stck[lr][y][3]			
		}
	case 3:
		if y == 2 {
			sm +=  stck[lr+1][0][4] + stck[lr+1][1][4] + stck[lr+1][2][4] + stck[lr+1][3][4] + stck[lr+1][4][4]
		} else {
			sm += stck[lr][y][2]
		}
		sm += stck[lr][y][4]
	case 4:	
		sm += stck[lr-1][2][3]
		sm += stck[lr][y][3]		
	}

	switch y {
	case 0: 
		sm += stck[lr-1][1][2]
		sm += stck[lr][1][x]
	case 1:
		if x == 2 {
			sm +=  stck[lr+1][0][0] + stck[lr+1][0][1] + stck[lr+1][0][2] + stck[lr+1][0][3] + stck[lr+1][0][4]
		} else {
			sm += stck[lr][2][x]
		}
		sm += stck[lr][0][x]
	case 2: 
		sm += stck[lr][1][x]
		sm += stck[lr][3][x]			
	case 3:
		if x == 2 {
			sm +=  stck[lr+1][4][0] + stck[lr+1][4][1] + stck[lr+1][4][2] + stck[lr+1][4][3] + stck[lr+1][4][4]
		} else {
			sm += stck[lr][2][x]
		}
		sm += stck[lr][4][x]
	case 4:	
		sm += stck[lr-1][3][2]
		sm += stck[lr][3][x]		
	}

	return
}

// print for debugging
func dump(stck map[int][][]int)	{
	sz := len(stck)/2 - 2
	for ix := -sz; ix <=sz; ix++ {
		fmt.Println("Layer",ix)
		for y := 0; y < 5; y ++ {
			for x := 0; x < 5; x++ {
				if x == 2 && y == 2 {
					fmt.Print("?")
				} else {
					fmt.Printf("%c", '.' - byte(stck[ix][y][x]) * ('.' - '#'))
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

// iterates one minute
func step(stck map[int][][]int) (nstck map[int][][]int) {
	sz    := len(stck)/2
	nsz   := sz + 1
	nstck  = newStck(nsz)
	ns    := 0

	// central layer
	for y := 0; y < 5; y ++ {
		for x := 0; x < 5; x++ {
			ns = nSum(0, x, y, stck) 
			if ns == 1 || ns == 2 && stck[0][y][x] == 0 {
				nstck[0][y][x] = 1
			}
		}
	}

	// added layers
	for lr := 1; lr <= (sz - 1); lr++ { 
		for y := 0; y < 5; y ++ {
			for x := 0; x < 5; x++ {
				ns = nSum(lr, x, y, stck) 
				if ns == 1 || ns == 2 && stck[lr][y][x] == 0 {
					nstck[lr][y][x] = 1
				}
				ns = nSum(-lr, x, y, stck) 
				if ns == 1 || ns == 2 && stck[-lr][y][x] == 0 {
					nstck[-lr][y][x] = 1
				}
			}
		}
	}

	return
}

// counts bugs
func count(stck map[int][][]int) (sm int) {
	for _, lr := range(stck) {
		for _, ln := range(lr) {
			for _, c := range(ln) {
				sm += c
			}
		}
	}
	return

}

func main () {
	start := time.Now()

	inp  := readTxtFile("d24.input.txt")
	stck := initial(inp)	

	for i := 0; i < 200; i++ {
		stck = step(stck)		
	}

	fmt.Println("Bug count after 200 iterations:", count(stck))

	fmt.Println("Execution time: ", time.Since(start))
}
