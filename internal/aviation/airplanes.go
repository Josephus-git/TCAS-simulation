package aviation

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/josephus-git/DEV_ACTUAL_1/TCAS-simulation/internal/config"
)

type plane struct {
	serial        string
	planeInFlight bool
	cruiseSpeed   float64
	flightLog     []Flight
}

func createPlane(planeCount int) plane {
	return plane{
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

func getNumberPlanes(conf *config.Config) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to TCAS-simulator")
	notValidInput := true

	for i := 0; notValidInput; i++ {

		fmt.Print("Input the number of planes for the simulation > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		num, err := strconv.Atoi(input[0])
		if err != nil {
			fmt.Println("Please input a valid integer")
			continue
		}
		if num < 2 {
			fmt.Println("Please input a valid integer greater than 1")
			continue
		}

		conf.noOfAirplanes = num
		notValidInput = false
	}

}
