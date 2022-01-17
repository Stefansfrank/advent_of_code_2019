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

// -------------------------------------- Repair Robot ----------------------------------------------------

// simple point struct
type pnt struct {
	x,y int
}

// delivers point objects next t the current location
// input is the direction (1-4)
func (p *pnt) next(direction int) pnt {
	switch direction {
	case 1:
		return pnt{x:p.x, y:p.y-1}
	case 2:
		return pnt{x:p.x, y:p.y+1}
	case 3:
		return pnt{x:p.x-1, y:p.y}
	case 4:
		return pnt{x:p.x+1, y:p.y}
	}
	return pnt{x:0, y:0}
}

// delivers the opposite direction number
// to the input number
func back(direction int) int {
	switch direction {
	case 1: return 2
	case 2: return 1
	case 3: return 4
	case 4: return 3
	}
	return 0
}

// the repair robot
type repair struct {
	loc     pnt              // current location
	cpu     ac19cpu.Machine  // brains
	locOxy  pnt 			 // loc of the oxygen
	cnt     int 			 // counter of moves since start
}

// kick off exploration
func (r *repair) startExploration() {

	r.loc     = pnt{x:0,y:0}
	r.explore()
}

// attempts to move into direction dir
// returns fals if there was a wall 
func (r *repair) move(dir int) (moved bool) {

	// new location to explore
	nLoc := r.loc.next(dir)

	// run CPU and update Map
	r.cpu.AddInput(dir)
	r.cpu.Execute(false)
	rMap[nLoc] = r.cpu.ConsumeOutput()

	// move if > 0 
	if rMap[nLoc] > 0 {
		moved = true
		r.loc = nLoc

		// mark oxygen if encountered
		if rMap[nLoc] == 2 {
			r.locOxy = nLoc
		}			
	}

	// keep track of moves
	r.cnt++
	return moved
}

// recursive exploration loop
func (r *repair) explore() {

	// looks in all four directions and lists
	// lists unknown map tiles in best
	unknown := []int{}
	for i := 1; i < 5; i++ {

		_, exists := rMap[r.loc.next(i)]
		if !exists { 
			unknown = append(unknown, i)
		} 
	}

	// kicks off further exploration 
	// for all unknown directions
	for _, ix := range unknown {
		moved := r.move(ix)

		// recursion
		if moved {
			r.explore()

			// needed to keep the cpu in sync
			// throughout recursion
			r.move(back(ix))
		}
	}
}

// one path through the labyrinth
type path struct {
	oxy  bool         // whether oxygen was found
	cLoc pnt          // the current location of the exploration
	vis  map[pnt]bool // the trail of visited locations
	trk  []pnt        // track of points
}

func (p *path) steps() int {
	return len(p.vis)-1
}
// this starts a path exploration from origin
// it returns the array fo found pathes
// if stopOxy is set, pathes that encounter Oxygen end there
func startPaths(origin pnt, stopOxy bool) []path {

	paths := []path{}
	avail := []int{}

	// detect feasible directions around origin
	for i := 1; i < 5; i++ {
		if rMap[origin.next(i)] > 0 {
			avail = append(avail, i)
		}
	}

	// create a path for each identified direction
	// and kick of recursive exploration
	for i := 0; i < len(avail); i++ {
			vis := make(map[pnt]bool)
			vis[origin] = true
			trk := []pnt{origin}
			paths = append(paths, path{cLoc:origin, vis:vis, trk:trk})
			paths, _ = buildPaths(paths, avail[i], len(paths)-1, 0, stopOxy)		
	}

	return paths
}

