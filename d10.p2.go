package main

import (
	"fmt"
	"os"
	"bufio"
	"math"
	"sort"
	"time"
)


// no error handling ...
func readTxtFile(name string) (data []string) {
	
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {	
		data = append(data, scanner.Text())
	}

	return
}

func abs(v int) int {
	if (v < 0) {
		return -v
	}
	return v
}

// build internal representation of map
// so that first index is x, second is y
// NOTE: that means amap[0] is the first COLUMN not row	
func buildAMap(matrix []string) (amap [][]bool) {

	width  := len(matrix[0])
	height := len(matrix)
	amap = make([][]bool, width)
	for i := range amap {
		amap[i] = make([]bool, height)
	}

	for j, s := range matrix {
		for i := 0; i < width; i++ {

			// 35 is ascii for #
			if s[i] == 35 {
				amap[i][j] = true
			}
		}
	}
	return
}

type vector struct {
	x,y 	int
	mDist 	int
	angle  	float64
	rotIx 	int // which rotation is that objct visible
}

func getVector(x, y, xc, yc int) (vec vector) {

	// negative in order to have lowest values for straight up (0, -n) 
	vec.angle = - math.Atan2(float64(x-xc), float64(y-yc))

	// the distance does not need to be exact
	// it's only used to compare two with the SAME angle
	// so it only needs to be monotone -> Manhattan Distance
	vec.mDist = abs(x-xc)+abs(y-yc)
	vec.x = x
	vec.y = y
	return
}

func pushElem(vecMap []vector, ix int) {

	vec := vecMap[ix] // the element to be pushed to the end	
	vec.rotIx += 1 	  // increment rotation index
	copy(vecMap[ix:], vecMap[ix+1:]) 
	vecMap[len(vecMap)-1] = vec
}

func firingSequence(amap [][]bool, xc, yc int) (vecMap []vector) {

	// temporary removal of center asteroid for counting
	// so it does not "fire on itself"
	amap[xc][yc] = false
	var vec, pVec vector

	// create vector representations of all asteroids
	for x, col := range amap {
		for y, bol := range col {
			if bol {
				vec = getVector(x, y, xc, yc)
				vecMap = append(vecMap, vec)
			} 
		}
	}

	// sorting by angle/distance
	sort.Slice(vecMap, func(i, j int) bool {
		if vecMap[i].angle < vecMap[j].angle {
			return true
		}
		return vecMap[i].angle == vecMap[j].angle && vecMap[i].mDist < vecMap[j].mDist
	})

	// go throug and push repeats that are hit in the same rotation
	// to the end of the line incresing the rotation index
	pVec = vecMap[0]		
	for i := 1; i < len(vecMap); i++ {
		if vecMap[i].angle == pVec.angle && vecMap[i].rotIx == pVec.rotIx {
			pushElem(vecMap, i)
			i--
		} else {
			pVec = vecMap[i]
		}		
	}

	// bringing the removed asteroid back
	amap[xc][yc] = true
	return
}

func main() {

	start  := time.Now()
	matrix := readTxtFile("d10.input.txt")
	xc     := 23 // from part 1
	yc     := 20 // from part 1


	amap   := buildAMap(matrix)
	seq    := firingSequence(amap, xc, yc)
	for i, v := range seq {
		if i == 199 {
			fmt.Printf("\n(%v)[%v,%v]-[a:%v,d:%v,r:%v]\n", i+1, v.x, v.y, v.angle, v.mDist, v.rotIx)
			fmt.Println("Result:", v.x*100 + v.y, "\n") 
		}
	}
	fmt.Printf("Execution time: %v\n", time.Since(start))

}

