# Tetris Game — Go + tview

A terminal-based Tetris game with an AI auto-player and a head-to-head AI Showdown mode. Built with Go and [tview](https://github.com/rivo/tview) (tcell v2).

## Features

### Core Gameplay
- Classic 10×20 Tetris board with standard 7 tetrominoes
- NES/Game Boy scoring formula (single/double/triple/tetris line clears)
- Level progression: speed increases every 10 lines cleared
- 7-bag randomizer for fair piece distribution
- Ghost piece showing where the current piece will land
- Background music (Korobeiniki) and sound effects

### Controls

| Key | Action |
|-----|--------|
| Left / Right Arrow | Move piece |
| Down Arrow | Soft drop |
| Space | Hard drop |
| Z | Rotate counter-clockwise |
| X | Rotate clockwise |
| G | Toggle ghost piece |
| A | Toggle AI auto-player |
| P | Pause / Resume |
| Q / Esc | Quit |

### AI Auto-Player
An on-board AI evaluates every possible placement using a weighted scoring function:

- **Aggregate height** — penalizes tall stacks
- **Holes** — penalizes covered empty cells
- **Bumpiness** — penalizes uneven column heights
- **Wells** — penalizes deep narrow gaps

The AI looks one piece ahead (`FindBestMoveWithNext`) and executes moves step-by-step each frame.

Toggle with the `A` key during a game.

### AI Showdown Mode
Watch two AI bots compete side-by-side on independent boards.

```bash
go run ./cmd/tetris --showdown
go run ./cmd/tetris --showdown --bot-a=aggressive --bot-b=conservative
```

**Available presets:**

| Preset | Style |
|--------|-------|
| `aggressive` | Plays fast, accepts risk |
| `conservative` | Prioritizes a clean, flat board |
| `balanced` | Default balanced weights |
| `speedrun` | Speed-focused, moderate safety |
| `chaos` | Near-random — nearly ignores heuristics |

**Showdown controls:**

| Key | Action |
|-----|--------|
| + / - | Increase / decrease speed (1–5) |
| R | Restart both boards |
| Q | Quit |

Layout: `Board A (22) | Stats (26) | Board B (22)` — 70 columns total.

## Prerequisites

- Go 1.21 or newer
- A terminal with 256-color support (`TERM=xterm-256color` or equivalent)
- macOS / Linux (or Windows Terminal with ANSI support)

## Installation & Running

```bash
# Install dependencies
go mod tidy

# Run from source
go run ./cmd/tetris

# Build binary
go build -o tetris ./cmd/tetris
./tetris

# Run showdown
go run ./cmd/tetris --showdown --bot-a=aggressive --bot-b=conservative
```

## Running Tests

```bash
go test ./...          # All tests
go test -race ./...    # With race detector
go fmt ./...           # Format
go vet ./...           # Lint
```

## Architecture

```
cmd/tetris/          — entry point, tview layout, rendering, key handling
internal/model/      — all game logic (no UI dependencies)
  board.go           — 10×20 grid, (0,0) at bottom-left
  gamestate.go       — orchestrates game flow
  piece.go           — tetromino shapes and rotation matrices
  randomizer.go      — 7-bag randomizer
  autoplay.go        — AI player with configurable weights
  presets.go         — named weight presets
  showdown.go        — two-board showdown state machine
internal/audio/      — background music and sound effects (gopxl/beep)
internal/testutil/   — shared test helpers
```
