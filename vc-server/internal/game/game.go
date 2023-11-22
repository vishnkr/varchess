package game

import (
	"bytes"
	"encoding/json"
)

type classicMoveType uint8

type Color string

type GameObjective uint8

type pieceType string

type result uint8

type variantMoveType uint8

type allowedJumpMoveType uint

const (
	ColorBlack Color = "black"
	ColorWhite Color = "white"

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
	x int
	y int
}

type piece struct {
	color    Color
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
	Turn            Color           `json:"turn"`
	PieceType       pieceType       `json:"pieceType"`
	PieceNotation   rune            `json:"pieceNotation"`
	ClassicMoveType classicMoveType `json:"classicMoveType"`
	VariantMoveType variantMoveType `json:"variantMoveType,omitempty"`
}

type Objective struct {
	ObjectiveType  GameObjective `json:"objectiveType"`
	ObjectiveProps interface{}   `json:"objectiveProps,omitempty"`
}

func CreateGame(gameConfigJSON string) (*Game, error) {
	var gameConfig GameConfig
	err := json.Unmarshal([]byte(gameConfigJSON), gameConfig)
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

func (v *variant) getGameConfig()(GameConfig,error){
	gameConfig := GameConfig{VariantType: v.variantType,Objective: v.Objective}
	gameConfig.Position = v.getPosition()
	return gameConfig,nil
}

func (v *variant) getPosition() Position{
	position := Position{Dimensions: v.dimensions,PieceProps: v.pieceProps}
	var fen bytes.Buffer
	for i:=0;i<v.dimensions.Files;i++{
		for j:=0;j<v.dimensions.Ranks;j++{
			id:= v.toPos(i,j)
			if piece,ok := v.position.pieceLocations[id]; ok{
				fen.WriteRune(piece.notation)
			} else if _, ok:= v.position.disabledSquares[id];ok{
				fen.WriteString(".")
			}
		}
		if (i!=v.dimensions.Files-1){fen.WriteString("/")}
	}
	if v.turn == ColorWhite{
		fen.WriteString(" w")
	} else { fen.WriteString(" b")}
	fen.WriteString(" KQkq - 0 1")
	position.Fen = fen.String()
	return position
}