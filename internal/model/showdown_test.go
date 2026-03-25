package model

import (
	"testing"
	"time"
)

// ── AutoPlayer per-instance weights ─────────────────────────────────────────

func TestAutoPlayer_PerInstanceWeights_Independent(t *testing.T) {
	apA := NewAutoPlayer()
	apB := NewAutoPlayer()

	apA.SetWeights(map[string]float64{"aggregateHeight": -0.05, "holes": -0.10, "bumpiness": -0.03, "wells": -0.02})
	apB.SetWeights(map[string]float64{"aggregateHeight": -0.90, "holes": -1.50, "bumpiness": -0.60, "wells": -0.45})

	wa := apA.GetWeights()
	wb := apB.GetWeights()

	if wa["aggregateHeight"] == wb["aggregateHeight"] {
		t.Error("Expected different aggregateHeight weights between bots")
	}
	if wa["aggregateHeight"] != -0.05 {
		t.Errorf("BotA aggregateHeight = %.2f, want -0.05", wa["aggregateHeight"])
	}
	if wb["aggregateHeight"] != -0.90 {
		t.Errorf("BotB aggregateHeight = %.2f, want -0.90", wb["aggregateHeight"])
	}
}

// ── TetrisCount ──────────────────────────────────────────────────────────────

func TestTetrisCount_IncrementOn4Lines(t *testing.T) {
	gs := NewGameState()

	// Fill bottom 4 rows completely
	for y := 0; y < 4; y++ {
		for x := 0; x < 10; x++ {
			gs.Board.Set(x, y, 1)
		}
	}
	gs.ClearedLines = []int{3, 2, 1, 0}
	gs.ClearAnimFrame = 1
	gs.ClearAnimIndex = 0

	// Drive animation to completion
	for gs.IsClearAnimating() {
		gs.UpdateClearAnimation()
	}

	if gs.TetrisCount != 1 {
		t.Errorf("TetrisCount = %d, want 1 after 4-line clear", gs.TetrisCount)
	}
}

func TestTetrisCount_NoIncrementOnFewerLines(t *testing.T) {
	for lines := 1; lines <= 3; lines++ {
		gs := NewGameState()

		cleared := make([]int, lines)
		for i := 0; i < lines; i++ {
			cleared[i] = i
			for x := 0; x < 10; x++ {
				gs.Board.Set(x, i, 1)
			}
		}
		gs.ClearedLines = cleared
		gs.ClearAnimFrame = 1
		gs.ClearAnimIndex = 0

		for gs.IsClearAnimating() {
			gs.UpdateClearAnimation()
		}

		if gs.TetrisCount != 0 {
			t.Errorf("%d-line clear: TetrisCount = %d, want 0", lines, gs.TetrisCount)
		}
	}
}

// ── GameState Reset ──────────────────────────────────────────────────────────

func TestGameState_Reset_ClearsAnimState(t *testing.T) {
	gs := NewGameState()

	// Set non-zero state
	gs.ClearedLines = []int{5, 4}
	gs.ClearAnimFrame = 7
	gs.ClearAnimIndex = 1
	gs.autoMoveStep = 2
	gs.TetrisCount = 3

	gs.Reset()

	if len(gs.ClearedLines) != 0 {
		t.Errorf("ClearedLines not cleared after Reset: %v", gs.ClearedLines)
	}
	if gs.ClearAnimFrame != 0 {
		t.Errorf("ClearAnimFrame = %d, want 0", gs.ClearAnimFrame)
	}
	if gs.ClearAnimIndex != 0 {
		t.Errorf("ClearAnimIndex = %d, want 0", gs.ClearAnimIndex)
	}
	if gs.autoMoveStep != 0 {
		t.Errorf("autoMoveStep = %d, want 0", gs.autoMoveStep)
	}
	if gs.TetrisCount != 0 {
		t.Errorf("TetrisCount = %d, want 0", gs.TetrisCount)
	}
	if !gs.CanHold {
		t.Error("CanHold should be true after Reset")
	}
}

// ── DefaultWeights ───────────────────────────────────────────────────────────

func TestDefaultWeights_HasAllFourKeys(t *testing.T) {
	w := DefaultWeights()
	for _, key := range []string{"aggregateHeight", "holes", "bumpiness", "wells"} {
		if _, ok := w[key]; !ok {
			t.Errorf("DefaultWeights() missing key: %s", key)
		}
	}
}

// ── Presets ──────────────────────────────────────────────────────────────────

