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
}

// Detects an intersection between two moves
// if found, the distance to 0,0 is returned
func intersect (mov, refMov move) (intersect bool, distance int) {
	if (mov.vert != refMov.vert) {
		intersect = between(mov.from, mov.to, refMov.cnst) && between(refMov.from, refMov.to, mov.cnst)
		if intersect {
			distance = abs(mov.cnst) + abs(refMov.cnst)
		} 
	} else {
		// not yet handling parallel overlap, assuming no intersection
		intersect = false
	}
	return
}

// Detects whether cnst is between from and to
// assumes from smaller than to
func between (from, to, cnst int) bool {
	return cnst >= from && cnst <= to
}

// parse the moves into fully qualified moves
// enforces from to always be smaller than to
func parseMove (curX, curY int, directive string) (newX, newY int, newMove move) {
	switch string(directive[0]) { 
	case "R":
		newY 		 = curY
		newX         = curX + atoi(directive[1:])
		newMove.vert = false
		newMove.cnst = curY
		newMove.from = curX
		newMove.to   = newX
	case "L":
		newY 		 = curY
		newX         = curX - atoi(directive[1:])
		newMove.vert = false
		newMove.cnst = curY
		newMove.from = newX
		newMove.to   = curX
	case "U":
		newY 		 = curY + atoi(directive[1:])
		newX         = curX 
		newMove.vert = true
		newMove.cnst = curX
		newMove.from = curY
		newMove.to   = newY
	case "D":
		newY 		 = curY - atoi(directive[1:])
		newX         = curX 
		newMove.vert = true
		newMove.cnst = curX
		newMove.from = newY
		newMove.to   = curY
	}
	return
}


func main () {
	
	start := time.Now()
	
	input := readCsvFile("d3.input.txt")

	curX := 0
	curY := 0
	refMoves := []move{}
	var curMov move

	// build move list of first wire for reference 
	for _, directive := range input[0] {
		curX, curY, curMov = parseMove(curX, curY, directive)
		refMoves = append(refMoves, curMov) 
	}

	// loop through second wire and keep track of minimum distance intersect
	curX = 0
	curY = 0
	minDist := 999999
	for _, directive := range input[1] {
		curX, curY, curMov = parseMove(curX, curY, directive)

		// check for intersections 
		for _, refMov := range refMoves {
			cross, dist := intersect(curMov, refMov)
			if cross && dist < minDist {
				if dist > 0 {
					minDist = dist
				}
			}
		}
	}

	fmt.Println("\nFinal distance from port is:", minDist,"\n")
	fmt.Println("Execution time:", time.Since(start))
}