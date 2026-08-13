package com.jdovey.rangecalc.rangemath

import kotlin.math.abs
import kotlin.math.sin
import kotlin.math.tan
import kotlin.math.PI

/** Selects the calculation approach. */
enum class Method {
    /** range = baseline * cot(theta) where theta is the circular bearing difference. */
    SIMPLE,

    /** Law of sines with a known baseline bearing. */
    FULL,
}

/** Ranges from each reference and the angle at the target. */
data class RangeResult(
    val rangeFromRef1: Double,
    val rangeFromRef2: Double,
    val angleAtTarget: Double,
    val method: Method,
)

/** Simple and full triangulation range estimates. Ported from RangeCalcCon's rangemath package. */
object RangeMath {

    /** Smallest angle between two compass bearings, in the open-closed interval (0, 180]. */
    fun circularAngleDiff(a: Double, b: Double): Double {
        var diff = abs(a - b)
        if (diff > 180) diff = 360 - diff
        return diff
    }

    private fun degToRad(deg: Double): Double = deg * PI / 180

    /** Approximate range from two bearings and baseline length. Valid when 0 < theta < 90. */
    fun simple(b1: Double, b2: Double, baseline: Double): Result<RangeResult> {
        if (baseline <= 0) {
            return Result.failure(IllegalArgumentException("baseline distance must be greater than zero"))
        }
        val theta = circularAngleDiff(b1, b2)
        if (theta == 0.0) {
            return Result.failure(IllegalArgumentException("bearings are identical — cannot calculate range"))
        }
        if (theta >= 90) {
            return Result.failure(
                IllegalArgumentException("angle difference too large for simple mode (need < 90°); try full triangulation")
            )
        }

        // cot(theta) = tan(90 - theta)
        val rangeVal = baseline * tan(degToRad(90 - theta))
        return Result.success(
            RangeResult(
                rangeFromRef1 = rangeVal,
                rangeFromRef2 = rangeVal,
                angleAtTarget = theta,
                method = Method.SIMPLE,
            )
        )
    }

    /** Ranges via the law of sines using baseline bearing (Ref1→Ref2). */
    fun full(b1: Double, b2: Double, baseline: Double, baselineBearing: Double): Result<RangeResult> {
        if (baseline <= 0) {
            return Result.failure(IllegalArgumentException("baseline distance must be greater than zero"))
        }

        val angleAtTarget = circularAngleDiff(b1, b2)
        if (angleAtTarget == 0.0) {
            return Result.failure(IllegalArgumentException("bearings are identical — cannot calculate range"))
        }

        // Interior angle at Ref1: between baseline (Ref1→Ref2) and LOB (Ref1→Target)
        val angleAtRef1 = circularAngleDiff(baselineBearing, b1)
        val angleAtRef2 = 180 - angleAtTarget - angleAtRef1

        if (angleAtRef1 <= 0 || angleAtRef2 <= 0) {
            return Result.failure(IllegalArgumentException("invalid geometry — target not possible with these bearings"))
        }

        val sinTarget = sin(degToRad(angleAtTarget))
        if (abs(sinTarget) < 1e-12) {
            return Result.failure(IllegalArgumentException("degenerate triangle (angle at target near 0° or 180°)"))
        }

        val r1 = baseline * sin(degToRad(angleAtRef2)) / sinTarget
        val r2 = baseline * sin(degToRad(angleAtRef1)) / sinTarget
        if (r1 <= 0 || r2 <= 0) {
            return Result.failure(IllegalArgumentException("invalid geometry — computed range is not positive"))
        }

        return Result.success(
            RangeResult(
                rangeFromRef1 = r1,
                rangeFromRef2 = r2,
                angleAtTarget = angleAtTarget,
                method = Method.FULL,
            )
        )
    }

    /** Short label for the method. */
    fun methodName(m: Method): String = when (m) {
        Method.FULL -> "Full triangulation"
        Method.SIMPLE -> "Simple"
    }
}
