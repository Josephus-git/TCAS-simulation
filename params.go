package main

import (
	"time"
)

// Point represents a 3D coordinate
type coord struct {
	X, Y, Z float64
}

// Airport represents an airport with its location
type airport struct {
	serial        string
	location      coord
	planeCapacity int
	runway        runway
}

type runway struct {
	numberOfRunway  int
	noOfRunwayinUse int
}

type plane struct {
	serial        string
	planeInFlight bool
	speed         float64
	location      planeLocation
}

// Flight represents a single flight from departure to arrival
// *** implement the climb / decent
type flight struct {
	flightID         string // inthe format {fromairport/toairport/index in digit}
	departure        airport
	arrival          airport
	takeoffTime      time.Time
	landingTime      time.Time
	cruisingAltitude float64 // Meters
}

// PlaneState represents the position and time of a plane
type planeLocation struct {
	point coord
	time  time.Time
}

// CoincidenceResult to contain
type coincidenceResult struct {
	flight1     flight
	flight2     flight
	closestTime time.Time
	minDistance float64
}
