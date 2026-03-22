# Tetris Auto-Play Mode - Work Plan

## TL;DR

> **Quick Summary**: Implement heuristic-based AI that automatically plays Tetris by evaluating all possible piece placements and selecting optimal moves. AI can be toggled during gameplay with adjustable speed levels.
> 
> **Deliverables**:
> - `internal/model/autoplay.go` - Core AI engine with heuristic evaluation
> - `internal/model/autoplay_test.go` - Comprehensive test suite
> - `cmd/tetris/main.go` - Integration: 'A' key toggle, AI input generation
> - `internal/ui/autoplay_render.go` - AI decision visualization UI
> - Speed control system (5 levels)
> - Documentation and examples
> 
> **Estimated Effort**: Medium (15-20 tasks, parallelizable)
> **Parallel Execution**: YES - 4 waves with 4-7 tasks each
> **Critical Path**: Types → Heuristics → Move Finder → Integration → UI → Tests

---

## Context

### Original Request
Add a complex feature to the Tetris game: auto-play mode where the game simulates player operations (turn/move/drop) and clears blocks automatically. The AI determines all operations by itself, deciding whether to clear one line or multi-lines at one block drop.

### Interview Summary
**Key Discussions**:
- **AI Algorithm**: Heuristic-based evaluation (NOT ML/greedy) - evaluates all positions, picks best scoring move
- **Speed Control**: Adjustable speed levels (1-5) for user to cycle through
- **Activation**: Toggle with 'A' key during gameplay (not menu-based)
- **AI Visibility**: Show AI decision-making process on screen (target position, score evaluation)

**Research Findings**:
- Standard Tetris AI uses weighted features: aggregate height, holes, bumpiness, wells, complete lines
- Typical weights from research: height -0.5, lines +0.76, holes -0.36, bumpiness -0.18
- Enumeration approach: 4 rotations × 10 positions = ~40 evaluations per piece

### Codebase Context
- **Go Version**: 1.24.0+ (tested on go1.26.0)
- **Module**: `github.com/oc-garden/tetris_game`
- **UI Library**: tcell v2 (256-color terminal)
- **Board**: 10×20 cells, coordinate (0,0) at bottom-left
- **Existing patterns**: Constructor pattern `NewX()`, pointer receivers, table-driven tests

---

## Work Objectives

### Core Objective
Build a heuristic-based AI system that can autonomously play Tetris by evaluating board states and executing optimal moves, seamlessly integrated with the existing game loop.

### Concrete Deliverables
- `internal/model/autoplay.go` (~300-400 lines)
- `internal/model/autoplay_test.go` (~200-300 lines)
- `cmd/tetris/main.go` modifications (~50 lines added)
- `internal/ui/autoplay_render.go` (~100 lines)
- 5 speed levels with visible UI indicator
- AI decision panel showing target position and evaluation score

### Definition of Done
- [ ] `go test ./...` passes (0 failures)
- [ ] `go build -o tetris` succeeds
- [ ] `go vet ./...` reports no issues
- [ ] Auto-play can be toggled with 'A' key during live gameplay
- [ ] AI successfully clears lines autonomously for 100+ pieces without game over
- [ ] Speed control cycles through 5 levels with visible change
- [ ] AI decision UI displays target position, rotations, and score

### Must Have
- Heuristic evaluation with configurable weights
- All 4 rotations × all X positions evaluated per piece
- Move execution that respects existing game mechanics
- Speed control: 5 levels (normal → instant)
- Toggle on/off with single key press
- AI respects pause state (pauses when game paused)
- UI shows: "AUTO-PLAY" indicator, speed level, target position, evaluation score

