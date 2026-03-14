# Fix Autoplay Algorithm - Multi-Line Clear Incentive

## TL;DR

> **Quick Summary**: The autoplay AI plays too conservatively because heuristic weights don't properly incentivize multi-line clears. Need to rebalance weights to favor Tetris (4-line) setups over safe, low-height play.
> 
> **Deliverables**: 
> - Updated heuristic weights in `autoplay.go`
> - Modified `evaluateBoard()` to apply exponential bonus for multi-line clears
> - Passing tests demonstrating improved line-clear rate
> 
> **Estimated Effort**: Short (30-60 min)
> **Parallel Execution**: NO - sequential (single file changes)
> **Critical Path**: Weight change → Evaluation function → Test verification

---

## Context

### Original Request
User reported: "the algorithm always leaves enough space for 'Hero' block, but does not chose the latter 'Hero' block to clear multi-lines. it looks strange."

### Interview Summary
**Key Findings from Code Analysis**:
- Current `completeLines` weight: 0.76 (linear scaling)
- Current `aggregateHeight` weight: -0.50 (heavy penalty)
- Algorithm evaluates board AFTER piece placement but BEFORE line clearing
- No exponential bonus for multiple lines (Tetris scoring uses 1→40, 2→100, 3→300, 4→1200)

### Research Findings
**Standard Tetris AI Heuristics** (from Dellacherie, Fahey, and other competitive AIs):
- Line clears should be the DOMINANT factor
- Common approach: Exponential or step-function bonus for multi-line clears
- Height penalty should be moderate, not dominant
- Typical weight ratios: Lines bonus should be 3-5× stronger than height penalty

---

## Work Objectives

### Core Objective
Rebalance autoplay heuristic weights to properly incentivize multi-line clears (especially Tetris/4-line) while maintaining reasonable safety (hole avoidance).

### Concrete Deliverables
1. Modified `autoplay.go` with:
   - Exponential line-clear bonus function
   - Rebalanced heuristic weights
   - Comment documenting the weight rationale
2. Updated test demonstrating improved line-clear rate (target: 15-20% vs current 11%)
3. All existing tests still passing

### Definition of Done
- [ ] `go test ./internal/model/...` passes all tests
- [ ] `TestAutoPlay_ClearsLines` shows ≥15% line-clear rate (currently ~11%)
- [ ] Manual observation: AI visibly sets up and executes Tetris moves
- [ ] Code builds: `go build ./...` succeeds

### Must Have
- Exponential or step-function bonus for multi-line clears (not linear)
- Reduced `aggregateHeight` penalty relative to line-clear bonus
- Maintain hole/well penalties to prevent self-destructive play

