package main

import (
	"flag"
)

func main() {
	argCost := flag.Bool("c", false, "set to calculate the cost of the program instead of running it")
	argLimit := flag.Int("l", 5000, "set the upper limit of simulation steps. Default: 5000")

	flag.Parse()

	if *argCost {
		calculateCost()
	} else {
		simulateProgram(*argLimit)
	}
}
