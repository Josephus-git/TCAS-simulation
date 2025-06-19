package main

import (
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
