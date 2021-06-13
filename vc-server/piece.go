package main

import (
	"unicode"
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
 
type Piece struct{
	Type Type
	Color rune //w or b
}

func (board *Board) getPieceColor(row int,col int) rune{
	if(board.tiles[row][col]!=0){
		if(unicode.IsUpper(board.tiles[row][col])){
			return rune('w')
		} else{
			return rune('b')
		}
	}
	return 0
}	

