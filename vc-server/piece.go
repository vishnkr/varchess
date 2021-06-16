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
}

type SqColor uint8 
const (
	Dark SqColor = iota
	Light
)

type Square struct{
	SqColor SqColor
	Piece Piece
	IsEmpty bool
}

type KingPiece struct{
	HasMoved bool
	InCheck bool
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

func (p Piece) promotableTo() bool {
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

var typeToRuneMap = map[Type]rune{Pawn:'p', Knight:'n', Bishop:'b', Rook:'r', Queen:'q', King:'k'}
var strToTypeMap = map[string]Type{"p":Pawn,"n":Knight,"b":Bishop, "r":Rook, "q": Queen, "k":King}

func (board *Board) isPieceStartPosValid(piece Piece, row int, col int) bool{
	return  board.getPieceColor(row,col) == piece.Color && board.Tiles[row][col].Piece.Type == piece.Type
}

func (board *Board) getPieceColor(row int,col int) Color{
	if(!board.Tiles[row][col].IsEmpty){
		return board.Tiles[row][col].Piece.Color
	}
	return EmptyTile
}	

