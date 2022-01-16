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

// decoding transparency and such (part 2)
func decodeLayers(layers []layer, width int) string {

	img  := layers[0].data
	imgS := "\n"
	for i, lyr := range layers {
		if i == 0 { continue }

		for j := 0; j < len(lyr.data); j++ {
			if img[j] == 2 {
				img[j] = lyr.data[j]
			}
		}
	}

	for i := 0; i < len(img); i += width {
		for _, j := range img[i:i+width] {
			if j == 0 { imgS += " " } else { imgS += "*"}
		}
		imgS +="\n"
	}

	return imgS

}

// 
func main() {

	start  := time.Now()
	width  := 25
	height := 6

	layers := readDigitFile("d8.input.txt", width, height)

	fmt.Println(decodeLayers(layers, width))
	fmt.Printf("Execution time: %v\n", time.Since(start))
}