### Must NOT Have (Guardrails)
- No hardcoded piece-specific logic (keep general heuristic approach)
- No lookahead beyond current piece (that's a separate feature)
- No breaking changes to function signatures (maintain API compatibility)
- No removing existing heuristics entirely (tune, don't delete)

---

## Verification Strategy

> **ZERO HUMAN INTERVENTION** — ALL verification is agent-executed.

### Test Decision
- **Infrastructure exists**: YES (Go testing framework)
- **Automated tests**: TDD (tests first to establish baseline)
- **Framework**: `go test`

### QA Policy
Every task MUST include agent-executed QA scenarios.

- **Library/Module**: Use Bash (`go test`) — Run tests, assert pass/fail, capture output

---

## Execution Strategy

### Sequential Tasks (Single File Focus)

```
Wave 1 (Foundation — test baseline):
├── Task 1: Add test capturing current line-clear rate [quick]
└── Task 2: Add test for Tetris incentive analysis [quick]

Wave 2 (Implementation — weight changes):
├── Task 3: Implement exponential line-clear bonus [unspecified-high]
├── Task 4: Rebalance heuristic weights [unspecified-high]
└── Task 5: Update function documentation [writing]

Wave 3 (Verification):
├── Task 6: Run all tests, verify improvement [quick]
└── Task 7: Build and lint check [quick]

Critical Path: Task 1 → Task 3 → Task 4 → Task 6
```

### Agent Dispatch Summary

- **Wave 1**: 2 tasks → `quick` (test additions)
- **Wave 2**: 3 tasks → `unspecified-high` (algorithm changes), `writing` (docs)
- **Wave 3**: 2 tasks → `quick` (verification)

---

## TODOs

> Implementation + Test = ONE Task. Never separate.

- [ ] 1. Add Baseline Test for Line-Clear Rate

  **What to do**:
  - Add new test `TestAutoplay_LineClearRate` in `autoplay_integration_test.go`
  - Run 100 pieces, record line-clear percentage
  - Assert current baseline (~11%) is captured
  - Add comment noting this test will be updated after weight changes
  
  **Must NOT do**:
  - Don't modify any algorithm code yet
  - Don't change existing tests
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: None needed (simple test addition)
  
  **Parallelization**:
  - **Can Run In Parallel**: NO (foundational baseline)
  - **Blocked By**: None
  
  **Acceptance Criteria**:
  - [ ] Test file compiles: `go test ./internal/model/... -run TestAutoplay_LineClearRate -v`
  - [ ] Test logs line-clear percentage
  
  **QA Scenarios**:
  ```
  Scenario: Test runs and captures baseline
    Tool: Bash
    Preconditions: In project directory
    Steps:
      1. Run: go test ./internal/model -run TestAutoplay_LineClearRate -v
      2. Parse output for "Lines cleared" log line
      3. Extract percentage value
    Expected Result: Test passes, logs line-clear rate (expect ~11-12%)
    Evidence: .sisyphus/evidence/task-1-baseline.txt
  ```

- [ ] 2. Add Tetris Incentive Analysis Test

  **What to do**:
  - Add test `TestWeightAnalysis_TetrisIncentive` in `autoplay_test.go`
  - Calculate whether current weights favor Tetris over conservative play
  - Log scenario comparisons (conservative vs Tetris setup)
  - Assert whether Tetris is currently incentivized (expect: NO)
  
  **Must NOT do**:
  - Don't fix the weights yet
  - Don't modify production code
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: None needed
  
  **Parallelization**:
  - **Can Run In Parallel**: YES (with Task 1)
  - **Parallel Group**: Wave 1 (with Task 1)
  - **Blocked By**: None
  
  **Acceptance Criteria**:
  - [ ] Test compiles and runs
  - [ ] Test logs weight analysis
  - [ ] Test clearly shows Tetris is NOT incentivized (negative finding)
  
  **QA Scenarios**:
  ```
  Scenario: Test demonstrates current weight imbalance
    Tool: Bash
    Steps:
      1. Run: go test ./internal/model -run TestWeightAnalysis_TetrisIncentive -v
      2. Capture test output showing weight comparison
    Expected Result: Test logs show Tetris score < Conservative score
    Evidence: .sisyphus/evidence/task-2-weight-analysis.txt
  ```

- [ ] 3. Implement Exponential Line-Clear Bonus

  **What to do**:
  - Replace linear `countCompleteLines()` with `evaluateLineClears(lines int) float64`
  - Implement step-function bonus matching official Tetris scoring:
    - 0 lines: 0
    - 1 line: 40
    - 2 lines: 100
    - 3 lines: 300
    - 4 lines: 1200
  - Apply a scaling factor (e.g., 0.01) to keep scores in reasonable range
  - Update `evaluateBoard()` to use new function
  
  **Must NOT do**:
  - Don't change function signature of `evaluateBoard()`
  - Don't remove other heuristics (height, holes, bumpiness, wells)
  
  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: None needed (algorithm implementation)
  
  **Parallelization**:
  - **Can Run In Parallel**: NO (depends on baseline tests)
  - **Blocked By**: Tasks 1, 2
  
  **References**:
  - `autoplay.go:306-323` — Current `evaluateBoard()` implementation
  - `autoplay.go:181-190` — Current `countCompleteLines()` function
  - Official Tetris scoring: 1→40, 2→100, 3→300, 4→1200 (linear is WRONG)
  
  **Acceptance Criteria**:
  - [ ] New `evaluateLineClears()` function exists
  - [ ] Function returns exponentially scaled values
  - [ ] `evaluateBoard()` calls new function
  - [ ] Code compiles: `go build ./...`
  
  **QA Scenarios**:
  ```
  Scenario: Exponential bonus calculation
    Tool: Bash
    Steps:
      1. Add temp test calling evaluateLineClears(0), (1), (2), (3), (4)
      2. Run: go test -run TestEvaluateLineClears -v
      3. Verify outputs: 0, ~0.4, ~1.0, ~3.0, ~12.0 (scaled)
    Expected Result: Outputs show exponential growth, not linear
    Evidence: .sisyphus/evidence/task-3-exponential-test.txt
  ```

- [ ] 4. Rebalance Heuristic Weights

  **What to do**:
  - Update `heuristicWeights` map in `autoplay.go`:
    - `aggregateHeight`: -0.50 → -0.25 (reduce height penalty by 50%)
    - `completeLines`: 0.76 → remove (replaced by exponential function)
    - `holes`: -0.36 → -0.30 (slightly reduce, keep important)
    - `bumpiness`: -0.18 → -0.15 (minor reduction)
    - `wells`: -0.12 → -0.10 (minor reduction)
  - Add comment explaining weight rationale
  - Add comment documenting Tetris incentive calculation
  
  **Must NOT do**:
  - Don't make heights penalty positive (still want to discourage tall stacks)
  - Don't remove hole/well penalties (causes self-blocking)
  - Don't make weights extreme (causes erratic behavior)
  
  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: None needed
  
  **Parallelization**:
  - **Can Run In Parallel**: NO (depends on Task 3)
  - **Blocked By**: Task 3
  
  **References**:
  - `autoplay.go:281-288` — Current `heuristicWeights` map
  - Dellacherie AI weights (research standard): Lines weight ~5× height penalty
  
  **Acceptance Criteria**:
  - [ ] Weights updated with new values
  - [ ] Comment explains Tetris incentive math
  - [ ] Code compiles: `go build ./...`
  
  **QA Scenarios**:
  ```
  Scenario: Weight verification
    Tool: Bash
    Steps:
      1. Run: go test -run TestGetWeights -v
      2. Verify weights match expected values
      3. Check comments exist explaining rationale
    Expected Result: Weights updated, comments present
    Evidence: .sisyphus/evidence/task-4-weights.txt
  ```

- [ ] 5. Update Function Documentation

  **What to do**:
  - Add doc comment to `evaluateBoard()` explaining the heuristic approach
  - Add comment above `heuristicWeights` explaining:
    - Why exponential line bonus
    - Weight balance rationale
    - Expected behavior (Tetris-seeking vs conservative)
  - Update package-level comment if needed
  
  **Must NOT do**:
  - Don't add excessive inline comments (code should be clear)
  - Don't document obvious things (focus on WHY, not WHAT)
  
  **Recommended Agent Profile**:
  - **Category**: `writing`
  - **Skills**: None needed
  
  **Parallelization**:
  - **Can Run In Parallel**: YES (with Task 4)
  - **Parallel Group**: Wave 2 (with Task 4)
  - **Blocked By**: Task 3
  
  **Acceptance Criteria**:
  - [ ] `evaluateBoard()` has clear doc comment
  - [ ] `heuristicWeights` has rationale comment
  - [ ] Comments pass `go vet` (no malformed comments)
  
  **QA Scenarios**:
  ```
  Scenario: Documentation check
    Tool: Bash
    Steps:
      1. Run: go vet ./internal/model/...
      2. Run: go doc internal/model.evaluateBoard
      3. Verify documentation is readable and accurate
    Expected Result: go vet passes, documentation exists
    Evidence: .sisyphus/evidence/task-5-docs.txt
  ```

- [ ] 6. Run All Tests and Verify Improvement

  **What to do**:
  - Run full test suite: `go test ./internal/model/... -v`
  - Verify all existing tests pass
  - Check `TestAutoPlay_ClearsLines` shows improved rate (target: ≥15%)
  - Run `TestAutoPlay_Survival50Pieces` — should survive better with aggressive clearing
  - Update baseline test expectation if significantly improved
  
  **Must NOT do**:
  - Don't adjust test expectations downward if improvement is modest
  - Don't ignore failing tests
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: None needed
  
  **Parallelization**:
  - **Can Run In Parallel**: NO (depends on all implementation)
  - **Blocked By**: Tasks 3, 4, 5
  
  **Acceptance Criteria**:
  - [ ] All tests pass: `go test ./internal/model/...`
  - [ ] Line-clear rate improved by ≥4 percentage points (11% → ≥15%)
  - [ ] `TestAutoPlay_Survival50Pieces` passes (survival maintained)
  
  **QA Scenarios**:
  ```
  Scenario: Full test suite passes with improvement
    Tool: Bash
    Steps:
      1. Run: go test ./internal/model/... -v 2>&1 | tee /tmp/test-output.txt
      2. Grep for "Lines cleared" in output
      3. Extract final percentage
      4. Verify all PASS results
    Expected Result: All tests PASS, line-clear rate ≥15%
    Evidence: .sisyphus/evidence/task-6-test-results.txt
  ```

- [ ] 7. Build and Lint Check

  **What to do**:
  - Build: `go build ./...`
  - Run: `go vet ./...`
  - Run: `go fmt ./...`
  - Verify binary works: `./tetris` (manual, just check it compiles)
  
  **Must NOT do**:
  - Don't skip any verification steps
  
  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: None needed
  
  **Parallelization**:
  - **Can Run In Parallel**: YES (with Task 6)
  - **Parallel Group**: Wave 3 (with Task 6)
  - **Blocked By**: Tasks 3, 4, 5
  
  **Acceptance Criteria**:
  - [ ] `go build ./...` succeeds (no errors)
  - [ ] `go vet ./...` passes (no warnings)
  - [ ] `go fmt ./...` makes no changes (code is formatted)
  
  **QA Scenarios**:
  ```
  Scenario: Build and lint pass
    Tool: Bash
    Steps:
      1. Run: go build ./... 2>&1
      2. Run: go vet ./... 2>&1
      3. Run: gofmt -d . | head -20
    Expected Result: All commands succeed with no output/errors
    Evidence: .sisyphus/evidence/task-7-build-lint.txt
  ```

---

## Final Verification Wave

- [ ] F1. **Plan Compliance Audit** — `oracle`
  Verify all "Must Have" features implemented, "Must NOT Have" respected, evidence files exist.

- [ ] F2. **Code Quality Review** — `unspecified-high`
  Run `go build`, `go vet`, `go test`. Check for `as any`, empty catches, unused imports.

- [ ] F3. **Real Manual QA** — `unspecified-high`
  Run game with autoplay enabled (`a` key), observe for 50+ pieces that AI:
  - Sets up multi-line clears visibly
  - Executes Tetris (4-line) moves when possible
  - Doesn't self-block with holes

- [ ] F4. **Scope Fidelity Check** — `deep`
  Verify only `autoplay.go` and test files modified, no unrelated changes.

---

## Commit Strategy

- **1**: `refactor(autoplay): implement exponential line-clear bonus` — autoplay.go, autoplay_test.go, autoplay_integration_test.go
- **2**: `chore(autoplay): update heuristic weights for Tetris incentive` — autoplay.go
- **3**: `docs(autoplay): add heuristic rationale comments` — autoplay.go
  - All in one commit: `refactor(autoplay): rebalance AI weights for multi-line clears`

---

## Success Criteria

### Verification Commands
```bash
go test ./internal/model/... -v                           # Expected: All PASS
go test ./internal/model/... -run TestAutoPlay_ClearsLines # Expected: ≥15% line rate
go build ./...                                            # Expected: Success
go vet ./...                                              # Expected: No warnings
```

### Final Checklist
- [ ] All "Must Have" present (exponential bonus, rebalanced weights)
- [ ] All "Must NOT Have" absent (no hardcoded piece logic, no breaking changes)
- [ ] All tests pass
- [ ] Line-clear rate improved from ~11% to ≥15%
- [ ] Code compiles and runs
- [ ] Documentation explains weight rationale

### Success Metrics

**Before Fix** (baseline):
- Line-clear rate: ~11% (6 lines / 53 pieces)
- Tetris (4-line) incentive: Negative (algorithm prefers conservative play)
- `completeLines` weight: 0.76 (linear)

**After Fix** (target):
- Line-clear rate: ≥15% (improved by ≥4 percentage points)
- Tetris (4-line) incentive: Positive (algorithm seeks multi-line setups)
- `completeLines` function: Exponential (0, 40, 100, 300, 1200 scaled)

---

## Risk Mitigation

### Potential Issues

1. **Over-correction: AI becomes too aggressive**
   - Mitigation: Conservative weight changes first, iterate if needed
   - Test: Verify `TestAutoPlay_Survival50Pieces` still passes

2. **Line-clear rate doesn't improve significantly**
   - Mitigation: Adjust scaling factor on exponential bonus
   - Test: Run 200+ piece simulation to reduce variance

3. **Breaking existing tests**
   - Mitigation: Run full test suite after each change
   - Test: TDD approach — baseline tests first

### Rollback Plan
If algo becomes worse:
1. Revert weight changes to original values
2. Keep exponential bonus (it's more correct than linear)
3. Adjust scaling factor rather than abandoning approach
