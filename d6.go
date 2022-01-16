package main

import (
	"fmt"
	"time"
	"os"
	"bufio"
)

type obj struct{
	name      string
	orbiting  string
	orbitedBy []string
	indOrb    int
}

// no error handling ...
func readTxtFile (name string) (lines []string) {
	
	file, _ := os.Open(name)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {		
		lines = append(lines, scanner.Text())
	}

	return

}

// parses the Map into go map
// setting orbiting relationships on the way
func parseMap (starMap []string) {

	var cen, obj string
	for _, line := range starMap {
		cen = line[0:3]
		obj = line[4:7]
		createObj(cen)
		createObj(obj)
		objMap[obj].addCenter(cen)
	}

	fmt.Printf("Map Parsed: %v objects created, %v center(s) detected\n", len(objMap), len(cenMap))
	if len(cenMap) == 1 {
		for name, _ := range cenMap {
			center = name
		}
	}
}

// adds a center to an object
func (o *obj) addCenter(cen string) {

	// can't resist a little error management ...
	if len(o.orbiting) > 0 {
		fmt.Printf("Object %v already orbits %v, can't switch to %v\n", o.name, o.orbiting, cen)
		return
	}

	// orbiting relationships in both directions
	o.orbiting = cen

	if len(objMap[cen].orbitedBy) > 0 {
		objMap[cen].orbitedBy = append(objMap[cen].orbitedBy, o.name)
	} else {
		objMap[cen].orbitedBy = []string{o.name}
	}

    // delete as possible center
	delete(cenMap, o.name)

}

// creates new object if it does not exist
func createObj(name string) {

	_, exists := objMap[name]
	if !exists {
		objMap[name] = &obj{name: name}
		cenMap[name] = true
	}
}

// detect number of orbiting relationships
// recursively going up the tree
func detectOrbNum() (orbNum int) {

	for name, _ := range cenMap {

		// indirect orbits for centers are -1 since
		// - cancels out the lack of direct orbit in the sum of all orbits
		// - makes sure the first orbiters have an indirect count of 0
		orbNum += recObjCnt(name, -1)
	}

	return
}

// recursion function determines indirect orbits
// and adds 1 to total count for the direct orbit
func recObjCnt(nm string, ind int) (cnt int) {

	objMap[nm].indOrb = ind
	cnt += 1 + ind

	for _, nm := range objMap[nm].orbitedBy {
		cnt += recObjCnt(nm, ind+1)
	}

	return
}

// detect shortest link between two objects
func shortestLink(nmFrom, nmTo string) (hops int) {
	
	fromMap  := buildOrbMap(nmFrom, center)
	curObj   := objMap[nmTo]
	distance := 0

	for curObj.name != center {
		for fromObj, fromDist := range fromMap {
			if fromObj == curObj.name {
				return fromDist + distance - 2 // reducing since I count the from/to objects themselves
			}
		}
		curObj = objMap[curObj.orbiting]
		distance++
	}

	return -1 // no link found (should not happen if there is a common center)

}

// build map for objects
func buildOrbMap(name string, cnt string) (orbMap map[string]int) {

	curObj   := objMap[name]
	distance := 0
	orbMap   = make(map[string]int)

	for curObj.name != cnt {
		orbMap[curObj.name] = distance
		curObj = objMap[curObj.orbiting]
		distance ++
	}
	orbMap[center] = distance
	return
}

// Being lazy, I declare the map globally
var objMap map[string]*obj

// This is a map of all centers - should be one at the end
var cenMap map[string]bool

// absolute center if there is one
var center string

func main() {
	start := time.Now()

	objMap = make(map[string]*obj)
	cenMap = make(map[string]bool)

	starMap := readTxtFile("d6.input.txt")
	parseMap(starMap)

	fmt.Printf("\nNumber of orbits: %v\n", detectOrbNum())
	fmt.Printf("Number of hops to Santa: %v\n", shortestLink("YOU", "SAN"))
	fmt.Printf("\nSpecial Objects:\n%v\n%v\n%v\n", *objMap["YOU"], *objMap["SAN"], *objMap[center])

	fmt.Printf("Execution time: %v\n", time.Since(start))
}
