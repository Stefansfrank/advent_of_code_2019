package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
	//"strconv"
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

// general point definition
type pnt struct {
	x int
	y int
}

// parse the maze into the mz array 
// mz contains the following ints
// 0 = wall, 1 = path, 2 = start, 3 = end, 10+ = portal jumps that are possible from the pnt
// the portal targets will be in the portals with the same index than in mz
func parse(inp []string) (mz map[pnt]int, jumps map[int]pnt, dim []int, start pnt) {

	mz    = map[pnt]int{}
	jumps = map[int]pnt{}
	dim   = []int{ len(inp[0]), len(inp)}

	// basic map building
	// keeping a list of letters encountered
	lets := []pnt{}
	for y, line := range inp {
		for x, char := range line {
			switch char {
			case ' ','#':
				mz[pnt{x:x, y:y}] = 0
			case '.':
				mz[pnt{x:x, y:y}] = 1
			default:
				mz[pnt{x:x, y:y}] = int(char)
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
			// is there an open path next to this letter?
			if mz[np] == 1 {
				op := lp.stp((i+2) % 4)
				if mz[lp] < mz[op] {
					portIx = fmt.Sprintf("%c%c", mz[lp], mz[op])
				} else {
					portIx = fmt.Sprintf("%c%c", mz[op], mz[lp])
				}
				ports[portIx] = append(ports[portIx], np)
				mz[op] = 0
				mz[lp] = 0
			}
		}
	}

	// now identify start / end / potential jumps
	jmpCnt := 10
	for ix, pts := range ports {
		if ix == "AA" {
			mz[pts[0]] = 2
			start      = pts[0]
		} else if ix == "ZZ" {
			mz[pts[0]] = 3
		} else {
			mz[pts[0]] = jmpCnt
			jumps[jmpCnt] = pts[1]
			jmpCnt += 1
			mz[pts[1]] = jmpCnt
			jumps[jmpCnt] = pts[0]
			jmpCnt += 1
		}
	}

	return
}

// loop through neighbours
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

// maze printing (debug)
func dump(mz map[pnt]int, dim []int) {
	for y := 0; y < dim[1]; y++ {
		for x := 0; x < dim[0]; x++ {
			cc := mz[pnt{x:x, y:y}]
			switch cc {
			case 0:
				fmt.Print("#")
			case 1:
				fmt.Print(".")
			case 2:
				fmt.Print("@")
			case 3:
				fmt.Print("$")
			default:
				fmt.Printf("%c", cc + 55)
			}
		}
		fmt.Print("\n")
	}
}

// Path Building ------------------------------------------------------------

// Basic Path struct (paths are ways to get from a location to any key)
// This is a depth first search with a cache of open paths
type path struct {
	ln  int           	// length of the path
	cur pnt             // curr end of path
	vld bool			// valid if it ended at ZZ
}

// function to determine all potential paths from key of index org to a key
func detectPaths(mz map[pnt]int, jumps map[int]pnt, dim []int, start pnt) (paths []path) {

	paths   = []path{}
	curP   := 0
	vis    := map[pnt]int{}

	// starting point
	paths      = []path{path{ln:0, cur: start}}
	vis[start] = -1

	// try to complete branches
	cnt := 0
	for (curP < len(paths) && cnt < 10000) {
		cnt += 1

		p := paths[curP]

		// reached the exit
		if mz[p.cur] == 3 {
			paths[curP].vld = true
			curP += 1
			continue
		}

		next := []pnt{}

		// where next on foot
		for i := 0; i < 4; i++ {
			np := p.cur.stp(i)
			if mz[np] > 0 && (vis[np] == 0 || (p.ln + 1) < vis[np]) {
				next = append(next, np)
			} 
		}

		// jump possible
		if mz[p.cur] > 9 {
			np := jumps[mz[p.cur]]
			if vis[np] == 0 || (p.ln + 1) < vis[np] {
				next = append(next, np)
			} 
		}

		// branch?
		if len(next) > 1 {
			// set up new branches
			for i := 1; i < len(next); i++ {
				paths = append(paths, path{ln:p.ln+1, cur:next[i]})
				vis[next[i]] = p.ln+1
			}
		}

		// continue you this path (not using 'p' for write access because of by-value reference in Go)
		if len(next) > 0 {
			paths[curP].ln += 1
			paths[curP].cur = next[0]
			vis[next[0]] = p.ln+1
	
		// stuck, go to next path w/o a valid destination on this
		} else {
			curP += 1
		}
	}
	return
}

func main() {
	start := time.Now()
	mz, jumps, dim, entry := parse(readTxtFile("d20." + os.Args[1] + ".txt"))

	paths := detectPaths(mz, jumps, dim, entry)
	for _,p := range paths {
		if p.vld {
			fmt.Println("Length: ", p.ln)
		}
	}

	fmt.Printf("Execution time: %v\n", time.Since(start))
} 