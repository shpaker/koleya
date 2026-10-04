package koleya

import "math"

// phase is what an axis is doing.
type phase uint8

const (
	idle phase = iota
	driving
	docking
)

// limits are the effective rates of one axis for one tick.
type limits struct {
	cap   float64 // top speed
	accel float64
	decel float64
	dock  float64 // lowest speed while docking
}

// newLimits scales the profile for one axis: a diagonal splits every rate
// between the two axes.
func newLimits(p Profile, s Surface, scale float64) limits {
	return limits{
		cap:   p.MaxSpeed * scale,
		accel: s.limit(p.Accel) * scale,
		decel: s.limit(p.Decel) * scale,
		dock:  p.dockSpeed() * scale,
	}
}

// axisState is the motion along one axis: the one-dimensional solver
// everything else is built on.
type axisState struct {
	phase  phase
	sign   int     // direction of motion while not idle
	vel    float64 // speed, never negative
	target float64 // the node to rest on while docking
}

// velocity returns the signed velocity.
func (a *axisState) velocity() float64 {
	return float64(a.sign) * a.vel
}

// update moves the axis by one tick towards the wanted velocity (signed,
// zero to let go) and reports whether it came to rest on a node.
func (a *axisState) update(pos *float64, want float64, lim limits, grid Lattice, dt float64) bool {
	wantSign := sign(want)
	switch {
	case wantSign != 0 && (a.phase == idle || a.sign == wantSign):
		a.drive(pos, wantSign, lim, dt)
		return false
	case a.phase == driving:
		// Let go, or the driver turned around: rest first
		a.beginDock(*pos, lim, grid)
	}
	if a.phase != docking {
		return false
	}
	landed, _ := a.dock(pos, lim, dt)
	return landed
}

// drive gains speed up to the cap and moves. A lower cap (a turn into a
// diagonal) takes effect at once: the speed turns, it does not grow.
func (a *axisState) drive(pos *float64, wantSign int, lim limits, dt float64) {
	a.phase = driving
	a.sign = wantSign
	a.vel = rise(a.vel, lim.accel, dt, lim.cap)
	*pos += float64(a.sign) * a.vel * dt
}

// beginDock picks the node to rest on: the first one beyond the braking
// distance.
func (a *axisState) beginDock(pos float64, lim limits, grid Lattice) {
	braking := 0.0
	if !math.IsInf(lim.decel, 1) {
		braking = a.vel * a.vel / (2 * lim.decel)
	}
	a.target = grid.NodeAhead(pos+float64(a.sign)*braking, a.sign)
	a.phase = docking
}

// remaining returns the distance left to the node.
func (a *axisState) remaining(pos float64) float64 {
	return math.Abs(a.target - pos)
}

// dock brakes so the speed would run out right on the node, never below
// the docking speed, and lands on the node exactly. It returns whether the
// axis landed and the distance covered.
func (a *axisState) dock(pos *float64, lim limits, dt float64) (bool, float64) {
	dist := a.remaining(*pos)
	if a.vel < lim.dock {
		a.vel = rise(a.vel, lim.accel, dt, lim.dock)
	} else if dist > 0 {
		need := a.vel * a.vel / (2 * dist)
		a.vel = math.Max(lim.dock, a.vel-need*dt)
	}
	travel := a.vel * dt
	if travel >= dist {
		a.land(pos)
		return true, dist
	}
	*pos += float64(a.sign) * travel
	return false, travel
}

// land puts the axis on its node at rest.
func (a *axisState) land(pos *float64) {
	*pos = a.target
	a.halt()
}

// halt stops the axis where it is.
func (a *axisState) halt() {
	a.phase = idle
	a.sign = 0
	a.vel = 0
}

// rise gains speed at the rate for dt, up to top; an Instant rate reaches
// top at once, even in a tick of zero length.
func rise(vel, rate, dt, top float64) float64 {
	if math.IsInf(rate, 1) {
		return top
	}
	return math.Min(top, vel+rate*dt)
}

// sign returns -1, 0 or 1.
func sign(v float64) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}
