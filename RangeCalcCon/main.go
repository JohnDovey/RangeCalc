// Copyright (c) 2026 John Dovey <dovey.john@gmail.com>
//
// Command rangecalccon is a terminal UI that estimates range to a target
// from two compass bearings and a known baseline (simple or full triangulation).
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jdovey/rangecalccon/internal/ui"
)

// version is set at link time via -X main.version=… (see scripts/build-windows.ps1).
var version = "dev"

func main() {
	model := ui.New(version)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "rangecalccon: %v\n", err)
	os.Exit(1)
}
