package aviation

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Plane struct {
	Serial        string
	PlaneInFlight bool
	CruiseSpeed   float64
	FlightLog     []Flight
}

func createPlane(planeCount int) Plane {
	return Plane{
		Serial:        generateSerialNumber(planeCount, "p"),
		PlaneInFlight: false,
		CruiseSpeed:   0.1,
		FlightLog:     []Flight{},
	}
}

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

// Get planes position
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

// resultant distance obtained by getting the magnitude of distance btw the two coordinates
func distance(p1, p2 Coordinate) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2) + math.Pow(p1.Z-p2.Z, 2))
}

// get the time at which the planes will be closest and the distance at this time
func (f1 Flight) GetClosestApproachDetails(f2 Flight) (closestTime time.Time, distanceBetweenPlanesatCA float64) {
	// get time When planes will get to coincidence point
	flight1ClosestCoord, flight2ClosestCoord := FindClosestApprachDuringTransit(f1.FlightSchedule, f2.FlightSchedule)

	flight1Distance := distance(f1.FlightSchedule.Depature, f1.FlightSchedule.Arrival)
	distBtwDepatureAndClosestApproachForFlight1 := distance(f1.FlightSchedule.Depature, flight1ClosestCoord)

	f1fractionofCA := distBtwDepatureAndClosestApproachForFlight1 / flight1Distance

	totalFlightDuration1 := f1.LandingTime.Sub(f1.TakeoffTime)
	closestTime = f1.TakeoffTime.Add(time.Duration(float64(totalFlightDuration1) * f1fractionofCA))

	distanceBetweenPlanesatCA = distance(flight1ClosestCoord, flight2ClosestCoord)

	return closestTime, distanceBetweenPlanesatCA
}
