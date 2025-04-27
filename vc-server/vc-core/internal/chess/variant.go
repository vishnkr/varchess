package chess

import (
	"vc-server/vc-core/internal/models"
)

type Variant interface {
    MakeMove(Move)
    UnmakeMove(Move)
    GetPseudoLegalMoves(color Color, capturesOnly bool) []Move
    GetLegalMoves() []Move
    IsLegalMove(move Move) bool
    PerformMove(move Move) (bool, error) 
    IsCheck() bool
}

type variant struct {
	*Position
	variantType
	recentMove RecentMoveInfo
	possibleLegalMoves []Move
	gameResult result
	isGameOverBool bool
}

func (v *variant) GetPieceAt(pos int) rune{
	var capturedPiece rune
	for p, bitboard := range v.Pieces {
        if bitboard.HasBit(pos) {
            capturedPiece = p
            break
        }
    }
	return capturedPiece
}
func (v *variant) CurrentTurn() Color{
    if v.Position.WhiteToMove{
        return White
    } 
    return Black
}

func opponent(color Color) Color{
    if color==Black{
        return White
    } 
    return Black
}

func (v *variant) SwitchTurn(){
	v.WhiteToMove = !v.WhiteToMove
}
func NewVariant(gameConfig models.GameConfig) (Variant, error) {
	var variantType = Checkmate
	var newVariant Variant
    fen := gameConfig.FEN
	position, err := ParseFEN(fen)
	if err != nil {
		return nil, err
	}
	var variant = variant{
		variantType: variantType,
		Position:    position,
	}
	switch variantType {
	case Checkmate:
		newVariant = &StandardVariant{variant}
	case Antichess:
		//newVariant = &AntichessVariant{variant}
	default:
		//newVariant = &CheckmateVariant{variant}
	}

	return newVariant, nil
}