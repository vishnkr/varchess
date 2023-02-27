package store

type User struct{
	ID int 
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type Template struct{
	ID int
	Name string
	FEN string
	CustomPieces []Piece
}

type Piece struct{
	ID int
	Name string
	MovePatterns []MovePattern
}

type MovePattern struct{
	X int
	Y int
	Type string
}