// Package koleya moves things along a grid the way classic games did, with
// a modern feel.
//
// A vehicle accelerates, coasts and brakes, yet always comes to rest exactly
// on a node of its [Lattice], so it slips into corridors as easily as in the
// 8-bit classics. Braking picks the first node beyond the braking distance and
// tunes the deceleration to stop right on it.
//
// The game owns the state ([Mover]) and the tuning ([Profile]) and calls
// [Step] every tick with what the driver wants ([Intent]) and the ground
// under the vehicle ([Surface]). Collisions stay with the game: it rolls the
// position back and calls [Mover.Halt] or [Mover.HaltAxis].
package koleya

import "math"

// Instant is an acceleration or a deceleration with no limit: the speed
// changes within one tick.
var Instant = math.Inf(1)

// eps is the tolerance of comparisons with lattice nodes, in units of
// the lattice step.
const eps = 1e-9

// Events are the facts of one [Step], as bit flags.
type Events uint8

const (
	// Started: the vehicle was at rest and began to move.
	Started Events = 1 << iota
	// Turned: the vehicle faces a new direction.
	Turned
	// Docked: an axis came to rest on a lattice node.
	Docked
	// Stopped: the vehicle was moving and came to rest.
	Stopped
)

// Has reports whether all the flags of f are set.
func (e Events) Has(f Events) bool {
	return e&f == f
}
