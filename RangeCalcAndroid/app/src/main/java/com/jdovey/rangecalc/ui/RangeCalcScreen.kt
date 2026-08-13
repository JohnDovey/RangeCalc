package com.jdovey.rangecalc.ui

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SegmentedButton
import androidx.compose.material3.SegmentedButtonDefaults
import androidx.compose.material3.SingleChoiceSegmentedButtonRow
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.TopAppBarDefaults
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.jdovey.rangecalc.rangemath.Method
import com.jdovey.rangecalc.ui.theme.RangeCalcGreen

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun RangeCalcScreen(viewModel: RangeCalcViewModel = viewModel()) {
    val uiState by viewModel.uiState.collectAsState()
    var showAbout by remember { mutableStateOf(false) }

    if (showAbout) {
        AboutDialog(onDismiss = { showAbout = false })
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("RangeCalc") },
                actions = {
                    TextButton(onClick = { showAbout = true }) {
                        Text("About", color = MaterialTheme.colorScheme.onPrimary)
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.primary,
                    titleContentColor = MaterialTheme.colorScheme.onPrimary,
                ),
            )
        },
    ) { padding ->
        LazyColumn(
            modifier = Modifier
                .fillMaxWidth()
                .padding(padding),
            contentPadding = PaddingValues(16.dp),
            verticalArrangement = Arrangement.spacedBy(16.dp),
        ) {
            item {
                InputCard(uiState = uiState, viewModel = viewModel)
            }

            uiState.error?.let { error ->
                item {
                    Text(
                        text = error,
                        color = MaterialTheme.colorScheme.error,
                        style = MaterialTheme.typography.bodyMedium,
                    )
                }
            }

            uiState.result?.let { result ->
                item {
                    ResultCard(result = result, method = uiState.method)
                }
            }

            if (uiState.history.isNotEmpty()) {
                item {
                    Text("History", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.Bold)
                }
                items(uiState.history) { entry ->
                    Column {
                        Text(entry.summary, fontWeight = FontWeight.Bold, style = MaterialTheme.typography.bodyMedium)
                        Text(entry.detail, style = MaterialTheme.typography.bodySmall)
                    }
                    HorizontalDivider(modifier = Modifier.padding(top = 8.dp))
                }
            }
        }
    }
}

@Composable
private fun InputCard(uiState: RangeCalcUiState, viewModel: RangeCalcViewModel) {
    Card(modifier = Modifier.fillMaxWidth()) {
        Column(modifier = Modifier.padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            SingleChoiceSegmentedButtonRow(modifier = Modifier.fillMaxWidth()) {
                SegmentedButton(
                    selected = uiState.method == Method.SIMPLE,
                    onClick = { viewModel.onMethodChange(Method.SIMPLE) },
                    shape = SegmentedButtonDefaults.itemShape(index = 0, count = 2),
                ) { Text("Simple") }
                SegmentedButton(
                    selected = uiState.method == Method.FULL,
                    onClick = { viewModel.onMethodChange(Method.FULL) },
                    shape = SegmentedButtonDefaults.itemShape(index = 1, count = 2),
                ) { Text("Full") }
            }

            OutlinedTextField(
                value = uiState.bearing1,
                onValueChange = viewModel::onBearing1Change,
                label = { Text("Bearing from Ref 1 (°)") },
                placeholder = { Text("e.g. 25") },
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Decimal),
                singleLine = true,
                modifier = Modifier.fillMaxWidth(),
            )

            OutlinedTextField(
                value = uiState.bearing2,
                onValueChange = viewModel::onBearing2Change,
                label = { Text("Bearing from Ref 2 (°)") },
                placeholder = { Text("e.g. 350") },
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Decimal),
                singleLine = true,
                modifier = Modifier.fillMaxWidth(),
            )

            OutlinedTextField(
                value = uiState.baseline,
                onValueChange = viewModel::onBaselineChange,
                label = { Text("Baseline distance (m)") },
                placeholder = { Text("e.g. 250") },
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Decimal),
                singleLine = true,
                modifier = Modifier.fillMaxWidth(),
            )

            AnimatedVisibility(visible = uiState.method == Method.FULL) {
                OutlinedTextField(
                    value = uiState.baselineBearing,
                    onValueChange = viewModel::onBaselineBearingChange,
                    label = { Text("Baseline bearing, Ref1→Ref2 (°)") },
                    placeholder = { Text("e.g. 90") },
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Decimal),
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth(),
                )
            }

            Row(horizontalArrangement = Arrangement.spacedBy(8.dp), modifier = Modifier.fillMaxWidth()) {
                Button(onClick = viewModel::calculate) { Text("Calculate") }
                OutlinedButton(onClick = { viewModel.clearForm(keepHistory = true) }) { Text("Clear") }
                TextButton(onClick = { viewModel.clearForm(keepHistory = false) }) { Text("Clear all") }
            }
        }
    }
}

@Composable
private fun ResultCard(result: com.jdovey.rangecalc.rangemath.RangeResult, method: Method) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = RangeCalcGreen.copy(alpha = 0.12f)),
    ) {
        Column(modifier = Modifier.padding(16.dp), verticalArrangement = Arrangement.spacedBy(4.dp)) {
            Text("Result", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.Bold)
            Text("Angle at target: %.2f°".format(result.angleAtTarget))
            if (method == Method.SIMPLE) {
                Text("Range: %.2f m".format(result.rangeFromRef1), style = MaterialTheme.typography.headlineSmall, color = RangeCalcGreen)
            } else {
                Text("Range from Ref 1: %.2f m".format(result.rangeFromRef1), style = MaterialTheme.typography.headlineSmall, color = RangeCalcGreen)
                Spacer(modifier = Modifier.height(2.dp))
                Text("Range from Ref 2: %.2f m".format(result.rangeFromRef2), style = MaterialTheme.typography.headlineSmall, color = RangeCalcGreen)
            }
        }
    }
}
