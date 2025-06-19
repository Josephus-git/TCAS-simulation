package main

import (
	"fmt"
	"math"
	"testing"
)

// Test helpers for float comparison
const FloatEqualityThreshold = 1e-5

func FloatEquals(a, b float64) bool {
	return math.Abs(a-b) < FloatEqualityThreshold
}

func CoordEquals(c1, c2 coord) bool {
	return FloatEquals(c1.X, c2.X) && FloatEquals(c1.Y, c2.Y) && FloatEquals(c1.Z, c2.Z)
}

func TestFindClosestApproachDuringTransit(t *testing.T) {
	tests := []struct {
		name    string
		fp1     flightPath
		fp2     flightPath
		wantFp1 coord
		wantFp2 coord
	}{
		{
			name: "Intersecting Paths",
			fp1: flightPath{
				start: coord{X: 0, Y: 0, Z: 0},
				end:   coord{X: 10, Y: 0, Z: 0},
			},
			fp2: flightPath{
				start: coord{X: 5, Y: -5, Z: 0},
				end:   coord{X: 5, Y: 5, Z: 0},
			},
			wantFp1: coord{X: 5, Y: 0, Z: 0},
			wantFp2: coord{X: 5, Y: 0, Z: 0},
		},
		{
			name: "Parallel Paths (non-overlapping)",
			fp1: flightPath{
				start: coord{X: 0, Y: 0, Z: 0},
				end:   coord{X: 10, Y: 0, Z: 0},
			},
			fp2: flightPath{
				start: coord{X: 0, Y: 1, Z: 0},
				end:   coord{X: 10, Y: 1, Z: 0},
			},
			wantFp1: coord{X: 0, Y: 0, Z: 0},
			wantFp2: coord{X: 0, Y: 1, Z: 0},
		},
		{
			name: "Skew Paths (non-intersecting, 3D)",
			fp1: flightPath{
				start: coord{X: 0, Y: 0, Z: 0},
				end:   coord{X: 10, Y: 0, Z: 0},
			},
			fp2: flightPath{
				start: coord{X: 0, Y: 10, Z: 10},
				end:   coord{X: 10, Y: 10, Z: 0},
			},
			wantFp1: coord{X: 10, Y: 0, Z: 0},
			wantFp2: coord{X: 10, Y: 10, Z: 0},
		},
		{
			name: "Endpoint to Endpoint (closest is an end point)",
			fp1: flightPath{
				start: coord{X: 0, Y: 0, Z: 0},
				end:   coord{X: 1, Y: 0, Z: 0},
			},
			fp2: flightPath{
				start: coord{X: 10, Y: 0, Z: 0},
				end:   coord{X: 11, Y: 0, Z: 0},
			},
			wantFp1: coord{X: 1, Y: 0, Z: 0},
			wantFp2: coord{X: 10, Y: 0, Z: 0},
		},
		{
			name: "Identical Paths (should return start points)",
			fp1: flightPath{
				start: coord{X: 0, Y: 0, Z: 0},
				end:   coord{X: 10, Y: 0, Z: 0},
			},
			fp2: flightPath{
				start: coord{X: 0, Y: 0, Z: 0},
				end:   coord{X: 10, Y: 0, Z: 0},
			},
			wantFp1: coord{X: 0, Y: 0, Z: 0},
			wantFp2: coord{X: 0, Y: 0, Z: 0},
		},
		{
			name: "Collinear, Overlapping (Segment 1 contains Segment 2)",
			fp1: flightPath{
				start: coord{X: 0, Y: 0, Z: 0},
				end:   coord{X: 10, Y: 0, Z: 0},
			},
			fp2: flightPath{
				start: coord{X: 2, Y: 0, Z: 0},
				end:   coord{X: 8, Y: 0, Z: 0},
			},
			wantFp1: coord{X: 2, Y: 0, Z: 0},
			wantFp2: coord{X: 2, Y: 0, Z: 0},
		},
		{
			name: "Collinear, Non-overlapping (Segment 1 before Segment 2)",
			fp1: flightPath{
				start: coord{X: 0, Y: 0, Z: 0},
				end:   coord{X: 5, Y: 0, Z: 0},
			},
			fp2: flightPath{
				start: coord{X: 7, Y: 0, Z: 0},
				end:   coord{X: 10, Y: 0, Z: 0},
			},
			wantFp1: coord{X: 5, Y: 0, Z: 0},
			wantFp2: coord{X: 7, Y: 0, Z: 0},
		},
		{
			name: "Perpendicular, not intersecting, one endpoint is closest",
			fp1: flightPath{
				start: coord{X: 0, Y: 0, Z: 0},
				end:   coord{X: 5, Y: 0, Z: 0},
			},
			fp2: flightPath{
				start: coord{X: 0, Y: 5, Z: 0},
				end:   coord{X: 0, Y: 10, Z: 0},
			},
			wantFp1: coord{X: 0, Y: 0, Z: 0},
			wantFp2: coord{X: 0, Y: 5, Z: 0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotFp1, gotFp2 := FindClosestApprachDuringTransit(test.fp1, test.fp2)

			if !CoordEquals(gotFp1, test.wantFp1) {
				t.Errorf("\nFor %s\n fp1Closest: got %v, want %v\n", test.name, gotFp1, test.wantFp1)
			} else {
				fmt.Printf("\n\nSuccess For %s\n fp1Closest: got %v, want %v\n", test.name, gotFp1, test.wantFp1)
			}

			if !CoordEquals(gotFp2, test.wantFp2) {
				t.Errorf("\nFor %s\n fp2Closest: got %v, want %v\n", test.name, gotFp2, test.wantFp2)
			} else {
				fmt.Printf("\nSuccess For %s\n fp2Closest: got %v, want %v\n", test.name, gotFp2, test.wantFp2)
			}
		})
	}

}
