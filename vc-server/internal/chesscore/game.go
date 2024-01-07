package chesscore

import (
	"bytes"
)

type classicMoveType uint8

type Color string

type GameObjective string

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
	DoublePawnPush
	EnPassantMove
	PromotionMove

	DuckPlacement variantMoveType = iota
	Teleport

	Checkmate GameObjective = "checkmate"
	NCheck GameObjective = "ncheck"
	Antichess GameObjective = "antichess"
	Targetsquare GameObjective = "target"
	Capture GameObjective = "capture"

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
	notation string
	pieceType
}

type Dimensions struct {
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
	PieceType       pieceType       `json:"piece_type"`
	PieceNotation   string            `json:"piece_notation"`
	ClassicMoveType classicMoveType `json:"classic_move_type"`
	VariantMoveType variantMoveType `json:"variant_move_type,omitempty"`
	AdditionalData interface{} 		`json:"additional_data,omitempty"`
}

type Objective struct {
	Type  GameObjective `json:"type"`
	ObjectiveProps interface{}   `json:"objective_props,omitempty"`
}

func CreateGame(gameConfig GameConfig) (*Game, error) {
	/*var gameConfig GameConfig
	err := json.Unmarshal([]byte(gameConfigJSON), gameConfig)
	if err != nil {
		return nil, err
	}*/
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
		switch gameConfig.Objective.Type {
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

func IsValidConfig (config GameConfig) (bool,error){
	return true,nil
}

func (v *variant) getPosition() Position{
	position := Position{Dimensions: v.dimensions,PieceProps: v.pieceProps}
	var fen bytes.Buffer
	for i:=0;i<v.dimensions.Files;i++{
		for j:=0;j<v.dimensions.Ranks;j++{
			id:= v.toPos(i,j)
			if piece,ok := v.position.pieceLocations[id]; ok{
				fen.WriteString(piece.notation)
			} else if _, ok:= v.position.wallSquares[id];ok{
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

/*func (g *Game) PerformMove(Move) (error){

}*/