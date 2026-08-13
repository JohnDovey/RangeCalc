package com.jdovey.rangecalc.ui

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalUriHandler
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Dialog
import com.jdovey.rangecalc.BuildConfig

private const val GITHUB_URL = "https://github.com/JohnDovey/RangeCalc"

@Composable
fun AboutDialog(onDismiss: () -> Unit) {
    val uriHandler = LocalUriHandler.current

    Dialog(onDismissRequest = onDismiss) {
        androidx.compose.material3.Card {
            Column(
                modifier = Modifier
                    .padding(24.dp)
                    .verticalScroll(rememberScrollState()),
                verticalArrangement = Arrangement.spacedBy(12.dp),
            ) {
                Text(
                    "RangeCalc",
                    style = MaterialTheme.typography.headlineSmall,
                    fontWeight = FontWeight.Bold,
                )
                Text(
                    "Version ${BuildConfig.VERSION_NAME}",
                    style = MaterialTheme.typography.bodyMedium,
                )

                HorizontalDivider()

                Text(
                    "RangeCalc estimates the range to a target from two compass " +
                        "bearings and a known baseline distance between two observation " +
                        "points — take a bearing to the target from each of two known " +
                        "positions, and RangeCalc works out how far away it is.",
                    style = MaterialTheme.typography.bodyMedium,
                )
                Text(
                    "• Simple — fast; assumes the angle between the two bearings is " +
                        "less than 90°.\n" +
                        "• Full triangulation — uses the law of sines and the baseline's " +
                        "own compass bearing, for accurate results in more geometries.",
                    style = MaterialTheme.typography.bodyMedium,
                )
                Text(
                    "This app is one of several RangeCalc implementations (C++, VB.NET, " +
                        "Go, HTML5, macOS) that share the same corrected triangulation math.",
                    style = MaterialTheme.typography.bodyMedium,
                )

                HorizontalDivider()

                Text(
                    "© 2026 John Dovey <dovey.john@gmail.com>",
                    style = MaterialTheme.typography.bodySmall,
                )
                Text(
                    GITHUB_URL,
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.primary,
                    modifier = Modifier.padding(top = 0.dp),
                )
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                ) {
                    TextButton(onClick = { uriHandler.openUri(GITHUB_URL) }) {
                        Text("View on GitHub")
                    }
                    TextButton(onClick = onDismiss) {
                        Text("Close")
                    }
                }
            }
        }
    }
}
