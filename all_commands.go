package main

import (
	"fmt"

	"log"
	"os"

	"github.com/josephus-git/TCAS-simulation/internal/aviation"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func startSimulation() error {
	fmt.Println("Simulation has started")

	return nil
}

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
func logAirplanesDetails(simState *aviation.SimulationState) {
	logFilePath := "logs/airPlaneDetails.txt"
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
	fmt.Println("\n--- Logging selected fields for each plane ---")
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

}

func logAirportDetails(simState *aviation.SimulationState) {
	logFilePath := "logs/airPortDetails.txt"
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
}

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

func commandExit() error {
	fmt.Println("Closing TCAS-simulator... Goodbye!")
	os.Exit(0)
	return nil
}

func helpFunc(simState *aviation.SimulationState, argument2 string) error {
	fmt.Print("Welcome to TCAS-simulator!\nUsage\n\n")
	for key := range getCommand(simState, argument2) {
		fmt.Printf("%s: %s\n", getCommand(simState, argument2)[key].name, getCommand(simState, argument2)[key].description)
	}
	return nil
}

func getCommand(simState *aviation.SimulationState, argument2 string) map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the TCAS-simulator",
			callback: func() error {
				return commandExit()
			},
		},
		"help": {
			name:        "help",
			description: "Display usage of the application",
			callback: func() error {
				return helpFunc(simState, argument2)
			},
		},
		"start": {
			name:        "start",
			description: "Initializes and starts the application",
			callback: func() error {
				return startSimulation()
			},
		},
		"get": {
			name:        "get",
			description: "prints details of the simulation such as airports, Planes and flights to the console",
			callback: func() error {
				return getDetails(simState, argument2)
			},
		},
		"log": {
			name:        "log",
			description: "logs details of the simulation such as airports, Planes and flights to an appropriate file",
			callback: func() error {
				return logDetails(simState, argument2)
			},
		},
	}
	return commands
}
