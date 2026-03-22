package model

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
