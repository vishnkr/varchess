package game

import (
	"fmt"
	"strings"
)

type Game struct {
	Board    *Board
	MoveList []string
	Pgn      string
	Result   string
	Turn     string
}

type Board struct {
	Tiles              [][]Square
	Rows               int
	Cols               int
	BlackKing          KingPiece
	WhiteKing          KingPiece
	CustomMovePatterns []MovePatterns
}

//customMove patterns: {'stringpiece' : {jump:[[]], slide:[[]]}}
type MovePatterns struct {
	PieceName    string  `json:"piece"`
	JumpPattern  [][]int `json:"jumpPattern"`
	SlidePattern [][]int `json:"slidePattern"`
}

type Move struct {
	SrcRow    int    `json:"srcRow"`
	SrcCol    int    `json:"srcCol"`
	DestRow   int    `json:"destRow"`
	DestCol   int    `json:"destCol"`
	Promote   Type   `json:"promote,omitempty"`
	Castle    bool   `json:"castle,omitempty"`
	PieceType string `json:"piece"`
	RoomId    string `json:"roomId"`
	Color     string `json:"color"`
	Captured  Piece
}

func (board *Board) IsEmpty(row int, col int) bool {
	return board.isSquareInBoardRange(row, col) && board.Tiles[row][col].IsEmpty
}

//visual helper function
func DisplayBoardState(board *Board) {
	var piece string
	for i := 0; i < board.Rows; i++ {
		for j := 0; j < board.Cols; j++ {
			if !board.Tiles[i][j].IsEmpty {
				if board.Tiles[i][j].Piece.CustomPiece != nil {
					piece = board.Tiles[i][j].Piece.CustomPiece.PieceName
				} else {
					piece = board.Tiles[i][j].Piece.String()
				}
				if board.Tiles[i][j].Piece.Color == White {
					piece = strings.ToUpper(piece)
				}
			} else {
				piece = "-"
			}
			fmt.Print(piece, " ")
		}
		fmt.Print("\n")
	}
}

func (board *Board) isSameColorPieceAtDest(color Color, destRow int, destCol int) bool {
	if board.isSquareInBoardRange(destRow, destCol) && !board.IsEmpty(destRow, destCol) && color == board.getPieceColor(destRow, destCol) {
		return true
	}
	return false
}

//IsValidMove: checks if a move made by given piece is valid, return string is just used for debugging purposes
func (board *Board) IsValidMove(piece Piece, move *Move) (bool, string) {
	if !board.isPieceStartPosValid(piece, move.SrcRow, move.SrcCol) {
		return false, "start pos not valid for given piece"
	}
	//check if same piece color exists at destination
	if board.isSameColorPieceAtDest(piece.Color, move.DestRow, move.DestCol) {
		return false, "same color piece at dest"
	}
	if board.willCauseCheck(piece, move) {
		return false, "under check/discovered check"
	}
	switch piece.Type {
	case Rook:
		return isRookMoveValid(piece, board, move)
	case Bishop:
		return isBishopMoveValid(piece, board, move)
	case Knight:
		if Abs((move.SrcRow-move.DestRow)*(move.SrcCol-move.DestCol)) == 2 {
			return true, "valid knight"
		} else {
			return false, "invalid knight"
		}
	case Pawn:
		return isPawnMoveValid(piece, board, move)
	case Queen:
		rookCheck, res := isRookMoveValid(piece, board, move)
		if !rookCheck {
			return isBishopMoveValid(piece, board, move)
		}
		return rookCheck, res
	case King:
		return isKingMoveValid(piece, board, move)
	case Custom:
		return isCustomMoveValid(piece, board, move)

	}
	return false, "something's wrong"
}

func (board *Board) isDestOccupied(color Color, destRow int, destCol int) bool {
	return !board.Tiles[destRow][destCol].IsEmpty && board.Tiles[destRow][destCol].Piece.Color == GetOpponentColor(color)

}

