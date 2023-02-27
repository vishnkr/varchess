package variant

import (
	"errors"
	"log"
	"varchess/pkg/game"
)

type DuckVariant struct{
	variantProps VariantProps
}

type DuckVariantData struct{
    duckPosition int
}

func NewDuckVariant() *DuckVariant {
    duck:= &DuckVariant{
        variantProps: VariantProps{
            variantType: Duck,
            board: setupDuckVariantBoard(),
            turn: game.White,
            validMoves: make(map[game.Piece][]VariantMove),
            gameData: DuckVariantData{
                duckPosition: -1,
            },
            result: InProgress,
        },
    }
    duck.UpdateValidMoves()
    return duck
}

func setupDuckVariantBoard() *game.Board{
	const startFen string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 1 7"
	return game.ConvertFENtoBoard(startFen)
}

func (duck *DuckVariant) UpdateValidMoves(){
    duck.variantProps.validMoves = duck.GetPseudoLegalMoves(duck.variantProps.turn)
    filterLegalVariantMoves(&duck.variantProps.validMoves,duck.IsLegalMove)
}

func (duck *DuckVariant) GetPseudoLegalMoves(color game.Color) map[game.Piece][]VariantMove{
    var board *game.Board = duck.variantProps.board
    ret := make(map[game.Piece][]VariantMove)
    for rowIndex, row := range board.Tiles {
		for colIndex, tile := range row {
			if !tile.IsEmpty && tile.Piece.Color == color {
				for k, v := range board.GenPieceMoves(tile.Piece, rowIndex, colIndex) {
                    if len(ret[v]) !=0{
                        ret[v] = append(ret[v], 
                            VariantMove{
                                Src: toPos(k.SrcRow,k.SrcCol),
                                Dest: toPos(k.DestRow,k.DestCol),
                                MoveType: duck.getMoveType(k),
                                VariantMoveType: NullVariantMove,
                            })
                    }
				}
			}
		}
	}
    
    return ret
}

func (duck *DuckVariant) getMoveType(mv *game.Move) ClassicMoveType{
    if mv.Promote != game.Empty { return Promote }
    if mv.Castle { return Castle }
    if !duck.variantProps.board.IsEmpty(mv.DestRow,mv.DestCol) { return Capture }

    return Quiet
}

func (duck *DuckVariant) IsLegalMove(mv VariantMove) bool{
    valid := false
    duck.makeMove(mv)
    if !duck.isDuckPositionValid(mv){
        valid = false
    }
    err := duck.unMakeMove(mv)
    if err != nil{
        log.Fatalf("Invalid Variant Unmake %+v\n",mv)
    }
    return valid
}

func (duck *DuckVariant) makeMove(mv VariantMove) error{
    if mv.VariantMoveType != DuckPlacement {
        return errors.New("DuckChess: make move - missing duck placement position")
    }
    srcRow,srcCol := toRowCol(mv.Src)
    destRow,destCol:= toRowCol(mv.Dest)
    if mv.MoveType != NullClassicMove {
        piece:= duck.variantProps.board.Tiles[srcRow][srcCol].Piece
        duck.variantProps.board.PerformMove(piece, game.Move{
            SrcRow: srcRow,
            SrcCol: srcCol,
            DestRow: destRow,
            DestCol: destCol,
        })
    } else{
        return errors.New("DuckChess: make move - invalid classic move")
    }

    curDuckPos := duck.variantProps.gameData.(DuckVariantData).duckPosition
    newDuckPos := mv.VariantMoveInfo.(DuckVariantData).duckPosition
    if  curDuckPos != -1{
        duckRow,duckCol := toRowCol(curDuckPos)
        duck.variantProps.board.Tiles[duckRow][duckCol].IsDisabled = false
        duck.variantProps.board.Tiles[duckRow][duckCol].IsEmpty = true
    }
    newDuckRow,newDuckCol := toRowCol(newDuckPos)
    duck.variantProps.board.Tiles[newDuckRow][newDuckCol].IsDisabled = true
    duck.variantProps.board.Tiles[newDuckRow][newDuckCol].IsEmpty = true
    
    return nil
}



func (duck *DuckVariant) unMakeMove(mv VariantMove) error{
    if mv.VariantMoveType != DuckPlacement {
        return errors.New("DuckChess: unmake move - missing duck placement position")
    }
    srcRow,srcCol := toRowCol(mv.Src)
    destRow,destCol:= toRowCol(mv.Dest)
    if mv.MoveType != NullClassicMove{
        piece:= duck.variantProps.board.Tiles[destRow][destCol].Piece
        duck.variantProps.board.PerformMove(piece, game.Move{
            SrcRow: destRow,
            SrcCol: destCol,
            DestRow: srcRow,
            DestCol: srcCol,
        })
    }  else{
        return errors.New("DuckChess: unmake move - invalid classic move")
    }
    oldDuckPos := duck.variantProps.gameData.(DuckVariantData).duckPosition
    curDuckPos := mv.VariantMoveInfo.(DuckVariantData).duckPosition
    if  oldDuckPos != -1{
        duckRow,duckCol := toRowCol(oldDuckPos)
        duck.variantProps.board.Tiles[duckRow][duckCol].IsDisabled = true
        duck.variantProps.board.Tiles[duckRow][duckCol].IsEmpty = true
    }
    curDuckRow,curDuckCol := toRowCol(curDuckPos)
    duck.variantProps.board.Tiles[curDuckRow][curDuckCol].IsDisabled = false
    duck.variantProps.board.Tiles[curDuckRow][curDuckCol].IsEmpty = true
    return nil
}

func (duck *DuckVariant) PerformMove(mv VariantMove){
    duck.makeMove(mv)
    destRow, destCol := toRowCol(mv.Dest)
    if mv.MoveType == Capture && duck.variantProps.board.Tiles[destRow][destCol].Piece.Type == game.King{
        switch duck.variantProps.turn{
        case game.White: 
            duck.variantProps.result = White
        case game.Black:
            duck.variantProps.result = Black
        }
    }
    duck.variantProps.switchTurn()
    duck.UpdateValidMoves()
}

func (duck *DuckVariant) isDuckPositionValid(mv VariantMove) bool{
    if mv.VariantMoveType == DuckPlacement{
        pos, ok :=  mv.VariantMoveInfo.(int)
        if !ok{
            return false
        }
        row, col:= toRowCol(pos)
        return duck.variantProps.board.IsEmpty(row,col) && pos != duck.variantProps.gameData.(DuckVariantData).duckPosition
    }
    return false
}

func (duck *DuckVariant) GetGameResult() Result{
    return duck.variantProps.result
}