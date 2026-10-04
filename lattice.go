package koleya

import "math"

// Lattice is the grid vehicles come to rest on: nodes lie at
// Origin + k·Step on both axes. A Step of zero or less means no grid:
// vehicles stop wherever their braking ends.
type Lattice struct {
	// Step is the distance between nodes, in pixels.
	Step float64
	// Origin shifts the nodes on both axes, in pixels.
	Origin float64
}

// Nearest returns the node closest to the coordinate.
func (l Lattice) Nearest(coord float64) float64 {
	if l.Step <= 0 {
		return coord
	}
	return l.Origin + math.Round((coord-l.Origin)/l.Step)*l.Step
}

// NodeAhead returns the first node at or beyond the coordinate in the
// direction of sign; for a zero sign, the nearest node.
func (l Lattice) NodeAhead(coord float64, sign int) float64 {
	if l.Step <= 0 {
		return coord
	}
	k := (coord - l.Origin) / l.Step
	switch {
	case sign > 0:
		k = math.Ceil(k - eps)
	case sign < 0:
		k = math.Floor(k + eps)
	default:
		k = math.Round(k)
	}
	return l.Origin + k*l.Step
}

// Near reports whether the coordinate is within tolerance of a node.
func (l Lattice) Near(coord, tolerance float64) bool {
	return math.Abs(coord-l.Nearest(coord)) <= tolerance
}

// On reports whether the coordinate lies on a node.
func (l Lattice) On(coord float64) bool {
	return l.Near(coord, eps*math.Max(l.Step, 1))
}
