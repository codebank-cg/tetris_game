package model

import (
	tu "github.com/oc-garden/tetris_game/internal/testutil"
	"testing"
)

func TestBoardInfrastructure_Pattern(t *testing.T) {
	_ = tu.NewTestBoard()
	p := tu.NewTestPiece("I")
	if p.Type != "I" || p.X != 0 || p.Y != 0 {
		t.Fatalf("unexpected test piece: %+v", p)
	}
}
