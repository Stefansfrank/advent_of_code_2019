package main

import (
	"fmt"
	"time"
	"strconv"
	"ac19/ac19cpu"
	"sort"
)

// ------------------------------- Helpers ---------------------------------

type pnt struct {
	x,y int
}

type rng struct {
	from,to int
}

func (p *pnt) next(d int) pnt {
	switch d {
	case 0:	return pnt{x:p.x, y:p.y-1}
	case 1:	return pnt{x:p.x+1, y:p.y}
	case 2: return pnt{x:p.x, y:p.y+1}
	default: return pnt{x:p.x-1, y:p.y}
	}
}

func back(d int) int {
	switch d {
	case 0: return 2
	case 1: return 3
	case 2: return 0
	default: return 1
	}
}

func left(d int) int {
	if (d == 0) { return 3 }
	return d-1
}

func right(d int) int {
	if (d == 3) { return 0 }
	return d+1
}

// converts a string into []int with ascii codes terminated with \n
func cmdInp(s string) []int {
	inp := []int{}
	for _, c := range s {
		inp = append(inp, int(c))
	}
	inp = append(inp, 10)
	return inp
}

// consumes all output and converts back from ascii
func getOut(cpu *ac19cpu.Machine) (s string) {
	s   = ""
	ln := len(cpu.Output)
	for i := 0; i < ln; i++ {
		c := cpu.ConsumeOutput()
		s = s + fmt.Sprintf("%c",c)
	}
	return 
}

// ---------------------------------- core ---------------------------------------

// parses the cpu output into an [][]int
// stored in the global car sMap
func buildMap(out []int) {

	// determine width and height
	// note: throws away uncompleted 
	// lines at the end if input not
	// complete
	sMap = make(map[pnt]int)
	wth  := 0
	for i, o := range out {
		if o == 10 {
			wth = i
			break
		}
	}
	hgt  := len(out)/(wth+1)
	x,y  := 0,0

	for _, t := range out {
		switch t {
		case 10:
			y++
			if y >= hgt {
				break
			}
			x = 0
		// matches '^' assuming known
		// initial state 
 		case 94:
 			vLoc = pnt{x:x,y:y}
			sMap[pnt{x:x,y:y}] = t
			x++
		default:
			sMap[pnt{x:x,y:y}] = t
			x++
		}
	}

	// save width and height within the map 
	sMap[pnt{x:-2,y:0}] = wth
	sMap[pnt{x:-2,y:1}] = hgt
}

// print the map
func dumpMap() {
	wth := sMap[pnt{x:-2,y:0}] 
	hgt := sMap[pnt{x:-2,y:1}] 

	fmt.Println()
	for y := 0; y < hgt; y++ {
		for x := 0; x < wth; x++ {
			fmt.Printf("%c", sMap[pnt{x:x,y:y}])
		}
		fmt.Printf("\n")		
	}
}

// alignment (part 1)
func alignment() int {
	wth := sMap[pnt{x:-2,y:0}] 
	hgt := sMap[pnt{x:-2,y:1}] 

	al := 0
	for y := 1; y < hgt -1; y++ {
		for x :=1; x < wth - 1; x++ {
			if sMap[pnt{x:x,y:y}] != 46 && cross(x, y) { al += x*y } 
		}
	}
	return al
} 

// detects crossing
func cross(x,y int) bool {
	return sMap[pnt{x:x-1,y:y}] != 46 && sMap[pnt{x:x+1,y:y}] != 46 &&
			sMap[pnt{x:x,y:y-1}] != 46 && sMap[pnt{x:x,y:y+1}] != 46
}

// Part 2 --- Finding the Path --------------------------------------

// detects if a point is empty
// returns true also for points
// outside the known map
func empty(p pnt) bool {
	return sMap[p] == 46 || sMap[p] == 0
}

// get shortest path for the concrete
// scaffolding given (starting with ^ at vLoc)
// (no T crossings)
func path() string {
	cur  := vLoc
	cDir := 3
	nDir := 0
	dStr := ""
	ret  := "L" // first move needed
	cnt  := 0

	for { 
		dStr, nDir = nextDir(cur, cDir)
		switch nDir - cDir {
		case -2:
			ret += "," + strconv.Itoa(cnt)
			return ret
		case 0:
			cnt++
			cur = cur.next(nDir)
		default:
			ret += "," + strconv.Itoa(cnt) + dStr
			cnt = 0
			cDir = nDir
		}
	}

	return ret
}

