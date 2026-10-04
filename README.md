# koleya

*Old-school looks, modern feel.*

koleya is grid-aligned movement for games that look old but feel new. Vehicles glide,
accelerate, and coast with modern smoothness, yet always come to rest exactly on the
grid, so they slip into corridors just as easily as they did in the 8-bit classics.
Pixel positions are fractional, so it works at any modern screen resolution.

*koleya (колея) — a rut, a track worn into the road.*

## Used in

- [tnk9x](https://github.com/shpaker/tnk9x) — tanks from the 90s the way we
  remember them, with real-time 2D lighting; the `Classic` feel.

## Quick start

```
go get github.com/shpaker/koleya
```

```go
grid := koleya.Lattice{Step: 4}
profile := koleya.Profile{
	Steering:  koleya.Four,
	MaxSpeed:  72,  // px/s
	Accel:     260, // px/s²; koleya.Instant for no ramp-up
	Decel:     220,
	DockSpeed: 20, // the last stretch to the node never gets slower
}
m := koleya.NewMover(koleya.Vec2{X: 0, Y: 0}, koleya.Up, grid)

// Every tick
prev := m.Pos
ev := koleya.Step(&m, koleya.Intent{Dir: held}, profile, koleya.Ground, dt)

// Collisions are the game's: roll back and halt
if world.Blocked(m.Pos) {
	m.Pos = prev
	m.Halt()
}
if ev.Has(koleya.Started) {
	// engine on
}
if ev.Has(koleya.Turned) {
	// sprite for m.Facing()
}
```

The game keeps the `Mover` inside its entity and the `Profile` as data: a light
tank, a heavy tank and a truck are three profiles.

## How it moves

- **Driving.** Holding a direction gains speed at `Accel` up to `MaxSpeed`.
- **Docking.** Letting go picks the first node beyond the braking distance and
  brakes so the speed runs out right on it, never slower than `DockSpeed`. A soft
  `Decel` stops a node or two later; `Instant` stops on the first node ahead.
- **Four.** One axis at a time. To turn, the vehicle first rests on a node, then sets
  off the new way — the axis it does not move along always lies on a node.
- **Eight.** Diagonals too, with the speed split between the axes. An axis the driver
  lets go of docks while the other keeps going; letting go of both brakes along the
  straight line to the node, and both axes land together.
- **Ice.** `Surface{Grip: 300}` caps acceleration and braking for the tick: the vehicle
  sets off slower and slides farther, a fast one farther than a slow one — and still
  rests on a node.
- **Walls.** `Halt` stops in place, off the grid; `HaltAxis` stops one axis so an
  Eight vehicle slides along a wall.
- **Turning in place.** `Face` turns a vehicle at rest without moving it.

`koleya.Classic(speed)` is the feel of the 8-bit tank games: four directions, full speed
at once, rolling to the first node ahead at full speed.

Helpers: `Lattice.NodeAhead`, `Lattice.Nearest`, `Lattice.Near` (e.g. an AI that turns
only near a node: `koleya.Lattice{Step: 8}.Near(pos.X, 2)`).

## The promise

At rest, unless halted, both coordinates lie on lattice nodes, and the speed never
exceeds `MaxSpeed` — for any profile, any input and any frame times. A fuzz test
(`just fuzz`) holds koleya to it.

## License

[MIT](LICENSE)
