package main

import (
	"fmt"
	"os"
	"bufio"
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

// detects not only the absolute value of an int
// but also returns -1 or 1 for polarity
func absPol(v int) (abs int, pol int) {
	if (v < 0) {
		return -v, -1
	} else {
		return v, 1
	}
	return
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

// location for the station
type location struct {
	x,y int
	vis int // amount of asteroids visible
}

// detects the best location
func findLocation(amap [][]bool) (loc location) {
	
	tmp := 0
	for x, col := range amap {
		for y, bol := range col {
			if bol {
				tmp = countVisible(amap, x, y)
				if tmp > loc.vis {
					loc.x = x
					loc.y = y
					loc.vis = tmp
				}
			}
		}
	}
	return
}

func countVisible(amap [][]bool, xc, yc int) (vis int) {

	// temporary removal of center asteroid for counting
	// so it does not "see itself"

	// TODO unused optimization potential:
	// ever pair of asteroids is detected twice
	// I could keep a map of previously detected visibilities around
	amap[xc][yc] = false

	for x, col := range amap {
		for y, bol := range col {
			if bol {
				if isVisible(amap, xc, yc, x, y) {
					vis++
				}
			} 
		}
	}

	// bringing the removed asteroid back
	amap[xc][yc] = true
	return
}

func isVisible(amap [][]bool, xc, yc, x, y int) bool {

	dx, px := absPol(x-xc)
	dy, py := absPol(y-yc)
	var xi, yi int

	// straight line cases are simpler
	if dx == 0 {
		for yi = 1; yi < dy; yi++ {
			if amap[xc][yc + py * yi] { return false }
		}
		return true
	}
	if dy == 0 {
		for xi = 1; xi < dx; xi++ {
			if amap[xc + px * xi][yc] { return false }
		}
		return true
	}

	// diagonal cases separated by wich coordinate is smaller
	// iterating through the smaller one
	// trying to stay with integer algorithm and modulo checks
	if (dx <= dy) {
		for xi = 1; xi < dx; xi++ {
			if dy * xi % dx == 0 {
				if amap[xc + px * xi][yc + py * dy * xi / dx] {return false}
			}
		}
		return true
	}
	if (dx > dy) {
		for yi = 1; yi < dy; yi++ {
			if dx * yi % dy == 0 {
				if amap[xc + px * dx * yi / dy][yc + py * yi] {return false}
			}
		}
		return true
	}

	// never hit
	return true
}

func main() {

	start := time.Now()
	matrix := readTxtFile("d10.input.txt")

	amap   := buildAMap(matrix)
	loc    := findLocation(amap)

	fmt.Printf("Best location: [%v,%v] with line of view to %v asteroids\n", loc.x, loc.y, loc.vis)
	fmt.Printf("Execution time: %v\n", time.Since(start))

}

