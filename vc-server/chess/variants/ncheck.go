package variants

import (
	"encoding/json"
	"fmt"
	"unicode"

	"vc-server/chess"
)

// NCheckRules implements the N-Check variant (e.g. 3-Check).
// Win condition: give the opponent N checks.
type NCheckRules struct {
	TargetN     int
	ChecksGiven map[chess.Color]int // how many checks each color has given
}

func NewNCheckRules(n int) *NCheckRules {
	return &NCheckRules{
		TargetN:     n,
		ChecksGiven: map[chess.Color]int{chess.White: 0, chess.Black: 0},
	}
}

func (r *NCheckRules) ExtraMoves(_ *chess.Position, _ chess.Color) []chess.Move { return nil }
func (r *NCheckRules) FilterMoves(moves []chess.Move, _ *chess.Position) []chess.Move {
	return moves
}
func (r *NCheckRules) OnBeforeMove(_ chess.Move, _ *chess.Position) error { return nil }

// OnAfterMove increments the check counter for the moving side if the opponent is in check.
func (r *NCheckRules) OnAfterMove(move chess.Move, pos *chess.Position, result *chess.MoveResult) {
	if !result.IsCheck {
		return
	}
	// The side that just moved is the one who gave check.
	// After MakeMove, WhiteToMove is flipped, so the mover is the opponent of WhiteToMove.
	var moverColor chess.Color
	if pos.WhiteToMove {
		moverColor = chess.Black // black just moved
	} else {
		moverColor = chess.White // white just moved
	}
	r.ChecksGiven[moverColor]++
	if r.ChecksGiven[moverColor] >= r.TargetN {
		winner := "white"
		if moverColor == chess.Black {
			winner = "black"
		}
		result.IsGameOver = true
		result.Result = &chess.EngineResult{
			Winner: winner,
			Reason: fmt.Sprintf("%d-check", r.TargetN),
		}
	}
}

// CheckTermination is called after OnAfterMove; we handle it there instead.
func (r *NCheckRules) CheckTermination(_ *chess.Position) *chess.EngineResult { return nil }

// nCheckState is used for JSON serialisation of variant state.
type nCheckState struct {
	TargetN     int            `json:"targetN"`
	ChecksGiven map[string]int `json:"checksGiven"`
}

func (r *NCheckRules) VariantState() interface{} {
	return nCheckState{
		TargetN: r.TargetN,
		ChecksGiven: map[string]int{
			"white": r.ChecksGiven[chess.White],
			"black": r.ChecksGiven[chess.Black],
		},
	}
}

func (r *NCheckRules) LoadVariantState(s interface{}) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	var state nCheckState
	if err := json.Unmarshal(b, &state); err != nil {
		return err
	}
	r.TargetN = state.TargetN
	r.ChecksGiven = map[chess.Color]int{
		chess.White: state.ChecksGiven["white"],
		chess.Black: state.ChecksGiven["black"],
	}
	return nil
}

func (r *NCheckRules) Clone() chess.VariantRules {
	cloned := *r
	cloned.ChecksGiven = map[chess.Color]int{
		chess.White: r.ChecksGiven[chess.White],
		chess.Black: r.ChecksGiven[chess.Black],
	}
	return &cloned
}

// countPiecesOfColor counts total pieces for a color on the board.
func countPiecesOfColor(pos *chess.Position, color chess.Color) int {
	total := 0
	for p, bb := range pos.Pieces {
		isWhite := unicode.IsUpper(p)
		if (color == chess.White && isWhite) || (color == chess.Black && !isWhite) {
			total += len(bb.GetSetBits())
		}
	}
	return total
}
