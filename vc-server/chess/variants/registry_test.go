package variants_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"vc-server/chess"
	"vc-server/chess/testutil"
)

func TestRegistry_AllKnownVariantsResolve(t *testing.T) {
	cases := []struct {
		name       string
		customData map[string]interface{}
	}{
		{"checkmate", nil},
		{"", nil}, // empty string = standard
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
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := chess.NewEngine(chess.GameConfig{
				FEN:         testutil.StartFEN,
				VariantType: tc.name,
				CustomData:  tc.customData,
			})
			require.NoError(t, err, "variant %q must construct without error", tc.name)
		})
	}
}

func TestRegistry_UnknownVariantReturnsError(t *testing.T) {
	_, err := chess.NewEngine(chess.GameConfig{
		FEN:         testutil.StartFEN,
		VariantType: "this_variant_does_not_exist",
	})
	require.Error(t, err, "unknown variant type must return an error")
}
