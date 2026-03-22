# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build & Run
go build -o tetris ./cmd/tetris    # Build binary
go run ./cmd/tetris                # Run from source

# Test
go test ./...                      # All tests
go test -v -run TestBoardIsEmpty ./internal/model/  # Single test
go test -race ./...                # Race detector

# Format & Lint
go fmt ./...
go vet ./...
```

## Architecture

**Entry point**: `cmd/tetris/main.go` — game loop, input handling, UI rendering using tview.

**`internal/model/`** — all game logic (no UI dependencies):
- `gamestate.go` — `GameState` struct orchestrates the game; methods: `MovePiece`, `RotatePiece`, `SoftDrop`, `DropPiece`, `UpdateScore`
- `board.go` — 10×20 grid; coordinate system is (0,0) at bottom-left, Y increases upward
- `piece.go` — `Tetromino` with 4×4 rotation matrices; `RotateClockwise`/`RotateCounterClockwise`
- `randomizer.go` — 7-bag randomizer for fair piece distribution
- `autoplay.go` — AI player; `FindBestMoveWithNext` evaluates placements, `ExecuteMove` steps through moves

**`internal/audio/`** — `music.go` plays background music (Korobeiniki) and sound effects via `gopxl/beep`.

**`internal/testutil/`** — shared test helpers.

**Key design facts**:
- UI uses `tview` (built on `tcell v2`) with a fixed 3-column layout: auto-play panel | game board | next/info panel
- Rendering is done via custom `SetDrawFunc` callbacks on tview boxes — each cell is 2 characters wide (`██`)
- Ghost piece uses `░` characters; rendered only when `ghostEnabled && !autoPlayer.IsEnabled()`
- Game loop runs in a goroutine at ~60fps (16ms sleep); calls `app.Draw()` only when state changes
- Scoring: original Nintendo NES/Game Boy formula; level-up every 10 lines
- Drop interval: `1500 - (level-1)*100` ms, minimum 100ms

## Conventions

- All new game logic goes in `internal/model/`, not in `main.go`
- Constructor pattern: `NewBoard()`, `NewTetromino()`, `NewGameState()`
- Pointer receivers on all methods for consistency
- Table-driven tests; test function naming: `Test[Functionality][Scenario]`
- Run `go test ./...` after any code change before considering work done
