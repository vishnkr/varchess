package main

import (
	//"github.com/gorilla/websocket"
	"unicode"
	"fmt"
	
)

type Game struct {
	p1       *Client
	p2       *Client
	board    *Board
	moveList []string
	pgn      string
	result   string
}

type Board struct {
	tiles [][]Square
	rows  int
	cols  int
	turn rune
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

type Move struct{
	SrcRow int `json:"srcRow"`
	SrcCol int `json:"srcCol"`
	DestRow int `json:"destRow"`
	DestCol int `json:"destCol"`
	Promote Type `json:"promote,omitempty"`
}

func (board *Board) IsEmpty(row int,col int) bool{
	return board.tiles[row][col].IsEmpty
}

func DisplayBoardState(board *Board){
	var piece string
	//fmt.Println(board.rows,board.cols)
	for i:=0;i<board.rows;i++{
		for j:=0;j<board.cols;j++{
			if (!board.tiles[i][j].IsEmpty){
				//fmt.Println("go",board.tiles[i][j].Piece,typeToRuneMap[board.tiles[i][j].Piece.Type],typeToRuneMap)
				if (board.tiles[i][j].Piece.Color == White){
					piece = string(unicode.ToUpper(typeToRuneMap[board.tiles[i][j].Piece.Type]))
				} else { piece = string(typeToRuneMap[board.tiles[i][j].Piece.Type]) }
			} else{piece = "-"}
			fmt.Print(piece," ")
		}
		fmt.Print("\n")
	}
}


//split this function up later once complete logic is done
func (piece Piece) isValidMove(board *Board,move *Move) (bool,string){
	if (!board.isPieceStartPosValid(piece,move.SrcRow,move.SrcCol)){ return false,"start pos not valid for given piece" }
	//check if same piece color exists at destination
	if (board.getPieceColor(move.SrcRow,move.SrcCol) == board.getPieceColor(move.DestRow,move.DestCol)){
		fmt.Println(board.getPieceColor(move.SrcRow,move.SrcCol),board.getPieceColor(move.DestRow,move.DestCol))
		return false,"same color piece at dest"
	}

	switch piece.Type{
		//doesn't check if move opens up a discovery check yet
		case Rook: 
			return isRookMoveValid(piece,board,move)
		case Bishop:
			return isBishopMoveValid(piece,board,move)
		case Knight:
			if (Abs((move.SrcRow-move.DestRow) * (move.SrcCol-move.DestCol)) ==2){
				return true,"valid knight"
			} else{ return false, "invalid knight"}
		case Pawn:
			return isPawnMoveValid(piece,board,move)
		case Queen:
			rookCheck,res:= isRookMoveValid(piece,board,move)
			if (!rookCheck){
				return isBishopMoveValid(piece,board,move)
			}
			return rookCheck,res

	}
	return true,"good to go"
}

func isRookMoveValid(piece Piece, board *Board, move *Move) (bool,string){
	//horizontal or vertical block
	if (move.SrcCol==move.DestCol && move.SrcRow!=move.DestRow){
		//vertical move
		start,end := Min(move.SrcRow,move.DestRow),Max(move.SrcRow,move.DestRow)
		for i:=start+1;i<end;i++{
			if(!board.IsEmpty(i,move.SrcCol)){
				return false,("rook path blocked")
			}
		}
		return true,"valid rook move"
	} else if(move.SrcRow==move.DestRow && move.SrcCol!=move.DestCol){
		//horizontal move
		start,end := Min(move.SrcCol,move.DestCol),Max(move.SrcCol,move.DestCol)
		for i:=start+1;i<end;i++{
			if(!board.IsEmpty(move.SrcRow,i)){
				return false,("rook path blocked")
			}
		}
		return true,"valid rook move"
	} else { return false,"invalid rook move"}
}

func isBishopMoveValid(piece Piece, board *Board, move *Move) (bool,string){
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
}

func isPawnMoveValid(piece Piece, board *Board, move *Move) (bool,string){
	//not considering en passant yet
	var doubleMoveStartRank,rowOffset,promoteDestRow int
	if piece.Color==Black {
		doubleMoveStartRank = 1
		rowOffset = 1
		promoteDestRow = board.rows-1
	} else {
		doubleMoveStartRank = board.rows - 2
		rowOffset = -1
		promoteDestRow = 0
	}
	if (move.DestCol== move.SrcCol && move.SrcRow!=move.DestRow){
		if (Abs(move.SrcRow-move.DestRow)==2 && move.SrcRow==doubleMoveStartRank){
			if ( (piece.Color==Black && board.IsEmpty(move.SrcRow+1,move.SrcCol)) || (piece.Color==White && board.IsEmpty(move.SrcRow-1,move.SrcCol))) {
				return true,"double pawn move allowed"
			} else{ return false,"double move blocked"}
		} else if (Abs(move.SrcRow-move.DestRow)==1 && !piece.isBackwardPawnMove(move)){
			if (board.IsEmpty(move.SrcRow-1,move.SrcCol)){
				if(move.Promote!=0 && move.DestRow==promoteDestRow){
					return true,("pawn promoted to"+string(move.Promote))
				}
				return true,"valid single pawn move"
			} else { return false, "dest blocked"}
		} 
		
	} else if(move.DestCol!= move.SrcCol && move.SrcRow!=move.DestRow){
		//check if col row +-1 logic
			if(Abs(move.SrcCol-move.DestCol)==1 && move.DestRow==move.SrcRow+rowOffset && !board.IsEmpty(move.DestRow,move.DestCol)){
				return true,"valid pawn capture"
			}
			return false,"invalid pawn capture"
	}
	return false, "not a valid pawn move"
}