package variants

import "vc-server/chess"

// HiddenQueenRules: each player has a hidden queen placed at game start that
// is invisible to the opponent until it moves or captures.
// TODO: implement hidden queen tracking, reveal mechanics, and server-side masking.
type HiddenQueenRules struct {
	HiddenSquares map[chess.Color]int // color → square of hidden queen (-1 if revealed)
}

func (r *HiddenQueenRules) ExtraMoves(_ *chess.Position, _ chess.Color) []chess.Move { return nil }
func (r *HiddenQueenRules) FilterMoves(moves []chess.Move, _ *chess.Position) []chess.Move {
	return moves
}
func (r *HiddenQueenRules) OnBeforeMove(_ chess.Move, _ *chess.Position) error            { return nil }
func (r *HiddenQueenRules) OnAfterMove(_ chess.Move, _ *chess.Position, _ *chess.MoveResult) {}
func (r *HiddenQueenRules) CheckTermination(_ *chess.Position) *chess.EngineResult        { return nil }
func (r *HiddenQueenRules) VariantState() interface{}                                      { return nil }
func (r *HiddenQueenRules) LoadVariantState(_ interface{}) error                           { return nil }
func (r *HiddenQueenRules) Clone() chess.VariantRules {
	cloned := *r
	if r.HiddenSquares != nil {
		cloned.HiddenSquares = make(map[chess.Color]int, len(r.HiddenSquares))
		for k, v := range r.HiddenSquares {
			cloned.HiddenSquares[k] = v
		}
	}
	return &cloned
}
