package main

import (
	"testing"
	//"fmt"
)

func TestPawnCapture(t *testing.T) {
	//white pawn capture
	exp := true
	board:= ConvertFENtoBoard("8/p6p/1Q6/8/3B1p2/4P3/1P6/8 w - - 0 1")
	piece:=&Piece{Type:Pawn,Color:White}
	move:=&Move{
		SrcRow:5,
		SrcCol:4,
		DestRow:4,
		DestCol:5,
	}
    result,reason := piece.isValidMove(board,move)
    if exp != result {
       t.Errorf("result 1 incorrect, got: %v, %s, want: %v.",result,reason,exp )
    }
	//black pawn capture
	exp = true
	piece =&Piece{Type:Pawn,Color:Black}
	move =&Move{
		SrcRow:1,
		SrcCol:0,
		DestRow:2,
		DestCol:1,
	}
    result,reason = piece.isValidMove(board,move)
    if exp != result {
       t.Errorf("result 2 incorrect, got: %v, %s, want: %v.",result,reason,exp )
    }
	//fmt.Println(result,reason)
}

func TestDoublePawnMove(t *testing.T){
	exp := false
	board:= ConvertFENtoBoard("8/p6p/1Q6/8/3B1p2/4P3/1P6/8 w - - 0 1")
	piece:=&Piece{Type:Pawn,Color:White}
	move:=&Move{
		SrcRow:1,
		SrcCol:7,
		DestRow:3,
		DestCol:7,
	}
    result,reason := piece.isValidMove(board,move)
    if exp != result {
       t.Errorf("result 1 incorrect, got: %v, %s, want: %v.",result,reason,exp )
    }

	exp = true
	piece =&Piece{Type:Pawn,Color:White}
	move =&Move{
		SrcRow:6,
		SrcCol:1,
		DestRow:4,
		DestCol:1,
	}
    result,reason = piece.isValidMove(board,move)
    if exp != result {
       t.Errorf("result 2 incorrect, got: %v, %s, want: %v.",result,reason,exp )
    }
}

func TestPawnPromotion(t *testing.T){
	exp := false
	board:= ConvertFENtoBoard("8/p2P3p/1Q6/8/3B1p2/4P3/1P6/8 w - - 0 1")
	piece:=&Piece{Type:Pawn,Color:Black}
	move:=&Move{
		SrcRow:1,
		SrcCol:0,
		DestRow:0,
		DestCol:0,
		Promote: Queen,
	}
    result,reason := piece.isValidMove(board,move)
    if exp != result {
       t.Errorf("result incorrect, got: %v, %s, want: %v.",result,reason,exp )
    }
}

func TestQueenDiagonal(t *testing.T){
	exp := true
	board:= ConvertFENtoBoard("8/8/8/8/3Q4/8/8/8 w - - 0 1")
	piece:=&Piece{Type:Queen,Color:White}
	move:=&Move{
		SrcRow:4,
		SrcCol:3,
		DestRow:1,
		DestCol:0,
		Promote: Queen,
	}
    result,reason := piece.isValidMove(board,move)
    if exp != result {
       t.Errorf("result incorrect, got: %v, %s, want: %v.",result,reason,exp )
    }
}

func TestQueenHorizontal(t *testing.T){
	exp := true
	board:= ConvertFENtoBoard("8/8/8/8/3Q4/8/8/8 w - - 0 1")
	piece:=&Piece{Type:Queen,Color:White}
	move:=&Move{
		SrcRow:4,
		SrcCol:3,
		DestRow:4,
		DestCol:7,
		Promote: Queen,
	}
    result,reason := piece.isValidMove(board,move)
    if exp != result {
       t.Errorf("result incorrect, got: %v, %s, want: %v.",result,reason,exp )
    }
}