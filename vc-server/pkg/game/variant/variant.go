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

type VariantMoveType uint8

const (
	NullVariantMove VariantMoveType = iota
	DuckPlacement
)

type Result uint8

const (
	InProgress Result = iota
	White
	Black
	Draw
)

type VariantProps struct {
	variantType VariantType
	board       *game.Board
	turn        game.Color
	validMoves  map[game.Piece][]VariantMove
	gameData    interface{}
	result      Result
}

func (props *VariantProps) switchTurn() {
	if props.turn == game.White {
		props.turn = game.Black
	} else {
		props.turn = game.Black
	}
}

type VariantMove struct {
	Src             int             `json:"src"`
	Dest            int             `json:"dest"`
	MoveType        ClassicMoveType `json:"classicMoveType"`
	VariantMoveType VariantMoveType `json:"variantMoveType,omitempty"`
	VariantMoveInfo interface{}     `json:"variantMoveInfo,omitempty"`
}

type Variant interface {
	IsLegalMove(VariantMove) bool
	GetGameResult() Result
	makeMove(VariantMove) error
	unmakeMove(VariantMove)
	GetPseudoLegalMoves(color game.Color) map[game.Piece][]VariantMove
	UpdateValidMoves()
	PerformMove() error
}

func filterLegalVariantMoves(pseudoMoves *map[game.Piece][]VariantMove, isLegalMove func(VariantMove) bool) {
	for piece, moves := range *pseudoMoves {
		mvs := make([]VariantMove, 0)
		for _, mv := range moves {
			if !isLegalMove(mv) {
				mvs = append(mvs, mv)
			}
		}
		(*pseudoMoves)[piece] = mvs
	}
}

func toPos(row int, col int) int {
	return 8*row + col
}

func toRowCol(pos int) (int, int) {
	return pos % 8, pos / 8
}
