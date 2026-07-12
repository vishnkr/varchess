package chess

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeMove(t *testing.T) {
	pos, err := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	require.NoError(t, err)

	v := &StandardVariant{variant{Position: pos}}

	// Move 1: White pawn e2→e4 (double push).
	// e2 = file 4, chess rank 2, internal rank 8-2=6 → 6*8+4=52
	// e4 = file 4, chess rank 4, internal rank 8-4=4 → 4*8+4=36
	// e3 = file 4, chess rank 3, internal rank 8-3=5 → 5*8+4=44 (en passant square)
	from, _ := v.Position.Square("e2")
	to, _ := v.Position.Square("e4")
	ep, _ := v.Position.Square("e3")

	move1 := Move{From: from, To: to, Piece: 'P', ClassicMoveType: DoublePawnPush}
	v.MakeMove(move1)
	require.False(t, v.Position.Pieces['P'].HasBit(from), "Pawn should have moved from e2")
	require.True(t, v.Position.Pieces['P'].HasBit(to), "Pawn should be on e4")
	require.Equal(t, ep, v.Position.EnPassant, "En passant square should be e3")

	// Move 2: Black pawn d7→d5 (double push).
	// d7 = file 3, chess rank 7, internal rank 8-7=1 → 1*8+3=11
	// d5 = file 3, chess rank 5, internal rank 8-5=3 → 3*8+3=27
	// d6 = file 3, chess rank 6, internal rank 8-6=2 → 2*8+3=19
	from, _ = v.Position.Square("d7")
	to, _ = v.Position.Square("d5")
	ep, _ = v.Position.Square("d6")
	move2 := Move{From: from, To: to, Piece: 'p', ClassicMoveType: DoublePawnPush}
	v.MakeMove(move2)
	require.False(t, v.Position.Pieces['p'].HasBit(from), "Pawn should have moved from d7")
	require.True(t, v.Position.Pieces['p'].HasBit(to), "Pawn should be on d5")
	require.Equal(t, ep, v.Position.EnPassant, "En passant square should be d6")
}

func TestIsCheckmate(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	pos, _ := ParseFEN(fen)
	cv := &StandardVariant{variant{Position: pos}}

	require.False(t, cv.IsCheckmate(), "Starting position is not checkmate")

	// Fool's Mate: 1.f3 e5 2.g4 Qh4#
	// Rank-from-top convention:
	// f2 = file 5, chess rank 2, internal rank 6 → 6*8+5=53
	// f3 = file 5, chess rank 3, internal rank 5 → 5*8+5=45
	// e7 = file 4, chess rank 7, internal rank 1 → 1*8+4=12
	// e5 = file 4, chess rank 5, internal rank 3 → 3*8+4=28
	// g2 = file 6, chess rank 2, internal rank 6 → 6*8+6=54
	// g4 = file 6, chess rank 4, internal rank 4 → 4*8+6=38
	// d8 = file 3, chess rank 8, internal rank 0 → 0*8+3=3
	// h4 = file 7, chess rank 4, internal rank 4 → 4*8+7=39

	_, _ = cv.PerformMove(Move{From: 53, To: 45, Piece: 'P', ClassicMoveType: QuietMove})    // 1. f3
	_, _ = cv.PerformMove(Move{From: 12, To: 28, Piece: 'p', ClassicMoveType: DoublePawnPush}) // 1... e5
	_, _ = cv.PerformMove(Move{From: 54, To: 38, Piece: 'P', ClassicMoveType: DoublePawnPush}) // 2. g4
	_, _ = cv.PerformMove(Move{From: 3, To: 39, Piece: 'q', ClassicMoveType: QuietMove})       // 2... Qh4#

	require.True(t, cv.IsCheckmate(), "Position after Fool's Mate should be checkmate")
}
