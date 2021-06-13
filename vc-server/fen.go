package main

import (
	"strconv"
	"strings"
	"unicode"
	"bytes"
	"fmt"
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
	fmt.Println(colCount,rowsData)
	board := &Board{
		tiles: make([][]rune, len(rowsData)),
		rows:  len(rowsData),
		cols:  colCount,
	}
	var col int = 0
	for rowIndex, row := range rowsData {
		col = 0
		board.tiles[rowIndex] = make([]rune, board.cols)
		for _, char := range row {
			fmt.Println(rowIndex,col,char,unicode.IsNumber(rune(char)))
			if unicode.IsNumber(rune(char)) {
					count,_ := strconv.Atoi(string(char))
					col+=count-1
			} else{board.tiles[rowIndex][col] = rune(char) }
			col++
		}
	}
	return board
}

func ConvertBoardtoFEN(board *Board) string{
	var fen bytes.Buffer
	for row:=0;row<board.rows; row++{
		var empty int = 0
		for col:=0;col<board.cols;col++{
			if board.tiles[row][col]==0{
				empty+=1
			} else{
				if empty>0{
					str:=strconv.Itoa(empty)
					fen.WriteString(str)
					empty=0
				}
				fen.WriteString(string(board.tiles[row][col]))
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