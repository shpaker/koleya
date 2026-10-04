package koleya

import "testing"

func TestNodeAhead(t *testing.T) {
	grid := Lattice{Step: 4}
	cases := []struct {
		coord float64
		sign  int
		want  float64
	}{
		{5.5, 1, 8},
		{5.5, -1, 4},
		{8, 1, 8},
		{8, -1, 8},
		{-1, 1, 0},
		{-1, -1, -4},
		{5.5, 0, 4},
	}
	for _, c := range cases {
		if got := grid.NodeAhead(c.coord, c.sign); got != c.want {
			t.Errorf("NodeAhead(%v, %d) = %v, want %v", c.coord, c.sign, got, c.want)
		}
	}
}

func TestNodeAheadWithOrigin(t *testing.T) {
	grid := Lattice{Step: 8, Origin: 2}
	if got := grid.NodeAhead(3, 1); got != 10 {
		t.Errorf("NodeAhead = %v, want 10", got)
	}
}

func TestNear(t *testing.T) {
	grid := Lattice{Step: 8}
	if !grid.Near(14, 2) || !grid.Near(17.5, 2) {
		t.Error("14 and 17.5 are within 2 of 16")
	}
	if grid.Near(13, 2) {
		t.Error("13 is 3 away from 16")
	}
}

func TestNoGrid(t *testing.T) {
	grid := Lattice{}
	if got := grid.NodeAhead(5.5, 1); got != 5.5 {
		t.Errorf("NodeAhead without a grid = %v, want 5.5", got)
	}
}
