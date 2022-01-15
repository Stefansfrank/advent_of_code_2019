package main

import (
	"fmt"
	"time"
	"ac19/ac19cpu"
)

// the global map and initialization using intcode
var mp [][]int

func initMap() (cnt int) {
	mp = make([][]int, 100)
	for y:=0; y<100; y++ {
		mp[y] = make([]int, 100)
		for x:=0; x<100; x++ {
			mp[y][x] = compMp(x,y)
			if (x < 50 && y < 50) {
				cnt += mp[y][x]
			}
		}
	}
	return
}

func compMp(x,y int) (ret int) {
	cpu.Input = []int{x,y}
	cpu.Execute(true)
	ret = cpu.ConsumeOutput()
	cpu.ResetProgram()
	return	
}

// the core logic
// starts with a grid of 100x100 and check whether the top left point has 100 points below and to the right
// if not, one line is added to the right and the bottom of the map and the points immediatly next to the top left are checked
// they form a mirrored L which is now made larger and larger while the map itself always stays 100 points ahead of the mirror L
func iterate() (result int) {

	sol   := [][]int{}    // a list of points that have 100x100 points below / to the left
	src   := 0            // how far is the mirrored L from the top left?  
	end   := 10000        // an upper end

	// pushing the mirrored L of possible corner points out in this loop
	for (src < end) {

		// check for all the points on the mirrored L moving out from the top left
		// whether they have 100 points below them and to the right of them
		for i := 0; i <= src; i++ {

			// first part of the mirror L
			if check(i, src) {
				sol = append(sol, []int{i, src}) // found a solution
				if end == 10000 {
					end = src + 10  // if I found a solution, add 10 more iterations of the mirrored L
								   // in order to make sure the shortest is among them
				}
			}		

			// second part of the mirror L
			if (i != src && check(src, i)) {
				sol = append(sol, []int{src, i})
				if end == 10000 {
					end = src + 10
				}
			}		
		}

		// extend map by 1 to the right and down
		nMp := make([][]int, len(mp)+1)
		for i, mpl := range mp {
			nMp[i] = append(mpl, compMp(100 + src, i))
		}
		nMpl := []int{}
		for i := 0; i < len(nMp[0]); i++ {
			nMpl = append(nMpl, compMp(i, 100 + src))
		}
		nMp[len(nMp)-1] = nMpl
		mp = nMp
		src += 1

	}

	// computes the shortest of the solutions found
	ln     := -1
	for _, s := range sol {
		ln2 := s[0]*s[0] + s[1]*s[1]
		if ln == -1 || ln2 < ln {
			ln = ln2
			result = s[0]*10000 + s[1]
		}
	}

	return
}

// checks whether the point itself, the point 100 down and the point 100 left are in the beam
// if all three are in the beam, the whole ship is in the beam
func check(x, y int) bool {
	if mp[y][x] == 0 || mp[y+99][x] == 0 || mp[y][x+99] == 0 {
		return false
	}
	return true
}

// -------------------------------------- Main ------------------------------------------------------------
var cpu ac19cpu.Machine

func main() {

	start := time.Now()

	cpu = ac19cpu.Machine{}
	cpu.LoadProgramFromCsv("d19.input.txt")

	fmt.Println("\nNumber of points affected in 50x50 area:",initMap())
	fmt.Println("Closest point of full coverage: ", iterate(),"\n")

	fmt.Printf("Execution time: %v\n", time.Since(start))
}