package chess_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"vc-server/chess"
	"vc-server/chess/testutil"
)

// --- IsLegalMove ---

func TestIsLegalMove_ValidAndInvalid(t *testing.T) {
	eng := testutil.MustStandard(t)
	pos := eng.GetPosition()

	from := testutil.MustSquare(t, pos, "e2")
	to := testutil.MustSquare(t, pos, "e4")
	to3 := testutil.MustSquare(t, pos, "e3")

	legalMoves := eng.GetLegalMoves()
	move, ok := chess.MatchMove(from, to, 0, legalMoves)
	require.True(t, ok, "e2e4 should be a legal move")
	require.True(t, eng.IsLegalMove(move))

	// Construct an illegal move (e2→e5, skipping 3 squares).
	illegalMove := chess.Move{From: from, To: to3 - 8, Piece: 'P', ClassicMoveType: chess.QuietMove}
	_ = illegalMove
	// Any move not in legal-move list is illegal.
	bad := chess.Move{From: from, To: testutil.MustSquare(t, pos, "e5"), Piece: 'P'}
	require.False(t, eng.IsLegalMove(bad), "e2e5 should be illegal")
}

func TestPerformMove_IllegalMoveReturnsError(t *testing.T) {
	eng := testutil.MustStandard(t)
	pos := eng.GetPosition()
	from := testutil.MustSquare(t, pos, "e2")
	to := testutil.MustSquare(t, pos, "e5")
	_, err := eng.PerformMove(chess.Move{From: from, To: to, Piece: 'P'})
	require.Error(t, err, "illegal move must return an error")
}

// --- MoveResult fields ---

func TestPerformMove_MoveResult_Capture(t *testing.T) {
	// Scaffold a position where White can immediately capture on e5.
	// FEN: white pawn on d4, black pawn on e5.
	eng := testutil.MustNewEngine(t, chess.GameConfig{
		VariantType: "checkmate",
		FEN:         "rnbqkbnr/pppp1ppp/8/4p3/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 2",
	})
	result := testutil.MustPerformMove(t, eng, "d4e5")
	require.True(t, result.IsCapture, "d4xe5 should be a capture")
}

func TestPerformMove_MoveResult_IsCheck(t *testing.T) {
	// After Fool's Mate sequence Qh4 the king is in check.
	eng := testutil.MustStandard(t)
	testutil.PlaySequence(t, eng, "f2f3", "e7e5", "g2g4")
	// 2...Qh4+ — this move gives check.
	result := testutil.MustPerformMove(t, eng, "d8h4")
	require.True(t, result.IsCheck, "Qh4 gives check")
	require.True(t, result.IsGameOver, "Qh4# ends the game (checkmate)")
	require.NotNil(t, result.Result)
	require.Equal(t, "black", result.Result.Winner)
	require.Equal(t, "checkmate", result.Result.Reason)
}

// --- IsGameOver ---

func TestIsGameOver_Checkmate(t *testing.T) {
	eng := testutil.MustStandard(t)
	testutil.PlaySequence(t, eng, "f2f3", "e7e5", "g2g4", "d8h4")
	over, res := eng.IsGameOver()
	require.True(t, over)
	require.Equal(t, "checkmate", res.Reason)
	require.Equal(t, "black", res.Winner)
}

func TestIsGameOver_Stalemate(t *testing.T) {
	// White queen on d6, White king on c6, Black king on a8, White to move.
	// White plays Qd6→c7: after this move a8-king has no legal squares (a7/b7/b8 all
	// attacked by Qc7) and is not in check → stalemate.
	eng := testutil.MustNewEngine(t, chess.GameConfig{
		VariantType: "checkmate",
		FEN:         "k7/8/2KQ4/8/8/8/8/8 w - - 0 1",
	})
	result := testutil.MustPerformMove(t, eng, "d6c7")
	require.True(t, result.IsGameOver, "Qc7 should cause stalemate")
	require.NotNil(t, result.Result)
	require.Equal(t, "stalemate", result.Result.Reason)
	require.Equal(t, "draw", result.Result.Winner)

	over, res := eng.IsGameOver()
	require.True(t, over)
	require.Equal(t, "stalemate", res.Reason)
}

// --- En passant ---

