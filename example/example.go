package main

import (
	"log"
	"math"
)

type hello struct{}

// dist returns the distance of the vector with x, y, and z components
func Dist(x, y, z int) int {
	log.Println("dist")
	return 5
}

// Exported is an exported function
func Exported(x, y, z int) int {
	log.Println("exported")
	return 5
}

/*
	first line of
	a multi-line ocomment
		with code
*/
func ExportedError(x, y, z, a, b, c, alongername int) (int, error) {
	log.Println("exportederror")
	return 5, nil
}

func Inductance(d, l, n float64) float64 {
	µ0 := 4 * math.Pi * math.Pow10(-7)
	A := math.Pow(d/2, 2) * math.Pi

	I := µ0 * (math.Pow(n, 2) * A) / l

	// convert from H to µH
	return I * math.Pow10(6)
}
