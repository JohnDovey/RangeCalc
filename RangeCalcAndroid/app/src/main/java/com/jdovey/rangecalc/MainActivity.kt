package com.jdovey.rangecalc

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.animation.Crossfade
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.core.splashscreen.SplashScreen.Companion.installSplashScreen
import com.jdovey.rangecalc.ui.RangeCalcScreen
import com.jdovey.rangecalc.ui.SplashScreen
import com.jdovey.rangecalc.ui.theme.RangeCalcTheme
import kotlinx.coroutines.delay

/** How long the in-app illustrated splash (SplashScreen.kt) stays up. */
private const val SPLASH_DURATION_MS = 1100L

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        // Must be called before super.onCreate() / setContent().
        installSplashScreen()
        super.onCreate(savedInstanceState)

        setContent {
            var showSplash by remember { mutableStateOf(true) }
            LaunchedEffect(Unit) {
                delay(SPLASH_DURATION_MS)
                showSplash = false
            }

            RangeCalcTheme {
                Crossfade(targetState = showSplash, label = "splash-to-app") { splashVisible ->
                    if (splashVisible) SplashScreen() else RangeCalcScreen()
                }
            }
        }
    }
}
