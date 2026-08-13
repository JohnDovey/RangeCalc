package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const ProgVer = "2.0.1 - Full Triangulation"

type TriangulationResult struct {
	RangeFromRef1 float64
	RangeFromRef2 float64
	AngleAtTarget float64
}

func main() {
	printHeader()

	bearing1 := getBearing("Bearing from Ref One to Target")
	bearing2 := getBearing("Bearing from Ref Two to Target")
	baseline := getPositiveDistance("Baseline Distance (Ref1 to Ref2 in meters)")

	fmt.Println("\nChoose calculation method:")
	fmt.Println("1. Simple (fast, assumes perpendicular baseline)")
	fmt.Println("2. Full Triangulation (more accurate with Baseline Bearing)")
	fmt.Print("Enter choice (1 or 2): ")
	choice := readInput()

	var result TriangulationResult
	if choice == "2" {
		baselineBearing := getBearing("Baseline Bearing (direction from Ref1 to Ref2)")
		result = calculateFullTriangulation(bearing1, bearing2, baseline, baselineBearing)
	} else {
		result = calculateSimple(bearing1, bearing2, baseline)
	}

	fmt.Println("\n=== RESULTS ===")
	fmt.Printf("Range from Ref One : %.2f meters\n", result.RangeFromRef1)
	fmt.Printf("Range from Ref Two : %.2f meters\n", result.RangeFromRef2)
	fmt.Printf("Angle at Target    : %.1f°\n", result.AngleAtTarget)
}

func printHeader() {
	fmt.Printf("+------------------------------[%s]----------------------------+\n", ProgVer)
	fmt.Println("|        Range Calculator - Full Triangulation Edition           |")
	fmt.Println("| View the code on GitHub [https://github.com/JohnDovey/RangeCalc] |")
	fmt.Println("+--------------------------------------------------------------------+\n")
	fmt.Println("Enter the values as prompted (0 to exit)\n")
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
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

		fmt.Printf("\t(%.0f meters)\n", value)
		return value
	}
}

func calculateSimple(b1, b2, baseline float64) TriangulationResult {
	diff := math.Abs(b1 - b2)
	if diff > 180 {
		diff = 360 - diff
	}
	if diff == 0 {
		fmt.Println("Error: Bearings are identical.")
		os.Exit(1)
	}

	rad := (90.0 - diff) * math.Pi / 180.0
	rangeVal := baseline * math.Tan(rad)

	return TriangulationResult{
		RangeFromRef1: rangeVal,
		RangeFromRef2: rangeVal,
		AngleAtTarget: diff,
	}
}

func calculateFullTriangulation(b1, b2, baseline, baselineBearing float64) TriangulationResult {
	angleAtTarget := math.Abs(b1 - b2)
	if angleAtTarget > 180 {
		angleAtTarget = 360 - angleAtTarget
	}

	if angleAtTarget == 0 {
		fmt.Println("Error: Bearings are identical.")
		os.Exit(1)
	}

	angleAtRef1 := math.Abs(baselineBearing - b1)
	if angleAtRef1 > 180 {
		angleAtRef1 = 360 - angleAtRef1
	}

	angleAtRef2 := 180.0 - angleAtTarget - angleAtRef1

	if angleAtRef2 <= 0 {
		fmt.Println("Error: Invalid geometry (target position not possible with these bearings).")
		os.Exit(1)
	}

	// Law of Sines
	sinTarget := math.Sin(angleAtTarget * math.Pi / 180)
	rangeFromRef1 := baseline * math.Sin(angleAtRef2*math.Pi/180) / sinTarget
	rangeFromRef2 := baseline * math.Sin(angleAtRef1*math.Pi/180) / sinTarget

	return TriangulationResult{
		RangeFromRef1: rangeFromRef1,
		RangeFromRef2: rangeFromRef2,
		AngleAtTarget: angleAtTarget,
	}
}
