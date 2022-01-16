package main

import (
	"fmt"
	"os"
	"time"
)

// Utility Functions ----------------------------------------------------------------------

type layer struct {
	data   []byte
	digNum []int
}

// no error handling ...
func readDigitFile (name string, width, height int) (layers []layer) {
	
	file, _ := os.Open(name)
	defer file.Close()

	bLength := width * height
	var bTemp []byte
	var lyr layer

	for {
		bTemp   = make([]byte, bLength)
	    n, err := file.Read(bTemp)
	    if err != nil {
	        break
	    } else if n < bLength {
	    	fmt.Printf("Input length %v is not a multiple of %v(%v x %v)\n", n, bLength, width, height)
	    }
	    lyr = layer{data:bTemp, digNum:make([]int,10)}

	    // processing
	    for i,_ := range lyr.data {
	    	lyr.data[i] -= 48
	    	lyr.digNum[lyr.data[i]]++
	    }
	    layers = append(layers, lyr)
	  
	}

	return
}

// detect 1s * 2s for min 0s (part 1)
func detectMinPrd(layers []layer) (result int) {

	min    := 200
	for _, lyr := range layers {
		if lyr.digNum[0] < min {
			result = lyr.digNum[1] * lyr.digNum[2]
			min = lyr.digNum[0]
		}
	}
	return 

}

// 
func main() {

	start  := time.Now()

	layers := readDigitFile("d8.input.txt", 25, 6)
	fmt.Printf("\nResult: %v\n\n", detectMinPrd(layers))
	fmt.Printf("Execution time: %v\n", time.Since(start))
}

