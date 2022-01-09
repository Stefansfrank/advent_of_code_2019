// Note: this solution is completely driven by part 2 and computes part 1 in less than microsecond on the side
// my original part 1 solution was a simulation like (most likely) anybody else ...

package main

import (
	_ "embed"
	"fmt"
	"time"
	"regexp"
	"strconv"
	"strings"
	"math/big"
)

//go:embed d22.input.txt
var inp string

var size *big.Int

const (
	rev = 0
	inc = 1
	cut = 2
)

const (
	bold = "\033[1;31m"
	norm = "\033[0m"
)

type step struct {
	tp int
	pm *big.Int
}

// parsing of the input 
func parse(inp []string) (stps []step) {

	reDl := regexp.MustCompile(`deal into new stack`)
	reDi := regexp.MustCompile(`deal with increment (\d+)`)
	reCt := regexp.MustCompile(`cut (-?\d+)`)
	stps  = []step{}

	for _, line := range inp {
		match := reDl.FindStringSubmatch(line)
		if match != nil {
			stps = append(stps, step{rev, big.NewInt(0)})
			continue
		}
		match  = reDi.FindStringSubmatch(line)
		if match != nil {
			stps = append(stps, step{inc, big.NewInt(int64(atoi(match[1])))})
			continue
		}
		match  = reCt.FindStringSubmatch(line)
		if match != nil {
			stps = append(stps, step{cut, big.NewInt(int64(atoi(match[1])))})
		}
	}

	return
}

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

// creates the initial deck
func initial() (deck []int) {
	sz := int(size.Int64())
	deck = make([]int, sz)

	for i := 0; i < sz; i++ {
	 	deck[i] = i
	}
	return
}

// representation of one coefficient pair a,b 
type coeff struct {
	a,b *big.Int
}

// calculates the coefficients a, b necessary in order to express
// a shuffle move as linear equations on the position of a card: 
// 		npos = a * pos + b
// this is possible if we assume the deck to be cyclic and we can fold
// back values into the actual deck using (% size) operations
// ------------------------------------------------------------------------
// we also use the fact that the parameters a, b themselves can be reduced
// to a size lower than the deck size by subtracting / adding the deck size
// from either parameter since that changes the total (a * pos + b) only in
// multiples of the deck size which are NOPs in a cyclic deck
// ------------------------------------------------------------------------
// the cyclic deck also has the advantage that we do not have to treat 
// negative and positive cuts any different
func getCoeffs(stps []step) (cofs []coeff) {

	var a,b *big.Int
	t := big.NewInt(0)
	cofs = make([]coeff, 0, len(stps))

	for _, stp := range stps {
		switch stp.tp {
		case rev:
			a = big.NewInt(-1)
			b = t.Sub(size, big.NewInt(1)) 
		case inc:
			a = stp.pm 
			b = big.NewInt(0)
		case cut:
			a = big.NewInt(1)
			b = t.Sub(size, stp.pm)
		}
		cofs = append(cofs, coeff{ nmod(a), nmod(b) })
	}
	return
}

// this function uses the fact that two linear equations y = a * x + b 
// executed one after each other can be represented as one linear equation:
// y = (a1 * a2) * x + (a2 * b1 + b2) -> new a = a1 * a2 / new b = a2 * b1 + b2
// thus an arbitrary amount of shuffle moves can be shortened to one coefficient pair a,b
// that allows a super simple linear calculation for an individual card
func simplify(cofs []coeff) (cof coeff) {
	t := big.NewInt(0)
	cof.a = big.NewInt(1)
	cof.b = big.NewInt(0)
	for _, c := range cofs {
		cof.a.Set(nmod(t.Mul(cof.a, c.a)))
		cof.b.Set(nmod(t.Add(t.Mul(cof.b, c.a), c.b)))
	}
	return
}

// basically a modulo function but one that properly projects negative parameters
// into the positive space using the deck size
func nmod(v *big.Int) *big.Int {
	t := big.NewInt(0)
	if v.Sign() < 0 {
		return t.Add(size, v.Mod(v, size)) 
	} else {
		return t.Mod(v, size)
	}
}

