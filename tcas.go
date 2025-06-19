package main

import (
	"math"
)

// represents flight path as a 3D line segment
type flightPath struct {
	start coord
	end   coord
}

func Distance(p1, p2 coord) float64 {
	//resultant distance obtained by getting the magnitude of distance btw the two points
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2) + math.Pow(p1.Z-p2.Z, 2))
}
