package game

import (
	"fmt"
	"testing"
)

func TestNewPosition(t *testing.T) {
	gameConfig := GameConfig{
		VariantType: Custom,
		gameConfigPosition: gameConfigPosition{
			Dimensions: dimensions{Ranks: 8, Files: 8},
			Fen:       "4r3/8/8/8/8/8/8/R1B1K2R w KQ - 0 1", //"rnbkqbnr/pppppppp/q7/8/8/8/1PPPPPPP/RNBKQBNR w Qkq - 0 1",
		},
		Objective: Objective{ObjectiveType: Checkmate},
	}

	position, err := newPosition(gameConfig)
	if err != nil {
		t.Errorf("Error creating new position: %v", err)
	}

	if position.castlingRights.whiteKingSide {
		t.Errorf("Expected white king-side castling rights to be false")
	}
	v, err := newVariant(gameConfig)
	moves := v.getPseudoLegalMoves(ColorWhite, false)
	fmt.Print(moves)
	legal := v.GetLegalMoves()
	fmt.Print(legal)
	// Add more assertions as needed
}
