// Copyright (c) 2026 John Dovey <dovey.john@gmail.com>
//
// Package ui is the bubbletea + lipgloss terminal UI for RangeCalcCon.
// Visual language matches GetWattPad / ServiceMonitor (title bar, rounded
// panels, keyboard help strip).
package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jdovey/rangecalccon/internal/rangemath"
)

const leftWidth = 42

const (
	fieldB1 = iota
	fieldB2
	fieldBaseline
	fieldBaselineBearing
	fieldCount
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("24")).
			Padding(0, 1)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	panelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("117"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("62")).
			Bold(true)

	normalItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	upStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("78")).Bold(true)
	errStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	helpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	dimStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	keyStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("117")).Bold(true)
	boldStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("231"))
	resultStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("78"))
)

type tickMsg time.Time

// historyEntry is one past calculation shown in the results panel.
type historyEntry struct {
	summary string
	detail  string
}

// Model is the root bubbletea model.
type Model struct {
	version string

	method rangemath.Method
	focus  int // fieldB1 … fieldBaselineBearing, or -1 when method row focused
	inputs []textinput.Model

	width  int
	height int
	ready  bool

	result  *rangemath.Result
	history []historyEntry
	detail  viewport.Model

	flash   string
	flashAt time.Time
	status  string
}

