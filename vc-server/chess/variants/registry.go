// Package variants provides concrete VariantRules implementations for all
// supported chess variants. Import this package (blank or named) to register
// the variant factory with the chess engine:
//
//	import _ "vc-server/chess/variants"
package variants

import (
	"fmt"

	"vc-server/chess"
)

func init() {
	chess.RegisterVariantFactory(NewVariantRules)
}

// NewVariantRules resolves a variant type string to its VariantRules implementation.
// customData carries variant-specific configuration (e.g. targetChecks for NCheck).
func NewVariantRules(variantType string, customData map[string]interface{}) (chess.VariantRules, error) {
	switch variantType {
	case "", "checkmate":
		return &StandardRules{}, nil

	case "antichess":
		return &AntichessRules{}, nil

	case "ncheck":
		n := 3 // default
		if customData != nil {
			if v, ok := customData["targetChecks"]; ok {
				switch val := v.(type) {
				case int:
					n = val
				case float64:
					n = int(val)
				}
			}
		}
		return NewNCheckRules(n), nil

	case "archerchess":
		return &ArcherRules{}, nil

	case "wormhole", "teleport":
		pairs := parseWormholePairs(customData)
		return NewTeleportRules(pairs), nil

	case "poisonedpawn":
		return &PoisonedPawnRules{
			PoisonDuration: 3,
			PoisonedPieces: make(map[int]int),
		}, nil

	case "hiddenqueen":
		return &HiddenQueenRules{
			HiddenSquares: map[chess.Color]int{chess.White: -1, chess.Black: -1},
		}, nil

	case "footballchess":
		return &FootballChessRules{BallSquare: -1}, nil

	case "spychess":
		return &SpyChessRules{
			SpySquares: map[chess.Color]int{chess.White: -1, chess.Black: -1},
		}, nil

	default:
		return nil, fmt.Errorf("unknown variant type: %q", variantType)
	}
}
