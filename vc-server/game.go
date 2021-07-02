package main

import (
	"strings"
	"fmt"
	
)

type Game struct {
	P1       *Client
	P2       *Client
	Board    *Board
	MoveList []string
	Pgn      string
	Result   string
	Turn string
}

type Board struct {
	Tiles [][]Square
	Rows  int
	Cols  int
	BlackKing KingPiece
	WhiteKing KingPiece
}


type Move struct{
	SrcRow int `json:"srcRow"`
	SrcCol int `json:"srcCol"`
	DestRow int `json:"destRow"`
	DestCol int `json:"destCol"`
	Promote Type `json:"promote,omitempty"`
	Castle bool `json:"castle,omitempty"`
	PieceType string `json:"piece"`
	RoomId string `json:"roomId"`
	Color string `json:"color"`
}

func (board *Board) IsEmpty(row int,col int) bool{
	return board.Tiles[row][col].IsEmpty
}

//visual helper function
func DisplayBoardState(board *Board){
	var piece string
	for i:=0;i<board.Rows;i++{
		for j:=0;j<board.Cols;j++{
			if (!board.Tiles[i][j].IsEmpty){
				if (board.Tiles[i][j].Piece.Color == White){
					piece = strings.ToUpper(board.Tiles[i][j].Piece.String())
				} else { piece = board.Tiles[i][j].Piece.String()} 
			} else{piece = "-"}
			fmt.Print(piece,board.Tiles[i][j].Id," ")
		}
		fmt.Print("\n")
	}
}


//isValidMove: checks if a move made by given piece is valid, return string is just used for debugging purposes
func (board *Board) isValidMove(piece *Piece,move *Move) (bool,string){
	if (!board.isPieceStartPosValid(piece,move.SrcRow,move.SrcCol)){ return false,"start pos not valid for given piece" }
	//check if same piece color exists at destination
	if (board.getPieceColor(move.SrcRow,move.SrcCol) == board.getPieceColor(move.DestRow,move.DestCol)){
		fmt.Println(board.getPieceColor(move.SrcRow,move.SrcCol),board.getPieceColor(move.DestRow,move.DestCol))
		return false,"same color piece at dest"
	}
	//if (board.willCauseDiscoveredCheck(piece,move)){ return false,"discovery check"}
	switch piece.Type{
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
		case King:
			return isKingMoveValid(piece,board,move)

	}
	return false,"something's wrong"
}

func (board *Board) performMove(piece *Piece,move *Move){
	
	board.Tiles[move.DestRow][move.DestCol].IsEmpty = false
	board.Tiles[move.DestRow][move.DestCol].Piece = board.Tiles[move.SrcRow][move.SrcCol].Piece
	board.Tiles[move.SrcRow][move.SrcCol].IsEmpty = true
	board.Tiles[move.SrcRow][move.SrcCol].Piece.Type = Empty
	board.Tiles[move.SrcRow][move.SrcCol].Piece.Color = EmptyTile
	if (piece.Type==King){
		if (piece.Color==Black){
			board.BlackKing.HasMoved = true
		} else {board.WhiteKing.HasMoved=true}
	}
	if (move.Castle){
		var oldRookPos, newRookPos int
		if (move.DestCol<move.SrcCol){
			oldRookPos = 0
			newRookPos = move.SrcCol-1
		} else{
			oldRookPos = board.Cols-1
			newRookPos = move.SrcCol+1
		}
		board.Tiles[move.SrcRow][newRookPos].IsEmpty = false
		board.Tiles[move.SrcRow][newRookPos].Piece = Piece{Type:Rook,Color:piece.Color}
		board.Tiles[move.SrcRow][oldRookPos].IsEmpty = true
		board.Tiles[move.SrcRow][oldRookPos].Piece.Type = Empty
		board.Tiles[move.SrcRow][oldRookPos].Piece.Color = EmptyTile
	}
}

func (board *Board) willCauseDiscoveredCheck(piece *Piece, move *Move) bool {
	//get opponent attacking squares, if king is present in that then it causes discovery check
	//squareId : bool hashmap 
	var copyBoard *Board = deepCopyBoard(board)
	copyBoard.performMove(piece,move)
	DisplayBoardState(copyBoard)
	attackedSquares := copyBoard.getSquaresAttackedBy(piece.Color)
	fmt.Print(attackedSquares)
	return false
}

func (board *Board) getSquaresAttackedBy(color Color) map[int]bool{
	var attackedSquares = make(map[int]bool)
	for rowIndex,row:= range board.Tiles{
		for colIndex,tile:= range row{
			if (!tile.IsEmpty){
				for _,id:= range board.genPieceMoves(&tile.Piece,rowIndex,colIndex){
					attackedSquares[id]=true
				}
			}
		}
	}
	return attackedSquares
}


