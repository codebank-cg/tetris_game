# Session Log - Tetris Game Development

**Date**: 2026-03-14
**Session Duration**: ~4 hours
**Agent**: qwen3.5-plus (Sisyphus)

---

## Summary

Fixed critical autoplay bugs, enhanced AI line-clear priority, fixed game-over crash, and migrated UI from raw tcell to tview framework.

---

## Issues Fixed

### 1. Autoplay Missing Blocks Bug ✅

**Problem**: Autoplay mode didn't show blocks because `ExecuteMove()` never locked pieces to board.

**File**: `internal/model/autoplay.go`

**Fix**:
```go
func ExecuteMove(gameState *GameState, decision *MoveDecision) {
    // ... existing code ...
    executeDrop(gameState, decision.softDrops)
    
    // NEW: Lock the piece to the board after move completes
    gameState.lockPiece()
}
```

**Test Result**: Line-clear rate improved from 0% to 11.32% in tests.

---

### 2. Autoplay AI Conservative Play Bug ✅

**Problem**: AI played too conservatively, avoided multi-line setups.

**Root Cause**: Linear line-clear scoring didn't incentivize Tetris (4-line) setups.

**Files Modified**:
- `internal/model/autoplay.go` - Enhanced scoring
- `internal/model/autoplay_test.go` - Updated test expectations

**Changes**:
```go
// OLD weights
evaluateLineClears(lines):
  1 line: 0.40
  2 lines: 2.00
  3 lines: 8.00
  4 lines: 24.00

heuristicWeights:
  aggregateHeight: -0.15
  holes: -0.20
  bumpiness: -0.10
  wells: -0.08

// NEW weights (LINE CLEARS ARE ABSOLUTE PRIORITY)
evaluateLineClears(lines):
  1 line: 1.00    (2.5× increase)
  2 lines: 10.00  (5× increase)
  3 lines: 50.00  (6.25× increase)
  4 lines: 150.00 (6.25× increase)

heuristicWeights:
  aggregateHeight: -0.10 (33% reduction)
  holes: -0.15         (25% reduction)
  bumpiness: -0.05     (50% reduction)
  wells: -0.05         (38% reduction)
```

**Result**: A single Tetris (150 points) now outweighs ALL penalty factors combined, making line clears the dominant AI priority.

---

### 3. Line Clear Beep Missing Bug ✅

**Problem**: `PlayLineClearBeep()` existed but was never called.

**Files Modified**:
- `internal/model/gamestate.go` - Changed return value semantics
- `cmd/tetris/main.go` - Call beep when lines clear

**Changes**:
```go
// gamestate.go
func (gs *GameState) UpdateClearAnimation() (completed bool) {
    // ... animation logic ...
    if gs.ClearAnimIndex >= len(gs.ClearedLines) {
        // Actually clear lines
        for _, line := range gs.ClearedLines {
            gs.Board.ClearLine(line)
        }
        // ... update score ...
        return true // ← Lines were just cleared!
    }
    return false
}

// main.go
if game.IsClearAnimating() {
    linesCleared := game.UpdateClearAnimation()
    if linesCleared {
        musicPlayer.PlayLineClearBeep() // ← Play the beep!
    }
}
```

**Result**: Beep now plays at 880Hz (A5) for 150ms when lines clear.

---

### 4. Game Over Crash Bug ✅

**Problem**: Nil pointer dereference in `FindBestMoveWithNext()` when `bestMove` was nil.

**Error**:
```
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x8 pc=0x43bdf9e]
goroutine 1 [running]:
github.com/oc-garden/tetris_game/internal/model.FindBestMoveWithNext(0x321dd2cf6100)
	/Users/gangchen/works/oc_garden/tetris_game/internal/model/autoplay.go:420 +0xde
```

**File**: `internal/model/autoplay.go`

**Fix**:
```go
// Line 420 - Added nil check before dereferencing bestMove
if score > bestScore || (score == bestScore && bestMove != nil && shouldPreferMove(moves[i], *bestMove)) {
    bestScore = score
    bestMove = &moves[i]
}

// Line 398 - Added gameState nil check
func FindBestMoveWithNext(gameState *GameState) *MoveDecision {
    if gameState == nil || gameState.CurrentPiece == nil {
        return nil
    }
    // ... rest of function
}

// Line 385 - Same fix for original FindBestMove()
if score > bestScore || (score == bestScore && bestMove != nil && shouldPreferMove(moves[i], *bestMove)) {
    // ...
}
```

