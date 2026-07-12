package variants

import "vc-server/chess"

// PoisonedPawnRules: a designated pawn is poisoned. Any piece that captures it
// becomes "poisoned" and is removed from the board after a configurable number of moves.
// TODO: implement poison tracking and delayed removal.
type PoisonedPawnRules struct {
	PoisonedSquare int            // square of the poisoned pawn
	PoisonedPieces map[int]int    // square → moves remaining before removal
	PoisonDuration int            // how many moves before a poisoned piece dies
}

func (r *PoisonedPawnRules) ExtraMoves(_ *chess.Position, _ chess.Color) []chess.Move { return nil }
func (r *PoisonedPawnRules) FilterMoves(moves []chess.Move, _ *chess.Position) []chess.Move {
	return moves
}
func (r *PoisonedPawnRules) OnBeforeMove(_ chess.Move, _ *chess.Position) error { return nil }
func (r *PoisonedPawnRules) OnAfterMove(_ chess.Move, _ *chess.Position, _ *chess.MoveResult) {}
func (r *PoisonedPawnRules) CheckTermination(_ *chess.Position) *chess.EngineResult { return nil }
func (r *PoisonedPawnRules) VariantState() interface{}                               { return nil }
func (r *PoisonedPawnRules) LoadVariantState(_ interface{}) error                    { return nil }
func (r *PoisonedPawnRules) Clone() chess.VariantRules {
	cloned := *r
	if r.PoisonedPieces != nil {
		cloned.PoisonedPieces = make(map[int]int, len(r.PoisonedPieces))
		for k, v := range r.PoisonedPieces {
			cloned.PoisonedPieces[k] = v
		}
	}
	return &cloned
}