func TestGetPreset_AllFivePresets(t *testing.T) {
	tests := []struct {
		name            string
		wantAggHeight   float64
		wantHoles       float64
		wantBumpiness   float64
		wantWells       float64
	}{
		{"aggressive", -0.05, -0.10, -0.03, -0.02},
		{"conservative", -0.90, -1.50, -0.60, -0.45},
		{"balanced", -0.54, -0.90, -0.36, -0.24},
		{"speedrun", -0.24, -0.45, -0.15, -0.09},
		{"chaos", -0.01, -0.02, -0.01, -0.01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := GetPreset(tt.name)
			if err != nil {
				t.Fatalf("GetPreset(%q) error: %v", tt.name, err)
			}
			if w["aggregateHeight"] != tt.wantAggHeight {
				t.Errorf("aggregateHeight = %.2f, want %.2f", w["aggregateHeight"], tt.wantAggHeight)
			}
			if w["holes"] != tt.wantHoles {
				t.Errorf("holes = %.2f, want %.2f", w["holes"], tt.wantHoles)
			}
			if w["bumpiness"] != tt.wantBumpiness {
				t.Errorf("bumpiness = %.2f, want %.2f", w["bumpiness"], tt.wantBumpiness)
			}
			if w["wells"] != tt.wantWells {
				t.Errorf("wells = %.2f, want %.2f", w["wells"], tt.wantWells)
			}
		})
	}
}

func TestGetPreset_UnknownNameReturnsError(t *testing.T) {
	_, err := GetPreset("unknown-preset")
	if err == nil {
		t.Error("Expected error for unknown preset name, got nil")
	}
}

func TestGetPreset_Conservative_ScaledValues(t *testing.T) {
	cons, _ := GetPreset("conservative")
	agg, _ := GetPreset("aggressive")

	// Conservative should have significantly higher magnitude than aggressive
	// (conservative -0.90 vs aggressive -0.05 = 18x ratio)
	if cons["aggregateHeight"] >= agg["aggregateHeight"] {
		t.Errorf("conservative aggregateHeight (%.2f) should be more negative than aggressive (%.2f)",
			cons["aggregateHeight"], agg["aggregateHeight"])
	}
	magRatio := cons["aggregateHeight"] / agg["aggregateHeight"]
	if magRatio < 10.0 {
		t.Errorf("conservative/aggressive aggregateHeight ratio = %.2f, want >10x", magRatio)
	}
}

// ── NewShowdownState ─────────────────────────────────────────────────────────

func TestNewShowdownState_BotWeights_AreIndependentAndDifferent(t *testing.T) {
	ss := NewShowdownState("aggressive", "conservative")

	wa := ss.BotA.Player.GetWeights()
	wb := ss.BotB.Player.GetWeights()

	if wa["aggregateHeight"] == wb["aggregateHeight"] {
		t.Error("BotA and BotB should have different weights")
	}
	if wa["aggregateHeight"] != -0.05 {
		t.Errorf("BotA aggregateHeight = %.2f, want -0.05", wa["aggregateHeight"])
	}
	if wb["aggregateHeight"] != -0.90 {
		t.Errorf("BotB aggregateHeight = %.2f, want -0.90", wb["aggregateHeight"])
	}
}

func TestNewShowdownState_InitialState(t *testing.T) {
	ss := NewShowdownState("aggressive", "conservative")

	if !ss.Running {
		t.Error("ShowdownState should be Running=true initially")
	}
	if ss.SpeedLevel != 1 {
		t.Errorf("SpeedLevel = %d, want 1", ss.SpeedLevel)
	}
	if ss.WinnerBot != "" {
		t.Errorf("WinnerBot = %q, want empty", ss.WinnerBot)
	}
	if ss.BotA.State == nil || ss.BotB.State == nil {
		t.Error("Both bot game states should be initialized")
	}
}

// ── Tick ─────────────────────────────────────────────────────────────────────

func TestShowdownState_Tick_SkipsAIWhenAnimating(t *testing.T) {
	ss := NewShowdownState("aggressive", "conservative")

	// Force BotA into animation state
	ss.BotA.State.ClearedLines = []int{5}
	ss.BotA.State.ClearAnimFrame = 1
	ss.BotA.State.ClearAnimIndex = 0
	initialScore := ss.BotA.State.Score
	scoresBefore := ss.BotA.State.LinesCleared

	// Tick with enough time elapsed
	ss.Tick(time.Now().Add(2 * time.Second))

	// BotA should have progressed animation, not placed a piece
	// Score should not have changed (no new piece placed)
	if ss.BotA.State.Score != initialScore {
		t.Error("BotA score changed during animation (expected no new piece placement)")
	}
	_ = scoresBefore
}

func TestShowdownState_Tick_NotRunning_DoesNothing(t *testing.T) {
	ss := NewShowdownState("aggressive", "conservative")
	ss.Running = false
	beforeA := ss.BotA.State.PieceCount
	beforeB := ss.BotB.State.PieceCount

	ss.Tick(time.Now().Add(10 * time.Second))

	if ss.BotA.State.PieceCount != beforeA || ss.BotB.State.PieceCount != beforeB {
		t.Error("Tick should not advance bots when not Running")
	}
}

// ── Game-over ────────────────────────────────────────────────────────────────

