// Package testutil provides helpers for constructing chess engines and positions
// in tests. It blank-imports chess/variants so the factory is always registered.
package testutil

import (
	"testing"

	"vc-server/chess"
	_ "vc-server/chess/variants"
)

const StartFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// MustNewEngine creates a chess.Engine from cfg or immediately fails the test.
func MustNewEngine(t *testing.T, cfg chess.GameConfig) chess.Engine {
	t.Helper()
	eng, err := chess.NewEngine(cfg)
	if err != nil {
		t.Fatalf("NewEngine(%+v): %v", cfg, err)
	}
	return eng
}

// MustStandard returns a StandardEngine at the starting position.
func MustStandard(t *testing.T) chess.Engine {
	t.Helper()
	return MustNewEngine(t, chess.GameConfig{
		VariantType: "checkmate",
		FEN:         StartFEN,
	})
}

// MustSquare converts an algebraic square name to an internal index or fails
// the test.
func MustSquare(t *testing.T, pos *chess.Position, sq string) int {
	t.Helper()
	idx, err := pos.Square(sq)
	if err != nil {
		t.Fatalf("Square(%q): %v", sq, err)
	}
	return idx
}

// MustPerformMove performs a move from a legal-move list matched by algebraic
// square string (e.g. "e2e4") and fails the test if not found or illegal.
func MustPerformMove(t *testing.T, eng chess.Engine, algebraic string) chess.MoveResult {
	t.Helper()
	pos := eng.GetPosition()
	from, to, promo, err := chess.ParseAlgebraicMove(algebraic, pos)
	if err != nil {
		t.Fatalf("ParseAlgebraicMove(%q): %v", algebraic, err)
	}
	move, ok := chess.MatchMove(from, to, promo, eng.GetLegalMoves())
	if !ok {
		t.Fatalf("move %q not found in legal moves", algebraic)
	}
	result, err := eng.PerformMove(move)
	if err != nil {
		t.Fatalf("PerformMove(%q): %v", algebraic, err)
	}
	return result
}

// PlaySequence performs a sequence of algebraic moves and fails on any error.
func PlaySequence(t *testing.T, eng chess.Engine, moves ...string) {
	t.Helper()
	for _, m := range moves {
		MustPerformMove(t, eng, m)
	}
}
