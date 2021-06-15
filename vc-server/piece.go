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

func (piece Piece) isBackwardPawnMove(move *Move) bool{
	if (piece.Color == Black && move.DestRow<move.SrcRow) || (piece.Color == White && move.DestRow>move.SrcRow) {
			return true
	}
	return false
}

var typeToRuneMap = map[Type]rune{Pawn:'p', Knight:'n', Bishop:'b', Rook:'r', Queen:'q', King:'k'}
var strToTypeMap = map[string]Type{"p":Pawn,"n":Knight,"b":Bishop, "r":Rook, "q": Queen, "k":King}

func (board *Board) isPieceStartPosValid(piece Piece, row int, col int) bool{
	return  board.getPieceColor(row,col) == piece.Color && board.tiles[row][col].Piece.Type == piece.Type
}

func (board *Board) getPieceColor(row int,col int) Color{
	if(!board.tiles[row][col].IsEmpty){
		return board.tiles[row][col].Piece.Color
	}
	return EmptyTile
}	

