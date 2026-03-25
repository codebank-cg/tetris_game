package model

import "time"

// WinnerBoard tracks the session's best result across restarts.
// Values are copied at game-over so Reset() does not affect them.
type WinnerBoard struct {
	BestBot      string        // "A" or "B"
	BestScore    int           // winner's score at game-over
	BestDuration time.Duration // RunStart to first game-over
	Runs         int           // total completed runs
}

// BotState holds per-bot state for a showdown.
type BotState struct {
	State           *GameState
	Player          *AutoPlayer
	Timer           time.Time
	CurrentDecision *MoveDecision
}

// ShowdownState orchestrates two bots playing simultaneously.
type ShowdownState struct {
	BotA       *BotState
	BotB       *BotState
	Running    bool
	WinnerBot  string // "A", "B", "DRAW", or ""
	SpeedLevel int    // 1-5; sole speed authority in showdown mode
	Winner     WinnerBoard
	RunStart   time.Time
}

// NewShowdownState creates a fresh showdown with the given preset names.
// If a preset name is unrecognized, DefaultWeights() is used as fallback.
func NewShowdownState(presetA, presetB string) *ShowdownState {
	weightsA, err := GetPreset(presetA)
	if err != nil {
		weightsA = DefaultWeights()
	}
	weightsB, err := GetPreset(presetB)
	if err != nil {
		weightsB = DefaultWeights()
	}

	playerA := NewAutoPlayer()
	playerA.SetWeights(weightsA)

	playerB := NewAutoPlayer()
	playerB.SetWeights(weightsB)

	now := time.Now()
	return &ShowdownState{
		BotA: &BotState{
			State:  NewGameState(),
			Player: playerA,
			Timer:  now,
		},
		BotB: &BotState{
			State:  NewGameState(),
			Player: playerB,
			Timer:  now,
		},
		Running:    true,
		SpeedLevel: 1,
		RunStart:   now,
	}
}

// Tick advances showdown state by one game-loop step.
// Call at ~60fps. now is the current time.
func (ss *ShowdownState) Tick(now time.Time) {
	if !ss.Running {
		return
	}

	ss.tickBot(ss.BotA, now)
	ss.tickBot(ss.BotB, now)

	ss.checkGameOver(now)
}

// tickBot advances one bot for this tick.
func (ss *ShowdownState) tickBot(bot *BotState, now time.Time) {
	if bot.State.GameOver {
		return
	}

	if bot.State.IsClearAnimating() {
		bot.State.UpdateClearAnimation()
		return
	}

	delay := GetDelayForSpeed(bot.State.GetDropInterval(), ss.SpeedLevel)
	if now.Sub(bot.Timer) < time.Duration(delay)*time.Millisecond {
		return
	}

	if bot.CurrentDecision == nil {
		bot.CurrentDecision = FindBestMoveWithNext(bot.State, bot.Player.GetWeights())
	}

	if bot.CurrentDecision != nil {
		done := ExecuteMove(bot.State, bot.CurrentDecision)
		// At max speed drain all drop steps in one tick
		if !done && ss.SpeedLevel == 5 && IsInDropPhase(bot.State) {
			for !done {
				done = ExecuteMove(bot.State, bot.CurrentDecision)
			}
		}
		if done {
			bot.CurrentDecision = nil
		}
	}

	bot.Timer = now
}

// checkGameOver detects when a bot's game ends and updates state accordingly.
func (ss *ShowdownState) checkGameOver(now time.Time) {
	aOver := ss.BotA.State.GameOver
	bOver := ss.BotB.State.GameOver

	if !aOver && !bOver {
		return
	}

	ss.Running = false

	if aOver && bOver {
		ss.WinnerBot = "DRAW"
	} else if aOver {
		ss.WinnerBot = "B"
	} else {
		ss.WinnerBot = "A"
	}

	ss.Winner.Runs++
	duration := now.Sub(ss.RunStart)

	// Determine winner's score for WinnerBoard comparison
	var winnerScore int
	var winnerBot string
	if ss.WinnerBot == "A" {
		winnerScore = ss.BotA.State.Score
		winnerBot = "A"
	} else if ss.WinnerBot == "B" {
		winnerScore = ss.BotB.State.Score
		winnerBot = "B"
	} else {
		// DRAW: pick the higher score
		if ss.BotA.State.Score >= ss.BotB.State.Score {
			winnerScore = ss.BotA.State.Score
			winnerBot = "A"
		} else {
			winnerScore = ss.BotB.State.Score
			winnerBot = "B"
		}
	}

	if ss.Winner.Runs == 1 || winnerScore > ss.Winner.BestScore {
		ss.Winner.BestBot = winnerBot
		ss.Winner.BestScore = winnerScore
		ss.Winner.BestDuration = duration
	}
}

// Reset restarts both bots' game states without recreating the ShowdownState.
// AutoPlayer weights and WinnerBoard are preserved.
func (ss *ShowdownState) Reset() {
	ss.BotA.State.Reset()
	ss.BotB.State.Reset()
	ss.BotA.CurrentDecision = nil
	ss.BotB.CurrentDecision = nil
	now := time.Now()
	ss.BotA.Timer = now
	ss.BotB.Timer = now
	ss.Running = true
	ss.WinnerBot = ""
	ss.RunStart = now
}
