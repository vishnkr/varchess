package game

import "encoding/json"

type classicMoveType uint8

type color string

type GameObjective uint8

type pieceType string

type result uint8

type variantMoveType uint8

type allowedJumpMoveType uint

const (
	ColorBlack color = "black"
	ColorWhite color = "white"

	Pawn        pieceType = "pawn"
	Knight      pieceType = "knight"
	Bishop      pieceType = "bishop"
	Rook        pieceType = "rook"
	Queen       pieceType = "queen"
	King        pieceType = "king"
	Duck        pieceType = "duck"
	CustomPiece pieceType = "custom"

	WhiteWins result = iota
	BlackWins
	Stalemate

	NullMove classicMoveType = iota
	CastleMove
	CaptureMove
	QuietMove
	EnPassantMove
	PromotionMove

	DuckPlacement variantMoveType = iota
	Teleport

	Checkmate GameObjective = iota
	NCheck
	Antichess
	Targetsquare
	Capture

	CaptureOnlyJump allowedJumpMoveType = iota
	QuietOnlyJump
	AllJumps
)

type moveOffset struct {
	xOffset int
	yOffset int
}

type piece struct {
	color    color
	notation rune
	pieceType
}

type dimensions struct {
	Ranks int `json:"ranks"`
	Files int `json:"files"`
}

type Game struct {
	History []Move
	variant Variant
	Result  result
}

type Move struct {
	Source          int             `json:"source"`
	Target          int             `json:"target"`
	Turn            color           `json:"turn"`
	PieceType       pieceType       `json:"pieceType"`
	PieceNotation   rune            `json:"pieceNotation"`
	ClassicMoveType classicMoveType `json:"classicMoveType"`
	VariantMoveType variantMoveType `json:"variantMoveType,omitempty"`
}

type Objective struct {
	ObjectiveType  GameObjective `json:"type"`
	ObjectiveProps interface{}   `json:"objectiveProps"`
}

func CreateGame(gameConfigString string) (*Game, error) {
	var gameConfig GameConfig
	err := json.Unmarshal([]byte(gameConfigString), gameConfig)
	if err != nil {
		return nil, err
	}
	variant, err := newVariant(gameConfig)
	if err != nil {
		return nil, err
	}
	game := Game{
		variant: variant,
		History: []Move{},
	}
	return &game, nil
}

func newVariant(gameConfig GameConfig) (Variant, error) {
	var variantType = gameConfig.VariantType
	var newVariant Variant
	position, err := newPosition(gameConfig)
	if err != nil {
		return nil, err
	}
	var variant = variant{
		Objective:   gameConfig.Objective,
		variantType: Custom,
		position:    position,
	}
	switch variantType {
	case Custom:
		switch gameConfig.Objective.ObjectiveType {
		case Antichess:
			newVariant = &AntichessVariant{variant}
		case NCheck:
			newVariant = &NCheckVariant{variant: variant, blackKingCheckCount: 0, whiteKingCheckCount: 0, targetChecks: 3}
		default:
			newVariant = &CheckmateVariant{variant}
		}
	default:
		newVariant = &CheckmateVariant{variant}
	}

	return newVariant, nil
}
