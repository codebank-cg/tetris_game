# Tetris Code Cleanup Plan

## TL;DR

> **Quick Summary**: Remove empty placeholder files, unused UI package files, and consolidate redundant code to simplify the codebase
> 
> **Deliverables**:
> - Remove 4 empty placeholder files
> - Remove unused ui package (all rendering in main.go)
> - Remove unused assets package
> - Clean up unused helper functions in model package
> 
> **Estimated Effort**: Quick
> **Parallel Execution**: NO - sequential file deletions
> **Critical Path**: Backup → Delete → Build → Verify

---

## Context

### Current State Analysis

**Empty Files (0 code):**
- `internal/assets/placeholder.go` - 1 line (package declaration only)
- `internal/ui/placeholder.go` - 1 line (package declaration only)
- `internal/model/placeholder.go` - 1 line (package declaration only)
- `internal/testutil/placeholder.go` - 1 line (package declaration only)

**Unused Packages:**
- `internal/ui/*` - All rendering now done in main.go via SetDrawFunc
- `internal/assets/*` - Not imported by main.go, UI elements rendered directly

**Unused Functions:**
- `getMaxHeight()` in autoplay.go - No external callers found
- `getColHeight()` in autoplay.go - Only used by getMaxHeight (also unused)

### Interview Summary

**Key Findings**:
- Development process left placeholder files
- UI rendering migrated to main.go but old files kept
- Assets package not used in current implementation
- No breaking changes - all code is internal/dead

---

## Work Objectives

### Core Objective
Remove all dead code and unused files to create a cleaner, more maintainable codebase.

### Concrete Deliverables
- Zero empty placeholder files
- Zero unused packages
- Zero dead functions
- Clean build with no warnings

### Definition of Done
- All identified files removed
- `go build` succeeds
- `go test ./...` passes
- Binary runs without errors

### Must Have
- Remove all 4 placeholder files
- Remove entire internal/ui directory
- Remove internal/assets directory
- Remove getMaxHeight and getColHeight functions

### Must NOT Have
- No changes to working game logic
- No modifications to main.go (unless removing imports)
- No breaking changes to public APIs

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: YES
- **Automated tests**: Run `go test ./...`
- **Manual QA**: Run `./tetris` and play for 30 seconds

### QA Policy
- Automated: All existing tests must pass
- Manual: Game runs normally

---

## Execution Strategy

### Sequential Steps

```
Step 1: Create backups [quick]
├── cp -r cmd/tetris/main.go cmd/tetris/main.go~
└── Create backup before any changes

Step 2: Remove empty placeholder files [quick]
├── rm internal/assets/placeholder.go
├── rm internal/ui/placeholder.go
├── rm internal/model/placeholder.go
└── rm internal/testutil/placeholder.go

Step 3: Remove unused packages [quick]
├── rm -rf internal/ui/
└── rm -rf internal/assets/

Step 4: Clean up unused functions [quick]
├── Remove getMaxHeight() from autoplay.go
├── Remove getColHeight() from autoplay.go
└── Update any code that calls these functions

Step 5: Build and test [quick]
├── go build -o tetris ./cmd/tetris
├── go test ./...
└── ./tetris (manual test)
```

---

## TODOs

- [ ] 1. Create backup of main.go

  **What to do**:
  ```bash
  cp cmd/tetris/main.go cmd/tetris/main.go~
  ```

  **Must NOT do**:
  - Do not skip backup step
  - Do not proceed without verified backup

  **Recommended Agent Profile**:
  - **Category**: `quick` - Simple file operation
  - **Skills**: None needed

  **Parallelization**:
  - **Sequential**: Must complete first

  **Acceptance Criteria**:
  - [ ] Backup file exists: `cmd/tetris/main.go~`
  - [ ] Backup matches original: `diff main.go main.go~` shows no differences

  **Commit**: NO (temporary backup)

- [ ] 2. Remove empty placeholder files

  **What to do**:
  ```bash
  rm internal/assets/placeholder.go
  rm internal/ui/placeholder.go
  rm internal/model/placeholder.go
  rm internal/testutil/placeholder.go
  ```

  **Must NOT do**:
  - Do not remove non-placeholder files
  - Do not remove directories yet

  **Recommended Agent Profile**:
  - **Category**: `quick` - File deletion
  - **Skills**: None needed

  **Parallelization**:
  - **Sequential**: After backup

  **Acceptance Criteria**:
  - [ ] All 4 placeholder files deleted
  - [ ] No other files affected

  **Commit**: YES
  - Message: `chore: remove empty placeholder files`
  - Files: `internal/*/placeholder.go`

