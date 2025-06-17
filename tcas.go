package main

import (
	"math"
	"time"
)

func Distance(p1, p2 Coord) float64 {
	//resultant distance obtained by getting the magnitude of distance btw the two points
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2) + math.Pow(p1.Z-p2.Z, 2))
}

func FindClosestApproach(f1, f2 Flight) (time.Time, float64, error) {

}
