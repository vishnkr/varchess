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
	var secDigit = 0
	//count columns
	for index, char := range rowsData[0] {
		if unicode.IsNumber(rune(char)) {
			count, _ := strconv.Atoi(string(char))
			if (index+1<len(rowsData[0]) && unicode.IsNumber(rune(rowsData[0][index+1]))){
				secDigit,_ = strconv.Atoi(string(char))
				fmt.Println("secdigit is",secDigit)
			} else{ 
				if (secDigit!=0){
					colCount+=secDigit*10+count
				} else {colCount += count}
			} 
		} else { colCount += 1}
	}
	board := &Board{
		Tiles: make([][]Square, len(rowsData)),
		Rows:  len(rowsData),
		Cols:  colCount,
	}
	var col,id int = 0,0
	var colEnd int = 0
	for rowIndex, row := range rowsData {
		col = 0
		board.Tiles[rowIndex] = make([]Square, board.Cols)
		secDigit = 0
		for index, char := range row {
			if unicode.IsNumber(rune(char)) {
				if (index+1<len(row) && unicode.IsNumber(rune(row[index+1]))){
					secDigit,_ = strconv.Atoi(string(char))
					fmt.Println("secdigit is",secDigit)
				} else{ 
					count,_ := strconv.Atoi(string(char))
					if (secDigit!=0){
						colEnd = secDigit*10+count
						secDigit = 0
					} else { colEnd = count}
					i:= col
					for (col < i+colEnd){
						fmt.Println("row-",rowIndex,"col-",col,colEnd,i)
						board.Tiles[rowIndex][col] = Square{IsEmpty:true,Id:id}
						col++	
						id+=1
					}
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
					Id:id,
				}
				val,ok := strToTypeMap[string(unicode.ToLower(char))]
				
				if (!ok){
					customPiece:=&CustomPiece{PieceName:string(char)}
					board.Tiles[rowIndex][col].Piece = Piece{Color: color,Type:Custom}
					board.Tiles[rowIndex][col].Piece.CustomPiece = customPiece
				} else {
					board.Tiles[rowIndex][col].Piece = Piece{Color: color,Type:val}
					board.Tiles[rowIndex][col].Piece.Type = val
					if (val==King){
						if (color==Black){
							board.BlackKing.Position = []int{rowIndex,col}
						} else {
							board.WhiteKing.Position = []int{rowIndex,col}
						}
					}
				}
				col++
				id+=1
			}
		}
	}
	return board
}

func ConvertBoardtoFEN(board *Board) string{
	var fen bytes.Buffer
	var name string
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
				if (board.Tiles[row][col].Piece.Type==Custom){
					name = board.Tiles[row][col].Piece.CustomPiece.PieceName
				} else { name = typeToStrMap[board.Tiles[row][col].Piece.Type]}
				if (board.Tiles[row][col].Piece.Color==White){
					name = strings.ToUpper(name)
				}
				fen.WriteString(name)
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