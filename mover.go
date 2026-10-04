package koleya

// Mover is the motion state of one vehicle. It is a plain value: the game
// keeps it inside its entity and passes it to [Step] every tick.
type Mover struct {
	// Pos is the position in pixels. The game may roll it back after a
	// collision, followed by Halt or HaltAxis.
	Pos Vec2

	lattice Lattice
	facing  Dir
	pending Dir
	axes    [2]axisState
}

// NewMover returns a vehicle at rest at pos, facing the given direction,
// that comes to rest on the nodes of the lattice. pos should lie on a node.
func NewMover(pos Vec2, facing Dir, lattice Lattice) Mover {
	return Mover{Pos: pos, lattice: lattice, facing: facing}
}

// Intent is what the driver wants this tick.
type Intent struct {
	// Dir is the direction to move in; None lets go.
	Dir Dir
}

// Facing returns the direction the vehicle faces: where it moves, or
// where it moved last.
func (m *Mover) Facing() Dir {
	return m.facing
}

// Pending returns the direction a Four vehicle sets off in once it has
// come to rest on a node; None when there is no such turn.
func (m *Mover) Pending() Dir {
	return m.pending
}

// Lattice returns the grid the vehicle comes to rest on.
func (m *Mover) Lattice() Lattice {
	return m.lattice
}

// Vel returns the velocity in pixels per second.
func (m *Mover) Vel() Vec2 {
	return Vec2{X: m.axes[X].velocity(), Y: m.axes[Y].velocity()}
}

// Moving reports whether the vehicle is not at rest.
func (m *Mover) Moving() bool {
	return m.axes[X].phase != idle || m.axes[Y].phase != idle
}

// Docking reports whether the vehicle is braking towards a node.
func (m *Mover) Docking() bool {
	return m.axes[X].phase == docking || m.axes[Y].phase == docking
}

// Target returns the node the vehicle comes to rest on along the axis
// while docking.
func (m *Mover) Target(a Axis) (float64, bool) {
	if m.axes[a].phase != docking {
		return 0, false
	}
	return m.axes[a].target, true
}

// Face turns a vehicle at rest to the direction, in place. It reports
// whether the vehicle turned; a moving vehicle does not.
func (m *Mover) Face(d Dir) bool {
	if d == None || d == m.facing || m.Moving() {
		return false
	}
	m.facing = d
	return true
}

// Halt stops the vehicle where it is, without coming to rest on a node:
// it ran into something.
func (m *Mover) Halt() {
	m.axes[X].halt()
	m.axes[Y].halt()
	m.pending = None
}

// HaltAxis stops the motion along one axis where it is, without coming to
// rest on a node; the other axis keeps going, so an Eight vehicle slides
// along a wall.
func (m *Mover) HaltAxis(a Axis) {
	m.axes[a].halt()
	if !m.Moving() {
		m.pending = None
	}
}
