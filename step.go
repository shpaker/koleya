package koleya

import "math"

// Step moves the vehicle by one tick of dt seconds and returns what
// happened. It reads nothing but its arguments.
func Step(m *Mover, in Intent, p Profile, s Surface, dt float64) Events {
	wasMoving := m.Moving()
	facing := m.facing

	var landed bool
	if p.Steering == Eight {
		landed = stepEight(m, in.Dir, p, s, dt)
	} else {
		landed = stepFour(m, in.Dir, p, s, dt)
	}

	var ev Events
	if landed {
		ev |= Docked
	}
	if m.facing != facing {
		ev |= Turned
	}
	switch moving := m.Moving(); {
	case !wasMoving && moving:
		ev |= Started
	case wasMoving && !moving:
		ev |= Stopped
	}
	return ev
}

// stepFour moves along one axis at a time.
func stepFour(m *Mover, d Dir, p Profile, s Surface, dt float64) bool {
	d = m.straight(d)
	lim := newLimits(p, s, 1)

	active, ok := m.activeAxis()
	if !ok {
		m.pending = None
		if d == None {
			return false
		}
		m.facing = d
		active = d.Axis()
	}

	want := 0.0
	if d == m.facing {
		want = float64(d.Sign(active)) * lim.cap
		m.pending = None
	} else {
		m.pending = d
	}

	landed := m.axes[active].update(m.Pos.at(active), want, lim, m.lattice, dt)
	if landed && m.pending != None {
		// Set off along the new direction from the next tick
		m.facing = m.pending
		m.pending = None
		next := &m.axes[m.facing.Axis()]
		next.phase = driving
		next.sign = m.facing.Sign(m.facing.Axis())
	}
	return landed
}

// straight reduces a diagonal for Four steering: the current heading when
// it is part of the diagonal, the vertical part otherwise.
func (m *Mover) straight(d Dir) Dir {
	if !d.Diagonal() {
		return d
	}
	if !m.facing.Diagonal() && m.facing != None {
		a := m.facing.Axis()
		if m.facing.Sign(a) == d.Sign(a) {
			return m.facing
		}
	}
	return DirOf(0, d.Sign(Y))
}

// activeAxis returns the axis a Four vehicle moves along.
func (m *Mover) activeAxis() (Axis, bool) {
	switch {
	case m.axes[X].phase != idle:
		return X, true
	case m.axes[Y].phase != idle:
		return Y, true
	}
	return X, false
}

// stepEight moves along both axes at once.
func stepEight(m *Mover, d Dir, p Profile, s Surface, dt float64) bool {
	if d != None {
		m.facing = d
	}
	m.pending = None

	sx, sy := d.Signs()
	signs := [2]int{sx, sy}

	// A diagonal splits the rates between the axes
	scale := 1.0
	if (sx != 0 || m.axes[X].phase != idle) && (sy != 0 || m.axes[Y].phase != idle) {
		scale = diagonal
	}
	lim := newLimits(p, s, scale)

	if m.brakes(X, signs[X]) && m.brakes(Y, signs[Y]) {
		return m.dockTogether(p, s, lim, dt)
	}

	landed := false
	for _, a := range [2]Axis{X, Y} {
		axisLim := lim
		if other := &m.axes[a.Other()]; m.brakes(a.Other(), signs[a.Other()]) {
			// The other axis brakes: speed it still has is not ours to gain
			used := max(other.vel, lim.dock)
			axisLim.cap = min(lim.cap, math.Sqrt(max(0, p.MaxSpeed*p.MaxSpeed-used*used)))
		}
		want := float64(signs[a]) * lim.cap
		if m.axes[a].update(m.Pos.at(a), want, axisLim, m.lattice, dt) {
			landed = true
		}
	}
	return landed
}

// brakes reports whether the axis is moving and the driver no longer wants
// it to go on the same way: it comes to rest on a node this tick.
func (m *Mover) brakes(a Axis, want int) bool {
	return m.axes[a].phase != idle && m.axes[a].sign != want
}

// dockTogether brings both axes to rest along the straight line to their
// nodes, as one braking along that line: they land on the same tick and
// the speed never grows.
func (m *Mover) dockTogether(p Profile, s Surface, lim limits, dt float64) bool {
	for _, a := range [2]Axis{X, Y} {
		if m.axes[a].phase == driving {
			m.axes[a].beginDock(*m.Pos.at(a), lim, m.lattice)
		}
	}

	dx := m.axes[X].target - m.Pos.X
	dy := m.axes[Y].target - m.Pos.Y
	dist := math.Hypot(dx, dy)
	line := axisState{
		phase:  docking,
		sign:   1,
		vel:    math.Hypot(m.axes[X].vel, m.axes[Y].vel),
		target: dist,
	}
	var along float64
	landed, travel := line.dock(&along, newLimits(p, s, 1), dt)
	if landed || dist == 0 {
		m.axes[X].land(&m.Pos.X)
		m.axes[Y].land(&m.Pos.Y)
		return true
	}

	share := travel / dist
	m.Pos.X += dx * share
	m.Pos.Y += dy * share
	m.axes[X].vel = line.vel * math.Abs(dx) / dist
	m.axes[Y].vel = line.vel * math.Abs(dy) / dist
	return false
}