// New constructs the UI model.
func New(version string) Model {
	mk := func(placeholder string) textinput.Model {
		ti := textinput.New()
		ti.Placeholder = placeholder
		ti.CharLimit = 16
		ti.Width = 18
		ti.Prompt = ""
		return ti
	}

	inputs := make([]textinput.Model, fieldCount)
	inputs[fieldB1] = mk("e.g. 25")
	inputs[fieldB2] = mk("e.g. 350")
	inputs[fieldBaseline] = mk("e.g. 250")
	inputs[fieldBaselineBearing] = mk("e.g. 90")

	// Start on first numeric field with focus.
	inputs[fieldB1].Focus()

	return Model{
		version: version,
		method:  rangemath.MethodSimple,
		focus:   fieldB1,
		inputs:  inputs,
		detail:  viewport.New(60, 20),
		status:  "ready",
		history: nil,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tickCmd())
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layout()
		m.ready = true
		m.refreshDetail()
		return m, nil

	case tickMsg:
		if time.Since(m.flashAt) > 5*time.Second {
			m.flash = ""
		}
		return m, tickCmd()

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Forward to focused text input for typing / cursor blink.
	if m.focus >= 0 && m.focus < fieldCount {
		var cmd tea.Cmd
		m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "tab", "down", "j":
		return m, m.moveFocus(1)

	case "shift+tab", "up", "k":
		return m, m.moveFocus(-1)

	case "1":
		// Method shortcuts only when the method row is focused so typing
		// digits into bearing/distance fields still works.
		if m.focus < 0 {
			m.method = rangemath.MethodSimple
			m.flash = upStyle.Render("✓ simple mode")
			m.flashAt = time.Now()
			m.refreshDetail()
			return m, nil
		}

	case "2":
		if m.focus < 0 {
			m.method = rangemath.MethodFull
			m.flash = upStyle.Render("✓ full triangulation")
			m.flashAt = time.Now()
			m.refreshDetail()
			return m, nil
		}

	case "left", "h":
		if m.focus < 0 {
			m.method = rangemath.MethodSimple
			m.flash = upStyle.Render("✓ simple mode")
			m.flashAt = time.Now()
			m.refreshDetail()
			return m, nil
		}

	case "right", "l":
		if m.focus < 0 {
			m.method = rangemath.MethodFull
			m.flash = upStyle.Render("✓ full triangulation")
			m.flashAt = time.Now()
			m.refreshDetail()
			return m, nil
		}

	case "m":
		// Toggle method.
		if m.method == rangemath.MethodSimple {
			m.method = rangemath.MethodFull
			m.flash = upStyle.Render("✓ full triangulation")
		} else {
			m.method = rangemath.MethodSimple
			m.flash = upStyle.Render("✓ simple mode")
			// Baseline bearing field is hidden in simple mode.
			if m.focus == fieldBaselineBearing {
				m.setFocus(fieldBaseline)
			}
		}
		m.flashAt = time.Now()
		m.refreshDetail()
		return m, nil

	case "enter", "c":
		m.calculate()
		m.refreshDetail()
		return m, nil

	case "r":
		// Clear form + current result (keep history).
		for i := range m.inputs {
			m.inputs[i].SetValue("")
		}
		m.result = nil
		m.status = "ready"
		m.flash = dimStyle.Render("form cleared")
		m.flashAt = time.Now()
		m.setFocus(fieldB1)
		m.refreshDetail()
		return m, textinput.Blink

	case "R":
		// Clear history as well.
		for i := range m.inputs {
			m.inputs[i].SetValue("")
		}
		m.result = nil
		m.history = nil
		m.status = "ready"
		m.flash = dimStyle.Render("form + history cleared")
		m.flashAt = time.Now()
		m.setFocus(fieldB1)
		m.refreshDetail()
		return m, textinput.Blink

	case "pgup", "ctrl+u":
		m.detail.HalfViewUp()
		return m, nil
	case "pgdown", "ctrl+d":
		m.detail.HalfViewDown()
		return m, nil
	}

	// Typing into focused field.
	if m.focus >= 0 && m.focus < fieldCount {
		var cmd tea.Cmd
		m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) moveFocus(delta int) tea.Cmd {
	// Focus order: method (-1) → b1 → b2 → baseline → [baselineBearing if full]
	order := m.focusOrder()
	// Find current index in order
	idx := 0
	for i, f := range order {
		if f == m.focus {
			idx = i
			break
		}
	}
	idx = (idx + delta + len(order)) % len(order)
	return m.setFocus(order[idx])
}

func (m Model) focusOrder() []int {
	order := []int{-1, fieldB1, fieldB2, fieldBaseline}
	if m.method == rangemath.MethodFull {
		order = append(order, fieldBaselineBearing)
	}
	return order
}

func (m *Model) setFocus(field int) tea.Cmd {
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
	m.focus = field
	if field >= 0 && field < fieldCount {
		m.inputs[field].Focus()
		return textinput.Blink
	}
	return nil
}

func (m *Model) calculate() {
	b1, err := parseBearing(m.inputs[fieldB1].Value())
	if err != nil {
		m.flash = errStyle.Render("✗ bearing 1: " + err.Error())
		m.flashAt = time.Now()
		m.status = "error"
		return
	}
	b2, err := parseBearing(m.inputs[fieldB2].Value())
	if err != nil {
		m.flash = errStyle.Render("✗ bearing 2: " + err.Error())
		m.flashAt = time.Now()
		m.status = "error"
		return
	}
	baseline, err := parsePositive(m.inputs[fieldBaseline].Value(), "baseline")
	if err != nil {
		m.flash = errStyle.Render("✗ " + err.Error())
		m.flashAt = time.Now()
		m.status = "error"
		return
	}

	var res rangemath.Result
	if m.method == rangemath.MethodFull {
		bb, err := parseBearing(m.inputs[fieldBaselineBearing].Value())
		if err != nil {
			m.flash = errStyle.Render("✗ baseline bearing: " + err.Error())
			m.flashAt = time.Now()
			m.status = "error"
			return
		}
		res, err = rangemath.Full(b1, b2, baseline, bb)
		if err != nil {
			m.flash = errStyle.Render("✗ " + err.Error())
			m.flashAt = time.Now()
			m.status = "error"
			m.result = nil
			return
		}
	} else {
		res, err = rangemath.Simple(b1, b2, baseline)
		if err != nil {
			m.flash = errStyle.Render("✗ " + err.Error())
			m.flashAt = time.Now()
			m.status = "error"
			m.result = nil
			return
		}
	}

	m.result = &res
	m.status = "calculated"
	m.flash = upStyle.Render("✓ range calculated")
	m.flashAt = time.Now()

	summary := fmt.Sprintf("%.0f m  (θ %.1f° · %s)",
		(res.RangeFromRef1+res.RangeFromRef2)/2, res.AngleAtTarget, rangemath.MethodName(res.Method))
	detail := fmt.Sprintf("B1 %.1f°  B2 %.1f°  base %.0f m → R1 %.2f m  R2 %.2f m",
		b1, b2, baseline, res.RangeFromRef1, res.RangeFromRef2)
	m.history = append([]historyEntry{{summary: summary, detail: detail}}, m.history...)
	if len(m.history) > 20 {
		m.history = m.history[:20]
	}
}

func parseBearing(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("required")
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("not a number")
	}
	// Allow 0 and 360 as north.
	if v < 0 || v > 360 {
		return 0, fmt.Errorf("must be 0–360°")
	}
	return v, nil
}

func parsePositive(s, name string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("%s required", name)
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("%s not a number", name)
	}
	if v <= 0 {
		return 0, fmt.Errorf("%s must be > 0", name)
	}
	return v, nil
}

func (m *Model) layout() {
	rightW := m.width - leftWidth - 4
	if rightW < 36 {
		rightW = 36
	}
	h := m.height - 3
	if h < 8 {
		h = 8
	}
	m.detail.Width = max(20, rightW-2)
	m.detail.Height = max(5, h-2)
	for i := range m.inputs {
		m.inputs[i].Width = max(12, leftWidth-14)
	}
}

func (m *Model) refreshDetail() {
	m.detail.SetContent(m.renderDetail())
	m.detail.GotoTop()
}

