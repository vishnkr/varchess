package chesscore

import (
	"fmt"
	"strconv"
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
	fmt.Println("incoming ---------- ", gameConfig)
	variant, err := newVariant(gameConfig)
	if err != nil {
		return nil, err
	}
	game := Game{
		variant: variant,
		History: []Move{},
	}
	ok,_:=game.GetGameConfig()
	fmt.Println("outgoing ---------- ", ok)
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
		variantType: Checkmate,
		position:    position,
	}
	switch variantType {
	case Checkmate:
		newVariant = &CheckmateVariant{variant}
	default:
		newVariant = &CheckmateVariant{variant}
	}

	return newVariant, nil
}

func (g *Game) GetGameConfig()(GameConfig,error){
	return g.variant.GetGameConfig()
}

func (v *variant) GetGameConfig()(GameConfig,error){
	pos:= v.GetPosition()
	gameConfig := GameConfig{
		VariantType: v.variantType,
		Fen: pos.Fen,
		Dimensions: pos.Dimensions,
	}

	return gameConfig,nil
}

func IsValidConfig (config GameConfig) (bool,error){
	return true,nil
}

func (v *variant) GetPosition() Position{
	position := Position{Dimensions: v.dimensions,PieceProps: v.pieceProps}
	var fen string
	for i:=0;i<v.dimensions.Ranks;i++{
		empty := 0
		for j:=0;j<v.dimensions.Files;j++{
			id:= v.toPos(i,j)
			if piece,ok := v.position.pieceLocations[id]; ok{
				if empty>0{
					fen += strconv.Itoa(empty)
				}
				empty = 0
				fen += piece.notation
			} else if _, ok:= v.position.wallSquares[id];ok{
				if empty>0{
					fen += strconv.Itoa(empty)
				}
				empty = 0
				fen += "."
			} else{
				empty+=1
			}
		}
		if (i!=v.dimensions.Ranks-1){
			if empty>0{
				fen += strconv.Itoa(empty)
			}
			empty = 0
			fen+="/"
		}
	}
	if v.turn == ColorWhite{
		fen+=" w"
	} else { fen+=" b"}
	fen+=" KQkq - 0 1"
	position.Fen = fen
	return position
}

/*func (g *Game) PerformMove(Move) (error){

}*/