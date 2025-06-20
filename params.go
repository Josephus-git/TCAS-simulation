package main

import (
	"time"
)

// Point represents a 3D coordinate
// may be changed to latitude logitude altitude
type coord struct {
	X, Y, Z float64
}

// Airport represents an airport with its location
type airport struct {
	serial        string
	location      coord
	planeCapacity int
	runway        runway
	planes        []plane
}

type runway struct {
	numberOfRunway  int
	noOfRunwayinUse int
}

type plane struct {
	serial        string
	planeInFlight bool
	cruiseSpeed   float64
	flightLog     []flight
}

// Flight represents a single flight from departure to arrival
// *** implement the climb / decent
type flight struct {
	flightID         string
	flightSchedule   flightPath
	takeoffTime      time.Time
	landingTime      time.Time
	cruisingAltitude float64 // Meters
}

// Flight path to store the movement of plane from one location to the other
type flightPath struct {
	depature coord
	arrival  coord
}
