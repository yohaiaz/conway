package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type ConwayGame struct {
	world           [][]int
	generationCount int
	infectAfter     int
	max             int
}

type mutation struct {
	row int
	col int
	val int
}

func NewGame(width int, height int, infectAfter int, max int) *ConwayGame {

	matrix := make([][]int, height)

	for i := range matrix {
		matrix[i] = make([]int, width)
	}

	return &ConwayGame{
		world:       matrix,
		infectAfter: infectAfter,
		max:         max,
	}
}

func (g *ConwayGame) Generate() chan string {

	output := make(chan string, 10)

	g.init()

	output <- g.display()

	go func() {
		defer close(output)

		for i := 0; i < g.max; i++ {
			g.generationCount++
			g.genNext()
			output <- g.display()
		}
	}()

	return output
}

func (g *ConwayGame) init() {

	rand.Seed(time.Now().UnixNano())

	for i := range g.world {
		for j := range g.world[i] {
			g.world[i][j] = rand.Intn(2)
		}
	}
}

func (g *ConwayGame) genNext() {

	sm := sync.Map{}

	wg := sync.WaitGroup{}

	sem := make(chan bool, 10)

	for i := 0; i < len(g.world); i++ {
		for j := 0; j < len(g.world[0]); j++ {
			sem <- true
			wg.Add(1)

			go func(i int, j int) {

				defer func() {
					wg.Done()
					<-sem
				}()

				neighbors := g.getNeighbors(i, j)

				currentValueCell := g.world[i][j]

				infect := 0

				if g.generationCount >= g.infectAfter {
					infect = 1
				}

				nextValueCell := g.decideLiveOrDead(currentValueCell, neighbors, infect)

				if currentValueCell != nextValueCell {

					m := &mutation{
						i,
						j,
						nextValueCell,
					}

					sm.Store(fmt.Sprintf("%d-%d", i, j), m)
				}
			}(i, j)
		}
	}

	wg.Wait()

	g.dump(sm)
}

func (g *ConwayGame) getNeighbors(row, col int) []int {
	neighbors := make([]int, 0)

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			if g.isOnMap(row+i, col+j) {
				neighbors = append(neighbors, g.world[row+i][col+j])
			}
		}
	}

	return neighbors
}

func (g *ConwayGame) isOnMap(x int, y int) bool {

	return x >= 0 && y >= 0 && x < len(g.world) && y < len(g.world[0])
}

func (g *ConwayGame) decideLiveOrDead(cell int, neighbors []int, infected int) int {

	total := 0

	for _, n := range neighbors {
		total += n
	}

	if infected == 0 {
		// live cell with fewer than 2 - under population
		if cell == 1 && total < 2 {
			return 0
		}

		// live cell with 2 or 3
		if cell == 1 && (total == 2 || total == 3) {
			return 1
		}

		// live cell with more than 3 - over population
		if cell == 1 && total > 3 {
			return 0
		}

		// dead cell with 3
		if cell == 0 && total == 3 {
			return 1
		}
	} else {
		// dead cell with 3
		if cell == 0 && total == 1 {
			return 1
		}

		// live cell with 2 or 3
		if cell == 1 && total == 0 {
			return 0
		}
	}

	return cell
}

func (g *ConwayGame) dump(sm sync.Map) {
	sm.Range(func(k, v interface{}) bool {
		m := v.(*mutation)
		g.world[m.row][m.col] = m.val
		return true
	})
}

func (g *ConwayGame) display() string {

	var sb = strings.Builder{}

	for _, row := range g.world {
		for _, v := range row {
			sb.WriteString(fmt.Sprintf("%d ", v))
		}
		//sb.WriteString(fmt.Sprintf("\n"))
	}

	return sb.String()
}
