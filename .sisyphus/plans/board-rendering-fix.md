# Fix Board Rendering - Block Spacing and Edge Overlap

## TL;DR

> **Quick Summary**: Fix the board rendering so blocks are properly centered within cells and don't touch/overlap the borders.
> 
> **Deliverables**: 
> - Adjusted cell width calculation in `renderBoard()`
> - Proper border positioning with 1-cell padding
> - Blocks centered within their cells, not touching edges
> 
> **Estimated Effort**: Quick
> **Parallel Execution**: NO - single file change
> **Critical Path**: Edit `cmd/tetris/main.go:renderBoard()` → Build → Test

---

## Context

### Original Request
"adjust game board, make the block fit in the board and do not touch or overlap the edges"

### Current Issue Analysis
The `renderBoard()` function in `cmd/tetris/main.go` has rendering issues:
- `cellWidth := 2` but blocks are drawn at exact cell boundaries
- Left border at `boardX` (col 2), right border at `boardX+10*cellWidth` (col 22)
- Blocks rendered at `boardX+x*cellWidth` positions (2,4,6,8...20)
- Result: Blocks touch or appear to overlap the vertical borders

### Code Analysis
**Current rendering logic** (lines 126-194):
```go
boardX := 2
boardY := 2
cellWidth := 2

// Borders drawn at exact boundaries
screen.SetContent(boardX, boardY+y, '║', ...)              // Left border at col 2
screen.SetContent(boardX+10*cellWidth, boardY+y, '║', ...) // Right border at col 22

// Blocks drawn at cell start positions
screen.SetContent(boardX+x*cellWidth, screenY, '█', ...)   // Block at cols 2,4,6...
screen.SetContent(boardX+x*cellWidth+1, screenY, '█', ...) // Block extends to cols 3,5,7...
```

**Problem**: With cellWidth=2 and 10 columns:
- Cell 0: columns 2-3 (touches left border at col 2!)
- Cell 9: columns 20-21 (right border at col 22, but block at 20-21 is too close)

---

## Work Objectives

### Core Objective
Adjust the board rendering so each block is centered within its cell with visible spacing from borders.

### Concrete Deliverables
- Modified `renderBoard()` function in `cmd/tetris/main.go`
- Proper cell width and border positioning
- Visual verification that blocks don't touch edges

### Definition of Done
- [ ] Build succeeds: `go build -o tetris`
- [ ] Game runs: `./tetris` shows properly spaced blocks
- [ ] Blocks are visibly centered within cells
- [ ] Clear gap between blocks and all four borders

### Must Have
- Each cell should have visible padding on all sides
- Blocks should not visually merge with borders
- Maintain the 10×20 board dimensions

### Must NOT Have (Guardrails)
- Do NOT change game logic or collision detection
- Do NOT change board dimensions (still 10×20)
- Do NOT modify piece movement speed or controls
- Do NOT change the UI panel layout (next piece, hold, score area)

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: YES (standard Go testing)
- **Automated tests**: NO (this is a visual/UI fix)
- **Framework**: N/A
- **Agent-Executed QA**: YES (mandatory)

### QA Policy
Every task MUST include agent-executed QA scenarios.
Evidence saved to `.sisyphus/evidence/task-{N}-{scenario-slug}.{ext}`.

- **TUI/CLI**: Use `interactive_bash` (tmux) — Run game, observe rendering
- **Visual verification**: Screenshot captures showing proper block spacing

---

## Execution Strategy

### Sequential Execution (Single Task)

```
Wave 1 (Single task):
└── Task 1: Fix board rendering in renderBoard() [visual-engineering]

Critical Path: Task 1 → Build → Visual QA
```

### Agent Dispatch Summary
- **1**: **1** — T1 → `visual-engineering` (UI rendering fix)

---

## TODOs

