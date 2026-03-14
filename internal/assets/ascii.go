package assets

const (
	BorderTL = "┌"
	BorderTR = "┐"
	BorderBL = "└"
	BorderBR = "┘"
	BorderH  = "─"
	BorderV  = "│"
)

const (
	LabelNext  = "NEXT"
	LabelHold  = "HOLD"
	LabelScore = "SCORE"
	LabelLevel = "LEVEL"
)

var (
	BorderTop    = BorderTL + BorderH + BorderH + BorderH + BorderH + BorderH + BorderH + BorderTR
	BorderMiddle = BorderV + "                  " + BorderV
	BorderBottom = BorderBL + BorderH + BorderH + BorderH + BorderH + BorderH + BorderH + BorderBR
)

type Piece [4][4]string

var (
	PieceI = Piece{
		{"  ", "  ", "  ", "  "},
		{"##", "##", "##", "##"},
		{"  ", "  ", "  ", "  "},
		{"  ", "  ", "  ", "  "},
	}
	PieceO = Piece{
		{"##", "##", "  ", "  "},
		{"##", "##", "  ", "  "},
		{"  ", "  ", "  ", "  "},
		{"  ", "  ", "  ", "  "},
	}
	PieceT = Piece{
		{"##", "  ", "  ", "  "},
		{"##", "##", "  ", "  "},
		{"  ", "##", "  ", "  "},
		{"  ", "  ", "  ", "  "},
	}
	PieceS = Piece{
		{"  ", "##", "##", "  "},
		{"##", "##", "  ", "  "},
		{"  ", "  ", "  ", "  "},
		{"  ", "  ", "  ", "  "},
	}
	PieceZ = Piece{
		{"##", "##", "  ", "  "},
		{"  ", "##", "##", "  "},
		{"  ", "  ", "  ", "  "},
		{"  ", "  ", "  ", "  "},
	}
	PieceJ = Piece{
		{"##", "  ", "  ", "  "},
		{"##", "  ", "  ", "  "},
		{"##", "##", "  ", "  "},
		{"  ", "  ", "  ", "  "},
	}
	PieceL = Piece{
		{"  ", "##", "  ", "  "},
		{"  ", "##", "  ", "  "},
		{"##", "##", "  ", "  "},
		{"  ", "  ", "  ", "  "},
	}
)

var (
	HUDBorderTopLine    = BorderTop
	HUDBorderMidLine    = BorderMiddle
	HUDBorderBottomLine = BorderBottom
)

var (
	GameOverScreen = BorderTop + BorderH + BorderH + BorderH + BorderH + BorderH + BorderH + BorderTR + "\n" +
		"|" + "         GAME OVER         " + "|" + "\n" +
		BorderBottom
	PauseScreen = BorderTop + BorderH + BorderH + BorderH + BorderH + BorderH + BorderH + BorderTR + "\n" +
		"|" + "         PAUSED            " + "|" + "\n" +
		BorderBottom
)
