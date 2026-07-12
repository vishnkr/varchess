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