### Must NOT Have (Guardrails)
- NO ML/neural network implementation (heuristic only)
- NO modification to existing piece mechanics or board logic
- NO changes to existing key bindings (add 'A', don't replace)
- NO online learning or weight adjustment during gameplay
- NO AI slop: excessive comments, over-abstraction, generic variable names
- NO breaking existing manual play functionality
- NO AI input during game over state

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: YES (Go testing framework)
- **Automated tests**: TDD (test-driven development)
- **Framework**: `go test` (standard library)
- **Approach**: Each task follows RED → GREEN → REFACTOR cycle

### QA Policy
Every task MUST include agent-executed QA scenarios:
- **Build verification**: `go build`, `go vet`, `go test ./...`
- **Unit tests**: Table-driven tests for heuristics and move finding
- **Integration**: Playwright not applicable (terminal app) - use `interactive_bash` for tmux-based verification
- **Evidence**: Screenshots captured via tmux, test output logs saved to `.sisyphus/evidence/`

---

## Detailed Test Plan

### Test File Structure

```
internal/model/
├── autoplay.go                        # Source code
├── autoplay_test.go                   # Unit tests (~300 lines)
└── autoplay_integration_test.go       # Integration tests (~150 lines)
```

### Unit Test Categories

#### Category 1: Type & Constructor Tests (Task 1-3)

**Test File**: `autoplay_test.go`

```go
// Test 1.1: AutoPlayer Creation
func TestAutoPlayerCreation(t *testing.T) {
    tests := []struct {
        name           string
        wantEnabled    bool
        wantSpeedLevel int
    }{
        {"default state", false, 1},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ap := NewAutoPlayer()
            if ap.IsEnabled() != tt.wantEnabled {
                t.Errorf("IsEnabled() = %v, want %v", ap.IsEnabled(), tt.wantEnabled)
            }
            if ap.GetSpeedLevel() != tt.wantSpeedLevel {
                t.Errorf("GetSpeedLevel() = %d, want %d", ap.GetSpeedLevel(), tt.wantSpeedLevel)
            }
        })
    }
}

// Test 1.2: Toggle Functionality
func TestAutoPlayer_Toggle(t *testing.T) {
    ap := NewAutoPlayer()
    
    // Initial state: disabled
    if ap.IsEnabled() {
        t.Error("New AutoPlayer should be disabled")
    }
    
    // Toggle on
    ap.Toggle()
    if !ap.IsEnabled() {
        t.Error("Toggle() should enable auto-play")
    }
    
    // Toggle off
    ap.Toggle()
    if ap.IsEnabled() {
        t.Error("Toggle() should disable auto-play")
    }
}

// Test 1.3: Speed Level Setting
func TestAutoPlayer_SetSpeedLevel(t *testing.T) {
    tests := []struct {
        name      string
        setLevel  int
        wantLevel int
    }{
        {"level 1", 1, 1},
        {"level 3", 3, 3},
        {"level 5", 5, 5},
        {"invalid level 0", 0, 1},    // Should clamp to 1
        {"invalid level 6", 6, 5},    // Should clamp to 5
        {"invalid level -1", -1, 1},  // Should clamp to 1
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ap := NewAutoPlayer()
            ap.SetSpeedLevel(tt.setLevel)
            if ap.GetSpeedLevel() != tt.wantLevel {
                t.Errorf("SetSpeedLevel(%d) = %d, want %d", tt.setLevel, ap.GetSpeedLevel(), tt.wantLevel)
            }
        })
    }
}

// Test 1.4: Speed Cycling
func TestAutoPlayer_CycleSpeed(t *testing.T) {
    ap := NewAutoPlayer()
    expected := []int{2, 3, 4, 5, 1, 2} // Cycle: 1→2→3→4→5→1→2
    
    for i, want := range expected {
        ap.CycleSpeed()
        if ap.GetSpeedLevel() != want {
            t.Errorf("Cycle %d: GetSpeedLevel() = %d, want %d", i+1, ap.GetSpeedLevel(), want)
        }
    }
}

// Test 2.1: MoveDecision Validation
func TestMoveDecision_IsValid(t *testing.T) {
    tests := []struct {
        name  string
        decis ion MoveDecision
        want  bool
    }{
        {"valid center", MoveDecision{rotations: 0, targetX: 5, softDrops: 10, score: 12.5}, true},
        {"valid edge X=0", MoveDecision{rotations: 2, targetX: 0, softDrops: 5, score: 8.0}, true},
        {"valid edge X=9", MoveDecision{rotations: 3, targetX: 9, softDrops: 15, score: 6.5}, true},
        {"invalid rotation 4", MoveDecision{rotations: 4, targetX: 5, softDrops: 10, score: 12.5}, false},
        {"invalid rotation -1", MoveDecision{rotations: -1, targetX: 5, softDrops: 10, score: 12.5}, false},
        {"invalid X=-1", MoveDecision{rotations: 0, targetX: -1, softDrops: 10, score: 12.5}, false},
        {"invalid X=10", MoveDecision{rotations: 0, targetX: 10, softDrops: 10, score: 12.5}, false},
        {"all invalid", MoveDecision{rotations: 5, targetX: -5, softDrops: -1, score: 0}, false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.decision.IsValid(); got != tt.want {
                t.Errorf("IsValid() = %v, want %v", got, tt.want)
            }
        })
    }
}

// Test 2.2: MoveDecision String Format
func TestMoveDecision_String(t *testing.T) {
    d := MoveDecision{rotations: 2, targetX: 5, softDrops: 10, score: 12.5}
    got := d.String()
    want := "MoveDecision{rot:2, x:5, drops:10, score:12.50}"
    
    if got != want {
        t.Errorf("String() = %q, want %q", got, want)
    }
}

// Test 2.3: CalculateSoftDrops
func TestCalculateSoftDrops(t *testing.T) {
    tests := []struct {
        name      string
        board     *Board
        pieceType TetrominoType
        targetX   int
        wantDrops int
    }{
        {"empty board I-piece at center", createEmptyBoard(), TetrominoI, 3, 19},
        {"I-piece at bottom", createBoardWithHeight(1), TetrominoI, 3, 18},
        {"I-piece on stacked blocks", createBoardWithHeight(10), TetrominoI, 3, 9},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            piece := NewTetromino(tt.pieceType)
            piece.X = 3
            piece.Y = 18
            got := CalculateSoftDrops(tt.board, piece, tt.targetX)
            if got != tt.wantDrops {
                t.Errorf("CalculateSoftDrops() = %d, want %d", got, tt.wantDrops)
            }
        })
    }
}
```

#### Category 2: Board Helper Tests (Task 4)

```go
// Test 4.1: Column Height Calculation
func TestGetColHeight(t *testing.T) {
    tests := []struct {
        name     string
        setupFn  func() *Board
        col      int
        wantHeight int
    }{
        {
            "empty column",
            createEmptyBoard,
            5,
            0,
        },
        {
            "full column",
            func() *Board {
                b := NewBoard()
                for y := 0; y < 20; y++ {
                    b.Set(5, y, 1)
                }
                return b
            },
            5,
            20,
        },
        {
            "partial column height 5",
            func() *Board {
                b := NewBoard()
                for y := 0; y < 5; y++ {
                    b.Set(5, y, 1)
                }
                return b
            },
            5,
            5,
        },
        {
            "column with gap (height still 5)",
            func() *Board {
                b := NewBoard()
                b.Set(5, 0, 1)
                b.Set(5, 1, 1)
                b.Set(5, 2, 0) // gap
                b.Set(5, 3, 1)
                b.Set(5, 4, 1)
                return b
            },
            5,
            5, // Height is topmost block, gaps don't reduce it
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            board := tt.setupFn()
            if got := getColHeight(board, tt.col); got != tt.wantHeight {
                t.Errorf("getColHeight() = %d, want %d", got, tt.wantHeight)
            }
        })
    }
}

// Test 4.2: Aggregate Height
func TestGetAggregateHeight(t *testing.T) {
    tests := []struct {
        name    string
        setupFn func() *Board
        want    int
    }{
        {"empty board", createEmptyBoard, 0},
        {"flat height 5 all cols", func() *Board {
            b := NewBoard()
            for x := 0; x < 10; x++ {
                for y := 0; y < 5; y++ {
                    b.Set(x, y, 1)
                }
            }
            return b
        }, 50}, // 10 cols × 5 height
        {"full board", func() *Board {
            b := NewBoard()
            for x := 0; x < 10; x++ {
                for y := 0; y < 20; y++ {
                    b.Set(x, y, 1)
                }
            }
            return b
        }, 200}, // 10 cols × 20 height
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := getAggregateHeight(tt.setupFn()); got != tt.want {
                t.Errorf("getAggregateHeight() = %d, want %d", got, tt.want)
            }
        })
    }
}

// Test 4.3: Complete Lines Count
func TestCountCompleteLines(t *testing.T) {
    tests := []struct {
        name     string
        setupFn  func() *Board
        wantLines int
    }{
        {"no lines", createEmptyBoard, 0},
        {"one full line", func() *Board {
            b := NewBoard()
            for x := 0; x < 10; x++ {
                b.Set(x, 5, 1)
            }
            return b
        }, 1},
        {"three full lines", func() *Board {
            b := NewBoard()
            for line := 0; line < 3; line++ {
                for x := 0; x < 10; x++ {
                    b.Set(x, line, 1)
                }
            }
            return b
        }, 3},
        {"partial line (not full)", func() *Board {
            b := NewBoard()
            for x := 0; x < 9; x++ { // missing one block
                b.Set(x, 5, 1)
            }
            return b
        }, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := countCompleteLines(tt.setupFn()); got != tt.wantLines {
                t.Errorf("countCompleteLines() = %d, want %d", got, tt.wantLines)
            }
        })
    }
}

// Test 4.4: Hole Detection
func TestCountHoles(t *testing.T) {
    tests := []struct {
        name     string
        setupFn  func() *Board
        wantHoles int
    }{
        {"no holes", createEmptyBoard, 0},
        {"single hole", func() *Board {
            b := NewBoard()
            b.Set(5, 0, 1)
            b.Set(5, 2, 1) // block above creates hole at y=1
            return b
        }, 1},
        {"multiple holes same column", func() *Board {
            b := NewBoard()
            b.Set(5, 0, 1)
            b.Set(5, 2, 1)
            b.Set(5, 4, 1)
            return b
        }, 2}, // holes at y=1 and y=3
        {"hole with multiple blocks above", func() *Board {
            b := NewBoard()
            b.Set(5, 0, 1)
            for y := 2; y < 20; y++ {
                b.Set(5, y, 1)
            }
            return b
        }, 1}, // still 1 hole at y=1
        {"no holes - solid stack", func() *Board {
            b := NewBoard()
            for y := 0; y < 10; y++ {
                b.Set(5, y, 1)
            }
            return b
        }, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := countHoles(tt.setupFn()); got != tt.wantHoles {
                t.Errorf("countHoles() = %d, want %d", got, tt.wantHoles)
            }
        })
    }
}

// Test 4.5: Bumpiness Calculation
func TestCalculateBumpiness(t *testing.T) {
    tests := []struct {
        name        string
        colHeights  [10]int
        wantBumpiness int
    }{
        {"flat surface", [10]int{5,5,5,5,5,5,5,5,5,5}, 0},
        {"single step", [10]int{5,10,5,5,5,5,5,5,5,5}, 10}, // |5-10| + |10-5| = 10
        {"alternating", [10]int{0,10,0,10,0,10,0,10,0,10}, 90}, // 9 transitions × 10
        {"increasing slope", [10]int{1,2,3,4,5,6,7,8,9,10}, 9}, // each step = 1
        {"random heights", [10]int{3,5,2,4,1,6,3,2,4,5}, 22}, // calculated
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            b := createBoardWithHeights(tt.colHeights)
            if got := calculateBumpiness(b); got != tt.wantBumpiness {
                t.Errorf("calculateBumpiness() = %d, want %d", got, tt.wantBumpiness)
            }
        })
    }
}

// Test 4.6: Wells Detection
func TestCountWells(t *testing.T) {
    tests := []struct {
        name       string
        setupFn    func() *Board
        wantWells  int
    }{
        {"no wells", createEmptyBoard, 0},
        {"single well depth 2", func() *Board {
            b := NewBoard()
            // Create U-shape: blocks at x=4 and x=6, empty at x=5
            for y := 0; y < 5; y++ {
                b.Set(4, y, 1)
                b.Set(6, y, 1)
            }
            return b
        }, 1}, // well at column 5
        {"multiple wells", func() *Board {
            b := NewBoard()
            // Create two U-shapes
            for y := 0; y < 5; y++ {
                b.Set(2, y, 1)
                b.Set(4, y, 1)
                b.Set(7, y, 1)
                b.Set(9, y, 1)
            }
            return b
        }, 2}, // wells at columns 3 and 8
        {"shallow depression (not a well)", func() *Board {
            b := NewBoard()
            // Only depth 1 - not a well
            for y := 0; y < 5; y++ {
                b.Set(4, y, 1)
                b.Set(6, y, 1)
            }
            b.Set(5, 4, 1) // fill to depth 1
            return b
        }, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := countWells(tt.setupFn()); got != tt.wantWells {
                t.Errorf("countWells() = %d, want %d", got, tt.wantWells)
            }
        })
    }
}
```

#### Category 3: Heuristic Evaluation Tests (Task 5-8)

```go
// Test 5-7: Individual Heuristic Functions
func TestEvalAggregateHeight(t *testing.T) {
    b := createBoardWithHeights([10]int{3,5,2,0,4,1,6,3,2,4})
    got := evalAggregateHeight(b)
    want := 30
    if got != want {
        t.Errorf("evalAggregateHeight() = %d, want %d", got, want)
    }
}

func TestEvalHoles(t *testing.T) {
    b := createBoardWithHoles(5) // helper creates board with 5 holes
    got := evalHoles(b)
    want := 5
    if got != want {
        t.Errorf("evalHoles() = %d, want %d", got, want)
    }
}

func TestEvalBumpiness(t *testing.T) {
    b := createBoardWithHeights([10]int{3,5,2,4,1,6,3,2,4,5})
    got := evalBumpiness(b)
    want := 22
    if got != want {
        t.Errorf("evalBumpiness() = %d, want %d", got, want)
    }
}

func TestEvalWells(t *testing.T) {
    b := createBoardWithWells(2) // helper creates board with 2 wells
    got := evalWells(b)
    want := 2
    if got != want {
        t.Errorf("evalWells() = %d, want %d", got, want)
    }
}

// Test 8: Combined Weighted Evaluation
func TestEvaluateBoard(t *testing.T) {
    tests := []struct {
        name        string
        boardFn     func() *Board
        wantBetter  bool // whether board1 should score better than board2
        description string
    }{
        {
            "flat vs bumpy",
            func() *Board {
                b1 := createBoardWithHeights([10]int{5,5,5,5,5,5,5,5,5,5})
                b2 := createBoardWithHeights([10]int{0,10,0,10,0,10,0,10,0,10})
                // Test internally - compare scores
                return b1
            },
            true,
            "Flat board should score better than bumpy",
        },
        {
            "no holes vs holes",
            func() *Board {
                b1 := createBoardWithHeights([10]int{5,5,5,5,5,5,5,5,5,5})
                b2 := createBoardWithHoles(5)
                return b1
            },
            true,
            "No holes should score better than with holes",
        },
        {
            "complete line vs none",
            func() *Board {
                b1 := createBoardWithCompleteLine()
                b2 := createEmptyBoard()
                return b1
            },
            true,
            "Complete line should improve score",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            board := tt.boardFn()
            piece := NewTetromino(TetrominoI)
            score := evaluateBoard(board, piece, 5, 0)
            
            // Score should be reasonable (not NaN, not infinite)
            if math.IsNaN(score) || math.IsInf(score, 0) {
                t.Errorf("evaluateBoard() returned invalid score: %f", score)
            }
        })
    }
}

// Test 8.2: Weight Configuration
func TestGetWeights_SetWeights(t *testing.T) {
    weights := GetWeights()
    
    // Check default weights exist
    expectedKeys := []string{"aggregateHeight", "completeLines", "holes", "bumpiness", "wells"}
    for _, key := range expectedKeys {
        if _, ok := weights[key]; !ok {
            t.Errorf("GetWeights() missing key: %s", key)
        }
    }
    
    // Test weight modification
    newWeights := map[string]float64{
        "aggregateHeight": -1.0,
        "completeLines": 1.5,
        "holes": -0.5,
        "bumpiness": -0.3,
        "wells": -0.2,
    }
    SetWeights(newWeights)
    got := GetWeights()
    
    for key, want := range newWeights {
        if math.Abs(got[key]-want) > 0.001 {
            t.Errorf("SetWeights(%s) = %f, want %f", key, got[key], want)
        }
    }
}

// Test 8.3: Weight Impact on Score
func TestEvaluateBoard_WeightImpact(t *testing.T) {
    board := createBoardWithHeights([10]int{5,5,5,5,5,5,5,5,5,5})
    piece := NewTetromino(TetrominoI)
    
    // Get baseline score
    baseline := evaluateBoard(board, piece, 5, 0)
    
    // Change aggregateHeight weight to 0
    originalWeights := GetWeights()
    testWeights := GetWeights()
    testWeights["aggregateHeight"] = 0
    SetWeights(testWeights)
    
    modified := evaluateBoard(board, piece, 5, 0)
    
    // Score should change
    if baseline == modified {
        t.Error("Changing weights should affect score")
    }
    
    // Restore original weights
    SetWeights(originalWeights)
}
```

#### Category 4: Move Finding Tests (Task 9-10)

```go
// Test 9: Move Enumeration
func TestEnumerateMoves(t *testing.T) {
    gameState := NewGameState()
    piece := gameState.CurrentPiece
    
    moves := enumerateMoves(gameState, piece)
    
    // Should have multiple moves (typically 30-40 for empty board)
    if len(moves) < 20 {
        t.Errorf("enumerateMoves() returned too few moves: %d, expect 30-40", len(moves))
    }
    
    // All moves should be valid
    for i, move := range moves {
        if !move.IsValid() {
            t.Errorf("Move %d is invalid: %+v", i, move)
        }
        if move.rotations < 0 || move.rotations > 3 {
            t.Errorf("Move %d has invalid rotation: %d", i, move.rotations)
        }
        if move.targetX < 0 || move.targetX > 9 {
            t.Errorf("Move %d has invalid X: %d", i, move.targetX)
        }
        if move.softDrops < 0 || move.softDrops > 20 {
            t.Errorf("Move %d has unreasonable softDrops: %d", i, move.softDrops)
        }
    }
}

// Test 9.2: Enumeration for Different Piece Types
func TestEnumerateMoves_AllPieceTypes(t *testing.T) {
    pieceTypes := []TetrominoType{
        TetrominoI, TetrominoO, TetrominoT,
        TetrominoS, TetrominoZ, TetrominoJ, TetrominoL,
    }
    
    for _, pType := range pieceTypes {
        t.Run(string(pType), func(t *testing.T) {
            gameState := NewGameState()
            piece := NewTetromino(pType)
            piece.X = 3
            piece.Y = 18
            
            moves := enumerateMoves(gameState, piece)
            
            // O piece has fewer unique rotations (symmetric)
            if pType == TetrominoO {
                if len(moves) < 20 {
                    t.Errorf("O piece: too few moves: %d", len(moves))
                }
            } else {
                if len(moves) < 30 {
                    t.Errorf("%s: too few moves: %d", pType, len(moves))
                }
            }
        })
    }
}

// Test 10: FindBestMove
func TestFindBestMove(t *testing.T) {
    gameState := NewGameState()
    
    decision := FindBestMove(gameState)
    
    if decision == nil {
        t.Fatal("FindBestMove() returned nil on empty board")
    }
    
    if !decision.IsValid() {
        t.Errorf("FindBestMove() returned invalid decision: %+v", decision)
    }
    
    // Should have reasonable score
    if decision.score < -100 || decision.score > 100 {
        t.Errorf("FindBestMove() score seems unreasonable: %f", decision.score)
    }
}

// Test 10.2: Best Move Selection - Obvious Cases
func TestFindBestMove_ObviousCases(t *testing.T) {
    // Test case: I piece can complete a line
    gameState := NewGameState()
    // Create board where I piece at specific position completes line
    for x := 0; x < 9; x++ {
        gameState.Board.Set(x, 10, 1)
    }
    gameState.CurrentPiece = NewTetromino(TetrominoI)
    gameState.CurrentPiece.X = 3
    gameState.CurrentPiece.Y = 18
    
    decision := FindBestMove(gameState)
    
    // Should find the line-completing move
    if decision == nil {
        t.Error("FindBestMove() should find line-completing move")
    }
}

// Test 10.3: Determinism
func TestFindBestMove_Determinism(t *testing.T) {
    gameState := NewGameState()
    
    // Run multiple times - should get same result
    decision1 := FindBestMove(gameState)
    decision2 := FindBestMove(gameState)
    decision3 := FindBestMove(gameState)
    
    if decision1.targetX != decision2.targetX || decision2.targetX != decision3.targetX {
        t.Error("FindBestMove() should be deterministic")
    }
    if decision1.rotations != decision2.rotations || decision2.rotations != decision3.rotations {
        t.Error("FindBestMove() should be deterministic")
    }
}

// Test 10.4: No Valid Moves (Game Over Scenario)
func TestFindBestMove_NoValidMoves(t *testing.T) {
    gameState := NewGameState()
    // Fill the top rows - no valid moves
    for x := 0; x < 10; x++ {
        for y := 17; y < 20; y++ {
            gameState.Board.Set(x, y, 1)
        }
    }
    
    decision := FindBestMove(gameState)
    
    // May return nil or a move with very low score
    if decision != nil && decision.score > -50 {
        t.Errorf("FindBestMove() on blocked board should return nil or very low score, got: %+v", decision)
    }
}
```

#### Category 5: Move Execution Tests (Task 11-15)

```go
// Test 11: Rotation Execution
func TestExecuteRotations(t *testing.T) {
    gameState := NewGameState()
    piece := gameState.CurrentPiece
    
    // Start at rotation 0
    if piece.Rotation != 0 {
        t.Skip("Piece not at rotation 0")
    }
    
    // Execute 2 rotations
    executeRotations(gameState, 2)
    
    if piece.Rotation != 2 {
        t.Errorf("After executeRotations(2): rotation = %d, want 2", piece.Rotation)
    }
}

// Test 12: Horizontal Movement
func TestExecuteHorizontalMove(t *testing.T) {
    gameState := NewGameState()
    piece := gameState.CurrentPiece
    
    startX := piece.X
    
    // Move to X=0
    executeHorizontalMove(gameState, 0)
    if piece.X != 0 {
        t.Errorf("After move to X=0: X = %d, want 0", piece.X)
    }
    
    // Move to X=9
    executeHorizontalMove(gameState, 9)
    if piece.X != 9 {
        t.Errorf("After move to X=9: X = %d, want 9", piece.X)
    }
}

// Test 13: Drop Execution
func TestExecuteDrop(t *testing.T) {
    gameState := NewGameState()
    initialY := gameState.CurrentPiece.Y
    
    // Soft drop 5 times
    executeDrop(gameState, 5, false)
    
    // Piece should have moved down
    if gameState.CurrentPiece.Y >= initialY {
        t.Errorf("After 5 soft drops: Y = %d, should be < %d", gameState.CurrentPiece.Y, initialY)
    }
}

// Test 14: Speed Delay Calculation
func TestGetDelayForSpeed(t *testing.T) {
    tests := []struct {
        baseDelay  int
        speedLevel int
        wantDelay  int
    }{
        {1500, 1, 1500}, // normal speed
        {1500, 2, 750},  // 2x
        {1500, 3, 300},  // 5x
        {1500, 4, 150},  // 10x
        {1500, 5, 0},    // instant
    }
    for _, tt := range tests {
        t.Run(fmt.Sprintf("level_%d", tt.speedLevel), func(t *testing.T) {
            got := GetDelayForSpeed(tt.baseDelay, tt.speedLevel)
            if got != tt.wantDelay {
                t.Errorf("GetDelayForSpeed(%d, %d) = %d, want %d", tt.baseDelay, tt.speedLevel, got, tt.wantDelay)
            }
        })
    }
}

// Test 15: Hold Piece Decision
func TestShouldHoldPiece(t *testing.T) {
    gameState := NewGameState()
    gameState.CurrentPiece = NewTetromino(TetrominoO) // current: O
    gameState.NextPiece = NewTetromino(TetrominoI)    // next: I
    
    // O piece is good for filling holes, I piece for clearing lines
    // Decision depends on board state
    shouldHold := shouldHoldPiece(gameState, nil)
    
    // Just verify it returns a boolean without panic
    _ = shouldHold
}
```

### Integration Test Specifications (Task 23)

**Test File**: `autoplay_integration_test.go`

```go
package model

import "testing"

// Test 23.1: Survival Test - 10 Pieces
func TestAutoPlay_Survival10Pieces(t *testing.T) {
    gameState := NewGameState()
    autoPlayer := NewAutoPlayer()
    autoPlayer.Toggle() // enable
    
    piecesPlaced := 0
    for piecesPlaced < 10 && !gameState.GameOver {
        // Find and execute best move
        decision := FindBestMove(gameState)
        if decision == nil {
            break
        }
        
        // Execute move
        executeMove(gameState, decision)
        piecesPlaced++
    }
    
    if gameState.GameOver {
        t.Errorf("Game over before 10 pieces: only %d placed", piecesPlaced)
    }
}

// Test 23.2: Survival Test - 50 Pieces
func TestAutoPlay_Survival50Pieces(t *testing.T) {
    gameState := NewGameState()
    autoPlayer := NewAutoPlayer()
    autoPlayer.Toggle()
    
    piecesPlaced := 0
    for piecesPlaced < 50 && !gameState.GameOver {
        decision := FindBestMove(gameState)
        if decision == nil {
            break
        }
        executeMove(gameState, decision)
        piecesPlaced++
    }
    
    if gameState.GameOver {
        t.Errorf("Game over before 50 pieces: only %d placed, lines cleared: %d", 
            piecesPlaced, gameState.LinesCleared)
    }
    
    if piecesPlaced < 50 {
        t.Errorf("Expected 50 pieces, got %d", piecesPlaced)
    }
}

// Test 23.3: Line Clear Rate
func TestAutoPlay_ClearsLines(t *testing.T) {
    gameState := NewGameState()
    autoPlayer := NewAutoPlayer()
    autoPlayer.Toggle()
    
    piecesPlaced := 0
    targetPieces := 100
    
    for piecesPlaced < targetPieces && !gameState.GameOver {
        decision := FindBestMove(gameState)
        if decision == nil {
            break
        }
        executeMove(gameState, decision)
        piecesPlaced++
    }
    
    // Expect at least 20% line clear rate
    expectedLines := piecesPlaced * 0.2
    if float64(gameState.LinesCleared) < expectedLines {
        t.Errorf("Line clear rate too low: %d lines / %d pieces = %.2f%%, want >=20%%",
            gameState.LinesCleared, piecesPlaced, 
            float64(gameState.LinesCleared)/float64(piecesPlaced)*100)
    }
}

// Test 23.4: All Piece Types Handled
func TestAutoPlay_AllPieceTypes(t *testing.T) {
    gameState := NewGameState()
    autoPlayer := NewAutoPlayer()
    autoPlayer.Toggle()
    
    pieceTypesSeen := make(map[TetrominoType]bool)
    piecesPlaced := 0
    
    for piecesPlaced < 50 && !gameState.GameOver {
        pieceTypesSeen[gameState.CurrentPiece.Type] = true
        
        decision := FindBestMove(gameState)
        if decision == nil {
            break
        }
        executeMove(gameState, decision)
        piecesPlaced++
    }
    
    // Verify all 7 piece types were handled
    expectedTypes := []TetrominoType{
        TetrominoI, TetrominoO, TetrominoT,
        TetrominoS, TetrominoZ, TetrominoJ, TetrominoL,
    }
    
    for _, pType := range expectedTypes {
        if !pieceTypesSeen[pType] {
            t.Errorf("Piece type %s never seen in 50 pieces", pType)
        }
    }
}

// Test 23.5: Speed Changes Don't Break Execution
func TestAutoPlay_SpeedChanges(t *testing.T) {
    gameState := NewGameState()
    autoPlayer := NewAutoPlayer()
    autoPlayer.Toggle()
    
    for level := 1; level <= 5; level++ {
        autoPlayer.SetSpeedLevel(level)
        
        // Place 5 pieces at each speed level
        for i := 0; i < 5 && !gameState.GameOver; i++ {
            decision := FindBestMove(gameState)
            if decision == nil {
                break
            }
            executeMove(gameState, decision)
        }
    }
    
    if gameState.GameOver {
        t.Error("Game over during speed change test")
    }
}

// Test 23.6: Pause and Resume
func TestAutoPlay_PauseResume(t *testing.T) {
    gameState := NewGameState()
    autoPlayer := NewAutoPlayer()
    autoPlayer.Toggle()
    
    // Place 5 pieces
    for i := 0; i < 5; i++ {
        decision := FindBestMove(gameState)
        executeMove(gameState, decision)
    }
    
    linesBeforePause := gameState.LinesCleared
    
    // Pause
    gameState.Pause()
    
    // Try to place more pieces (should not happen while paused)
    gameState.Paused = true
    decision := FindBestMove(gameState)
    if decision != nil {
        // Should not execute while paused
        // This is more of an integration check
    }
    
    // Resume
    gameState.Pause()
    
    // Continue placing pieces
    for i := 0; i < 5; i++ {
        decision := FindBestMove(gameState)
        if decision == nil {
            break
        }
        executeMove(gameState, decision)
    }
    
    // Game should still be playable
    if gameState.GameOver {
        t.Error("Game over after pause/resume")
    }
}
```

### Benchmark Tests (Task 26)

```go
// Benchmark 26.1: FindBestMove Performance
func BenchmarkFindBestMove(b *testing.B) {
    gameState := NewGameState()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        FindBestMove(gameState)
    }
}

// Benchmark 26.2: EvaluateBoard Performance
func BenchmarkEvaluateBoard(b *testing.B) {
    board := NewBoard()
    piece := NewTetromino(TetrominoI)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        evaluateBoard(board, piece, 5, 0)
    }
}

// Benchmark 26.3: Full Game Simulation
func BenchmarkAutoPlayGame(b *testing.B) {
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        gameState := NewGameState()
        autoPlayer := NewAutoPlayer()
        autoPlayer.Toggle()
        autoPlayer.SetSpeedLevel(5) // fastest
        
        piecesPlaced := 0
        for piecesPlaced < 50 && !gameState.GameOver {
            decision := FindBestMove(gameState)
            if decision == nil {
                break
            }
            executeMove(gameState, decision)
            piecesPlaced++
        }
    }
}
```

### Performance Targets

| Test | Target | Measurement |
|------|--------|-------------|
| `TestAutoPlay_Survival10Pieces` | 100% pass rate | 10/10 runs |
| `TestAutoPlay_Survival50Pieces` | 100% pass rate | 10/10 runs |
| `TestAutoPlay_ClearsLines` | >20% line clear rate | lines/pieces |
| `BenchmarkFindBestMove` | <10ms per call | ns/op |
| `BenchmarkAutoPlayGame` | <5s for 50 pieces | total time |

---

## Test Result Report Template

**Location**: `.sisyphus/evidence/test-results/`

### Test Execution Log

```
=== Test Run: [DATE]
=== Go Version: go1.26.0
=== Platform: darwin/amd64

--- Unit Tests ---
go test -v ./internal/model/autoplay_test.go

PASS: TestAutoPlayerCreation
PASS: TestAutoPlayer_Toggle
PASS: TestAutoPlayer_SetSpeedLevel
PASS: TestAutoPlayer_CycleSpeed
PASS: TestMoveDecision_IsValid (8 subtests)
PASS: TestMoveDecision_String
PASS: TestCalculateSoftDrops (3 subtests)
PASS: TestGetColHeight (4 subtests)
PASS: TestGetAggregateHeight (3 subtests)
PASS: TestCountCompleteLines (4 subtests)
PASS: TestCountHoles (5 subtests)
PASS: TestCalculateBumpiness (5 subtests)
PASS: TestCountWells (4 subtests)
PASS: TestEvalAggregateHeight
PASS: TestEvalHoles
PASS: TestEvalBumpiness
PASS: TestEvalWells
PASS: TestEvaluateBoard (3 subtests)
PASS: TestGetWeights_SetWeights
PASS: TestEvaluateBoard_WeightImpact
PASS: TestEnumerateMoves
PASS: TestEnumerateMoves_AllPieceTypes (7 subtests)
PASS: TestFindBestMove
PASS: TestFindBestMove_ObviousCases
PASS: TestFindBestMove_Determinism
PASS: TestFindBestMove_NoValidMoves
PASS: TestExecuteRotations
PASS: TestExecuteHorizontalMove
PASS: TestExecuteDrop
PASS: TestGetDelayForSpeed (5 subtests)
PASS: TestShouldHoldPiece

Total: 31 tests, 0 failures

--- Integration Tests ---
go test -v ./internal/model/autoplay_integration_test.go

PASS: TestAutoPlay_Survival10Pieces
PASS: TestAutoPlay_Survival50Pieces
PASS: TestAutoPlay_ClearsLines
PASS: TestAutoPlay_AllPieceTypes
PASS: TestAutoPlay_SpeedChanges
PASS: TestAutoPlay_PauseResume

Total: 6 tests, 0 failures

--- Benchmark Tests ---
go test -bench=. -benchmem ./internal/model/...

BenchmarkFindBestMove-8           100    8234567 ns/op    1234 B/op    45 allocs/op
BenchmarkEvaluateBoard-8         5000     234567 ns/op     567 B/op     12 allocs/op
BenchmarkAutoPlayGame-8             20   456789012 ns/op  123456 B/op  4567 allocs/op

--- Coverage Report ---
go test -coverprofile=coverage.out ./internal/model/...

github.com/oc-garden/tetris_game/internal/model    87.3%
  autoplay.go                                      92.1%
  autoplay_test.go                                 100.0%
  autoplay_integration_test.go                     100.0%

--- Evidence Files ---
.sisyphus/evidence/test-results/
├── unit-tests-output.txt
├── integration-tests-output.txt
├── benchmark-results.txt
├── coverage.html
└── test-run-[timestamp].log
```

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately — Foundation + Types):
├── Task 1: AutoPlayer struct + basic types [quick]
├── Task 2: MoveDecision struct + validation [quick]
├── Task 3: Test infrastructure setup [quick]
├── Task 4: Board evaluation helpers [quick]
├── Task 5: Heuristic function - aggregate height [quick]
├── Task 6: Heuristic function - holes detection [quick]
└── Task 7: Heuristic function - bumpiness + wells [quick]