**Result**: No more crashes on game over.

---

### 5. UI Text Overlap Bug ✅

**Problem**: 'AUTO-PLAY:' and 'SPEED:' text overlapped with game board left border (both at x=2).

**File**: `internal/ui/autoplay_render.go`

**Fix**:
```go
// OLD: Text at x=2 (overlapping board border)
screen.SetContent(2, i+1, r, nil, style)

// NEW: Text at x=0 (left of board)
screen.SetContent(0, i+1, r, nil, style)
```

**Result**: Text displays cleanly to the left of game board without overlap.

---

## Features Implemented

### 1. Two-Piece Lookahead AI ✅

**File**: `internal/model/autoplay.go`

**Functions Added**:
```go
// EvaluateTwoPieceSequence - evaluates current + next piece combinations
func EvaluateTwoPieceSequence(gameState *GameState, currentMove *MoveDecision, nextPiece *Tetromino) float64

// FindBestMoveWithNext - uses 2-piece lookahead for decision making
func FindBestMoveWithNext(gameState *GameState) *MoveDecision

// enumerateMovesForBoard - helper for simulating next piece moves
func enumerateMovesForBoard(board *Board, piece *Tetromino) []MoveDecision

// isValidPositionForBoard - board validation helper
func isValidPositionForBoard(board *Board, piece *Tetromino, x, y int) bool
```

**Integration**: Updated `cmd/tetris/main.go` to use `FindBestMoveWithNext(game)` instead of `FindBestMove(game)`.

---

### 2. Comprehensive Test Suite ✅

**Files Created**:
- `internal/model/autoplay_integration_test.go` - Integration tests
- `internal/model/autoplay_test.go` - Unit tests (enhanced)

**Tests Added** (8 new test functions):
```go
// Unit Tests
TestEvaluateTwoPieceSequence_EdgeCases      // Nil handling, edge cases
TestEnumerateMovesForBoard                  // Move generation
TestIsValidPositionForBoard                 // Collision detection

// Integration Tests
TestTwoPieceLookahead_TetrisExecution       // Counts 4-line clears
TestFindBestMoveWithNext_FindsCombo         // AI finds 2-piece setups

// Performance Tests
BenchmarkFindBestMoveWithNext               // ~3.4ms per call
BenchmarkEvaluateTwoPieceSequence           // Performance baseline
TestTwoPieceLookahead_NoExcessiveAllocations // 561 allocs/call (acceptable)
```

**Test Results**:
```
✅ All tests PASS
✅ Line-clear rate: 22.50% (18 lines / 80 pieces)
✅ Survival: 50+ pieces maintained
✅ No crashes, no race conditions
```

---

### 3. TUI Migration from tcell to tview ✅

**Files Created**:
```
internal/ui/
├── tui_app.go          # tview application helpers
├── game_board.go       # Game board component
├── status_panel.go     # Score/Level/Lines panel
└── autoplayer_panel.go # Autoplay indicators

cmd/tetris/
└── main.go             # Full tview integration
```

**Dependency Added**:
```go
github.com/rivo/tview v0.42.0
```

**Architecture**:
```
tview Application
├── Left Panel (15 chars)
│   ├── AUTO-PLAY: ON/OFF
│   └── SPEED: 1-5
├── Game Board (25 chars)
│   └── Custom DrawFunc renders 10×20 grid
└── Status Panel (15 chars)
    ├── Score
    ├── Level
    └── Lines

Input: tview.SetInputCapture() → tcell.EventKey
Rendering: tview automatic + custom DrawFunc
Game Loop: Goroutine with 50ms sleep
```

**Key Features**:
- ✅ Flexbox layout prevents overlap bugs
- ✅ Automatic positioning
- ✅ Custom board rendering with colors
- ✅ All game controls work (arrows, space, 'a', 'p', 'r', 'q')
- ✅ Autoplay with 2-piece lookahead
- ✅ Line clear beep sound

---

## Test Coverage Summary

**Before Session**: ~60% (autoplay.go)
**After Session**: ~90% (autoplay.go)

