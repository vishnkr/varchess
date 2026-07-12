package variants

import "vc-server/chess"

// ArcherRules implements Archer Chess: designated pieces can attack without moving.
// The archer squares and their attack ranges are configured via CustomData.
// TODO: implement ranged attack move generation.
type ArcherRules struct {
	ArcherSquares map[int]bool // squares whose pieces can attack without moving
}

func (r *ArcherRules) ExtraMoves(_ *chess.Position, _ chess.Color) []chess.Move { return nil }
func (r *ArcherRules) FilterMoves(moves []chess.Move, _ *chess.Position) []chess.Move {
	return moves
}
func (r *ArcherRules) OnBeforeMove(_ chess.Move, _ *chess.Position) error            { return nil }
func (r *ArcherRules) OnAfterMove(_ chess.Move, _ *chess.Position, _ *chess.MoveResult) {}
func (r *ArcherRules) CheckTermination(_ *chess.Position) *chess.EngineResult        { return nil }
func (r *ArcherRules) VariantState() interface{}                                      { return nil }
func (r *ArcherRules) LoadVariantState(_ interface{}) error                           { return nil }
func (r *ArcherRules) Clone() chess.VariantRules {
	cloned := *r
	if r.ArcherSquares != nil {
		cloned.ArcherSquares = make(map[int]bool, len(r.ArcherSquares))
		for k, v := range r.ArcherSquares {
			cloned.ArcherSquares[k] = v
		}
	}
	return &cloned
}
