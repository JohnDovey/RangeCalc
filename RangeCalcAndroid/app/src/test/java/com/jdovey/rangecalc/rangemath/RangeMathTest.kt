package com.jdovey.rangecalc.rangemath

import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import kotlin.math.atan2
import kotlin.math.hypot
import kotlin.math.tan
import kotlin.math.PI

class RangeMathTest {

    @Test
    fun circularAngleDiff() {
        assertEquals(35.0, RangeMath.circularAngleDiff(25.0, 350.0), 1e-9)
        assertEquals(35.0, RangeMath.circularAngleDiff(350.0, 25.0), 1e-9)
        assertEquals(10.0, RangeMath.circularAngleDiff(10.0, 20.0), 1e-9)
        assertEquals(180.0, RangeMath.circularAngleDiff(0.0, 180.0), 1e-9)
        assertEquals(0.0, RangeMath.circularAngleDiff(1.0, 1.0), 1e-9)
    }

    @Test
    fun simpleReadmeExample() {
        // Tower 52/53 example: bearings 25° and 350°, baseline 250 m.
        // Fixed math: theta=35°, range = 250 * cot(35°) ≈ 357.04
        val r = RangeMath.simple(25.0, 350.0, 250.0).getOrThrow()
        assertEquals(35.0, r.angleAtTarget, 1e-9)

        val want = 250 * tan((90 - 35) * PI / 180)
        assertEquals(want, r.rangeFromRef1, 0.01)

        // Must not equal the old buggy value (degrees fed to tan as radians).
        val buggy = 250 * tan(90 - 325.0)
        assertTrue(kotlin.math.abs(r.rangeFromRef1 - buggy) >= 1)
    }

    @Test
    fun simpleRejectsLargeAngle() {
        assertTrue(RangeMath.simple(0.0, 100.0, 100.0).isFailure)
    }

    @Test
    fun fullKnownGeometry() {
        // R1(0,0), R2(100,0) baseline bearing 90° (east), T(50,100).
        // Compass bearing = atan2(east, north).
        var b1 = atan2(50.0, 100.0) * 180 / PI
        if (b1 < 0) b1 += 360
        var b2 = atan2(-50.0, 100.0) * 180 / PI
        if (b2 < 0) b2 += 360

        val r = RangeMath.full(b1, b2, 100.0, 90.0).getOrThrow()
        val expected = hypot(50.0, 100.0)
        assertEquals(expected, r.rangeFromRef1, 0.01)
        assertEquals(expected, r.rangeFromRef2, 0.01)
    }

    @Test
    fun fullIdenticalBearings() {
        assertTrue(RangeMath.full(45.0, 45.0, 100.0, 90.0).isFailure)
    }
}