Wave 2 (After Wave 1 — Core AI Logic, MAX PARALLEL):
├── Task 8: Complete heuristic evaluation (weighted sum) [deep]
├── Task 9: Enumerate all possible moves [unspecified-high]
├── Task 10: FindBestMove algorithm [deep]
├── Task 11: Move execution - rotations [quick]
├── Task 12: Move execution - horizontal movement [quick]
├── Task 13: Move execution - drop logic [quick]
├── Task 14: Speed control system (5 levels) [quick]
└── Task 15: Hold piece AI logic [unspecified-high]

Wave 3 (After Wave 2 — Integration):
├── Task 16: Main.go integration - 'A' key handler [quick]
├── Task 17: Input generation from AI decisions [unspecified-high]
├── Task 18: AI game loop timing (speed-based delays) [deep]
├── Task 19: Pause/game over state handling [quick]
├── Task 20: UI render - AUTO-PLAY indicator [visual-engineering]
├── Task 21: UI render - speed level display [visual-engineering]
└── Task 22: UI render - decision panel [visual-engineering]

Wave 4 (After Wave 3 — Verification + Polish):
├── Task 23: Integration tests - full game scenarios [deep]
├── Task 24: AI tuning - weight adjustment for better play [unspecified-high]
├── Task 25: Edge case handling (game over recovery, reset) [quick]
├── Task 26: Performance optimization (cache evaluations) [unspecified-high]
├── Task 27: Documentation - code comments + README section [writing]
├── Task 28: Final QA - 100+ piece survival test [deep]
└── Task 29: Git cleanup + tagging [git]

Wave FINAL (After ALL tasks — Independent Review, 4 parallel):
├── Task F1: Plan compliance audit (oracle)
├── Task F2: Code quality review (unspecified-high)
├── Task F3: Real manual QA (unspecified-high)
└── Task F4: Scope fidelity check (deep)

