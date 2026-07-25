package variants_test

import (
	"testing"

	"vc-server/chess"
	"vc-server/chess/variants"
	_ "vc-server/chess/variants"

	"github.com/stretchr/testify/require"
)

func TestWormholeTeleportAndCooldown(t *testing.T) {
	// 8x8 empty-ish: white knight on b1, portals on c3 (18) and f6 (45) in large-8 indexing.
	// rank0 file2 = 2; wait: FileRankToLargeIndex
	// Internal rank: FEN rank 8 is internal 0.
	// c3 chess = file 2, chess rank 3 → internal rank = 8-3 = 5 → index 5*8+2 = 42
	// f6 chess = file 5, chess rank 6 → internal rank = 2 → index 2*8+5 = 21
	// Simpler: use engine squares from a custom FEN with knight that can jump onto a portal.

	// Place white knight on b1 (index 1), portal on a3 (16) and h3 (23) — knight can't reach those.
	// Use king on e1 and portal adjacent: king e1 (4), portal e2 (12), exit e4 (28).
	fen := "4k3/8/8/8/8/8/4P3/4K3 w - - 0 1"
	// Actually use a piece that can move to a portal square.
	// White rook on a1, portal on a4 (index 24), exit on h4 (index 31).
	fen = "4k3/8/8/8/7p/8/8/R3K3 w - - 0 1"
	// a4 = file0 rank (8-4)=4 → 4*8+0 = 32
	// h4 = file7 rank4 → 4*8+7 = 39
	// black pawn on h4 so teleport captures.
	a4 := 4*8 + 0
	h4 := 4*8 + 7

	eng, err := chess.NewEngine(chess.GameConfig{
		FEN:         fen,
		VariantType: "wormhole",
		CustomData: map[string]interface{}{
			"wormholePairs": []interface{}{
				[]interface{}{a4, h4},
			},
		},
	})
	require.NoError(t, err)

	from := 7*8 + 0 // a1
	var teleport *chess.Move
	for _, m := range eng.GetLegalMoves() {
		if m.From == from && m.To == h4 && m.VariantMoveType == "teleport" {
			cp := m
			teleport = &cp
			break
		}
	}
	require.NotNil(t, teleport, "expected rook a1→portal a4→exit h4 teleport capture")
	require.True(t, teleport.Capture)

	portal, ok := variants.TeleportPortalOf(*teleport)
	require.True(t, ok)
	require.Equal(t, a4, portal)

	// Client-style: click the portal entrance, not the exit.
	_, err = eng.PerformMove(chess.Move{From: from, To: a4, Piece: 'R'})
	require.NoError(t, err, "clicking portal entrance should resolve to teleport")
	pos := eng.GetPosition()
	require.True(t, pos.Pieces['R'].HasBit(h4), "rook should be on exit h4")
	require.False(t, pos.Pieces['R'].HasBit(a4), "rook should have left portal a4")
	require.False(t, pos.Pieces['p'].HasBit(h4), "black pawn should be captured")

	// Pair on cooldown: black should not be able to teleport via this pair.
	// Put black king move elsewhere; verify no white teleport available next... after black moves.
	// Black to move — find any legal move and play it, then white should still see cooldown (1 left after black).
	blackMoves := eng.GetLegalMoves()
	require.NotEmpty(t, blackMoves)
	_, err = eng.PerformMove(blackMoves[0])
	require.NoError(t, err)

	// White to move again — pair still cooling (1→0 after this white non-teleport, but during generation cooldown is still >0 from after black's decrement...)
	// After teleport: cooldown=2. After black OnAfterMove: cooldown=1. White GetLegalMoves: pair still closed.
	for _, m := range eng.GetLegalMoves() {
		require.NotEqual(t, "teleport", m.VariantMoveType, "pair should still be on cooldown")
	}
}
