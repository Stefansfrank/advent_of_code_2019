package main

import (
	"fmt"
	"strconv"
	"os"
	"encoding/csv"
	"time"
	"ac19/ac19cpu"
)

// Utility Functions ----------------------------------------------------------------------

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

// simple int min function (for clarity)
func min (v1, v2 int) int {
	if v2 < v1 { return v2 }
	return v1
}

// simple int max function (for clarity)
func max (v1, v2 int) int {
	if v2 > v1 { return v2 }
	return v1
}

// Paint robot -------------------------------------------------------------------------
type pRobot struct {
	x,y    int  				// position
	xmin, xmax, ymin, ymax int 	// max values
	dir    int  				// direction (0: up, 1: right, 2:down, 3:left)
	pntPnl map[string]int  		// map of painted panels

	cpu    ac19cpu.Machine		// instance of the ac19cpu
}

// load program into the cpu
func (r *pRobot) LoadProgram(program []int) {
	r.cpu.LoadProgram(program)
}

// start painting
func (r *pRobot) Go(p2 bool) {

	stat := "NEW"
	r.pntPnl = make(map[string]int)

	// start from white panel for part 2
	if p2 {
		r.pntPnl["[0,0]"] = 1
	}

	for stat != "END" {
		r.cpu.AddInput(r.pntPnl[r.enc(r.x,r.y)])
		stat = r.cpu.Execute(false)
		r.paint(r.cpu.ConsumeOutput())
		r.move(r.cpu.ConsumeOutput())
	}
}

// turns the robot and moves
func (r *pRobot) move(dir int) {

	// turn
	if dir == 1 {
		r.dir = (r.dir + 1) % 4
	} else {
		r.dir -= 1
		if r.dir == -1 { r.dir = 3 }
	}

	// move 
	switch r.dir {
	case 0: 
		r.y--
		r.ymin = min(r.y, r.ymin)
	case 1: 
		r.x++
		r.xmax = max(r.x, r.xmax)
	case 2: 
		r.y++
		r.ymax = max(r.y, r.ymax)
	case 3: 
		r.x--
		r.xmin = min(r.x, r.xmin)
	}
}

// paints the current panel
func (r *pRobot) paint(col int) {

	r.pntPnl[r.enc(r.x,r.y)] = col
}

// paints a picture of the artwork
func (r *pRobot) DumpResult() {
	fmt.Println()
	for y := r.ymin; y <= r.ymax; y++ {
		for x := r.xmin; x <= r.xmax; x++ {
			if r.pntPnl[r.enc(x,y)] == 0 {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

// encodes x,y into a unique string
func (r *pRobot) enc(x, y int) string {
	return fmt.Sprintf("[%v,%v]", x, y)
}

// counts how many panels have been touched
// not that it does not double count mutiple paints of the same panel
func (r *pRobot) Count() int {
	return len(r.pntPnl)
}

// Main control programs ---------------------------------------------------------------

// MAIN ----
func main () {

	start := time.Now()

	program  := readCsvFile2Int("d11.input.txt")

	pr := pRobot{}	
	pr.LoadProgram(program)
	pr.Go(false)
	fmt.Printf("\nPainted panels - first attempt: %v\n", pr.Count())

	pr = pRobot{}
	pr.LoadProgram(program)
	pr.Go(true)
	pr.DumpResult()

	fmt.Printf("Execution time: %v\n", time.Since(start))
}