- [ ] 3. Remove unused packages (ui, assets)

  **What to do**:
  ```bash
  rm -rf internal/ui/
  rm -rf internal/assets/
  ```

  **Why safe to remove**:
  - UI rendering: All done in main.go via SetDrawFunc
  - Assets: Not imported anywhere
  - No external dependencies on these packages

  **Must NOT do**:
  - Do not remove other internal packages
  - Do not modify main.go yet

  **Recommended Agent Profile**:
  - **Category**: `quick` - Directory deletion
  - **Skills**: None needed

  **Parallelization**:
  - **Sequential**: After Step 2

  **Acceptance Criteria**:
  - [ ] internal/ui/ directory removed
  - [ ] internal/assets/ directory removed
  - [ ] `go build` still succeeds

  **Commit**: YES
  - Message: `chore: remove unused ui and assets packages`
  - Files: `internal/ui/*`, `internal/assets/*`

- [ ] 4. Remove unused helper functions

  **What to do**:
  - Open `internal/model/autoplay.go`
  - Remove `getMaxHeight()` function
  - Remove `getColHeight()` function
  - Verify no callers remain

  **Must NOT do**:
  - Do not remove other helper functions
  - Do not change function signatures

  **Recommended Agent Profile**:
  - **Category**: `quick` - Code cleanup
  - **Skills**: None needed

  **Parallelization**:
  - **Sequential**: After Step 3

  **References**:
  - `internal/model/autoplay.go` - Contains functions to remove

  **Acceptance Criteria**:
  - [ ] getMaxHeight() removed
  - [ ] getColHeight() removed
  - [ ] Build succeeds
  - [ ] No test failures

  **Commit**: YES
  - Message: `refactor: remove unused helper functions from autoplay.go`
  - Files: `internal/model/autoplay.go`

- [ ] 5. Build, test, and verify

  **What to do**:
  ```bash
  go build -o tetris ./cmd/tetris
  go test ./...
  ./tetris  # Manual test
  ```

  **QA Scenarios**:
  ```
  Scenario: Build succeeds
    Tool: Bash
    Steps:
      1. Run: go build -o tetris ./cmd/tetris
      2. Check exit code is 0
      3. Verify tetris binary exists
    Expected Result: Build succeeds with exit code 0
    Evidence: .sisyphus/evidence/build-output.txt

  Scenario: All tests pass
    Tool: Bash
    Steps:
      1. Run: go test ./...
      2. Check all packages pass
    Expected Result: 0 failures
    Evidence: .sisyphus/evidence/test-output.txt

  Scenario: Manual gameplay test
    Tool: interactive_bash (tmux)
    Steps:
      1. Run: timeout 30 ./tetris
      2. Play game for 30 seconds
      3. Check process exits cleanly
    Expected Result: Game runs smoothly, no crashes
    Evidence: .sisyphus/evidence/manual-test.txt
  ```

  **Acceptance Criteria**:
  - [ ] Build succeeds
  - [ ] All tests pass
  - [ ] Game runs without errors
  - [ ] No new warnings from `go vet`

  **Commit**: YES
  - Message: `chore: verify cleanup - build and test pass`

---

## Final Verification Wave

- [ ] F1. **Cleanup Audit** — `oracle`
  Verify all 4 placeholders removed, ui/assets packages gone, no dead functions

- [ ] F2. **Build Verification** — `quick`
  `go build` succeeds, no warnings

- [ ] F3. **Test Suite** — `unspecified-high`
  `go test ./...` passes completely

- [ ] F4. **Manual QA** — `unspecified-high`
  Game runs for 60 seconds, no crashes

- [ ] F5. **Code Quality** — `met is`
  Run `go vet ./...` - zero warnings

---

## Commit Strategy

**Sequential commits for safety:**

1. `chore: create backup before cleanup` (backup file)
2. `chore: remove empty placeholder files` (4 files)
3. `chore: remove unused ui and assets packages` (2 directories)
4. `refactor: remove unused helper functions` (autoplay.go)
5. `chore: verify cleanup passes all tests` (verification)

---

## Success Criteria

### File Count Reduction

| Before | After | Removed |
|--------|-------|---------|
| ~24 .go files | ~18 .go files | 6 files |
| 4 empty files | 0 empty files | 4 files |
| 2 unused packages | 0 unused packages | 2 packages |

### Verification Commands
```bash
go build -o tetris ./cmd/tetris   # Expected: success
go test ./...                      # Expected: all pass
go vet ./...                       # Expected: no warnings
./tetris                           # Manual: runs normally
```

### Final Checklist
- [ ] All placeholders removed
- [ ] ui/ package removed
- [ ] assets/ package removed
- [ ] getMaxHeight() removed
- [ ] getColHeight() removed
- [ ] Build succeeds
- [ ] Tests pass
- [ ] Game runs
- [ ] No warnings