// PerformMove: modify board state after a move has been made
func (board *Board) PerformMove(piece Piece, move Move) {

	if board.isDestOccupied(piece.Color, move.DestRow, move.DestCol) {
		move.Captured = board.Tiles[move.DestRow][move.DestCol].Piece
	}
	board.Tiles[move.DestRow][move.DestCol].IsEmpty = false
	board.Tiles[move.DestRow][move.DestCol].Piece = board.Tiles[move.SrcRow][move.SrcCol].Piece
	board.Tiles[move.SrcRow][move.SrcCol].IsEmpty = true
	if piece.Type == King {
		if piece.Color == Black {
			board.BlackKing.MoveCount += 1
			board.BlackKing.Position = []int{move.DestRow, move.DestCol}
		} else {
			board.WhiteKing.Position = []int{move.DestRow, move.DestCol}
			board.WhiteKing.MoveCount += 1
		}
	}
	if move.Castle {
		var oldRookPos, newRookPos int
		if move.DestCol < move.SrcCol {
			oldRookPos = 0
			newRookPos = move.SrcCol - 1
		} else {
			oldRookPos = board.Cols - 1
			newRookPos = move.SrcCol + 1
		}
		board.Tiles[move.SrcRow][newRookPos].IsEmpty = false
		board.Tiles[move.SrcRow][newRookPos].Piece = Piece{Type: Rook, Color: piece.Color}
		board.Tiles[move.SrcRow][oldRookPos].IsEmpty = true
	}
}

//umakeMove: revert board to original state before perform move
func (board *Board) unmakeMove(piece Piece, move Move) {
	board.Tiles[move.SrcRow][move.SrcCol].IsEmpty = false
	if (move.Captured != Piece{}) {
		board.Tiles[move.DestRow][move.DestCol].Piece = move.Captured
	}
	if piece.Type == King {
		if piece.Color == Black {
			board.BlackKing.MoveCount -= 1
			board.BlackKing.Position = []int{move.SrcRow, move.SrcCol}
		} else {
			board.WhiteKing.Position = []int{move.SrcRow, move.SrcCol}
			board.WhiteKing.MoveCount -= 1
		}
	}
	board.Tiles[move.SrcRow][move.SrcCol].Piece = board.Tiles[move.DestRow][move.DestCol].Piece
	board.Tiles[move.DestRow][move.DestCol].IsEmpty = true
}

func (board *Board) IsGameOver(color Color) (bool, string) {
	check := board.IsKingUnderCheck(color)
	moves := board.GetAllValidMoves(color)
	if len(moves) == 0 {
		if check {
			return true, GetOpponentColor(color).String()
		} else {
			return true, "draw"
		}
	}
	return false, "nah"
}

func (board *Board) willCauseCheck(piece Piece, move *Move) bool {
	copyBoard := deepCopyBoard(board)
	copyBoard.PerformMove(piece, *move)
	postCheck := copyBoard.IsKingUnderCheck(piece.Color)
	return postCheck
}

//IsKingUnderCheck: returns whether player's king is under check or not
func (board *Board) IsKingUnderCheck(color Color) bool {
	var kingPos []int
	if color == White {
		kingPos = board.WhiteKing.Position
	} else {
		kingPos = board.BlackKing.Position
	}
	oppMoves := board.getAllPseudoLegalMoves(GetOpponentColor(color))
	for move := range oppMoves {
		if kingPos[0] == move.DestRow && kingPos[1] == move.DestCol {
			return true
		}
	}
	return false
}

func (board *Board) isSquareInBoardRange(row int, col int) bool {
	return row >= 0 && col >= 0 && row < board.Rows && col < board.Cols
}

// getAllValidMoves: filter out illegal moves from pseudo legal move map
func (board *Board) GetAllValidMoves(color Color) map[*Move]Piece {
	movesList := board.getAllPseudoLegalMoves(color)
	for move, piece := range movesList {
		if board.isDestOccupied(piece.Color, move.DestRow, move.DestCol) {
			move.Captured = board.Tiles[move.DestRow][move.DestCol].Piece
		}
		copyBoard := deepCopyBoard(board)
		copyBoard.PerformMove(piece, *move)
		if copyBoard.IsKingUnderCheck(color) {
			delete(movesList, move)
		}
		copyBoard.unmakeMove(piece, *move)
	}
	return movesList
}

// getAllPseudoLegalMoves: compiles all pseudo legal moves for given player
func (board *Board) getAllPseudoLegalMoves(color Color) map[*Move]Piece {
	var validMoves map[*Move]Piece = make(map[*Move]Piece)
	for rowIndex, row := range board.Tiles {
		for colIndex, tile := range row {
			if !tile.IsEmpty && tile.Piece.Color == color {
				for k, v := range board.genPieceMoves(tile.Piece, rowIndex, colIndex) {
					validMoves[k] = v
				}
			}
		}
	}
	return validMoves
}

