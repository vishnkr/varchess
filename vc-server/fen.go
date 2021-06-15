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
		tiles: make([][]Square, len(rowsData)),
		rows:  len(rowsData),
		cols:  colCount,
	}
	var col int = 0
	for rowIndex, row := range rowsData {
		col = 0
		board.tiles[rowIndex] = make([]Square, board.cols)
		for _, char := range row {
			
			if unicode.IsNumber(rune(char)) {
					count,_ := strconv.Atoi(string(char))
					i:= col
					for (col < i+count){
						board.tiles[rowIndex][col] = Square{IsEmpty:true}
						col++	
					}
			} else{
				var color Color
				if(unicode.IsUpper(rune(char))){
					color = White
				} else{
					color = Black
				}
				
				board.tiles[rowIndex][col] = Square{
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
	for row:=0;row<board.rows; row++{
		var empty int = 0
		for col:=0;col<board.cols;col++{
			if board.tiles[row][col].IsEmpty{
				empty+=1
			} else{
				if empty>0{
					str:=strconv.Itoa(empty)
					fen.WriteString(str)
					empty=0
				}
				if (board.tiles[row][col].Piece.Color==White){
					fen.WriteString(string(unicode.ToUpper(typeToRuneMap[board.tiles[row][col].Piece.Type])))
				} else{
					fen.WriteString(string(typeToRuneMap[board.tiles[row][col].Piece.Type]))
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