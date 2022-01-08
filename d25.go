package main

import (
	"fmt"
	"time"
	"ac19/ac19cpu"
	"strings"
	//"bufio" uncomment if uncommenting the interactive game below
	//"os"    uncomment if uncommenting the interactive game below
)


// Converts input string into []int
func cmdInp(s string) []int {
	inp := []int{}
	for _, c := range s {
		inp = append(inp, int(c))
	}
	return inp
}

// Get and removes the Output
func getOut(cpu *ac19cpu.Machine) (s string) {
	s   = ""
	ln := len(cpu.Output)
	for i := 0; i < ln; i++ {
		c := cpu.ConsumeOutput()
		s = s + fmt.Sprintf("%c",c)
	}
	return 
}

// -------------------------------------- Main ------------------------------------------------------------
func main() {

	start := time.Now()

	// After playing for a while, I determined the walk through:
	// this navigates through the ship and collects all
	// inventory that might be useful (and not the ones loosing the game)
	auto := []string{   "west\n",
						"take ornament\n",
						"west\n",
						"take astrolabe\n",
						"north\n",
						"take fuel cell\n",
						"south\n",
						"south\n",
						"take hologram\n",
						"north\n",
						"east\n",
						"south\n",
						"east\n",
						"take weather machine\n",
						"west\n",
						"north\n",
						"east\n",
						"east\n",
						"take mug\n",
						"north\n",
						"take monolith\n",
						"south\n",
						"south\n",
						"west\n",
						"north\n",
						"west\n",
						"take bowl of rice\n",
						"north\n",
						"west\n",
						"north\n"}

    // these indexes are used to vary the compbination of inventory pieces
    // with the goal of hitting the correct weight for the pressure sensor
	inv := []string{"ornament", "astrolabe", "fuel cell", "hologram", "weather machine", "mug", "monolith", "bowl of rice"}
	hlp := []int{1,2,4,8,16,32,64,128}

	cpu := ac19cpu.Machine{}
	cpu.LoadProgramFromCsv("d25.input.txt")

	tmp := ""
	end := false
	for !end {

		fmt.Println(" Day 25 (2019) Game")
		fmt.Println("---------------------")

		// start the game
		cpu.Execute(false)
		tmp = getOut(&cpu)

		// execute the automatic commands 
		for _, cmd := range auto {
			cpu.AddInput(cmdInp(cmd)...)
			cpu.Execute(false)
			tmp = getOut(&cpu)
		}

		// trying all combinations of inventory and step onto
		// the pressure sensor until I get a different answer than
		// "heavier" and "lighter"
		for by := 0; by < 256; by++ {
			for i := 0; i < 8; i ++ {
				if hlp[i] & by > 0 {
					cpu.AddInput(cmdInp("take " + inv[i] + "\n")...)
				} else {
					cpu.AddInput(cmdInp("drop " + inv[i] + "\n")...)
				}
				cpu.Execute(false)
				tmp = getOut(&cpu)
			}

			// go north and see what the analysis says
			cpu.AddInput(cmdInp("north\n")...)
			cpu.Execute(false)
			tmp = getOut(&cpu)

			// detect whether a different output is found
			if strings.Index(tmp, "Droids on this ship are lighter") > -1 {
				fmt.Printf("(%03v/256) - Too Heavy\n", by+1)
			} else if strings.Index(tmp, "Droids on this ship are heavier") > -1 {
				fmt.Printf("(%03v/256) - Too Light\n", by+1)
			} else {
				// yay!
				fmt.Println(tmp)
				end = true
				break
			}
		}

		// UNCOMMENT IF YOU WANT TO PLAY INTERACTIVELY
		// (if uncommented, don't forget to adapt the os and bufio imports)
		// reader := bufio.NewReader(os.Stdin)
		// inner: for !end {
		// 	text, _ := reader.ReadString('\n')
		// 	if strings.Compare("exit\n", text) == 0 {
		// 		end = true
		// 	} else if strings.Compare("restart\n", text) == 0 {
		// 		cpu.ResetProgram()
		// 		break inner
		// 	} else {
		// 		cpu.AddInput(cmdInp(text)...)
		// 		cpu.Execute(false)
		// 		fmt.Println(getOut(&cpu))
		// 	}
		// }

	}

	fmt.Printf("Execution time: %v\n", time.Since(start))
}