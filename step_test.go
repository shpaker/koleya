package koleya

import (
	"math"
	"testing"
)

const dt = 1.0 / 60

var grid = Lattice{Step: 4}

// run steps the mover with the same intent until it rests or the ticks
// run out, and returns the events of all ticks together.
func run(m *Mover, d Dir, p Profile, s Surface, ticks int) Events {
	var all Events
	for range ticks {
		all |= Step(m, Intent{Dir: d}, p, s, dt)
		if !m.Moving() {
			break
		}
	}
	return all
}

func TestClassicDocksToNodeAhead(t *testing.T) {
	m := NewMover(Vec2{X: 0, Y: 0}, Right, grid)
	p := Classic(48)

	// 7 ticks at 0.8 px: 5.6 px
	for range 7 {
		Step(&m, Intent{Dir: Right}, p, Ground, dt)
	}
	if math.Abs(m.Pos.X-5.6) > 1e-9 {
		t.Fatalf("X = %v, want 5.6", m.Pos.X)
	}

	ev := Step(&m, Intent{}, p, Ground, dt)
	if !m.Docking() || m.Vel().X != 48 {
		t.Fatalf("the first docking tick keeps full speed, vel = %v", m.Vel())
	}
	ev |= run(&m, None, p, Ground, 100)
	if m.Pos.X != 8 || m.Moving() {
		t.Fatalf("rests at %v, want 8", m.Pos.X)
	}
	if !ev.Has(Docked | Stopped) {
		t.Errorf("events %b, want Docked and Stopped", ev)
	}
}

func TestClassicStartsAtFullSpeed(t *testing.T) {
	m := NewMover(Vec2{}, Up, grid)
	ev := Step(&m, Intent{Dir: Right}, Classic(48), Ground, dt)
	if m.Vel().X != 48 {
		t.Errorf("vel = %v, want 48 at once", m.Vel())
	}
	if !ev.Has(Started | Turned) {
		t.Errorf("events %b, want Started and Turned", ev)
	}
}

func TestFourTurnsAfterDocking(t *testing.T) {
	m := NewMover(Vec2{}, Right, grid)
	p := Classic(48)
	for range 3 {
		Step(&m, Intent{Dir: Right}, p, Ground, dt)
	}

	// Turn up: X rolls to 4 first, Y does not move meanwhile
	for m.Facing() == Right {
		Step(&m, Intent{Dir: Up}, p, Ground, dt)
		if m.Pos.Y != 0 {
			t.Fatalf("Y moved before X came to rest: %v", m.Pos)
		}
	}
	if m.Pos.X != 4 || m.Pending() != None {
		t.Fatalf("turned at %v, pending %v", m.Pos, m.Pending())
	}

	Step(&m, Intent{Dir: Up}, p, Ground, dt)
	if m.Pos.Y >= 0 || m.Pos.X != 4 {
		t.Errorf("after the turn moves up from X = 4: %v", m.Pos)
	}
}

func TestFourTurnsAroundAfterDocking(t *testing.T) {
	m := NewMover(Vec2{}, Right, grid)
	p := Classic(48)
	Step(&m, Intent{Dir: Right}, p, Ground, dt)

	for m.Facing() == Right {
		Step(&m, Intent{Dir: Left}, p, Ground, dt)
	}
	if m.Pos.X != 4 {
		t.Fatalf("turned around at %v, want 4", m.Pos.X)
	}
}

func TestFourMoveCancelsDocking(t *testing.T) {
	m := NewMover(Vec2{}, Right, grid)
	p := Classic(48)
	Step(&m, Intent{Dir: Right}, p, Ground, dt)
	Step(&m, Intent{}, p, Ground, dt)
	if !m.Docking() {
		t.Fatal("letting go docks")
	}
	Step(&m, Intent{Dir: Right}, p, Ground, dt)
	if m.Docking() || !m.Moving() {
		t.Error("pressing again drives on")
	}
}

func TestFourLettingGoDropsTheTurn(t *testing.T) {
	m := NewMover(Vec2{}, Right, grid)
	p := Classic(48)
	Step(&m, Intent{Dir: Right}, p, Ground, dt)
	Step(&m, Intent{Dir: Up}, p, Ground, dt)
	run(&m, None, p, Ground, 100)
	if m.Facing() != Right || m.Pos.X != 4 || m.Pos.Y != 0 {
		t.Errorf("rests at %v facing %v, want (4, 0) facing Right", m.Pos, m.Facing())
	}
}

func TestFourDiagonalKeepsHeading(t *testing.T) {
	m := NewMover(Vec2{}, Right, grid)
	p := Classic(48)
	Step(&m, Intent{Dir: Right}, p, Ground, dt)
	Step(&m, Intent{Dir: UpRight}, p, Ground, dt)
	if m.Facing() != Right || m.Docking() {
		t.Errorf("a diagonal with the heading keeps driving, facing %v", m.Facing())
	}
}

func TestBrakingDistancePicksFartherNode(t *testing.T) {
	p := Profile{Steering: Four, MaxSpeed: 60, Accel: Instant, Decel: 100, DockSpeed: 5}
	m := NewMover(Vec2{}, Right, grid)
	Step(&m, Intent{Dir: Right}, p, Ground, dt) // X = 1, braking distance 18

	Step(&m, Intent{}, p, Ground, dt)
	target, ok := m.Target(X)
	if !ok || target != 20 {
		t.Fatalf("target = %v, want 20: the first node beyond 1 + 18", target)
	}

	prev := m.Vel().X
	for m.Moving() {
		Step(&m, Intent{}, p, Ground, dt)
		if v := m.Vel().X; v > prev {
			t.Fatalf("speeds up while braking: %v -> %v", prev, v)
		} else {
			prev = v
		}
	}
	if m.Pos.X != 20 {
		t.Errorf("rests at %v, want 20", m.Pos.X)
	}
}

