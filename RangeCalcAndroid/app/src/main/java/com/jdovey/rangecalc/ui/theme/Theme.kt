package com.jdovey.rangecalc.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

/** Deep blue used for app chrome and primary actions. */
val RangeCalcBlue = Color(0xFF123C69)

/** Spring green used to highlight calculated results (matches the TUI's result colour). */
val RangeCalcGreen = Color(0xFF2AA876)

private val RangeCalcLightColorScheme = lightColorScheme(
    primary = RangeCalcBlue,
    onPrimary = Color.White,
    primaryContainer = RangeCalcBlue,
    onPrimaryContainer = Color.White,
    secondary = RangeCalcGreen,
    onSecondary = Color.White,
    tertiary = RangeCalcGreen,
    onTertiary = Color.White,
)

@Composable
fun RangeCalcTheme(content: @Composable () -> Unit) {
    MaterialTheme(
        colorScheme = RangeCalcLightColorScheme,
        content = content,
    )
}
