package chess

import (
	"testing"

	"github.com/stretchr/testify/require"
)


func TestMakeMove(t *testing.T) {
    
	pos,err := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	require.Equal(t,err,nil)
    
	v := &StandardVariant{variant{Position: pos}}
    // Move 1: White Pawn moves two squares forward
	from,_:=v.Position.Square("e2")
	to,_ := v.Position.Square("e4")
	ep,_ := v.Position.Square("e3")
    move1 := Move{From: from, To: to, Piece: 'P'}
    v.MakeMove(move1)
    require.False(t, v.Position.Pieces['P'].HasBit(from), "Pawn should have moved from e2")
    require.True(t, v.Position.Pieces['P'].HasBit(to), "Pawn should be on e4")
    require.Equal(t, ep, v.Position.EnPassant, "En passant square should be set")
	from,_=v.Position.Square("d7")
	to,_ = v.Position.Square("d5")
	ep,_ = v.Position.Square("d6")
	move2  := Move{From: from, To: to, Piece: 'p'}
	v.MakeMove(move2)
	require.False(t, v.Position.Pieces['p'].HasBit(from), "Pawn should have moved from d7")
    require.True(t, v.Position.Pieces['p'].HasBit(to), "Pawn should be on d5")
	require.Equal(t, ep, v.Position.EnPassant, "En passant square should be updated")

	/*
    // Move 2: Black Pawn moves two squares forward
    move2 := Move{From: v.Position.Square("d7"), To: v.Position.Square("d5"), Piece: 'p'}
    v.MakeMove(move2)

    require.Equal(t, v.Position.Square("d6"), v.Position.EnPassant, "En passant square should be updated")

    // Move 3: White captures en passant
    move3 := Move{From: v.Position.Square("e4"), To: v.Position.Square("d6"), Piece: 'P'}
    v.MakeMove(move3)

    require.False(t, v.Position.Pieces['P'].HasBit(v.Position.Square("e4")), "Pawn should have moved from e4")
    require.True(t, v.Position.Pieces['P'].HasBit(v.Position.Square("d6")), "Pawn should be on d6")
    require.False(t, v.Position.Pieces['p'].HasBit(v.Position.Square("d5")), "Captured pawn should be removed from d5")

    // Castling
    v, _ = NewVariant(gameConfig)
    move4 := Move{From: v.Position.Square("e1"), To: v.Position.Square("g1"), Piece: 'K'}
    v.MakeMove(move4)

    require.False(t, v.Position.Pieces['K'].HasBit(v.Position.Square("e1")), "King should have moved from e1")
    require.True(t, v.Position.Pieces['K'].HasBit(v.Position.Square("g1")), "King should be on g1")
    require.False(t, v.Position.Pieces['R'].HasBit(v.Position.Square("h1")), "Rook should have moved from h1")
    require.True(t, v.Position.Pieces['R'].HasBit(v.Position.Square("f1")), "Rook should be on f1")

    // Promotion
    v, _ = NewVariant(gameConfig)
    v.Position.Pieces['P'].SetBit(v.Position.Square("h7"))
    move5 := Move{From: v.Position.Square("h7"), To: v.Position.Square("h8"), Piece: 'P', Promotion: 'Q'}
    v.MakeMove(move5)

    require.False(t, v.Position.Pieces['P'].HasBit(v.Position.Square("h7")), "Pawn should have moved from h7")
    require.False(t, v.Position.Pieces['P'].HasBit(v.Position.Square("h8")), "Pawn should be replaced by promoted piece")
    require.True(t, v.Position.Pieces['Q'].HasBit(v.Position.Square("h8")), "Pawn should have promoted to Queen")

    // Capture resets half-move counter
    v, _ = NewVariant(gameConfig)
    v.Position.Pieces['N'].SetBit(v.Position.Square("f3"))
    v.Position.Pieces['p'].SetBit(v.Position.Square("d4"))

    move6 := Move{From: v.Position.Square("f3"), To: v.Position.Square("d4"), Piece: 'N'}
    v.MakeMove(move6)

    require.False(t, v.Position.Pieces['N'].HasBit(v.Position.Square("f3")), "Knight should have moved from f3")
    require.True(t, v.Position.Pieces['N'].HasBit(v.Position.Square("d4")), "Knight should be on d4")
    require.False(t, v.Position.Pieces['p'].HasBit(v.Position.Square("d4")), "Pawn should be captured")

    require.Equal(t, 0, v.Position.HalfMove, "Half-move counter should reset on capture")*/
}

func TestIsCheckmate(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	pos, _ := ParseFEN(fen)
	cv := &StandardVariant{variant{Position: pos}}

	require.False(t, cv.IsCheckmate())

	// Fool's Mate (1. f3 e5 2. g4 Qh4#)
	cv.PerformMove(Move{From: FileRankToIndex(5, 1, 8), To: FileRankToIndex(5, 3, 8)}) // White f3
	cv.PerformMove(Move{From: FileRankToIndex(4, 6, 8), To: FileRankToIndex(4, 4, 8)}) // Black e5
	cv.PerformMove(Move{From: FileRankToIndex(6, 1, 8), To: FileRankToIndex(6, 3, 8)}) // White g4
	cv.PerformMove(Move{From: FileRankToIndex(3, 7, 8), To: FileRankToIndex(7, 3, 8)}) // Black Qh4#

	require.True(t, cv.IsCheckmate())
}

