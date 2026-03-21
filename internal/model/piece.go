package model

import "fmt"

// Tetromino represents a Tetris piece.
type Tetromino struct {
	Type     TetrominoType
	Color    int
	X        int
	Y        int
	Rotation int     // 0-3
	Matrix   [][]int // current rotation state (4x4)
}

// pieceDefinitions holds the base shapes for all tetrominoes.
var pieceDefinitions = map[TetrominoType][][][]int{
	TetrominoI: {
		// 0
		{
			{0, 0, 0, 0},
			{1, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		// R
		{
			{0, 0, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 1, 0},
		},
		// 2
		{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{1, 1, 1, 1},
			{0, 0, 0, 0},
		},
		// L
		{
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
		},
	},
	TetrominoO: {
		// O piece has 1 state (symmetric)
		{
			{0, 1, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
	},
	TetrominoT: {
		// State 0: T pointing up
		{
			{0, 1, 0, 0},
			{1, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		// State 1: T pointing right
		{
			{0, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
		},
		// State 2: T pointing down
		{
			{0, 0, 0, 0},
			{1, 1, 1, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
		},
		// State 3: T pointing left
		{
			{0, 1, 0, 0},
			{1, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
		},
	},
	TetrominoS: {
		{
			{0, 1, 1, 0},
			{1, 1, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0},
			{0, 1, 1, 0},
			{1, 1, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{1, 0, 0, 0},
			{1, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
		},
	},
	TetrominoZ: {
		{
			{1, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 0, 1, 0},
			{0, 1, 1, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0},
			{1, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 0, 0},
			{1, 1, 0, 0},
			{1, 0, 0, 0},
			{0, 0, 0, 0},
		},
	},
	TetrominoJ: {
		{
			{1, 0, 0, 0},
			{1, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 1, 0},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0},
			{1, 1, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{1, 1, 0, 0},
			{0, 0, 0, 0},
		},
	},
	TetrominoL: {
		{
			{0, 0, 1, 0},
			{1, 1, 1, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0},
			{1, 1, 1, 0},
			{1, 0, 0, 0},
			{0, 0, 0, 0},
		},
		{
			{1, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
		},
	},
}

// pieceColors maps tetromino types to colors (1-7).
var pieceColors = map[TetrominoType]int{
	TetrominoI: 1, // Cyan
	TetrominoO: 2, // Yellow
	TetrominoT: 3, // Magenta
	TetrominoS: 4, // Green
	TetrominoZ: 5, // Red
	TetrominoJ: 6, // Blue
	TetrominoL: 7, // Orange
}

// NewTetromino creates a new tetromino of the specified type.
func NewTetromino(t TetrominoType) *Tetromino {
	def := pieceDefinitions[t]
	matrix := make([][]int, 4)
	for i := range def[0] {
		matrix[i] = make([]int, 4)
		copy(matrix[i], def[0][i])
	}
	return &Tetromino{
		Type:     t,
		Color:    pieceColors[t],
		X:        3, // spawn in center
		Y:        18,
		Rotation: 0,
		Matrix:   matrix,
	}
}

// GetMatrix returns the current rotation matrix.
func (t *Tetromino) GetMatrix() [][]int {
	return t.Matrix
}

// RotateClockwise rotates the piece 90 degrees clockwise.
func (t *Tetromino) RotateClockwise() {
	t.Rotation = (t.Rotation + 1) % 4
	t.updateMatrix()
}

// RotateCounterClockwise rotates the piece 90 degrees counter-clockwise.
func (t *Tetromino) RotateCounterClockwise() {
	t.Rotation = (t.Rotation - 1 + 4) % 4
	t.updateMatrix()
}

// Move moves the piece by the specified offset.
func (t *Tetromino) Move(dx, dy int) {
	t.X += dx
	t.Y += dy
}

// GetPosition returns the current position.
func (t *Tetromino) GetPosition() (int, int) {
	return t.X, t.Y
}

// srsWallKicks holds the SRS wall kick offset tests for each rotation transition.
// Key format: "fromRotation>toRotation", values are (dx, dy) offset pairs to try.
// Based on the official Tetris SRS specification.
var srsWallKicks = map[string][][2]int{
	// JLSTZ pieces
	"0>1": {{0, 0}, {-1, 0}, {-1, 1}, {0, -2}, {-1, -2}},
	"1>0": {{0, 0}, {1, 0}, {1, -1}, {0, 2}, {1, 2}},
	"1>2": {{0, 0}, {1, 0}, {1, -1}, {0, 2}, {1, 2}},
	"2>1": {{0, 0}, {-1, 0}, {-1, 1}, {0, -2}, {-1, -2}},
	"2>3": {{0, 0}, {1, 0}, {1, 1}, {0, -2}, {1, -2}},
	"3>2": {{0, 0}, {-1, 0}, {-1, -1}, {0, 2}, {-1, 2}},
	"3>0": {{0, 0}, {-1, 0}, {-1, -1}, {0, 2}, {-1, 2}},
	"0>3": {{0, 0}, {1, 0}, {1, 1}, {0, -2}, {1, -2}},
}

// srsWallKicksI holds SRS wall kick offsets for the I piece.
var srsWallKicksI = map[string][][2]int{
	"0>1": {{0, 0}, {-2, 0}, {1, 0}, {-2, -1}, {1, 2}},
	"1>0": {{0, 0}, {2, 0}, {-1, 0}, {2, 1}, {-1, -2}},
	"1>2": {{0, 0}, {-1, 0}, {2, 0}, {-1, 2}, {2, -1}},
	"2>1": {{0, 0}, {1, 0}, {-2, 0}, {1, -2}, {-2, 1}},
	"2>3": {{0, 0}, {2, 0}, {-1, 0}, {2, 1}, {-1, -2}},
	"3>2": {{0, 0}, {-2, 0}, {1, 0}, {-2, -1}, {1, 2}},
	"3>0": {{0, 0}, {1, 0}, {-2, 0}, {1, -2}, {-2, 1}},
	"0>3": {{0, 0}, {-1, 0}, {2, 0}, {-1, 2}, {2, -1}},
}

// GetWallKicks returns the SRS kick offsets for a rotation transition.
func (t *Tetromino) GetWallKicks(fromRot, toRot int) [][2]int {
	key := fmt.Sprintf("%d>%d", fromRot, toRot)
	if t.Type == TetrominoI {
		if kicks, ok := srsWallKicksI[key]; ok {
			return kicks
		}
	} else if t.Type == TetrominoO {
		// O piece never needs wall kicks
		return [][2]int{{0, 0}}
	} else {
		if kicks, ok := srsWallKicks[key]; ok {
			return kicks
		}
	}
	return [][2]int{{0, 0}}
}

// updateMatrix updates the matrix based on current rotation.
func (t *Tetromino) updateMatrix() {
	def := pieceDefinitions[t.Type]
	matrix := make([][]int, 4)
	for i := range def[t.Rotation] {
		matrix[i] = make([]int, 4)
		copy(matrix[i], def[t.Rotation][i])
	}
	t.Matrix = matrix
}
