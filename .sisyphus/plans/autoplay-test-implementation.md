# Autoplay Optimization - Test Implementation Plan

## Status: READY FOR EXECUTION

This document provides the exact implementation steps for all TODO test cases. Execute with `/start-work autoplay-test-implementation`.

---

## TODO Tests to Implement (9 total)

### Priority 1: Unit Tests (3 tests)

#### Test 1: `TestEvaluateTwoPieceSequence_EdgeCases`
**File**: `internal/model/autoplay_test.go`
**Location**: After line 682 (after `TestEvaluateTwoPieceSequence_ComboBonus`)

**Implementation**:
```go
func TestEvaluateTwoPieceSequence_EdgeCases(t *testing.T) {
	gameState := NewGameState()
	decision := &MoveDecision{rotations: 0, targetX: 5, softDrops: 10}
	nextPiece := NewTetromino(TetrominoI)

	t.Run("NilGameState", func(t *testing.T) {
		score := EvaluateTwoPieceSequence(nil, decision, nextPiece)
		if score != -999999.0 {
			t.Errorf("Expected -999999.0 for nil gameState, got %.2f", score)
		}
	})

	t.Run("NilDecision", func(t *testing.T) {
		score := EvaluateTwoPieceSequence(gameState, nil, nextPiece)
		if score != -999999.0 {
			t.Errorf("Expected -999999.0 for nil decision, got %.2f", score)
		}
	})

	t.Run("NilNextPiece", func(t *testing.T) {
		score := EvaluateTwoPieceSequence(gameState, decision, nil)
		if score != -999999.0 {
			t.Errorf("Expected -999999.0 for nil nextPiece, got %.2f", score)
		}
	})

	t.Run("NoValidNextMoves", func(t *testing.T) {
		// Fill board completely except 1 row
		for y := 0; y < 19; y++ {
			for x := 0; x < 10; x++ {
				gameState.Board.Set(x, y, 1)
			}
		}
		score := EvaluateTwoPieceSequence(gameState, decision, nextPiece)
		if score != -999999.0 {
			t.Errorf("Expected -999999.0 for no valid moves, got %.2f", score)
		}
	})
}
```

**Acceptance Criteria**:
- [ ] All 4 sub-tests pass
- [ ] Nil parameters return -999999.0
- [ ] No valid moves returns -999999.0
- [ ] Test runs in <10ms

---

#### Test 2: `TestEnumerateMovesForBoard`
**File**: `internal/model/autoplay_test.go`
**Location**: After Test 1

**Implementation**:
```go
func TestEnumerateMovesForBoard(t *testing.T) {
	board := NewBoard()

	t.Run("EmptyBoard_AllPieces", func(t *testing.T) {
		pieceTypes := []TetrominoType{TetrominoI, TetrominoO, TetrominoT, TetrominoS, TetrominoZ, TetrominoJ, TetrominoL}
		for _, pt := range pieceTypes {
			piece := NewTetromino(pt)
			moves := enumerateMovesForBoard(board, piece)
			if len(moves) == 0 {
				t.Errorf("Expected moves for %v on empty board, got 0", pt)
			}
		}
	})

	t.Run("NarrowGap_OnlyIPieceFits", func(t *testing.T) {
		// Create 1-wide well at column 5
		for x := 0; x < 10; x++ {
			for y := 0; y < 5; y++ {
				if x != 5 {
					board.Set(x, y, 1)
				}
			}
		}

		iPiece := NewTetromino(TetrominoI)
		iMoves := enumerateMovesForBoard(board, iPiece)
		if len(iMoves) == 0 {
			t.Error("Expected I-piece moves in 1-wide well")
		}

		oPiece := NewTetromino(TetrominoO)
		oMoves := enumerateMovesForBoard(board, oPiece)
		if len(oMoves) > 0 {
			t.Errorf("Expected 0 O-piece moves in 1-wide well, got %d", len(oMoves))
		}
	})

	t.Run("FullBoard_NoMoves", func(t *testing.T) {
		fullBoard := NewBoard()
		for y := 0; y < 20; y++ {
			for x := 0; x < 10; x++ {
				fullBoard.Set(x, y, 1)
			}
		}

		piece := NewTetromino(TetrominoI)
		moves := enumerateMovesForBoard(fullBoard, piece)
		if len(moves) != 0 {
			t.Errorf("Expected 0 moves on full board, got %d", len(moves))
		}
	})
}
```

**Acceptance Criteria**:
- [ ] Empty board: all 7 piece types return moves
- [ ] 1-wide well: I-piece has moves, O-piece has 0 moves
- [ ] Full board: all pieces return 0 moves
- [ ] Test runs in <50ms

