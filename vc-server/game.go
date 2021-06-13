package main

import (
	//"github.com/gorilla/websocket"
	"fmt"
	
)

type Game struct {
	p1       *Client
	p2       *Client
	fen      string
	moveList []string
	pgn      string
	result   string
}

type Board struct {
	tiles [][]rune
	rows  int
	cols  int
	turn rune
}

type Move struct{
	SrcRow int `json:"srcRow"`
	SrcCol int `json:"srcCol"`
	DestRow int `json:"destRow"`
	DestCol int `json:"destCol"`
	Promote Type `json:"promote,omitempty"`
}

func (board *Board) IsEmpty(row int,col int) bool{
	return board.tiles[row][col]==0
}

func DisplayBoardState(board *Board){
	var piece string
	for i:=0;i<board.rows;i++{
		for j:=0;j<board.cols;j++{
			if (board.tiles[i][j] !=0){
				piece =string(board.tiles[i][j])
			} else{piece ="-"}
			fmt.Print(piece," ")
		}
		fmt.Print("\n")
	}
}


func (piece Piece) isValidMove(board *Board,move *Move) (bool,string){
	//check if theres blocks to dest
	//opens up a discovery check?
	
	if(board.getPieceColor(move.SrcRow,move.SrcCol) == board.getPieceColor(move.DestRow,move.DestCol)){
		fmt.Println(board.getPieceColor(move.SrcRow,move.SrcCol),board.getPieceColor(move.DestRow,move.DestCol))
		return false,"same color piece at dest"
	}
	switch piece.Type{
		case Rook: 
			//horizontal or vertical block
			if(move.SrcCol==move.DestCol && move.SrcRow!=move.DestRow){
				//vertical move
				start,end := Min(move.SrcRow,move.DestRow),Max(move.SrcRow,move.DestRow)
				for i:=start+1;i<end;i++{
					if(!board.IsEmpty(i,move.SrcCol)){
						return false,("path blocked by "+string(board.tiles[i][move.SrcCol]))
					}
				}
				return true,"valid rook move"
			} else if(move.SrcRow==move.DestRow && move.SrcCol!=move.DestCol){
				//horizontal move
				start,end := Min(move.SrcCol,move.DestCol),Max(move.SrcCol,move.DestCol)
				for i:=start+1;i<end;i++{
					if(!board.IsEmpty(move.SrcRow,i)){
						return false,("path blocked by "+string(board.tiles[i][move.SrcCol]))
					}
				}
				return true,"valid rook move"
			} else { return false,"invalid rook move"}

		case Bishop:
			pathLength:= Abs(move.SrcRow - move.DestRow)
			if pathLength!= Abs(move.SrcCol - move.DestCol){
				return false, "not diagonal"
			}
			for i := 1; i < pathLength; i++{
				x := move.SrcRow + i;
				y := move.SrcCol + i;
				if (!board.IsEmpty(x, y)){
					// Obstacle found before reaching target: the move is invalid
					return false,"obstacle in bishop path" 
				} 
			}
			return true,"valid"

		case Knight:
			if (Abs((move.SrcRow-move.DestRow) * (move.SrcCol-move.DestCol)) ==2){
				return true,"valid knight"
			} else{ return false, "invalid knight"}
		
		case Pawn:
			//not considering en passant yet
			var doubleMoveStartRank int
			if piece.Color=='b' {
				doubleMoveStartRank = 1
			} else {
				doubleMoveStartRank = board.rows - 2
			}
			if(move.DestCol== move.SrcCol && move.SrcRow!=move.DestRow){
				if (Abs(move.SrcRow-move.DestRow)==2 && move.SrcRow==doubleMoveStartRank){
					if ( (piece.Color=='b' && board.IsEmpty(move.SrcRow+1,move.SrcCol)) || (piece.Color=='w' && board.IsEmpty(move.SrcRow-1,move.SrcCol))) {
						fmt.Println(board.tiles[move.SrcRow-1][move.SrcCol],board.IsEmpty(move.SrcRow-1,move.SrcCol))
						return true,"double pawn move allowed"
					} else{ return false,"double move blocked"}
				} else if (Abs(move.SrcRow-move.DestRow)==1){
					if (board.IsEmpty(move.SrcRow-1,move.SrcCol)){
						return true,"valid single pawn move"
					} else { return false, "dest blocked"}
				} 
				return true,"more logic to handle here"
			}

	}
	return true,"good to go"
}