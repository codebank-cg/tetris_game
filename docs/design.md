# Tetris Game Design Document

## 1. Game Overview

### Objective
Tetris is a puzzle game where the player manipulates falling geometric shapes called tetrominoes. The objective is to create complete horizontal lines across the playing field without gaps. When a line is completed, it disappears, and any blocks above it fall down. The game ends when the blocks stack up to the top of the playing field and new pieces can no longer enter.

### Basic Gameplay
- Tetrominoes fall from the top of the playing field
- Player can move pieces left/right, rotate them, and accelerate their descent
- When a tetromino lands, it becomes part of the static block structure
- Complete horizontal lines are cleared, causing blocks above to fall
- Game speed increases as the player clears more lines
- Player earns points for clearing lines and performing advanced maneuvers

## 2. Tetromino Definitions

### The Seven Tetrominoes
Each tetromino is composed of four square blocks. The seven tetrominoes are named after the letters they resemble:

#### I-Tetromino (Cyan)
```
Shape: ####
Color: Cyan (#00FFFF)
Description: Straight line of 4 blocks
Rotation States: 2 (horizontal and vertical)
```

#### O-Tetromino (Yellow)
```
Shape: ##
       ##
Color: Yellow (#FFFF00)
Description: 2x2 square
Rotation States: 1 (symmetric, no visible rotation)
```

#### T-Tetromino (Magenta)
```
Shape: ###
        #
Color: Magenta (#FF00FF)
Description: T-shaped piece
Rotation States: 4
```

#### S-Tetromino (Green)
```
Shape:  ##
       ##
Color: Green (#00FF00)
Description: S-shaped piece
Rotation States: 4
```

#### Z-Tetromino (Red)
```
Shape: ##
        ##
Color: Red (#FF0000)
Description: Z-shaped piece (mirror of S)
Rotation States: 4
```

#### J-Tetromino (Blue)
```
Shape: ###
       #
Color: Blue (#0000FF)
Description: J-shaped piece
Rotation States: 4
```

#### L-Tetromino (Orange)
```
Shape: ###
         #
Color: Orange (#FFA500)
Description: L-shaped piece (mirror of J)
Rotation States: 4
```

### Rotation System
For this implementation, we use a **Basic Rotation System** (simplified from the Super Rotation System):
- All pieces rotate clockwise around their center
- Counter-clockwise rotation is also supported
- No wall kicks in MVP (rotation fails if blocked)
- O-piece has no visible rotation due to symmetry

## 3. Board Specifications

### Playing Field Dimensions
- **Width**: 10 columns
- **Height**: 20 rows (visible playing area)
- **Total Cells**: 200 playable cells

### Coordinate System
```
Origin: Bottom-left corner (0, 0)
X-axis: Increases from left (0) to right (9)
Y-axis: Increases from bottom (0) to top (19)

Coordinate Format: (x, y)
Example: (5, 10) = column 5, row 10
```

### Spawn Position
New tetrominoes spawn at the top-center of the board:
- **Default Spawn Y**: Row 18-19 (varies by piece)
- **Default Spawn X**: Columns 3-6 (centered)
- **Spawn Area**: 2 rows above the visible playing field

### Cell States
Each cell can be in one of two states:
- **Empty**: No block present (rendered as dark/empty)
- **Filled**: Block present (rendered with piece color)

## 4. Game Mechanics

### Movement

#### Left/Right Movement
- Player can move piece one column per input
- Movement is blocked by walls or existing blocks
- Input: Left Arrow / Right Arrow

#### Soft Drop
- Player can accelerate piece descent
- Holds piece at bottom while input is held
- Earns 1 point per cell dropped
- Input: Down Arrow

#### Hard Drop
- Instantly drops piece to bottom
- Piece locks immediately upon landing
- Earns 2 points per cell dropped
- Input: Space Bar

### Rotation

#### Clockwise Rotation
- Rotates piece 90 degrees clockwise
- Fails if rotation would cause collision
- Input: X key or Up Arrow

#### Counter-Clockwise Rotation
- Rotates piece 90 degrees counter-clockwise
- Fails if rotation would cause collision
- Input: Z key

### Collision Detection
Collision occurs when:
- Piece moves outside board boundaries
- Piece overlaps with locked blocks
- Rotation would place blocks in occupied cells

### Line Clearing
When all 10 cells in a row are filled:
1. The complete line is removed
2. All lines above shift down by one row
3. Player earns points based on number of lines cleared simultaneously

## 5. Scoring System

### Line Clear Scoring
Points are awarded based on the number of lines cleared at once, multiplied by the current level:

| Lines Cleared | Base Points | Formula |
|---------------|-------------|---------|
| 1 (Single)    | 100         | 100 × level |
| 2 (Double)    | 300         | 300 × level |
| 3 (Triple)    | 500         | 500 × level |
| 4 (Tetris)    | 800         | 800 × level |

### Drop Scoring
- **Soft Drop**: 1 point per cell the piece is moved down while holding Down Arrow
- **Hard Drop**: 2 points per cell when piece instantly drops to bottom

### Combo System
Clearing lines in consecutive drops earns combo bonuses:
- **Combo Bonus**: 50 × combo count × level
- Combo count resets when a drop doesn't clear any lines

### T-Spin Scoring (Optional Enhancement)
When T-tetromino is rotated into a tight space:
- **T-Spin Single**: 800 × level
- **T-Spin Double**: 1200 × level
- **T-Spin Triple**: 1600 × level

### Back-to-Back Bonus
Clearing multiple lines (Tetris) or T-spins consecutively:
- **Bonus Multiplier**: 1.5× base score
- Resets when clearing 1-2 lines without T-spin

