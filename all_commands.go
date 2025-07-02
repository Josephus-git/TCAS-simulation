package main

import (
	"github.com/josephus-git/TCAS-simulation/internal/aviation"
)

// cliCommand defines the structure for a command-line interface command.
type cliCommand struct {
	name        string
	description string
	callback    func() error
}

// getCommand returns a map of available CLI commands for the TCAS-simulator.
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
				return startSimulationInit(simState, durationMinutes)
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
