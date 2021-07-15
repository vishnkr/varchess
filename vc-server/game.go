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
	CustomMovePatterns []MovePatterns
}
//customMove patterns: {'stringpiece' : {jump:[[]], slide:[[]]}}
type MovePatterns struct {
	PieceName string `json:"piece"`
	JumpPattern [][]int `json:"jumpPattern"`
	SlidePattern [][]int `json:"slidePattern"`
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
				//fmt.Println("before custom",board.Tiles[i][j].Piece.CustomPiece)
				if (board.Tiles[i][j].Piece.CustomPiece!=nil){
						//fmt.Println("found custom",board.Tiles[i][j].Piece.CustomPiece)
						piece = board.Tiles[i][j].Piece.CustomPiece.PieceName
				} else {
					piece = board.Tiles[i][j].Piece.String()
				}
				if (board.Tiles[i][j].Piece.Color == White){
					piece = strings.ToUpper(piece)
				}
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
	if (board.willCauseDiscoveredCheck(piece,move)){ return false,"discovery check"}
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
		case Custom:
			return isCustomMoveValid(piece,board,move)

	}
	return false,"something's wrong"
}

// performMove modifies board state to account for position changes after a valid move has been made
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
/*func (board *Board) findKingSquareId(color Color) int{
	for 
}*/

func (board *Board) willCauseDiscoveredCheck(piece *Piece, move *Move) bool {
	//get opponent attacking squares, if king is present in that then it causes discovery check
	//squareId : bool hashmap 
	var copyBoard *Board = deepCopyBoard(board)
	copyBoard.performMove(piece,move)
	DisplayBoardState(copyBoard)
	//copyBoard.isKingUnderCheck(piece,)
	return false
}

// getSquaresAttackedBy returns a collection of squares attacked by pieces of the given color
func (board *Board) getSquaresAttackedBy(color Color) map[int]bool{
	var attackedSquares = make(map[int]bool)
	for rowIndex,row:= range board.Tiles{
		for colIndex,tile:= range row{
			if (!tile.IsEmpty && tile.Piece.Color==color){
				board.genPieceAttacks(&tile.Piece,rowIndex,colIndex,attackedSquares)
			}
		}
	}
	return attackedSquares
}

func (board *Board) isKingUnderCheck(piece *Piece,squareId int) bool{
	opponentCol := piece.getOpponentColor()
	attackedSquares := board.getSquaresAttackedBy(opponentCol)
	if _,underAttack := attackedSquares[squareId]; underAttack{
		return true
	} 
	return false
}

func (board *Board)isSquareInBoardRange(row int, col int) bool{
	return row>=0 && col>=0 && row<board.Rows && col<board.Cols
}	

// genPieceMoves returns a list of square positions (using its unique Square Id) that the given piece attacks
// this includes the position of the opponent king since it would be useful to detect checks/discovered checks
func (board *Board) genPieceAttacks(piece *Piece,srcRow int, srcCol int,attackedSquares map[int]bool){
	switch piece.Type{
	case Bishop:
		board.genBishopMoves(piece,srcRow,srcCol,attackedSquares)
	case Rook:
		board.genRookMoves(piece,srcRow,srcCol,attackedSquares)
	case Queen:
		board.genRookMoves(piece,srcRow,srcCol,attackedSquares)
		board.genBishopMoves(piece,srcRow,srcCol,attackedSquares)
	case King:
		for row:=-1;row<=1;row+=1{
			for col:=-1;col<=1;col+=1{
				if (!(row==0)&& !(col==0)){
					continue
				}
				if (board.isSquareInBoardRange(srcRow+row,srcCol+col) && board.IsEmpty(srcRow+row,srcCol+col)){
					attackedSquares[board.Tiles[srcRow+row][srcCol+col].Id] = true
				}
			}
		}
	case Knight:
		jumpSquares:= [][]int{{-2,-1},{-2,1},{-1,2},{-1,-2},{1,-2},{1,2},{2,-1},{2,1}}
		for _, pair:= range jumpSquares{
			targetRow:= srcRow+pair[0]
			targetCol:= srcCol+pair[1]
			if (board.isSquareInBoardRange(targetRow,targetCol) && board.IsEmpty(targetRow,targetCol)){
				attackedSquares[board.Tiles[targetRow][targetCol].Id] = true
			}
		}
	case Pawn:
		var rowOffset int
		if (piece.Color==White){
			rowOffset=-1
		} else {rowOffset=1}
		if (board.isSquareInBoardRange(srcRow+rowOffset,srcCol-1) && board.IsEmpty(srcRow+rowOffset,srcCol-1)){
			attackedSquares[board.Tiles[srcRow+rowOffset][srcCol-1].Id] = true
		}
		if (board.isSquareInBoardRange(srcRow+rowOffset,srcCol+1) && board.IsEmpty(srcRow+rowOffset,srcCol+1)){
			attackedSquares[board.Tiles[srcRow+rowOffset][srcCol+1].Id] = true
		}

	}
}

