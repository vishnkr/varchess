package variants

import "vc-server/chess"

// StandardRules is the default no-op VariantRules for checkmate chess.
// All hooks fall through to the core engine's standard behaviour.
type StandardRules struct{}

func (r *StandardRules) ExtraMoves(_ *chess.Position, _ chess.Color) []chess.Move { return nil }
func (r *StandardRules) FilterMoves(moves []chess.Move, _ *chess.Position) []chess.Move {
	return moves
}
func (r *StandardRules) OnBeforeMove(_ chess.Move, _ *chess.Position) error            { return nil }
func (r *StandardRules) OnAfterMove(_ chess.Move, _ *chess.Position, _ *chess.MoveResult) {}
func (r *StandardRules) CheckTermination(_ *chess.Position) *chess.EngineResult        { return nil }
func (r *StandardRules) VariantState() interface{}                                      { return nil }
func (r *StandardRules) LoadVariantState(_ interface{}) error                           { return nil }
func (r *StandardRules) Clone() chess.VariantRules                                      { return &StandardRules{} }
