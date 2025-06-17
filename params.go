package main

import (
	"time"
)

// Point represents a 3D coordinate
type Coord struct {
	X, Y, Z float64
}

// Airport represents an airport with its location
type Airport struct {
	Serial        string
	Location      Coord
	PlaneCapacity int
	Runway        Runway
}

type Runway struct {
	NumberOfRunway  int
	NoOfRunwayinUse int
}

type Plane struct {
	Serial        string
	PlaneInFlight bool
	Speed         float64
}

// Flight represents a single flight from departure to arrival
// *** implement the climb / decent
type Flight struct {
	FlightID         string // inthe format {fromairport/toairport/index in digit}
	Departure        Airport
	Arrival          Airport
	TakeoffTime      time.Time
	LandingTime      time.Time
	CruisingAltitude float64 // Meters
}

// PlaneState represents the position and time of a plane
type Planelocation struct {
	Point Coord
	Time  time.Time
}

// CoincidenceResult to contain
type CoincidenceResult struct {
	Flight1     Flight
	Flight2     Flight
	ClosestTime time.Time
	MinDistance float64
}
