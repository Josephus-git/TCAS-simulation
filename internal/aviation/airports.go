package aviation

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/josephus-git/TCAS-simulation/internal/config"
	"github.com/josephus-git/TCAS-simulation/internal/util"
)

// Airport represents an Airport with its location
type Airport struct {
	Serial        string
	Location      Coordinate
	PlaneCapacity int
	Runway        runway
	Planes        []Plane
	Mu            sync.Mutex
}

// runway represents the state of an airport's runways.
type runway struct {
	numberOfRunway  int
	noOfRunwayinUse int
}

// InitializeAirports creates appropriate amount of airports and airplanes
func InitializeAirports(conf *config.Config, simState *SimulationState) {

	planesCreated := 0
	airportsCreated := 0

	for i := 0; planesCreated < conf.NoOfAirplanes; i++ {
		newAirport := createAirport(airportsCreated, planesCreated, conf.NoOfAirplanes)
		planesGenerated := planesCreated
		for range newAirport.PlaneCapacity {
			newPlane := createPlane(planesGenerated)
			newAirport.Planes = append(newAirport.Planes, newPlane)
			planesGenerated += 1
		}
		simState.Airports = append(simState.Airports, &newAirport)
		planesCreated += newAirport.PlaneCapacity
		airportsCreated = i + 1
	}

	listOfAirportCoordinates := generateCoordinates(len(simState.Airports))

	for i := range simState.Airports {
		newLocation := Coordinate{listOfAirportCoordinates[i].X, listOfAirportCoordinates[i].Y, 0.0}
		simState.Airports[i].Location = newLocation
	}

	fmt.Printf("Initialized: %d airports, %d planes distributed among airports.\n",
		len(simState.Airports), conf.NoOfAirplanes)
}

// createAirport initializes and returns a new Airport struct.
// It generates a serial number, plane capacity, and runway details for the airport.
func createAirport(airportCount, planecount, totalNumPlanes int) Airport {
	return Airport{
		Serial:        util.GenerateSerialNumber(airportCount, "ap"),
		PlaneCapacity: generatePlaneCapacity(totalNumPlanes, planecount),
		Runway:        generateRunway(),
	}
}

// generateRunway creates and returns a new runway configuration.
func generateRunway() runway {
	randomNumber := rand.Intn(3) + 1
	return runway{
		numberOfRunway:  randomNumber,
		noOfRunwayinUse: 0,
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
