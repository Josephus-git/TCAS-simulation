package aviation

import (
	"fmt"

	"github.com/josephus-git/TCAS-simulation/internal/config"
)

// Airport represents an Airport with its location
type Airport struct {
	serial        string
	location      Coordinate
	planeCapacity int
	runway        runway
	planes        []Plane
}

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
		for range newAirport.planeCapacity {
			newPlane := createPlane(planesGenerated)
			newAirport.planes = append(newAirport.planes, newPlane)
			planesGenerated += 1
		}
		planesCreated += newAirport.planeCapacity
		simState.Airports = append(simState.Airports, newAirport)
		airportsCreated = i + 1
	}

	listOfAirportCoordinates := generateCoordinates(len(simState.Airports))

	for i := range simState.Airports {
		newLocation := Coordinate{listOfAirportCoordinates[i].X, listOfAirportCoordinates[i].Y, 0.0}
		simState.Airports[i].location = newLocation
	}

	fmt.Printf("planes created: %d\n", conf.NoOfAirplanes)
}

func createAirport(airportCount, planecount, totalNumPlanes int) Airport {
	return Airport{
		serial:        generateSerialNumber(airportCount, "ap"),
		planeCapacity: generatePlaneCapacity(totalNumPlanes, planecount),
		runway:        generateRunway(),
	}
}

// generateRunway creates and returns a new runway configuration.
func generateRunway() runway {
	return runway{
		numberOfRunway:  1,
		noOfRunwayinUse: 0,
	}
}
