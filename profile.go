package koleya

import (
	"errors"
	"math"
)

// Steering is how a vehicle turns the driver's direction into motion.
type Steering uint8

const (
	// Four moves along one axis at a time. To turn, the vehicle first
	// comes to rest on a node of the current axis, then sets off along
	// the new one; a diagonal intent keeps the current heading when it
	// is part of the diagonal and takes the vertical part otherwise.
	Four Steering = iota
	// Eight moves along both axes at once, diagonals included. An axis
	// the driver lets go of comes to rest on a node while the other one
	// keeps going; when both are let go at once, they reach their nodes
	// together, along a straight line.
	Eight
)

// Profile is the tuning of a vehicle. Speeds are in pixels per second,
// accelerations in pixels per second squared.
type Profile struct {
	// Steering is the set of directions the vehicle moves in.
	Steering Steering
	// MaxSpeed is the top speed.
	MaxSpeed float64
	// Accel is how fast the vehicle gains speed; Instant for no ramp-up.
	Accel float64
	// Decel is how hard the vehicle brakes. A softer brake means a longer
	// braking distance and a farther node to stop on; Instant stops on
	// the first node ahead.
	Decel float64
	// DockSpeed is the lowest speed of the last stretch to the node, so
	// the vehicle does not crawl the final pixels. Equal to MaxSpeed, it
	// rolls to the node at full speed. Zero or less means MaxSpeed.
	DockSpeed float64
}

// Classic is the feel of the 8-bit tank games: four directions, full speed
// at once, and rolling to the first node ahead at full speed.
func Classic(speed float64) Profile {
	return Profile{
		Steering:  Four,
		MaxSpeed:  speed,
		Accel:     Instant,
		Decel:     Instant,
		DockSpeed: speed,
	}
}

// Validate reports a profile a vehicle cannot move with.
func (p Profile) Validate() error {
	switch {
	case !(p.MaxSpeed > 0) || math.IsInf(p.MaxSpeed, 0):
		return errors.New("koleya: MaxSpeed must be positive and finite")
	case !(p.Accel > 0):
		return errors.New("koleya: Accel must be positive")
	case !(p.Decel > 0):
		return errors.New("koleya: Decel must be positive")
	case p.Steering != Four && p.Steering != Eight:
		return errors.New("koleya: unknown Steering")
	}
	return nil
}

// dockSpeed returns the effective DockSpeed.
func (p Profile) dockSpeed() float64 {
	if p.DockSpeed <= 0 || p.DockSpeed > p.MaxSpeed {
		return p.MaxSpeed
	}
	return p.DockSpeed
}

// Surface is the ground under the vehicle this tick.
type Surface struct {
	// Grip caps both the acceleration and the braking, in pixels per
	// second squared: on ice the vehicle sets off slower and slides
	// farther, and a faster vehicle slides farther than a slow one.
	// Zero means no cap: solid ground.
	Grip float64
}

// Ground is solid ground: the profile alone decides.
var Ground = Surface{}

// limit caps a rate by the grip.
func (s Surface) limit(rate float64) float64 {
	if s.Grip > 0 && s.Grip < rate {
		return s.Grip
	}
	return rate
}
