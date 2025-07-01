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
	planes        []plane
}

type runway struct {
	numberOfRunway  int
	noOfRunwayinUse int
}

// InitializeAirports creates appropriate amount of airports and airplanes
func InitializeAirports(conf *config.Config) {

	planesCreated := 0
	airportsCreated := 0

	for i := 0; planesCreated < conf.noOfAirplanes; i++ {
		newAirport := createAirport(airportsCreated, planesCreated, conf.noOfAirplanes)
		planesGenerated := planesCreated
		for range newAirport.planeCapacity {
			newPlane := createPlane(planesGenerated)
			newAirport.planes = append(newAirport.planes, newPlane)
			planesGenerated += 1
		}
		planesCreated += newAirport.planeCapacity
		conf.listAirports = append(conf.listAirports, newAirport)
		airportsCreated = i + 1
	}

	listOfAirportCoordinates := generateCoordinates(len(conf.listAirports))

	for i := range conf.listAirports {
		newLocation := Coordinate{listOfAirportCoordinates[i].X, listOfAirportCoordinates[i].Y, 0.0}
		conf.listAirports[i].location = newLocation
	}

	fmt.Printf("planes created: %d\n", conf.noOfAirplanes)
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
