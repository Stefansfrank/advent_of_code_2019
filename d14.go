package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"strings"
	"sort"
)

// no error handling ...
func readTxtFile(name string) (reacts []string) {
	
	file, _ := os.Open(name)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {		
		reacts = append(reacts, scanner.Text())
	}

	return

}

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

// -------------------------------- parsing --------------------------------

// list of input elements and output element
type reaction struct {
	out element
	in  []element
}

type element struct {
	amt    int    // amount needed
	nam    string // name
}

func parseReacts(reaLns []string) () {

	// global reaction map since I am lazy
	reaMap = make(map[string]reaction)
	var rea reaction 
	var ix  int

	for _, reaLn := range reaLns {

		rea  = reaction{}

		// split line into before and after =>
		ix     = strings.Index(reaLn, " => ")
		rea.in  = parseElements(reaLn[0:ix])
		rea.out = parseElements(reaLn[ix+4:])[0]
		reaMap[rea.out.nam] = rea
	}

	return
}

// parse the trimmed element lists before or after =>
func parseElements(elemLn string) (elems []element) {

	elems  = []element{}
	ix    := strings.Index(elemLn, ", ")
	ix2   := 0
	tmp   := ""

	// looping through found commas
	for ix > -1 {
		tmp = elemLn[0:ix]
		ix2 = strings.Index(tmp, " ")
		elems = append(elems, element{ amt:atoi(tmp[0:ix2]), nam:tmp[ix2+1:] })
		elemLn = elemLn[ix+2:]
		ix  = strings.Index(elemLn, ", ")	
	}

	// last value after the comma
	tmp = elemLn[0:]
	ix2 = strings.Index(tmp, " ")
	elems = append(elems, element{ amt:atoi(tmp[0:ix2]), nam:tmp[ix2+1:] })

	return
}

//  --------------------------------------------------------------------------

// calculates the amount of ORE to produce elmAmt of elmName
func calcNeed(elmNam string, elmAmt int) int {

	needMap         := make(map[string]int)
	needMap[elmNam]  = elmAmt
	onlyOre         := false
	runs            := 0
	rea             := reaction{}

	for !onlyOre {
		// assume no element but Ore 
		// until one was hit by loop
		onlyOre = true

		// loop through needed elements
		for nam, amt := range needMap {

			// ignore ore
			if nam != "ORE" && amt > 0 {
				onlyOre = false

				// identify the necessary reaction
				// and determine how often it needs to be run
				rea  = reaMap[nam]
				runs = amt / rea.out.amt
				if amt % rea.out.amt > 0 { runs++ }

				// not that a negative need makes sense here
				needMap[nam] -= runs * rea.out.amt 

				// adds the new needs from that reaction
				for _, elem := range rea.in {
					needMap[elem.nam] += elem.amt * runs
				}


			}
		}
	}
	return needMap["ORE"]
}

// calculate the amount of element elmName that is 
// produced by oreAmt of ore
func calcFuel(elmNam string, oreAmt int) int {

	// assuming at least 1 ore per fuel (max guess = oreAmt)
	// using binary search
	return sort.Search(oreAmt, func(n int) bool {
		return calcNeed(elmNam, n+1) > oreAmt
	} )
}

// MAIN ----

// global reaction map parsed from input
// shows what can produced by what combination
var reaMap map[string]reaction

func main () {

	start := time.Now()

	parseReacts(readTxtFile("d14.input.txt"))

	fmt.Printf("\nOre needed: %v\n", calcNeed("FUEL", 1))
	fmt.Printf("Fuel produced: %v\n\n", calcFuel("FUEL", 1000000000000))
	fmt.Printf("Execution time: %v\n", time.Since(start))
}

