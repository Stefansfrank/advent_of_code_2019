// This does run long - 13 sec on my MB Pro

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

// a point coordinate pair
type pnt struct {
	x int
	y int
}

// a portal
type port struct {
	in  pnt
	out pnt
	nam string
}

var mz     [][]int       // the map of the maze (without any portals, just walls)
var portMp map[pnt]port  // the list of portals
var start  pnt  		 // start point
var end    pnt           // end point

// parsing into a [][]int map that uses 0 for walls, 1 for a path
// portal info is in a seperate structure portMp mapping a pnt to the portal
// also sets start and end (all global since I am lazy)
func parse(inp []string) {

	mz     = make([][]int, len(inp))
	portMp = map[pnt]port{}

	// basic map building
	// keeping a list of letters encountered
	lets := []pnt{}
	for y, line := range inp {
		mz[y] = make([]int, len(inp[0]))
		for x, char := range line {
			switch char {
			case ' ','#':
				mz[y][x] = 0
			case '.':
				mz[y][x] = 1
			default:
				mz[y][x] = int(char)
				lets = append(lets, pnt{x:x, y:y})
			}
		}
	}

	// go through letters and identify the list of portals
	ports  := map[string][]pnt{}
	portIx := ""
	for _, lp := range lets {

		// go through four directions
		for i := 0; i < 4; i++ {

			np := lp.stp(i)
			if out(np) {
				continue
			}
			// is there an open path next to this letter?
			if mz[np.y][np.x] == 1 {
				op := lp.stp((i+2) % 4)
				if mz[lp.y][lp.x] < mz[op.y][op.x] {
					portIx = fmt.Sprintf("%c%c", mz[lp.y][lp.x], mz[op.y][op.x])
				} else {
					portIx = fmt.Sprintf("%c%c", mz[op.y][op.x], mz[lp.y][lp.x])
				}
				ports[portIx] = append(ports[portIx], np)
				mz[op.y][op.x] = 0
				mz[lp.y][lp.x] = 0
			}
		}
	}

	// now identify start / end / potential jumps
	for ix, pts := range ports {
		if ix == "AA" {
			start = pts[0]
		} else if ix == "ZZ" {
			end   = pts[0]
		} else {
			var pr port
			if outer(pts[0]) {
				pr = port{in:pts[1], out:pts[0], nam:ix}
			} else {
				pr = port{in:pts[0], out:pts[1], nam:ix}
			}
			portMp[pts[0]] = pr
			portMp[pts[1]] = pr
		}
	}

	return
}

// determins whether the current point is on an outer side of the ring
// helping to decide whether a portal jump is possible and which level
// it leads to.
func outer(p pnt) bool {
	return p.x == 2 || p.y == 2 || p.x == len(mz[0]) - 3 || p.y == len(mz) - 3
}

// determins invalid coordinates
func out(p pnt) bool {
	return p.x < 0 || p.y < 0 || p.x >= len(mz[0]) || p.y >= len(mz)
}

// returns ccoridnates after a step in direction 'dir'
func (p *pnt) stp(dir int) pnt {
	switch dir {
	case 0:
		return pnt{x:p.x, y:p.y-1}
	case 1:
		return pnt{x:p.x+1, y:p.y}
	case 2:
		return pnt{x:p.x, y:p.y+1}
	case 3:
		return pnt{x:p.x-1, y:p.y}
	}
	return pnt{}
}

// -------------------- let's build possible pathes 
type path struct {
	ln  int
	cur pnt
	lvl int
	trc string
}

// whenever we jump to a new inner ring, we are adding a layer (in my terminology)
func addLayer(visd [][][]bool) [][][]bool {
	lvl := len(visd)
	visd = append(visd, make([][]bool, len(mz)))
	for i := 0; i < len(mz); i++ {
		visd[lvl][i] = make([]bool, len(mz[lvl]))
	}
	return visd	
}

// a breadth first path search
func shortestPath() {

	visd    := [][][]bool{}
	visd     = addLayer(visd)
	paths   := []path{ path{ln:0, cur:start, lvl:0} }
	visd[0][start.y][start.x] = true

	cont := true
	for cont {

		cont = false
		for pix, pth := range paths {

			// exit scenario (since I use breadth first in a graph
			// with identical edge length 1, the first scenario found
			// is the shortest) 
			if pth.cur == end && pth.lvl == 0 {
				fmt.Println("Length of solution: ", pth.ln)
				return
			}

			// potential flat next steps (no jumps)
			next := []pnt{}
			nLvl := []int{}
			for i := 0; i < 4; i++ {
				np := pth.cur.stp(i)

				if mz[np.y][np.x] == 0 || visd[pth.lvl][np.y][np.x] {
					continue
				}

				next = append(next, np)
				nLvl = append(nLvl, pth.lvl)
			}

			// check for a portal jump
			lnk, found := portMp[pth.cur]
			if found {

				// Reaching an inner portal
				if pth.cur == lnk.in {
					if pth.lvl == len(visd) - 1 {
						visd = addLayer(visd)
					}
					next = append(next, lnk.out)
					nLvl = append(nLvl, pth.lvl + 1)

				// an outer portal
				} else {
					if pth.lvl > 0 {
						next = append(next, lnk.in)
						nLvl = append(nLvl, pth.lvl - 1)					
					}
				}
			}

			// add all the newly identified pathes
			for i, nx := range next {
				cont = true
				if i == 0 {
					paths[pix].cur = nx
					paths[pix].ln += 1
					paths[pix].lvl = nLvl[i]
				} else {
					paths = append(paths, path{cur: nx, ln: pth.ln + 1, lvl: nLvl[i]})
				}
				visd[nLvl[i]][nx.y][nx.x] = true
			} 
		}
	}
}

func main() {
	startTime := time.Now()
	parse(readTxtFile("d20." + os.Args[1] + ".txt"))

	shortestPath()

	fmt.Printf("Execution time: %v\n", time.Since(startTime))
} 


