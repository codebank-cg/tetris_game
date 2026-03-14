# Autoplay Optimization - Comprehensive Testing Plan

## Overview

This document provides detailed test cases and testing strategy for the autoplay optimization featuring two-piece lookahead and enhanced multi-line priority.

---

## Test Categories

### Category 1: Unit Tests (Function-Level)

**Purpose**: Verify individual functions work correctly in isolation.

### Category 2: Integration Tests (AI Behavior)

**Purpose**: Verify AI makes correct decisions in realistic game scenarios.

### Category 3: Performance Tests

**Purpose**: Ensure algorithm meets performance requirements.

### Category 4: Regression Tests

**Purpose**: Ensure new features don't break existing functionality.

---

## Test Case Details

### **UNIT TESTS**

#### **Test 1.1: evaluateLineClears_EnhancedRatios**

**File**: `autoplay_test.go`

**Purpose**: Verify enhanced multi-line scoring ratios are correct.

**Test Cases**:
```go
testCases := []struct {
    lines      int
    wantScore  float64
    multiplier float64  // Relative to single-line
}{
    {0, 0.0, 0},
    {1, 0.40, 1.0},    // Base value
    {2, 2.00, 5.0},    // 5× single-line
    {3, 8.00, 20.0},   // 20× single-line
    {4, 24.00, 60.0},  // 60× single-line (Tetris!)
}
```

**Assertions**:
- ✅ All values match specification
- ✅ Exponential growth pattern verified
- ✅ 4-line > 3× 3-line (encourages waiting for Tetris)
- ✅ 3-line > 4× 2-line (encourages 3-line setups)

**Status**: ✅ IMPLEMENTED

---

#### **Test 1.2: EvaluateTwoPieceSequence_Basic**

**File**: `autoplay_test.go`

**Purpose**: Verify two-piece sequence evaluation returns valid scores.

**Test Scenarios**:

**Scenario A: Empty Board**
```go
gameState := NewGameState()
decision := &MoveDecision{rotations: 0, targetX: 5, softDrops: 10}
score := EvaluateTwoPieceSequence(gameState, decision, gameState.NextPiece)
```
**Expected**: score > -999999.0 (valid score, not error value)

**Scenario B: No Valid Next Moves**
```go
// Fill board completely except 1 row
for y := 0; y < 19; y++ {
    for x := 0; x < 10; x++ {
        board.Set(x, y, 1)
    }
}
```
**Expected**: score = -999999.0 (error value, no valid moves)

**Scenario C: Combo Setup**
```go
// Setup: 3 rows with 2-block gap (O + I can clear 4 lines)
for y := 0; y < 3; y++ {
    for x := 0; x < 10; x++ {
        if x < 4 || x > 5 {
            board.Set(x, y, 1)
        }
    }
}
```
**Expected**: score > single-piece evaluation (combo bonus applied)

**Status**: ✅ IMPLEMENTED

---

#### **Test 1.3: EvaluateTwoPieceSequence_EdgeCases**

**File**: `autoplay_test.go`

**Purpose**: Verify edge case handling.

**Edge Cases**:

1. **Nil Parameters**
   ```go
   EvaluateTwoPieceSequence(nil, decision, nextPiece)  // Should return -999999.0
   EvaluateTwoPieceSequence(gameState, nil, nextPiece) // Should return -999999.0
   EvaluateTwoPieceSequence(gameState, decision, nil)  // Should return -999999.0
   ```

2. **Game Over After Current Piece**
   ```go
   // Fill board to top, current piece will cause game over
   decision := EvaluateTwoPieceSequence(...)
   ```
   **Expected**: Returns -999999.0 (invalid sequence)

3. **Next Piece Causes Game Over**
   ```go
   // Current piece OK, but no space for next piece
   ```
   **Expected**: Lower score (penalized but not -999999.0)

**Status**: ⏳ TODO - Add to autoplay_test.go

---

#### **Test 1.4: enumerateMovesForBoard**

**File**: `autoplay_test.go`

**Purpose**: Verify move enumeration helper works correctly.

**Test Scenarios**:

