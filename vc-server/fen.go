package main

import (
	"strconv"
	"strings"
	"unicode"
	"bytes"
	//"fmt"
)

func ConvertFENtoBoard(fen string) *Board {
	//TODO: add enpassant, castling, turn to FEN
	boardData := strings.Split(fen, " ")
	rowsData := strings.Split(boardData[0], "/")
	var colCount int = 0
	//count columns
	for _, char := range rowsData[0] {
		if unicode.IsNumber(rune(char)) {
			empty, _ := strconv.Atoi(string(char))
			colCount += empty
		} else {
			colCount += 1
		}
	}
	board := &Board{
		Tiles: make([][]Square, len(rowsData)),
		Rows:  len(rowsData),
		Cols:  colCount,
	}
	var col int = 0
	for rowIndex, row := range rowsData {
		col = 0
		board.Tiles[rowIndex] = make([]Square, board.Cols)
		for _, char := range row {
			
			if unicode.IsNumber(rune(char)) {
					count,_ := strconv.Atoi(string(char))
					i:= col
					for (col < i+count){
						board.Tiles[rowIndex][col] = Square{IsEmpty:true}
						col++	
					}
			} else{
				var color Color
				if(unicode.IsUpper(rune(char))){
					color = White
				} else{
					color = Black
				}
				
				board.Tiles[rowIndex][col] = Square{
											IsEmpty:false, 
											Piece: Piece{
												Type: strToTypeMap[string(unicode.ToLower(char))],
												Color: color,
											},
										}
				col++
			}
			
		}
	}
	return board
}

func ConvertBoardtoFEN(board *Board) string{
	var fen bytes.Buffer
	for row:=0;row<board.Rows; row++{
		var empty int = 0
		for col:=0;col<board.Cols;col++{
			if board.Tiles[row][col].IsEmpty{
				empty+=1
			} else{
				if empty>0{
					str:=strconv.Itoa(empty)
					fen.WriteString(str)
					empty=0
				}
				if (board.Tiles[row][col].Piece.Color==White){
					fen.WriteString(string(unicode.ToUpper(typeToRuneMap[board.Tiles[row][col].Piece.Type])))
				} else{
					fen.WriteString(string(typeToRuneMap[board.Tiles[row][col].Piece.Type]))
				}
				
			}
			
		}
		if empty>0{
			str:=strconv.Itoa(empty)
			fen.WriteString(str)
			empty=0
		}
		fen.WriteString("/")
	}
	fenString := fen.String()[:len(fen.String())-1]
	fenString += " w KQkq - 0 1"
	return fenString
}