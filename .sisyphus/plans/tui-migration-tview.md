# TUI Implementation Plan - Migrate to tview

## TL;DR

> **Quick Summary**: Replace raw tcell drawing with tview (TUI framework built on tcell) for proper UI components, better layout management, and cleaner code.
> 
> **Deliverables**: 
> - Replace manual rendering with tview components (Flex, Frame, TextView, Box)
> - Clean separation of UI concerns
> - Proper borders, colors, and layouts
> - No visual regression - same game appearance
> 
> **Estimated Effort**: Medium (1-2 hours)
> **Breaking Changes**: None (internal UI refactor only)

---

## Context

### Current State
- **Library**: Raw tcell (low-level terminal manipulation)
- **Issues**:
  - Manual border drawing (lines 198-211 in main.go)
  - Manual text positioning (prone to overlap bugs like we just fixed)
  - No layout management
  - Hard-coded coordinates throughout
  - Manual clear/redraw logic

### Target State
- **Library**: tview (built on tcell, provides high-level components)
- **Benefits**:
  - Proper layout system (Flexbox-like)
  - Built-in borders, padding, margins
  - Automatic repositioning
  - Cleaner, more maintainable code
  - Easier to extend (add panels, help screens, etc.)

---

## Work Objectives

### Core Objective
Migrate from raw tcell drawing to tview while maintaining identical visual appearance and game functionality.

### Concrete Deliverables
1. New `internal/ui/game_view.go` - Main game board component using tview
2. New `internal/ui/status_panel.go` - Side panel with score, level, lines
3. New `internal/ui/autoplayer_panel.go` - Autoplay status and AI decisions
4. Updated `cmd/tetris/main.go` - Integration with tview Application
5. go.mod updated with `github.com/rivo/tview` dependency

### Definition of Done
- [ ] Game renders identically to before (board, pieces, borders)
- [ ] No text overlap issues
- [ ] All game functionality works (input, scoring, game over)
- [ ] Autoplay mode displays correctly
- [ ] `go build ./...` succeeds
- [ ] `go run ./cmd/tetris` runs without errors

### Must Have
- tview for UI components
- Maintain 10×20 board with 2-char wide cells
- Keep all existing colors and styles
- Preserve autoplay indicator positioning
- No visual regression

### Must NOT Have (Guardrails)
- No changing game logic
- No changing board dimensions
- No changing piece colors
- No breaking existing input handling
- No removing sound effects

---

## Execution Strategy

### Sequential Tasks

```
Wave 1 (Setup — install dependencies):
├── Task 1: Add tview dependency [quick]
└── Task 2: Create basic tview app structure [quick]

Wave 2 (UI Components):
├── Task 3: Implement GameBoard component [unspecified-high]
├── Task 4: Implement StatusPanel component [unspecified-high]
└── Task 5: Implement AutoPlayerPanel component [unspecified-high]

Wave 3 (Integration):
├── Task 6: Integrate tview with main game loop [deep]
└── Task 7: Remove old tcell rendering code [quick]

Wave 4 (Verification):
├── Task 8: Visual testing - verify appearance [quick]
└── Task 9: Run all tests [quick]

Critical Path: Task 1 → Task 3 → Task 4 → Task 6 → Task 8
```

---

## TODOs

- [ ] 1. Add tview Dependency

  **What to do**:
  - Run: `go get github.com/rivo/tview`
  - Verify: `go mod tidy`
  - Check: `go mod why github.com/rivo/tview`
  
  **Must NOT do**:
  - Don't remove tcell (tview depends on it)
  - Don't update other dependencies
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: None needed
  
  **Acceptance Criteria**:
  - [ ] `go.mod` includes `github.com/rivo/tview`
  - [ ] `go mod tidy` succeeds
  - [ ] No dependency conflicts
  
  **QA Scenarios**:
  ```
  Scenario: Dependency added successfully
    Tool: Bash
    Steps:
      1. Run: go get github.com/rivo/tview
      2. Run: go mod tidy
      3. Run: go build ./...
    Expected Result: All commands succeed
    Evidence: .sisyphus/evidence/task-1-dependency.txt
  ```