1. **Empty Board - All Pieces**
   ```go
   for _, pieceType := range []TetrominoType{I, O, T, S, Z, J, L} {
       piece := NewTetromino(pieceType)
       moves := enumerateMovesForBoard(board, piece)
       assert(len(moves) > 0)
   }
   ```

2. **Narrow Gap - Only I-Piece Fits**
   ```go
   // Create 1-wide well
   moves := enumerateMovesForBoard(board, IPiece)
   assert(len(moves) > 0)
   
   moves = enumerateMovesForBoard(board, OPiece)
   assert(len(moves) == 0)  // O-piece too wide
   ```

3. **No Valid Moves**
   ```go
   // Fill board completely
   moves := enumerateMovesForBoard(fullBoard, piece)
   assert(len(moves) == 0)
   ```

**Status**: ⏳ TODO - Add to autoplay_test.go

---

#### **Test 1.5: isValidPositionForBoard**

**File**: `autoplay_test.go`

**Purpose**: Verify board collision detection.

**Test Matrix**:
| Test | Board State | Piece | Position | Expected |
|------|-------------|-------|----------|----------|
| Empty board, valid | All empty | I | x=5, y=18 | ✅ true |
| Empty board, out of bounds | All empty | I | x=10, y=18 | ❌ false |
| Collision test | Block at (5,16) | I | x=5, y=17 | ❌ false |
| Within bounds, clear | Block at (0,0) | I | x=9, y=18 | ✅ true |
| Rotation collision | Block at (6,18) | I (rotated) | x=5, y=18 | ❌ false |

**Status**: ⏳ TODO - Add to autoplay_test.go

---

### **INTEGRATION TESTS**

#### **Test 2.1: FindBestMoveWithNext_FindsCombo**

**File**: `autoplay_integration_test.go`

**Purpose**: Verify AI finds 2-piece Tetris setups.

**Setup**:
```go
// Create classic Tetris setup: 3 rows with 1-column well
for y := 0; y < 3; y++ {
    for x := 0; x < 10; x++ {
        if x != 4 {  // Leave column 4 open
            board.Set(x, y, 1)
        }
    }
}
gameState.CurrentPiece = NewTetromino(TetrominoO)  // O-piece fills adjacent
gameState.NextPiece = NewTetromino(TetrominoI)     // I-piece clears Tetris
```

**Expected Behavior**:
- AI places O-piece to preserve the well (column 4)
- O-piece should NOT fill column 4
- Setup should enable I-piece Tetris on next move

**Assertion**:
```go
decision := FindBestMoveWithNext(gameState)
assert(decision.targetX != 4)  // Don't fill the well!
assert(decision.targetX == 2 || decision.targetX == 6)  // Build around well
```

**Status**: ⏳ TODO - Replace TestTwoPieceLookahead_Capability with this

---

#### **Test 2.2: FindBestMoveWithNext_PrefersTetris**

**File**: `autoplay_integration_test.go`

**Purpose**: Verify AI chooses Tetris setup over multiple single clears.

**Setup**:
```go
// Option A: Place piece to clear 1 line now
// Option B: Place piece to set up 4-line Tetris with next piece

// Board state where both options are valid
AI should choose Option B (higher long-term score)
```

**Expected**: Decision leads to Tetris setup, not immediate single-line clear

**Metric**: Score(Tetris setup) > Score(Single clear × 4)

**Status**: ⏳ TODO

---

#### **Test 2.3: TwoPieceLookahead_Survival**

**File**: `autoplay_integration_test.go`

**Purpose**: Verify 2-piece lookahead doesn't reduce survival.

**Test**:
```go
func TestTwoPieceLookahead_Survival50Pieces(t *testing.T) {
    gameState := NewGameState()
    piecesPlaced := 0
    
    for piecesPlaced < 50 && !gameState.GameOver {
        decision := FindBestMoveWithNext(gameState)  // Use 2-piece lookahead
        ExecuteMove(gameState, decision)
        // Process line clears...
        piecesPlaced++
    }
    
    assert(piecesPlaced >= 40)  // Allow some variance, but maintain survival
}
```

**Status**: ✅ EXISTING (TestAutoPlay_Survival50Pieces) - Just verifies it still passes

---

#### **Test 2.4: LineClearRate_Improvement**