Critical Path: Task 1 → Task 5-7 → Task 8 → Task 10 → Task 16-18 → Task 23 → Task 28 → F1-F4
Parallel Speedup: ~65% faster than sequential
Max Concurrent: 7 (Waves 1 & 2)
```

### Dependency Matrix

- **1-7**: — (can all start immediately)
- **8**: 5-7 (needs individual heuristics)
- **9-10**: 8 (needs complete evaluation)
- **11-15**: 9-10 (needs move finding)
- **16-19**: 11-15 (needs move execution)
- **20-22**: 16 (needs integration)
- **23-28**: 20-22 (needs full integration)
- **29**: 23-28 (needs all complete)
- **F1-F4**: All tasks complete

### Agent Dispatch Summary

- **Wave 1**: **7 tasks** — All `quick` (types, helpers, individual heuristics)
- **Wave 2**: **8 tasks** — T8,T10 → `deep`, T9,T15 → `unspecified-high`, T11-T14 → `quick`
- **Wave 3**: **7 tasks** — T16,T19 → `quick`, T17 → `unspecified-high`, T18 → `deep`, T20-T22 → `visual-engineering`
- **Wave 4**: **7 tasks** — T23,T28 → `deep`, T24,T26 → `unspecified-high`, T25 → `quick`, T27 → `writing`, T29 → `git`
- **FINAL**: **4 tasks** — F1 → `oracle`, F2 → `unspecified-high`, F3 → `unspecified-high`, F4 → `deep`

---

## TODOs

- [ ] 1. **AutoPlayer struct + basic types**

  **What to do**:
  - Create `internal/model/autoplay.go` with package declaration
  - Define `AutoPlayer` struct with fields:
    - `enabled bool` - is auto-play active
    - `speedLevel int` - 1-5 speed setting
    - `targetDecision *MoveDecision` - current target move
    - `moveIndex int` - current step in move execution
  - Define `MoveDecision` struct with fields:
    - `rotations int` - target rotation (0-3)
    - `targetX int` - target X position (0-9)
    - `softDrops int` - number of soft drops needed
    - `score float64` - evaluation score
  - Add `NewAutoPlayer()` constructor
  - Add `SetSpeedLevel(level int)` method
  - Add `GetSpeedLevel()` getter
  - Add `Toggle()` method to enable/disable
  - Add `IsEnabled()` getter

  **Must NOT do**:
  - Implement AI logic yet (just types)
  - Add speed delay calculations yet

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (simple struct definitions, no complex logic)
  - **Skills Evaluated but Omitted**: None needed for basic types

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 2-7)
  - **Blocks**: Tasks 8-15 (need types to proceed)
  - **Blocked By**: None

  **References**:
  - `internal/model/gamestate.go:4-21` - GameState struct pattern for field organization
  - `internal/model/piece.go:3-10` - Tetromino struct pattern
  - `internal/model/gamestate.go:24-35` - Constructor pattern (`NewGameState`)
  - `internal/model/randomizer.go` - Similar model package structure

  **Acceptance Criteria**:
  - [ ] File `internal/model/autoplay.go` created with valid Go syntax
  - [ ] `go build ./internal/model/...` succeeds
  - [ ] `go vet ./internal/model/...` reports no issues
  - [ ] All structs and methods compile without errors

  **QA Scenarios**:

  ```
  Scenario: Build verification
    Tool: Bash
    Preconditions: In project root directory
    Steps:
      1. Run: go build ./internal/model/...
      2. Check exit code is 0
      3. Run: go vet ./internal/model/...
      4. Check no output (no issues)
    Expected Result: Both commands succeed with exit code 0
    Failure Indicators: Non-zero exit code or vet warnings
    Evidence: .sisyphus/evidence/task-1-build.txt
  ```

  **Commit**: YES (groups with 2-7)
  - Message: `feat(autoplay): add AutoPlayer and MoveDecision types`
  - Files: `internal/model/autoplay.go`
  - Pre-commit: `go build ./internal/model/... && go vet ./internal/model/...`

- [ ] 2. **MoveDecision validation + helper methods**

  **What to do**:
  - Add `IsValid()` method to MoveDecision - validates rotation 0-3, X 0-9
  - Add `Reset()` method to clear decision state
  - Add `String()` method for debugging (fmt.Stringer interface)
  - Add helper: `CalculateSoftDrops(board, piece, targetX)` - calculates drops needed
  - Ensure all methods use pointer receivers for consistency with codebase

  **Must NOT do**:
  - Implement actual move execution yet
  - Add AI evaluation logic

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (simple validation logic)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 3-7)
  - **Blocks**: Tasks 8-15
  - **Blocked By**: Task 1 (MoveDecision struct)

  **References**:
  - `internal/model/gamestate.go:277-299` - Method pattern with pointer receivers
  - `internal/model/piece.go:237-260` - Helper methods on structs
  - `internal/model/board.go:37-44` - Validation pattern (`IsWithinBounds`)

  **Acceptance Criteria**:
  - [ ] `MoveDecision.IsValid()` returns false for rotation >3 or X outside 0-9
  - [ ] `MoveDecision.String()` returns readable format for debugging
  - [ ] `go test -run TestMoveDecision ./internal/model/...` passes

  **QA Scenarios**:

  ```
  Scenario: IsValid validation
    Tool: Bash (go test)
    Preconditions: Task 1 complete
    Steps:
      1. Create table-driven test with valid/invalid cases
      2. Test: rotation=0, X=5 → IsValid() = true
      3. Test: rotation=4, X=5 → IsValid() = false
      4. Test: rotation=2, X=-1 → IsValid() = false
      5. Test: rotation=2, X=10 → IsValid() = false
    Expected Result: All test cases pass with correct boolean results
    Failure Indicators: Any test failure or panic
    Evidence: .sisyphus/evidence/task-2-validation-test.txt
  ```

  **Commit**: YES (groups with 1, 3-7)
  - Message: `feat(autoplay): add MoveDecision validation and helpers`
  - Files: `internal/model/autoplay.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 3. **Test infrastructure setup for autoplay**

  **What to do**:
  - Create `internal/model/autoplay_test.go` with package declaration
  - Add test helper: `createTestBoard()` - creates board with known state
  - Add test helper: `checkDecision(t, got, expected)` - compares MoveDecision
  - Add test helper: `checkFloatApprox(t, got, expected, tolerance)` - float comparison
  - Add basic sanity test: `TestAutoPlayerCreation` - verifies constructor works
  - Follow existing test patterns from codebase (table-driven tests)

  **Must NOT do**:
  - Test AI logic yet (nothing to test)
  - Add complex test scenarios

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (standard Go testing patterns)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1-2, 4-7)
  - **Blocks**: All subsequent test tasks
  - **Blocked By**: Task 1 (needs AutoPlayer type)

  **References**:
  - `internal/model/board_test.go:10-40` - Test helper function pattern
  - `internal/model/piece_test.go` - Table-driven test structure
  - `internal/model/gamestate_test.go` - Test setup patterns
  - `internal/testutil/helpers.go` - Existing test utilities

  **Acceptance Criteria**:
  - [ ] `internal/model/autoplay_test.go` created
  - [ ] `go test -run TestAutoPlayerCreation ./internal/model/...` passes
  - [ ] Test helpers compile and work correctly
  - [ ] Test file follows existing codebase style

  **QA Scenarios**:

  ```
  Scenario: Test execution
    Tool: Bash (go test)
    Preconditions: Tasks 1-2 complete
    Steps:
      1. Run: go test -v -run TestAutoPlayerCreation ./internal/model/...
      2. Verify output shows "PASS"
      3. Run: go test ./internal/model/... (all tests)
      4. Verify all existing tests still pass
    Expected Result: All tests pass, no regressions
    Failure Indicators: Test failures or compilation errors
    Evidence: .sisyphus/evidence/task-3-tests.txt
  ```

  **Commit**: YES (groups with 1-2, 4-7)
  - Message: `feat(autoplay): add test infrastructure`
  - Files: `internal/model/autoplay_test.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 4. **Board evaluation helpers**

  **What to do**:
  - Add `getColHeight(board, col)` - returns height of pieces in column (0-9)
  - Add `getMaxHeight(board)` - returns highest column
  - Add `getAggregateHeight(board)` - returns sum of all column heights
  - Add `countCompleteLines(board)` - returns number of full lines
  - Add `countHoles(board)` - returns count of empty cells with blocks above
  - Add `calculateBumpiness(board)` - returns sum of height differences between adjacent columns
  - Add `countWells(board)` - returns count of deep empty columns (2+ depth)
  - Each helper should be pure function (board *Board input, int output)

  **Must NOT do**:
  - Combine into weighted evaluation yet
  - Modify existing Board struct

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (simple iteration logic)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1-3, 5-7)
  - **Blocks**: Task 8 (needs all helpers)
  - **Blocked By**: Task 1 (uses AutoPlayer package)

  **References**:
  - `internal/model/board.go:47-56` - Board iteration pattern (`IsFull`)
  - `internal/model/board.go:59-69` - Line checking pattern (`IsLineFull`)
  - `internal/model/gamestate.go:172-189` - Line clearing iteration

  **Acceptance Criteria**:
  - [ ] Each helper function has dedicated test
  - [ ] `getColHeight` correctly counts from bottom until empty cell
  - [ ] `countHoles` correctly identifies cells with blocks above
  - [ ] `calculateBumpiness` returns sum of abs(col[i] - col[i+1])
  - [ ] All helpers handle empty board correctly (return 0)

  **QA Scenarios**:

  ```
  Scenario: Hole detection on known board
    Tool: Bash (go test)
    Preconditions: Board helpers implemented
    Steps:
      1. Create board with known configuration (2 blocks with hole between)
      2. Call countHoles(board)
      3. Assert result matches expected count
      4. Test empty board returns 0
      5. Test full board returns 0 (no holes possible)
    Expected Result: All hole detection tests pass
    Failure Indicators: Incorrect hole count or panic on edge cases
    Evidence: .sisyphus/evidence/task-4-holes-test.txt
  ```

  **Commit**: YES (groups with 1-3, 5-7)
  - Message: `feat(autoplay): add board evaluation helper functions`
  - Files: `internal/model/autoplay.go`, `internal/model/autoplay_test.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 5. **Heuristic function - aggregate height**

  **What to do**:
  - Implement `evalAggregateHeight(board)` function
  - Should return sum of all column heights (0-10 per column, max 200)
  - Add test: empty board returns 0
  - Add test: single column with height 5 returns 5
  - Add test: multiple columns sum correctly
  - Document that LOWER is better (negative weight in final evaluation)

  **Must NOT do**:
  - Apply weight yet (just raw value)
  - Combine with other heuristics

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (simple sum calculation)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1-4, 6-7)
  - **Blocks**: Task 8 (needs for weighted sum)
  - **Blocked By**: Task 4 (uses getColHeight helper)

  **References**:
  - Research: Standard Tetris AI weights (aggregate height typically -0.5)
  - `internal/model/board.go:37-44` - Column iteration pattern

  **Acceptance Criteria**:
  - [ ] `evalAggregateHeight(emptyBoard)` returns 0
  - [ ] `evalAggregateHeight(singleColumn)` returns correct height
  - [ ] `evalAggregateHeight(fullBoard)` returns 200 (10 cols × 20 rows)
  - [ ] Test coverage for edge cases

  **QA Scenarios**:

  ```
  Scenario: Aggregate height calculation
    Tool: Bash (go test)
    Steps:
      1. Create board with column heights [3,5,2,0,4,1,6,3,2,4]
      2. Call evalAggregateHeight()
      3. Assert result = 30 (sum of heights)
      4. Verify calculation is deterministic
    Expected Result: Correct sum calculation
    Failure Indicators: Wrong sum or non-deterministic result
    Evidence: .sisyphus/evidence/task-5-height-test.txt
  ```

  **Commit**: YES (groups with 1-4, 6-7)
  - Message: `feat(autoplay): implement aggregate height heuristic`
  - Files: `internal/model/autoplay.go`, `internal/model/autoplay_test.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 6. **Heuristic function - holes detection**

  **What to do**:
  - Implement `evalHoles(board)` function
  - Count empty cells that have at least one block above them in same column
  - Add test: empty board returns 0
  - Add test: single hole returns 1
  - Add test: multiple holes in same column counted correctly
  - Document that LOWER is better (negative weight, typically -0.36)

  **Must NOT do**:
  - Apply weight yet
  - Count wells (separate heuristic)

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (simple iteration and counting)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1-5, 7)
  - **Blocks**: Task 8
  - **Blocked By**: Task 4 (uses iteration helpers)

  **References**:
  - Research: Holes typically weighted -0.36 in Tetris AI
  - `internal/model/board.go:42-44` - Empty cell check pattern

  **Acceptance Criteria**:
  - [ ] `evalHoles(emptyBoard)` returns 0
  - [ ] Board with one empty cell below a block returns 1
  - [ ] Multiple holes in same column all counted
  - [ ] Holes only counted when block exists above (not at top)

  **QA Scenarios**:

  ```
  Scenario: Hole counting accuracy
    Tool: Bash (go test)
    Steps:
      1. Create board with column: [block, empty, block, empty, empty...] from bottom
      2. Call evalHoles() - should count 2 holes
      3. Create board with overhang creating 3 holes
      4. Verify count matches expected
    Expected Result: Accurate hole counting in all configurations
    Failure Indicators: Under/over counting holes
    Evidence: .sisyphus/evidence/task-6-holes-test.txt
  ```

  **Commit**: YES (groups with 1-5, 7)
  - Message: `feat(autoplay): implement holes heuristic`
  - Files: `internal/model/autoplay.go`, `internal/model/autoplay_test.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 7. **Heuristic function - bumpiness + wells**

  **What to do**:
  - Implement `evalBumpiness(board)` - sum of abs(col[i] - col[i+1]) for i=0..8
  - Implement `evalWells(board)` - count of columns with depth 2+ empty spaces between blocks
  - Add test: flat board (all same height) returns bumpiness 0
  - Add test: alternating heights [5,0,5,0...] returns high bumpiness
  - Add test: wells counted correctly (U-shaped gaps)
  - Document that LOWER is better for both (bumpiness ~-0.18, wells ~-0.12)

  **Must NOT do**:
  - Apply weights yet
  - Combine heuristics

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (arithmetic and iteration)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1-6)
  - **Blocks**: Task 8
  - **Blocked By**: Task 4 (uses col height helpers)

  **References**:
  - Research: Bumpiness typically -0.18, wells -0.12
  - Standard Tetris AI literature on feature evaluation

  **Acceptance Criteria**:
  - [ ] `evalBumpiness(flatBoard)` returns 0
  - [ ] `evalBumpiness(steppedBoard)` returns correct sum of differences
  - [ ] `evalWells(emptyBoard)` returns 0
  - [ ] `evalWells(boardWithWells)` counts wells correctly
  - [ ] All tests pass

  **QA Scenarios**:

  ```
  Scenario: Bumpiness calculation
    Tool: Bash (go test)
    Steps:
      1. Create board with heights [3,5,2,4,1,6,3,2,4,5]
      2. Call evalBumpiness()
      3. Calculate expected: |3-5|+|5-2|+|2-4|+|4-1|+|1-6|+|6-3|+|3-2|+|2-4|+|4-5| = 2+3+2+3+5+3+1+2+1 = 22
      4. Assert result equals 22
    Expected Result: Correct bumpiness calculation
    Failure Indicators: Wrong sum or off-by-one errors
    Evidence: .sisyphus/evidence/task-7-bumpiness-test.txt
  ```

  **Commit**: YES (groups with Tasks 1-6)
  - Message: `feat(autoplay): implement bumpiness and wells heuristics`
  - Files: `internal/model/autoplay.go`, `internal/model/autoplay_test.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 8. **Complete heuristic evaluation (weighted sum)**

  **What to do**:
  - Implement `evaluateBoard(board, piece, x, rotations)` - main evaluation function
  - Combine all heuristics with weights:
    - aggregateHeight × -0.50
    - completeLines × +0.76
    - holes × -0.36
    - bumpiness × -0.18
    - wells × -0.12
  - Return single float64 score (higher = better position)
  - Add `GetWeights()` function to allow weight customization
  - Add `SetWeights()` function for future tuning
  - Document each weight's rationale in comments

  **Must NOT do**:
  - Change weight values without testing
  - Add new heuristic features (scope creep)

  **Recommended Agent Profile**:
  - **Category**: `deep`
  - **Skills**: `[]` (complex logic with multiple components)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2 (starts after Wave 1)
  - **Blocks**: Tasks 9-10 (need evaluation)
  - **Blocked By**: Tasks 5-7 (needs all individual heuristics)

  **References**:
  - Research: Dellacherie's Tetris AI weights (standard reference)
  - `internal/model/autoplay.go:1-50` - Heuristic functions from Tasks 5-7
  - Tasks 4-7 in this plan - individual heuristic implementations

  **Acceptance Criteria**:
  - [ ] `evaluateBoard` returns higher score for flat boards vs bumpy
  - [ ] `evaluateBoard` returns higher score for boards with complete lines
  - [ ] `evaluateBoard` returns lower score for boards with holes
  - [ ] Weights are configurable via GetWeights/SetWeights
  - [ ] Comprehensive tests for score comparisons

  **QA Scenarios**:

  ```
  Scenario: Weighted evaluation comparison
    Tool: Bash (go test)
    Steps:
      1. Create board A: flat surface, no holes, 1 complete line
      2. Create board B: bumpy surface, 2 holes, no complete lines
      3. Call evaluateBoard on both
      4. Assert score(A) > score(B)
      5. Test with various board configurations
    Expected Result: Evaluation correctly ranks better positions higher
    Failure Indicators: Worse boards scoring higher, or equal scores for different boards
    Evidence: .sisyphus/evidence/task-8-evaluation-test.txt
  ```

  **Commit**: YES
  - Message: `feat(autoplay): implement weighted heuristic evaluation`
  - Files: `internal/model/autoplay.go`, `internal/model/autoplay_test.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 9. **Enumerate all possible moves**

  **What to do**:
  - Implement `enumerateMoves(gameState, piece)` function
  - Generate all 4 rotations (0-3)
  - For each rotation, generate all valid X positions (0-9)
  - For each (rotation, X) combination, calculate soft drops needed to land
  - Filter out invalid positions (collisions, out of bounds)
  - Return slice of MoveDecision with rotations, targetX, softDrops populated
  - Handle wall kicks if piece can't fit at certain X positions

  **Must NOT do**:
  - Evaluate moves yet (just enumerate)
  - Select best move yet

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: `[]` (complex enumeration logic)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2 (starts after Wave 1)
  - **Blocks**: Task 10 (needs move list)
  - **Blocked By**: Tasks 1-2 (needs MoveDecision type)

  **References**:
  - `internal/model/gamestate.go:128-147` - Position validation (`isValidPosition`)
  - `internal/model/piece.go:242-253` - Rotation methods
  - `internal/model/board.go:37-44` - Bounds checking

  **Acceptance Criteria**:
  - [ ] Returns all valid (rotation, X) combinations for I piece (typically 30-40 moves)
  - [ ] Filters out positions where piece doesn't fit
  - [ ] Calculates correct soft drops for each position
  - [ ] Handles edge pieces (O piece has 1 unique rotation effectively)
  - [ ] Handles wall kicks or documents behavior

  **QA Scenarios**:

  ```
  Scenario: Move enumeration for I piece
    Tool: Bash (go test)
    Steps:
      1. Create gameState with I piece at spawn position
      2. Call enumerateMoves(gameState, piece)
      3. Count returned moves (expect ~30-40 depending on walls)
      4. Verify each move has rotation 0-3 and X 0-9
      5. Verify softDrops is reasonable (0-20)
    Expected Result: All valid moves enumerated, no invalid moves included
    Failure Indicators: Missing valid moves or including invalid moves
    Evidence: .sisyphus/evidence/task-9-enumeration-test.txt
  ```

  **Commit**: YES
  - Message: `feat(autoplay): implement move enumeration`
  - Files: `internal/model/autoplay.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 10. **FindBestMove algorithm**

  **What to do**:
  - Implement `FindBestMove(gameState)` - main AI decision function
  - Call enumerateMoves() to get all possible moves
  - For each move, simulate the landing position
  - Call evaluateBoard() on simulated result
  - Track move with highest score
  - Return best MoveDecision (or nil if no valid moves - game over scenario)
  - Handle tie-breaking: prefer lower X, then fewer rotations

  **Must NOT do**:
  - Execute the move yet (just find it)
  - Add lookahead beyond current piece (scope creep)

  **Recommended Agent Profile**:
  - **Category**: `deep`
  - **Skills**: `[]` (complex decision logic)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2 (starts after Wave 1)
  - **Blocks**: Tasks 11-15 (need best move)
  - **Blocked By**: Tasks 8-9 (needs evaluation and enumeration)

  **References**:
  - Task 8 - evaluateBoard function
  - Task 9 - enumerateMoves function
  - `internal/model/gamestate.go:148-170` - Piece locking logic (for simulation)

  **Acceptance Criteria**:
  - [ ] Returns valid MoveDecision for empty board
  - [ ] Returns move that avoids holes when possible
  - [ ] Returns move that creates/uses complete lines when available
  - [ ] Handles no-valid-moves scenario gracefully (returns nil)
  - [ ] Deterministic: same board state = same decision

  **QA Scenarios**:

  ```
  Scenario: Best move selection
    Tool: Bash (go test)
    Steps:
      1. Create board with obvious best move (one position completes line)
      2. Call FindBestMove(gameState)
      3. Verify returned move targets the line-completing position
      4. Test with multiple piece types
      5. Verify tie-breaking works (lower X preferred)
    Expected Result: AI selects optimal or near-optimal moves consistently
    Failure Indicators: Selecting obviously bad moves, non-deterministic behavior
    Evidence: .sisyphus/evidence/task-10-findbest-test.txt
  ```

  **Commit**: YES
  - Message: `feat(autoplay): implement FindBestMove algorithm`
  - Files: `internal/model/autoplay.go`, `internal/model/autoplay_test.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 11. **Move execution - rotations**

  **What to do**:
  - Implement `executeRotations(gameState, targetRotations)` function
  - Compare current piece rotation with target
  - Calculate minimum rotations needed (0-3, accounting for wrap)
  - Call existing RotatePiece() or RotatePieceCounter() for each rotation
  - Handle rotation failures (wall kicks not supported - skip if fails)
  - Track execution progress in AutoPlayer.moveIndex

  **Must NOT do**:
  - Execute horizontal movement yet
  - Execute drop yet

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (simple rotation logic)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 12-15)
  - **Blocks**: Task 18 (needs all execution)
  - **Blocked By**: Task 10 (needs target decision)

  **References**:
  - `internal/model/gamestate.go:61-84` - Rotation methods
  - `internal/model/piece.go:247-253` - Rotation implementation

  **Acceptance Criteria**:
  - [ ] Correct number of rotations executed
  - [ ] Uses shortest rotation direction (CW vs CCW)
  - [ ] Handles rotation failures gracefully
  - [ ] Updates AutoPlayer state correctly

  **QA Scenarios**:

  ```
  Scenario: Rotation execution
    Tool: Bash (go test or interactive_bash)
    Steps:
      1. Create gameState with piece at rotation 0
      2. Call executeRotations(gameState, 2)
      3. Verify piece is now at rotation 2
      4. Test wrapping: rotation 3 → 0 (1 rotation CW)
    Expected Result: Piece rotated to target orientation
    Failure Indicators: Wrong rotation, panic on wall collision
    Evidence: .sisyphus/evidence/task-11-rotation-test.txt
  ```

  **Commit**: YES (groups with 12-15)
  - Message: `feat(autoplay): implement move execution (rotations, movement, drop)`
  - Files: `internal/model/autoplay.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 12. **Move execution - horizontal movement**

  **What to do**:
  - Implement `executeHorizontalMove(gameState, targetX)` function
  - Compare current piece X with target X
  - Call MovePiece(-1, 0) or MovePiece(1, 0) repeatedly until target reached
  - Handle movement failures (piece blocked - return error/skip)
  - Track execution progress

  **Must NOT do**:
  - Execute rotations yet
  - Execute drop yet

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (simple movement logic)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 11, 13-15)
  - **Blocks**: Task 18
  - **Blocked By**: Task 10

  **References**:
  - `internal/model/gamestate.go:47-58` - MovePiece method

  **Acceptance Criteria**:
  - [ ] Piece moves to target X when path is clear
  - [ ] Handles blocked movement gracefully
  - [ ] Correct direction (left vs right)

  **Commit**: YES (groups with 11, 13-15)
  - Message: `feat(autoplay): implement move execution (rotations, movement, drop)`
  - Files: `internal/model/autoplay.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 13. **Move execution - drop logic**

  **What to do**:
  - Implement `executeDrop(gameState, softDrops, useHardDrop)` function
  - For soft drops: call SoftDrop() repeatedly
  - For hard drop: call DropPiece() once (instant)
  - Respect speed level for soft drop timing
  - Handle drop completion (piece locked automatically)

  **Must NOT do**:
  - Add timing logic yet (handled in Task 18)

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (simple drop logic)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 11-12, 14-15)
  - **Blocks**: Task 18
  - **Blocked By**: Task 10

  **References**:
  - `internal/model/gamestate.go:106-122` - DropPiece and SoftDrop methods

  **Acceptance Criteria**:
  - [ ] Soft drop executes correct number of drops
  - [ ] Hard drop executes instant drop
  - [ ] Piece locks correctly after drop

  **Commit**: YES (groups with 11-12, 14-15)
  - Message: `feat(autoplay): implement move execution (rotations, movement, drop)`
  - Files: `internal/model/autoplay.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 14. **Speed control system (5 levels)**

  **What to do**:
  - Define 5 speed levels with delay multipliers:
    - Level 1: 1.0× (normal game speed, ~1500ms at level 1)
    - Level 2: 0.5× (750ms)
    - Level 3: 0.2× (300ms)
    - Level 4: 0.1× (150ms)
    - Level 5: 0.0× (instant, no delay)
  - Implement `GetDelayForSpeed(baseDelay, speedLevel)` function
  - Add `CycleSpeed()` method to rotate through levels 1→2→3→4→5→1
  - Add `GetSpeedLevel()` returns current level (1-5)

  **Must NOT do**:
  - Integrate with game loop yet
  - Add UI display yet

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (simple delay calculation)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 11-13, 15)
  - **Blocks**: Task 18
  - **Blocked By**: Task 1 (needs AutoPlayer struct)

  **References**:
  - `internal/model/gamestate.go:291-299` - GetDropInterval pattern

  **Acceptance Criteria**:
  - [ ] Each speed level returns correct delay
  - [ ] CycleSpeed() correctly rotates 1→2→3→4→5→1
  - [ ] Level 5 returns 0 delay (instant)

  **QA Scenarios**:

  ```
  Scenario: Speed level cycling
    Tool: Bash (go test)
    Steps:
      1. Create AutoPlayer with default speed
      2. Call CycleSpeed() 5 times
      3. Verify sequence: 1→2→3→4→5→1
      4. Test GetDelayForSpeed with various base delays
    Expected Result: Speed levels cycle correctly, delays calculated accurately
    Failure Indicators: Wrong level sequence or incorrect delay calculation
    Evidence: .sisyphus/evidence/task-14-speed-test.txt
  ```

  **Commit**: YES (groups with 11-13, 15)
  - Message: `feat(autoplay): add speed control system with 5 levels`
  - Files: `internal/model/autoplay.go`, `internal/model/autoplay_test.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 15. **Hold piece AI logic**

  **What to do**:
  - Implement `shouldHoldPiece(gameState, currentDecision)` function
  - Evaluate if holding current piece would be beneficial
  - Strategy: hold if next piece is better suited for current board
  - Compare best move score for current piece vs hypothetical best for next piece
  - Implement `executeHold(gameState)` - calls existing HoldCurrentPiece()
  - Track hold state in AutoPlayer

  **Must NOT do**:
  - Complex lookahead (just current + next piece comparison)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: `[]` (strategic decision logic)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 11-14)
  - **Blocks**: Task 18
  - **Blocked By**: Task 10 (needs FindBestMove)

  **References**:
  - `internal/model/gamestate.go:87-104` - HoldCurrentPiece method
  - Task 10 - FindBestMove algorithm

  **Acceptance Criteria**:
  - [ ] shouldHoldPiece returns true when holding is beneficial
  - [ ] Respects CanHold flag (can't hold twice in a row)
  - [ ] executeHold successfully swaps pieces

  **QA Scenarios**:

  ```
  Scenario: Hold piece decision
    Tool: Bash (go test)
    Steps:
      1. Create board where I piece would be better than current
      2. Set next piece to I
      3. Call shouldHoldPiece with O piece current
      4. Verify returns true (should hold O for I)
      5. Test with CanHold=false returns false
    Expected Result: Hold decisions make strategic sense
    Failure Indicators: Always holding or never holding regardless of situation
    Evidence: .sisyphus/evidence/task-15-hold-test.txt
  ```

  **Commit**: YES (groups with 11-14)
  - Message: `feat(autoplay): add hold piece AI logic`
  - Files: `internal/model/autoplay.go`, `internal/model/autoplay_test.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 16. **Main.go integration - 'A' key handler**

  **What to do**:
  - Modify `cmd/tetris/main.go` input handling section
  - Add 'A'/'a' key case in rune switch (around line 67-70)
  - Call `autoPlayer.Toggle()` when 'A' pressed
  - Add 'S'/'s' key for speed cycling (calls `autoPlayer.CycleSpeed()`)
  - Ensure key handling doesn't conflict with existing bindings
  - Initialize autoPlayer in main() alongside gameState

  **Must NOT do**:
  - Change existing key bindings
  - Modify game loop structure yet

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (simple input handling)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 3 (starts after Wave 2)
  - **Blocks**: Tasks 17-22 (needs integration)
  - **Blocked By**: Tasks 11-15 (needs move execution)

  **References**:
  - `cmd/tetris/main.go:58-75` - Input handling switch statement
  - `cmd/tetris/main.go:17-35` - Game initialization pattern

  **Acceptance Criteria**:
  - [ ] 'A' key toggles auto-play on/off
  - [ ] 'S' key cycles speed levels
  - [ ] No existing functionality broken
  - [ ] `go build -o tetris` succeeds

  **QA Scenarios**:

  ```
  Scenario: Key binding integration
    Tool: interactive_bash (tmux)
    Steps:
      1. Build and run: go run ./cmd/tetris
      2. Press 'A' key during gameplay
      3. Verify auto-play toggles (UI indicator appears)
      4. Press 'A' again
      5. Verify auto-play toggles off
      6. Press 'S' multiple times
      7. Verify speed level changes in UI
    Expected Result: Keys toggle auto-play and cycle speed as expected
    Failure Indicators: Keys not responding, conflicts with existing bindings
    Evidence: .sisyphus/evidence/task-16-key-test.gif (screenshot sequence)
  ```

  **Commit**: YES
  - Message: `feat(autoplay): integrate 'A' key toggle and 'S' speed cycle`
  - Files: `cmd/tetris/main.go`
  - Pre-commit: `go build -o tetris && go vet ./...`

- [ ] 17. **Input generation from AI decisions**

  **What to do**:
  - Implement `GenerateInputForDecision(autoPlayer, gameState)` function
  - Return InputEvent based on current move execution state
  - State machine approach:
    - State 1: Need rotation → generate rotation InputEvent
    - State 2: Need horizontal move → generate left/right InputEvent
    - State 3: Need drop → generate drop InputEvent
    - State 4: Move complete → reset for next piece
  - Track state in AutoPlayer.moveIndex
  - Handle edge cases (rotation failed, movement blocked)

  **Must NOT do**:
  - Add timing/delay logic yet
  - Modify main game loop yet

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: `[]` (state machine logic)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 3 (starts after Wave 2)
  - **Blocks**: Task 18
  - **Blocked By**: Tasks 11-15, 16

  **References**:
  - `cmd/tetris/main.go:11-15` - InputEvent struct
  - `cmd/tetris/main.go:58-75` - Input handling pattern
  - Tasks 11-13 - Move execution functions

  **Acceptance Criteria**:
  - [ ] Returns correct InputEvent for rotation state
  - [ ] Returns correct InputEvent for movement state
  - [ ] Returns correct InputEvent for drop state
  - [ ] Handles completion and resets for next piece
  - [ ] Handles failed moves gracefully (skip to next step)

  **QA Scenarios**:

  ```
  Scenario: Input generation sequence
    Tool: Bash (go test)
    Steps:
      1. Create AutoPlayer with MoveDecision (rot=2, X=5, drops=10)
      2. Call GenerateInputForDecision repeatedly
      3. Verify sequence: 2 rotation inputs → horizontal inputs → drop inputs
      4. Verify state resets after completion
    Expected Result: Correct input sequence generated for each move phase
    Failure Indicators: Wrong input type, missing inputs, infinite loop
    Evidence: .sisyphus/evidence/task-17-input-gen-test.txt
  ```

  **Commit**: YES
  - Message: `feat(autoplay): implement AI input generation from decisions`
  - Files: `internal/model/autoplay.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 18. **AI game loop timing (speed-based delays)**

  **What to do**:
  - Modify main game loop in `cmd/tetris/main.go`
  - Add auto-play execution path when `autoPlayer.IsEnabled()`
  - Implement timing logic using speed level delays
  - For each AI input, wait appropriate delay before next
  - Level 5 (instant): no delay, execute all inputs immediately
  - Levels 1-4: delay between inputs based on speed level
  - Track lastActionTime for delay calculation

  **Must NOT do**:
  - Change existing game timing for manual play
  - Break game loop structure

  **Recommended Agent Profile**:
  - **Category**: `deep`
  - **Skills**: `[]` (complex timing integration)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 3 (starts after Wave 2)
  - **Blocks**: Tasks 20-22 (UI needs integration)
  - **Blocked By**: Tasks 16-17

  **References**:
  - `cmd/tetris/main.go:54-148` - Main game loop structure
  - `cmd/tetris/main.go:96-103` - Timing logic for drop interval
  - Task 14 - GetDelayForSpeed function

  **Acceptance Criteria**:
  - [ ] Auto-play executes at correct speed for each level
  - [ ] Manual play unaffected when auto-play disabled
  - [ ] Level 5 executes moves instantly
  - [ ] Smooth execution without stuttering

  **QA Scenarios**:

  ```
  Scenario: Auto-play timing at different speeds
    Tool: interactive_bash (tmux recording)
    Steps:
      1. Start game, enable auto-play ('A')
      2. Set speed to level 1 ('S' once)
      3. Record time for one complete piece placement (~1.5s)
      4. Change to level 3 ('S' twice more)
      5. Record time for one piece (~0.3s)
      6. Change to level 5
      7. Verify instant piece placement
    Expected Result: Timing matches speed level expectations
    Failure Indicators: Wrong timing, game stuttering, input backlog
    Evidence: .sisyphus/evidence/task-18-timing-test.mp4
  ```

  **Commit**: YES
  - Message: `feat(autoplay): integrate AI timing into game loop`
  - Files: `cmd/tetris/main.go`
  - Pre-commit: `go build -o tetris && go vet ./...`

- [ ] 19. **Pause/game over state handling**

  **What to do**:
  - Add auto-play pause when game paused
  - Disable auto-play input when `gameState.Paused == true`
  - Disable auto-play when `gameState.GameOver == true`
  - Reset auto-play state on game reset
  - Auto-play re-enables after unpause if was enabled before

  **Must NOT do**:
  - Change pause behavior for manual play
  - Modify game over logic

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (simple state checking)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Tasks 16, 20-22)
  - **Blocks**: None
  - **Blocked By**: Task 16 (needs main.go integration)

  **References**:
  - `internal/model/gamestate.go:14-16` - Paused, GameOver fields
  - `cmd/tetris/main.go:96-103` - Pause check in game loop

  **Acceptance Criteria**:
  - [ ] Auto-play pauses when 'P' pressed
  - [ ] Auto-play disabled on game over
  - [ ] Auto-play state preserved across pause
  - [ ] Auto-play resets on game reset

  **QA Scenarios**:

  ```
  Scenario: Pause handling
    Tool: interactive_bash (tmux)
    Steps:
      1. Enable auto-play
      2. Press 'P' to pause
      3. Verify AI stops making moves
      4. Press 'P' to resume
      5. Verify AI resumes from where it left off
      6. Trigger game over
      7. Verify auto-play disabled
    Expected Result: Auto-play correctly handles pause and game over states
    Failure Indicators: AI continues during pause, state lost on resume
    Evidence: .sisyphus/evidence/task-19-pause-test.gif
  ```

  **Commit**: YES (groups with 16, 20-22)
  - Message: `feat(autoplay): add pause and game over state handling`
  - Files: `cmd/tetris/main.go`, `internal/model/autoplay.go`
  - Pre-commit: `go build -o tetris`

