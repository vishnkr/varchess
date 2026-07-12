package variants

import "vc-server/chess"

// SpyChessRules: each player places one hidden "spy" piece in the opponent's territory
// at game start. The spy can be revealed and transformed back at any time as a free move.
// TODO: implement spy placement, reveal moves, and hidden piece server-side masking.
type SpyChessRules struct {
	SpySquares map[chess.Color]int // color → square of that color's spy in opponent territory (-1 if revealed)
}

func (r *SpyChessRules) ExtraMoves(_ *chess.Position, _ chess.Color) []chess.Move { return nil }
func (r *SpyChessRules) FilterMoves(moves []chess.Move, _ *chess.Position) []chess.Move {
	return moves
}
func (r *SpyChessRules) OnBeforeMove(_ chess.Move, _ *chess.Position) error            { return nil }
func (r *SpyChessRules) OnAfterMove(_ chess.Move, _ *chess.Position, _ *chess.MoveResult) {}
func (r *SpyChessRules) CheckTermination(_ *chess.Position) *chess.EngineResult        { return nil }
func (r *SpyChessRules) VariantState() interface{}                                      { return nil }
func (r *SpyChessRules) LoadVariantState(_ interface{}) error                           { return nil }
func (r *SpyChessRules) Clone() chess.VariantRules {
	cloned := *r
	if r.SpySquares != nil {
		cloned.SpySquares = make(map[chess.Color]int, len(r.SpySquares))
		for k, v := range r.SpySquares {
			cloned.SpySquares[k] = v
		}
	}
	return &cloned
}
