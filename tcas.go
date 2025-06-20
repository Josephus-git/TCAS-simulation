package main

import (
	"fmt"
	"math"
	"time"
)

// resultant distance obtained by getting the magnitude of distance btw the two coordinates
func distance(p1, p2 coord) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2) + math.Pow(p1.Z-p2.Z, 2))
}

// get the time at which the planes will be closest and the distance at this time
func (f1 flight) GetClosestApproachDetails(f2 flight) (closestTime time.Time, distanceBetweenPlanesatCA float64) {
	// get time When planes will get to coincidence point
	flight1ClosestCoord, flight2ClosestCoord := FindClosestApprachDuringTransit(f1.flightSchedule, f2.flightSchedule)

	flight1Distance := distance(f1.flightSchedule.depature, f1.flightSchedule.arrival)
	distBtwDepatureAndClosestApproachForFlight1 := distance(f1.flightSchedule.depature, flight1ClosestCoord)

	f1fractionofCA := distBtwDepatureAndClosestApproachForFlight1 / flight1Distance

	totalFlightDuration1 := f1.landingTime.Sub(f1.takeoffTime)
	closestTime = f1.takeoffTime.Add(time.Duration(float64(totalFlightDuration1) * f1fractionofCA))

	distanceBetweenPlanesatCA = distance(flight1ClosestCoord, flight2ClosestCoord)

	return closestTime, distanceBetweenPlanesatCA
}

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
	departureLocation := f.flightSchedule.depature
	arivalLocation := f.flightSchedule.arrival

	// Calculate the intermediate point
	pX := departureLocation.X + (departureLocation.X-arivalLocation.X)*progress
	pY := departureLocation.Y + (departureLocation.Y-arivalLocation.Y)*progress
	pZ := f.cruisingAltitude

	return coord{X: pX, Y: pY, Z: pZ}, nil
}
