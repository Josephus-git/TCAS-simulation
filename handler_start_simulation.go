package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/josephus-git/TCAS-simulation/internal/aviation"
)

// startInit parses the duration string and initializes the simulation.
// Handles input validation, ensuring a positive integer for simulation duration.
func startInit(simState *aviation.SimulationState, durationMinutesString string) {
	durationMinutes, err := strconv.Atoi(durationMinutesString)
	if err != nil {
		fmt.Println("usage: start <integer> (integer represents time in minutes)")
		return
	}
	if durationMinutes < 1 {
		fmt.Println("Please input a valid integer greater than 0")
		return
	}
	startSimulationInit(simState, time.Duration(durationMinutes))
}

// Simulation parameters

// AirportLaunchIntervalMin is the min random delay before an airport tries to launch a plane
const AirportLaunchIntervalMin = 5 * time.Second

// AirportLaunchIntervalMax is the max random delay before an airport tries to launch a plane
const AirportLaunchIntervalMax = 10 * time.Second

// FlightMonitorInterval is how often the monitor checks planes for landing time
const FlightMonitorInterval = 500 * time.Millisecond

// startSimulationInit initializes and starts the TCAS simulation, managing goroutines for takeoffs and landings.
// It sets up a context for graceful shutdown and waits for all simulation activities to complete.
func startSimulationInit(simState *aviation.SimulationState, durationMinutes time.Duration) {
	log.Printf("--- TCAS Simulation Started for %d minutes ---", durationMinutes)

	// WaitGroup to keep track of running goroutines
	var wg sync.WaitGroup

	// Context for graceful shutdown of goroutines after SimulationDuration
	simulationDuration := time.Duration(durationMinutes) * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), simulationDuration)
	defer cancel() // Ensure cancel is called when Start() exits

	// Start the takeoff simulation (using your provided startSimulation function)
	// Pass ctx and wg to startSimulation so airport goroutines can respect shutdown
	startSimulation(simState, ctx, &wg)

	// --- Start Flight Monitoring Goroutine (for landings) ---
	log.Printf("\n--- Starting Flight Landing Monitor ---")
	wg.Add(1) // Add for the monitor goroutine
	go func(globalSimState *aviation.SimulationState, ctx context.Context) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done(): // Check if the main simulation context is done
				log.Printf("Flight monitor stopping.")
				return // Exit goroutine
			default:
				// Continue monitoring
			}

			time.Sleep(FlightMonitorInterval) // Sleep to avoid busy-waiting and reduce CPU usage

			// We need to safely access and potentially modify globalSimState.PlanesInFlight.
			// It's safer to copy the list of planes to be processed, then release the lock,
			// and then process the copy. This prevents deadlocks if Land() tries to acquire
			// other locks (like airport.Mu) while globalSimState.Mu is held.
			globalSimState.Mu.Lock()
			planesToLand := []aviation.Plane{}
			currentTime := time.Now()

			for _, p := range globalSimState.PlanesInFlight {
				if len(p.FlightLog) > 0 {
					currentFlight := p.FlightLog[len(p.FlightLog)-1]
					// Check if current time is past or at the plane's scheduled landing time
					if currentTime.After(currentFlight.LandingTime) || currentTime.Equal(currentFlight.LandingTime) {
						planesToLand = append(planesToLand, p)
					}
				}
			}
			globalSimState.Mu.Unlock() // Release lock on global state after identifying planes

			// Process the planes that are ready to land
			for _, p := range planesToLand {
				// Find the corresponding destination airport object
				currentFlight := p.FlightLog[len(p.FlightLog)-1]
				var destinationAirport *aviation.Airport = nil
				for i := range globalSimState.Airports {
					ap := globalSimState.Airports[i]
					// Match airport by location, using Epsilon for robust float comparison
					if aviation.Distance(ap.Location, currentFlight.FlightSchedule.Arrival) < aviation.Epsilon {
						destinationAirport = ap
						break
					}
				}

				if destinationAirport != nil {
					// Call the Land function. It handles its own internal locking for runway use
					// and updates globalSimState.PlanesInFlight by removing the landed plane.
					// The Land function itself acquires the necessary simState.Mu lock for its modification.
					err := destinationAirport.Land(p, globalSimState)
					if err != nil {
						// This error could be due to runway busy. The plane remains in PlanesInFlight
						// and will be retried in the next monitor interval.
					}
				} else {
					log.Printf("Monitor Error: Destination airport not found for plane %s (arrival coord: %s)", p.Serial, currentFlight.FlightSchedule.Arrival.String())
				}
			}
		}
	}(simState, ctx)

	// Wait for all goroutines (airport launchers and flight monitor) to finish.
	// This will happen when ctx.Done() is closed after simulationDuration.
	wg.Wait()

	log.Printf("\n--- All simulation goroutines have stopped. ---")
	log.Printf("Final Simulation State Summary:")
	simState.Mu.Lock() // Acquire lock to safely read final count of planes in flight
	log.Printf("  Planes currently in flight: %d", len(simState.PlanesInFlight))
	simState.Mu.Unlock()

	for i := range simState.Airports {
		ap := simState.Airports[i]
		ap.Mu.Lock() // Acquire lock for each airport to safely read its parked planes count
		log.Printf("  Airport %s has %d planes parked.", ap.Serial, len(ap.Planes))
		ap.Mu.Unlock()
	}
	log.Printf("--- TCAS Simulation Ended ---")
}

// startSimulation launches goroutines for each airport to handle takeoffs.
func startSimulation(simState *aviation.SimulationState, ctx context.Context, wg *sync.WaitGroup) {
	log.Printf("--- Starting Airport Launch Operations ---")
	for i := range simState.Airports {
		ap := simState.Airports[i] // Get a pointer to the airport
		wg.Add(1)                  // Add to WaitGroup for each airport goroutine
		go func(airport *aviation.Airport) {
			defer wg.Done()
			airportRand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)*1000)) // Unique seed for each airport

			for {
				select {
				case <-ctx.Done(): // Check if the main simulation context is done
					log.Printf("Airport %s stopping launch operations.", airport.Serial)
					return // Exit goroutine
				default:
					// Continue operation
				}

				time.Sleep(time.Duration(airportRand.Intn(int(AirportLaunchIntervalMax.Seconds()-AirportLaunchIntervalMin.Seconds())+1)+int(AirportLaunchIntervalMin.Seconds())) * time.Second) // Wait 5-10 seconds

				airport.Mu.Lock() // Lock airport to safely check and pick a plane
				if len(airport.Planes) > 0 {
					planeToTakeOff := airport.Planes[0] // Pick the first available plane for simplicity
					airport.Mu.Unlock()                 // Unlock airport before calling TakeOff

					// IMPORTANT: Pass the global simState here.
					_, err := airport.TakeOff(planeToTakeOff, simState) // Pass the simState from main
					if err != nil {
						// log.Printf("error taking off from %s: %v", airport.Serial, err)
					}
				} else {
					airport.Mu.Unlock() // Always ensure lock is released
					// log.Printf("Airport %s has no planes to take off.", airport.Serial)
				}
			}
		}(ap) // Pass airport pointer
	}
}
