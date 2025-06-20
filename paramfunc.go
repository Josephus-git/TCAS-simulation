package main

import (
	"fmt"
	"math/rand"
)

// create appropriate amount of airports and airplanes
func initializeAirports(conf *apiConfig) {
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

	fmt.Printf("planes created: %d\n", conf.noOfAirplanes)
}

func createAirport(airportCount, planecount, totalNumPlanes int) airport {
	return airport{
		serial:        generateSerialNumber(airportCount, "ap"),
		location:      coord{0, 0, 0},
		planeCapacity: generatePlaneCapacity(totalNumPlanes, planecount),
		runway:        generateRunway(),
	}
}

func createPlane(planeCount int) plane {
	return plane{
		serial:        generateSerialNumber(planeCount, "p"),
		planeInFlight: false,
		cruiseSpeed:   0.1,
		flightLog:     []flight{},
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

func generateRunway() runway {
	return runway{
		numberOfRunway:  1,
		noOfRunwayinUse: 0,
	}
}

func generateSerialNumber(count int, paramType string) string {
	var serialNumber string
	adjustedCount := count - 1
	blockIndex := adjustedCount / 999

	letter := string('A' + rune(blockIndex))

	numericalPart := (adjustedCount % 999) + 1
	formatedNumericPart := fmt.Sprintf("%03d", numericalPart)

	if paramType == "p" {
		serialNumber = fmt.Sprintf("P_%s%s", letter, formatedNumericPart)
	} else if paramType == "ap" {
		serialNumber = fmt.Sprintf("AP_%s%s", letter, formatedNumericPart)
	} else if paramType == "f" {
		serialNumber = fmt.Sprintf("F_%s%s", letter, formatedNumericPart)
	}

	return serialNumber
}