- [ ] 2. Create Basic tview Application Structure

  **What to do**:
  - Create `internal/ui/app.go`
  - Set up tview.Application with basic layout
  - Create Flex container for board + side panel
  - Add minimal test to verify tview renders
  
  **Must NOT do**:
  - Don't implement full game board yet
  - Don't remove existing rendering
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: None needed
  
  **Acceptance Criteria**:
  - [ ] Basic tview app compiles
  - [ ] Simple "Hello World" renders in terminal
  - [ ] App can be started and stopped cleanly
  
  **QA Scenarios**:
  ```
  Scenario: tview app renders
    Tool: Bash
    Steps:
      1. Create minimal tview app
      2. Run: go run ./cmd/tetris
      3. Verify tview window appears
    Expected Result: tview renders without errors
    Evidence: .sisyphus/evidence/task-2-basic-app.txt
  ```

- [ ] 3. Implement GameBoard Component

  **What to do**:
  - Create `internal/ui/game_board.go`
  - Use tview.Box or custom primitive for board
  - Implement Draw() method to render:
    - Board borders (double-line style)
    - Cells with colors (use tcell colors)
    - Current piece
    - Line clear animation (flash white)
  - Maintain 10×20 grid with 2-char wide cells
  
  **Must NOT do**:
  - Don't change board dimensions
  - Don't change piece colors
  - Don't change border style
  
  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: [`frontend-ui-ux`] if available
  
  **References**:
  - `cmd/tetris/main.go:189-259` — Current renderBoard() implementation
  - tview docs: https://pkg.go.dev/github.com/rivo/tview#Box
  
  **Acceptance Criteria**:
  - [ ] Board renders with correct borders (╔═╗║╚═╝)
  - [ ] Pieces render with correct colors
  - [ ] Line clear animation works
  - [ ] No text overlap issues
  
  **QA Scenarios**:
  ```
  Scenario: Game board renders correctly
    Tool: Bash
    Steps:
      1. Start game
      2. Verify board appears at correct position
      3. Move pieces, verify colors match
      4. Clear lines, verify flash animation
    Expected Result: Board looks identical to tcell version
    Evidence: .sisyphus/evidence/task-3-game-board.png (screenshot)
  ```

- [ ] 4. Implement StatusPanel Component

  **What to do**:
  - Create `internal/ui/status_panel.go`
  - Use tview.TextView or Flex with TextViews
  - Display:
    - Score
    - Level
    - Lines cleared
    - Next piece preview
  - Right side of game board
  
  **Must NOT do**:
  - Don't change information displayed
  - Don't change layout (right side)
  
  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: [`frontend-ui-ux`] if available
  
  **References**:
  - `cmd/tetris/main.go:260-290` — Current renderUI() implementation
  
  **Acceptance Criteria**:
  - [ ] Score, Level, Lines displayed
  - [ ] Next piece preview renders
  - [ ] Panel positioned to right of board
  - [ ] No text wrapping issues
  
  **QA Scenarios**:
  ```
  Scenario: Status panel displays correctly
    Tool: Bash
    Steps:
      1. Start game
      2. Verify score/level/lines update
      3. Check next piece preview
    Expected Result: Panel displays all info correctly
    Evidence: .sisyphus/evidence/task-4-status-panel.png
  ```

- [ ] 5. Implement AutoPlayerPanel Component

  **What to do**:
  - Create `internal/ui/autoplayer_panel.go`
  - Use tview.TextView
  - Display:
    - "AUTO-PLAY: ON/OFF" indicator (left side)
    - "SPEED: 1-5" indicator (left side)
    - AI decision panel (target X, rotation, drops, score)
  
  **Must NOT do**:
  - Don't change text content
  - Don't change positioning (left side for indicators)
  
  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: [`frontend-ui-ux`] if available
  
  **References**:
  - `internal/ui/autoplay_render.go` — Current autoplay rendering
  - `cmd/tetris/main.go:132-136` — Current autoplay indicator calls
  
  **Acceptance Criteria**:
  - [ ] AUTO-PLAY indicator displays (left side, no overlap)
  - [ ] SPEED indicator displays (left side)
  - [ ] AI decision panel shows when autoplay enabled
  - [ ] No text overlap with game board
  
  **QA Scenarios**:
  ```
  Scenario: Autoplay panel displays correctly
    Tool: Bash
    Steps:
      1. Enable autoplay (press 'a')
      2. Verify indicators appear on left
      3. Verify no overlap with board
      4. Check AI decision panel updates
    Expected Result: All autoplay UI displays without overlap
    Evidence: .sisyphus/evidence/task-5-autoplayer-panel.png
  ```