- [ ] 1. Fix board rendering - adjust cell width and border positioning

  **What to do**:
  - Modify `renderBoard()` function in `cmd/tetris/main.go`
  - Change cell width calculation to properly center blocks:
    - Option A: Increase `cellWidth` from 2 to 3 (each block gets 3 chars: space-block-space)
    - Option B: Add 1-char padding on left side, adjust right border position
  - Update border positions to account for new cell width:
    - Left border stays at `boardX`
    - Right border moves to `boardX + 10*newCellWidth`
    - OR add explicit padding: left border at `boardX+1`, right at `boardX+1+10*2`
  - Update block rendering to use new cell positions
  - Update top/bottom border loops to match new width
  - Update corner positions to match new dimensions
  - Test build: `go build -o tetris`

  **Must NOT do**:
  - Do NOT change game logic in `internal/model/`
  - Do NOT modify collision detection
  - Do NOT change board dimensions (10×20)
  - Do NOT modify UI panel positions (score, next, hold areas)

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: This is a UI rendering fix requiring visual precision
  - **Skills**: [`dev-browser`]
    - `dev-browser`: Can be used for visual verification via screenshots if needed
  - **Skills Evaluated but Omitted**:
    - `playwright`: Not needed - this is a terminal UI, not web
    - `git-master`: Not needed - user will handle commits

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Sequential (only task)
  - **Blocks**: None
  - **Blocked By**: None (can start immediately)

  **References**:
  - `cmd/tetris/main.go:126-194` - Current `renderBoard()` implementation to modify
  - `cmd/tetris/main.go:132-145` - Border drawing logic to adjust
  - `cmd/tetris/main.go:147-158` - Block rendering to adjust
  - `cmd/tetris/main.go:177-193` - Current piece rendering to adjust

  **Acceptance Criteria**:
  - [ ] `go build -o tetris` succeeds with no errors
  - [ ] `go vet ./...` passes with no warnings
  - [ ] Game launches: `./tetris` runs without crash
  - [ ] Visual: Board has clear left/right padding (blocks don't touch vertical borders)
  - [ ] Visual: Board has clear top/bottom padding (blocks don't touch horizontal borders)
  - [ ] Visual: All 7 tetromino types render with proper spacing when tested

  **QA Scenarios**:

  ```
  Scenario: Verify block spacing from left border
    Tool: interactive_bash (tmux)
    Preconditions: Game built and running in tmux session
    Steps:
      1. Launch game in tmux: send-keys "cd /path && ./tetris" Enter
      2. Wait 2 seconds for game to initialize
      3. Capture screenshot: capture-pane -S 200 -t tetris-session
      4. Inspect leftmost column of board (x=0 in game coords)
    Expected Result: Block at game x=0 renders at screen column 3 (not 2), with space between block and left border
    Failure Indicators: Block character '█' appears immediately adjacent to '║' border character
    Evidence: .sisyphus/evidence/task-1-left-spacing.png

  Scenario: Verify block spacing from right border
    Tool: interactive_bash (tmux)
    Preconditions: Game running, piece moved to right edge
    Steps:
      1. Send right arrow key multiple times to move piece to rightmost position
      2. Capture screenshot
      3. Inspect rightmost column of board (x=9 in game coords)
    Expected Result: Block at game x=9 has visible space before right border '║'
    Failure Indicators: Block character touches or overlaps with right border
    Evidence: .sisyphus/evidence/task-1-right-spacing.png

  Scenario: Verify block spacing from top border
    Tool: interactive_bash (tmux)
    Preconditions: New piece spawned at top of board
    Steps:
      1. Observe freshly spawned piece at y=18-19
      2. Capture screenshot
      3. Inspect top row rendering
    Expected Result: Blocks at top rows have space below top border '═'
    Failure Indicators: Block appears to merge with top border
    Evidence: .sisyphus/evidence/task-1-top-spacing.png

  Scenario: Verify all four borders have proper spacing
    Tool: interactive_bash (tmux)
    Preconditions: Game running with piece in center
    Steps:
      1. Capture full board screenshot
      2. Verify all corners show clear separation
      3. Check that board looks visually centered with even padding
    Expected Result: Board appears as a framed box with blocks floating inside, not touching edges
    Failure Indicators: Any block visually merges with any border
    Evidence: .sisyphus/evidence/task-1-full-board.png
  ```

  **Evidence to Capture**:
  - [ ] Screenshot showing left edge spacing
  - [ ] Screenshot showing right edge spacing
  - [ ] Screenshot showing full board with all borders visible
  - [ ] Terminal output from `go build` and `go vet`

  **Commit**: YES
  - Message: `fix(ui): adjust board rendering for proper block spacing`
  - Files: `cmd/tetris/main.go`
  - Pre-commit: `go build -o tetris && go vet ./...`

---

## Final Verification Wave

- [ ] F1. **Plan Compliance Audit** — `oracle`
  Verify board dimensions unchanged (10×20), only rendering adjusted. Check evidence files exist.
  Output: `Visual spacing [VERIFIED] | Board dims [10×20] | VERDICT: APPROVE/REJECT`

- [ ] F2. **Code Quality Review** — `unspecified-high`
  Run `go build -o tetris && go vet ./...`. Check no `//nolint`, no unused imports.
  Output: `Build [PASS] | Vet [PASS] | VERDICT`

- [ ] F3. **Real Manual QA** — `unspecified-high` + `dev-browser`
  Run game, verify visual spacing on all sides. Test with each piece type.
  Output: `Spacing [4/4 sides OK] | All pieces [7/7 verified] | VERDICT`

- [ ] F4. **Scope Fidelity Check** — `deep`
  Verify only `renderBoard()` modified, no game logic changes.
  Output: `Files [1 modified] | Game logic [UNCHANGED] | VERDICT`

---

## Commit Strategy

- **1**: `fix(ui): adjust board rendering for proper block spacing` — cmd/tetris/main.go, go build && go vet

---

## Success Criteria

### Verification Commands
```bash
go build -o tetris && ./tetris    # Expected: Game launches, blocks properly spaced
go vet ./...                       # Expected: No warnings
```

### Final Checklist
- [ ] All borders have visible padding (blocks don't touch)
- [ ] Build succeeds without errors
- [ ] Game runs without crash
- [ ] Visual evidence captured for all 4 sides
