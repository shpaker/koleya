package koleya

import "math"

// Dir is one of the eight directions of the screen, or None.
type Dir uint8

const (
	// None is no direction: the driver lets go.
	None Dir = iota
	Up
	Down
	Left
	Right
	UpLeft
	UpRight
	DownLeft
	DownRight
)

// diagonal is the length of each part of a unit diagonal, √½.
var diagonal = math.Sqrt2 / 2

// DirOf returns the direction with the given signs along X and Y
// (each -1, 0 or 1; anything else counts by its sign).
func DirOf(sx, sy int) Dir {
	switch {
	case sx == 0 && sy < 0:
		return Up
	case sx == 0 && sy > 0:
		return Down
	case sx < 0 && sy == 0:
		return Left
	case sx > 0 && sy == 0:
		return Right
	case sx < 0 && sy < 0:
		return UpLeft
	case sx > 0 && sy < 0:
		return UpRight
	case sx < 0 && sy > 0:
		return DownLeft
	case sx > 0 && sy > 0:
		return DownRight
	}
	return None
}

// Signs returns the signs of the direction along X and Y.
func (d Dir) Signs() (sx, sy int) {
	switch d {
	case Up:
		return 0, -1
	case Down:
		return 0, 1
	case Left:
		return -1, 0
	case Right:
		return 1, 0
	case UpLeft:
		return -1, -1
	case UpRight:
		return 1, -1
	case DownLeft:
		return -1, 1
	case DownRight:
		return 1, 1
	}
	return 0, 0
}

// Sign returns the sign of the direction along the axis.
func (d Dir) Sign(a Axis) int {
	sx, sy := d.Signs()
	if a == X {
		return sx
	}
	return sy
}

// Vec returns the unit vector of the direction; zero for None.
func (d Dir) Vec() Vec2 {
	sx, sy := d.Signs()
	if d.Diagonal() {
		return Vec2{X: float64(sx) * diagonal, Y: float64(sy) * diagonal}
	}
	return Vec2{X: float64(sx), Y: float64(sy)}
}

// Diagonal reports whether the direction moves along both axes.
func (d Dir) Diagonal() bool {
	sx, sy := d.Signs()
	return sx != 0 && sy != 0
}

// Axis returns the axis of a straight direction; for a diagonal or None
// the result means nothing.
func (d Dir) Axis() Axis {
	if sx, _ := d.Signs(); sx != 0 {
		return X
	}
	return Y
}

// Opposite returns the direction turned around.
func (d Dir) Opposite() Dir {
	sx, sy := d.Signs()
	return DirOf(-sx, -sy)
}

// String returns the name of the direction.
func (d Dir) String() string {
	switch d {
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	case UpLeft:
		return "UpLeft"
	case UpRight:
		return "UpRight"
	case DownLeft:
		return "DownLeft"
	case DownRight:
		return "DownRight"
	}
	return "None"
}