// returns next Dir for the scaffolding traversal rule:
// Always go straight if possible which works for 
// the specific scaffolding given 
func nextDir(c pnt, d int) (string, int) {
	if !empty(c.next(d)) {
		return "", d
	}
	if !empty(c.next(left(d))) {
		return ",L", left(d)
	}
	if !empty(c.next(right(d))) {
		return ",R", right(d)
	}
	return "E", -1 // END
}

// ---- finding the right splits for the path to optimize travel

// determines all possible sub patterns with a length between min and max
// when a certain fragmentation is already given (initial fragmentation is 0 to len(dir))
func patterns(dir string, frags []rng, min, max int) (pats []string) {
	pats = []string{}

	// loop through length of patterns
	for ln := max; ln >= min; ln-- {

		// loop through fragments
		for _, f := range frags {

			// loop through index for patterns
			for patI := 0; patI <= (f.to - f.from - ln); patI++ {
				pat := dir[f.from+patI:f.from+patI+ln]

				// throw away invalid patterns
				// pattern should not start or end with comma
				if pat[0] == ',' || pat[ln-1] == ',' { continue }

				// but pattern should be between commas
				if patI > 0 && dir[f.from + patI - 1] != ',' { continue }
				if f.from+patI+ln < len(dir) && dir[f.from + patI + ln] != ',' { continue }

				// see whether it already exists
				ex := false
				for _, p := range pats {
					if pat == p { ex = true}
				}

				if !ex { pats = append(pats, pat) }
			}
		}
	} 
	return
}

// determines all matches i.e patterns that can be found
// multiple times in the fragments provided
func matches(pats []string, dir string, frags []rng) (matches []match) {
	matches  = []match{}

	// loop through patterns
	for _, pat := range pats {

		ln  := len(pat)
		cnt := 0
		loc := []rng{}

		// loop through fragments
		for _, f := range frags {

			// loop through potential match index
			for matI := 0; matI <= (f.to - f.from - ln); matI++ {

				// matched
				if pat == (dir[f.from:f.to])[matI:matI+ln] {
					cnt++
					loc = append(loc, rng{from:(f.from+matI), to:(f.from+matI+ln)})
				}
			}
		}

		// more than one hit across fragments
		if cnt > 1 {

			// detect overlaps in ranges
			for i := 1; i < cnt; i++ {
				if loc[i].from - loc[i-1].from < ln {

					// cut loc
					if i == cnt {
						loc = loc[0:i]
					} else {
						loc = append(loc[0:i],loc[i+1:]...)
					}
					cnt--
					i--
				}
			}
		}

		// if still more than one match
		if cnt > 1 {
			matches = append(matches, match{cnt:cnt, loc:loc, ln:ln})
		}

	} 
	return
}

// builds a new fragmentation based on a subset of fragments given
func buildFragments(dir string, loc []rng) []rng {
	frags := []rng{}

	if loc[0].from > 0 {
		frags = append(frags, rng{from:0, to:loc[0].from})
	}

	for i := 0; i < len(loc) - 1; i++ {
		frags = append(frags, rng{from:loc[i].to, to:loc[i+1].from})
	}

	if loc[len(loc)-1].to < len(dir) {
		frags = append(frags, rng{from:loc[len(loc)-1].to, to:len(dir)})
	}

	return frags
}

// other helpers to support the main split finder below

// append a new range and sort all
func mergeLocs(la, lb []rng) (ll []rng) {
	ll = la
	ll = append(ll, lb...)

	sort.Slice(ll, func(i, j int) bool {
    	return ll[i].from < ll[j].from
	})
	return
}

// helper structure to represent a split
type split struct {
	ma, mb, mc match
	tot        []rng
}

// calculates how much is not covered but a given split (end goal -> 0)
func (s * split) remain(dir string) int {
	// the last two summands are to account for commas between matched ranges
	return len(dir) - (s.ma.cnt * s.ma.ln + s.mb.cnt * s.mb.ln + s.mc.cnt * s.mc.ln) - len(s.tot) + 1
}

// A given match
type match struct{
	loc []rng
	cnt int
	ln  int
}