// genPieceMoves: generate pseudo legal moves for given piece
func (board *Board) genPieceMoves(piece Piece, srcRow int, srcCol int) map[*Move]Piece {
	var validMoves = make(map[*Move]Piece)
	switch piece.Type {
	case Bishop:
		validMoves = board.genBishopMoves(piece, srcRow, srcCol)
	case Rook:
		validMoves = board.genRookMoves(piece, srcRow, srcCol)
	case Queen:
		validMoves = board.genRookMoves(piece, srcRow, srcCol)
		for k, v := range board.genBishopMoves(piece, srcRow, srcCol) {
			validMoves[k] = v
		}
	case King:
		for row := -1; row <= 1; row += 1 {
			for col := -1; col <= 1; col += 1 {
				if row == 0 && col == 0 {
					continue
				}
				if board.isSquareInBoardRange(srcRow+row, srcCol+col) && (!board.isSameColorPieceAtDest(piece.Color, srcRow+row, srcCol+col) || board.IsEmpty(srcRow+row, srcCol+col)) {
					move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: srcRow + row, DestCol: srcCol + col}
					validMoves[move] = piece
				}
			}
		}
	case Knight:
		jumpSquares := [][]int{{-2, -1}, {-2, 1}, {-1, 2}, {-1, -2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}
		for _, pair := range jumpSquares {
			targetRow := srcRow + pair[0]
			targetCol := srcCol + pair[1]
			if board.isSquareInBoardRange(targetRow, targetCol) && !board.isSameColorPieceAtDest(piece.Color, targetRow, targetCol) {
				move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: targetRow, DestCol: targetCol}
				validMoves[move] = piece
			}
		}
	case Pawn:
		var rowOffset, doubleMoveStartRank int
		if piece.Color == White {
			doubleMoveStartRank = board.Rows - 2
			rowOffset = -1
		} else {
			doubleMoveStartRank = 1
			rowOffset = 1
		}
		targetRow := srcRow + rowOffset
		for i := -1; i <= 1; i++ {
			//non-capture moves
			if i == 0 {
				if board.IsEmpty(targetRow, srcCol) {
					move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: targetRow, DestCol: srcCol}
					validMoves[move] = piece
				}
				if board.IsEmpty(targetRow+rowOffset, srcCol) && srcRow == doubleMoveStartRank {
					move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: targetRow + rowOffset, DestCol: srcCol}
					validMoves[move] = piece
				}
			}
			//capture moves
			if board.isSquareInBoardRange(targetRow, srcCol+i) && !board.IsEmpty(targetRow, srcCol+i) && !board.isSameColorPieceAtDest(piece.Color, targetRow, srcCol+i) {
				move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: targetRow, DestCol: srcCol + i, Captured: board.Tiles[targetRow][srcCol+i].Piece}
				validMoves[move] = piece
			}
		}
	}
	return validMoves
}

// genBishopMoves: generate pseudo-legal bishop moves
func (board *Board) genBishopMoves(piece Piece, srcRow int, srcCol int) map[*Move]Piece {
	var validMoves = make(map[*Move]Piece)
	diagonals := [][]int{{1, 1, board.Rows - 1}, {1, -1, board.Rows - 1}, {-1, 1, 0}, {-1, -1, 0}}
	for _, value := range diagonals {
		xOffset, yOffset, endRow := value[0], value[1], value[2]
		pathLength := Abs(srcRow - endRow)
		for i := 1; i <= pathLength; i++ {
			x := srcRow + i*xOffset
			y := srcCol + i*yOffset
			if board.isSquareInBoardRange(x, y) && !board.isSameColorPieceAtDest(piece.Color, x, y) {
				move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: x, DestCol: y}
				validMoves[move] = piece
				if !board.IsEmpty(x, y) {
					break
				}
			} else {
				break
			}
		}
	}
	return validMoves
}

