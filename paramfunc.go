package main

import (
	"fmt"
	"time"
)

// Get planes position
func (p plane) getPlanePosition(f flight, t time.Time) (coord, error) {
	if t.Before(f.takeoffTime) || t.After(f.landingTime) { //*** return here to check incase plane shold go on another flight
		return coord{}, fmt.Errorf("time %v is outside flight %s duration", t, f.flightID)
	}

	// Calculate fraction of flight completed (normalized time 0-1)
	totalDuration := f.landingTime.Sub(f.takeoffTime)
	elapsed := t.Sub(f.takeoffTime)
	progress := float64(elapsed) / float64(totalDuration)

	// get the arival and departure location
	departureLocation := f.departure.location
	arivalLocation := f.arrival.location

	// Calculate the intermediate point
	pX := departureLocation.X + (departureLocation.X-arivalLocation.X)*progress
	pY := departureLocation.Y + (departureLocation.Y-arivalLocation.Y)*progress
	pZ := f.cruisingAltitude

	return coord{X: pX, Y: pY, Z: pZ}, nil
}
