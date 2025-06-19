package main

import (
	"fmt"
	"time"
)

func main() {
	start()
}

func start() {
	baseTime := time.Date(2025, time.June, 19, 10, 0, 0, 0, time.UTC)

	flight1 := flight{
		flightSchedule: flightPath{depature: coord{-10, -10, 0}, arrival: coord{10, 10, 10}},
		takeoffTime:    baseTime,
		landingTime:    baseTime.Add(20 * time.Minute),
	}
	flight2 := flight{
		flightSchedule: flightPath{depature: coord{-5, 5, 5}, arrival: coord{5, -5, -5}},
		takeoffTime:    baseTime,
		landingTime:    baseTime.Add(10 * time.Minute),
	}

	//expectedF1ClosestTime := baseTime // (0/100) * 10min = 0min, so takeoff time
	//expectedF2ClosestTime := baseTime // (0/100) * 10min = 0min, so takeoff time
	//expectedDistanceBetweenPlanesCA := 5.0

	closestTime, distanceBetweenPlanesatCA := flight1.GetClosestApproachDetails(flight2)
	fmt.Println(closestTime)
	fmt.Println(distanceBetweenPlanesatCA)
}