func (board *Board) genRookMoves(piece Piece, srcRow int, srcCol int) map[*Move]Piece {
	directions := []int{-1, 1}
	//check horizontal
	var validMoves = make(map[*Move]Piece)
	for _, xOffset := range directions {
		for i := srcCol + xOffset; i >= 0 && i < board.Cols; i += xOffset {
			if !board.isSameColorPieceAtDest(piece.Color, srcRow, i) {
				move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: srcRow, DestCol: i}
				validMoves[move] = piece
			} else {
				break
			}
			if !board.IsEmpty(srcRow, i) {
				break
			}
		}
	}
	//check vertical
	for _, yOffset := range directions {
		for i := srcRow + yOffset; i >= 0 && i < board.Rows; i += yOffset {
			if !board.isSameColorPieceAtDest(piece.Color, i, srcCol) {
				move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: i, DestCol: srcCol}
				validMoves[move] = piece
			} else {
				break
			}
			if !board.IsEmpty(i, srcCol) {
				break
			}
		}
	}
	return validMoves
}

// isKingMoveValid: determine if a king move is valid
func isKingMoveValid(piece Piece, board *Board, move *Move) (bool, string) {
	if move.Castle {
		if !board.hasKingMoved(piece.Color) {
			if move.SrcRow == move.DestRow {
				var tile int
				if move.SrcCol+2 == move.DestCol && board.Tiles[move.SrcRow][board.Cols-1].Piece.Type == Rook { //castle to the right
					tile = move.SrcCol + 1
					for tile < board.Cols-1 {
						if !board.Tiles[move.SrcRow][tile].IsEmpty {
							return false, "castle path blocked"
						}
						tile += 1
					}
					return true, "valid castle"
				} else if move.SrcCol-2 == move.DestCol && board.Tiles[move.SrcRow][0].Piece.Type == Rook {
					tile = move.SrcCol - 1
					for tile > 0 {
						if !board.Tiles[move.SrcRow][tile].IsEmpty {
							return false, "castle path blocked"
						}
						tile -= 1
					}
					return true, "valid castle"
				} else {
					return false, "invalid castle dest"
				}
			}
		} else {
			return false, "king has already moved"
		}
	} else {
		if Abs(move.SrcRow-move.DestRow) <= 1 && Abs(move.SrcCol-move.DestCol) <= 1 {
			for row := -1; row <= 1; row += 1 {
				for col := -1; col <= 1; col += 1 {
					if (row == 0) && (col == 0) {
						continue
					}
					if board.isSquareInBoardRange(move.DestRow+row, move.DestCol+col) && board.Tiles[move.DestRow+row][move.DestCol+col].Piece.Type == King &&
						board.Tiles[move.DestRow+row][move.DestCol+col].Piece.Color == GetOpponentColor(piece.Color) {
						return false, "moving into check"
					}
				}
			}
			return true, "valid king move"
		}
	}
	return false, "invalid king move"
}

func (board *Board) hasKingMoved(color Color) bool {
	if color == White && board.WhiteKing.MoveCount > 0 || (color == Black && board.BlackKing.MoveCount > 0) {
		return true
	} else {
		return false
	}
}

// isRoomMoveValid determines if a rook move is valid
func isRookMoveValid(piece Piece, board *Board, move *Move) (bool, string) {
	//horizontal or vertical block
	if move.SrcCol == move.DestCol && move.SrcRow != move.DestRow {
		//vertical move
		start, end := Min(move.SrcRow, move.DestRow), Max(move.SrcRow, move.DestRow)
		for i := start + 1; i < end; i++ {
			if !board.IsEmpty(i, move.SrcCol) {
				return false, ("rook path blocked")
			}
		}
		return true, "valid rook move"
	} else if move.SrcRow == move.DestRow && move.SrcCol != move.DestCol {
		//horizontal move
		start, end := Min(move.SrcCol, move.DestCol), Max(move.SrcCol, move.DestCol)
		for i := start + 1; i < end; i++ {
			if !board.IsEmpty(move.SrcRow, i) {
				return false, ("rook path blocked")
			}
		}
		return true, "valid rook move"
	} else {
		return false, "invalid rook move"
	}
}

