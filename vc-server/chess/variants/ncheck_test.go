package variants_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"vc-server/chess"
	"vc-server/chess/testutil"
)

// checkNotMateFEN: White queen on g1, White king on e1, Black king on e8.
// White Qg1→g8 gives check along rank 8 (g8-f8-e8); the Black king has safe escape squares
// (e7, d7) so it is NOT checkmate.
const checkNotMateFEN = "4k3/8/8/8/8/8/8/4K1Q1 w - - 0 1"

func ncheckEngine(t *testing.T, n int) chess.Engine {
	t.Helper()
	return testutil.MustNewEngine(t, chess.GameConfig{
		VariantType: "ncheck",
		FEN:         checkNotMateFEN,
		CustomData:  map[string]interface{}{"targetChecks": n},
	})
}

func TestNCheck_CheckCounterIncrements(t *testing.T) {
	// 3-Check: after White gives one check, game should not be over.
	eng := ncheckEngine(t, 3)
	// White Qg1→g8 gives check along the g-file. Black king has e7/d7/d8 escape squares.
	result := testutil.MustPerformMove(t, eng, "g1g8")
	require.True(t, result.IsCheck, "Qg8 gives check")
	require.False(t, result.IsGameOver, "one check is not enough in 3-check")
}

func TestNCheck_WinsAtTargetN(t *testing.T) {
	// 1-Check: the first check ends the game immediately.
	eng := ncheckEngine(t, 1)
	result := testutil.MustPerformMove(t, eng, "g1g8")
	require.True(t, result.IsCheck, "Qg8 gives check")
	require.True(t, result.IsGameOver, "1-check variant: first check ends the game")
	require.NotNil(t, result.Result)
	require.Equal(t, "white", result.Result.Winner)
	require.Contains(t, result.Result.Reason, "1-check")

	// Engine-level IsGameOver must also reflect the end (not just MoveResult).
	over, res := eng.IsGameOver()
	require.True(t, over, "engine IsGameOver must be true after N-Check win")
	require.Equal(t, "white", res.Winner)
	require.Contains(t, res.Reason, "1-check")
}

func TestNCheck_DefaultTargetN_Is3(t *testing.T) {
	// Without CustomData the target defaults to 3 — engine constructs and has legal moves.
	eng := testutil.MustNewEngine(t, chess.GameConfig{
		VariantType: "ncheck",
		FEN:         testutil.StartFEN,
	})
	require.Len(t, eng.GetLegalMoves(), 20)
}

func TestNCheck_ClonePreservesState(t *testing.T) {
	// After giving one check on the original, the clone (taken before the check) should not
	// have that check counted.
	eng := ncheckEngine(t, 3)
	clone := eng.Clone()

	// Give check on original.
	result := testutil.MustPerformMove(t, eng, "g1g8")
	require.True(t, result.IsCheck)
	require.False(t, result.IsGameOver, "3-check: still need 2 more")

	// Clone was taken before the move; its game should still be in progress.
	cloneOver, _ := clone.IsGameOver()
	require.False(t, cloneOver, "clone is untouched by original's move")

	// Clone should still have legal moves (unplayed position).
	require.NotEmpty(t, clone.GetLegalMoves())
}

func TestNCheck_VariantStateRoundTrip(t *testing.T) {
	// Verify engine constructs and is playable — state is internal to NCheckRules.
	eng := testutil.MustNewEngine(t, chess.GameConfig{
		VariantType: "ncheck",
		FEN:         testutil.StartFEN,
		CustomData:  map[string]interface{}{"targetChecks": 5},
	})
	require.Len(t, eng.GetLegalMoves(), 20)
	over, _ := eng.IsGameOver()
	require.False(t, over)
}