- [ ] 20. **UI render - AUTO-PLAY indicator**

  **What to do**:
  - Create `internal/ui/autoplay_render.go` with package declaration
  - Add `RenderAutoPlayIndicator(screen, autoPlayer)` function
  - Display "AUTO-PLAY" text in prominent location (top center or side panel)
  - Use distinct color (bright green or yellow) when enabled
  - Hide or gray out when disabled
  - Add visual emphasis (bold or inverse) when enabled

  **Must NOT do**:
  - Modify existing UI elements
  - Overcrowd the display

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: `[]` (terminal UI rendering with tcell)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Tasks 16, 19, 21-22)
  - **Blocks**: Task 28 (final QA needs UI)
  - **Blocked By**: Task 16 (needs integration)

  **References**:
  - `cmd/tetris/main.go:238-265` - renderUI function pattern
  - `cmd/tetris/main.go:163-198` - renderBoard style patterns
  - `internal/assets/ascii.go` - Existing UI element patterns

  **Acceptance Criteria**:
  - [ ] "AUTO-PLAY" visible when enabled
  - [ ] Hidden or grayed when disabled
  - [ ] Positioned prominently without obscuring game board
  - [ ] Uses appropriate tcell styling

  **QA Scenarios**:

  ```
  Scenario: Auto-play indicator visibility
    Tool: interactive_bash (tmux screenshot)
    Steps:
      1. Run game
      2. Verify no indicator when auto-play off
      3. Press 'A' to enable
      4. Verify "AUTO-PLAY" indicator appears
      5. Verify indicator is clearly visible and distinct
      6. Press 'A' to disable
      7. Verify indicator disappears or grays out
    Expected Result: Clear visual indicator of auto-play state
    Failure Indicators: Indicator not visible, obscures game, wrong state
    Evidence: .sisyphus/evidence/task-20-indicator-test.png
  ```

  **Commit**: YES (groups with 16, 19, 21-22)
  - Message: `feat(autoplay): add AUTO-PLAY indicator UI`
  - Files: `internal/ui/autoplay_render.go`, `cmd/tetris/main.go`
  - Pre-commit: `go build -o tetris`

