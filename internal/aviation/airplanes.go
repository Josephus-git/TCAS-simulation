package aviation

import (
	"fmt"
	"math"
	"time"

	"github.com/josephus-git/TCAS-simulation/internal/util"
)

// Plane represents an aircraft with its key operational details and flight history.
type Plane struct {
	Serial        string
	PlaneInFlight bool
	CruiseSpeed   float64
	FlightLog     []Flight
}

// createPlane initializes and returns a new Plane struct with a generated serial number.
func createPlane(planeCount int) Plane {
	return Plane{
		Serial:        util.GenerateSerialNumber(planeCount, "p"),
		PlaneInFlight: false,
		CruiseSpeed:   5,
		FlightLog:     []Flight{},
	}
}

// getPlanePosition calculates the plane's interpolated coordinates (X, Y, Z) at a given time during a specific flight.
// It returns an error if the provided time is outside the flight's duration.
func (p Plane) getPlanePosition(f Flight, t time.Time) (Coordinate, error) {
	if t.Before(f.TakeoffTime) || t.After(f.ExpectedLandingTime) { //*** return here to check incase plane shold go on another Flight
		return Coordinate{}, fmt.Errorf("time %v is outside Flight %s duration", t, f.FlightID)
	}

	// Calculate fraction of Flight completed (normalized time 0-1)
	totalDuration := f.ExpectedLandingTime.Sub(f.TakeoffTime)
	elapsed := t.Sub(f.TakeoffTime)
	progress := float64(elapsed) / float64(totalDuration)

	// get the arival and departure location
	departureLocation := f.FlightSchedule.Depature
	arivalLocation := f.FlightSchedule.Destination

	// Calculate the intermediate point
	pX := departureLocation.X + (departureLocation.X-arivalLocation.X)*progress
	pY := departureLocation.Y + (departureLocation.Y-arivalLocation.Y)*progress
	pZ := f.CruisingAltitude

	return Coordinate{X: pX, Y: pY, Z: pZ}, nil
}

// Distance calculates the Euclidean Distance between two 3D coordinates.
func Distance(p1, p2 Coordinate) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2) + math.Pow(p1.Z-p2.Z, 2))
}

// GetClosestApproachDetails calculates the time and minimum Distance at which two planes will be closest during their respective flights.
func (f1 Flight) GetClosestApproachDetails(f2 Flight) (closestTime time.Time, distanceBetweenPlanesatCA float64) {
	flight1ClosestCoord, flight2ClosestCoord := FindClosestApproachDuringTransit(f1.FlightSchedule, f2.FlightSchedule)

	flight1Distance := Distance(f1.FlightSchedule.Depature, f1.FlightSchedule.Destination)
	distBtwDepatureAndClosestApproachForFlight1 := Distance(f1.FlightSchedule.Depature, flight1ClosestCoord)

	f1fractionofCA := distBtwDepatureAndClosestApproachForFlight1 / flight1Distance

	totalFlightDuration1 := f1.ExpectedLandingTime.Sub(f1.TakeoffTime)
	closestTime = f1.TakeoffTime.Add(time.Duration(float64(totalFlightDuration1) * f1fractionofCA))

	distanceBetweenPlanesatCA = Distance(flight1ClosestCoord, flight2ClosestCoord)

	return closestTime, distanceBetweenPlanesatCA
}
