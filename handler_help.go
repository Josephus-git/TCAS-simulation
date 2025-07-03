package main

import (
	"fmt"

	"github.com/josephus-git/TCAS-simulation/internal/aviation"
)

// helpFunc displays a welcome message and lists all available commands with their descriptions.
func helpFunc(simState *aviation.SimulationState, argument2 string) {
	fmt.Print("Welcome to TCAS-simulator!\nUsage\n\n")
	for key := range getCommand(simState, argument2) {
		fmt.Printf("%s: %s\n", getCommand(simState, argument2)[key].name, getCommand(simState, argument2)[key].description)
	}
}