---

#### Test 3: `TestIsValidPositionForBoard`
**File**: `internal/model/autoplay_test.go`
**Location**: After Test 2

**Implementation**:
```go
func TestIsValidPositionForBoard(t *testing.T) {
	board := NewBoard()
	piece := NewTetromino(TetrominoI)

	t.Run("EmptyBoard_Valid", func(t *testing.T) {
		if !isValidPositionForBoard(board, piece, 5, 18) {
			t.Error("Expected valid position on empty board")
		}
	})

	t.Run("OutOfBounds_Invalid", func(t *testing.T) {
		if isValidPositionForBoard(board, piece, 10, 18) {
			t.Error("Expected invalid position (x=10 out of bounds)")
		}
		if isValidPositionForBoard(board, piece, -1, 18) {
			t.Error("Expected invalid position (x=-1 out of bounds)")
		}
	})

	t.Run("Collision_Invalid", func(t *testing.T) {
		// Place block where piece would land
		board.Set(5, 16, 1)
		if isValidPositionForBoard(board, piece, 5, 17) {
			t.Error("Expected invalid position (collision at y=16)")
		}
	})

	t.Run("ClearPosition_Valid", func(t *testing.T) {
		// Block far away, position should still be valid
		board.Set(0, 0, 1)
		if !isValidPositionForBoard(board, piece, 9, 18) {
			t.Error("Expected valid position (block at 0,0 doesn't interfere)")
		}
	})

	t.Run("RotationCollision_Invalid", func(t *testing.T) {
		// Clear board first
		board = NewBoard()
		// Block adjacent to piece rotation
		board.Set(6, 18, 1)
		// Rotate I-piece (it will extend to x=6)
		rotatedPiece := NewTetromino(TetrominoI)
		rotatedPiece.RotateClockwise()
		if isValidPositionForBoard(board, rotatedPiece, 5, 18) {
			t.Error("Expected invalid position (rotation collision)")
		}
	})
}
```

**Acceptance Criteria**:
- [ ] All 5 sub-tests pass
- [ ] Empty board returns true
- [ ] Out of bounds returns false
- [ ] Collision returns false
- [ ] Clear position returns true
- [ ] Rotation collision returns false

---

### Priority 2: Integration Tests (2 tests)

#### Test 4: `TestTwoPieceLookahead_TetrisExecution`
**File**: `internal/model/autoplay_integration_test.go`
**Location**: After line 285 (after `BenchmarkAutoPlayGame`)

**Implementation**:
```go
// TestTwoPieceLookahead_TetrisExecution verifies AI executes 4-line clears (Tetris)
// Target: At least 1 Tetris in 100 pieces with 2-piece lookahead
func TestTwoPieceLookahead_TetrisExecution(t *testing.T) {
	gameState := NewGameState()
	tetrisCount := 0
	piecesPlaced := 0

	for piecesPlaced < 100 && !gameState.GameOver {
		linesBefore := gameState.LinesCleared
		decision := FindBestMoveWithNext(gameState)
		if decision == nil {
			break
		}
		ExecuteMove(gameState, decision)

		// Simulate game loop: process line clear animation
		for gameState.IsClearAnimating() {
			gameState.UpdateClearAnimation()
		}

		linesCleared := gameState.LinesCleared - linesBefore
		if linesCleared == 4 {
			tetrisCount++
		}
		piecesPlaced++
	}

	t.Logf("Executed %d Tetrises in %d pieces (%.2f%%)",
		tetrisCount, piecesPlaced, float64(tetrisCount)/float64(piecesPlaced)*100)

	// Target: At least 1 Tetris in 100 pieces
	if tetrisCount < 1 {
		t.Errorf("Expected ≥1 Tetris, got %d in %d pieces", tetrisCount, piecesPlaced)
	}

	// Verify survival
	if piecesPlaced < 50 {
		t.Errorf("Game ended too early: only %d pieces", piecesPlaced)
	}
}
```

**Acceptance Criteria**:
- [ ] At least 1 Tetris executed in 100 pieces
- [ ] Survival: 50+ pieces
- [ ] Test logs Tetris frequency percentage
- [ ] Test runs in <2 seconds

---

#### Test 5: `TestFindBestMoveWithNext_FindsCombo`
**File**: `internal/model/autoplay_integration_test.go`
**Location**: After Test 4