func TestEnPassantCapture(t *testing.T) {
	// After 1.e4 e5 2.e5→ wait, set up: white pawn on e5, black just pushed d7→d5.
	eng := testutil.MustNewEngine(t, chess.GameConfig{
		VariantType: "checkmate",
		// White pawn on e5, Black just pushed d7→d5 so en passant square is d6.
		FEN: "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3",
	})
	pos := eng.GetPosition()
	epSq := testutil.MustSquare(t, pos, "d6")
	require.Equal(t, epSq, pos.EnPassant, "en passant square should be d6")

	result := testutil.MustPerformMove(t, eng, "e5d6")
	require.True(t, result.IsCapture, "en passant should count as a capture")
}

// --- Promotion ---

func TestPromotion(t *testing.T) {
	// White pawn on e7, Black rook on d8 (capturable), Black king on g8.
	// e7 can promote quietly to e8 (clear) or by capturing d8 rook.
	eng := testutil.MustNewEngine(t, chess.GameConfig{
		VariantType: "checkmate",
		FEN:         "3r2k1/4P3/8/8/8/8/8/4K3 w - - 0 1",
	})
	lm := eng.GetLegalMoves()
	found := false
	for _, m := range lm {
		if m.ClassicMoveType == chess.PromotionMove && m.Promotion == 'Q' {
			found = true
			_, err := eng.PerformMove(m)
			require.NoError(t, err)
			break
		}
	}
	require.True(t, found, "at least one queen promotion must be available")
}

// --- Castling ---

func TestCastling_KingSide(t *testing.T) {
	// White has king-side castling rights, path is clear.
	eng := testutil.MustNewEngine(t, chess.GameConfig{
		VariantType: "checkmate",
		FEN:         "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
	})
	lm := eng.GetLegalMoves()
	found := false
	for _, m := range lm {
		if m.ClassicMoveType == chess.CastleMove {
			found = true
			break
		}
	}
	require.True(t, found, "castling move must be available in clear position")
}

func TestCastling_BlockedWhenInCheck(t *testing.T) {
	// White king in check from a rook on e8 — cannot castle.
	eng := testutil.MustNewEngine(t, chess.GameConfig{
		VariantType: "checkmate",
		FEN:         "4r3/8/8/8/8/8/8/R3K2R w KQ - 0 1",
	})
	lm := eng.GetLegalMoves()
	for _, m := range lm {
		require.NotEqual(t, chess.CastleMove, m.ClassicMoveType, "cannot castle while in check")
	}
}

// --- Clone isolation (deep) ---

func TestClone_DoesNotShareBitboards(t *testing.T) {
	orig := testutil.MustStandard(t)
	clone := orig.Clone()

	testutil.PlaySequence(t, clone, "e2e4", "e7e5")

	origPos := orig.GetPosition()
	clonePos := clone.GetPosition()

	origE4 := testutil.MustSquare(t, origPos, "e4")
	cloneE4 := testutil.MustSquare(t, clonePos, "e4")

	// e4 should be empty on the original but occupied on the clone.
	require.False(t, origPos.Pieces['P'].HasBit(origE4), "original: e4 still empty")
	require.True(t, clonePos.Pieces['P'].HasBit(cloneE4), "clone: e4 has white pawn")
}

// --- Display helpers ---

func TestSquareToAlgebraic_RoundTrip(t *testing.T) {
	pos, err := chess.ParseFEN(testutil.StartFEN)
	require.NoError(t, err)
	cases := []string{"a1", "a8", "h1", "h8", "e4", "d5"}
	for _, sq := range cases {
		idx, err := pos.Square(sq)
		require.NoError(t, err)
		got := chess.SquareToAlgebraic(idx, pos.Files, pos.Ranks)
		require.Equal(t, sq, got, "round-trip for %s", sq)
	}
}

func TestPositionToASCII_DoesNotPanic(t *testing.T) {
	eng := testutil.MustStandard(t)
	s := chess.PositionToASCII(eng.GetPosition())
	require.NotEmpty(t, s)
	require.Contains(t, s, "R", "white rook on starting board")
	require.Contains(t, s, "r", "black rook on starting board")
}

func TestParseAlgebraicMove_Valid(t *testing.T) {
	pos, _ := chess.ParseFEN(testutil.StartFEN)
	from, to, promo, err := chess.ParseAlgebraicMove("e2e4", pos)
	require.NoError(t, err)
	e2 := testutil.MustSquare(t, pos, "e2")
	e4 := testutil.MustSquare(t, pos, "e4")
	require.Equal(t, e2, from)
	require.Equal(t, e4, to)
	require.Equal(t, rune(0), promo)
}
