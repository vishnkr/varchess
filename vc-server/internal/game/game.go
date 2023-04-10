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
	Turn     Color
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

func IsSameMove(move1 Move, move2 Move) bool {
	return move1.SrcCol == move2.SrcCol && move1.SrcRow == move2.SrcRow && move1.DestCol == move2.DestCol && move1.DestRow == move2.DestRow
}

func (g *Game) ChangeTurn() {
	if g.Turn == White {
		g.Turn = Black
	} else {
		g.Turn = White
	}
}

//visual helper function
func DisplayBoardState(board *Board) {
	var piece string
	for i := 0; i < board.Rows; i++ {
		for j := 0; j < board.Cols; j++ {
			if board.Tiles[i][j].IsDisabled {
				piece = "X"
			} else if !board.Tiles[i][j].IsEmpty {
				if board.Tiles[i][j].Piece.CustomPieceName != "" {
					piece = board.Tiles[i][j].Piece.CustomPieceName
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
func (board *Board) UnmakeMove(piece Piece, move Move) {
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
	oppMoves := board.GetAllPseudoLegalMoves(GetOpponentColor(color))
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
	movesList := board.GetAllPseudoLegalMoves(color)
	for move, piece := range movesList {
		if board.isDestOccupied(piece.Color, move.DestRow, move.DestCol) {
			move.Captured = board.Tiles[move.DestRow][move.DestCol].Piece
		}
		copyBoard := deepCopyBoard(board)
		copyBoard.PerformMove(piece, *move)
		if copyBoard.IsKingUnderCheck(color) {
			delete(movesList, move)
		}
	}
	return movesList
}

// getAllPseudoLegalMoves: compiles all pseudo legal moves for given player
func (board *Board) GetAllPseudoLegalMoves(color Color) map[*Move]Piece {
	var validMoves map[*Move]Piece = make(map[*Move]Piece)
	for rowIndex, row := range board.Tiles {
		for colIndex, tile := range row {
			if !tile.IsEmpty && tile.Piece.Color == color {
				for k, v := range board.GenPieceMoves(tile.Piece, rowIndex, colIndex) {

					validMoves[k] = v
				}
			}
		}
	}
	return validMoves
}

// genPieceMoves: generate pseudo legal moves for given piece
func (board *Board) GenPieceMoves(piece Piece, srcRow int, srcCol int) map[*Move]Piece {
	var validMoves = make(map[*Move]Piece)
	switch piece.Type {
	case Bishop:
		dirs := [][]int{{1, -1}, {1, 1}, {-1, 1}, {-1, -1}}
		validMoves = board.genSlideMoves(dirs, piece, srcRow, srcCol)
	case Rook:
		dirs := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
		validMoves = board.genSlideMoves(dirs, piece, srcRow, srcCol)
	case Queen:
		dirs := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}
		validMoves = board.genSlideMoves(dirs, piece, srcRow, srcCol)
	case King:
		for row := -1; row <= 1; row += 1 {
			for col := -1; col <= 1; col += 1 {
				if row == 0 && col == 0 {
					continue
				}
				destRow, destCol := srcRow+row, srcCol+col
				if board.isSquareInBoardRange(destRow, destRow) && (!board.isSameColorPieceAtDest(piece.Color, destRow, destRow) || board.IsEmpty(srcRow+row, srcCol+col)) && !board.isSquareDisabled(destRow, destRow) {
					move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: destRow, DestCol: destCol}
					validMoves[move] = piece
				}
			}
		}
		//validMoves board.genCastleMoves(piece,srcRow,srcCol)
		if !board.hasKingMoved(piece.Color) {
			var queenSideCastle = &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: srcRow, DestCol: srcCol - 2, Castle: true}
			var kingSideCastle = &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: srcRow, DestCol: srcCol + 2, Castle: true}
			if board.isValidCastle(queenSideCastle) {
				validMoves[queenSideCastle] = piece
			}
			if board.isValidCastle(kingSideCastle) {
				validMoves[kingSideCastle] = piece
			}
		}

	case Knight:
		jumpSquares := [][]int{{-2, -1}, {-2, 1}, {-1, 2}, {-1, -2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}
		for _, pair := range jumpSquares {
			targetRow := srcRow + pair[0]
			targetCol := srcCol + pair[1]
			if board.isSquareInBoardRange(targetRow, targetCol) && !board.isSameColorPieceAtDest(piece.Color, targetRow, targetCol) && !board.isSquareDisabled(targetRow, targetCol) {
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
				if board.IsEmpty(targetRow, srcCol) && !board.isSquareDisabled(targetRow, srcCol) {
					move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: targetRow, DestCol: srcCol}
					validMoves[move] = piece
				}
				if board.IsEmpty(targetRow+rowOffset, srcCol) && srcRow == doubleMoveStartRank && !board.isSquareDisabled(targetRow+rowOffset, srcCol) {
					move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: targetRow + rowOffset, DestCol: srcCol}
					validMoves[move] = piece
				}
			}
			//capture moves
			if board.isSquareInBoardRange(targetRow, srcCol+i) && !board.IsEmpty(targetRow, srcCol+i) && !board.isSquareDisabled(targetRow, srcCol+i) && !board.isSameColorPieceAtDest(piece.Color, targetRow, srcCol+i) {
				move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: targetRow, DestCol: srcCol + i, Captured: board.Tiles[targetRow][srcCol+i].Piece}
				validMoves[move] = piece
			}
		}
	case Custom:
		return board.genCustomMoves(piece, srcRow, srcCol)
	}

	return validMoves
}

