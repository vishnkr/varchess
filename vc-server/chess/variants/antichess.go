package variants

import (
	"unicode"

	"vc-server/chess"
)

// AntichessRules implements the "lose all pieces" / giveaway variant.
// Rules:
//   - Captures are mandatory when available.
//   - Win condition: be the first to lose all pieces (or have no legal moves).
type AntichessRules struct{}

func (r *AntichessRules) ExtraMoves(_ *chess.Position, _ chess.Color) []chess.Move { return nil }

// FilterMoves enforces mandatory captures: if any capture exists, only captures are legal.
func (r *AntichessRules) FilterMoves(moves []chess.Move, _ *chess.Position) []chess.Move {
	var captures []chess.Move
	for _, m := range moves {
		if m.Capture {
			captures = append(captures, m)
		}
	}
	if len(captures) > 0 {
		return captures
	}
	return moves
}

func (r *AntichessRules) OnBeforeMove(_ chess.Move, _ *chess.Position) error { return nil }
func (r *AntichessRules) OnAfterMove(_ chess.Move, _ *chess.Position, _ *chess.MoveResult) {}

// CheckTermination returns a win if either side has lost all pieces, or a draw if
// only same-coloured bishops remain (dead position).
func (r *AntichessRules) CheckTermination(pos *chess.Position) *chess.EngineResult {
	whitePieces, blackPieces := 0, 0
	for p, bb := range pos.Pieces {
		count := len(bb.GetSetBits())
		if unicode.IsUpper(p) {
			whitePieces += count
		} else {
			blackPieces += count
		}
	}
	if whitePieces == 0 {
		return &chess.EngineResult{Winner: "white", Reason: "no-pieces"}
	}
	if blackPieces == 0 {
		return &chess.EngineResult{Winner: "black", Reason: "no-pieces"}
	}

	// Dead-position: only same-coloured bishops remain for both sides.
	if r.isBishopDraw(pos) {
		return &chess.EngineResult{Winner: "draw", Reason: "bishop-draw"}
	}
	return nil
}

func (r *AntichessRules) isBishopDraw(pos *chess.Position) bool {
	lbd := pos.LargestDimension
	darkWhite, lightWhite, darkBlack, lightBlack := 0, 0, 0, 0
	for p, bb := range pos.Pieces {
		if unicode.ToLower(p) != 'b' {
			if len(bb.GetSetBits()) > 0 {
				return false // non-bishop piece present
			}
			continue
		}
		for _, sq := range bb.GetSetBits() {
			f, rank := sq%lbd, sq/lbd
			isDark := (f+rank)%2 == 0
			if unicode.IsUpper(p) {
				if isDark {
					darkWhite++
				} else {
					lightWhite++
				}
			} else {
				if isDark {
					darkBlack++
				} else {
					lightBlack++
				}
			}
		}
	}
	// Draw only if bishops of different colours cannot capture each other.
	if (darkWhite > 0 && darkBlack > 0) || (lightWhite > 0 && lightBlack > 0) {
		return false
	}
	return true
}

func (r *AntichessRules) VariantState() interface{}              { return nil }
func (r *AntichessRules) LoadVariantState(_ interface{}) error   { return nil }
func (r *AntichessRules) Clone() chess.VariantRules              { return &AntichessRules{} }