// recursive path building loop
// inputs are direction of the next move, the index of the current path
// and an execution counter (in order to limit recursion for debugging)
// the paths array needs to be handed in and is returned again
func buildPaths(paths []path, dir, ix, cnt int, stopOxy bool) ([]path, int) {

	/*if ix == 0 && len(paths[ix].vis) < 3 {
		dumpTrack(paths[ix])
		fmt.Println(paths[ix].vis)
		fmt.Println(dir, ix, cnt)
	} //*/

	// move to the next tile
	nLoc := paths[ix].cLoc.next(dir)
	paths[ix].cLoc = nLoc
	paths[ix].vis[nLoc] = true
	paths[ix].trk = append(paths[ix].trk, nLoc)
	cnt++

	// detect Oxygen
	if rMap[nLoc] == 2 {
		paths[ix].oxy = true
		if stopOxy {
			return paths, cnt
		}
	}

	// detect next possible directions
	avail := []int{}
	for i := 1; i < 5; i++ {

		// only directions that have no wall 
		// and are not yet visited 
		val := rMap[nLoc.next(i)]
		if val > 0 && !paths[ix].vis[nLoc.next(i)] {
			avail = append(avail, i)
		}
	}

	// dead end
	if len(avail) == 0 {
		return paths, cnt
	}

	// safety valve
	if cnt > 100000 { 
		fmt.Println("Limit of iterations exceeded !!")
		return paths, cnt 
	}

	// this starts at 1 so it runs only if there is more
	// then one available direction to go
	for i := 1; i < len(avail); i++ {

		// copy the visited map of this path
		// for use in the newly created path
		cpVis := make(map[pnt]bool)
		for p, v := range paths[ix].vis {
			cpVis[p] = v
		}
		cpTrk := make([]pnt, len(paths[ix].trk))
		copy(cpTrk, paths[ix].trk)

		// create new path and kick off recursive path building
		paths = append(paths, path{ cLoc: paths[ix].cLoc, vis:cpVis, trk:cpTrk})
		paths, cnt = buildPaths(paths, avail[i], len(paths)-1, cnt, stopOxy)
	} //*/

	// continue on this path recursively
	paths, cnt = buildPaths(paths, avail[0], ix, cnt, stopOxy)
	return paths, cnt
}

// prints out all paths in a readable form
func dumpPaths(paths []path) {
	for i, p := range paths {
		fmt.Printf("Path[%v]: [Steps: %v] [Oxygen: %v] [Location: %v]\n", i, p.steps(), p.oxy, p.cLoc)
	}
}

func dumpTrack(p path) {
	for _,p := range p.trk {
		fmt.Printf("%3v\n", p)
	}

}

// prints out all paths in a readable form
func dumpOxyPaths(paths []path) {
	for i, p := range paths {
		if (p.oxy) {
			fmt.Printf("Path[%v]: [Steps: %v] [Oxygen: %v] [Location: %v]\n", i, p.steps(), p.oxy, p.cLoc)
		}
	}
}

// dumps the map to the terminal
func dumpMap() {

	// detect range
	var xmin, xmax, ymin, ymax int
	for key := range rMap {
		xmin = min(xmin, key.x)
		xmax = max(xmax, key.x)
		ymin = min(ymin, key.y)
		ymax = max(ymax, key.y)
	}

	for y:= ymin; y <= ymax; y++ {
		for x:= xmin; x <= xmax; x++ {
			if x == 0 && y == 0 {
				fmt.Print("Z")
				continue
			}
			dot, exists := rMap[pnt{x:x, y:y}]
			if !exists { dot = -1}
			switch dot {
			case -1: fmt.Print(" ")
			case 0: fmt.Print("*")
			case 1: fmt.Print(".")
			case 2: fmt.Print("O")
			}
		}
		fmt.Print("\n")
	}
}

// -------------------------------------- Main ------------------------------------------------------------

var rMap map[pnt]int

func main() {

	start := time.Now()
	rMap   = make(map[pnt]int)

	r := repair{ cpu:ac19cpu.Machine{} }
	r.cpu.LoadProgramFromCsv("d15.input.txt")

	// Improvement potential:
	// currently recursive exploration stops whenever
	// no unknown fields are encountered. The shortest path
	// algorithm is very similar to exploration but only stops on
	// dead ends. Both could be combined into one algorithm where
	// it switches from explore mode to path finding mode if no 
	// unknown fields are encountered but viable continuation is there

	// Build Map (for both parts)
	r.startExploration()

	// Part 1
	paths := startPaths(pnt{x:0, y:0}, true)
	minStp := 10000000
	for _, path := range paths {
		if path.oxy {
			minStp = min(minStp, path.steps())
		}
	}
	fmt.Printf("\nMinimum steps to Oxygen: %v\n", minStp)

	// Part 2
	paths = startPaths(r.locOxy, false)
	maxStp := 0
	for _, path := range paths {
		maxStp = max(maxStp, path.steps())
	}
	fmt.Printf("Minutes (maximum steps in paths) to fill Oxygen: %v\n\n", maxStp)

	fmt.Printf("Execution time: %v\n", time.Since(start))
}