// shuffling now becomes somewhat trivial as we have a simple linear index conversion 
func shuffle(deck []int, cof coeff) (ndeck []int) {

	ndeck = make([]int, len(deck))
	t    := big.NewInt(0) 
	for ix, c := range deck {
	  	ndeck[ int(nmod(t.Add(t.Mul(cof.a, big.NewInt(int64(ix))), cof.b)).Int64()) ] = c
	}
	return
}

// this is a trick to execute many iterations of the same linear coefficients. The idea is
// that simplify above allows to express multiple nested linear transformations as one linear 
// transformation. Thus we first compute the coefficients that express the application of the
// original coefficient pair 1x, 2x, 4x, 8x ... 
// We then express the total number of iterations as a binary and add the according c1x, c2x, c4x
// coefficient pairs if that bit is set. 
// For example: if n would be 11 (0b1011) we would apply c1x, c2x and c4x instead of 11 times c1x
// ... of course, we simplify these down to one coefficient pair again ...
// I like the idea of caluclating these 101741582076661 shuffles by just two linear coefficients :)
func simplMany(cof coeff, n *big.Int) coeff {

	num   := n.BitLen()
	t     := big.NewInt(0) 

	// calculate the coefficients representing 1x, 2x, 4x ... application of 'cof'
	bcofs := make([]coeff, num) 
	bcofs[0] = cof
	for i := 1; i < num; i++ {
		bcofs[i] = simplify([]coeff{ bcofs[i-1], bcofs[i-1] })
	}

	// apply to binary representation of n
	mcofs := make([]coeff, 0, num)
	for i := 0; i < num; i++ {
		if t.Mod(n, big.NewInt(2)).Sign() == 1 {
			mcofs = append(mcofs, bcofs[i])
		}
		n.Rsh(n, uint(1))
	}
	return simplify(mcofs)
}

// MAIN --------------------------------------------

func main () {
	start := time.Now()
	t := big.NewInt(0)

	// Part 1 (step by step) =============================
	size = big.NewInt(10007)

	// parses the input into the steps format 
	stps := parse(strings.Split(strings.TrimSpace(inp), "\n"))

	// converts the steps into linear coefficient pairs
	cofs := getCoeffs(stps)

	// converts list of coefficients for all steps into one coefficient pair
	cof  := simplify(cofs)

	// shuffles the deck using that one coefficient pair
	deck := shuffle(initial(), cof)

	// searches for card 2019
	for ix, crd := range deck {
		if crd == 2019 {
			fmt.Println("\n2019 is at position" + bold, ix, norm)
			break
		}
	}

	// Part 2 (step by step) ============================
	size   = big.NewInt(119315717514047)
	count := big.NewInt(101741582076661)

	// recomputing the coefficients in order to ensure that 
	// none of the above changed the starting situation 
	// (since I use *big.Int, there is always the danger I accidentally change values)
	cofs = getCoeffs(stps)
	cof  = simplify(cofs)

	// this is the last trick we use: Since the question is which card is in position 2020 after 
	// shuffling, the problem here is to invert all computations. This is possible by solving npos = pos * a + b
	// for pos -> pos = npos * (b/a) - (1/a) but I could not find a performant way to find n, m so that 
	// (1 + n*size)/a and (b + m*size)/a are integer. 
	// I tried to just brut force adding size and check but that is of the order of O(n^2) and 
	// with the ~50 bit numbers given I ended up with ~10^100 operations :( (well, it never 'ended' ;)
	//
	// Instead I am using Fermat's little theorem that says that if x is a prime number,
	// applying the linear coefficents a,b (x - 1) times gives an outcome of a = (1 % x), b = (0 % x) 
	// In other words, since 'size' is a prime number, the deck is in it's original state after (size - 1)
	// identical shuffles no matter what the shuffle is (as the deck is cyclic and (% size) is a NOP)
	//
	// This means instead of inverting the shuffles, we just 'continue' the shuffles until (size - 1) and see
	// where the card now in pos 2020 will end up (i.e. we apply the shuffle (size - count -1) times)
	//
	// Note: simplMany simplifies many identical shuffles into one coeff pair
	ncof  := simplMany(cof, t.Sub(t.Sub(size, count), big.NewInt(1)))
	fmt.Println("Card ending up in position 2020:" + bold, nmod(t.Add(t.Mul(ncof.a, big.NewInt(2020)) ,ncof.b)), norm + "\n")

	fmt.Println("Execution time: ", time.Since(start))
}