func findSplits(dir string) split {

	frags   := []rng{}    // a list of fragments of dir to analyse
	pats    := []string{} // a list of patterns to match
	mats    := []match{}  // a list of matches
	splits  := []split{}  // an ongoing list of ranges

	// determine potential pattern A
	frags   = []rng{ rng{from:0, to:len(dir)} }
	pats    = patterns(dir, frags, 10, 20)
	mats    = matches(pats, dir, frags)
	for _, m := range mats {
		s := split{ ma:m, mb:match{}, mc: match{}, tot:m.loc }
		splits = append(splits, s)
	}

	// determine matching patterns B
	tmpSplits := splits
	for i, s := range tmpSplits {
		frags = buildFragments (dir, s.tot)
		pats  = patterns(dir, frags, 5, s.ma.ln)
		mats  = matches(pats, dir, frags)
		for j, mm := range mats {
			tot := mergeLocs(s.tot, mm.loc)
			if j == 0 {
				splits[i].mb  = mm
				splits[i].tot = tot
			} else {
				splits = append(splits, split{ ma:s.ma, mb:mm, mc: match{}, tot:tot })
			}
		}
	}

	// determine pattern C
	tmpSplits = splits
	for i, s := range tmpSplits {
		frags = buildFragments (dir, s.tot)
		pats  = patterns(dir, frags, 3, s.mb.ln)
		mats  = matches(pats, dir, frags)
		for j, mm := range mats {
			tot := mergeLocs(s.tot, mm.loc)
			if j == 0 {
				splits[i].mc  = mm
				splits[i].tot = tot
			} else {
				splits = append(splits, split{ ma:s.ma, mb:s.mb, mc:mm, tot:tot })
			}
		}
	}

	// kick out splits where no 3 ranges were found
	tmpSplits = []split{}
	for _, s := range splits {
		if s.mc.ln > 0 {
			tmpSplits = append(tmpSplits, s)
		} 
	}
	splits = tmpSplits

	// return the split with no remaining elements
	for _, s := range splits {
		if s.remain(dir) == 0 {
			return s
		}
	}

	return split{}
}

// converting the inherent sequencing of the sub-matches A,B,C into
// the format required by the cpu as a direction
func detTopPat(spl split) string {

	pat := ""
	lp:for _,r := range spl.tot {
		for _,r2 := range spl.ma.loc {
			if r.from == r2.from {
				pat += "A,"
				continue lp
			}
		}
		for _,r2 := range spl.mb.loc {
			if r.from == r2.from {
				pat += "B,"
				continue lp
			}
		}
		for _,r2 := range spl.mc.loc {
			if r.from == r2.from {
				pat += "C,"
				continue lp
			}
		}
	}
	return pat[0:len(pat)-1]
}


// -------------------------------------- Main ------------------------------------------------------------

var cpu  ac19cpu.Machine
var sMap map[pnt]int
var vLoc pnt

func main() {

	start := time.Now()

	// Part 1
	cpu.LoadProgramFromCsv("d17.input.txt")
	cpu.Execute(false)
	buildMap(cpu.Output)
	dumpMap()
	fmt.Printf("\nAlignment (part 1): %v\n", alignment())

	// Part 2:
	// Path determination 
	pth := path()

	// finding splits 
	spl := findSplits(pth)

	// Run the robot again
	cpu.LoadProgramFromCsv("d17.input.txt")
	cpu.ManMemory(0, 2)
	cpu.Execute(false)
	getOut(&cpu)

	// add the overall pattern scheme
	cpu.AddInput(cmdInp(detTopPat(spl))...)
	cpu.Execute(false)
	getOut(&cpu)

	// add the three individual patterns A, B, C
	cpu.AddInput(cmdInp(pth[spl.ma.loc[0].from:spl.ma.loc[0].to])...)
	cpu.Execute(false)
	getOut(&cpu)
	cpu.AddInput(cmdInp(pth[spl.mb.loc[0].from:spl.mb.loc[0].to])...)
	cpu.Execute(false)
	getOut(&cpu)
	cpu.AddInput(cmdInp(pth[spl.mc.loc[0].from:spl.mc.loc[0].to])...)
	cpu.Execute(false)
	getOut(&cpu)

	// no camera feed needed
	cpu.AddInput([]int{ 110, 10 }...)
	cpu.Execute(false)

	// and go 
	fmt.Println("After the robot walked the optimized path, it collected", cpu.Output[len(cpu.Output)-1], "dust\n")

	fmt.Printf("Execution time: %v\n", time.Since(start))
}