func TestShowdown_GameOver_BotADies(t *testing.T) {
	ss := NewShowdownState("aggressive", "conservative")
	ss.BotA.State.GameOver = true
	ss.BotB.State.Score = 1000

	ss.Tick(time.Now())

	if ss.Running {
		t.Error("Running should be false when BotA dies")
	}
	if ss.WinnerBot != "B" {
		t.Errorf("WinnerBot = %q, want \"B\"", ss.WinnerBot)
	}
}

func TestShowdown_GameOver_BotBDies(t *testing.T) {
	ss := NewShowdownState("aggressive", "conservative")
	ss.BotB.State.GameOver = true
	ss.BotA.State.Score = 2000

	ss.Tick(time.Now())

	if ss.Running {
		t.Error("Running should be false when BotB dies")
	}
	if ss.WinnerBot != "A" {
		t.Errorf("WinnerBot = %q, want \"A\"", ss.WinnerBot)
	}
}

func TestShowdown_GameOver_Draw(t *testing.T) {
	ss := NewShowdownState("aggressive", "conservative")
	ss.BotA.State.GameOver = true
	ss.BotB.State.GameOver = true

	ss.Tick(time.Now())

	if ss.Running {
		t.Error("Running should be false on draw")
	}
	if ss.WinnerBot != "DRAW" {
		t.Errorf("WinnerBot = %q, want \"DRAW\"", ss.WinnerBot)
	}
}

// ── WinnerBoard ──────────────────────────────────────────────────────────────

func TestWinnerBoard_CapturesValuesNotPointer(t *testing.T) {
	ss := NewShowdownState("aggressive", "conservative")
	ss.BotA.State.Score = 5000
	ss.BotB.State.GameOver = true

	ss.Tick(time.Now())

	capturedScore := ss.Winner.BestScore
	capturedBot := ss.Winner.BestBot

	// Reset game states (simulate R key)
	ss.Reset()

	// WinnerBoard should be unchanged
	if ss.Winner.BestScore != capturedScore {
		t.Errorf("WinnerBoard.BestScore changed after Reset: got %d, want %d", ss.Winner.BestScore, capturedScore)
	}
	if ss.Winner.BestBot != capturedBot {
		t.Errorf("WinnerBoard.BestBot changed after Reset: got %q, want %q", ss.Winner.BestBot, capturedBot)
	}
}

func TestWinnerBoard_PreservesSessionBest(t *testing.T) {
	ss := NewShowdownState("aggressive", "conservative")

	// Run 1: BotA wins with score 10000
	ss.BotA.State.Score = 10000
	ss.BotB.State.GameOver = true
	ss.Tick(time.Now())
	if ss.Winner.BestScore != 10000 {
		t.Fatalf("After run 1: BestScore = %d, want 10000", ss.Winner.BestScore)
	}

	// Reset for run 2
	ss.Reset()

	// Run 2: BotB wins with score 5000 (lower)
	ss.BotB.State.Score = 5000
	ss.BotA.State.GameOver = true
	ss.Tick(time.Now())

	if ss.Winner.BestScore != 10000 {
		t.Errorf("BestScore = %d, want 10000 (run 1 was better)", ss.Winner.BestScore)
	}
	if ss.Winner.BestBot != "A" {
		t.Errorf("BestBot = %q, want \"A\" (run 1 winner)", ss.Winner.BestBot)
	}
	if ss.Winner.Runs != 2 {
		t.Errorf("Runs = %d, want 2", ss.Winner.Runs)
	}
}

// ── Reset preserves weights and WinnerBoard ───────────────────────────────────

func TestShowdown_Reset_AutoPlayerWeightsSurviveReset(t *testing.T) {
	ss := NewShowdownState("aggressive", "conservative")
	waBefore := ss.BotA.Player.GetWeights()
	wbBefore := ss.BotB.Player.GetWeights()

	ss.BotA.State.GameOver = true
	ss.Tick(time.Now())
	ss.Reset()

	waAfter := ss.BotA.Player.GetWeights()
	wbAfter := ss.BotB.Player.GetWeights()

	if waAfter["aggregateHeight"] != waBefore["aggregateHeight"] {
		t.Errorf("BotA weights changed after Reset")
	}
	if wbAfter["aggregateHeight"] != wbBefore["aggregateHeight"] {
		t.Errorf("BotB weights changed after Reset")
	}
}

func TestShowdown_Reset_WinnerBoardSurvives(t *testing.T) {
	ss := NewShowdownState("aggressive", "conservative")
	ss.BotA.State.Score = 7500
	ss.BotB.State.GameOver = true
	ss.Tick(time.Now())

	savedScore := ss.Winner.BestScore
	savedBot := ss.Winner.BestBot
	savedRuns := ss.Winner.Runs

	ss.Reset()

	if ss.Winner.BestScore != savedScore {
		t.Errorf("WinnerBoard.BestScore changed after Reset")
	}
	if ss.Winner.BestBot != savedBot {
		t.Errorf("WinnerBoard.BestBot changed after Reset")
	}
	if ss.Winner.Runs != savedRuns {
		t.Errorf("WinnerBoard.Runs changed after Reset")
	}
}
