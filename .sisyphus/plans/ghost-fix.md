# Ghost Mode Deadlock and Missing INFO Text Fix

## TL;DR

> **Quick Summary**: Fix two critical issues: (1) INFO panel missing "Ghost (manual only)" text during dynamic updates, (2) Potential deadlock/freeze when ghost mode is enabled due to unsafe rendering
> 
> **Deliverables**:
> - Add "GHOST" line to dynamic INFO panel updates
> - Add comprehensive safety checks in ghost rendering to prevent crashes
> - Verify no blocking calls in game loop
> 
> **Estimated Effort**: Quick
> **Critical Path**: Fix INFO text → Add safety checks → Test ghost mode

---

## Context

### Original Request
User reports:
1. Cannot see "Ghost (manual only)" text in INFO panel
2. App freezes/deadlocks when ghost mode is switched on

### Analysis Findings

**Issue 1: Missing INFO Text**
- Initial setup (line 182): HAS "GHOST [white](manual only)" ✅
- Dynamic update (lines 415-428): MISSING ghost line ❌
- When score/board changes, entire text replaced without ghost control

**Issue 2: Potential Deadlock Causes**
- GetGhostY() loops through board positions
- Ghost rendering iterates through matrix without comprehensive bounds checks
- Potential for negative indexing or out-of-bounds access
- No blocking calls found, but safety checks incomplete

---

## Work Objectives

### Core Objective
Fix the ghost mode implementation to be safe and ensure INFO panel always shows all controls.

### Concrete Deliverables
- Updated INFO panel text with ghost line in dynamic updates
- Safe ghost rendering with comprehensive bounds checking
- No app freezes when toggling ghost mode

### Definition of Done
- Ghost mode toggle works without freezing
- INFO panel always shows "G Ghost (manual only)" line
- No deadlocks or crashes in game loop

### Must Have
- All bounds checks for ghost piece rendering
- Matrix nil checks
- Y position range validation
- Ghost only renders when safe (manual mode, valid state)

### Must NOT Have
- Blocking calls in render loop
- Array out-of-bounds access
- Infinite loops in GetGhostY()

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: YES
- **Automated tests**: Run existing tests
- **Manual QA**: Test ghost toggle in-game

### QA Policy
- Manual testing required for UI rendering
- Verify no freeze after 30 seconds of ghost mode usage

---

## Execution Strategy

### Sequential Steps

```
Step 1: Fix INFO panel dynamic update [quick]
├── Add "G Ghost (manual only)" line to SetText call
└── Ensure consistency with initial setup

Step 2: Add comprehensive safety checks [quick]
├── Validate ghostY range before rendering
├── Check matrix not nil
├── Validate each cell's calculated position
└── Add board bounds verification

Step 3: Test and verify [quick]
├── go build
├── Manual test: toggle G key
├── Verify no freeze
└── Verify INFO text visible
```

---

## TODOs

- [ ] 1. Fix INFO panel to include Ghost control line in dynamic updates

  **What to do**:
  - Update the `infoBox.SetText()` call in the game loop (lines 415-428)
  - Add the line: `"[#FFFF00]G[white] Ghost [white](manual only)\n"+`
  - Insert it between "Space Hard Drop" and "P Pause" lines
  - Ensure format matches initial setup exactly

  **Must NOT do**:
  - Do not change other control lines
  - Do not alter the formatting style

  **Recommended Agent Profile**:
  - **Category**: `quick` - Simple text change
  - **Skills**: None needed

  **Parallelization**:
  - **Sequential**: Complete before testing

  **References**:
  - `cmd/tetris/main.go:172-186` - Initial INFO text setup
  - `cmd/tetris/main.go:415-428` - Dynamic update to fix

  **Acceptance Criteria**:
  - [ ] INFO panel shows "G Ghost (manual only)" after score changes
  - [ ] Text format matches initial setup

  **QA Scenarios**:
  ```
  Scenario: Verify ghost text after score update
    Tool: interactive_bash (tmux)
    Steps:
      1. Run game and clear some lines (change score)
      2. Check INFO panel shows "G Ghost (manual only)"
    Expected: Ghost line visible in INFO panel
  ```

  **Commit**: YES