// isBishopMoveValid determines if a Bishop move is valid
func isBishopMoveValid(piece Piece, board *Board, move *Move) (bool, string) {
	pathLength := Abs(move.SrcRow - move.DestRow)
	if pathLength != Abs(move.SrcCol-move.DestCol) {
		return false, "not diagonal"
	}
	var xOffset, yOffset int
	if move.SrcRow < move.DestRow && move.SrcCol < move.DestCol {
		xOffset = 1
		yOffset = 1
	} else if move.SrcRow > move.DestRow && move.SrcCol < move.DestCol {
		xOffset = -1
		yOffset = 1
	} else if move.SrcRow < move.DestRow && move.SrcCol > move.DestCol {
		xOffset = 1
		yOffset = -1
	} else {
		xOffset = -1
		yOffset = -1
	}
	for i := 1; i < pathLength; i++ {
		x := move.SrcRow + i*xOffset
		y := move.SrcCol + i*yOffset
		if !board.IsEmpty(x, y) {
			// Obstacle found before reaching target: the move is invalid
			return false, "obstacle in bishop path"
		}
	}
	return true, "valid"
}

// isBishopMoveValid determines if a bishop move is valid
func isPawnMoveValid(piece Piece, board *Board, move *Move) (bool, string) {
	//not considering en passant yet
	var doubleMoveStartRank, rowOffset, promoteDestRow int
	if piece.Color == Black {
		doubleMoveStartRank = 1
		rowOffset = 1
		promoteDestRow = board.Rows - 1
	} else {
		doubleMoveStartRank = board.Rows - 2
		rowOffset = -1
		promoteDestRow = 0
	}
	if move.DestCol == move.SrcCol && move.SrcRow != move.DestRow {
		if Abs(move.SrcRow-move.DestRow) == 2 && move.SrcRow == doubleMoveStartRank {
			if (piece.Color == Black && board.IsEmpty(move.SrcRow+1, move.SrcCol)) || (piece.Color == White && board.IsEmpty(move.SrcRow-1, move.SrcCol)) {
				return true, "double pawn move allowed"
			} else {
				return false, "double move blocked"
			}
		} else if Abs(move.SrcRow-move.DestRow) == 1 && !piece.isBackwardPawnMove(move) {
			if board.IsEmpty(move.DestRow, move.DestCol) {
				if move.Promote != 0 && move.DestRow == promoteDestRow {
					return true, "pawn promoted"
				}
				return true, "valid single pawn move"
			} else {

				return false, "dest blocked"
			}
		}

	} else if move.DestCol != move.SrcCol && move.SrcRow != move.DestRow {
		//check if col row +-1 logic
		if Abs(move.SrcCol-move.DestCol) == 1 && move.DestRow == move.SrcRow+rowOffset && !board.IsEmpty(move.DestRow, move.DestCol) {
			return true, "valid pawn capture"
		}
		return false, "invalid pawn capture"
	}
	return false, "not a valid pawn move"
}

func ChangeTurn(turn string) string {
	if turn == "w" {
		return "b"
	} else {
		return "w"
	}
}

func isCustomMoveValid(piece Piece, board *Board, move *Move) (bool, string) {
	//Check jump moves followed by slide moves
	var jumpPattern, slidePattern [][]int
	//find pattern for the piece, will change this to a hashmap later
	for _, movePatterns := range board.CustomMovePatterns {
		if movePatterns.PieceName == strings.ToLower(piece.CustomPiece.PieceName) {
			jumpPattern = movePatterns.JumpPattern
			slidePattern = movePatterns.SlidePattern
			break
		}
	}

	var rowDiff, colDiff, multiplier int
	if piece.Color == White{
		rowDiff,colDiff = move.DestRow - move.SrcRow, move.DestCol - move.SrcCol
		multiplier = 1
	} else {
		rowDiff,colDiff = move.SrcRow-move.DestRow, move.SrcCol - move.DestCol
		multiplier = -1
	}
	if len(jumpPattern) != 0 {		
		for _, pair := range jumpPattern {
			if pair[0] == rowDiff && pair[1] == colDiff {
				return true, "valid custom jump move"
			}
		}
	}

	for _, direction := range slidePattern {
		var dx,dy int = direction[0]*multiplier, direction[1]*multiplier 
		var tempRow, tempCol int = move.SrcRow+dx, move.SrcCol + dy
		for tempRow >= 0 && tempCol >= 0 && tempRow < board.Rows && tempCol < board.Cols {
			if tempRow == move.DestRow && tempCol == move.DestCol {
				return true, "valid custom slide move"
			} else if !board.IsEmpty(tempRow, tempCol) {
				break
			} else {
				tempRow += dx
				tempCol += dy
			}
		}
	}

	return false, "no"
}
