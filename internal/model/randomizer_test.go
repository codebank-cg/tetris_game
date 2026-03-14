package model

import (
	"reflect"
	"testing"
)

// Test that each generated bag contains all seven tetromino types exactly once.
func TestRandomizer_BagContainsAllSeven(t *testing.T) {
	r := NewRandomizer()
	seen := map[TetrominoType]bool{}
	for i := 0; i < 7; i++ {
		p := r.NextPiece()
		seen[p] = true
	}
	if len(seen) != 7 {
		t.Fatalf("expected 7 unique pieces in a bag, got %d: %v", len(seen), seen)
	}
}

// Ensure that two random sequences with different seeds are not identical.
func TestRandomizer_RandomSequencesDiffer(t *testing.T) {
	r1 := NewRandomizer()
	r2 := NewRandomizer()
	// draw some pieces from both
	s1 := make([]TetrominoType, 14)
	s2 := make([]TetrominoType, 14)
	for i := 0; i < 14; i++ {
		s1[i] = r1.NextPiece()
		s2[i] = r2.NextPiece()
	}
	if reflect.DeepEqual(s1, s2) {
		t.Fatalf("expected different sequences for different randomizers, got identical: %v", s1)
	}
}

// Ensure reproducibility with a fixed seed.
func TestRandomizer_SeedReproducible(t *testing.T) {
	seed := int64(12345)
	rA := NewRandomizer()
	rA.SetSeed(seed)
	rB := NewRandomizer()
	rB.SetSeed(seed)

	n := 12
	a := make([]TetrominoType, n)
	b := make([]TetrominoType, n)
	for i := 0; i < n; i++ {
		a[i] = rA.NextPiece()
		b[i] = rB.NextPiece()
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("expected sequences to be equal for same seed: %v vs %v", a, b)
	}
}

// Test GetNextPieces (peek functionality).
func TestRandomizer_GetNextPiecesPeek(t *testing.T) {
	r := NewRandomizer()
	first := r.NextPiece()
	if len(r.GetNextPieces(3)) != 3 {
		t.Fatalf("expected 3 upcoming pieces in bag after first draw")
	}
	peek1 := r.GetNextPieces(3)
	// consume one piece and ensure peek reflects the new upcoming pieces
	_ = first
	// Now perform a piece draw and compare with a fresh peek
	r.NextPiece()
	peek2 := r.GetNextPieces(3)
	if reflect.DeepEqual(peek1, peek2) {
		t.Fatalf("GetNextPieces should reflect state change after consuming a piece")
	}
}
