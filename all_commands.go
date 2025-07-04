package main

import (
	"github.com/josephus-git/TCAS-simulation/internal/aviation"
)

// cliCommand defines the structure for a command-line interface command.
type cliCommand struct {
	name        string
	description string
	callback    func()
}

// getCommand returns a map of available CLI commands for the TCAS-simulator.
func getCommand(simState *aviation.SimulationState, argument2 string) map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the TCAS-simulator",
			callback: func() {
				commandExit()
			},
		},
		"help": {
			name:        "help",
			description: "Display usage of the application",
			callback: func() {
				helpFunc(simState, argument2)
			},
		},
		"start": {
			name:        "start",
			description: "Initializes and starts the application",
			callback: func() {
				go startInit(simState, argument2)
			},
		},
		"get": {
			name:        "get",
			description: "prints details of the simulation such as airports, Planes and flights to the console",
			callback: func() {
				getDetails(simState, argument2)
			},
		},
		"log": {
			name:        "log",
			description: "logs details of the simulation such as airports, Planes and flights to an appropriate file",
			callback: func() {
				logDetails(simState, argument2)
			},
		},
		"q": {
			name:        "emergency stop",
			description: "Immediately halts the active simulation.",
			callback: func() {
				emergencyStop(simState)
			},
		},
	}
	return commands
}
