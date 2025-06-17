package main

import (
	"fmt"
	"time"
)

// Get planes position
func getPlanePosition(f Flight, t time.Time) (Coord, error) {
	if t.Before(f.TakeoffTime) || t.After(f.LandingTime) { //*** return here to check incase plane shold go on another flight
		return Coord{}, fmt.Errorf("time %v is outside flight %s duration", t, f.FlightID)
	}

	// Calculate fraction of flight completed (normalized time 0-1)
	totalDuration := f.LandingTime.Sub(f.TakeoffTime)
	elapsed := t.Sub(f.TakeoffTime)
	progress := float64(elapsed) / float64(totalDuration)

	// get the arival and departure location
	departureLocation := f.Departure.Location
	arivalLocation := f.Arrival.Location

	// Calculate the intermediate point
	pX := departureLocation.X + (departureLocation.X-arivalLocation.X)*progress
	pY := departureLocation.Y + (departureLocation.Y-arivalLocation.Y)*progress
	pZ := f.CruisingAltitude

	return Coord{X: pX, Y: pY, Z: pZ}, nil
}
