package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const ProgVer = "2.0.1 - Full Triangulation"

// stdin is shared across all prompts. bufio.Reader buffers ahead of what it's
// asked for, so creating a fresh one per prompt (as this used to do) silently
// discards any input already buffered past the first line — piped input past
// the first prompt would vanish.
var stdin = bufio.NewReader(os.Stdin)

type TriangulationResult struct {
	RangeFromRef1 float64
	RangeFromRef2 float64
	AngleAtTarget float64
}

func main() {
	b1Flag := flag.Float64("b1", 0, "Bearing from Ref One to Target, degrees (1-360)")
	b2Flag := flag.Float64("b2", 0, "Bearing from Ref Two to Target, degrees (1-360)")
	dFlag := flag.Float64("d", 0, "Baseline distance between Ref One and Ref Two, meters")
	bpFlag := flag.Float64("bp", 0, "Baseline bearing, Ref One to Ref Two, degrees (1-360) — selects Full Triangulation")
	oFlag := flag.String("o", "", "Write a formatted report to this file instead of printing bare results to stdout")
	flag.Parse()

	set := map[string]bool{}
	flag.Visit(func(f *flag.Flag) { set[f.Name] = true })

	if len(set) > 0 {
		runNonInteractive(*b1Flag, *b2Flag, *dFlag, *bpFlag, *oFlag, set["b1"], set["b2"], set["d"], set["bp"])
		return
	}

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

// runNonInteractive handles the -b1/-b2/-d/-bp/-o flag-driven path, for
// scripting and piping. -b1, -b2, and -d are required together; -bp is
// optional and selects Full Triangulation in place of Simple mode; -o is
// optional and redirects a formatted report to a file instead of printing
// bare key=value results to stdout.
func runNonInteractive(b1, b2, baseline, bp float64, outFile string, b1Set, b2Set, dSet, bpSet bool) {
	if !b1Set || !b2Set || !dSet {
		fmt.Fprintln(os.Stderr, "Error: -b1, -b2, and -d are all required together.")
		flag.Usage()
		os.Exit(2)
	}
	if b1 < 1 || b1 > 360 {
		fmt.Fprintln(os.Stderr, "Error: -b1 must be between 1 and 360 degrees.")
		os.Exit(2)
	}
	if b2 < 1 || b2 > 360 {
		fmt.Fprintln(os.Stderr, "Error: -b2 must be between 1 and 360 degrees.")
		os.Exit(2)
	}
	if baseline <= 0 {
		fmt.Fprintln(os.Stderr, "Error: -d must be greater than zero.")
		os.Exit(2)
	}
	if bpSet && (bp < 1 || bp > 360) {
		fmt.Fprintln(os.Stderr, "Error: -bp must be between 1 and 360 degrees.")
		os.Exit(2)
	}

	method := "Simple"
	var result TriangulationResult
	if bpSet {
		method = "Full Triangulation"
		result = calculateFullTriangulation(b1, b2, baseline, bp)
	} else {
		result = calculateSimple(b1, b2, baseline)
	}

	if outFile == "" {
		fmt.Printf("RangeFromRef1=%.2f\n", result.RangeFromRef1)
		fmt.Printf("RangeFromRef2=%.2f\n", result.RangeFromRef2)
		fmt.Printf("AngleAtTarget=%.1f\n", result.AngleAtTarget)
		return
	}

	report := buildReport(b1, b2, baseline, bp, bpSet, method, result)
	if err := os.WriteFile(outFile, []byte(report), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing report to %s: %v\n", outFile, err)
		os.Exit(1)
	}
	fmt.Printf("Report written to %s\n", outFile)
}

const reportWidth = 72

func boxRule(width int) string {
	return "+" + strings.Repeat("-", width-2) + "+\n"
}

func boxCentered(text string, width int) string {
	pad := width - 2 - len(text)
	if pad < 0 {
		pad = 0
	}
	left := pad / 2
	right := pad - left
	return "|" + strings.Repeat(" ", left) + text + strings.Repeat(" ", right) + "|\n"
}

// buildReport formats a full result report for -o file output: an ASCII-box
// heading followed by the inputs used and the calculated results.
func buildReport(b1, b2, baseline, bp float64, bpSet bool, method string, result TriangulationResult) string {
	var sb strings.Builder
	sb.WriteString(boxRule(reportWidth))
	sb.WriteString(boxCentered(fmt.Sprintf("RangeCalc v%s", ProgVer), reportWidth))
	sb.WriteString(boxRule(reportWidth))
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("Method            : %s\n", method))
	sb.WriteString(fmt.Sprintf("Bearing 1         : %.1f°\n", b1))
	sb.WriteString(fmt.Sprintf("Bearing 2         : %.1f°\n", b2))
	sb.WriteString(fmt.Sprintf("Baseline Distance : %.1f m\n", baseline))
	if bpSet {
		sb.WriteString(fmt.Sprintf("Baseline Bearing  : %.1f°\n", bp))
	}
	sb.WriteString("\n")

	sb.WriteString(boxRule(reportWidth))
	sb.WriteString(boxCentered("RESULTS", reportWidth))
	sb.WriteString(boxRule(reportWidth))
	sb.WriteString(fmt.Sprintf("Range from Ref One : %.2f meters\n", result.RangeFromRef1))
	sb.WriteString(fmt.Sprintf("Range from Ref Two : %.2f meters\n", result.RangeFromRef2))
	sb.WriteString(fmt.Sprintf("Angle at Target    : %.1f°\n", result.AngleAtTarget))

	return sb.String()
}

func printHeader() {
	fmt.Printf("+------------------------------[%s]----------------------------+\n", ProgVer)
	fmt.Println("|        Range Calculator - Full Triangulation Edition           |")
	fmt.Println("| View the code on GitHub [https://github.com/JohnDovey/RangeCalc] |")
	fmt.Println("+--------------------------------------------------------------------+\n")
	fmt.Println("Enter the values as prompted (0 to exit)\n")
}

func readInput() string {
	input, err := stdin.ReadString('\n')
	trimmed := strings.TrimSpace(input)
	if err != nil && trimmed == "" {
		fmt.Println("\nNo more input available. Exiting.")
		os.Exit(0)
	}
	return trimmed
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
