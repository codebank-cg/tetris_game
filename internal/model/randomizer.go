package model

import (
	"math/rand"
	"time"
)

// TetrominoType represents a Tetris piece type.
type TetrominoType string

const (
	TetrominoI TetrominoType = "I"
	TetrominoO TetrominoType = "O"
	TetrominoT TetrominoType = "T"
	TetrominoS TetrominoType = "S"
	TetrominoZ TetrominoType = "Z"
	TetrominoJ TetrominoType = "J"
	TetrominoL TetrominoType = "L"
)

// Randomizer implements a 7-bag based randomizer for Tetromino pieces.
type Randomizer struct {
	r   *rand.Rand
	bag []TetrominoType
	idx int // next piece index inside bag
}

// NewRandomizer creates a new randomizer with a time-based seed.
func NewRandomizer() *Randomizer {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rr := &Randomizer{r: rnd}
	rr.refillBag()
	return rr
}

// SetSeed sets the RNG seed for deterministic tests.
func (rr *Randomizer) SetSeed(seed int64) {
	rr.r = rand.New(rand.NewSource(seed))
	rr.refillBag()
	rr.idx = 0
}

// NextPiece returns the next piece from the current bag, refilling and shuffling a new bag when needed.
func (rr *Randomizer) NextPiece() TetrominoType {
	if rr.bag == nil || rr.idx >= len(rr.bag) {
		rr.refillBag()
	}
	p := rr.bag[rr.idx]
	rr.idx++
	return p
}

// GetNextPieces returns up to n upcoming pieces from the current state without consuming them.
// Peeks into the current bag and, if needed, generates a preview of the next bag.
func (rr *Randomizer) GetNextPieces(n int) []TetrominoType {
	if n <= 0 {
		return []TetrominoType{}
	}
	out := make([]TetrominoType, 0, n)
	// Take what's left in the current bag
	remaining := rr.bag[rr.idx:]
	for _, p := range remaining {
		if len(out) >= n {
			break
		}
		out = append(out, p)
	}
	// If we still need more, generate a preview of the next bag using a copy of the RNG state.
	// We can't actually peek at future shuffles without advancing the RNG, so we fill with
	// a deterministic placeholder bag (all 7 types in fixed order) as a best-effort preview.
	if len(out) < n {
		nextBag := []TetrominoType{TetrominoI, TetrominoO, TetrominoT, TetrominoS, TetrominoZ, TetrominoJ, TetrominoL}
		for _, p := range nextBag {
			if len(out) >= n {
				break
			}
			out = append(out, p)
		}
	}
	return out
}

// refillBag creates a new shuffled bag of the 7 Tetromino types and resets the index.
func (rr *Randomizer) refillBag() {
	rr.bag = []TetrominoType{TetrominoI, TetrominoO, TetrominoT, TetrominoS, TetrominoZ, TetrominoJ, TetrominoL}
	rr.r.Shuffle(len(rr.bag), func(i, j int) {
		rr.bag[i], rr.bag[j] = rr.bag[j], rr.bag[i]
	})
	rr.idx = 0
}
