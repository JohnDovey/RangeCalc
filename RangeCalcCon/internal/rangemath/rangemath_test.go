// Copyright (c) 2026 John Dovey <dovey.john@gmail.com>
package rangemath

import (
	"math"
	"testing"
)

func TestCircularAngleDiff(t *testing.T) {
	cases := []struct {
		a, b, want float64
	}{
		{25, 350, 35},
		{350, 25, 35},
		{10, 20, 10},
		{0, 180, 180},
		{1, 1, 0},
	}
	for _, tc := range cases {
		got := CircularAngleDiff(tc.a, tc.b)
		if math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("CircularAngleDiff(%v,%v)=%v want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestSimpleREADMEExample(t *testing.T) {
	// Tower 52/53 example: bearings 25° and 350°, baseline 250 m.
	// Fixed math: theta=35°, range = 250 * cot(35°) ≈ 357.04
	r, err := Simple(25, 350, 250)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(r.AngleAtTarget-35) > 1e-9 {
		t.Errorf("angle=%v want 35", r.AngleAtTarget)
	}
	want := 250 * math.Tan((90-35)*math.Pi/180)
	if math.Abs(r.RangeFromRef1-want) > 0.01 {
		t.Errorf("range=%v want %v", r.RangeFromRef1, want)
	}
	// Must not equal the old buggy value (degrees fed to tan as radians).
	buggy := 250 * math.Tan(90-325)
	if math.Abs(r.RangeFromRef1-buggy) < 1 {
		t.Errorf("range still matches old buggy value %v", buggy)
	}
}

func TestSimpleRejectsLargeAngle(t *testing.T) {
	_, err := Simple(0, 100, 100)
	if err == nil {
		t.Fatal("expected error for theta >= 90")
	}
}

func TestFullKnownGeometry(t *testing.T) {
	// R1(0,0), R2(100,0) baseline bearing 90° (east), T(50,100).
	// Compass bearing = atan2(east, north).
	b1 := math.Atan2(50, 100) * 180 / math.Pi
	if b1 < 0 {
		b1 += 360
	}
	b2 := math.Atan2(-50, 100) * 180 / math.Pi
	if b2 < 0 {
		b2 += 360
	}
	r, err := Full(b1, b2, 100, 90)
	if err != nil {
		t.Fatal(err)
	}
	expected := math.Hypot(50, 100)
	if math.Abs(r.RangeFromRef1-expected) > 0.01 {
		t.Errorf("r1=%v want %v", r.RangeFromRef1, expected)
	}
	if math.Abs(r.RangeFromRef2-expected) > 0.01 {
		t.Errorf("r2=%v want %v", r.RangeFromRef2, expected)
	}
}

func TestFullIdenticalBearings(t *testing.T) {
	_, err := Full(45, 45, 100, 90)
	if err == nil {
		t.Fatal("expected error")
	}
}
