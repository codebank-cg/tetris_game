# TODOS

## AI Showdown

### Verbose debug mode (-v flag)

**What:** Add a `-v` / `--verbose` flag that logs each bot's top-3 move candidates and their scores to stderr during showdown mode.

**Why:** When both bots play identically, there's no way to know if the weights threading is working correctly or if the AI is seeing different positions. Debug visibility is near-zero without this.

**Context:** After the `heuristicWeights` refactor to per-instance, a missing weights param silently produces zero-value weights — bots play identically with no error. This flag makes AI behavior transparent for tuning new presets and catching silent correctness bugs. Log format: `BotA tick: top move (x=4, rot=1, score=12.3), alt (x=3, rot=0, score=11.1)`.

**Effort:** S
**Priority:** P3
**Depends on:** AI Showdown feature shipped

---

### ELO tracking across sessions

**What:** Write win/loss results to `~/.tetris-showdown-elo.json` so bot rankings accumulate over time.

**Why:** The in-session winner board resets on exit. ELO tracking lets you run hundreds of rounds across many sessions and see which personality actually dominates over time.

**Context:** Winner board (best winner, score, duration) is already implemented per-session. ELO adds persistence. Format: `{"aggressive": {"wins": 12, "losses": 8, "elo": 1523}, "conservative": {...}}`. Update on each game-over before displaying winner screen. Should be atomic write to prevent corruption if user Ctrl+C.

**Effort:** S
**Priority:** P2
**Depends on:** AI Showdown feature shipped

---

### Interactive weight tuner (TUI "synth knobs")

**What:** Add an interactive panel (press E in showdown mode) to nudge each bot's heuristic weights up/down with arrow keys, seeing strategy change live.

**Why:** Adding new presets currently requires editing Go source. A live tuner makes the AI legible and experimentable — like a synthesizer with knobs. Watch the bot change behavior in real time as you increase the hole penalty.

**Context:** After the per-instance weights refactor, `AutoPlayer.SetWeights()` can be called at any time and the next `FindBestMoveWithNext` call will use the new values. The tuner panel overlays the stats panel; left/right selects the weight to tune, up/down nudges by 0.05. Display the current value as a decimal. Press E again to dismiss. Save tuned preset with a name if desired.

**Effort:** M
**Priority:** P2
**Depends on:** AI Showdown feature shipped; per-instance weights refactor (Step 0)

## Completed

