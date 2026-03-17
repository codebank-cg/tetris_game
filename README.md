# 🎮 Tetris Game - Go + tcell v2

> **📢 Announcement**: This is an open-source terminal-based Tetris game implementation. All code was AI-generated using OpenCode with the qwen3.5-plus model. Built with Go and the tcell v2 library for cross-platform terminal support.
>
> **Repository**: [github.com/codebank-cg/tetris_game](https://github.com/codebank-cg/tetris_game.git)
>
> ---

Go version detected: go1.26.0 darwin/amd64

Prerequisites
- Go 1.21 or newer
- A POSIX-compatible terminal (macOS/Linux) or Windows Terminal with ANSI color support
- A terminal configured to support 256 colors (TERM=xterm-256color or equivalent)

Installation
- Initialize modules (if not already done)
  - go mod tidy
- Install dependencies
  - go get github.com/gdamore/tcell/v2
- Ensure you have a main package. If your repo uses a different structure, adjust accordingly.

How to Run
- If this project contains a main package at the root:
  - go run .
- Build to a binary:
  - go build -o tetris
  - ./tetris
- If the main package is under a subdirectory (e.g., cmd/tetris):
  - go run ./cmd/tetris

Controls Reference
------------------
| Key(s)           | Action                |
- Left Arrow       | Move piece left       |
- Right Arrow      | Move piece right      |
- Down Arrow       | Soft drop             |
- Space            | Hard drop / slam      |
- Z                | Rotate left           |
- X                | Rotate right          |
- P                | Pause/Resume          |
- Q                | Quit                  |
- Esc              | Quit (alternative)    |

Notes
- This README is a starting template. If your project uses a different run path or module path, update the commands accordingly.
- The tcell docs are available at https://github.com/gdamore/tcell and https://pkg.go.dev/github.com/gdamore/tcell/v2
