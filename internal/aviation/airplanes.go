package aviation

import (
	"fmt"
	"math"
	"math/rand"
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
		CruiseSpeed:   0.1,
		FlightLog:     []Flight{},
	}
}

// generatePlaneCapacity calculates a random number of planes to create,
// adjusting the quantity based on the total target and already generated planes.
func generatePlaneCapacity(totalPlanes, planeGenerated int) int {
	var randomNumber int
	if totalPlanes < 20 {
		planeToCreate := totalPlanes - planeGenerated
		if planeToCreate <= 3 {
			randomNumber = planeToCreate
		} else {
			randomNumber = rand.Intn(2) + 1
		}

	} else if totalPlanes < 100 {
		planeToCreate := totalPlanes - planeGenerated
		if planeToCreate <= 6 {
			randomNumber = planeToCreate
		} else {
			randomNumber = rand.Intn(5) + 1
		}

	} else {
		planeToCreate := totalPlanes - planeGenerated
		if planeToCreate <= 30 {
			randomNumber = planeToCreate
		} else {
			randomNumber = rand.Intn(20) + 10
		}

	}
	return randomNumber
}

// getPlanePosition calculates the plane's interpolated coordinates (X, Y, Z) at a given time during a specific flight.
// It returns an error if the provided time is outside the flight's duration.
func (p Plane) getPlanePosition(f Flight, t time.Time) (Coordinate, error) {
	if t.Before(f.TakeoffTime) || t.After(f.LandingTime) { //*** return here to check incase plane shold go on another Flight
		return Coordinate{}, fmt.Errorf("time %v is outside Flight %s duration", t, f.FlightID)
	}

	// Calculate fraction of Flight completed (normalized time 0-1)
	totalDuration := f.LandingTime.Sub(f.TakeoffTime)
	elapsed := t.Sub(f.TakeoffTime)
	progress := float64(elapsed) / float64(totalDuration)

	// get the arival and departure location
	departureLocation := f.FlightSchedule.Depature
	arivalLocation := f.FlightSchedule.Arrival

	// Calculate the intermediate point
	pX := departureLocation.X + (departureLocation.X-arivalLocation.X)*progress
	pY := departureLocation.Y + (departureLocation.Y-arivalLocation.Y)*progress
	pZ := f.CruisingAltitude

	return Coordinate{X: pX, Y: pY, Z: pZ}, nil
}

// distance calculates the Euclidean distance between two 3D coordinates.
func distance(p1, p2 Coordinate) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2) + math.Pow(p1.Z-p2.Z, 2))
}

// GetClosestApproachDetails calculates the time and minimum distance at which two planes will be closest during their respective flights.
func (f1 Flight) GetClosestApproachDetails(f2 Flight) (closestTime time.Time, distanceBetweenPlanesatCA float64) {
	flight1ClosestCoord, flight2ClosestCoord := FindClosestApproachDuringTransit(f1.FlightSchedule, f2.FlightSchedule)

	flight1Distance := distance(f1.FlightSchedule.Depature, f1.FlightSchedule.Arrival)
	distBtwDepatureAndClosestApproachForFlight1 := distance(f1.FlightSchedule.Depature, flight1ClosestCoord)

	f1fractionofCA := distBtwDepatureAndClosestApproachForFlight1 / flight1Distance

	totalFlightDuration1 := f1.LandingTime.Sub(f1.TakeoffTime)
	closestTime = f1.TakeoffTime.Add(time.Duration(float64(totalFlightDuration1) * f1fractionofCA))

	distanceBetweenPlanesatCA = distance(flight1ClosestCoord, flight2ClosestCoord)

	return closestTime, distanceBetweenPlanesatCA
}
