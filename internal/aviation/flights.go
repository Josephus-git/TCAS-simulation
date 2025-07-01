package aviation

import "time"

// Flight represents a single flight from departure to arrival
// *** implement the climb / decent
type Flight struct {
	flightID         string
	flightSchedule   FlightPath
	takeoffTime      time.Time
	landingTime      time.Time
	cruisingAltitude float64 // Meters
}

// FlightPath to store the movement of plane from one location to the other
type FlightPath struct {
	depature Coordinate
	arrival  Coordinate
}
