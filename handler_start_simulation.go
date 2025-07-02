package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/josephus-git/TCAS-simulation/internal/aviation"
)

// Simulation parameters (can be passed to Start or kept as package constants)
const AirportLaunchIntervalMin = 5 * time.Second     // Min random delay before an airport tries to launch a plane
const AirportLaunchIntervalMax = 10 * time.Second    // Max random delay before an airport tries to launch a plane
const FlightMonitorInterval = 500 * time.Millisecond // How often the monitor checks planes for landing time

// Your existing startSimulation function (assumed to be in main.go or accessible)
// This function launches goroutines for each airport to handle takeoffs.
func startSimulation(simState *aviation.SimulationState, ctx context.Context, wg *sync.WaitGroup) error {
	log.Printf("--- Starting Airport Launch Operations ---")
	for i := range simState.Airports {
		ap := &simState.Airports[i] // Get a pointer to the airport
		wg.Add(1)                   // Add to WaitGroup for each airport goroutine
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

				airport.mu.Lock() // Lock airport to safely check and pick a plane
				if len(airport.Planes) > 0 {
					planeToTakeOff := airport.Planes[0] // Pick the first available plane for simplicity
					airport.mu.Unlock()                 // Unlock airport before calling TakeOff

					// IMPORTANT: Pass the global simState here.
					_, err := airport.TakeOff(planeToTakeOff, simState) // Pass the simState from main
					if err != nil {
						// log.Printf("error taking off from %s: %v", airport.Serial, err)
					}
				} else {
					airport.mu.Unlock() // Always ensure lock is released
					// log.Printf("Airport %s has no planes to take off.", airport.Serial)
				}
			}
		}(ap) // Pass airport pointer
	}
	return nil
}

// Start initiates and runs the flight simulation for a given duration in minutes.
func startSimulationInit(simState *aviation.SimulationState, durationMinutes int) error {
	log.Printf("--- TCAS Simulation Started for %d minutes ---", durationMinutes)

	// WaitGroup to keep track of running goroutines
	var wg sync.WaitGroup

	// Context for graceful shutdown of goroutines after SimulationDuration
	simulationDuration := time.Duration(durationMinutes) * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), simulationDuration)
	defer cancel() // Ensure cancel is called when Start() exits

	// 3. Start the takeoff simulation (using your provided startSimulation function)
	// Pass ctx and wg to startSimulation so airport goroutines can respect shutdown
	err := startSimulation(&simState, ctx, &wg)
	if err != nil {
		log.Fatalf("Failed to start airport simulations: %v", err)
	}

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
			// other locks (like airport.mu) while globalSimState.mu is held.
			globalSimState.mu.Lock()
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
			globalSimState.mu.Unlock() // Release lock on global state after identifying planes

			// Process the planes that are ready to land
			for _, p := range planesToLand {
				// Find the corresponding destination airport object
				currentFlight := p.FlightLog[len(p.FlightLog)-1]
				var destinationAirport *aviation.Airport = nil
				for i := range globalSimState.Airports {
					ap := &globalSimState.Airports[i]
					// Match airport by location, using Epsilon for robust float comparison
					if aviation.distance(ap.Location, currentFlight.FlightSchedule.Arrival) < aviation.Epsilon {
						destinationAirport = ap
						break
					}
				}

				if destinationAirport != nil {
					// Call the Land function. It handles its own internal locking for runway use
					// and updates globalSimState.PlanesInFlight by removing the landed plane.
					// The Land function itself acquires the necessary simState.mu lock for its modification.
					err := destinationAirport.Land(p, globalSimState)
					if err != nil {
						// This error could be due to runway busy. The plane remains in PlanesInFlight
						// and will be retried in the next monitor interval.
						// log.Printf("Plane %s landing attempt error at Airport %s: %v", p.Serial, destinationAirport.Serial, err)
					}
				} else {
					log.Printf("Monitor Error: Destination airport not found for plane %s (arrival coord: %s)", p.Serial, currentFlight.FlightSchedule.Arrival.String())
				}
			}
		}
	}(&simState, ctx)

	// Wait for all goroutines (airport launchers and flight monitor) to finish.
	// This will happen when ctx.Done() is closed after simulationDuration.
	wg.Wait()

	log.Printf("\n--- All simulation goroutines have stopped. ---")
	log.Printf("Final Simulation State Summary:")
	simState.mu.Lock() // Acquire lock to safely read final count of planes in flight
	log.Printf("  Planes currently in flight: %d", len(simState.PlanesInFlight))
	simState.mu.Unlock()

	for i := range simState.Airports {
		ap := &simState.Airports[i]
		ap.mu.Lock() // Acquire lock for each airport to safely read its parked planes count
		log.Printf("  Airport %s has %d planes parked.", ap.Serial, len(ap.Planes))
		ap.mu.Unlock()
	}
	log.Printf("--- TCAS Simulation Ended ---")
}
