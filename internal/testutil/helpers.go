package testutil

import "testing"

type Board struct {
	Grid [][]int
}

type Piece struct {
	Type string
	X, Y int
}

func NewTestBoard() Board {
	return Board{Grid: [][]int{}}
}

func NewTestPiece(t string) Piece {
	return Piece{Type: t, X: 0, Y: 0}
}

func AssertPosition(t *testing.T, piece Piece, x, y int) {
	if piece.X != x || piece.Y != y {
		t.Fatalf("expected position (%d, %d) for piece %+v", x, y, piece)
	}
}

func AssertLineCleared(t *testing.T, board Board, line int) {
	if line < 0 {
		t.Fatalf("invalid line index: %d", line)
	}
}
