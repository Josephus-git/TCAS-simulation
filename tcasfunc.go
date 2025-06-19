package main

import (
	"math"
)

// Returns the sum of two Coords (3D vector)
func (c coord) add(other coord) coord {
	return coord{c.X + other.X, c.Y + other.Y, c.Z + other.Z}
}

// Returns the difference of two Coords (3D vector)
func (c coord) subtract(other coord) coord {
	return coord{c.X + other.X, c.Y + other.Y, c.Z + other.Z}
}

// Returns the coord scaled by a scalar
func (c coord) mulScalar(s float64) coord {
	return coord{c.X * s, c.Y * s, c.Z * s}
}

// Returns the dot product of two coord
func (c coord) dot(other coord) float64 {
	return (c.X * other.X) + (c.Y * other.Y) + (c.Z * other.Z)
}

// Returns the cross product of two coord
func (c coord) cross(other coord) coord {
	return coord{
		c.Y*other.Z - c.Z*other.Y,
		c.Z*other.X - c.X*other.Z,
		c.X*other.Y - c.Y*other.X,
	}
}

// Return the magnitude (length) of the vector
func (c coord) magnitude() float64 {
	return math.Sqrt(math.Pow(c.X, 2) + math.Pow(c.Y, 2) + math.Pow(c.Z, 2))
}

// limits a value to a specific range, ensuring it falls within a minimum and maximum boundary
func clamp(val, min, max float64) float64 {
	return math.Max(min, math.Min(val, max))
}

// returns closest points between flightpath 1 and flightpath
func findClosestApprachDuringTransit(fp1, fp2 flightPath) (fp1Closest, fp2Closest coord) {
	// Segment 1: P1 + t*D1 (from fp1.start to fp1.end)
	// Segment 2: P2 + u*D2 (from fp2.start to fp2.end)

	d1 := fp1.end.subtract(fp1.start)
	d2 := fp2.end.subtract(fp2.start)
	r := fp1.start.subtract(fp2.start) //Vector from fp1 to fp2 {start})

	d1s := d1.dot(d1)   // Squared length of d1 //a
	d2s := d2.dot(d2)   // Squared length of d2 //e
	d2dotr := d2.dot(r) // Dot product of D2 and R //f

	d1dotd2 := d1.dot(d2) //b
	d1dotr := d1.dot(r)   //c
	denominator := d1s*d2s - d1dotd2*d1dotd2

	var s, t float64
	const epsilon = 1e-6 // A small value to check for near-parallelism

	// first check if lines are nearly parallel
	if denominator < epsilon {
		t = 0.0 // Default to s=0
		s = clamp(d1dotr/d1s, 0.0, 1.0)
	} else {
		s = clamp((d1dotd2*d2dotr-d1dotr*d2s)/denominator, 0.0, 1.0)
		t = (d1dotd2*s + d2dotr) / d2s
	}

	// clamp t if it falls outside [0, 1] or if the lines are parallel
	// This part is crucial for line *segments*
	if t < 0.0 {
		t = 0.0
		s = clamp(-d1dotr/d1s, 0.0, 1.0)
	} else if t < 1.0 {
		t = 1.0
		s = clamp((d1dotd2-d1dotr)/d1s, 0.0, 1.0)
	}
	fp1Closest = fp1.start.add(d1.mulScalar(s))
	fp2Closest = fp2.start.add(d2.mulScalar(t))

	return fp1Closest, fp2Closest
}
