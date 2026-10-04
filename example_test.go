package koleya_test

import (
	"fmt"

	"github.com/shpaker/koleya"
)

func Example() {
	grid := koleya.Lattice{Step: 4}
	profile := koleya.Profile{
		Steering:  koleya.Four,
		MaxSpeed:  72,  // px/s
		Accel:     260, // px/s²; koleya.Instant for no ramp-up
		Decel:     220,
		DockSpeed: 20, // the last stretch to the node never gets slower
	}
	m := koleya.NewMover(koleya.Vec2{X: 0, Y: 0}, koleya.Up, grid)

	// Hold Right for half a second, then let go
	for tick := range 60 {
		held := koleya.None
		if tick < 30 {
			held = koleya.Right
		}

		prev := m.Pos
		ev := koleya.Step(&m, koleya.Intent{Dir: held}, profile, koleya.Ground, 1.0/60)

		// Collisions are the game's: roll back and halt
		if m.Pos.X > 1000 {
			m.Pos = prev
			m.Halt()
		}
		if ev.Has(koleya.Started) {
			fmt.Println("engine on, facing", m.Facing())
		}
		if ev.Has(koleya.Stopped) {
			fmt.Printf("engine off at %.1f\n", m.Pos.X)
		}
	}
	// Output:
	// engine on, facing Right
	// engine off at 40.0
}
