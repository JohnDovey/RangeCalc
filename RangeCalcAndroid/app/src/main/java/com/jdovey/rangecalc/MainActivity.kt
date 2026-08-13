package com.jdovey.rangecalc

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import com.jdovey.rangecalc.ui.RangeCalcScreen
import com.jdovey.rangecalc.ui.theme.RangeCalcTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            RangeCalcTheme {
                RangeCalcScreen()
            }
        }
    }
}
