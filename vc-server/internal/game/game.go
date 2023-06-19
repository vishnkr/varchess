package game

import "encoding/json"

type ClassicMoveType uint8

type Color string

type GameObjective uint8

type PieceType int

type Result uint8

type VariantMoveType uint8

type AllowedJumpMoveType uint

const (
	ColorBlack Color = "black"
	ColorWhite Color = "white"

	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
	Duck
	CustomPiece

	ResultInProgress Result = iota
	ResultWhite
	ResultBlack
	ResultDraw

	NullMove ClassicMoveType = iota
	CastleMove
	CaptureMove
	QuietMove
	EnPassantMove
	PromotionMove

	DuckPlacement VariantMoveType = iota
	Teleport

	Checkmate GameObjective = iota
	NCheck
	Antichess
	Targetsquare
	Capture

	CaptureOnlyJump AllowedJumpMoveType = iota
	QuietOnlyJump
	AllJumps
)

type MoveOffset struct{
	XOffset int
	YOffset int
}

type Piece struct{
	Color Color
	Name string
	Piecetype PieceType
}

type PieceProps struct{
	SlideOffsets map[MoveOffset]bool
	CanDoubleJump bool
	DoubleJumpSquares map[int]bool
	JumpProps map[MoveOffset]AllowedJumpMoveType
	PromotionProps PromotionProps
}

type PromotionProps struct{
	PromotionSquares []int
	CanPromoteTo []Piece
}

type Dimensions struct{
	Ranks int
	Files int
}

type Position struct{
	PieceLocations map[int]Piece
	PieceProps map[Piece]PieceProps
	Turn Color
	Dimensions Dimensions
	DisabledSquares map[int]bool
	CastlingRights uint8
}

type Game struct{
	History []Move
	Variant Variant
	Result Result
}


type Move struct{
	Source int `json:"source"`
	Target int `json:"target"`
	Turn Color `json:"turn"`
	PieceType PieceType `json:"pieceType"`
	PieceName string `json:"pieceName"`
	ClassicMoveType ClassicMoveType `json:"classicMoveType"`
	VariantMoveType VariantMoveType `json:"variantMoveType,omitempty"`
}


type Objective struct{
	ObjType GameObjective `json:"type"`
	ObjectiveProps interface{} `json:"objectiveProps"`
}


func CreateGame(gameConfig string) (*Game,error){
	var gameConfigMap map[string]interface{};
	err:= json.Unmarshal([]byte(gameConfig),gameConfigMap)
	if err!=nil{
		return nil,err
	}

	variant,err := newVariant(gameConfigMap)
	if err!=nil{
		return nil,err
	}
	game:= Game{
		Variant: variant,
		History: []Move{},
	}
	return &game,nil

}

func newPosition(gameConfigMap map[string]interface{}) (Position,error){
	var position = Position{}
	return position,nil
}
func newVariant(gameConfigMap map[string]interface{}) (Variant,error){
	var variantType = gameConfigMap["variantType"].(VariantType)
	var variant Variant
	switch variantType {
	case Custom:
		variant = CustomVariant{}
	case DuckChess:
		//
	case PoisonedPawn:
		//
	case Wormhole:
		//
	default:
		// Handle unrecognized variantType
	}
	return variant,nil
}