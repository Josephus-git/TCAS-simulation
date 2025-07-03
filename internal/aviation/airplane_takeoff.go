package aviation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/josephus-git/TCAS-simulation/internal/util"
)

// CruisingAltitude defines the standard cruising altitude for planes in meters.
const CruisingAltitude = 10000.0

// TakeoffDuration defines how long a takeoff operation physically lasts.
const TakeoffDuration = 5 * time.Second

// TakeOff prepares a plane for flight, simulates its takeoff, and updates the simulation state.
// It handles runway allocation, flight path generation, and state transitions for the plane and airport.
//
// Parameters:
//
//	plane: The Plane struct that is taking off. Note that this is passed by value;
//	       the modifications to this copy are then reflected when it's added to simState.PlanesInFlight.
//	simState: A pointer to the global SimulationState, allowing updates to the list of planes in flight.
//
// Returns:
//
//	*Flight: A pointer to the newly created Flight struct representing this takeoff.
//	error: An error if the takeoff cannot be initiated (e.g., no available runways, plane not found).
func (ap *Airport) TakeOff(plane Plane, simState *SimulationState) (*Flight, error) {
	// Acquire a lock only for checking/updating runway count
	ap.Mu.Lock()

	// Check if there's an available runway.
	if ap.Runway.noOfRunwayinUse >= ap.Runway.numberOfRunway {
		ap.Mu.Unlock() // Release lock immediately if no runway
		return nil, fmt.Errorf("airport %s has no available runways for takeoff (all %d in use)", ap.Serial, ap.Runway.numberOfRunway)
	}

	// Mark a runway as in use.
	ap.Runway.noOfRunwayinUse++
	ap.Mu.Unlock() // <<< IMPORTANT: Release the lock BEFORE the 3-second sleep

	// Simulate the physical takeoff duration. This does NOT hold the lock.
	// This allows other planes to acquire the lock and potentially start taking off
	// on another available runway immediately.
	time.Sleep(TakeoffDuration)

	// After the takeoff duration, re-acquire the lock to safely decrement the counter.
	ap.Mu.Lock()
	ap.Runway.noOfRunwayinUse--
	ap.Mu.Unlock() // Release the lock after updating

	// Find and remove the plane from this airport's list of parked planes.
	planeIndex := -1
	for i, p := range ap.Planes {
		if p.Serial == plane.Serial {
			planeIndex = i
			break
		}
	}

	if planeIndex == -1 {
		return nil, fmt.Errorf("plane %s not found at airport %s to initiate takeoff", plane.Serial, ap.Serial)
	}

	// Remove the plane from the airport's Planes slice.
	ap.Planes = append(ap.Planes[:planeIndex], ap.Planes[planeIndex+1:]...)

	// Select a random destination airport for the plane.
	destinationAirport, err := ap.getRandomDestinationAirport(simState.Airports)
	if err != nil {
		return nil, fmt.Errorf("failed to select destination airport for plane %s: %w", plane.Serial, err)
	}

	// Define the flight path from the current airport to the destination.
	flightPath := FlightPath{
		Depature: ap.Location,
		Arrival:  destinationAirport.Location,
	}

	// Calculate the total distance and estimated flight duration.
	flightDistance := Distance(flightPath.Depature, flightPath.Arrival)
	if plane.CruiseSpeed <= 0 {
		return nil, fmt.Errorf("plane %s has an invalid cruise speed (%.2f), cannot calculate flight duration", plane.Serial, plane.CruiseSpeed)
	}
	// Assuming CruiseSpeed is in units per second, and distance is in those same units.
	flightDuration := time.Duration(flightDistance/plane.CruiseSpeed) * time.Second

	takeoffTime := time.Now()
	landingTime := takeoffTime.Add(flightDuration)

	// Create a new Flight record with all its details.
	newFlight := Flight{
		FlightID:         util.GenerateSerialNumber(len(plane.FlightLog), "f"), // Generate unique ID for this specific flight
		FlightSchedule:   flightPath,
		TakeoffTime:      takeoffTime,
		LandingTime:      landingTime,
		CruisingAltitude: CruisingAltitude,
	}

	// Update the plane's internal state to reflect it's now in flight.
	plane.PlaneInFlight = true
	plane.FlightLog = append(plane.FlightLog, newFlight)

	// Add the updated plane to the global list of planes currently in flight.
	simState.PlanesInFlight = append(simState.PlanesInFlight, plane)

	fmt.Printf("Plane %s (Cruise Speed: %.2f) took off from Airport %s (%s), heading to Airport %s (%s). Estimated landing at %s.\n",
		plane.Serial, plane.CruiseSpeed, ap.Serial, ap.Location.String(), destinationAirport.Serial, destinationAirport.Location.String(), landingTime.Format("15:04:05"))

	return &newFlight, nil
}

// getRandomDestinationAirport selects a random airport from the list of all airports
// that is not the current airport (ap). This helps in simulating inter-airport travel.
func (ap *Airport) getRandomDestinationAirport(allAirports []*Airport) (*Airport, error) {
	eligibleAirports := []*Airport{}
	for _, otherAp := range allAirports {
		if otherAp.Serial != ap.Serial { // A plane cannot fly to the airport it just took off from
			eligibleAirports = append(eligibleAirports, otherAp)
		}
	}
	if len(eligibleAirports) == 0 {
		return nil, fmt.Errorf("no other airports available to serve as a destination")
	}

	// Use a new random source to ensure sufficient randomness, especially if called frequently.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(eligibleAirports))
	return eligibleAirports[randomIndex], nil
}