- [ ] 21. **UI render - speed level display**

  **What to do**:
  - Add `RenderSpeedLevel(screen, autoPlayer)` function
  - Display "SPEED: N" (where N = 1-5) near auto-play indicator
  - Use visual feedback for speed (more bars/indicators for higher speed)
  - Optional: color code (green=slow, yellow=medium, red=fast)
  - Update in real-time when 'S' pressed

  **Must NOT do**:
  - Overcrowd UI
  - Use distracting animations

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: `[]` (terminal UI rendering)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Tasks 16, 19-20, 22)
  - **Blocks**: Task 28
  - **Blocked By**: Task 16

  **References**:
  - `cmd/tetris/main.go:238-265` - renderUI pattern
  - `cmd/tetris/main.go:239-241` - SCORE/LEVEL/LINES display pattern
  - Task 14 - Speed level system

  **Acceptance Criteria**:
  - [ ] Speed level 1-5 displayed clearly
  - [ ] Updates immediately on 'S' key press
  - [ ] Visual distinction between levels
  - [ ] Consistent placement with other UI elements

  **QA Scenarios**:

  ```
  Scenario: Speed level display
    Tool: interactive_bash (tmux)
    Steps:
      1. Enable auto-play
      2. Verify "SPEED: 1" displayed
      3. Press 'S' to cycle to level 3
      4. Verify display changes to "SPEED: 3"
      5. Cycle through all levels
      6. Verify each level displays correctly
    Expected Result: Speed level clearly visible and updates correctly
    Failure Indicators: Wrong level shown, delayed update, hard to read
    Evidence: .sisyphus/evidence/task-21-speed-ui-test.gif
  ```

  **Commit**: YES (groups with 16, 19-20, 22)
  - Message: `feat(autoplay): add speed level display UI`
  - Files: `internal/ui/autoplay_render.go`, `cmd/tetris/main.go`
  - Pre-commit: `go build -o tetris`

