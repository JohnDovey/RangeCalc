// Copyright (c) 2026 John Dovey <dovey.john@gmail.com>
//
// Package rangemath implements simple and full triangulation range estimates.
package rangemath

import (
	"fmt"
	"math"
)

// Method selects the calculation approach.
type Method int

const (
	// MethodSimple uses range = baseline * cot(theta) where theta is the
	// circular bearing difference. Approximate; assumes favourable geometry.
	MethodSimple Method = iota
	// MethodFull uses the law of sines with a known baseline bearing.
	MethodFull
)

// Result holds ranges from each reference and the angle at the target.
type Result struct {
	RangeFromRef1 float64
	RangeFromRef2 float64
	AngleAtTarget float64
	Method        Method
}

// CircularAngleDiff returns the smallest angle between two compass bearings
// in the open-closed interval (0, 180].
func CircularAngleDiff(a, b float64) float64 {
	diff := math.Abs(a - b)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

// Simple computes an approximate range from two bearings and baseline length.
// Valid when 0 < circular(|b1-b2|) < 90.
func Simple(b1, b2, baseline float64) (Result, error) {
	if baseline <= 0 {
		return Result{}, fmt.Errorf("baseline distance must be greater than zero")
	}
	theta := CircularAngleDiff(b1, b2)
	if theta == 0 {
		return Result{}, fmt.Errorf("bearings are identical — cannot calculate range")
	}
	if theta >= 90 {
		return Result{}, fmt.Errorf("angle difference too large for simple mode (need < 90°); try full triangulation")
	}

	// cot(theta) = tan(90 - theta)
	rangeVal := baseline * math.Tan(degToRad(90-theta))
	return Result{
		RangeFromRef1: rangeVal,
		RangeFromRef2: rangeVal,
		AngleAtTarget: theta,
		Method:        MethodSimple,
	}, nil
}

// Full computes ranges via the law of sines using baseline bearing (Ref1→Ref2).
func Full(b1, b2, baseline, baselineBearing float64) (Result, error) {
	if baseline <= 0 {
		return Result{}, fmt.Errorf("baseline distance must be greater than zero")
	}

	angleAtTarget := CircularAngleDiff(b1, b2)
	if angleAtTarget == 0 {
		return Result{}, fmt.Errorf("bearings are identical — cannot calculate range")
	}

	// Interior angle at Ref1: between baseline (Ref1→Ref2) and LOB (Ref1→Target)
	angleAtRef1 := CircularAngleDiff(baselineBearing, b1)
	angleAtRef2 := 180 - angleAtTarget - angleAtRef1

	if angleAtRef1 <= 0 || angleAtRef2 <= 0 {
		return Result{}, fmt.Errorf("invalid geometry — target not possible with these bearings")
	}

	sinTarget := math.Sin(degToRad(angleAtTarget))
	if math.Abs(sinTarget) < 1e-12 {
		return Result{}, fmt.Errorf("degenerate triangle (angle at target near 0° or 180°)")
	}

	r1 := baseline * math.Sin(degToRad(angleAtRef2)) / sinTarget
	r2 := baseline * math.Sin(degToRad(angleAtRef1)) / sinTarget
	if r1 <= 0 || r2 <= 0 {
		return Result{}, fmt.Errorf("invalid geometry — computed range is not positive")
	}

	return Result{
		RangeFromRef1: r1,
		RangeFromRef2: r2,
		AngleAtTarget: angleAtTarget,
		Method:        MethodFull,
	}, nil
}

// MethodName returns a short label for the method.
func MethodName(m Method) string {
	switch m {
	case MethodFull:
		return "Full triangulation"
	default:
		return "Simple"
	}
}
