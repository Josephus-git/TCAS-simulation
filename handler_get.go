package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/josephus-git/TCAS-simulation/internal/aviation"
)

// getDetails displays specific simulation details (airports, airplanes, or flights) based on the provided argument.
// It prints usage instructions if an invalid option is given.
func getDetails(simState *aviation.SimulationState, argument2 string) {
	switch argument2 {
	case "airports":
		getAirportDetails(simState)
	case "airplanes":
		getAirPlanesDetails(simState)
	case "flights":
		getFlightDetails(simState)
	case "all":
		getAirportDetails(simState)
		getAirPlanesDetails(simState)
		getFlightDetails(simState)
	default:
		fmt.Println("usage: get <option>, options: airports, airplanes, flights, all")
	}
}

// getAirPlanesDetails prints selected details of all flights logged in all various planes
func getFlightDetails(simState *aviation.SimulationState) {
	var simTime time.Time
	if simState.SimStatus {
		simTime = time.Now()
	} else {
		simTime = simState.SimEndedTime
	}
	var flightLogs []aviation.Flight

	fmt.Println("\n--- Printing all recorded flights ---")

	for _, airport := range simState.Airports {
		for _, plane := range airport.Planes {
			if len(plane.FlightLog) == 0 {
				continue
			}
			flightLogs = append(flightLogs, plane.FlightLog...)
		}
	}

	for _, plane := range simState.PlanesInFlight {
		flightLogs = append(flightLogs, plane.FlightLog...)
	}

	if len(flightLogs) == 0 {
		fmt.Println("\n--- No flight recorded currently ---")
		return
	}
	sort.Slice(flightLogs, func(i, j int) bool {
		return flightLogs[i].FlightID < flightLogs[j].FlightID
	})
	for i, flight := range flightLogs {
		fmt.Printf("\nflight %d:\n", i)
		printFlightDetails(flight, simTime)
	}
	fmt.Println()
}

// getAirPlanesDetails prints selected details of all airplanes from the simulation state to the console.
func getAirPlanesDetails(simState *aviation.SimulationState) {
	var simTime time.Time
	if simState.SimStatus {
		simTime = time.Now()
	} else {
		simTime = simState.SimEndedTime
	}
	Planes := []aviation.Plane{}

	for _, airport := range simState.Airports {
		Planes = append(Planes, airport.Planes...)
	}

	Planes = append(Planes, simState.PlanesInFlight...)
	sort.Slice(Planes, func(i, j int) bool {
		return Planes[i].Serial < Planes[j].Serial
	})

	fmt.Println("\n--- Printing selected fields for each plane in airports ---")
	for i, plane := range Planes {
		fmt.Printf("Plane %d (Serial: %s):\n", i+1, plane.Serial)
		fmt.Printf("  In Flight: %t\n", plane.PlaneInFlight)
		fmt.Printf("  Cruise Speed: %.2f m/s\n", plane.CruiseSpeed)
		fmt.Println("  Flight Log:")
		if len(plane.FlightLog) == 0 {
			fmt.Println("    No flights recorded for this plane.")
		} else {
			for _, flight := range plane.FlightLog { // Looping to count flights, but not printing content if 'flight' is empty
				printFlightDetails(flight, simTime)
			}
		}
		fmt.Println("-------------------------------------------")
	}
	fmt.Println()
}

// getAirportDetails prints selected details of all airports from the simulation state to the console.
func getAirportDetails(simState *aviation.SimulationState) {
	fmt.Println("\n--- Printing selected fields for all airports ---")
	for i, airport := range simState.Airports {
		fmt.Printf("Airport %d (Serial: %s):\n", i+1, airport.Serial)
		fmt.Printf("  Location: %v\n", airport.Location)
		fmt.Printf("  Runway: %v\n", airport.Runway)
		fmt.Println("  Planes:")
		if len(airport.Planes) == 0 {
			fmt.Println("    No Planes currently.")
		} else {
			for j, plane := range airport.Planes {
				fmt.Printf("    %d. Serial: %s\n", j+1, plane.Serial)
			}
		}
		fmt.Println("-------------------------------------------")
	}
	fmt.Println()
}

// getFlightDetails prints all details for a given Flight struct,
func printFlightDetails(flight aviation.Flight, simTime time.Time) {
	fmt.Println("    --- Flight Details ---")
	fmt.Printf("    Flight ID: %s\n", flight.FlightID)
	fmt.Printf("    Takeoff Time: %s\n", flight.TakeoffTime.Format("15:04:05"))
	fmt.Printf("    Expected Landing Time: %s\n", flight.ExpectedLandingTime.Format("15:04:05"))
	fmt.Printf("    Cruising Altitude: %.2f meters\n", flight.CruisingAltitude)
	fmt.Printf("    Depature Airport: %s\n", flight.DepatureAirPort)
	fmt.Printf("    Destination Airport: %s\n", flight.ArrivalAirPort)
	var actualLandingTime string
	if flight.ActualLandingTime.IsZero() {
		actualLandingTime = "Plane is yet to land"
	} else {
		actualLandingTime = flight.ActualLandingTime.Format("15:04:05")
	}
	fmt.Printf("    Actual Landing Time: %s\n", actualLandingTime)

	// calculate progress
	progress := flight.GetFlightProgress(simTime)

	fmt.Printf("    Progress: %s\n", progress)
	fmt.Println("    ---------------------------------------")
}