**File**: `autoplay_integration_test.go`

**Purpose**: Measure line-clear rate improvement from optimizations.

**Test**:
```go
func TestOptimizedAI_LineClearRate(t *testing.T) {
    gameState := NewGameState()
    piecesPlaced := 0
    
    for piecesPlaced < 100 && !gameState.GameOver {
        decision := FindBestMoveWithNext(gameState)
        ExecuteMove(gameState, decision)
        piecesPlaced++
    }
    
    lineClearRate := gameState.LinesCleared / piecesPlaced
    assert(lineClearRate >= 0.20)  // Target: ≥20% (was ~15%)
}
```

**Status**: ✅ EXISTING (TestAutoPlay_BaselineLineClearRate) - Enhanced with new target

---

#### **Test 2.5: TetrisFrequency_Count**

**File**: `autoplay_integration_test.go`

**Purpose**: Count actual 4-line clears (Tetris) executed.

**Test**:
```go
func TestTwoPieceLookahead_TetrisExecution(t *testing.T) {
    gameState := NewGameState()
    tetrisCount := 0
    piecesPlaced := 0
    
    for piecesPlaced < 100 && !gameState.GameOver {
        linesBefore := gameState.LinesCleared
        decision := FindBestMoveWithNext(gameState)
        ExecuteMove(gameState, decision)
        
        linesCleared := gameState.LinesCleared - linesBefore
        if linesCleared == 4 {
            tetrisCount++
        }
        piecesPlaced++
    }
    
    assert(tetrisCount >= 1)  // At least 1 Tetris in 100 pieces
    t.Logf("Executed %d Tetrises in %d pieces", tetrisCount, piecesPlaced)
}
```

**Target**: 1-3 Tetrises per 100 pieces (up from 0-1)

**Status**: ⏳ TODO - NEW TEST

---

### **PERFORMANCE TESTS**

#### **Test 3.1: BenchmarkFindBestMoveWithNext**

**File**: `autoplay_test.go`

**Purpose**: Measure performance impact of 2-piece lookahead.

**Benchmark**:
```go
func BenchmarkFindBestMoveWithNext(b *testing.B) {
    gameState := NewGameState()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        FindBestMoveWithNext(gameState)
    }
}
```

**Target**: <10ms per call (acceptable for real-time play)

**Comparison**:
- Single-piece (FindBestMove): ~1-2ms
- Two-piece (FindBestMoveWithNext): Expected ~5-10ms
- Acceptable overhead: ≤5× single-piece

**Status**: ⏳ TODO - Add benchmark

---

#### **Test 3.2: BenchmarkEvaluateTwoPieceSequence**

**File**: `autoplay_test.go`

**Purpose**: Measure core evaluation function performance.

**Benchmark**:
```go
func BenchmarkEvaluateTwoPieceSequence(b *testing.B) {
    gameState := NewGameState()
    decision := &MoveDecision{rotations: 0, targetX: 5, softDrops: 10}
    nextPiece := NewTetromino(TetrominoI)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        EvaluateTwoPieceSequence(gameState, decision, nextPiece)
    }
}
```

**Target**: <1ms per evaluation

**Status**: ⏳ TODO - Add benchmark

---

#### **Test 3.3: MemoryAllocation_Test**

**File**: `autoplay_test.go`

**Purpose**: Ensure 2-piece lookahead doesn't cause excessive allocations.

**Test**:
```go
func TestTwoPieceLookahead_NoExcessiveAllocations(t *testing.T) {
    gameState := NewGameState()
    decision := &MoveDecision{rotations: 0, targetX: 5, softDrops: 10}
    
    allocs := testing.AllocsPerRun(100, func() {
        EvaluateTwoPieceSequence(gameState, decision, gameState.NextPiece)
    })
    
    // Should allocate ~2-3 boards for simulation
    if allocs > 10 {
        t.Errorf("Excessive allocations: %.1f per call (expected <10)", allocs)
    }
}
```

**Status**: ⏳ TODO - Add test

---

### **REGRESSION TESTS**

#### **Test 4.1: BackwardCompatibility_FindBestMove**

**File**: `autoplay_test.go`