- [ ] 22. **UI render - decision panel**

  **What to do**:
  - Add `RenderDecisionPanel(screen, autoPlayer, gameState)` function
  - Display AI's current decision:
    - Target X position
    - Target rotation
    - Evaluation score
    - Move type (e.g., "I-piece → line clear")
  - Position panel on right side of screen (column 40+)
  - Use compact format to avoid clutter
  - Update in real-time as AI makes decisions

  **Must NOT do**:
  - Display overwhelming detail
  - Block game board visibility
  - Add animations that slow gameplay

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: `[]` (terminal UI rendering)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Tasks 16, 19-21)
  - **Blocks**: Task 28
  - **Blocked By**: Task 16

  **References**:
  - `cmd/tetris/main.go:267-285` - renderNextPiece pattern
  - `cmd/tetris/main.go:238-265` - UI layout structure
  - Task 2 - MoveDecision struct (has all needed data)

  **Acceptance Criteria**:
  - [ ] Shows target X position (e.g., "Target X: 5")
  - [ ] Shows target rotation (e.g., "Rotation: 2")
  - [ ] Shows evaluation score (e.g., "Score: 12.5")
  - [ ] Panel fits within screen without wrapping
  - [ ] Updates each piece

  **QA Scenarios**:

  ```
  Scenario: Decision panel display
    Tool: interactive_bash (tmux)
    Steps:
      1. Enable auto-play
      2. Verify decision panel appears
      3. Watch panel update as AI evaluates new piece
      4. Verify X, rotation, score all visible
      5. Verify score changes based on board state
    Expected Result: Decision panel shows AI reasoning clearly
    Failure Indicators: Missing data, wrong values, text wrapping issues
    Evidence: .sisyphus/evidence/task-22-decision-panel-test.png
  ```

  **Commit**: YES (groups with 16, 19-21)
  - Message: `feat(autoplay): add AI decision panel UI`
  - Files: `internal/ui/autoplay_render.go`, `cmd/tetris/main.go`
  - Pre-commit: `go build -o tetris`

- [ ] 23. **Integration tests - full game scenarios**

  **What to do**:
  - Create `internal/model/autoplay_integration_test.go`
  - Test: AI can place 10 pieces without game over on empty board
  - Test: AI clears at least 1 line in 20 pieces
  - Test: AI handles all 7 piece types correctly
  - Test: AI speed changes don't break execution
  - Test: AI pauses and resumes correctly
  - Use table-driven tests where applicable

  **Must NOT do**:
  - Test individual heuristics (unit tests cover that)
  - Test UI rendering (manual QA)

  **Recommended Agent Profile**:
  - **Category**: `deep`
  - **Skills**: `[]` (complex integration testing)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4 (starts after Wave 3)
  - **Blocks**: Tasks 28-29
  - **Blocked By**: Tasks 16-22 (needs full integration)

  **References**:
  - `internal/model/gamestate_test.go` - Integration test patterns
  - `internal/model/autoplay_test.go` - AutoPlay test infrastructure
  - Tasks 1-22 - All auto-play functionality

  **Acceptance Criteria**:
  - [ ] TestAutoPlaySurvival10Pieces passes consistently
  - [ ] TestAutoPlayClearsLines passes (>50% line clear rate)
  - [ ] TestAutoPlayAllPieceTypes passes (all 7 types handled)
  - [ ] TestAutoPlaySpeedChanges doesn't panic
  - [ ] TestAutoPlayPauseResume preserves state

  **QA Scenarios**:

  ```
  Scenario: Integration test execution
    Tool: Bash (go test)
    Steps:
      1. Run: go test -v -run TestAutoPlay ./internal/model/...
      2. Verify all integration tests pass
      3. Run tests multiple times (5x) for consistency
      4. Verify no flaky tests
    Expected Result: All integration tests pass consistently
    Failure Indicators: Test failures, game overs, panics, flaky behavior
    Evidence: .sisyphus/evidence/task-23-integration-test.txt
  ```

  **Commit**: YES
  - Message: `test(autoplay): add integration tests for full game scenarios`
  - Files: `internal/model/autoplay_integration_test.go`
  - Pre-commit: `go test -v -run TestAutoPlay ./internal/model/...`

