package com.jdovey.rangecalc.ui

import androidx.lifecycle.ViewModel
import com.jdovey.rangecalc.rangemath.Method
import com.jdovey.rangecalc.rangemath.RangeMath
import com.jdovey.rangecalc.rangemath.RangeResult
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import java.util.Locale

data class HistoryEntry(val summary: String, val detail: String)

data class RangeCalcUiState(
    val method: Method = Method.SIMPLE,
    val bearing1: String = "",
    val bearing2: String = "",
    val baseline: String = "",
    val baselineBearing: String = "",
    val result: RangeResult? = null,
    val error: String? = null,
    val history: List<HistoryEntry> = emptyList(),
)

private fun fmt(v: Double): String = String.format(Locale.US, "%.2f", v)

class RangeCalcViewModel : ViewModel() {

    private val _uiState = MutableStateFlow(RangeCalcUiState())
    val uiState: StateFlow<RangeCalcUiState> = _uiState.asStateFlow()

    fun onMethodChange(method: Method) {
        _uiState.update { it.copy(method = method, error = null) }
    }

    fun onBearing1Change(value: String) {
        _uiState.update { it.copy(bearing1 = value, error = null) }
    }

    fun onBearing2Change(value: String) {
        _uiState.update { it.copy(bearing2 = value, error = null) }
    }

    fun onBaselineChange(value: String) {
        _uiState.update { it.copy(baseline = value, error = null) }
    }

    fun onBaselineBearingChange(value: String) {
        _uiState.update { it.copy(baselineBearing = value, error = null) }
    }

    fun calculate() {
        val s = _uiState.value
        val b1 = s.bearing1.toDoubleOrNull()
        val b2 = s.bearing2.toDoubleOrNull()
        val baseline = s.baseline.toDoubleOrNull()

        if (b1 == null || b2 == null || baseline == null) {
            _uiState.update { it.copy(error = "Enter valid numbers for both bearings and the baseline", result = null) }
            return
        }

        val outcome = when (s.method) {
            Method.SIMPLE -> RangeMath.simple(b1, b2, baseline)
            Method.FULL -> {
                val bb = s.baselineBearing.toDoubleOrNull()
                if (bb == null) {
                    _uiState.update { it.copy(error = "Enter a baseline bearing for full triangulation", result = null) }
                    return
                }
                RangeMath.full(b1, b2, baseline, bb)
            }
        }

        outcome.fold(
            onSuccess = { result ->
                val entry = HistoryEntry(
                    summary = "${RangeMath.methodName(result.method)} · ${fmt(result.rangeFromRef1)} m",
                    detail = buildString {
                        append("Bearings ${fmt(b1)}°/${fmt(b2)}°, baseline ${fmt(baseline)} m")
                        if (s.method == Method.FULL) append(", baseline bearing ${s.baselineBearing}°")
                        append(" — angle at target ${fmt(result.angleAtTarget)}°")
                        if (result.rangeFromRef1 != result.rangeFromRef2) {
                            append(", range from Ref2 ${fmt(result.rangeFromRef2)} m")
                        }
                    },
                )
                _uiState.update {
                    it.copy(result = result, error = null, history = listOf(entry) + it.history)
                }
            },
            onFailure = { e ->
                _uiState.update { it.copy(error = e.message ?: "Calculation failed", result = null) }
            },
        )
    }

    fun clearForm(keepHistory: Boolean = true) {
        _uiState.update {
            it.copy(
                bearing1 = "",
                bearing2 = "",
                baseline = "",
                baselineBearing = "",
                result = null,
                error = null,
                history = if (keepHistory) it.history else emptyList(),
            )
        }
    }
}
