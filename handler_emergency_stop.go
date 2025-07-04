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
		simState.SimStatus = false
		<-simState.SimStatusChannel
		if stopTrigger.Stop() {
		} else {
			log.Printf("Simulation has ended a while ago")
		}
		simulationCancelFunc = nil

	} else {
		log.Println("EmergencyStop: Simulation not running")
	}
}
