package main

import (
	"flag"
	"fmt"
)

func main() {

	width := flag.Int("w", 3, "matrix width")
	height := flag.Int("h", 3, "matrix height")
	infAfter := flag.Int("i", 2, "infected after number of generation")
	max := flag.Int("m", 5, "max generation")
	needHelp := flag.Bool("help", false, "print help")

	flag.Parse()

	if *needHelp && (flag.NFlag()+len(flag.Args()) == 0) {
		flag.Usage()
		return
	}

	fmt.Println("conway next generation game of life - Infection")
	fmt.Println(fmt.Sprintf("running with params %d, %d, %d, %d ", *width, *height, *infAfter, *max))

	c := NewGame(*width, *height, *infAfter, *max)

	genChan := c.Generate()

	for g := range genChan {
		fmt.Println(g)
	}
}
