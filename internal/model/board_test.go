package model

import (
	"testing"
)

func TestBoardInfrastructure_Pattern(t *testing.T) {
	b := NewBoard()
	if b == nil {
		t.Fatal("NewBoard() returned nil")
	}
}

func TestNewBoard(t *testing.T) {
	b := NewBoard()
	for y := 0; y < 20; y++ {
		for x := 0; x < 10; x++ {
			if b.Get(x, y) != 0 {
				t.Errorf("NewBoard() cell (%d,%d) = %d, want 0", x, y, b.Get(x, y))
			}
		}
	}
}

func TestBoardSetGet(t *testing.T) {
	b := NewBoard()
	b.Set(5, 10, 3)
	if got := b.Get(5, 10); got != 3 {
		t.Errorf("Get(5,10) = %d, want 3", got)
	}
}

func TestBoardBounds(t *testing.T) {
	b := NewBoard()
	if !b.IsWithinBounds(0, 0) {
		t.Error("IsWithinBounds(0,0) = false, want true")
	}
	if !b.IsWithinBounds(9, 19) {
		t.Error("IsWithinBounds(9,19) = false, want true")
	}
	if b.IsWithinBounds(-1, 0) {
		t.Error("IsWithinBounds(-1,0) = true, want false")
	}
	if b.IsWithinBounds(10, 0) {
		t.Error("IsWithinBounds(10,0) = true, want false")
	}
	if b.IsWithinBounds(0, -1) {
		t.Error("IsWithinBounds(0,-1) = true, want false")
	}
	if b.IsWithinBounds(0, 20) {
		t.Error("IsWithinBounds(0,20) = true, want false")
	}
}

func TestBoardIsEmpty(t *testing.T) {
	b := NewBoard()
	if !b.IsEmpty(5, 10) {
		t.Error("IsEmpty(5,10) = false, want true")
	}
	b.Set(5, 10, 1)
	if b.IsEmpty(5, 10) {
		t.Error("IsEmpty(5,10) after Set = true, want false")
	}
}

func TestBoardClear(t *testing.T) {
	b := NewBoard()
	b.Set(5, 10, 3)
	b.Clear(5, 10)
	if b.Get(5, 10) != 0 {
		t.Errorf("Clear(5,10) failed, got %d, want 0", b.Get(5, 10))
	}
}

func TestBoardLineOperations(t *testing.T) {
	b := NewBoard()
	for x := 0; x < 10; x++ {
		b.Set(x, 10, 1)
	}
	if !b.IsLineFull(10) {
		t.Error("IsLineFull(10) = false, want true")
	}
	b.ClearLine(10)
	if b.IsLineFull(10) {
		t.Error("After ClearLine, IsLineFull(10) = true, want false")
	}
}

func TestBoardLineDrop(t *testing.T) {
	b := NewBoard()
	// Fill line 5 and put a block at line 6 (above line 5)
	for x := 0; x < 10; x++ {
		b.Set(x, 5, 1)
		b.Set(x, 6, 3)
	}

	// Clear line 5 - line 6 should drop into line 5
	b.ClearLine(5)

	// Line 5 should now have what was in line 6 (color 3)
	if b.Get(0, 5) != 3 {
		t.Errorf("After clearing line 5, line 5 should have color from line 6, got %d", b.Get(0, 5))
	}

	// Line 6 should now have what was in line 7 (empty)
	if !b.IsEmpty(0, 6) {
		t.Error("Line 6 should be empty after shift")
	}

	// Top line should be empty
	if !b.IsEmpty(0, 19) {
		t.Error("Top line should be empty after shift")
	}
}
