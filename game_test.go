package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewGame(t *testing.T) {

	testTable := []struct {
		width       int
		height      int
		infectAfter int
		max         int
	}{
		{2, 3, 2, 5},
		{3, 2, 2, 5},
		{3, 3, 2, 5},
	}

	for _, test := range testTable {
		g := NewGame(test.width, test.height, test.infectAfter, test.max)

		if g.infectAfter != test.infectAfter {
			t.Errorf("param 'infectAfter' not initialized correctly")
		}

		if g.max != test.max {
			t.Errorf("param 'max' not initialized correctly")
		}

		if len(g.world[0]) != test.width {
			t.Errorf("param 'width' not initialized correctly")
		}

		if len(g.world) != test.height {
			t.Errorf("param 'height' not initialized correctly")
		}

		if g == nil {
			t.Errorf("failed to instantiate new game")
		}
	}
}

func TestIsOnMap(t *testing.T) {
	g := NewGame(3, 3, 99, 2)

	g.world = [][]int{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}

	testTable := []struct {
		x       int
		y       int
		isOnMap bool
	}{
		{0, 0, true},
		{2, 2, true},
		{3, 3, false},
		{1, 5, false},
		{5, 1, false},
	}

	for _, test := range testTable {
		if isOnMap := g.isOnMap(test.x, test.y); isOnMap != test.isOnMap {
			if test.isOnMap {
				t.Errorf("failed: (%d, %d) is on map", test.x, test.y)
			} else {
				t.Errorf("failed: (%d, %d) is not on map", test.x, test.y)
			}
		}
	}
}

func TestDump(t *testing.T) {

	testTable := []struct {
		x   int
		y   int
		val int
	}{
		{0, 0, 0},
		{2, 2, 1},
		{1, 2, 0},
		{2, 1, 1},
	}

	sm := sync.Map{}

	for _, test := range testTable {
		m := &mutation{
			test.x,
			test.y,
			test.val,
		}
		sm.Store(fmt.Sprintf("%d-%d", test.x, test.y), m)
	}

	g := NewGame(3, 3, 99, 2)

	g.world = [][]int{{1, 0, 0}, {0, 0, 1}, {0, 0, 0}}

	g.dump(sm)

	for _, test := range testTable {
		if g.world[test.x][test.y] != test.val {
			t.Errorf("failed to dump: (%d, %d) = %d", test.x, test.y, test.val)
		}
	}
}

func TestDecideLiveOrDead(t *testing.T) {

	testTable := []struct {
		cell      int
		neighbors []int
		infected  int
		state     int
		err       string
	}{
		{1, []int{0, 1, 0}, 0, 0, "expected value 0"},
		{1, []int{0, 1, 1}, 0, 1, "expected value 1"},
		{1, []int{1, 1, 1}, 0, 1, "expected value 1"},

		{1, []int{1, 1, 0, 0, 1, 0, 1}, 0, 0, "expected value 0"},
		{0, []int{1, 1, 1}, 0, 1, "expected value 1"},
		{0, []int{0, 1, 0}, 1, 1, "expected value 1"},
		{1, []int{0, 0, 0}, 1, 0, "expected value 0"},
	}

	g := NewGame(3, 3, 99, 2)

	for _, test := range testTable {
		if g.decideLiveOrDead(test.cell, test.neighbors, test.infected) != test.state {
			t.Error(test.err)
		}
	}
}

// go test
// go test -v
// go test -cover
// go test -cover -coverprofile=c.out
// go tool cover -html=c.out -o coverage.html
