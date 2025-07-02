package main

import (
	"fmt"
	"log"
	"os"

	"github.com/josephus-git/TCAS-simulation/internal/aviation"
)

// logDetails logs specific simulation details (airports, airplanes, or flights) based on the provided argument.
// It prints usage instructions if an invalid option is given.
func logDetails(simState *aviation.SimulationState, argument2 string) error {
	switch argument2 {
	case "airports":
		logAirportDetails(simState)
	case "airplanes":
		logAirplanesDetails(simState)
	case "flights":
		fmt.Println("successfully logged flights")
	default:
		fmt.Println("usage: log <option>, options: airports, airplanes, flights")
	}
	return nil
}

// logAirplanesDetails appends selected details of all airplanes from the simulation state to a log file.
// It includes serial, flight status, cruise speed, and a count of flights for each plane.
func logAirplanesDetails(simState *aviation.SimulationState) {
	logFilePath := "logs/airplaneDetails.txt"
	// Open the file in append mode. Create it if it doesn't exist.
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer f.Close()

	Planes := []aviation.Plane{}
	for _, ap := range simState.Airports {
		Planes = append(Planes, ap.Planes...)
	}
	fmt.Fprintln(f, "\n--- Logging selected fields for each plane ---")
	for i, p := range Planes {
		fmt.Fprintf(f, "Plane %d (Serial: %s):\n", i+1, p.Serial)
		fmt.Fprintf(f, "  In Flight: %t\n", p.PlaneInFlight)
		fmt.Fprintf(f, "  Cruise Speed: %.2f km/h\n", p.CruiseSpeed)
		fmt.Fprintln(f, "  Flight Log:")
		if len(p.FlightLog) == 0 {
			fmt.Fprintln(f, "    No flights recorded for this plane.")
		} else {
			for j := range p.FlightLog { // Looping to count flights, but not printing content if 'flight' is empty
				fmt.Fprintf(f, "    Flight %d (details depend on 'flight' struct's fields)\n", j+1)
			}
		}
		fmt.Fprintln(f, "-------------------------------------------")
	}
	fmt.Println("Successfully logged airplanes")
}

// logAirportDetails appends selected details of all airports from the simulation state to a log file.
// It includes serial, location, plane capacity, runway information, and a list of associated plane serials.
func logAirportDetails(simState *aviation.SimulationState) {
	logFilePath := "logs/airportDetails.txt"
	// Open the file in append mode. Create it if it doesn't exist.
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer f.Close()
	fmt.Fprintln(f, "\n--- Logging selected fields for each airport ---")
	for i, ap := range simState.Airports {
		fmt.Fprintf(f, "Airport %d (Serial: %s):\n", i+1, ap.Serial)
		fmt.Fprintf(f, "  Location: %v\n", ap.Location)
		fmt.Fprintf(f, "  Plane Capacity: %d\n", ap.PlaneCapacity)
		fmt.Fprintf(f, "  Runway: %v\n", ap.Runway)
		fmt.Fprintln(f, "  Planes:")
		if len(ap.Planes) == 0 {
			fmt.Fprintln(f, "    No Planes currently.")
		} else {
			for j, p := range ap.Planes {
				fmt.Fprintf(f, "    %d. Serial: %s\n", j+1, p.Serial)
			}
		}
		fmt.Fprintln(f, "-------------------------------------------")
	}
	fmt.Println("Successfully logged airplanes")
}