**Purpose**: Ensure original `FindBestMove()` still works.

**Test**:
```go
func TestFindBestMove_StillWorks(t *testing.T) {
    gameState := NewGameState()
    decision := FindBestMove(gameState)
    
    assert(decision != nil)
    assert(decision.IsValid())
}
```

**Status**: ✅ EXISTING (TestFindBestMove) - Verify still passes

---

#### **Test 4.2: Weights_GetWeights**

**File**: `autoplay_test.go`

**Purpose**: Verify weight functions work correctly.

**Test**:
```go
func TestGetWeights_AfterOptimization(t *testing.T) {
    weights := GetWeights()
    
    // Verify all expected keys exist
    expectedKeys := []string{"aggregateHeight", "holes", "bumpiness", "wells"}
    for _, key := range expectedKeys {
        _, exists := weights[key]
        assert(exists)
    }
    
    // Verify values are in reasonable ranges
    assert(weights["aggregateHeight"] < 0)  // Should be negative
    assert(weights["holes"] < 0)            // Should be negative
}
```

**Status**: ✅ EXISTING (TestGetWeights_SetWeights) - Enhanced

---

#### **Test 4.3: ExecuteMove_WithTwoPiecePlanning**

**File**: `autoplay_integration_test.go`

**Purpose**: Verify ExecuteMove works with 2-piece planned moves.

**Test**:
```go
func TestExecuteMove_TwoPiecePlanned(t *testing.T) {
    gameState := NewGameState()
    decision := FindBestMoveWithNext(gameState)
    
    // Execute should work the same regardless of how decision was made
    ExecuteMove(gameState, decision)
    
    assert(gameState.CurrentPiece != nil)  // New piece spawned
    // Piece should be locked to board
}
```

**Status**: ✅ COVERED (TestAutoPlay_Survival10Pieces)

---

## Test Execution Plan

### Phase 1: Unit Tests (First)
```bash
go test ./internal/model -run "TestEvaluateLineClears|TestEvaluateTwoPiece" -v
```
**Expected**: All pass before integration testing

### Phase 2: Integration Tests
```bash
go test ./internal/model -run "TestTwoPieceLookahead|TestFindBestMoveWithNext" -v
```
**Expected**: AI demonstrates 2-piece planning behavior

### Phase 3: Performance Tests
```bash
go test ./internal/model -bench="BenchmarkFindBestMove|BenchmarkEvaluateTwoPiece" -benchmem
```
**Expected**: Meet performance targets

### Phase 4: Full Suite
```bash
go test ./internal/model/... -v
go test ./internal/model -race
```
**Expected**: All pass, no race conditions

---

## Test Coverage Goals

| Component | Current Coverage | Target Coverage |
|-----------|-----------------|-----------------|
| `EvaluateTwoPieceSequence()` | ~60% | 90%+ |
| `FindBestMoveWithNext()` | ~50% | 90%+ |
| `enumerateMovesForBoard()` | 0% | 80%+ |
| `isValidPositionForBoard()` | 0% | 80%+ |
| `evaluateLineClears()` | 100% | 100% ✅ |

**Overall Target**: 85%+ coverage for autoplay.go

---

## Evidence Collection

For each test category, capture:

1. **Unit Tests**: Test output logs showing pass/fail
2. **Integration Tests**: AI decision logs (targetX, rotations, score comparisons)
3. **Performance Tests**: Benchmark results with allocs/op, ns/op
4. **Regression Tests**: Comparison with baseline metrics

**Evidence Location**: `.sisyphus/evidence/task-{N}-{test-name}.txt`

---

## Test Data Examples

### Board State Templates

```go
// Template 1: Classic Tetris Well (3 rows, 1-column gap)
func createTetrisWellBoard() *Board {
    board := NewBoard()
    for y := 0; y < 3; y++ {
        for x := 0; x < 10; x++ {
            if x != 4 {  // Column 4 is the well
                board.Set(x, y, 1)
            }
        }
    }
    return board
}

// Template 2: 2-Row Setup for O + I Combo
func createComboSetupBoard() *Board {
    board := NewBoard()
    for y := 0; y < 2; y++ {
        for x := 0; x < 10; x++ {
            if x < 4 || x > 5 {  // Columns 4-5 open
                board.Set(x, y, 1)
            }
        }
    }
    return board
}

// Template 3: Dangerous Well (depth 4+)
func createDangerousWellBoard() *Board {
    board := NewBoard()
    for x := 0; x < 10; x++ {
        height := 5
        if x == 5 {  // Well at column 5
            height = 1
        }
        for y := 0; y < height; y++ {
            board.Set(x, y, 1)
        }
    }
    return board
}
```