func(board *Board) genBishopMoves(piece *Piece,srcRow int, srcCol int,attackedSquares map[int]bool){
	diagonals:= [][]int{{1,1,board.Rows-1},{1,-1,board.Rows-1},{-1,1,0},{-1,-1,0}}
	for _,value:= range diagonals{
		xOffset,yOffset,endRow := value[0],value[1],value[2]
		pathLength:=Abs(srcRow-endRow)
		for i := 1; i <= pathLength; i++{
			x := srcRow + i*xOffset;
			y := srcCol + i*yOffset;
			if (board.isSquareInBoardRange(x,y) && board.IsEmpty(x,y) ){
				attackedSquares[board.Tiles[x][y].Id] = true
			} else if (board.isSquareInBoardRange(x,y) && board.Tiles[x][y].Piece.Type == King && board.Tiles[x][y].Piece.Color == piece.getOpponentColor()){
				attackedSquares[board.Tiles[x][y].Id] = true
				break
			}else { break}
		}
	}
}

func(board *Board) genRookMoves(piece *Piece,srcRow int,srcCol int,attackedSquares map[int]bool){
	directions:= []int{-1,1}
	//check horizontal
	for _,xOffset := range directions{
		for i:=srcCol+xOffset;i>=0 && i<board.Cols;i+=xOffset{
			if (board.IsEmpty(srcRow,i)){
				attackedSquares[board.Tiles[srcRow][i].Id] = true
			} else if (board.Tiles[srcRow][i].Piece.Type == King && board.Tiles[srcRow][i].Piece.Color == piece.getOpponentColor()){
				attackedSquares[board.Tiles[srcRow][i].Id] = true
				break
			} else{break}
		}
	}
	//check vertical
	for _,yOffset := range directions{
		for i :=srcRow+yOffset;i>=0 && i<board.Rows;i+=yOffset{
			if (board.IsEmpty(i,srcCol)){
				attackedSquares[board.Tiles[i][srcCol].Id] = true
			} else if (board.Tiles[i][srcCol].Piece.Type == King && board.Tiles[i][srcCol].Piece.Color == piece.getOpponentColor()){
				attackedSquares[board.Tiles[i][srcCol].Id] = true
				break
			} else {break}
		}
	}
}

// isKingMoveValid handles all logic to determine if a king move is valid
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

// isRoomMoveValid determines if a rook move is valid
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

// isBishopMoveValid determines if a Bishop move is valid
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

// isBishopMoveValid determines if a bishop move is valid
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
					return true,"pawn promoted"
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

func isCustomMoveValid(piece *Piece, board *Board, move *Move) (bool,string){
	//Check jump moves followed by slide moves
	var jumpPattern,slidePattern [][]int
	//find pattern for the piece, will change this to a hashmap later
	for _,movePatterns := range board.CustomMovePatterns{
		fmt.Println("mp",movePatterns,piece,piece.CustomPiece)
		if (movePatterns.PieceName == strings.ToLower(piece.CustomPiece.PieceName)){
			jumpPattern = movePatterns.JumpPattern
			slidePattern = movePatterns.SlidePattern
			break
		}
	}
	
	var rowDiff,colDiff int = move.SrcRow-move.DestRow, move.SrcCol-move.DestCol
	fmt.Println("diffs",rowDiff,colDiff,jumpPattern)
	if (len(jumpPattern)!=0){
		for _, pair := range jumpPattern{
			fmt.Println("pair",pair)
			if (pair[0]==rowDiff && pair[1]==colDiff){
				return true,"valid custom jump move"
			}
		}
	}

	for _, direction := range slidePattern{
		var rowOffset,colOffset int = direction[0], direction[1]
		var tempRow,tempCol int = move.SrcRow+rowOffset, move.SrcCol+colOffset
		for (tempRow>=0 && tempCol >=0 && tempRow<board.Rows && tempCol<board.Cols){
			if (tempRow==move.DestRow && tempCol==move.DestCol){
				return true,"valid custom slide move"
			} else if (!board.IsEmpty(tempRow,tempCol)){ break } else{
				tempRow+=rowOffset
				tempCol+=colOffset
			}
		}
	}
	
	return false,"no"
}