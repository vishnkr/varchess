package main

import (
	//"unicode"
	//"fmt"
)

type Type uint8
const (
	Custom Type = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
	Empty
)

type Color uint8
const (
	EmptyTile Color = iota
	White 
	Black
)	
 
type Piece struct{
	Type Type
	Color Color
	CustomPiece *CustomPiece
}

type CustomPiece struct{
	PieceName string
}

type SqColor uint8 
const (
	Dark SqColor = iota
	Light
)

type Square struct{
	SqColor SqColor
	Id int
	Piece Piece
	IsEmpty bool
}

type KingPiece struct{
	HasMoved bool
	InCheck bool
	Position []int
}

func (p Piece) String() string {
	switch p.Type {
	case King:
		return "k"
	case Queen:
		return "q"
	case Rook:
		return "r"
	case Bishop:
		return "b"
	case Knight:
		return "n"
	case Pawn:
		return "p"
	}
	return ""
}

func promotableTo(p Piece) bool {
	switch p.Type {
	case Queen, Rook, Bishop, Knight:
		return true
	}
	return false
}

func (piece Piece) isBackwardPawnMove(move *Move) bool{
	if (piece.Color == Black && move.DestRow<move.SrcRow) || (piece.Color == White && move.DestRow>move.SrcRow) {
			return true
	}
	return false
}

//var typeToRuneMap = map[Type]rune{Pawn:'p', Knight:'n', Bishop:'b', Rook:'r', Queen:'q', King:'k'}
var typeToStrMap = map[Type]string{Pawn:"p", Knight:"n", Bishop:"b", Rook:"r", Queen:"q", King:"k"}
var strToTypeMap = map[string]Type{"p":Pawn,"n":Knight,"b":Bishop, "r":Rook, "q": Queen, "k":King}

func (board *Board) isPieceStartPosValid(piece *Piece, row int, col int) bool{
	//fmt.Println("at row",row,"col",col,"color",piece.Color,"type",piece.Type)
	var validType bool
	if (board.Tiles[row][col].Piece.Type == Custom){
		validType = board.Tiles[row][col].Piece.CustomPiece.PieceName == piece.CustomPiece.PieceName
	} else {
		validType = board.Tiles[row][col].Piece.Type == piece.Type
	}
	return  board.getPieceColor(row,col) == piece.Color && validType
}

func (board *Board) getPieceColor(row int,col int) Color{
	if (!board.Tiles[row][col].IsEmpty){
		return board.Tiles[row][col].Piece.Color
	}
	return EmptyTile
}	

func getOpponentColor(color Color) Color{
	if (color==White){
		return Black
	}else{
		return White
	}
}

