package chess_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"vc-server/chess"
	"vc-server/chess/testutil"
	_ "vc-server/chess/variants"
)

func TestNewEngine_FENRequired(t *testing.T) {
	_, err := chess.NewEngine(chess.GameConfig{VariantType: "checkmate"})
	require.Error(t, err, "NewEngine with empty FEN must return an error")
}

func TestNewEngine_InvalidFEN(t *testing.T) {
	_, err := chess.NewEngine(chess.GameConfig{FEN: "not/a/fen"})
	require.Error(t, err)
}

func TestNewEngine_UnknownVariantErrors(t *testing.T) {
	// With the variants package registered, unknown types should error.
	_, err := chess.NewEngine(chess.GameConfig{
		FEN:         testutil.StartFEN,
		VariantType: "nonexistent_variant_xyz",
	})
	require.Error(t, err, "unknown variant type should return an error when factory is registered")
}

func TestNewEngine_StartPosition_LegalMoves(t *testing.T) {
	eng := testutil.MustStandard(t)
	require.Len(t, eng.GetLegalMoves(), 20, "20 legal moves from start position")
}

func TestNewEngine_GetPosition(t *testing.T) {
	eng := testutil.MustStandard(t)
	pos := eng.GetPosition()
	require.NotNil(t, pos)
	require.Equal(t, 8, pos.Files)
	require.Equal(t, 8, pos.Ranks)
	require.True(t, pos.WhiteToMove, "white moves first")
}

func TestNewEngine_IsGameOver_StartPosition(t *testing.T) {
	eng := testutil.MustStandard(t)
	over, _ := eng.IsGameOver()
	require.False(t, over, "game not over at start")
}

func TestNewEngine_Clone_Independence(t *testing.T) {
	orig := testutil.MustStandard(t)
	clone := orig.Clone()

	// Play a move on the clone only.
	testutil.MustPerformMove(t, clone, "e2e4")

	// Original should be unchanged.
	require.Len(t, orig.GetLegalMoves(), 20, "original still has 20 moves after clone moves")
	over, _ := orig.IsGameOver()
	require.False(t, over)
}

func TestNewEngine_AllKnownVariants_Construct(t *testing.T) {
	variants := []struct {
		name       string
		customData map[string]interface{}
	}{
		{"checkmate", nil},
		{"antichess", nil},
		{"ncheck", map[string]interface{}{"targetChecks": 3}},
		{"archerchess", nil},
		{"wormhole", nil},
		{"teleport", nil},
		{"poisonedpawn", nil},
		{"hiddenqueen", nil},
		{"footballchess", nil},
		{"spychess", nil},
	}
	for _, v := range variants {
		t.Run(v.name, func(t *testing.T) {
			_, err := chess.NewEngine(chess.GameConfig{
				FEN:         testutil.StartFEN,
				VariantType: v.name,
				CustomData:  v.customData,
			})
			require.NoError(t, err, "variant %q should construct without error", v.name)
		})
	}
}