func (m Model) View() string {
	if !m.ready {
		return "loading…"
	}

	header := titleStyle.Render(" RangeCalcCon ") + " " +
		dimStyle.Render("v"+m.version) + "  " +
		dimStyle.Render("range from two bearings")

	left := m.renderLeft()
	right := m.renderRight()
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right)

	help := helpStyle.Render(
		keyStyle.Render("tab") + "/" + keyStyle.Render("↑↓") + " field  " +
			keyStyle.Render("m") + " method  " +
			keyStyle.Render("enter") + "/" + keyStyle.Render("c") + " calc  " +
			keyStyle.Render("r") + " clear  " +
			keyStyle.Render("R") + " clear all  " +
			keyStyle.Render("q") + " quit",
	)
	flash := m.flash
	if flash != "" {
		flash = "  " + flash
	}
	status := dimStyle.Render(" " + m.status)
	return lipgloss.JoinVertical(lipgloss.Left, header, body, help+flash+status)
}

func (m Model) renderLeft() string {
	h := m.height - 3
	if h < 8 {
		h = 8
	}
	innerW := leftWidth - 2
	if innerW < 12 {
		innerW = 12
	}

	var b strings.Builder
	b.WriteString(panelTitleStyle.Render("Inputs") + "\n\n")

	// Method row
	methodLabel := "Method"
	simpleLab := " 1 Simple "
	fullLab := " 2 Full "
	if m.method == rangemath.MethodSimple {
		simpleLab = selectedStyle.Render(simpleLab)
		fullLab = normalItemStyle.Render(fullLab)
	} else {
		simpleLab = normalItemStyle.Render(simpleLab)
		fullLab = selectedStyle.Render(fullLab)
	}
	methodLine := dimStyle.Render(methodLabel) + "  " + simpleLab + " " + fullLab
	if m.focus < 0 {
		methodLine = selectedStyle.Render("› ") + methodLine
	} else {
		methodLine = "  " + methodLine
	}
	b.WriteString(methodLine + "\n\n")

	b.WriteString(m.fieldLine("Bearing Ref 1", fieldB1, "°") + "\n")
	b.WriteString(m.fieldLine("Bearing Ref 2", fieldB2, "°") + "\n")
	b.WriteString(m.fieldLine("Baseline dist", fieldBaseline, "m") + "\n")
	if m.method == rangemath.MethodFull {
		b.WriteString(m.fieldLine("Baseline brg", fieldBaselineBearing, "°") + "\n")
		b.WriteString("\n" + dimStyle.Render("Baseline bearing = direction") + "\n")
		b.WriteString(dimStyle.Render("from Ref 1 toward Ref 2.") + "\n")
	} else {
		b.WriteString("\n" + dimStyle.Render("Simple mode: approx. range") + "\n")
		b.WriteString(dimStyle.Render("using cot(θ). Press m for full.") + "\n")
	}

	_ = innerW
	return panelStyle.Width(leftWidth).Height(h).Render(b.String())
}

func (m Model) fieldLine(label string, field int, unit string) string {
	marker := "  "
	lab := dimStyle.Render(label)
	if m.focus == field {
		marker = selectedStyle.Render("› ")
		lab = boldStyle.Render(label)
	}
	return marker + lab + "\n  " + m.inputs[field].View() + " " + dimStyle.Render(unit)
}

func (m Model) renderRight() string {
	h := m.height - 3
	if h < 8 {
		h = 8
	}
	w := m.width - leftWidth - 4
	if w < 36 {
		w = 36
	}
	inner := panelTitleStyle.Render("Results") + "\n" + m.detail.View()
	return panelStyle.Width(w).Height(h).Render(inner)
}

func (m Model) renderDetail() string {
	var b strings.Builder
	dw := max(20, m.detail.Width)

	if m.result != nil {
		r := m.result
		b.WriteString(boldStyle.Render("Current calculation") + "\n")
		fmt.Fprintf(&b, "%s %s\n", dimStyle.Render("method"), rangemath.MethodName(r.Method))
		fmt.Fprintf(&b, "%s %s\n", dimStyle.Render("from Ref 1"), resultStyle.Render(fmt.Sprintf("%.2f m", r.RangeFromRef1)))
		fmt.Fprintf(&b, "%s %s\n", dimStyle.Render("from Ref 2"), resultStyle.Render(fmt.Sprintf("%.2f m", r.RangeFromRef2)))
		fmt.Fprintf(&b, "%s %.1f°\n", dimStyle.Render("angle at target"), r.AngleAtTarget)
		b.WriteByte('\n')
	} else {
		b.WriteString(dimStyle.Render("Enter bearings and baseline, then press Enter (or c) to calculate.") + "\n\n")
		b.WriteString(dimStyle.Render("Example (simple): 25° / 350° / 250 m → ~357 m") + "\n\n")
	}

	if len(m.history) > 0 {
		b.WriteString(panelTitleStyle.Render("History") + "\n")
		for i, h := range m.history {
			if i >= 12 {
				break
			}
			fmt.Fprintf(&b, "%s %s\n", dimStyle.Render(fmt.Sprintf("%2d.", i+1)), h.summary)
			// Wrap detail lightly
			d := h.detail
			if len(d) > dw {
				d = d[:dw-1] + "…"
			}
			b.WriteString(dimStyle.Render("    "+d) + "\n")
		}
	} else if m.result == nil {
		b.WriteString(dimStyle.Render("Previous calculations appear here.") + "\n")
	}

	return b.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
