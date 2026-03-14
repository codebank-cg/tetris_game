package model

import "testing"

func TestTetrominoCreation(t *testing.T) {
	pieces := []TetrominoType{TetrominoI, TetrominoO, TetrominoT, TetrominoS, TetrominoZ, TetrominoJ, TetrominoL}
	for _, p := range pieces {
		tet := NewTetromino(p)
		if tet.Type != p {
			t.Errorf("NewTetromino(%v).Type = %v, want %v", p, tet.Type, p)
		}
		if tet.Color == 0 {
			t.Errorf("NewTetromino(%v).Color = 0, want non-zero", p)
		}
		if tet.Matrix == nil {
			t.Errorf("NewTetromino(%v).Matrix = nil, want matrix", p)
		}
	}
}

func TestTetrominoRotation(t *testing.T) {
	tet := NewTetromino(TetrominoT)
	initialMatrix := tet.GetMatrix()
	for i := 0; i < 4; i++ {
		tet.RotateClockwise()
	}
	finalMatrix := tet.GetMatrix()
	for i := range initialMatrix {
		for j := range initialMatrix[i] {
			if initialMatrix[i][j] != finalMatrix[i][j] {
				t.Error("After 4 clockwise rotations, matrix should be same as initial")
			}
		}
	}
}

func TestTetrominoCounterRotation(t *testing.T) {
	tet := NewTetromino(TetrominoT)
	initialMatrix := tet.GetMatrix()
	for i := 0; i < 4; i++ {
		tet.RotateCounterClockwise()
	}
	finalMatrix := tet.GetMatrix()
	for i := range initialMatrix {
		for j := range initialMatrix[i] {
			if initialMatrix[i][j] != finalMatrix[i][j] {
				t.Error("After 4 counter-clockwise rotations, matrix should be same as initial")
			}
		}
	}
}

func TestTetrominoMove(t *testing.T) {
	tet := NewTetromino(TetrominoI)
	tet.Move(2, -3)
	x, y := tet.GetPosition()
	if x != 5 {
		t.Errorf("Move(2,-3) X = %d, want 5", x)
	}
	if y != 15 {
		t.Errorf("Move(2,-3) Y = %d, want 15", y)
	}
}

func TestTetrominoColor(t *testing.T) {
	expectedColors := map[TetrominoType]int{
		TetrominoI: 1,
		TetrominoO: 2,
		TetrominoT: 3,
		TetrominoS: 4,
		TetrominoZ: 5,
		TetrominoJ: 6,
		TetrominoL: 7,
	}
	for piece, expectedColor := range expectedColors {
		tet := NewTetromino(piece)
		if tet.Color != expectedColor {
			t.Errorf("NewTetromino(%v).Color = %d, want %d", piece, tet.Color, expectedColor)
		}
	}
}

func TestTPieceRotation(t *testing.T) {
	tet := NewTetromino(TetrominoT)

	// State 0: T pointing up
	// . T . .
	// T T T .
	// . . . .
	// . . . .
	expected0 := [][]int{
		{0, 1, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	checkMatrix(t, "T state 0 (up)", tet.GetMatrix(), expected0)

	// State 1: T pointing right
	// . T . .
	// . T T .
	// . T . .
	// . . . .
	tet.RotateClockwise()
	expected1 := [][]int{
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
	}
	checkMatrix(t, "T state 1 (right)", tet.GetMatrix(), expected1)

	// State 2: T pointing down
	// . . . .
	// T T T .
	// . T . .
	// . . . .
	tet.RotateClockwise()
	expected2 := [][]int{
		{0, 0, 0, 0},
		{1, 1, 1, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
	}
	checkMatrix(t, "T state 2 (down)", tet.GetMatrix(), expected2)

	// State 3: T pointing left
	// . T . .
	// T T . .
	// . T . .
	// . . . .
	tet.RotateClockwise()
	expected3 := [][]int{
		{0, 1, 0, 0},
		{1, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
	}
	checkMatrix(t, "T state 3 (left)", tet.GetMatrix(), expected3)
}

func checkMatrix(t *testing.T, name string, got, expected [][]int) {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if got[row][col] != expected[row][col] {
				t.Errorf("%s: matrix[%d][%d] = %d, want %d", name, row, col, got[row][col], expected[row][col])
			}
		}
	}
}
