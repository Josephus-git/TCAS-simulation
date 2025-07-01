package aviation

import "time"

// Flight represents a single flight from departure to arrival
// *** implement the climb / decent
type Flight struct {
	FlightID         string
	FlightSchedule   FlightPath
	TakeoffTime      time.Time
	LandingTime      time.Time
	CruisingAltitude float64 // Meters
}

// FlightPath to store the movement of plane from one location to the other
type FlightPath struct {
	Depature Coordinate
	Arrival  Coordinate
}