- [ ] 6. Integrate tview with Main Game Loop

  **What to do**:
  - Update `cmd/tetris/main.go`
  - Create tview.Application in main()
  - Set up input handling:
    - Arrow keys for piece movement
    - Space for hard drop
    - 'a' for autoplay toggle
    - 'q' for quit
  - Integrate game state updates with tview draw calls
  - Handle game over screen
  
  **Must NOT do**:
  - Don't change input mappings
  - Don't change game logic
  - Don't break existing key handlers
  
  **Recommended Agent Profile**:
  - **Category**: `deep`
  - **Skills**: None needed (complex integration)
  
  **References**:
  - `cmd/tetris/main.go:16-168` — Current main() and input handling
  
  **Acceptance Criteria**:
  - [ ] All keys work (arrows, space, 'a', 'q')
  - [ ] Game loop integrates with tview
  - [ ] Game over screen displays
  - [ ] Replay (r) and quit (q) work after game over
  
  **QA Scenarios**:
  ```
  Scenario: Full game loop with tview
    Tool: Bash
    Steps:
      1. Start game
      2. Move pieces with arrow keys
      3. Hard drop with space
      4. Toggle autoplay with 'a'
      5. Lose game, verify game over screen
      6. Press 'r' to replay or 'q' to quit
    Expected Result: All functionality works identically to tcell version
    Evidence: .sisyphus/evidence/task-6-integration.txt
  ```

- [ ] 7. Remove Old tcell Rendering Code

  **What to do**:
  - Delete or comment out old render functions:
    - `renderBoard()` in main.go
    - `renderUI()` in main.go
    - `renderGameOver()` in main.go
  - Remove unused imports
  - Clean up any tcell-only code
  
  **Must NOT do**:
  - Don't break compilation
  - Don't remove game state management
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: None needed
  
  **Acceptance Criteria**:
  - [ ] No references to old render functions
  - [ ] Code compiles without warnings
  - [ ] Binary size similar to before
  
  **QA Scenarios**:
  ```
  Scenario: Old code removed cleanly
    Tool: Bash
    Steps:
      1. Run: grep -n "renderBoard\|renderUI" cmd/tetris/main.go
      2. Run: go build ./...
    Expected Result: No old render calls, build succeeds
    Evidence: .sisyphus/evidence/task-7-cleanup.txt
  ```

- [ ] 8. Visual Testing - Verify Appearance

  **What to do**:
  - Run game manually
  - Compare visual appearance to tcell version
  - Check:
    - Board position and size
    - Colors match
    - Text positioning (no overlap)
    - Borders correct style
    - All panels visible
  
  **Must NOT do**:
  - Don't accept visual regressions
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`frontend-ui-ux`] if available
  
  **Acceptance Criteria**:
  - [ ] Board looks identical to tcell version
  - [ ] No text overlap
  - [ ] All colors match
  - [ ] Borders correct style
  
  **QA Scenarios**:
  ```
  Scenario: Visual comparison
    Tool: Manual
    Steps:
      1. Run tcell version (git stash current changes)
      2. Take screenshot of game in progress
      3. Run tview version
      4. Take screenshot of same game state
      5. Compare screenshots
    Expected Result: No visual differences
    Evidence: .sisyphus/evidence/task-8-visual-comparison.png
  ```

