package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
	"strconv"
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

// parse the maze into the mz array and build slices for keys and doors and set the origin
func parse(inp []string) {

	xmax   = len(inp[0])
	ymax   = len(inp)
	maxKey = byte(0)
	keys   = make(map[byte]pnt)
	doors  = make(map[byte]pnt)

	mz = make(map[pnt]byte)
	for x := 0; x < xmax; x++ {
		for y, in := range inp {
			ltr := byte(in[x])
			pt  := pnt{x:x, y:y}

			mz[pt] = ltr

			// keys
			tst1, key := pt.isKey()
			if tst1 {
				keys[key] = pt

				// maintain max key
				if key > maxKey {
					maxKey = key
				}

			// origin		
			} else if pt.isOrig() {
				keys[255] = pt

			// doors
			} else {
				tst2, dor := pt.isDoor()
				if tst2 {
					doors[dor] = pt
				}	
			}
		}
	} 
}

// this copies the trace variable of the path if a path is spawned
func copyTrace(trc map[pnt]bool) (ctrc map[pnt]bool) {
	ctrc = make(map[pnt]bool)
	for k, v := range trc {
		ctrc[k] = v 
	}
	return
}

// Navigational Tools ----------------------------------------------------
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

func (p *pnt) isWall() bool {
	return mz[*p] == '#'
}

func (p *pnt) isKey() (bool, byte) {
	return ((mz[*p] >= 'a') && (mz[*p] <= 'z')), mz[*p] - 'a'
}

func (p *pnt) isDoor() (bool, byte) {
	return ((mz[*p] >= 'A') && (mz[*p] <= 'Z')), mz[*p] - 'A'
}

func (p *pnt) isOrig() bool {
	return mz[*p] == '@'
}

// this returns the possible next points to go to
// and the direction pointing from each next point back 
func (p *pnt) nextStep(trc map[pnt]bool) (pnts []pnt) {
	pnts = []pnt{}
	for i := 0; i < 4; i++ {
		pt := p.stp(i)
		if trc[pt] {
			continue
		}
		if !pt.isWall() {
			pnts = append(pnts, pt)
		}
	}
	return
}

// Path Building ------------------------------------------------------------

// Basic Path struct (paths are ways to get from a location to any key)
type path struct {
	ln  int           	// length of the path
	org byte			// origin of the path (expressed in index of keys)
	des byte			// destination of the path (expressed in index of keys)
	drs int				// the doors encountered on the path (as bits in an int)
	cur pnt 			// the current location (not needed once the path is built)
	trc map[pnt]bool    // a trace of the locations covered in order to prevent circular movement
}

// function to determine all potential paths from key of index org to a key
func detectPaths(org byte) (prPaths []path) {
	paths  := []path{}
	curPth := 0

	// determin the first branches at the entry
	origin := keys[org]
	nPs    := origin.nextStep(make(map[pnt]bool))
	for _, nP := range nPs {
		trc := make(map[pnt]bool)
		trc[origin] = true
		trc[nP]     = true
		paths = append(paths, path{ln:1, org:org, des:255, drs:0, cur:nP, trc:trc})
	}

	for (curPth < len(paths)) {
		pth := paths[curPth] // careful, no writing access due to by-value !

		// skip if already terminated
		if pth.des != 255 {
			curPth += 1
			continue
		}

		// am I on a door?
		tst1, dor := pth.cur.isDoor()
		if tst1 {
			paths[curPth].drs += 1 << dor
		}

		// am I on a new key?
		tst2, key := pth.cur.isKey()
		if tst2 && key != org {

			// add an identical terminated path to the end and then later try to continue this to the next key
			trc := copyTrace(pth.trc)
			paths = append(paths, path{ln:pth.ln, org:pth.org, des:key, drs:pth.drs, cur:pth.cur, trc:trc})
		}

		// where next on this path?
		nPs = pth.cur.nextStep(pth.trc) 
	
		// branch?
		if len(nPs) > 1 {

			// set up new branches
			for i := 1; i < len(nPs); i++ {
				trc := copyTrace(pth.trc)
				trc[nPs[i]] = true
				paths = append(paths, path{ln:pth.ln+1, org:pth.org, des:255, drs:pth.drs, cur:nPs[i], trc:trc})
			}
		}

		// continue you this path (not using 'pth' for write access because of by-value reference in Go)
		if len(nPs) > 0 {

			paths[curPth].ln += 1
			paths[curPth].cur = nPs[0]
			paths[curPth].trc[nPs[0]] = true

		// stuck, go to next path w/o a valid destination on this
		} else {
			curPth += 1
		}
	}

	// Pruning of paths
	short := make(map[pnt]pnt) 	// this tracks the shortest path with the same destination and doors on the way
								// for the key, x contains the destination and y the doors on the way
								// fot the value, x contains the length of the shortest and y the index in []paths

	for i,p := range paths {
		key := pnt{x:int(p.des), y:p.drs}
		st  := short[key]
		if p.des != 255 && (st.x == 0 || p.ln < st.x) {
			short[key] = pnt{x:p.ln, y:i}
		}
	}

	prPaths = []path{}
	for _, v := range short {
		prPaths = append(prPaths, paths[v.y])
	}
	return
}

