package variant

import (
	"log"
	"varchess/pkg/game"
)

type DuckVariant struct{
	variantProps VariantProps
	additionalData interface{}
}

func NewDuckVariant() *DuckVariant {
    duck:= &DuckVariant{
        variantProps: VariantProps{
            variantType: Duck,
            board: setupDuckVariantBoard(),
            turn: game.White,
            validMoves: []VariantMove{},
            gameData: nil,
        },
        additionalData: nil,
    }
    duck.UpdateValidMoves()
    return duck
}

func setupDuckVariantBoard() *game.Board{
	const startFen string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 1 7"
	return game.ConvertFENtoBoard(startFen)
}

func (duck *DuckVariant) UpdateValidMoves(){
    duck.variantProps.validMoves = filterLegalVariantMoves(duck.GetPseudoLegalMoves(duck.variantProps.turn),duck.IsLegalMove)
}

func (duck *DuckVariant) GetPseudoLegalMoves(color game.Color) []VariantMove{
    //moves := duck.variantProps.board.GetAllPseudoLegalMoves(duck.variantProps.turn)
    ret := []VariantMove{}
    return ret
}

func (duck *DuckVariant) IsKingUnderCheck() bool{
    var kingPos int
	if duck.variantProps.turn == game.White {
		kingPos = toPos(duck.variantProps.board.WhiteKing.Position[0],duck.variantProps.board.WhiteKing.Position[1])
	} else {
		kingPos = toPos(duck.variantProps.board.BlackKing.Position[0],duck.variantProps.board.BlackKing.Position[1])
	}
    oppMoves := duck.GetPseudoLegalMoves(game.GetOpponentColor(duck.variantProps.turn))
	for _,move := range oppMoves {
		if kingPos == move.Dest{
			return true
		}
	}
	return false
}

func (duck *DuckVariant) IsLegalMove(mv VariantMove) bool{
    valid := false
    duck.MakeMove(mv)
    if !duck.IsKingUnderCheck(){
        valid = false
    }
    err := duck.UnMakeMove(mv)
    if err != nil{
        log.Fatalf("Invalid Variant Unmake %+v\n",mv)
    }
    return valid
}

func (duck *DuckVariant) MakeMove(mv VariantMove) error{
    if mv.MoveType != NullClassicMove {
        
    }
    if mv.VariantMoveType != NullVariantMove{

    }
    return nil
}

func (duck *DuckVariant) UnMakeMove(mv VariantMove) error{
    if mv.MoveType != NullClassicMove{
        
    }
    return nil
}

func (duck *DuckVariant) PerformMove(mv VariantMove){
    duck.MakeMove(mv)
    duck.variantProps.switchTurn()
    duck.UpdateValidMoves()
}