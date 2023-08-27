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
			Fen: "rnbkqbnr/pppppppp/q7/8/8/8/1PPPPPPP/RNBKQBNR w Qkq - 0 1",
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
    v,err:= newVariant(gameConfig)
    moves:= v.getPseudoLegalMoves(ColorWhite,false)
    fmt.Print(moves)
    // Add more assertions as needed
}