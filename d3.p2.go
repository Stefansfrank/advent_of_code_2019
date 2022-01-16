package main

import (
	"fmt"
	"strconv"
	"time"
	"os"
	"encoding/csv"
)

// no error handling ...
func readCsvFile (name string) [][]string {
	
	file, _ := os.Open(name)
	defer file.Close()

	dirs, _ := csv.NewReader(file).ReadAll()
	return dirs

}


// Integer abs
func abs (x int) int { 
	if x < 0 {
		return -x
	}
	return x
}

// Simpler atoi
func atoi (x string) int { 
	y, _ := strconv.Atoi(x)
	return y
}

// Move structure
type move struct {
	vert bool // is move vertical?
	from int  // where the move starts in the direction of the move (always smaller than to!)
	to   int  // where the move stops in the direction of the move
	cnst int  // the value of the constant dimension
	invt bool // inverted ("L" or "D" i.e. counting backwards)
	stps int  // the total steps at the beginning of move
}

// Detects an intersection between two moves
// if found, the new step total is returned
func intersect (mov, refMov move) (intersect bool, totalStps int) {
	if (mov.vert != refMov.vert) {

		intersect = between(mov.from, mov.to, refMov.cnst) && between(refMov.from, refMov.to, mov.cnst)
		if intersect {

			totalStps = mov.stps + refMov.stps

			// last partial move
			if mov.invt {
				totalStps += mov.to - refMov.cnst
			} else {
				totalStps += refMov.cnst - mov.from
			}

			// last partial refMove
			if refMov.invt {
				totalStps += refMov.to - mov.cnst
			} else {
				totalStps += mov.cnst - refMov.from
			}
		}

	} else {
		// not yet handling parallel overlap, assuming no intersection
		intersect = false
	}
	return
}

// Detects whether cnst is between from and to
// assume from smaller than to 
func between (from, to, cnst int) bool {
	return cnst >= from && cnst <= to
}

// parse the moves into fully qualified moves
// enforces from to always be smaller than to
func parseMove (curX, curY, curStps int, directive string) (newX, newY, newStps int, newMove move) {
	
	stps := atoi(directive[1:])

	switch string(directive[0]) { 
	case "R":
		newY 		 = curY
		newX         = curX + stps
		newMove.vert = false
		newMove.cnst = curY
		newMove.from = curX
		newMove.to   = newX
		newMove.stps = curStps
		newMove.invt = false
	case "L":
		newY 		 = curY
		newX         = curX - stps
		newMove.vert = false
		newMove.cnst = curY
		newMove.from = newX
		newMove.to   = curX
		newMove.stps = curStps
		newMove.invt = true		
	case "U":
		newY 		 = curY + stps
		newX         = curX 
		newMove.vert = true
		newMove.cnst = curX
		newMove.from = curY
		newMove.to   = newY
		newMove.stps = curStps
		newMove.invt = false
	case "D":
		newY 		 = curY - stps
		newX         = curX 
		newMove.vert = true
		newMove.cnst = curX
		newMove.from = newY
		newMove.to   = curY
		newMove.stps = curStps
		newMove.invt = true
	}
	newStps = curStps + stps
	return
}


func main () {
	
	start := time.Now()
	
	input := readCsvFile("d3.input.txt")

	curX := 0
	curY := 0
	refMoves := []move{}
	var curMov move
	curStps := 0

	// build move list for wire 1 for reference 
	for _, directive := range input[0] {
		curX, curY, curStps, curMov = parseMove(curX, curY, curStps, directive)
		refMoves = append(refMoves, curMov) 
	}

	// loop through second wire
	curX = 0
	curY = 0
	minStps := 999999
	curStps = 0

	for _, directive := range input[1] {
		curX, curY, curStps, curMov = parseMove(curX, curY, curStps, directive)

		// check for intersections 
		for _, refMov := range refMoves {
			cross, stps := intersect(curMov, refMov)
			if cross && stps < minStps {
				minStps = stps
			}
		}
	}

	fmt.Printf("\nSteps to first intersection: %v\n\n", minStps)
	fmt.Printf("Execution time: %v\n", time.Since(start))
}