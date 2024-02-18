package chesscore

import (
	"fmt"
	"runtime"
	"testing"
)

func MeasureMemoryUsage(t *testing.T, f func()) {
	var mStart, mEnd runtime.MemStats
	runtime.ReadMemStats(&mStart)
	f()
	runtime.ReadMemStats(&mEnd)
	memoryUsed := mEnd.TotalAlloc - mStart.TotalAlloc
	fmt.Printf("Memory used: %v bytes\n", memoryUsed)
}

func TestCastlingRights(t *testing.T) {
	gameConfig := GameConfig{
		VariantType: Checkmate,
		Dimensions: Dimensions{Ranks: 8, Files: 8},
		Fen:       "n3r3/8/7p/8/8/8/8/R1B1K2R w Kq - 0 1", //"rnbkqbnr/pppppppp/q7/8/8/8/1PPPPPPP/RNBKQBNR w Qkq - 0 1",
	}

	position, err := newPosition(gameConfig)
	if err != nil {
		t.Errorf("Error creating new position: %v", err)
	}

	if position.castlingRights.whiteQueenSide || !position.castlingRights.whiteKingSide{
		t.Errorf("Incorrect white king castling rights")
	}
	if position.castlingRights.blackKingSide || !position.castlingRights.blackQueenSide{
		t.Errorf("Incorrect black king castling rights")
	}
}

func TestCheckmatePieceMoves(t *testing.T) {
	
		gameConfig := GameConfig{
			VariantType: Checkmate,
			Dimensions: Dimensions{Ranks: 8, Files: 8},
			Fen:       "n3r3/7a/p7/8/8/8/3k4/3R4 b - - 0 1", //"rnbkqbnr/pppppppp/q7/8/8/8/1PPPPPPP/RNBKQBNR w Qkq - 0 1",
			PieceProps: map[string]PieceProperties{
				"a":{
					Name: "testPiece",
					JumpOffsets: []moveOffset{{x: -2,y: 2},{x:-1,y:1},{x:3,y:-1},{x:1,y:-2}},
				},
			},
		}

		v,err :=newVariant(gameConfig)
		if err!=nil{
			fmt.Println(err)
		}
		moves:= v.GetLegalMoves()
		fmt.Println(moves)
}

//rnbqkbnr/ppp1ppp1/7p/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq 19 0 1
func TestAntichessEP(t *testing.T){
	gc := GameConfig{
		VariantType: Checkmate,
		Dimensions: Dimensions{Ranks: 8, Files: 8},
		Fen:       "rnbqkbnr/ppp1ppp1/7p/3pP3/8/8/7R/8 w KQkq 19 0 1",
	}
	v, _ := newVariant(gc)
	moves := v.getPseudoLegalMoves(ColorWhite, false)
	fmt.Println(moves)
	legal := v.GetLegalMoves()
	fmt.Println(legal)

}

func TestMovegenPrint(t *testing.T){
	gameConfig := GameConfig{
		VariantType: Checkmate,
		Dimensions: Dimensions{Ranks: 8, Files: 8},
		Fen:       "n3r3/2B5/7p/8/8/8/8/R4K1R b Ke - 0 1",
	}
	v, _ := newVariant(gameConfig)
	moves := v.getPseudoLegalMoves(ColorBlack, false)
	fmt.Println(moves)
	legal := v.GetLegalMoves()
	fmt.Println(legal)
}


