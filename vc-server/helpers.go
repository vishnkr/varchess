package main 
func Max(x, y int) int {
    if x < y {
        return y
    }
    return x
}
func Abs(x int)int{
	if x<0{
		return -x
	}
	return x
}
func Min(x, y int) int {
    if x > y {
        return y
    }
    return x
}

func deepCopyBoard(board *Board) *Board{
    copy := &Board{
		Tiles: make([][]Square, board.Rows),
		Rows:  board.Rows,
		Cols:  board.Cols,
	}
    for i := range copy.Tiles {
        copy.Tiles[i] = make([]Square, board.Cols)
        for j := range copy.Tiles[i] {
            copy.Tiles[i][j].Id = board.Tiles[i][j].Id
            copy.Tiles[i][j].Piece = board.Tiles[i][j].Piece
            copy.Tiles[i][j].IsEmpty = board.Tiles[i][j].IsEmpty
            copy.BlackKing = board.BlackKing
            copy.WhiteKing = board.WhiteKing
        }
    }
    return copy
}