func TestAccelerationRampsUp(t *testing.T) {
	p := Profile{Steering: Four, MaxSpeed: 60, Accel: 120, Decel: Instant, DockSpeed: 60}
	m := NewMover(Vec2{}, Right, grid)
	Step(&m, Intent{Dir: Right}, p, Ground, dt)
	if v := m.Vel().X; math.Abs(v-2) > 1e-9 {
		t.Errorf("vel after one tick = %v, want 2", v)
	}
	for range 60 {
		Step(&m, Intent{Dir: Right}, p, Ground, dt)
	}
	if m.Vel().X != 60 {
		t.Errorf("vel after a second = %v, want 60", m.Vel().X)
	}
}

func TestIceSlidesFarther(t *testing.T) {
	rest := func(s Surface, speed float64) float64 {
		p := Classic(speed)
		m := NewMover(Vec2{}, Right, grid)
		for range 30 { // up to full speed even on ice
			Step(&m, Intent{Dir: Right}, p, s, dt)
		}
		run(&m, None, p, s, 1000)
		return m.Pos.X
	}
	ice := Surface{Grip: 288}

	ground, slow, fast := rest(Ground, 48), rest(ice, 48), rest(ice, 72)
	if !(slow > ground) {
		t.Errorf("on ice rests at %v, on ground at %v: want farther", slow, ground)
	}
	if !(fast > slow) {
		t.Errorf("a fast tank rests at %v, a slow one at %v: want farther", fast, slow)
	}
	for _, x := range []float64{ground, slow, fast} {
		if !grid.On(x) {
			t.Errorf("rests at %v, off the grid", x)
		}
	}
}

func TestHaltStopsInPlace(t *testing.T) {
	m := NewMover(Vec2{}, Right, grid)
	p := Classic(48)
	Step(&m, Intent{Dir: Right}, p, Ground, dt)
	Step(&m, Intent{Dir: Up}, p, Ground, dt)
	m.Halt()
	if m.Moving() || m.Pending() != None || m.Pos.X != 1.6 {
		t.Errorf("halted: moving %v, pending %v, at %v", m.Moving(), m.Pending(), m.Pos)
	}
}

func TestFaceTurnsOnlyAtRest(t *testing.T) {
	m := NewMover(Vec2{}, Right, grid)
	if !m.Face(Up) || m.Facing() != Up {
		t.Error("a vehicle at rest turns in place")
	}
	Step(&m, Intent{Dir: Up}, Classic(48), Ground, dt)
	if m.Face(Left) {
		t.Error("a moving vehicle does not turn in place")
	}
}

func TestEightDiagonalSpeed(t *testing.T) {
	p := Classic(48)
	p.Steering = Eight
	m := NewMover(Vec2{}, Up, grid)
	Step(&m, Intent{Dir: DownRight}, p, Ground, dt)
	v := m.Vel()
	if math.Abs(math.Hypot(v.X, v.Y)-48) > 1e-9 || v.X != v.Y {
		t.Errorf("diagonal vel = %v, want 48 split evenly", v)
	}
}

func TestEightDocksBothAxesTogether(t *testing.T) {
	p := Profile{Steering: Eight, MaxSpeed: 60, Accel: Instant, Decel: 200, DockSpeed: 5}
	m := NewMover(Vec2{}, Right, grid)
	for range 5 {
		Step(&m, Intent{Dir: DownRight}, p, Ground, dt)
	}
	m.Pos.Y += 1.5 // the axes are now at different distances from their nodes

	Step(&m, Intent{}, p, Ground, dt)
	for m.Moving() {
		if !m.Docking() || m.Vel().X == 0 || m.Vel().Y == 0 {
			t.Fatalf("an axis came to rest before the other: vel %v", m.Vel())
		}
		Step(&m, Intent{}, p, Ground, dt)
	}
	if !grid.On(m.Pos.X) || !grid.On(m.Pos.Y) {
		t.Errorf("rests at %v, off the grid", m.Pos)
	}
}

func TestEightSlidesAlongWall(t *testing.T) {
	p := Classic(48)
	p.Steering = Eight
	m := NewMover(Vec2{}, Right, grid)
	Step(&m, Intent{Dir: UpRight}, p, Ground, dt)
	m.Pos.Y = 0
	m.HaltAxis(Y)

	x := m.Pos.X
	Step(&m, Intent{Dir: UpRight}, p, Ground, dt)
	m.Pos.Y = 0 // the wall again
	m.HaltAxis(Y)
	if !m.Moving() || m.Pos.X <= x {
		t.Errorf("keeps going along X: %v", m.Pos)
	}
}

func TestEightLettingGoOfOneAxis(t *testing.T) {
	p := Classic(48)
	p.Steering = Eight
	m := NewMover(Vec2{}, Right, grid)
	Step(&m, Intent{Dir: DownRight}, p, Ground, dt)
	for range 20 {
		Step(&m, Intent{Dir: Right}, p, Ground, dt)
	}
	if m.Vel().Y != 0 || !grid.On(m.Pos.Y) || m.Vel().X != 48 {
		t.Errorf("Y rests on a node and X drives at full speed: pos %v, vel %v", m.Pos, m.Vel())
	}
}

func TestValidate(t *testing.T) {
	if err := Classic(48).Validate(); err != nil {
		t.Errorf("Classic: %v", err)
	}
	if err := (Profile{MaxSpeed: 10, Accel: 0, Decel: 1}).Validate(); err == nil {
		t.Error("zero Accel is invalid")
	}
}
