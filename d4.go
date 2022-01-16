package main

import (
	"fmt"
	"time"
)

// Goes through individual digits from lowest to highest 
// and checks contitions
func checkCode1 (code int) bool {
	nxtDig  := 0
	pairDet := false
	curDig  := code % 10
	code     = code / 10

	for i :=0; i < 5; i++ {
		nxtDig = code % 10
		code   = code / 10

		// next highest digit is higher than current digit
		if nxtDig > curDig { return false }

		// next highest digit is equal
		pairDet = pairDet || (nxtDig == curDig) 

		curDig = nxtDig
	}

	return pairDet 
}

// Goes through individual digits from lowest to highest 
// and checks contitions
func checkCode2 (code int) bool {
	nxtDig  := 0
	pairDet := false
	rptCnt  := 0
	curDig  := code % 10
	code     = code / 10

	for i :=0; i < 5; i++ {
		nxtDig = code % 10
		code   = code / 10

		// next highest digit is higher than current digit
		if nxtDig > curDig { return false }

		// next highest digit is equal
		if nxtDig == curDig { 
			rptCnt++ 
		
		// next highest digit is lower
		} else {

			// there was exactly one repetition before
			if rptCnt == 1 {
				pairDet = true
			}
			rptCnt = 0
		}

		curDig = nxtDig
	}

	// either a pair was detected before or 
	// the highes digit had exactly one repition 
	return pairDet || rptCnt == 1
}

func main() {

	start := time.Now()

	// My input: 183564-657474
	from  := 183564
	to    := 657474 

	count := 0
	for i := from; i <= to; i++ {
		if checkCode1(i) { 
			count ++ 
		}
	}
	fmt.Printf("\nTotal for part 1: %v\n", count) 

	count = 0
	for i := from; i <= to; i++ {
		if checkCode2(i) { 
			count ++ 
		}
	}
	fmt.Printf("Total for part 2: %v\n\n", count) 
	fmt.Println(time.Since(start))

}