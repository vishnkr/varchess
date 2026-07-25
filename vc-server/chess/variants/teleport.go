package variants

import (
	"encoding/json"
	"unicode"

	"vc-server/chess"
)

const teleportVariantMove = "teleport"

// TeleportRules implements Wormhole / Teleportation Chess.
//
// Rules:
//   - Landing on an open portal teleports the piece to the paired exit.
//   - Exit must be empty or hold an enemy (capture). Own piece on exit → illegal.
//   - Kings may use portals.
//   - After a pair is used, it cools down for one full turn (2 half-moves).
//
// Legal moves use To = exit (final square) for correct check detection, and stash
// the portal entrance in AdditionalData so clients can click the wormhole square.
type TeleportRules struct {
	WormholePairs [][2]int
	// PairCooldown[i] > 0 means pair i is closed for that many remaining half-moves.
	PairCooldown []int
}

func NewTeleportRules(pairs [][2]int) *TeleportRules {
	return &TeleportRules{
		WormholePairs: pairs,
		PairCooldown:   make([]int, len(pairs)),
	}
}

func (r *TeleportRules) ExtraMoves(_ *chess.Position, _ chess.Color) []chess.Move { return nil }

// FilterMoves rewrites destinations that land on an open portal to the exit square.
func (r *TeleportRules) FilterMoves(moves []chess.Move, pos *chess.Position) []chess.Move {
	if len(r.WormholePairs) == 0 {
		return moves
	}
	out := make([]chess.Move, 0, len(moves))
	for _, move := range moves {
		exit, pairIdx, ok := r.openExitOf(move.To)
		if !ok {
			out = append(out, move)
			continue
		}
		occupant := posPieceAt(pos, exit)
		if occupant != 0 && sameColor(occupant, move.Piece) {
			// Own piece blocks exit — teleport illegal.
			continue
		}
		rewritten := move
		portal := move.To
		rewritten.To = exit
		rewritten.Capture = occupant != 0 || move.Capture
		if rewritten.Capture {
			rewritten.ClassicMoveType = chess.CaptureMove
		}
		rewritten.VariantMoveType = teleportVariantMove
		rewritten.AdditionalData = map[string]interface{}{
			"portal": portal,
			"pair":   pairIdx,
		}
		out = append(out, rewritten)
	}
	return out
}

func (r *TeleportRules) OnBeforeMove(_ chess.Move, _ *chess.Position) error { return nil }

func (r *TeleportRules) OnAfterMove(move chess.Move, _ *chess.Position, _ *chess.MoveResult) {
	for i := range r.PairCooldown {
		if r.PairCooldown[i] > 0 {
			r.PairCooldown[i]--
		}
	}
	if move.VariantMoveType != teleportVariantMove {
		return
	}
	pairIdx := teleportPairIndex(move.AdditionalData)
	if pairIdx < 0 || pairIdx >= len(r.PairCooldown) {
		return
	}
	// One full turn = 2 half-moves; set after decrement so this ply stays closed.
	r.PairCooldown[pairIdx] = 2
}

func (r *TeleportRules) CheckTermination(_ *chess.Position) *chess.EngineResult { return nil }

type teleportState struct {
	WormholePairs [][2]int `json:"wormholePairs"`
	PairCooldown  []int    `json:"pairCooldown"`
}

func (r *TeleportRules) VariantState() interface{} {
	return teleportState{
		WormholePairs: r.WormholePairs,
		PairCooldown:  append([]int(nil), r.PairCooldown...),
	}
}

func (r *TeleportRules) LoadVariantState(s interface{}) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	var state teleportState
	if err := json.Unmarshal(b, &state); err != nil {
		return err
	}
	r.WormholePairs = state.WormholePairs
	r.PairCooldown = state.PairCooldown
	if len(r.PairCooldown) < len(r.WormholePairs) {
		padded := make([]int, len(r.WormholePairs))
		copy(padded, r.PairCooldown)
		r.PairCooldown = padded
	}
	return nil
}

func (r *TeleportRules) Clone() chess.VariantRules {
	cloned := *r
	if r.WormholePairs != nil {
		cloned.WormholePairs = make([][2]int, len(r.WormholePairs))
		copy(cloned.WormholePairs, r.WormholePairs)
	}
	if r.PairCooldown != nil {
		cloned.PairCooldown = make([]int, len(r.PairCooldown))
		copy(cloned.PairCooldown, r.PairCooldown)
	}
	return &cloned
}

func (r *TeleportRules) openExitOf(sq int) (exit int, pairIdx int, ok bool) {
	for i, p := range r.WormholePairs {
		if i < len(r.PairCooldown) && r.PairCooldown[i] > 0 {
			continue
		}
		if p[0] == sq {
			return p[1], i, true
		}
		if p[1] == sq {
			return p[0], i, true
		}
	}
	return 0, -1, false
}

func posPieceAt(pos *chess.Position, sq int) rune {
	for p, bb := range pos.Pieces {
		if bb != nil && bb.HasBit(sq) {
			return p
		}
	}
	return 0
}

func sameColor(a, b rune) bool {
	return unicode.IsUpper(a) == unicode.IsUpper(b)
}

func additionalInt(v interface{}, fallback int) int {
	switch n := v.(type) {
	case int:
		return n
	case int32:
		return int(n)
	case int64:
		return int(n)
	case float64:
		return int(n)
	case float32:
		return int(n)
	default:
		return fallback
	}
}

func teleportPairIndex(data interface{}) int {
	switch v := data.(type) {
	case map[string]interface{}:
		return additionalInt(v["pair"], -1)
	default:
		return additionalInt(data, -1)
	}
}

// TeleportPortalOf returns the portal entrance square for a teleport legal move.
func TeleportPortalOf(m chess.Move) (int, bool) {
	if m.VariantMoveType != teleportVariantMove {
		return 0, false
	}
	switch v := m.AdditionalData.(type) {
	case map[string]interface{}:
		p := additionalInt(v["portal"], -1)
		if p >= 0 {
			return p, true
		}
	}
	return 0, false
}

// parseWormholePairs reads customData["wormholePairs"] as [][2]int (JSON numbers).
func parseWormholePairs(customData map[string]interface{}) [][2]int {
	if customData == nil {
		return nil
	}
	raw, ok := customData["wormholePairs"]
	if !ok || raw == nil {
		return nil
	}
	// Prefer a round-trip through JSON so we accept [][]int, []interface{}, etc.
	b, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	var asInts [][]int
	if err := json.Unmarshal(b, &asInts); err == nil {
		var pairs [][2]int
		for _, p := range asInts {
			if len(p) < 2 || p[0] == p[1] {
				continue
			}
			pairs = append(pairs, [2]int{p[0], p[1]})
		}
		return pairs
	}
	return nil
}
