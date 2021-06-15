package main

import (
	"unicode"
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
//temporary - will be removed once tiles stores piece structs instead of runes
var typeToRuneMap = map[Type]rune{Pawn:'p', Knight:'n', Bishop:'b', Rook:'r', Queen:'q', King:'k'}

func (board *Board) isPieceStartPosValid(piece Piece, row int, col int) bool{
	//DisplayBoardState(board)
	//fmt.Println("bro",board.getPieceColor(row,col),piece.Color,unicode.ToLower(board.tiles[row][col]),typeToRuneMap[piece.Type])
	return  board.getPieceColor(row,col) == piece.Color && unicode.ToLower(board.tiles[row][col]) == typeToRuneMap[piece.Type]
}

func (board *Board) getPieceColor(row int,col int) Color{
	if(board.tiles[row][col]!=0){
		if(unicode.IsUpper(board.tiles[row][col])){
			return White
		} else{
			return Black
		}
	}
	return EmptyTile
}	