**Implementation**:
```go
// TestFindBestMoveWithNext_FindsCombo verifies AI finds 2-piece Tetris setups
func TestFindBestMoveWithNext_FindsCombo(t *testing.T) {
	gameState := NewGameState()

	// Create classic Tetris setup: 3 rows with 1-column well at column 4
	for y := 0; y < 3; y++ {
		for x := 0; x < 10; x++ {
			if x != 4 { // Leave column 4 open for I-piece
				gameState.Board.Set(x, y, 1)
			}
		}
	}

	// O-piece current (should build around well, not fill it)
	gameState.CurrentPiece = NewTetromino(TetrominoO)
	gameState.CurrentPiece.X = 3
	gameState.CurrentPiece.Y = 18

	// I-piece next (will clear Tetris)
	gameState.NextPiece = NewTetromino(TetrominoI)

	decision := FindBestMoveWithNext(gameState)
	if decision == nil {
		t.Fatal("FindBestMoveWithNext() returned nil")
	}

	t.Logf("AI decision: X=%d, rotations=%d, drops=%d",
		decision.GetTargetX(), decision.GetRotations(), decision.GetSoftDrops())

	// Verify AI doesn't fill the well (column 4)
	if decision.GetTargetX() == 4 {
		t.Error("AI filled the Tetris well! Should preserve column 4 for I-piece")
	}

	// Expected: AI places O-piece at x=2 or x=6 (building around well)
	validPositions := []int{2, 3, 5, 6} // Positions that build around well
	isValid := false
	for _, pos := range validPositions {
		if decision.GetTargetX() == pos {
			isValid = true
			break
		}
	}

	if !isValid {
		t.Errorf("Expected O-piece at %v, got X=%d", validPositions, decision.GetTargetX())
	}
}
```

**Acceptance Criteria**:
- [ ] AI doesn't fill column 4 (the well)
- [ ] AI places O-piece at position 2, 3, 5, or 6 (around well)
- [ ] Test logs AI decision details
- [ ] Test runs in <100ms

---

### Priority 3: Performance Tests (3 tests)

#### Test 6: `BenchmarkFindBestMoveWithNext`
**File**: `internal/model/autoplay_test.go`
**Location**: After existing benchmarks

**Implementation**:
```go
func BenchmarkFindBestMoveWithNext(b *testing.B) {
	gameState := NewGameState()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindBestMoveWithNext(gameState)
	}
}
```

**Target Metrics**:
- <10ms per call (acceptable for real-time play)
- <500KB allocs/op
- <100 allocs/op

**Comparison Baseline**:
```bash
# Run both benchmarks for comparison
go test -bench="BenchmarkFindBestMove$" ./internal/model -benchmem
go test -bench="BenchmarkFindBestMoveWithNext" ./internal/model -benchmem
```

Expected: 2-piece is 5-10× slower than 1-piece (acceptable trade-off)

---

#### Test 7: `BenchmarkEvaluateTwoPieceSequence`
**File**: `internal/model/autoplay_test.go`
**Location**: After Test 6

**Implementation**:
```go
func BenchmarkEvaluateTwoPieceSequence(b *testing.B) {
	gameState := NewGameState()
	decision := &MoveDecision{
		rotations: 0,
		targetX:   5,
		softDrops: 10,
	}
	nextPiece := NewTetromino(TetrominoI)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvaluateTwoPieceSequence(gameState, decision, nextPiece)
	}
}
```

**Target Metrics**:
- <1ms per evaluation
- <50KB allocs/op
- <20 allocs/op

---

#### Test 8: `TestTwoPieceLookahead_NoExcessiveAllocations`
**File**: `internal/model/autoplay_test.go`
**Location**: After Test 7

**Implementation**:
```go
func TestTwoPieceLookahead_NoExcessiveAllocations(t *testing.T) {
	gameState := NewGameState()
	decision := &MoveDecision{
		rotations: 0,
		targetX:   5,
		softDrops: 10,
	}

	allocs := testing.AllocsPerRun(100, func() {
		EvaluateTwoPieceSequence(gameState, decision, gameState.NextPiece)
	})

	// Should allocate ~2-3 boards for simulation (each board is ~200 bytes)
	// Allow up to 10 allocations for slice operations, maps, etc.
	if allocs > 10 {
		t.Errorf("Excessive allocations: %.1f per call (expected <10)", allocs)
	}
}
```

**Acceptance Criteria**:
- [ ] <10 allocations per call
- [ ] Board clones are efficient (not copying unnecessarily)
- [ ] Test runs in <500ms

---

### Priority 4: Regression Test (1 test)

#### Test 9: Already Implemented ✅

All regression tests are already covered:
- `TestFindBestMove_StillWorks` → `TestFindBestMove` (existing)
- `TestWeights_GetWeights` → `TestGetWeights_SetWeights` (existing, updated)
- `TestExecuteMove_WithTwoPiecePlanning` → `TestAutoPlay_Survival10Pieces` (existing)

