package model

import "testing"

func TestGameStateCreation(t *testing.T) {
	gs := NewGameState()
	if gs == nil {
		t.Fatal("NewGameState() returned nil")
	}
	if gs.CurrentPiece == nil {
		t.Error("CurrentPiece should not be nil")
	}
	if gs.NextPiece == nil {
		t.Error("NextPiece should not be nil")
	}
}

func TestEdgeCollision(t *testing.T) {
	gs := NewGameState()
	initialX := gs.CurrentPiece.X

	// Try to move left past left edge
	for i := 0; i < 10; i++ {
		gs.MovePiece(-1, 0)
	}
	x, _ := gs.CurrentPiece.GetPosition()
	if x < 0 {
		t.Errorf("Piece went past left edge: x=%d", x)
	}

	// Reset piece position
	gs.CurrentPiece.X = initialX

	// Try to move right past right edge
	for i := 0; i < 15; i++ {
		gs.MovePiece(1, 0)
	}
	x, _ = gs.CurrentPiece.GetPosition()
	if x > 9 {
		t.Errorf("Piece went past right edge: x=%d", x)
	}
}

func TestMoveBlockedByExistingBlocks(t *testing.T) {
	gs := NewGameState()

	// Block the area where the piece would land
	// Since piece type is random, block a wide area at multiple heights
	for y := 0; y < 18; y++ {
		for x := 0; x < 10; x++ {
			gs.Board.Set(x, y, 1)
		}
	}

	// Piece should not be able to move down at all from spawn
	moved := gs.MovePiece(0, -1)
	if moved {
		t.Error("Move down should be blocked by existing blocks")
	}
}

func TestRotateBlockedByWall(t *testing.T) {
	gs := NewGameState()

	// Move piece to right edge
	gs.CurrentPiece.X = 9

	// Try to rotate - should fail if piece would overlap wall
	gs.RotatePiece()

	// Piece should either not rotate or adjust position to fit
	x, _ := gs.CurrentPiece.GetPosition()
	if x > 9 {
		t.Errorf("Piece rotated outside bounds: x=%d", x)
	}
}
