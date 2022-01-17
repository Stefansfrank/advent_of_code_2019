package main

import (
	"fmt"
	"time"
	"ac19/ac19cpu"
)


// ------------------------------- Helpers ---------------------------------

// simle integer min function
func min(a, b int) int {
	if b < a { return b }
	return a
}

// simple integer max function
func max(a, b int) int {
	if b > a { return b }
	return a
}

// sign function (returns -1,0,1 depending on sign of value)
func sgn(a int) int {
	if a == 0 { 
		return 0 
	} 
	if a > 0 {
		return 1
	}
	return -1
}

// ------------------------------ Arcade -------------------------------------


// one point on the screen as struct
// used in order to use it as a hash in the map
type pnt struct {
	x,y int
}

type arcade struct {
	screen map[pnt]int          // screen (uses pnt structure as hash and contains object type)
	score  int                  // current score as returned by the machine
	cpu    ac19cpu.Machine      // instance of the intcode machine
	xmin, ymin, xmax, ymax int  // the borders of the screen
	xball, xpad  int            // x positions of the ball and the paddle after each iteration
}

// runs the arcade CPU one time (if parameter is true, prints extensive logging and the map)
func (a *arcade) run(log bool) (status string) {

	if log { fmt.Printf("Execute with input: %v\n", a.cpu.Input) }
	
	// executes the cpu
	status   = a.cpu.Execute(false)
	
	if log { fmt.Printf("Status: %v - Size of output: %v triples\n", status, len(a.cpu.Output)/3)}
	if log && len(a.cpu.Output)/3 < 13 {
		fmt.Printf("Output: %v\n", a.cpu.Output)
	}
	if log { fmt.Printf("Screen Update - ") }
	
	// converts the output
	a.updateScreen(log)
	if log { a.dumpScreen() }
	if log { fmt.Printf("\n")}
	return
}

// converts the output after a run
// updates the screen
// redetermines the screen limits every time (not optimized)
// also updates the ball and paddle positions if changed
func (a *arcade) updateScreen(log bool) {

	var x,y,v int

	for len(a.cpu.Output) > 0 {
		x = a.cpu.ConsumeOutput()
		y = a.cpu.ConsumeOutput()
		v = a.cpu.ConsumeOutput()

		// limit detection
		a.xmin = min(a.xmin, x)
		a.ymin = min(a.ymin, y)
		a.xmax = max(a.xmax, x)
		a.ymax = max(a.ymax, y)

		// score change?
		if x == -1 && y == 0 {
			a.score = v
			if log { fmt.Printf("Score to: %v ", a.score)}

		// everything else changes a screen cell
		} else {
			a.screen[pnt{x: x, y: y}] = v
		}

		// update x value of the ball if ball repainted
		if v == 4 {
			if log { fmt.Printf("XBall to: %v ", x)}
			a.xball = x

		// update x value of paddle if paddle repainted
		} else if v == 3 {
			if log { fmt.Printf("XPad to: %v ", x)}
			a.xpad  = x
		}
	}

	if log { fmt.Printf("\n") }

}

// ucomputes the total blocks still in play
func (a *arcade) countBlocks() (cnt int) {

	for _, v := range a.screen {
		if v == 2 {
			cnt++
		}
	}
	return
}

// paints the current screen in ASCII
func (a *arcade) dumpScreen() {

	syms := []string{" ", "*", "x", "-", "o", "E"}
	col  := 0

	for y := a.ymin; y <= a.ymax; y++ {
		for x := a.xmin; x <= a.xmax; x++ {
			col = a.screen[pnt{x:x, y:y}]
			if col > len(syms)-2 { col = len(syms)-1 }
			fmt.Print(syms[col])
		}
		fmt.Println("")
	}
	return
}

// -------------------------------------- Main ------------------------------------------------------------

func main() {

	start := time.Now()

	arc := arcade{ cpu:ac19cpu.Machine{}, screen:make(map[pnt]int) }
	arc.cpu.LoadProgramFromCsv("d13.input.txt")

	arc.run(false)
	fmt.Printf("\nBlocks after initial run: %v\n",arc.countBlocks())

	arc = arcade{ cpu:ac19cpu.Machine{}, screen:make(map[pnt]int) }
	arc.cpu.LoadProgramFromCsv("d13.input.txt")
	arc.cpu.ManMemory(0,2)

	status := "NEW" // detects whether the cpu ended waiting for input or ended
	ix     := 0     // iteration counter
	maxIx  := 10000 // max limit of iterations
	joy    := 0     // joystick movement

	for status != "END" {

		// run the machine
		status = arc.run(false)

		//determine joystick direction for next step
		joy = sgn(arc.xball - arc.xpad)
		arc.cpu.AddInput(joy)
		
		// max iteration
		ix++
		if ix > maxIx { break }
	}

	fmt.Printf("Score with joystick control: %v\n\n",arc.score)
	fmt.Printf("Execution time: %v\n", time.Since(start))
}