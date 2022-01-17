package main

import (
	"fmt"
	"time"
)

// returns -1, 0 , 1 depending on sign of input
func sign(val int) (sign int) {

	if val < 0 {
		return -1
	} else if val == 0 {
		return 0
	}
	return 1
}

// simple absolute function
func abs(val int) (sign int) {

	if val < 0 {
		return -val
	}
	return val
}

// greatest common divisor (GCD)
// Euclidean algorithm
func gcd(a, b int) int {
	
	tmp := 0
	for b != 0 {
		tmp = b
		b = a % b
		a = tmp
	}
	return a
}

// least common multiple (LCM) using GCD
func lcm(a, b int, ints ...int) (result int) {

	result = a * b / gcd(a, b)
	for i := 0; i < len(ints); i++ {
		result = lcm(result, ints[i])
	}
	return 
}

// simple []int compare
func equals(a []int, b []int) bool {
	for i := range a {
		if a[i] != b[i] { return false}
	}
	return true
}

// --------------------------------------------------------------------------------------------------------------

type moon struct {
	x,y,z    int // position
	vx,vy,vz int // velocity
	nm       string
	px,py,pz int // period of the system (only in first moon)
}

func (m *moon) dump() {
	fmt.Printf("pos=<x= %2v, y=%2v, z=%2v>, vel =<x= %2v, y=%2v, z=%2v\n", m.x, m.y, m.z, m.vx, m.vy, m.vz)
}

func (m *moon) energy() int {
	pot := abs(m.x)+abs(m.y)+abs(m.z)
	kin := abs(m.vx)+abs(m.vy)+abs(m.vz)
	return pot*kin
}

// simulates moons until repetition is detected in all three axis
func simulateRpt(moons []moon) {
	snap := snapshot(moons)
	for i := 0; true; i++ {
		computeVelocity(moons)
		computeLocation(moons)
		if detectPeriods(moons, snap, i) { return }
	}
}

func computeVelocity(moons []moon) {
	sgn := 0

	// go through all permutations
	for i := 0; i < (len(moons)-1); i++ {
		for j := i+1; j < len(moons); j++ {
			sgn = sign(moons[j].x - moons[i].x) 
			moons[j].vx += -sgn
			moons[i].vx += sgn
			sgn = sign(moons[j].y - moons[i].y) 
			moons[j].vy += -sgn
			moons[i].vy += sgn
			sgn = sign(moons[j].z - moons[i].z) 
			moons[j].vz += -sgn
			moons[i].vz += sgn
		}
	}
}

func computeLocation(moons []moon) {
	for i := range moons {
		moons[i].x += moons[i].vx
		moons[i].y += moons[i].vy
		moons[i].z += moons[i].vz
	}
}

// takes a snapshot for reference
// note that the order is with all x coordinates of all moons first then all y then all z
func snapshot(moons []moon) (snap []int) {

	snap = make([]int, 6*len(moons))
	for i, moon := range moons {
		snap[i]                = moon.x
		snap[i + 2*len(moons)] = moon.y
		snap[i + 4*len(moons)] = moon.z
		snap[i +   len(moons)] = moon.vx
		snap[i + 3*len(moons)] = moon.vy
		snap[i + 5*len(moons)] = moon.vz
	}
	return
}

// this is called every step and uses snap to identify
// periods for each axis. Once all periods are identified
// it returns true
func detectPeriods(moons []moon, start []int, ix int) bool {

	lm   := len(moons)
	snap := snapshot(moons)

	// compares the x coordinates of all moons	
	if equals(snap[0:2*lm], start[0:2*lm]) {
		if moons[0].px == 0 {
			moons[0].px = ix + 1
		}
		// checks whether that was the last missing period
		if moons[0].px > 0 && moons[0].py > 0 && moons[0].pz > 0 {
			return true
		}
	}

	// compares the y coordinates of all moons	
	if equals(snap[2*lm:4*lm], start[2*lm:4*lm]) {
		if moons[0].py == 0 {
			moons[0].py = ix + 1
		}
		// checks whether that was the last missing period
		if moons[0].px > 0 && moons[0].py > 0 && moons[0].pz > 0 {
			return true
		}
	}

	// compares the z coordinates of all moons	
	if equals(snap[4*lm:6*lm], start[4*lm:6*lm]) {
		if moons[0].pz == 0 {
			moons[0].pz = ix + 1
		}
		// checks whether that was the last missing period
		if moons[0].px > 0 && moons[0].py > 0 && moons[0].pz > 0 {
			return true
		}
	}

	return false
}

// ------------------------------------------------------------------------------------------------------

func main() {

	start := time.Now()
	moons := make([]moon, 4)
	moons[0] = moon{nm:"Io"}
	moons[1] = moon{nm:"Europa"}
	moons[2] = moon{nm:"Ganymede"}
	moons[3] = moon{nm:"Callisto"}

	/* Main Input
	  <x=19, y=-10, z=7>
	  <x=1, y=2, z=-3>
	  <x=14, y=-4, z=1>
	  <x=8, y=7, z=-6> */

    moons[0].x = 19
    moons[0].y = -10
    moons[0].z = 7
	moons[1].x = 1
	moons[1].y = 2
	moons[1].z = -3
	moons[2].x = 14
	moons[2].y = -4
	moons[2].z = 1
	moons[3].x = 8
	moons[3].y = 7
	moons[3].z = -6 

	simulateRpt(moons)
	fmt.Printf("\nPer Coordinate Periods: %v %v %v\n", moons[0].px, moons[0].py, moons[0].pz)
	fmt.Printf("Total Period: %v\n\n", lcm(moons[0].px, moons[0].py, moons[0].pz))
	fmt.Printf("Execution time: %v\n", time.Since(start))
}