// Debugging
func dumpPGraph(pn int) {
	for i, pp := range pGraph {
		if len(pp) <= 0 || (pn != -1 && pn != i) {
			continue
		}

		fmt.Printf("--- %2v Paths from key %c ---------------- ", len(pp), i + 'a')
		for j := maxKey; true; j-- {
			fmt.Printf("%c", j + 'A')
			if j == byte(0) {
				break
			}
		}
		fmt.Println()
		for _, p := range pp {
			fmt.Printf("> Path to key %c of length %3d with doors " + dumpKD(p.drs, false) + "\n", p.des + 'a', p.ln)
		}
		fmt.Println()
	}
}

// Debugging
func dumpSGraph() {
	fmt.Printf("--- %2v Seqences ---------------------------- ", len(sGraph))
	for j := maxKey; true; j-- {
		fmt.Printf("%c", j + 'a')
		if j == byte(0) {
			break
		}
	}
	fmt.Println()
	for _, ss := range sGraph {
		fmt.Printf("> Sequence to key %c of length %4d with keys %0" + strconv.Itoa(int(maxKey)+1) + "b\n", ss.dest + 'a', ss.len, ss.keys)
	}
	fmt.Println()
}

// debugging
func dumpKD(key int, tp bool) (s string) {
	offs := byte('A')
	if tp {
		offs = byte('a')
	}
	for i := byte(0); i <= maxKey; i++ {
		if ((1 << i) & key) > 0 {
			s += fmt.Sprintf("%c", offs + i)
		} else {
			s += " "
		}
	}
	return
}

// Sequence building -----------------------------------------------------------------

type sequ struct {
	name  string
	keys  int
	dest  byte
	len   int
}

// used as an index/key to access the visited log
type visK struct {
	loc   byte
	keys  int
}

// this builds valid sequences of paths
func buildSequences () (short int) {

	visit := make(map[visK]int) // the visited log remembers the shortest way to get to 'loc' and collect keys 'keys'
								// it is used to prune superflous solutions

	// initial sequences from the maze entry
	sGraph = []sequ{}
	for _, pth := range pGraph[255] {
		if pth.drs == 0 {
			sGraph = append(sGraph, sequ{name:string(pth.des + 'a'), keys:(1 << pth.des), dest:pth.des, len:pth.ln})
		}
	}
	
	currSq := 0

	// loop through all sequences and try to complete them
	for (currSq < len(sGraph)) {
		allPaths := pGraph[sGraph[currSq].dest]
		paths    := []path{} 

		// building a list of potential next paths by excluding invalid ones
		for _, pth := range allPaths {

			// Excludes solutions visiting a key more than once
			if (1 << byte(pth.des)) & sGraph[currSq].keys != 0 {
				continue
			}

			// Excludes solutions with doors with no key yet
			if pth.drs - (sGraph[currSq].keys & pth.drs) != 0 {
				continue
			}

			// Excludes solutions with no advantage over existing solutions 
			ll := visit[visK{loc:pth.des, keys: sGraph[currSq].keys}]
			if ll != 0 && ll <= sGraph[currSq].len + pth.ln {
				continue
			}

			// since this path will now be added, mark the visitor log ...
			visit[visK{loc:pth.des, keys: sGraph[currSq].keys}] = pth.ln + sGraph[currSq].len
			paths = append(paths, pth)
		}

		// branch if more than one path is possible
		if len(paths) > 1 {
			for i := 1; i < len(paths); i++ {
				sGraph = append(sGraph, sequ{   name: sGraph[currSq].name + string(paths[i].des + 'a'),
												keys: sGraph[currSq].keys | (1 << paths[i].des),
												dest: paths[i].des,
												len:  (sGraph[currSq].len + paths[i].ln)})
			}
		}

		// continue on this path if possible
		if len(paths) > 0 {
			sGraph[currSq].name  = sGraph[currSq].name + string(paths[0].des + 'a')
			sGraph[currSq].keys  = sGraph[currSq].keys | (1 << paths[0].des)
			sGraph[currSq].dest  = paths[0].des
			sGraph[currSq].len  += paths[0].ln

		// stuck - terminate the path
		} else {
			currSq += 1
		}
	}

	// prune the resulting sequences
	for _, ss := range sGraph {

		// incomplete solutions are kicked out
		if ss.keys != (1 << (maxKey + 1)) - 1 { 
			continue
		}

		if short == 0 || ss.len < short {
			short = ss.len
		}
	}

	return short
}

// globals (too lazy to handle all the back and forth ...)
var mz         map[pnt]byte
var xmax, ymax int
var maxKey     byte
var keys       map[byte]pnt   // note that origin is saved in keys[byte(255)]
var doors      map[byte]pnt
var pGraph     [][]path
var sGraph     []sequ

func main() {
	start := time.Now()
	parse(readTxtFile("d18." + os.Args[1] + ".txt"))

	// fmt.Println(keys, doors, maxKey)
	pGraph      = make([][]path, 256)
	pGraph[255] = detectPaths(255)
	for i := byte(0); i <= maxKey; i++ {
		pGraph[int(i)] = detectPaths(i)		
	}

	fmt.Printf("\nShortest path collecting all keys is of length %v.\n\n", buildSequences())	

	fmt.Printf("Execution time: %v\n", time.Since(start))
} 