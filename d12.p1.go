package main

import (
	"fmt"
)

func sign(val int) (sign int) {

	if val < 0 {
		return -1
	} else if val == 0 {
		return 0
	}
	return 1
}

func abs(val int) (sign int) {

	if val < 0 {
		return -val
	}
	return val
}

type moon struct {
	x,y,z    int // position
	vx,vy,vz int // velocity
	nm       string
}

func (m *moon) dump() {
	fmt.Printf("pos=<x= %2v, y=%2v, z=%2v>, vel =<x= %2v, y=%2v, z=%2v\n", m.x, m.y, m.z, m.vx, m.vy, m.vz)
}

func (m *moon) energy() int {
	pot := abs(m.x)+abs(m.y)+abs(m.z)
	kin := abs(m.vx)+abs(m.vy)+abs(m.vz)
	return pot*kin
}

// simulates the moons
func simulate(moons []moon, it int) {
	for i := 0; i < it; i++ {
		if i == it-1 {
			fmt.Printf("\nStep %v - ", i)
		}
		te := 0
		for _, moon := range moons {
		    te += moon.energy()
		}
		if i == it-1 {
			fmt.Printf("Total Energy: %v\n\n", te)
		}
		computeVelocity(moons)
		computeLocation(moons)
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

func main() {

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

	simulate(moons, 1001)
}