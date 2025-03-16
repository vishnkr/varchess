package chess

import (
	"fmt"
	"unicode"
)

type StandardVariant struct {
    variant
}

func (sv *StandardVariant) IsCheckmate() bool {
    return sv.IsCheck() && len(sv.GetLegalMoves()) == 0
}

func (sv *StandardVariant) IsStalemate() bool {
    return !sv.IsCheck() && len(sv.GetLegalMoves()) == 0
}


func (sv *StandardVariant) GetLegalMoves() []Move {
    pseudoMoves := sv.GetPseudoLegalMoves(sv.CurrentTurn(), false)
    legalMoves := []Move{}

    for _, move := range pseudoMoves {
        sv.MakeMove(move)
        if !sv.IsCheck() {
            legalMoves = append(legalMoves, move)
        }
        sv.UnmakeMove(move)
    }

    return legalMoves
}

func (sv *StandardVariant) PerformMove(move Move) (bool, error) {
    if !sv.IsLegalMove(move) {
        return false, fmt.Errorf("illegal move")
    }

    sv.MakeMove(move)

    if sv.IsCheckmate() {
        sv.isGameOverBool = true
        if sv.CurrentTurn() == White {
            sv.gameResult = BlackWins
        } else {
            sv.gameResult = WhiteWins
        }
    } else if sv.IsStalemate() {
        sv.isGameOverBool = true
        sv.gameResult = Stalemate
    }

    return true, nil
}

func (sv *StandardVariant) MakeMove(move Move) {
    piece := move.Piece
    from, to := move.From, move.To

    var capturedPiece rune
    enPassantSquare := -1
    capturedPiece = sv.GetPieceAt(to)

    if piece == 'P' && to == sv.Position.EnPassant {
        capturedPiece = 'P'
        if sv.CurrentTurn() == White{
            enPassantSquare = to - 8
        } else { enPassantSquare = to + 8}
    }

    sv.recentMove = RecentMoveInfo{
        Piece:         piece,
        squareId:      to,
        moveType:      move.ClassicMoveType,
        capturedPiece: capturedPiece,
        epSquare: enPassantSquare,
        prevCastling:  sv.Position.Castling,
        prevEnPassant: sv.Position.EnPassant,
        prevHalfMove:  sv.Position.HalfMove,
        prevFullMove:  sv.Position.FullMove,
    }

    sv.Position.Pieces[piece].ClearBit(from)
    sv.Position.ColorBitboards[sv.CurrentTurn()].ClearBit(from)

    if capturedPiece != 0 {
        if enPassantSquare != -1 {
            sv.Position.Pieces[capturedPiece].ClearBit(enPassantSquare)
            sv.Position.ColorBitboards[opponent(sv.CurrentTurn())].ClearBit(enPassantSquare)
        } else {
            sv.Position.Pieces[capturedPiece].ClearBit(to)
            sv.Position.ColorBitboards[opponent(sv.CurrentTurn())].ClearBit(to)
        }
    }

    sv.Position.Pieces[piece].SetBit(to)
    sv.Position.ColorBitboards[sv.CurrentTurn()].SetBit(to)

    if move.Promotion != 0 {
        sv.Position.Pieces[piece].ClearBit(to)
        sv.Position.Pieces[move.Promotion].SetBit(to)
    }

    updateGameState(sv, move)
    sv.SwitchTurn()
}


