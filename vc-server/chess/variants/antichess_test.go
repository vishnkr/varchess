package variants_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"vc-server/chess"
	"vc-server/chess/testutil"
)

func antichessEngine(t *testing.T, fen string) chess.Engine {
	t.Helper()
	return testutil.MustNewEngine(t, chess.GameConfig{
		VariantType: "antichess",
		FEN:         fen,
	})
}

func TestAntichess_MandatoryCapture(t *testing.T) {
	// White pawn on d4 can capture black pawn on e5.
	eng := antichessEngine(t, "rnbqkbnr/pppp1ppp/8/4p3/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 2")
	lm := eng.GetLegalMoves()
	require.NotEmpty(t, lm)
	for _, m := range lm {
		require.True(t, m.Capture, "when captures are available, all legal moves must be captures")
	}
}

func TestAntichess_NoCaptureWhenNoneAvailable(t *testing.T) {
	// Fresh start position — no captures available, all moves are normal.
	eng := antichessEngine(t, testutil.StartFEN)
	lm := eng.GetLegalMoves()
	require.NotEmpty(t, lm)
	hasCap := false
	for _, m := range lm {
		if m.Capture {
			hasCap = true
		}
	}
	require.False(t, hasCap, "no captures at start position")
}

func TestAntichess_NoPiecesWin(t *testing.T) {
	// White pawn on e5, Black pawn on d6 — no kings so no standard check/checkmate rules.
	// After exd6 Black has zero pieces → Black wins (first to lose all pieces wins).
	eng := antichessEngine(t, "8/8/3p4/4P3/8/8/8/8 w - - 0 1")
	result := testutil.MustPerformMove(t, eng, "e5d6")
	require.True(t, result.IsGameOver, "game should end when a side loses all pieces")
	require.NotNil(t, result.Result)
	require.Equal(t, "black", result.Result.Winner, "Black wins in antichess by losing all pieces")
	require.Equal(t, "no-pieces", result.Result.Reason)
}

// Castling rights set but no rooks on a small board used to nil-deref in MakeMove
// while GetLegalMoves simulated castling (seen on antichess minigames).
func TestAntichess_SmallBoardNoRookCastlingDoesNotPanic(t *testing.T) {
	eng := antichessEngine(t, "2bqk/3pp/5/2Q2/1NK2 w KQkq - 0 1")
	require.NotPanics(t, func() {
		_ = eng.GetLegalMoves()
	})
	lm := eng.GetLegalMoves()
	for _, m := range lm {
		require.NotEqual(t, chess.CastleMove, m.ClassicMoveType, "no castling without a rook")
	}
	require.NotPanics(t, func() {
		if len(lm) > 0 {
			_, _ = eng.PerformMove(lm[0])
		}
	})
}

// Non-8 boards pack squares on an 8-wide index grid. Decoding with LargestDimension
// (e.g. 5 on a 4×5 board) made queen slides illegal — including mandatory captures.
func TestAntichess_SmallBoardQueenCapture(t *testing.T) {
	// 4×5: White queen on b2 can capture black bishop on b5 (clear file).
	eng := antichessEngine(t, "1bqk/2pp/4/1Q2/1NK1 w - - 0 1")
	pos := eng.GetPosition()
	from := testutil.MustSquare(t, pos, "b2")
	to := testutil.MustSquare(t, pos, "b5")
	lm := eng.GetLegalMoves()
	found := false
	for _, m := range lm {
		if m.From == from && m.To == to && m.Capture {
			found = true
			_, err := eng.PerformMove(m)
			require.NoError(t, err)
			break
		}
	}
	require.True(t, found, "Qb2xb5 must be legal on 4×5 antichess; legal=%v", lm)
}

// In antichess the king is a normal piece: adjacent king must be captured when
// captures are mandatory (quiet moves like Qd2 are illegal).
func TestAntichess_ForcedKingCapture(t *testing.T) {
	// Black queen on c3, white king on c2 — only legal move is Qxc2.
	eng := antichessEngine(t, "5/5/2q2/2K2/4k b - - 0 1")
	pos := eng.GetPosition()
	from := testutil.MustSquare(t, pos, "c3")
	to := testutil.MustSquare(t, pos, "c2")

	lm := eng.GetLegalMoves()
	require.NotEmpty(t, lm)
	for _, m := range lm {
		require.True(t, m.Capture, "must capture when king is en prise: got %v", m)
		require.Equal(t, from, m.From)
		require.Equal(t, to, m.To)
	}
}