---

## Test Result Templates

### Unit Test Result
```
=== RUN   TestEvaluateLineClears_EnhancedRatios
=== RUN   TestEvaluateLineClears_EnhancedRatios/0_lines
=== RUN   TestEvaluateLineClears_EnhancedRatios/1_lines
=== RUN   TestEvaluateLineClears_EnhancedRatios/2_lines
=== RUN   TestEvaluateLineClears_EnhancedRatios/3_lines
=== RUN   TestEvaluateLineClears_EnhancedRatios/4_lines
--- PASS: TestEvaluateLineClears_EnhancedRatios (0.00s)
    PASS: All ratios correct (1×, 5×, 20×, 60×)
```

### Integration Test Result
```
=== RUN   TestTwoPieceLookahead_TetrisExecution
    autoplay_integration_test.go:XXX: Executed 3 Tetrises in 87 pieces
    autoplay_integration_test.go:XXX: Tetris frequency: 3.45%
--- PASS: TestTwoPieceLookahead_TetrisExecution (0.02s)
    PASS: ≥1 Tetris executed (target met)
```

### Performance Test Result
```
goos: darwin
goarch: arm64
BenchmarkFindBestMoveWithNext-10          100    8234567 ns/op    2567890 B/op    45 allocs/op
BenchmarkFindBestMove-10                 1000    1234567 ns/op     345678 B/op     8 allocs/op
```
**Verdict**: 2-piece is 6.7× slower (acceptable, <10ms per call)

---

## Continuous Testing

### Pre-commit Checks
```bash
go test ./internal/model -run "Test.*TwoPiece|Test.*LineClear" -v
go build ./...
go vet ./...
```

### CI/CD Integration
Add to `.github/workflows/test.yml`:
```yaml
- name: Run Autoplay Tests
  run: go test ./internal/model -run "TestAutoPlay" -v
  
- name: Run Performance Tests
  run: go test ./internal/model -bench="BenchmarkFindBest" -benchmem
```

---

## Test Maintenance

### When to Update Tests

1. **Algorithm Changes**: Update Test 2.1, 2.2 (integration tests)
2. **Weight Tuning**: Update Test 1.1 (line clear ratios)
3. **Performance Regression**: Update Test 3.1, 3.2 (benchmarks)
4. **New Features**: Add new test cases to appropriate category

### Test Debt Tracking

Monitor:
- ❌ Tests marked ⏳ TODO in this document
- ⚠️ Flaky tests (intermittent failures due to randomness)
- 📉 Performance degradation (>20% slower than baseline)

---

## Success Criteria Summary

### Must Pass Before Merge:
- ✅ All unit tests (Category 1)
- ✅ Core integration tests (Test 2.1, 2.3, 2.4)
- ✅ No performance regression (Test 3.1)
- ✅ All regression tests (Category 4)

### Target Metrics:
- Line-clear rate: ≥20% (up from ~15%)
- Tetris frequency: ≥1 per 100 pieces
- Survival: 50+ pieces maintained
- Decision time: <10ms per move

---

## Testing Checklist

**Before Starting Implementation**:
- [ ] Read this testing plan
- [ ] Identify which tests are already implemented
- [ ] Create test file structure for TODO tests

**During Implementation**:
- [ ] Write test BEFORE implementing function (TDD)
- [ ] Run unit tests after each function
- [ ] Capture evidence for each passing test

**After Implementation**:
- [ ] Run full test suite
- [ ] Verify all TODO tests implemented
- [ ] Document actual vs expected results
- [ ] Update this plan with outcomes

---

**Last Updated**: 2026-03-14
**Author**: Prometheus (Planning Agent)
**Status**: Ready for execution