func (sv *StandardVariant) UnmakeMove(move Move) {
    piece := move.Piece
    from, to := move.From, move.To

    sv.Position.Pieces[piece].ClearBit(to)
    sv.Position.ColorBitboards[sv.CurrentTurn()].ClearBit(to)

    sv.Position.Pieces[piece].SetBit(from)
    sv.Position.ColorBitboards[sv.CurrentTurn()].SetBit(from)

    if sv.recentMove.capturedPiece != 0 {
        capturedPiece := sv.recentMove.capturedPiece
        capturedSquare := sv.recentMove.squareId

        sv.Position.Pieces[capturedPiece].SetBit(capturedSquare)
        sv.Position.ColorBitboards[opponent(sv.CurrentTurn())].SetBit(capturedSquare)
    }

    if move.Promotion != 0 {
        sv.Position.Pieces[move.Promotion].ClearBit(to)
        sv.Position.Pieces[piece].SetBit(to)
    }


    sv.Position.Castling = sv.recentMove.prevCastling
    sv.Position.EnPassant = sv.recentMove.prevEnPassant
    sv.Position.HalfMove = sv.recentMove.prevHalfMove
    sv.Position.FullMove = sv.recentMove.prevFullMove

    sv.SwitchTurn()
}


func updateGameState(sv *StandardVariant, move Move) {

    sv.recentMove.moveType = move.ClassicMoveType

    if unicode.ToLower(move.Piece) == 'k' {

        if sv.CurrentTurn() == White {
            sv.Position.Castling &^= (KSCW | QSCW)
        } else {
            sv.Position.Castling &^= (KSCB | QSCB)
        }
    } else if unicode.ToLower(move.Piece) == 'r' {
        /*if move.From == H1 {
            sv.Position.Castling &^= KSCW
        } else if move.From == A1 {
            sv.Position.Castling &^= QSCW
        } else if move.From == H8 {
            sv.Position.Castling &^= KSCB
        } else if move.From == A8 { 
            sv.Position.Castling &^= QSCB
        }*/
    }

    if unicode.ToLower(move.Piece) == 'p' && abs(move.From-move.To) == sv.Files*2 {

        sv.Position.EnPassant = (move.From + move.To) / 2
    } else {
        sv.Position.EnPassant = -1
    }

    capturedPiece := rune(0)
    for p, bitboard := range sv.Position.Pieces {
        if bitboard.HasBit(move.To) {
            capturedPiece = p
            break
        }
    }
    sv.recentMove.capturedPiece = capturedPiece

    if move.Piece == 'P' || capturedPiece != 0 {
        sv.Position.HalfMove = 0
    } else {
        sv.Position.HalfMove++
    }
    if !sv.Position.WhiteToMove {
        sv.Position.FullMove++
    }
}


func restoreGameState(sv *StandardVariant, move Move) {
    sv.Position.Castling = sv.recentMove.prevCastling
    sv.Position.EnPassant = sv.recentMove.prevEnPassant
    sv.Position.HalfMove = sv.recentMove.prevHalfMove
    sv.Position.FullMove = sv.recentMove.prevFullMove
}



func (sv *StandardVariant)  GetPseudoLegalMoves(color Color, canCaptureKing bool) []Move {
    moves := []Move{}
    for piece, bitboard := range sv.Pieces {
        if color == White && unicode.IsUpper(piece) || color == Black && unicode.IsLower(piece) {
            piecePositions := bitboard.GetSetBits() 
			for _, pos := range piecePositions {
				pieceMoves := GenerateMovesForPiece(piece, pos, sv.Position, canCaptureKing)
				moves = append(moves, pieceMoves...)
			}
        }
    }
    return moves
}

func (sv *StandardVariant) IsLegalMove(move Move) bool {
    pseudoMoves := sv.GetPseudoLegalMoves(sv.CurrentTurn(), false)
    
    isPseudoLegal := false
    for _, pseudoMove := range pseudoMoves {
        if pseudoMove == move {
            isPseudoLegal = true
            break
        }
    }
    if !isPseudoLegal {
        return false
    }

    sv.MakeMove(move)
    inCheck := sv.IsCheck()
    sv.UnmakeMove(move)

    return !inCheck
}

func (sv *StandardVariant) IsCheck() bool {
    //kingPos := FindKing(sv.Position, sv.CurrentTurn())
    /*opponentMoves := sv.GetPseudoLegalMoves(Opponent(sv.CurrentTurn()), true)

    for _, move := range opponentMoves {
        if move.To == kingPos {
            return true
        }
    }*/
    return false
}