- [ ] 9. Run All Tests

  **What to do**:
  - Run: `go test ./...`
  - Verify all tests pass
  - Run: `go test ./internal/model/...` (game logic tests)
  - Ensure no test failures from UI changes
  
  **Must NOT do**:
  - Don't accept test failures
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: None needed
  
  **Acceptance Criteria**:
  - [ ] All tests pass
  - [ ] No test modifications needed (UI changes don't affect logic)
  
  **QA Scenarios**:
  ```
  Scenario: Full test suite passes
    Tool: Bash
    Steps:
      1. Run: go test ./... -v
      2. Capture output
      3. Verify all PASS
    Expected Result: All tests pass
    Evidence: .sisyphus/evidence/task-9-tests.txt
  ```

---

## Implementation Details

### tview Component Hierarchy

```
Application
└── Flex (horizontal)
    ├── TextView (left panel)
    │   ├── "AUTO-PLAY: ON/OFF"
    │   └── "SPEED: 1-5"
    ├── Box (game board - custom primitive)
    │   ├── Borders (double-line)
    │   ├── Cells (colored blocks)
    │   └── Current piece
    └── Flex (right panel - vertical)
        ├── TextView (Score)
        ├── TextView (Level)
        ├── TextView (Lines)
        └── TextView (Next piece preview)
```

### Key tview APIs

- **`tview.NewApplication()`**: Create application
- **`tview.NewFlex()`**: Flexbox layout container
- **`tview.NewTextView()`**: Text display with colors
- **`tview.NewBox()`**: Base for custom primitives
- **`SetDrawFunc()`**: Custom rendering for game board
- **`SetInputCapture()`**: Global key handler

### Color Mapping

Keep existing color mapping (tcell colors work with tview):

```go
colors := map[int]tcell.Color{
    1: tcell.ColorCyan,
    2: tcell.ColorYellow,
    3: tcell.ColorPurple,
    4: tcell.ColorGreen,
    5: tcell.ColorRed,
    6: tcell.ColorBlue,
    7: tcell.ColorOrange,
}
```

---

## Success Criteria

### Verification Commands
```bash
go get github.com/rivo/tview          # Add dependency
go mod tidy                            # Sync dependencies
go build ./...                         # Build succeeds
go run ./cmd/tetris                    # Game runs
go test ./...                          # All tests pass
```

### Final Checklist
- [ ] tview dependency added
- [ ] All UI components implemented
- [ ] Game renders correctly
- [ ] No text overlap
- [ ] All inputs work
- [ ] All tests pass
- [ ] No visual regression

### Success Metrics

**Before** (tcell):
- Manual coordinate management
- Text overlap bugs (as we experienced)
- ~325 lines in main.go
- No layout system

**After** (tview):
- Component-based layout
- Automatic positioning (no overlap)
- Cleaner separation of concerns
- Professional TUI framework

---

## Risk Mitigation

### Potential Issues

1. **tview doesn't support custom drawing well**
   - Mitigation: Use Box.SetDrawFunc() for custom board rendering
   - Fallback: Keep tcell for board, use tview for panels only

2. **Input handling conflicts**
   - Mitigation: Use SetInputCapture() for global keys, component input for rest
   - Test: All keys (arrows, space, 'a', 'q', 'r')

3. **Performance degradation**
   - Mitigation: Profile frame rate, ensure 60 FPS maintained
   - Benchmark: Compare tcell vs tview rendering time

4. **Visual differences**
   - Mitigation: Manual visual testing before merge
   - Acceptance: Must look identical to tcell version

### Rollback Plan
If tview integration fails:
1. `git stash` current changes
2. Keep tcell implementation
3. Document why tview didn't work
4. Consider alternative: tcell + custom layout manager

---

## Commit Strategy

**Single commit for clean migration**:
```
refactor(ui): migrate from tcell to tview TUI framework

- Add tview dependency
- Implement GameBoard, StatusPanel, AutoPlayerPanel components
- Replace manual tcell rendering with tview components
- Maintain identical visual appearance
- Fix text overlap issues with proper layout management

Benefits:
- Cleaner code organization
- Automatic layout management
- No coordinate-based positioning bugs
- Easier to extend with new UI elements
```

---

## Files to Create/Modify

### Create:
- `internal/ui/app.go` - tview Application setup
- `internal/ui/game_board.go` - Game board component
- `internal/ui/status_panel.go` - Score/level/lines panel
- `internal/ui/autoplayer_panel.go` - Autoplay indicators

### Modify:
- `cmd/tetris/main.go` - Integrate tview, remove old rendering
- `go.mod` - Add tview dependency
- `internal/ui/autoplay_render.go` - Adapt for tview

### Delete (or comment out):
- Old render functions in `cmd/tetris/main.go`

---

## Estimated Timeline

- **Wave 1** (Setup): 10 minutes
- **Wave 2** (Components): 45 minutes
- **Wave 3** (Integration): 30 minutes
- **Wave 4** (Verification): 15 minutes
- **Total**: ~100 minutes (1 hour 40 minutes)

---

**Created**: 2026-03-14
**Author**: Prometheus (Planning Agent)
**Status**: READY FOR EXECUTION