func (board *Board) genPieceMoves(piece *Piece,srcRow int, srcCol int) []int{
	//includes king
	attackedSquares:=[]int{}
	switch piece.Type{
	case Bishop:
		diagonals:= [][]int{{1,1,7},{1,-1,7},{-1,1,0},{-1,-1,0}}
		for _,value:= range diagonals{
			xOffset,yOffset,endRow := value[0],value[1],value[2]
			fmt.Println(xOffset,yOffset,endRow)
			pathLength:=Abs(srcRow-endRow)
			for i := 1; i <= pathLength; i++{
				x := srcRow + i*xOffset;
				y := srcCol + i*yOffset;
				if (x>=0 && y>=0 && x<8 && y<8 && (board.Tiles[x][y].IsEmpty || board.Tiles[x][y].Piece.Type == King)){
					attackedSquares = append(attackedSquares,board.Tiles[x][y].Id)
				} else{ break}
			}
		}
	case Rook:

	}
	return attackedSquares
}



func isKingMoveValid(piece *Piece, board *Board, move *Move) (bool,string){
	if (move.Castle && !board.hasKingMoved(piece.Color)){
		if (move.SrcRow==move.DestRow){
			var tile int 
			if (move.SrcCol+2==move.DestCol && board.Tiles[move.SrcRow][board.Cols-1].Piece.Type==Rook){ //castle to the right
				tile= move.SrcCol+1
				for (tile<board.Cols-1){
					if (!board.Tiles[move.SrcRow][tile].IsEmpty){return false,"castle path blocked"}
					tile+=1
				}
				return true, "valid castle"
			} else if (move.SrcCol-2==move.DestCol && board.Tiles[move.SrcRow][0].Piece.Type==Rook){
				tile= move.SrcCol-1
				for tile>0{
					if (!board.Tiles[move.SrcRow][tile].IsEmpty){return false,"castle path blocked"}
					tile-=1
				}
				return true, "valid castle"
			} else {return false,"invalid castle dest"}
		}
	} else if (move.Castle && board.hasKingMoved(piece.Color)) {return false,"king has already moved"}
	return true,"valid king move"
}

func (board* Board) hasKingMoved(color Color) bool{
	if(color==White && board.WhiteKing.HasMoved || (color==Black && board.BlackKing.HasMoved)){
		return true
	} else {return false}
}

func isRookMoveValid(piece *Piece, board *Board, move *Move) (bool,string){
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

func isBishopMoveValid(piece *Piece, board *Board, move *Move) (bool,string){
	pathLength:= Abs(move.SrcRow - move.DestRow)
	if pathLength!= Abs(move.SrcCol - move.DestCol){
		return false, "not diagonal"
	}
	var xOffset,yOffset int
	if (move.SrcRow < move.DestRow && move.SrcCol<move.DestCol){
		xOffset = 1
		yOffset = 1
	} else if (move.SrcRow > move.DestRow && move.SrcCol<move.DestCol){
		xOffset = -1
		yOffset = 1
	} else if (move.SrcRow < move.DestRow && move.SrcCol > move.DestCol){
		xOffset = 1
		yOffset = -1
	} else {
		xOffset = -1
		yOffset = -1
	}
	for i := 1; i < pathLength; i++{
		x := move.SrcRow + i*xOffset;
		y := move.SrcCol + i*yOffset;
		if (!board.IsEmpty(x, y)){
			// Obstacle found before reaching target: the move is invalid
			return false,"obstacle in bishop path" 
		} 
	}
	
	
	return true,"valid"
}

func isPawnMoveValid(piece *Piece, board *Board, move *Move) (bool,string){
	//not considering en passant yet
	var doubleMoveStartRank,rowOffset,promoteDestRow int
	if piece.Color==Black {
		doubleMoveStartRank = 1
		rowOffset = 1
		promoteDestRow = board.Rows-1
	} else {
		doubleMoveStartRank = board.Rows - 2
		rowOffset = -1
		promoteDestRow = 0
	}
	if (move.DestCol== move.SrcCol && move.SrcRow!=move.DestRow){
		if (Abs(move.SrcRow-move.DestRow)==2 && move.SrcRow==doubleMoveStartRank){
			if ( (piece.Color==Black && board.IsEmpty(move.SrcRow+1,move.SrcCol)) || (piece.Color==White && board.IsEmpty(move.SrcRow-1,move.SrcCol))) {
				return true,"double pawn move allowed"
			} else{ return false,"double move blocked"}
		} else if (Abs(move.SrcRow-move.DestRow)==1 && !piece.isBackwardPawnMove(move)){
			if (board.IsEmpty(move.DestRow,move.DestCol)){
				if(move.Promote!=0 && move.DestRow==promoteDestRow){
					return true,("pawn promoted to "+ string(typeToRuneMap[move.Promote]))
				}
				return true,"valid single pawn move"
			} else { 
				
				return false, "dest blocked"}
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

func changeTurn(turn string) string{
	if turn =="w"{
		return "b"
	} else { return "w"}
}