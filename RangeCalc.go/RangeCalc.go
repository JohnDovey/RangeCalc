package main

import (
	"fmt"
	"math"
)

// main function
func main() {
	var bearing1 float64      // bearing from left ref point to target
	var bearing2 float64      // bearing from right ref point to target
	var refdistance float64   // distance between two ref points _ use %u for unsigned number
	var rangetotarget float64 // Calculated range
	var tmpAngle float64

	// Println function is used to
	// display output in the next line
	fmt.Println("Enter the Compass bearing from the first point: ")

	// Taking input from user
	fmt.Scanln(&bearing1)
	fmt.Printf("\t(%v degrees bearing Ref 1 to target) ", bearing1)
	fmt.Println("")

	if bearing1 > 360 {
		fmt.Println("Error! No more than 360 Degrees allowed")
		panic(0)
	}
	if bearing1 < 1 {
		fmt.Println("Error! Bearing must be greater than Zero")
		panic(1)
	}

	fmt.Println("Enter bearing from second point to target: ", bearing2)
	fmt.Scanln(&bearing2)
	fmt.Printf("\t(%v degrees bearing Ref 2 to target) ", bearing2)
	fmt.Println("")

	if bearing2 > 360 {
		fmt.Println("Error! No more than 360 Degrees allowed")
		panic(3)
	}
	if bearing2 < 1 {
		fmt.Println("Error! Bearing must be greater than Zero")
		panic(4)
	}

	// Separation Distance
	fmt.Println("Distance between Ref One and Ref Two: ")
	fmt.Scanln(&refdistance)

	fmt.Printf("\t(%v meters between reference points) ", refdistance)
	fmt.Println("")

	// d = (Tan (90 - (A -B))) x Ref
	tmpAngle = bearing1 - bearing2
	if tmpAngle < 0 {
		tmpAngle = tmpAngle * (0 - 1)
	}

	rangetotarget = (math.Tan(90 - tmpAngle)) * refdistance

	fmt.Printf("\r\nRange to Target: %f\r\n", rangetotarget)
	fmt.Println("")

}
