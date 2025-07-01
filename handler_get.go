package main

import (
	"fmt"

	"github.com/josephus-git/TCAS-simulation/internal/aviation"
)

// getDetails displays specific simulation details (airports, airplanes, or flights) based on the provided argument.
// It prints usage instructions if an invalid option is given.
func getDetails(simState *aviation.SimulationState, argument2 string) error {
	switch argument2 {
	case "airports":
		getAirportDetails(simState)
	case "airplanes":
		getAirPlanesDetails(simState)
	case "flights":
		fmt.Println("successfully got flights")
	default:
		fmt.Println("usage: get <option>, options: airports, airplanes, flights")
	}
	return nil
}

// getAirPlanesDetails prints selected details of all airplanes from the simulation state to the console.
func getAirPlanesDetails(simstate *aviation.SimulationState) {
	Planes := []aviation.Plane{}

	for _, ap := range simstate.Airports {
		Planes = append(Planes, ap.Planes...)
	}
	fmt.Println("\n--- Printing selected fields for each plane ---")
	for i, p := range Planes {
		fmt.Printf("Plane %d (Serial: %s):\n", i+1, p.Serial)
		fmt.Printf("  In Flight: %t\n", p.PlaneInFlight)
		fmt.Printf("  Cruise Speed: %.2f km/h\n", p.CruiseSpeed)
		fmt.Println("  Flight Log:")
		if len(p.FlightLog) == 0 {
			fmt.Println("    No flights recorded for this plane.")
		} else {
			for j := range p.FlightLog { // Looping to count flights, but not printing content if 'flight' is empty
				fmt.Printf("    Flight %d (details depend on 'flight' struct's fields)\n", j+1)
			}
		}
		fmt.Println("-------------------------------------------")
	}
}

// getAirportDetails prints selected details of all airports from the simulation state to the console.
func getAirportDetails(simState *aviation.SimulationState) {
	fmt.Println("\n--- Printing selected fields for each airport ---")
	for i, ap := range simState.Airports {
		fmt.Printf("Airport %d (Serial: %s):\n", i+1, ap.Serial)
		fmt.Printf("  Location: %v\n", ap.Location)
		fmt.Printf("  Plane Capacity: %d\n", ap.PlaneCapacity)
		fmt.Printf("  Runway: %v\n", ap.Runway)
		fmt.Println("  Planes:")
		if len(ap.Planes) == 0 {
			fmt.Println("    No Planes currently.")
		} else {
			for j, p := range ap.Planes {
				fmt.Printf("    %d. Serial: %s\n", j+1, p.Serial)
			}
		}
		fmt.Println("-------------------------------------------")
	}
}
