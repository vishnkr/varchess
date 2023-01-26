package variant

import "varchess/pkg/game"

type VariantType uint8
const (
	Duck VariantType = iota
	PoisonedPawn
	Antichess
	Wormhole
)

type ClassicMoveType uint8
const (
	NullClassicMove ClassicMoveType = iota
	Castle 
	Capture
	Quiet
	EnPassant
	Promote
)

type VariantMoveType int
const (
	NullVariantMove VariantMoveType = iota
	DuckPlacement 
)

type VariantProps struct{
	variantType VariantType
	board *game.Board
	turn game.Color
	validMoves []VariantMove
	gameData interface{}
}

func (props *VariantProps) switchTurn(){
	if props.turn == game.White{
		props.turn = game.Black
	} else {
		props.turn = game.Black
	}
}

type VariantMove struct{
	Src int `json:"src"`
	Dest int `json:"dest"`
	MoveType ClassicMoveType `json:"classicMoveType"`
	VariantMoveType VariantMoveType `json:"variantmoveType,omitempty"`
	VariantMoveInfo interface{}
}

type Variant interface{
	IsLegalMove(VariantMove) bool
	IsGameOver() bool
	makeMove(VariantMove) error
	unmakeMove(VariantMove)
	GetPseudoLegalMoves(color game.Color) []VariantMove
	UpdateValidMoves()
	IsKingUnderCheck() bool
	PerformMove() error
}


func filterLegalVariantMoves(moves []VariantMove, isLegalMove func(VariantMove) bool) (ret []VariantMove){
	for _, mv := range moves{
		if isLegalMove(mv){
			ret = append(ret, mv)
		}
	}
	return
}

func toPos(row int, col int) int{
	return 8*row + col
}

func toRowCol(pos int) (int, int){
	return pos%8, pos/8
}
