# RangeCalcCon

Terminal UI for estimating **range to a target** from two compass bearings and a known baseline.

**Left panel** — method + input fields  
**Right panel** — current result and calculation history

Console style matches [GetWattPad](../../GetWattPad) / [ServiceMonitor](../../ServiceMonitor): **bubbletea** TUI, **lipgloss** panels, keyboard-driven workflow.

Math matches the fixed triangulation used across RangeCalc (circular bearing wrap, degrees→radians, law of sines for full mode).

## Requirements

- Go 1.26+
- **macOS / Linux / Windows**

## Quick start

### Windows

```powershell
cd RangeCalcCon
.\scripts\build-windows.ps1   # builds rangecalccon.exe
.\rangecalccon.exe
```

### macOS / Linux

```bash
cd RangeCalcCon
go build -ldflags "-s -w -X main.version=$(cat VERSION)" -o rangecalccon .
./rangecalccon
```

### Keyboard

| Key | Action |
|-----|--------|
| `tab` / `↑` `↓` / `j` `k` | Move between method + fields |
| `m` | Toggle Simple ↔ Full triangulation |
| `←` `→` | Simple / Full (when method row focused) |
| `enter` / `c` | Calculate |
| `r` | Clear form (keep history) |
| `R` | Clear form + history |
| `PgUp` / `PgDn` | Scroll results |
| `q` / `ctrl+c` | Quit |

### Methods

1. **Simple** — `range = baseline × cot(θ)` where `θ` is the smallest angle between the two bearings. Fast; assumes favourable geometry. Needs `θ < 90°`.
2. **Full triangulation** — law of sines. Also needs **baseline bearing** (compass direction from Ref 1 to Ref 2). Reports separate ranges from each reference.

### Example

- Bearing Ref 1: **25°**
- Bearing Ref 2: **350°**
- Baseline: **250 m**
- Simple result: **≈ 357.04 m** (angle at target 35°)

## Layout

```
RangeCalcCon/
  rangecalccon / .exe        # built TUI
  scripts/build-windows.ps1  # Windows build
  VERSION
  internal/
    rangemath/               # simple + full triangulation (+ tests)
    ui/                      # bubbletea + lipgloss TUI
```

## License

Same as the parent RangeCalc project.
