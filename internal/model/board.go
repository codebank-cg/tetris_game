package model

// Board represents the 10x20 Tetris playing field.
type Board struct {
	cells [20][10]int // 0=empty, 1-7=color
}

// NewBoard creates an empty board.
func NewBoard() *Board {
	return &Board{}
}

// Get returns the cell value at (x, y).
func (b *Board) Get(x, y int) int {
	if !b.IsWithinBounds(x, y) {
		return 0
	}
	return b.cells[y][x]
}

// Set sets the cell value at (x, y).
func (b *Board) Set(x, y, val int) {
	if !b.IsWithinBounds(x, y) {
		return
	}
	b.cells[y][x] = val
}

// Clear clears the cell at (x, y).
func (b *Board) Clear(x, y int) {
	if !b.IsWithinBounds(x, y) {
		return
	}
	b.cells[y][x] = 0
}

// IsWithinBounds checks if coordinates are valid.
func (b *Board) IsWithinBounds(x, y int) bool {
	return x >= 0 && x < 10 && y >= 0 && y < 20
}

// IsEmpty checks if cell is empty.
func (b *Board) IsEmpty(x, y int) bool {
	return b.Get(x, y) == 0
}

// IsFull checks if board is completely full.
func (b *Board) IsFull() bool {
	for y := 0; y < 20; y++ {
		for x := 0; x < 10; x++ {
			if b.cells[y][x] == 0 {
				return false
			}
		}
	}
	return true
}

// IsLineFull checks if a specific line is full.
func (b *Board) IsLineFull(y int) bool {
	if y < 0 || y >= 20 {
		return false
	}
	for x := 0; x < 10; x++ {
		if b.cells[y][x] == 0 {
			return false
		}
	}
	return true
}

// ClearLine removes a line and shifts all lines above down.
func (b *Board) ClearLine(y int) {
	if y < 0 || y >= 20 {
		return
	}
	// Shift lines down: lines above (higher y) move down to fill the gap
	for row := y; row < 19; row++ {
		b.cells[row] = b.cells[row+1]
	}
	// Clear the top line (y=19)
	for x := 0; x < 10; x++ {
		b.cells[19][x] = 0
	}
}

// ClearLines removes all lines at the given y-coordinates in a single pass.
// This is the correct multi-line clear: it avoids index shifting bugs that
// occur when calling ClearLine repeatedly for adjacent rows.
func (b *Board) ClearLines(lines []int) {
	if len(lines) == 0 {
		return
	}
	// Build a set of rows to remove for O(1) lookup
	toRemove := make(map[int]bool, len(lines))
	for _, y := range lines {
		toRemove[y] = true
	}
	// Collect surviving rows (bottom to top order)
	surviving := make([][10]int, 0, 20)
	for y := 0; y < 20; y++ {
		if !toRemove[y] {
			surviving = append(surviving, b.cells[y])
		}
	}
	// Rebuild board: surviving rows at bottom, empty rows at top
	for y := 0; y < 20; y++ {
		if y < len(surviving) {
			b.cells[y] = surviving[y]
		} else {
			b.cells[y] = [10]int{}
		}
	}
}
