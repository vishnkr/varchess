package variants

import "vc-server/chess"

// FootballChessRules: a "ball" piece must be moved to the opponent's goal squares to win.
// Pieces can "pass" the ball or "shoot" it. No king in this variant.
// TODO: implement ball tracking, pass/shoot move generation, and goal detection.
type FootballChessRules struct {
	BallSquare  int    // current square of the ball (-1 if not on board)
	GoalSquares [2]int // [0]=white goal, [1]=black goal
}

func (r *FootballChessRules) ExtraMoves(_ *chess.Position, _ chess.Color) []chess.Move { return nil }
func (r *FootballChessRules) FilterMoves(moves []chess.Move, _ *chess.Position) []chess.Move {
	return moves
}
func (r *FootballChessRules) OnBeforeMove(_ chess.Move, _ *chess.Position) error            { return nil }
func (r *FootballChessRules) OnAfterMove(_ chess.Move, _ *chess.Position, _ *chess.MoveResult) {}
func (r *FootballChessRules) CheckTermination(_ *chess.Position) *chess.EngineResult        { return nil }
func (r *FootballChessRules) VariantState() interface{}                                      { return nil }
func (r *FootballChessRules) LoadVariantState(_ interface{}) error                           { return nil }
func (r *FootballChessRules) Clone() chess.VariantRules                                      { cloned := *r; return &cloned }
