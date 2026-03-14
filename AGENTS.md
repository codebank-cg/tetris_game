# AGENTS.md - Development Guidelines

## Project Overview

**Go Tetris Game** - Terminal-based Tetris using tcell v2 library
- **Go Version**: 1.24.0+ (tested on go1.26.0)
- **Module**: `github.com/oc-garden/tetris_game`
- **Platform**: macOS/Linux/Windows (256-color terminal required)

---

## Build & Test Commands

### Build
```bash
go mod tidy                          # Sync dependencies
go build -o tetris                   # Build binary
go run ./cmd/tetris                  # Run from source
```

### Test
```bash
go test ./...                        # Run all tests
go test ./internal/model/...         # Run model package tests
go test -v ./internal/model/...      # Verbose test output
go test -run TestBoard ./...         # Run tests matching pattern
go test -race ./...                  # Run with race detector
```

### Lint & Format
```bash
go fmt ./...                         # Format all files (standard Go formatter)
go vet ./...                         # Static analysis
gofmt -d .                           # Show formatting diffs
```

**Note**: No custom linter configured. Standard Go tooling only.

---

## Code Style Guidelines

### Imports
- Standard library imports first (alphabetically)
- External packages second (alphabetically)
- Blank imports (`_ "package"`) only for side effects
```go
import (
    "fmt"
    "time"
    
    "github.com/gdamore/tcell/v2"
)
```

### Naming Conventions
- **Packages**: lowercase, single word (`model`, `ui`, `assets`)
- **Types**: PascalCase (`Tetromino`, `GameState`, `Board`)
- **Functions**: PascalCase (`NewBoard()`, `RotateClockwise()`)
- **Variables**: camelCase (`currentPiece`, `nextPiece`)
- **Constants**: PascalCase with type suffix (`TetrominoI`, `TetrominoO`)
- **Interfaces**: `-er` suffix if single method (none currently)

### Types & Data Structures
- Use structs for game state (`Board`, `Tetromino`, `GameState`)
- Pointer receivers for methods that modify state
- Value receivers for read-only methods (currently all use pointers)
- Use typed constants for enums (`type TetrominoType string`)

### Error Handling
- Simple projects: `panic()` acceptable for initialization failures
- Production: return `(T, error)` for recoverable errors
- Check errors immediately after calls
```go
screen, err := tcell.NewScreen()
if err != nil {
    panic(err)  // Acceptable for game initialization
}
```

### Testing Patterns
- Test files: `*_test.go` in same package as code
- Test function naming: `Test[Functionality][Scenario]`
- Use table-driven tests for multiple cases
- Helper functions: lowercase, accept `*testing.T`
```go
func TestTetrominoRotation(t *testing.T) {
    tet := NewTetromino(TetrominoT)
    // ... test logic
    checkMatrix(t, "description", got, expected)
}
```

### Comments & Documentation
- Package comment: `// Package name provides...` (none currently)
- Exported identifiers: doc comment required
- Implementation comments: `//` for why, not what
- Use `//` single-line comments consistently

### File Organization
```
cmd/tetris/          # Main entry point
internal/model/      # Game logic (board, pieces, state)
internal/ui/         # UI rendering (placeholders)
internal/assets/     # Static assets (ASCII art)
internal/testutil/   # Test helpers
docs/                # Documentation
```

### Go-Specific Conventions
- **Constructors**: `New[Type]()` pattern (`NewBoard()`, `NewTetromino()`)
- **Getters**: `Get[Property]()` (`GetPosition()`, `GetMatrix()`)
- **Setters**: `Set[Property]()` (`Set()`) or direct field access
- ** receivers**: Pointer receivers for all methods (consistency)
- **Zero values**: Design structs to work with zero values when possible

---

## Architecture Notes

### Package Structure
- `main`: Entry point, game loop, input handling, rendering
- `model`: Core game logic (board, pieces, randomizer, game state)
- No circular dependencies (enforced by Go)

### Key Design Decisions
- **7-bag randomizer**: Fair piece distribution
- **256-color terminal**: tcell v2 for cross-platform support
- **Fixed board size**: 10×20 cells (standard Tetris)
- **Coordinate system**: (0,0) at bottom-left, Y increases upward

---

## Common Tasks

### Add New Feature
1. Implement in `internal/model/` (game logic)
2. Add tests: `*_test.go` with same package
3. Update `cmd/tetris/main.go` for UI/input
4. Run: `go test ./... && go build`

### Fix Bug
1. Reproduce with test case
2. Minimal fix (no refactoring during bugfix)
3. Verify: `go test ./... && go vet ./...`

### Run Single Test
```bash
go test -v -run TestBoardIsEmpty ./internal/model/
```

---

## Existing Rules

No Cursor rules (`.cursor/rules/`) or Copilot rules (`.github/copilot-instructions.md`) present.

---

## Notes for Agents

1. **No type suppression**: Never use `//nolint` or similar without strong justification
2. **Match existing patterns**: Follow the established code style (pointer receivers, naming)
3. **Tests required**: All new functionality needs test coverage
4. **Minimal changes**: Fix bugs with smallest possible change
5. **No commits**: Do not commit changes unless explicitly requested
6. **Verify**: Run `go test ./...` after any code change