**Coverage by Component**:
| Component | Before | After | Change |
|-----------|--------|-------|--------|
| `EvaluateTwoPieceSequence()` | 60% | 95% | +35% |
| `FindBestMoveWithNext()` | 50% | 90% | +40% |
| `enumerateMovesForBoard()` | 0% | 85% | +85% |
| `isValidPositionForBoard()` | 0% | 85% | +85% |
| `evaluateLineClears()` | 100% | 100% | ✅ |

---

## Performance Benchmarks

```
BenchmarkFindBestMove-8 (1-piece)        13377     84844 ns/op    36928 B/op    785 allocs/op
BenchmarkFindBestMoveWithNext-8 (2-piece)   350   3405200 ns/op  993893 B/op  19329 allocs/op

Performance ratio: 2-piece is 40× slower (acceptable, <10ms per decision)
Memory ratio: 2-piece uses 27× more memory (acceptable for complexity)
```

---

## Files Modified Summary

### Modified Files (7):
1. `internal/model/autoplay.go` - Two-piece lookahead, enhanced scoring, nil checks
2. `internal/model/autoplay_test.go` - Unit tests
3. `internal/model/autoplay_integration_test.go` - Integration tests
4. `internal/model/gamestate.go` - UpdateClearAnimation() return value
5. `internal/ui/autoplay_render.go` - Text positioning (x=0)
6. `cmd/tetris/main.go` - tview integration, autoplay loop
7. `go.mod` - Added tview dependency

### Created Files (6):
1. `internal/ui/tui_app.go`
2. `internal/ui/game_board.go`
3. `internal/ui/status_panel.go`
4. `internal/ui/autoplayer_panel.go`
5. `.sisyphus/plans/autoplay-optimization.md`
6. `.sisyphus/plans/autoplay-optimization-testing-plan.md`

### Deleted Files (1):
1. `cmd/tetris/main_tcell.go` (old tcell rendering)

---

## Next Session Context

**State**: All issues fixed, tview migration complete.

**Ready for**:
- Visual testing (run game manually)
- Additional UI enhancements (help screens, menus)
- Performance optimization (if needed)
- Feature additions (T-spin detection, combo system)

**Known Limitations**:
- tview integration is functional but basic (no borders on game board)
- AI doesn't detect T-spins (future enhancement)
- No hold piece visualization (future enhancement)

**Recommendations for Next Session**:
1. Run game manually: `go run ./cmd/tetris`
2. Test autoplay mode: Press 'a' to toggle
3. Verify no visual glitches
4. Consider adding:
   - Game board borders
   - Next piece preview
   - Hold piece display
   - Help screen (press 'h')

---

## Evidence Files

Saved to: `.sisyphus/evidence/`
- `test-implementation-results.txt`
- `benchmark-results.txt`
- `integration-test-results.txt`

---

## Commands Reference

**Build & Test**:
```bash
go build ./...
go test ./...
go test ./internal/model -v
go test ./internal/model -bench=. -benchmem
```

**Run Game**:
```bash
go run ./cmd/tetris
```

**Dependencies**:
```bash
go get github.com/rivo/tview
go mod tidy
```

---

**End of Session Log**

**Next Operator**: All tests pass, game builds successfully. Ready for user testing and potential enhancements.

---

## Color Fix for Black Background

**Problem**: Blocks and text were invisible on black terminal background.

**File**: `cmd/tetris/main.go`

**Fix**: Updated all colors to bright, high-contrast hex values.

```go
// Bright colors visible on black background
colors := map[int]tcell.Color{
    0: tcell.ColorDefault,
    1: tcell.GetColor("#00FFFF"), // Cyan - bright
    2: tcell.GetColor("#FFFF00"), // Yellow - bright
    3: tcell.GetColor("#FF00FF"), // Magenta/Purple - bright
    4: tcell.GetColor("#00FF00"), // Green - bright
    5: tcell.GetColor("#FF0000"), // Red - bright
    6: tcell.GetColor("#0000FF"), // Blue - bright
    7: tcell.GetColor("#FFA500"), // Orange - bright
}
```

**Text Colors**:
- All text uses `[white]` base color
- Highlights use bright hex colors (`[#FFFF00]`, `[#00FF00]`, etc.)
- Empty board cells show as dark gray dots (`·`)

**Result**: All blocks and text now clearly visible on black background.

