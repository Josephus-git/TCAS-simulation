package main

import (
	"log"

	"github.com/josephus-git/TCAS-simulation/internal/aviation"
)

func emergencyStop(simState *aviation.SimulationState) {
	if simulationCancelFunc != nil {
		log.Println("\n--- EMERGENCY STOP ACTIVATED! Signaling all goroutines to stop... ---")
		simulationCancelFunc() // Trigger cancellation
		// Reset the cancel func to indicate no active simulation,
		// and prevent multiple calls to a potentially nil context if Start() finished.
		simulationCancelFunc = nil
		<-simState.SimStatusChannel

	} else {
		log.Println("EmergencyStop: Simulation not running or cancel function not initialized.")
	}
}