- [ ] 2. Add comprehensive bounds checking in ghost rendering

  **What to do**:
  - Add check: `ghostY < 20` (within board)
  - Add check: `pieceY >= 0 && pieceY < 20` (valid row)
  - Add matrix nil check before iteration
  - Validate calculated screen coordinates before SetContent

  **Must NOT do**:
  - Do not change the ghost rendering logic itself
  - Only add safety guards

  **Recommended Agent Profile**:
  - **Category**: `quick` - Defensive programming
  - **Skills**: None needed

  **Parallelization**:
  - **Sequential**: After Step 1

  **References**:
  - `cmd/tetris/main.go:55-73` - Ghost rendering code
  - `internal/model/gamestate.go:128-142` - GetGhostY() implementation

  **Acceptance Criteria**:
  - [ ] No crashes when toggling ghost mode
  - [ ] No out-of-bounds array access
  - [ ] Ghost renders correctly within board

  **QA Scenarios**:
  ```
  Scenario: Ghost mode stability
    Tool: interactive_bash (tmux)
    Steps:
      1. Start game
      2. Press 'G' to enable ghost
      3. Move piece around for 30 seconds
      4. Press 'G' to disable
      5. Verify no freeze or crash
    Expected: No app freeze, smooth operation
    Evidence: .sisyphus/evidence/task-2-ghost-stability.txt
  ```

  **Commit**: YES

- [ ] 3. Build and manual test

  **What to do**:
  - Run `go build -o tetris ./cmd/tetris`
  - Manually test ghost toggle
  - Verify INFO panel text
  - Test for 1 minute with ghost enabled

  **QA Scenarios**:
  ```
  Scenario: Full ghost mode test
    Tool: interactive_bash (tmux)
    Steps:
      1. Build: go build -o tetris ./cmd/tetris
      2. Run game
      3. Verify INFO shows "G Ghost (manual only)"
      4. Press 'G' - ghost appears
      5. Move piece - ghost updates position
      6. Press 'G' - ghost disappears
      7. Enable auto-play - verify ghost OFF
      8. Disable auto-play - ghost can be re-enabled
    Expected: All operations smooth, no freezes
  ```

  **Commit**: YES

---

## Final Verification Wave

- [ ] F1. **Plan Compliance Audit** — `oracle`
  Verify INFO text has ghost line, verify all safety checks present

- [ ] F2. **Build Verification** — `quick`
  `go build` succeeds, no errors

- [ ] F3. **Manual QA** — `unspecified-high`
  Test ghost mode for 1+ minute, no freezes

- [ ] F4. **Scope Fidelity** — `deep`
  Only changed INFO text and added safety checks

---

## Commit Strategy

- **1**: `fix(ghost): add comprehensive bounds checking to prevent crashes`
  - cmd/tetris/main.go (ghost rendering safety)
  - internal/model/gamestate.go (GetGhostY safety)

- **2**: `fix(ui): include Ghost control in INFO panel dynamic updates`
  - cmd/tetris/main.go (infoBox.SetText call)

---

## Success Criteria

### Verification Commands
```bash
go build -o tetris ./cmd/tetris   # Expected: success
# Manual: Run game, press G, verify no freeze
# Manual: Check INFO panel shows "G Ghost (manual only)"
```

### Final Checklist
- [ ] INFO panel has ghost line in dynamic updates
- [ ] Ghost Y validated (< 20)
- [ ] Matrix nil checked
- [ ] Piece Y validated (0-19)
- [ ] Screen bounds checked before SetContent
- [ ] No app freeze with ghost enabled
- [ ] Ghost toggles on/off correctly
