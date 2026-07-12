package variants

import "vc-server/chess"

// TeleportRules implements Teleportation / Wormhole Chess.
// Pieces that enter a wormhole square are teleported to the paired square.
// TODO: implement wormhole teleport move generation.
type TeleportRules struct {
	WormholePairs [][2]int // pairs of wormhole entrance/exit squares
}

func (r *TeleportRules) ExtraMoves(_ *chess.Position, _ chess.Color) []chess.Move { return nil }
func (r *TeleportRules) FilterMoves(moves []chess.Move, _ *chess.Position) []chess.Move {
	return moves
}
func (r *TeleportRules) OnBeforeMove(_ chess.Move, _ *chess.Position) error            { return nil }
func (r *TeleportRules) OnAfterMove(_ chess.Move, _ *chess.Position, _ *chess.MoveResult) {}
func (r *TeleportRules) CheckTermination(_ *chess.Position) *chess.EngineResult        { return nil }
func (r *TeleportRules) VariantState() interface{}                                      { return nil }
func (r *TeleportRules) LoadVariantState(_ interface{}) error                           { return nil }
func (r *TeleportRules) Clone() chess.VariantRules {
	cloned := *r
	if r.WormholePairs != nil {
		cloned.WormholePairs = make([][2]int, len(r.WormholePairs))
		copy(cloned.WormholePairs, r.WormholePairs)
	}
	return &cloned
}