func (board *Board) isValidCastle(move *Move) bool {
	var tile int
	if move.SrcCol+2 == move.DestCol && board.Tiles[move.SrcRow][board.Cols-1].Piece.Type == Rook { //castle to the right
		tile = move.SrcCol + 1
		for tile < board.Cols-1 {
			if !board.Tiles[move.SrcRow][tile].IsEmpty || board.isSquareDisabled(move.SrcRow, tile) {
				return false
			}
			tile += 1
		}
		return true
	} else if move.SrcCol-2 == move.DestCol && board.Tiles[move.SrcRow][0].Piece.Type == Rook {
		tile = move.SrcCol - 1
		for tile > 0 {
			if !board.Tiles[move.SrcRow][tile].IsEmpty || board.isSquareDisabled(move.SrcRow, tile) {
				return false
			}
			tile -= 1
		}
		return true
	}
	return false
}

func (board *Board) hasKingMoved(color Color) bool {
	if color == White && board.WhiteKing.MoveCount > 0 || (color == Black && board.BlackKing.MoveCount > 0) {
		return true
	} else {
		return false
	}
}

func (board *Board) genCustomMoves(piece Piece, srcRow int, srcCol int) map[*Move]Piece {
	var validMoves = make(map[*Move]Piece)
	for _, movePatterns := range board.CustomMovePatterns {
		if movePatterns.PieceName == strings.ToLower(piece.CustomPieceName) {
			multiplier := 1
			if piece.Color == Black {
				multiplier = -1
			}
			for _, pair := range movePatterns.JumpPattern {
				var destRow, destCol int = srcRow + (multiplier * pair[0]), srcCol + (multiplier * pair[1])
				if board.isSquareInBoardRange(destRow, destCol) && !(piece.Color == board.getPieceColor(destRow, destCol)) && !board.isSquareDisabled(destRow, destCol) {
					move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: destRow, DestCol: destCol}
					validMoves[move] = piece
				}
			}
			for k, v := range board.genSlideMoves(movePatterns.SlidePattern, piece, srcRow, srcCol) {
				validMoves[k] = v
			}
			break
		}
	}
	return validMoves
}

func (board *Board) genSlideMoves(directions [][]int, piece Piece, srcRow int, srcCol int) map[*Move]Piece {
	var validMoves = make(map[*Move]Piece)
	multiplier := 1
	if piece.Color == Black {
		multiplier = -1
	}
	for _, direction := range directions {
		var dx, dy int = direction[0] * multiplier, direction[1] * multiplier
		var tempRow, tempCol int = srcRow + dx, srcCol + dy
		for board.isSquareInBoardRange(tempRow, tempCol) {
			if !board.isSameColorPieceAtDest(piece.Color, tempRow, tempCol) && !board.isSquareDisabled(tempRow, tempCol) {
				move := &Move{SrcRow: srcRow, SrcCol: srcCol, DestRow: tempRow, DestCol: tempCol}
				validMoves[move] = piece
				if !board.IsEmpty(tempRow, tempCol) {
					break
				}
			} else {
				break
			}
			tempRow += dx
			tempCol += dy

		}
	}
	return validMoves
}

func (board *Board) isSquareDisabled(row int, col int) bool {
	return board.Tiles[row][col].IsDisabled
}