---

## Execution Plan

### Step 1: Add Unit Tests (10 minutes)
```bash
# Edit: internal/model/autoplay_test.go
# Add: Tests 1, 2, 3 (edge cases, enumeration, validation)
go test ./internal/model -run "TestEvaluateTwoPieceSequence_EdgeCases|TestEnumerateMovesForBoard|TestIsValidPositionForBoard" -v
```

### Step 2: Add Integration Tests (10 minutes)
```bash
# Edit: internal/model/autoplay_integration_test.go
# Add: Tests 4, 5 (Tetris execution, combo finding)
go test ./internal/model -run "TestTwoPieceLookahead_TetrisExecution|TestFindBestMoveWithNext_FindsCombo" -v
```

### Step 3: Add Performance Tests (5 minutes)
```bash
# Edit: internal/model/autoplay_test.go
# Add: Tests 6, 7, 8 (benchmarks, allocation test)
go test ./internal/model -bench="BenchmarkFindBestMoveWithNext|BenchmarkEvaluateTwoPiece" -benchmem
go test ./internal/model -run "TestTwoPieceLookahead_NoExcessiveAllocations" -v
```

### Step 4: Run Full Test Suite (5 minutes)
```bash
go test ./internal/model/... -v
go test ./internal/model -race
go build ./...
go vet ./...
```

### Step 5: Capture Evidence (5 minutes)
```bash
# Save test output to evidence files
go test ./internal/model -v > .sisyphus/evidence/task-5-unit-tests.txt
go test ./internal/model -run "TestTwoPieceLookahead" -v > .sisyphus/evidence/task-8-integration-tests.txt
go test ./internal/model -bench="BenchmarkFindBest" -benchmem > .sisyphus/evidence/task-9-benchmarks.txt
```

---

## Test Implementation Checklist

- [ ] Test 1: `TestEvaluateTwoPieceSequence_EdgeCases` - Unit test
- [ ] Test 2: `TestEnumerateMovesForBoard` - Unit test
- [ ] Test 3: `TestIsValidPositionForBoard` - Unit test
- [ ] Test 4: `TestTwoPieceLookahead_TetrisExecution` - Integration test
- [ ] Test 5: `TestFindBestMoveWithNext_FindsCombo` - Integration test
- [ ] Test 6: `BenchmarkFindBestMoveWithNext` - Performance test
- [ ] Test 7: `BenchmarkEvaluateTwoPieceSequence` - Performance test
- [ ] Test 8: `TestTwoPieceLookahead_NoExcessiveAllocations` - Performance test

**Total**: 8 tests to implement (1 regression test already covered)

---

## Expected Test Coverage After Implementation

| Component | Before | After |
|-----------|--------|-------|
| `EvaluateTwoPieceSequence()` | 60% | **95%** |
| `FindBestMoveWithNext()` | 50% | **90%** |
| `enumerateMovesForBoard()` | 0% | **85%** |
| `isValidPositionForBoard()` | 0% | **85%** |
| `evaluateLineClears()` | 100% | 100% ✅ |

**Overall Target**: 85%+ coverage for autoplay.go

---

## Success Criteria

All tests must pass:
```bash
go test ./internal/model -v 2>&1 | grep -E "(PASS|FAIL)" | tail -5
```

Expected output:
```
--- PASS: TestEvaluateTwoPieceSequence_EdgeCases
--- PASS: TestEnumerateMovesForBoard
--- PASS: TestIsValidPositionForBoard
--- PASS: TestTwoPieceLookahead_TetrisExecution
--- PASS: TestFindBestMoveWithNext_FindsCombo
PASS
ok  	github.com/oc-garden/tetris_game/internal/model
```

---

## Files to Modify

1. **`internal/model/autoplay_test.go`**
   - Add Tests 1, 2, 3, 6, 7, 8
   - Estimated lines: +150

2. **`internal/model/autoplay_integration_test.go`**
   - Add Tests 4, 5
   - Estimated lines: +80

**Total changes**: ~230 lines of test code

---

## Time Estimate

- **Implementation**: 25 minutes
- **Testing**: 10 minutes
- **Evidence capture**: 5 minutes
- **Total**: 40 minutes

---

## Next Step

To execute this plan:
```bash
/start-work autoplay-test-implementation
```

This will implement all 8 TODO test cases, run the full test suite, and capture evidence.

---

**Created**: 2026-03-14
**Author**: Prometheus (Planning Agent)
**Status**: READY FOR EXECUTION
