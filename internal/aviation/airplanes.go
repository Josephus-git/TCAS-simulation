package aviation

import (
	"math/rand"
)

type Plane struct {
	serial        string
	planeInFlight bool
	cruiseSpeed   float64
	flightLog     []Flight
}

func createPlane(planeCount int) Plane {
	return Plane{
		serial:        generateSerialNumber(planeCount, "p"),
		planeInFlight: false,
		cruiseSpeed:   0.1,
		flightLog:     []Flight{},
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