## 6. Game States

### Playing State
- Active gameplay
- Pieces fall at regular intervals
- Player controls are active
- Score and level are displayed

### Paused State
- Game is temporarily suspended
- No pieces fall
- Player controls are disabled (except unpause)
- Pause overlay is displayed
- Enter with P key, exit with P key

### Game Over State
- Triggered when new piece cannot spawn
- All player controls disabled
- Final score is displayed
- "Game Over" message shown
- Option to restart (R key) or quit (Q key)

## 7. Controls

### Movement Controls
| Key | Action |
|-----|--------|
| Left Arrow | Move piece left |
| Right Arrow | Move piece right |
| Down Arrow | Soft drop (accelerate descent) |
| Space Bar | Hard drop (instant drop) |

### Rotation Controls
| Key | Action |
|-----|--------|
| X | Rotate clockwise |
| Z | Rotate counter-clockwise |
| Up Arrow | Rotate clockwise (alternative) |

### Game Controls
| Key | Action |
|-----|--------|
| P | Pause/Resume game |
| Q | Quit game |
| Esc | Quit game (alternative) |
| R | Restart game (when game over) |

### Hold Piece (Optional Enhancement)
| Key | Action |
|-----|--------|
| C or Shift | Hold current piece, swap with held piece |

## 8. Level Progression

### Level System
- **Starting Level**: 1
- **Level Increase**: Every 10 lines cleared
- **Maximum Level**: 15 (or unlimited, depending on implementation)

### Speed Progression
As level increases, the gravity (drop speed) increases:

| Level | Lines to Next | Drop Delay (ms) |
|-------|---------------|-----------------|
| 1     | 10            | 1000            |
| 2     | 20            | 900             |
| 3     | 30            | 800             |
| 4     | 40            | 700             |
| 5     | 50            | 600             |
| 6     | 60            | 500             |
| 7     | 70            | 400             |
| 8     | 80            | 300             |
| 9     | 90            | 200             |
| 10+   | +10 each      | 100 (cap)       |

### Scoring Multiplier
- All scoring is multiplied by current level
- Higher levels = more points per action
- Higher levels = faster gameplay

## 9. Hold and Next Piece Features

### Next Piece Preview
- Shows the next tetromino that will spawn
- Displayed in a dedicated UI panel
- Allows player to plan ahead
- Essential for strategic play

### Hold Piece Feature
- Player can "hold" the current piece for later use
- Held piece is stored and can be swapped back
- Can only swap once per turn (cannot swap back immediately)
- Displayed in a dedicated UI panel
- Strategic tool for difficult situations

### 7-Bag Randomizer
To ensure fair piece distribution:
- All 7 tetrominoes are placed in a "bag"
- Pieces are drawn randomly from the bag
- When bag is empty, a new bag is created with all 7 pieces
- Guarantees each piece appears at least once every 7 pieces
- Prevents long droughts of specific pieces

## 10. Terminal UI Specifications

### Screen Layout
```
┌──────────────────────────────────────┐
│  TETRIS                              │
├──────────────────────────────────────┤
│  ┌────────┐    ┌────────┐           │
│  │  HOLD  │    │  NEXT  │           │
│  │        │    │        │           │
│  └────────┘    └────────┘           │
├──────────────────────────────────────┤
│           ┌──────────────┐           │
│           │              │           │
│           │   GAME       │           │
│           │   BOARD      │           │
│           │  (10x20)     │           │
│           │              │           │
│           └──────────────┘           │
├──────────────────────────────────────┤
│  SCORE: 00000    LEVEL: 1            │
│  LINES: 00       PIECES: 00          │
└──────────────────────────────────────┘
```

### Color Specifications (256-Color Terminal)
- **I-Tetromino**: Cyan (color 51)
- **O-Tetromino**: Yellow (color 226)
- **T-Tetromino**: Magenta (color 201)
- **S-Tetromino**: Green (color 46)
- **Z-Tetromino**: Red (color 196)
- **J-Tetromino**: Blue (color 21)
- **L-Tetromino**: Orange (color 208)

### ASCII Rendering
- Each block rendered as 2 characters wide (## for filled, spaces for empty)
- Creates square appearance in terminal (characters are ~2:1 aspect ratio)
- Border characters: ┌ ┐ └ ┘ ─ │

## 11. Technical Requirements

### Minimum Terminal Requirements
- **Terminal Emulator**: Any modern terminal (iTerm2, Terminal.app, xterm, etc.)
- **Colors**: 256-color support recommended
- **Size**: Minimum 80 columns × 30 rows
- **Encoding**: UTF-8 for box-drawing characters

### Performance Targets
- **Frame Rate**: 60 FPS for smooth rendering
- **Input Latency**: <50ms from keypress to action
- **Memory**: <10MB RAM usage

### Compatibility
- **Go Version**: 1.21 or higher
- **Platform**: macOS, Linux, Windows (via WSL or compatible terminal)
- **Library**: tcell v2 for cross-platform terminal handling

## 12. Edge Cases and Special Rules

### Wall Collision
- Pieces cannot move outside the 10-column boundary
- Left wall: x = 0
- Right wall: x = 9

### Floor Collision
- Pieces lock when they reach y = 0
- Pieces lock when they land on another piece

### Ceiling Game Over
- Game over triggers when new piece cannot spawn
- Spawn area must be clear for game to continue
- If spawn area blocked, game ends immediately

### Lock Delay
- Pieces have a brief moment to move after touching ground
- Prevents unfair locks in tight spaces
- Standard delay: 0.5 seconds

---

*Document Version: 1.0*
*Last Updated: 2026-03-11*
*Target Implementation: Go with tcell terminal UI*
