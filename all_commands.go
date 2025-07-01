package main

import (
	"fmt"
	"log"
	"os"

	"internal/config"
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

func logDetails(conf *config.Config, argument2 string) error {
	switch argument2 {
	case "airports":
		logAirportDetails(conf)
	case "airplanes":
		logAirplanesDetails(conf)
	case "flights":
		fmt.Println("successfully logged flights")
	default:
		fmt.Println("usage: log <option>, options: airports, airplanes, flights")
	}
	return nil
}
func logAirplanesDetails(conf *config.Config) {
	logFilePath := "logs/airPlaneDetails.txt"
	// Open the file in append mode. Create it if it doesn't exist.
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer f.Close()

	planes := []plane{}
	for _, ap := range conf.listAirports {
		planes = append(planes, ap.planes...)
	}
	fmt.Println("\n--- Logging selected fields for each plane ---")
	for i, p := range planes {
		fmt.Fprintf(f, "Plane %d (Serial: %s):\n", i+1, p.serial)
		fmt.Fprintf(f, "  In Flight: %t\n", p.planeInFlight)
		fmt.Fprintf(f, "  Cruise Speed: %.2f km/h\n", p.cruiseSpeed)
		fmt.Fprintln(f, "  Flight Log:")
		if len(p.flightLog) == 0 {
			fmt.Fprintln(f, "    No flights recorded for this plane.")
		} else {
			for j := range p.flightLog { // Looping to count flights, but not printing content if 'flight' is empty
				fmt.Fprintf(f, "    Flight %d (details depend on 'flight' struct's fields)\n", j+1)
			}
		}
		fmt.Fprintln(f, "-------------------------------------------")
	}

}

func logAirportDetails(conf *config.Config) {
	logFilePath := "logs/airPortDetails.txt"
	// Open the file in append mode. Create it if it doesn't exist.
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer f.Close()
	fmt.Fprintln(f, "\n--- Logging selected fields for each airport ---")
	for i, ap := range conf.listAirports {
		fmt.Fprintf(f, "Airport %d (Serial: %s):\n", i+1, ap.serial)
		fmt.Fprintf(f, "  Location: %v\n", ap.location)
		fmt.Fprintf(f, "  Plane Capacity: %d\n", ap.planeCapacity)
		fmt.Fprintf(f, "  Runway: %v\n", ap.runway)
		fmt.Fprintln(f, "  Planes:")
		if len(ap.planes) == 0 {
			fmt.Fprintln(f, "    No planes currently.")
		} else {
			for j, p := range ap.planes {
				fmt.Fprintf(f, "    %d. Serial: %s\n", j+1, p.serial)
			}
		}
		fmt.Fprintln(f, "-------------------------------------------")
	}
}

func getDetails(conf *config.Config, argument2 string) error {
	switch argument2 {
	case "airports":
		getAirportDetails(conf)
	case "airplanes":
		getAirPlanesDetails(conf)
	case "flights":
		fmt.Println("successfully got flights")
	default:
		fmt.Println("usage: get <option>, options: airports, airplanes, flights")
	}
	return nil
}

func getAirPlanesDetails(conf *config.Config) {
	planes := []plane{}

	for _, ap := range conf.listAirports {
		planes = append(planes, ap.planes...)
	}
	fmt.Println("\n--- Printing selected fields for each plane ---")
	for i, p := range planes {
		fmt.Printf("Plane %d (Serial: %s):\n", i+1, p.serial)
		fmt.Printf("  In Flight: %t\n", p.planeInFlight)
		fmt.Printf("  Cruise Speed: %.2f km/h\n", p.cruiseSpeed)
		fmt.Println("  Flight Log:")
		if len(p.flightLog) == 0 {
			fmt.Println("    No flights recorded for this plane.")
		} else {
			for j := range p.flightLog { // Looping to count flights, but not printing content if 'flight' is empty
				fmt.Printf("    Flight %d (details depend on 'flight' struct's fields)\n", j+1)
			}
		}
		fmt.Println("-------------------------------------------")
	}

}

func getAirportDetails(conf *config.Config) {
	fmt.Println("\n--- Printing selected fields for each airport ---")
	for i, ap := range conf.listAirports {
		fmt.Printf("Airport %d (Serial: %s):\n", i+1, ap.serial)
		fmt.Printf("  Location: %v\n", ap.location)
		fmt.Printf("  Plane Capacity: %d\n", ap.planeCapacity)
		fmt.Printf("  Runway: %v\n", ap.runway)
		fmt.Println("  Planes:")
		if len(ap.planes) == 0 {
			fmt.Println("    No planes currently.")
		} else {
			for j, p := range ap.planes {
				fmt.Printf("    %d. Serial: %s\n", j+1, p.serial)
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

func helpFunc(conf *config.Config, argument2 string) error {
	fmt.Print("Welcome to TCAS-simulator!\nUsage\n\n")
	for key := range getCommand(conf, argument2) {
		fmt.Printf("%s: %s\n", getCommand(conf, argument2)[key].name, getCommand(conf, argument2)[key].description)
	}
	return nil
}

func getCommand(conf *config.Config, argument2 string) map[string]cliCommand {
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
				return helpFunc(conf, argument2)
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
			description: "prints details of the simulation such as airports, planes and flights to the console",
			callback: func() error {
				return getDetails(conf, argument2)
			},
		},
		"log": {
			name:        "log",
			description: "logs details of the simulation such as airports, planes and flights to an appropriate file",
			callback: func() error {
				return logDetails(conf, argument2)
			},
		},
	}
	return commands
}
