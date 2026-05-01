package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const ProgVer = "1.0.13"

func main() {
	printHeader()

	fmt.Println("Enter the values as prompted (0 to exit)\n")

	bearing1 := getBearing("Bearing from Ref One to Target")
	bearing2 := getBearing("Bearing from Ref Two to Target")
	refDistance := getPositiveDistance("Distance between Ref One and Ref Two (meters)")

	rangeToTarget := calculateRange(bearing1, bearing2, refDistance)

	fmt.Println()
	fmt.Printf("Range to Target: %.2f meters\n\n", rangeToTarget)
}

func printHeader() {
	fmt.Printf("+------------------------------[%s]----------------------------+\n", ProgVer)
	fmt.Println("|        Range Calculator - by John Dovey <dovey.john@gmail.com>     |")
	fmt.Println("| View the code on GitHub [https://github.com/JohnDovey/RangeCalc]   |")
	fmt.Printf("+--------------------------------------------------------------------+\n\n")
}

func getBearing(prompt string) float64 {
	for {
		fmt.Printf("%s: ", prompt)
		input := readInput()

		if input == "0" {
			fmt.Println("Operation cancelled.")
			os.Exit(0)
		}

		value, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		if value < 1 || value > 360 {
			fmt.Println("Error: Bearing must be between 1 and 360 degrees.")
			continue
		}

		fmt.Printf("\t(%.0f°)\n", value)
		return value
	}
}

func getPositiveDistance(prompt string) float64 {
	for {
		fmt.Printf("%s: ", prompt)
		input := readInput()

		if input == "0" {
			fmt.Println("Operation cancelled.")
			os.Exit(0)
		}

		value, err := strconv.ParseFloat(input, 64)
		if err != nil || value < 1 {
			fmt.Println("Error: Please enter a positive number greater than zero.")
			continue
		}

		fmt.Printf("\t(%.0f meters between reference points)\n", value)
		return value
	}
}

func calculateRange(b1, b2, distance float64) float64 {
	diff := math.Abs(b1 - b2)

	if diff == 0 {
		fmt.Println("Error: Bearings are identical - cannot calculate range.")
		os.Exit(1)
	}
	if diff >= 90 {
		fmt.Println("Error: Angle difference too large (target likely behind baseline).")
		os.Exit(1)
	}

	// FIXED: Convert degrees to radians
	radians := (90.0 - diff) * math.Pi / 180.0
	return distance * math.Tan(radians)
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
