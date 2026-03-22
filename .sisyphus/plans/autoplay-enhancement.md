# Autoplay Algorithm Enhancement - Two-Piece Lookahead

## TL;DR

> **Quick Summary**: Enhance the autoplay AI to better utilize the Next piece preview for strategic two-piece planning
> 
> **Deliverables**:
> - Improved `EvaluateTwoPieceSequence()` with better combo weighting
> - Enhanced scoring that prioritizes setup plays for multi-line clears
> - Better balance between current move and next piece planning
> 
> **Estimated Effort**: Quick
> **Parallel Execution**: NO - sequential (single file change)
> **Critical Path**: EvaluateTwoPieceSequence modification → Build & Test

---

## Context

### Original Request
User wants the autoplay algorithm to better consider the 'Next' block by pre-calculating positions for both current and next pieces together.

### Current Implementation
The `FindBestMoveWithNext()` function in `autoplay.go` already has two-piece lookahead via `EvaluateTwoPieceSequence()`, but the scoring is basic:
- Only gives 1.5× bonus for combo clears
- Equal weighting between current and next piece scores
- No special bonus for Tetris (4-line) setups

### Interview Summary
**Key Discussions**:
- Algorithm should strategically plan across two pieces
- Prioritize setups that enable multi-line clears
- Better weighting of future board state

---

## Work Objectives

### Core Objective
Improve the two-piece evaluation to make smarter strategic decisions that set up better plays with the next piece.

### Concrete Deliverables
- Enhanced `EvaluateTwoPieceSequence()` function in `/Users/gangchen/works/oc_garden/tetris_game/internal/model/autoplay.go`

### Definition of Done
- Build passes: `go build -o tetris ./cmd/tetris`
- Tests pass: `go test ./internal/model/...`
- Autoplay demonstrates better setup behavior when tested manually

### Must Have
- Maintain backward compatibility (no API changes)
- Keep existing function signatures
- Preserve all existing test cases

### Must NOT Have (Guardrails)
- Do not change the MoveDecision struct
- Do not modify single-piece evaluation logic
- Do not break existing autoplay tests

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: YES
- **Automated tests**: YES (after) - Run existing tests
- **Framework**: Go standard testing
- **Agent-Executed QA**: Run autoplay integration tests

### QA Policy
- Run existing unit tests
- Manual QA: Observe autoplay behavior in-game

---

## Execution Strategy

### Sequential Steps

```
Step 1: Enhance EvaluateTwoPieceSequence() [deep]
├── Increase combo clear bonuses (2×-3× for multi-line setups)
├── Weight next piece score more heavily (60% next, 40% current)
├── Add special Tetris (4-line) setup bonus
└── Track best combo lines for extra scoring

Step 2: Build & Test [quick]
├── go build -o tetris ./cmd/tetris
├── go test ./internal/model/...
└── Manual test: Run autoplay mode and observe behavior
```

---

## TODOs

> Implementation + Test = ONE Task. Never separate.
> EVERY task MUST have: Recommended Agent Profile + Parallelization info + QA Scenarios.

- [ ] 1. Enhance EvaluateTwoPieceSequence() with better strategic weighting

  **What to do**:
  - Modify the combo clear bonus from 1.5× to tiered bonuses:
    - 4-line setup: 3.0× bonus
    - 3-line setup: 2.5× bonus
    - 2-line setup: 2.0× bonus
    - 1-line: 1.2× bonus
  - Change score weighting: 60% next piece + 40% current piece (instead of 50/50)
  - Track bestComboLines and add extra bonus when setup leads to multi-line clears
  - Add `bestComboLines` variable to track the best combo found

  **Must NOT do**:
  - Do not change function signature
  - Do not modify other evaluation functions
  - Do not break existing code structure

  **Recommended Agent Profile**:
  - **Category**: `deep` - Requires understanding game theory and scoring trade-offs
  - **Skills**: None needed - standard Go editing

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Sequential**: Must complete before testing
  - **Blocks**: Step 2 (Build & Test)
  - **Blocked By**: None

  **References**:
  - `internal/model/autoplay.go:446-520` - Current EvaluateTwoPieceSequence implementation
  - `internal/model/autoplay.go:195-209` - evaluateLineClears() function showing bonus ratios
  - `internal/model/autoplay.go:300-310` - Heuristic weights for context

  **Acceptance Criteria**:
  - [ ] Code compiles without errors
  - [ ] All existing tests pass
  - [ ] Combo bonus logic handles all cases (1, 2, 3, 4 lines)

  **QA Scenarios**:

  ```
  Scenario: Build verification
    Tool: Bash
    Preconditions: In project root directory
    Steps:
      1. Run: go build -o tetris ./cmd/tetris
      2. Check exit code is 0
      3. Verify tetris binary exists
    Expected Result: Build succeeds with exit code 0
    Evidence: .sisyphus/evidence/task-1-build-output.txt

  Scenario: Unit tests pass
    Tool: Bash
    Preconditions: Build successful
    Steps:
      1. Run: go test -v ./internal/model/...
      2. Check all tests pass
    Expected Result: All tests pass (0 failures)
    Evidence: .sisyphus/evidence/task-1-test-output.txt

  Scenario: Autoplay runs without errors
    Tool: interactive_bash (tmux)
    Preconditions: Binary built
    Steps:
      1. Run: timeout 5 go run ./cmd/tetris
      2. Send 'a' key to enable autoplay
      3. Wait 3 seconds
      4. Check process exits cleanly
    Expected Result: Autoplay mode activates, no panics/errors
    Evidence: .sisyphus/evidence/task-1-autoplay-run.txt
  ```

  **Commit**: YES
  - Message: `enhance(autoplay): improve two-piece lookahead scoring for better setup plays`
  - Files: `internal/model/autoplay.go`
  - Pre-commit: `go test ./internal/model/...`

---

## Final Verification Wave

- [ ] F1. **Plan Compliance Audit** — `oracle`
  Verify EvaluateTwoPieceSequence has: tiered combo bonuses, 60/40 weighting, bestComboLines tracking

- [ ] F2. **Code Quality Review** — `unspecified-high`
  Run `go build` + `go vet` + `go test ./...`

- [ ] F3. **Real Manual QA** — `unspecified-high`
  Test autoplay mode for 30 seconds, verify it makes strategic plays

- [ ] F4. **Scope Fidelity Check** — `deep`
  Verify only EvaluateTwoPieceSequence was modified, no other functions changed

---

## Commit Strategy

- **1**: `enhance(autoplay): improve two-piece lookahead scoring for better setup plays`
  - internal/model/autoplay.go
  - `go test ./internal/model/...`

---

## Success Criteria

### Verification Commands
```bash
go build -o tetris ./cmd/tetris      # Expected: success
go test ./internal/model/...         # Expected: all pass
./tetris                             # Manual: autoplay shows strategic behavior
```

### Final Checklist
- [ ] Combo bonus tiered (1.2×, 2.0×, 2.5×, 3.0×)
- [ ] Score weighting 60/40 (next/current)
- [ ] bestComboLines tracked and bonus applied
- [ ] All tests pass
- [ ] Autoplay runs without errors
