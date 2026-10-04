package koleya

// Axis is one of the two axes of the plane.
type Axis uint8

const (
	// X is the horizontal axis, growing to the right.
	X Axis = iota
	// Y is the vertical axis, growing downwards, as on the screen.
	Y
)

// Other returns the perpendicular axis.
func (a Axis) Other() Axis {
	return 1 - a
}

// String returns "X" or "Y".
func (a Axis) String() string {
	if a == X {
		return "X"
	}
	return "Y"
}

// Vec2 is a point or a vector in pixels.
type Vec2 struct {
	X, Y float64
}

// On returns the component along the axis.
func (v Vec2) On(a Axis) float64 {
	if a == X {
		return v.X
	}
	return v.Y
}

// at returns a pointer to the component along the axis.
func (v *Vec2) at(a Axis) *float64 {
	if a == X {
		return &v.X
	}
	return &v.Y
}