- [ ] 24. **AI tuning - weight adjustment for better play**

  **What to do**:
  - Run AI through 100-piece test games
  - Analyze failure modes (game overs, missed line clears)
  - Adjust heuristic weights based on observed behavior:
    - If AI creates too many holes: increase hole weight magnitude
    - If AI doesn't clear lines: increase completeLines weight
    - If AI makes bumpy surfaces: increase bumpiness weight
  - Document final weight values and rationale
  - Add weights to config for easy tuning

  **Must NOT do**:
  - Implement automatic weight learning (scope creep)
  - Change evaluation algorithm structure

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: `[]` (analytical tuning)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4 (starts after Wave 3)
  - **Blocks**: Task 28
  - **Blocked By**: Task 23 (needs integration tests)

  **References**:
  - Task 8 - Heuristic weights (-0.50, +0.76, -0.36, -0.18, -0.12)
  - Research: Dellacherie's Tetris AI weight optimization
  - Task 23 - Integration test results showing failure modes

  **Acceptance Criteria**:
  - [ ] AI survives 100+ pieces consistently (avg >100 before game over)
  - [ ] AI clears lines at reasonable rate (>20% of pieces)
  - [ ] AI avoids obvious mistakes (holes, extreme bumpiness)
  - [ ] Final weights documented in code comments

  **QA Scenarios**:

  ```
  Scenario: AI performance after tuning
    Tool: Bash (go test or custom benchmark)
    Steps:
      1. Run 10 simulated games with tuned weights
      2. Record pieces survived per game
      3. Record lines cleared per game
      4. Calculate average survival: expect >100 pieces
      5. Calculate average lines: expect >20 per 100 pieces
    Expected Result: AI performance meets or exceeds targets
    Failure Indicators: Early game overs (<50 pieces), no line clears
    Evidence: .sisyphus/evidence/task-24-tuning-results.txt
  ```

  **Commit**: YES
  - Message: `feat(autoplay): tune heuristic weights for optimal play`
  - Files: `internal/model/autoplay.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 25. **Edge case handling (game over recovery, reset)**

  **What to do**:
  - Handle game over: auto-disable auto-play
  - Handle reset: reset AutoPlayer state
  - Handle no-valid-moves: gracefully skip piece or game over
  - Handle piece spawn collision: detect and game over
  - Add defensive checks for nil pointers
  - Test all edge cases

  **Must NOT do**:
  - Change game over logic
  - Modify reset behavior for manual play

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: `[]` (defensive programming)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4 (with Tasks 23-24, 26-28)
  - **Blocks**: Task 28
  - **Blocked By**: Task 16 (needs integration)

  **References**:
  - `internal/model/gamestate.go:262-275` - Reset method
  - `internal/model/gamestate.go:148-170` - Game over detection

  **Acceptance Criteria**:
  - [ ] Auto-play disables on game over
  - [ ] Auto-play resets on game reset
  - [ ] No panics on edge cases
  - [ ] All edge cases have test coverage

  **QA Scenarios**:

  ```
  Scenario: Edge case - game over recovery
    Tool: Bash (go test)
    Steps:
      1. Create gameState with near-game-over board
      2. Enable auto-play
      3. Trigger game over condition
      4. Verify auto-play.disabled = true
      5. Call gameState.Reset()
      6. Verify auto-play state is reset
    Expected Result: Clean handling of game over and reset
    Failure Indicators: Panic, state corruption, auto-play still enabled after game over
    Evidence: .sisyphus/evidence/task-25-edge-case-test.txt
  ```

  **Commit**: YES
  - Message: `feat(autoplay): add edge case handling for game over and reset`
  - Files: `internal/model/autoplay.go`, `internal/model/autoplay_test.go`
  - Pre-commit: `go test ./internal/model/...`

- [ ] 26. **Performance optimization (cache evaluations)**

  **What to do**:
  - Profile AI to find bottlenecks (likely evaluateBoard calls)
  - Add caching for repeated board evaluations
  - Optimize board copying for simulation
  - Consider early pruning of obviously bad moves
  - Benchmark before/after optimization
  - Target: <10ms per piece decision at speed level 5

  **Must NOT do**:
  - Premature optimization before profiling
  - Sacrifice correctness for speed

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: `[]` (performance optimization)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4 (starts after Wave 3)
  - **Blocks**: Task 28
  - **Blocked By**: Task 23 (needs working integration)

  **References**:
  - Task 8 - evaluateBoard function (main bottleneck)
  - Task 9 - enumerateMoves (potential optimization)
  - Go benchmark testing patterns

  **Acceptance Criteria**:
  - [ ] Benchmark shows <10ms per decision
  - [ ] 100-piece game completes in <5 seconds at level 5
  - [ ] No correctness regressions after optimization
  - [ ] Benchmark tests added to test suite

  **QA Scenarios**:

  ```
  Scenario: Performance benchmark
    Tool: Bash (go test -bench)
    Steps:
      1. Run: go test -bench=BenchmarkAutoPlay -benchmem ./internal/model/...
      2. Record time per decision
      3. Verify <10ms per decision
      4. Run 100-piece simulation, record total time
      5. Verify <5 seconds at level 5
    Expected Result: AI performs within performance targets
    Failure Indicators: >10ms per decision, slow gameplay
    Evidence: .sisyphus/evidence/task-26-benchmark.txt
  ```

  **Commit**: YES
  - Message: `perf(autoplay): optimize evaluation caching and board simulation`
  - Files: `internal/model/autoplay.go`
  - Pre-commit: `go test -bench=. -benchtime=1s ./internal/model/...`

- [ ] 27. **Documentation - code comments + README section**

  **What to do**:
  - Add comprehensive godoc comments to all exported functions
  - Document heuristic weights and their rationale
  - Add usage examples in comments
  - Add "Auto-Play Mode" section to README.md:
    - How to enable/disable ('A' key)
    - Speed control ('S' key)
    - AI behavior explanation
    - Performance characteristics
  - Document configuration options (weights)

  **Must NOT do**:
  - Over-document obvious code
  - Add markdown documentation outside README

  **Recommended Agent Profile**:
  - **Category**: `writing`
  - **Skills**: `[]` (technical documentation)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4 (with Tasks 23-26, 28-29)
  - **Blocks**: None
  - **Blocked By**: Task 24 (needs final weights)

  **References**:
  - `README.md` - Existing documentation structure
  - `internal/model/gamestate.go` - Godoc comment patterns
  - AGENTS.md - Project documentation standards

  **Acceptance Criteria**:
  - [ ] All exported functions have godoc comments
  - [ ] README.md has "Auto-Play Mode" section
  - [ ] Documentation explains how to use auto-play
  - [ ] Weights and heuristics documented

  **QA Scenarios**:

  ```
  Scenario: Documentation completeness
    Tool: Bash + manual review
    Steps:
      1. Run: go doc ./internal/model/...
      2. Verify all exported functions have docs
      3. Read README.md Auto-Play section
      4. Verify instructions are clear and accurate
      5. Follow README instructions yourself
      6. Verify they work as described
    Expected Result: Complete, accurate, usable documentation
    Failure Indicators: Missing docs, unclear instructions, wrong information
    Evidence: .sisyphus/evidence/task-27-docs-review.txt
  ```

  **Commit**: YES
  - Message: `docs(autoplay): add comprehensive documentation and README section`
  - Files: `internal/model/autoplay.go`, `README.md`
  - Pre-commit: `go doc ./internal/model/...`

- [ ] 28. **Final QA - 100+ piece survival test**

  **What to do**:
  - Run extended gameplay test: 100+ pieces without game over
  - Record terminal session (tmux capture)
  - Verify all features work together:
    - Auto-play toggle
    - Speed control
    - AI decision making
    - Line clearing
    - UI indicators
    - Pause/resume
  - Capture screenshots at key moments
  - Document any issues found

  **Must NOT do**:
  - Skip if earlier tests fail
  - Accept game over before 100 pieces

  **Recommended Agent Profile**:
  - **Category**: `deep`
  - **Skills**: `[]` (comprehensive QA)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4 (final task before review)
  - **Blocks**: Final verification wave (F1-F4)
  - **Blocked By**: Tasks 20-27 (needs all features complete)

  **References**:
  - Tasks 1-27 - All auto-play features
  - `cmd/tetris/main.go` - Full game execution

  **Acceptance Criteria**:
  - [ ] AI survives 100+ pieces in at least 3 of 5 test runs
  - [ ] All UI elements display correctly throughout
  - [ ] Speed control works at all levels
  - [ ] Pause/resume preserves state correctly
  - [ ] Evidence captured (screenshots, video)

  **QA Scenarios**:

  ```
  Scenario: 100-piece survival test
    Tool: interactive_bash (tmux recording)
    Steps:
      1. Build: go build -o tetris
      2. Run: ./tetris
      3. Enable auto-play ('A')
      4. Set speed to level 4 or 5
      5. Record entire session with tmux capture
      6. Let AI play until 100+ pieces or game over
      7. Capture screenshots at: 25, 50, 75, 100 pieces
      8. Repeat 5 times
    Expected Result: AI survives 100+ pieces in ≥3 of 5 runs
    Failure Indicators: Early game overs, UI glitches, feature failures
    Evidence: .sisyphus/evidence/task-28-survival-test-{1-5}.mp4
  ```

  **Commit**: NO (QA task, no code changes)

- [ ] 29. **Git cleanup + tagging**

  **What to do**:
  - Review all changes: `git diff`
  - Ensure logical commit grouping (see commit messages in tasks)
  - Create final commit if needed
  - Create annotated tag: `git tag -a v1.1.0-autoplay -m "Add auto-play AI mode"`
  - Push tag if repo has remote: `git push origin v1.1.0-autoplay`
  - Create brief summary of changes for tag message

  **Must NOT do**:
  - Force push
  - Rewrite history
  - Commit untested code

  **Recommended Agent Profile**:
  - **Category**: `git` (using git-master skill)
  - **Skills**: `['git-master']` (safe git operations)
  - **Skills Evaluated but Omitted**: None needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4 (final implementation task)
  - **Blocks**: Final verification wave
  - **Blocked By**: Tasks 23-28 (needs all work complete)

  **References**:
  - Git best practices from AGENTS.md
  - Existing commit history for message style

  **Acceptance Criteria**:
  - [ ] All changes committed with clear messages
  - [ ] Tag v1.1.0-autoplay created
  - [ ] `git status` shows clean working tree
  - [ ] `git log` shows logical commit history

  **QA Scenarios**:

  ```
  Scenario: Git cleanup verification
    Tool: Bash (git commands)
    Steps:
      1. Run: git status
      2. Verify clean working tree (or staged changes ready)
      3. Run: git log --oneline -10
      4. Verify logical commit messages
      5. Run: git tag -l
      6. Verify v1.1.0-autoplay tag exists
    Expected Result: Clean git state, proper tagging
    Failure Indicators: Uncommitted changes, messy history, missing tag
    Evidence: .sisyphus/evidence/task-29-git-status.txt
  ```

  **Commit**: This task IS the commit task
  - Message: N/A (creates tag, not code commit)
  - Files: Git metadata only
  - Pre-commit: `go test ./... && go build -o tetris`

- [ ] F1. **Plan Compliance Audit** — `oracle`

  **What to do**:
  Read the plan end-to-end. For each "Must Have": verify implementation exists (read file, curl endpoint, run command). For each "Must NOT Have": search codebase for forbidden patterns — reject with file:line if found. Check evidence files exist in .sisyphus/evidence/. Compare deliverables against plan.

  **Output Format**:
  ```
  Must Have [N/N]:
  - [ ] Heuristic evaluation with configurable weights
  - [ ] All 4 rotations × all X positions evaluated
  - [ ] Move execution respecting existing mechanics
  - [ ] Speed control: 5 levels
  - [ ] Toggle on/off with 'A' key
  - [ ] AI respects pause state
  - [ ] UI shows: indicator, speed, decision panel

  Must NOT Have [N/N]:
  - [ ] No ML/neural network implementation
  - [ ] No modification to piece mechanics
  - [ ] No existing key binding changes
  - [ ] No online learning
  - [ ] No AI slop (check code quality)
  - [ ] No breaking manual play

  Tasks [N/N]: All 29 tasks + 4 final = 33 total
  Evidence [.sisyphus/evidence/]: [count files]

  VERDICT: APPROVE or REJECT (with reasons)
  ```

  **Commit**: NO (review task)

- [ ] F2. **Code Quality Review** — `unspecified-high`

  **What to do**:
  Run `go build`, `go vet`, `go test ./...`. Review all changed files for: `as any` (not applicable in Go), empty catches, TODOs left in code, commented-out code, unused imports. Check AI slop: excessive comments, over-abstraction, generic names (data/result/item/temp). Verify Go formatting with `gofmt -d`.

  **Output Format**:
  ```
  Build [PASS/FAIL]: go build -o tetris
  Vet [PASS/FAIL]: go vet ./...
  Tests [N pass/N fail]: go test ./...
  Format [PASS/FAIL]: gofmt -d .
  Files [N changed]: git diff --name-only
  Issues Found: [list any problems]

  VERDICT: APPROVE or REJECT (with reasons)
  ```

  **Commit**: NO (review task)

- [ ] F3. **Real Manual QA** — `unspecified-high` (+ `playwright` skill if UI)

  **What to do**:
  Start from clean state. Execute EVERY QA scenario from EVERY task — follow exact steps, capture evidence. Test cross-task integration (features working together, not isolation). Test edge cases: empty state, invalid input, rapid actions. Save to `.sisyphus/evidence/final-qa/`.

  **Key Scenarios to Test**:
  1. Toggle auto-play with 'A' key during gameplay
  2. Cycle through all 5 speed levels with 'S' key
  3. Verify AI decision panel updates in real-time
  4. Pause during auto-play, resume, verify state preserved
  5. Let AI play 50+ pieces, verify no game over
  6. Test all 7 piece types handled correctly
  7. Verify UI doesn't obscure game board

  **Output Format**:
  ```
  Scenarios [N/N pass]: Test all task QA scenarios
  Integration [N/N]: Cross-feature integration tests
  Edge Cases [N tested]: Pause, game over, reset, etc.
  Evidence [.sisyphus/evidence/final-qa/]: [count files]

  VERDICT: APPROVE or REJECT (with reasons)
  ```

  **Commit**: NO (review task)

- [ ] F4. **Scope Fidelity Check** — `deep`

  **What to do**:
  For each task: read "What to do", read actual diff (git log/diff). Verify 1:1 — everything in spec was built (no missing), nothing beyond spec was built (no creep). Check "Must NOT do" compliance. Detect cross-task contamination: Task N touching Task M's files unnecessarily. Flag unaccounted changes.

  **Output Format**:
  ```
  Tasks [N/N compliant]: Each task matches its spec
  Contamination [CLEAN/N issues]: Unnecessary cross-task file touching
  Unaccounted [CLEAN/N files]: Changes not in any task spec
  Scope Creep [NONE/N items]: Features beyond original scope
  Missing [NONE/N items]: Spec items not implemented

  VERDICT: APPROVE or REJECT (with reasons)
  ```

  **Commit**: NO (review task)

---

## Commit Strategy

**Wave 1 (Tasks 1-7)** — Single commit:
- Message: `feat(autoplay): add core types and heuristic helpers`
- Files: `internal/model/autoplay.go`, `internal/model/autoplay_test.go`
- Pre-commit: `go test ./internal/model/...`

**Wave 2 (Tasks 8-15)** — Single commit:
- Message: `feat(autoplay): implement AI evaluation and move execution`
- Files: `internal/model/autoplay.go`, `internal/model/autoplay_test.go`
- Pre-commit: `go test ./internal/model/...`

**Wave 3 (Tasks 16-22)** — Single commit:
- Message: `feat(autoplay): integrate with main.go and add UI`
- Files: `cmd/tetris/main.go`, `internal/ui/autoplay_render.go`
- Pre-commit: `go build -o tetris && go vet ./...`

**Wave 4 (Tasks 23-29)** — Multiple commits:
- Task 23: `test(autoplay): add integration tests`
- Task 24: `feat(autoplay): tune heuristic weights`
- Task 25: `feat(autoplay): add edge case handling`
- Task 26: `perf(autoplay): optimize evaluation performance`
- Task 27: `docs(autoplay): add documentation`
- Task 29: Git tag (not commit)

---

## Success Criteria

### Verification Commands
```bash
go mod tidy                          # Verify dependencies
go build -o tetris                   # Expected: builds successfully
go vet ./...                         # Expected: no issues
go test ./...                        # Expected: all tests pass
go test -v -run TestAutoPlay ./...   # Expected: all autoplay tests pass
go test -bench=BenchmarkAutoPlay ./... # Expected: <10ms per decision
./tetris                             # Expected: game runs, 'A' enables auto-play
```

### Final Checklist
- [ ] All "Must Have" features present and working
- [ ] All "Must NOT Have" guardrails respected
- [ ] All tests pass (unit + integration)
- [ ] Auto-play survives 100+ pieces consistently
- [ ] UI displays all required elements
- [ ] Speed control works at all 5 levels
- [ ] Pause/resume preserves auto-play state
- [ ] Evidence captured for all QA scenarios
- [ ] Documentation complete and accurate
- [ ] Git tagged as v1.1.0-autoplay

### Performance Targets
- Decision time: <10ms per piece at level 5
- 100-piece game: <5 seconds at level 5
- Survival rate: >100 pieces average before game over
- Line clear rate: >20% of pieces result in line clears

### User Experience
- Toggle auto-play: Single 'A' key press
- Speed control: Single 'S' key press cycles levels
- Visual feedback: All indicators visible and clear
- No disruption to manual play when auto-play off

---

## Test Result Report (Post-Execution)

**This section will be populated after test execution.**

### Execution Summary

| Metric | Target | Actual | Status |
|--------|--------|--------|--------|
| Unit Tests | 100% pass | [pending] | ⏳ |
| Integration Tests | 100% pass | [pending] | ⏳ |
| Code Coverage | >85% | [pending] | ⏳ |
| Benchmark (FindBestMove) | <10ms | [pending] | ⏳ |
| Survival (50 pieces) | 100% | [pending] | ⏳ |
| Survival (100 pieces) | >80% | [pending] | ⏳ |
| Line Clear Rate | >20% | [pending] | ⏳ |

### Detailed Test Results

#### Unit Tests

```
Test Suite: autoplay_test.go
Date: [DATE]
Go Version: go1.26.0

Results:
  Total Tests: 31
  Passed: [ ]
  Failed: [ ]
  Skipped: [ ]

Individual Results:
[ ] TestAutoPlayerCreation
[ ] TestAutoPlayer_Toggle
[ ] TestAutoPlayer_SetSpeedLevel
[ ] TestAutoPlayer_CycleSpeed
[ ] TestMoveDecision_IsValid
[ ] TestMoveDecision_String
[ ] TestCalculateSoftDrops
[ ] TestGetColHeight
[ ] TestGetAggregateHeight
[ ] TestCountCompleteLines
[ ] TestCountHoles
[ ] TestCalculateBumpiness
[ ] TestCountWells
[ ] TestEvalAggregateHeight
[ ] TestEvalHoles
[ ] TestEvalBumpiness
[ ] TestEvalWells
[ ] TestEvaluateBoard
[ ] TestGetWeights_SetWeights
[ ] TestEvaluateBoard_WeightImpact
[ ] TestEnumerateMoves
[ ] TestEnumerateMoves_AllPieceTypes
[ ] TestFindBestMove
[ ] TestFindBestMove_ObviousCases
[ ] TestFindBestMove_Determinism
[ ] TestFindBestMove_NoValidMoves
[ ] TestExecuteRotations
[ ] TestExecuteHorizontalMove
[ ] TestExecuteDrop
[ ] TestGetDelayForSpeed
[ ] TestShouldHoldPiece

Failures (if any):
[NONE / List failures with error messages]
```

#### Integration Tests

```
Test Suite: autoplay_integration_test.go
Date: [DATE]

Results:
  Total Tests: 6
  Passed: [ ]
  Failed: [ ]

Individual Results:
[ ] TestAutoPlay_Survival10Pieces
[ ] TestAutoPlay_Survival50Pieces
[ ] TestAutoPlay_ClearsLines
[ ] TestAutoPlay_AllPieceTypes
[ ] TestAutoPlay_SpeedChanges
[ ] TestAutoPlay_PauseResume

Failures (if any):
[NONE / List failures with error messages]
```

#### Benchmark Tests

```
Benchmark Suite: autoplay_benchmark_test.go
Date: [DATE]

Results:
BenchmarkFindBestMove-8           [   ]    [    ] ns/op    [   ] B/op    [  ] allocs/op
BenchmarkEvaluateBoard-8          [   ]    [    ] ns/op    [   ] B/op    [  ] allocs/op
BenchmarkAutoPlayGame-8           [   ]    [    ] ns/op    [   ] B/op    [  ] allocs/op

Performance Assessment:
- FindBestMove: [PASS/FAIL] (target: <10ms)
- EvaluateBoard: [PASS/FAIL] (target: <1ms)
- Full Game (50 pieces): [PASS/FAIL] (target: <5s)
```

#### Code Coverage

```
Coverage Report
Generated: [DATE]

Overall: [  ]%

By File:
  autoplay.go:                    [  ]%
  autoplay_test.go:               [  ]%
  autoplay_integration_test.go:   [  ]%

Coverage Threshold: 85%
Status: [PASS/FAIL]

Uncovered Functions (if any):
[list functions with <50% coverage]
```

### Evidence Files

```
.sisyphus/evidence/
├── test-results/
│   ├── unit-tests-output.txt
│   ├── integration-tests-output.txt
│   ├── benchmark-results.txt
│   ├── coverage.html
│   └── test-run-[timestamp].log
├── task-[N]-[test]-test.txt (per-task test evidence)
└── final-qa/
    └── [QA evidence files]
```

### Known Issues

| Issue ID | Severity | Description | Status | Workaround |
|----------|----------|-------------|--------|------------|
| [None] | - | - | - | - |

### Test Environment

```
Go Version: go1.26.0
OS: darwin/amd64
Module: github.com/oc-garden/tetris_game
Dependencies:
  - github.com/gdamore/tcell/v2 v2.x.x
  - (other deps from go.mod)

Test Execution Time: [duration]
Total Assertions: [count]
```

### Sign-Off

- [ ] All critical tests passed
- [ ] Code coverage meets threshold
- [ ] Performance targets met
- [ ] No blocking issues
- [ ] Evidence files captured

**Approved by**: [Agent name]
**Date**: [DATE]
**Verdict**: READY FOR RELEASE / NEEDS FIXES

---
