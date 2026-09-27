package com.jdovey.rangecalc.ui

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.res.painterResource
import com.jdovey.rangecalc.R

/**
 * The in-app branded splash, shown briefly after the OS-level cold-start
 * splash (see Theme.RangeCalc.Starting) while [RangeCalcScreen] is prepared.
 * Reuses the same illustration as the project's README/social-preview
 * artwork, without the phone-screenshot overlay.
 */
@Composable
fun SplashScreen() {
    AnimatedVisibility(visible = true, enter = fadeIn(), exit = fadeOut()) {
        Image(
            painter = painterResource(R.drawable.splash),
            contentDescription = null,
            contentScale = ContentScale.Crop,
            modifier = Modifier.fillMaxSize(),
        )
    }
}
