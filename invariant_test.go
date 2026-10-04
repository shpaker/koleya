package koleya

import (
	"math"
	"math/rand/v2"
	"testing"
)

// FuzzRestsOnGrid drives random vehicles with random input and checks the
// promise of the library: at rest, both coordinates lie on lattice nodes;
// a Four vehicle keeps the idle axis on a node while it moves; speeds stay
// within the profile.
func FuzzRestsOnGrid(f *testing.F) {
	f.Add(uint64(1), uint8(0))
	f.Add(uint64(2), uint8(1))
	f.Add(uint64(3), uint8(2))
	f.Fuzz(func(t *testing.T, seed uint64, kind uint8) {
		rng := rand.New(rand.NewPCG(seed, uint64(kind)))
		p, grid, s := randomSetup(rng, kind)
		if p.Validate() != nil {
			t.Skip()
		}

		m := NewMover(Vec2{X: grid.Step * 3, Y: grid.Step * 5}, Up, grid)
		d := None
		for tick := range 2000 {
			if rng.IntN(20) == 0 {
				d = Dir(rng.IntN(9))
			}
			step := dt
			if kind%3 == 2 {
				step = rng.Float64() / 20 // uneven frames
			}
			if rng.IntN(50) == 0 {
				step = 0 // a frame with no time
			}
			Step(&m, Intent{Dir: d}, p, s, step)

			v := m.Vel()
			if math.IsNaN(m.Pos.X) || math.IsNaN(m.Pos.Y) {
				t.Fatalf("tick %d: NaN position", tick)
			}
			if math.Hypot(v.X, v.Y) > p.MaxSpeed*(1+1e-9) {
				t.Fatalf("tick %d: speed %v above %v", tick, v, p.MaxSpeed)
			}
			if !m.Moving() && !(grid.On(m.Pos.X) && grid.On(m.Pos.Y)) {
				t.Fatalf("tick %d: rests at %v, off the grid", tick, m.Pos)
			}
			if p.Steering == Four && m.Moving() {
				a, _ := m.activeAxis()
				if !grid.On(m.Pos.On(a.Other())) {
					t.Fatalf("tick %d: moves along %v with %v off the grid", tick, a, m.Pos)
				}
			}
		}

		// Letting go always ends on the grid
		for range 100000 {
			if !m.Moving() {
				break
			}
			Step(&m, Intent{}, p, s, dt)
		}
		if m.Moving() || !grid.On(m.Pos.X) || !grid.On(m.Pos.Y) {
			t.Fatalf("after letting go: moving %v at %v", m.Moving(), m.Pos)
		}
	})
}

func randomSetup(rng *rand.Rand, kind uint8) (Profile, Lattice, Surface) {
	rate := func() float64 {
		if rng.IntN(3) == 0 {
			return Instant
		}
		return 20 + rng.Float64()*500
	}
	p := Profile{
		Steering:  Steering(kind % 2),
		MaxSpeed:  10 + rng.Float64()*150,
		Accel:     rate(),
		Decel:     rate(),
		DockSpeed: rng.Float64() * 80,
	}
	grid := Lattice{Step: float64(1 + rng.IntN(16))}
	s := Ground
	if rng.IntN(2) == 0 {
		s = Surface{Grip: 10 + rng.Float64()*300}
	}
	return p, grid, s
}
