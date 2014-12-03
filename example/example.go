package main

import (
	"log"
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
func ExportedError(x, y int) (int, error) {
	log.Println("exportederror")
	return 5, nil
}
