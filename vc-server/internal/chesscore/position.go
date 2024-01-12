package chesscore

import (
	"strconv"
	"strings"
	"unicode"
)

type castlingRights struct {
	whiteKingSide  bool
	whiteQueenSide bool
	blackKingSide  bool
	blackQueenSide bool
}


type promotionProps struct {
	promotionSquares      map[int]bool
	allowedPromotionTypes []pieceType
}

type AdditionalProps struct {
	kingCaptureAllowed bool
	blackKingMoved     bool
	whiteKingMoved     bool
	blackKingPos       int
	whiteKingPos       int
}

type position struct {
	pieceLocations map[int]piece
	//pieceProps stores move patterns for custom pieces
	pieceProps map[string]PieceProperties
	turn       Color
	dimensions Dimensions
	wallSquares map[int]bool
	castlingRights
	attackedSquares map[int]bool
	additionalProps AdditionalProps
	promotionProps promotionProps
	epSquare int
}

/*type JumpMove struct {
	Offset           moveOffset `json:"offset"`
	IsCaptureAllowed bool       `json:"isCaptureAllowed"`
}*/

type PieceProperties struct {
	Name         string       `json:"name"`
	SlideOffsets []moveOffset `json:"slide_offsets,omitempty"`
	JumpOffsets    []moveOffset   `json:"jump_offsets,omitempty"`
}

type Position struct {
	Dimensions Dimensions               `json:"dimensions"`
	Fen        string                   `json:"fen"`
	PieceProps map[string]PieceProperties `json:"piece_props,omitempty"`
}

type GameConfig struct {
	VariantType        variantType `json:"variant_type"`
	Dimensions Dimensions               `json:"dimensions"`
	Fen        string                   `json:"fen"`
	PieceProps map[string]PieceProperties `json:"piece_props,omitempty"`
	AdditionalData interface{} `json:"additonal_data,omitempty"`
}

func (p *position) addStandardPieceProps() {
	diagonals := []moveOffset{
		{x: -1, y: -1},
		{x: 1, y: -1},
		{x: 1, y: 1},
		{x: -1, y: 1},
	}
	nonDiagonals := []moveOffset{
		{x: -1, y: 0},
		{x: 1, y: 0},
		{x: 0, y: 1},
		{x: 0, y: -1},
	}
	props := map[string]PieceProperties{
		"q": {
			Name:         "Queen",
			SlideOffsets: append(nonDiagonals, diagonals...),
		},
		"b": {
			Name:         "Bishop",
			SlideOffsets: diagonals,
		},
		"r": {
			Name:         "Rook",
			SlideOffsets: nonDiagonals,
		},
		"n": {
			Name: "Knight",
			JumpOffsets: []moveOffset{
				{x: -2, y: 1},
				{x: -2, y: -1},
				{x: 2, y: 1},
				{x: 2, y: -1},
				{x: 1, y: -2},
				{x: 1, y: 2},
				{x: -1, y: 2},
				{x: -1, y: -2},
			},
		}}
	p.pieceProps = props
}
func newPosition(gameConfig GameConfig) (position, error) {
	boardData := strings.Split(gameConfig.Fen, " ")
	rowsData := strings.Split(boardData[0], "/")
	var standardPieceMap = map[string]pieceType{"p": Pawn,"k": King, "n": Knight, "q": Queen, "r": Rook, "b": Bishop}
	position := position{
		pieceLocations: make(map[int]piece),
		dimensions:     gameConfig.Dimensions,
		additionalProps: AdditionalProps{
			blackKingMoved:     false,
			whiteKingMoved:     false,
			kingCaptureAllowed: false,
		},
	}

	var col, id int = 0, 0
	var colEnd int = 0
	var secDigit = 0
	position.addStandardPieceProps()
	for _, row := range rowsData {
		col = 0
		secDigit = 0
		for index, char := range row {
			if unicode.IsNumber(char) {
				if index+1 < len(row) && unicode.IsNumber(rune(row[index+1])) {
					secDigit, _ = strconv.Atoi(string(char))
				} else {
					count, _ := strconv.Atoi(string(char))
					if secDigit != 0 {
						colEnd = secDigit*10 + count
						secDigit = 0
					} else {
						colEnd = count
					}
					i := col
					for col < i+colEnd {
						col++
						id += 1
					}
				}
			} else {
				if char == '.' {
					position.wallSquares[id] = true
					col++
					id += 1
					continue
				}
				var color Color
				if unicode.IsUpper(char) {
					color = ColorWhite
				} else {
					color = ColorBlack
				}
				pieceChar := string(char)
				
				if !isStandardPiece(pieceChar) {
					position.pieceLocations[id] = piece{
						color:     color,
						notation:  pieceChar,
						pieceType: CustomPiece,
					}
					position.pieceProps[pieceChar] = gameConfig.PieceProps[pieceChar]
				} else {
					piece := piece{
						color:     color,
						notation:  pieceChar,
						pieceType: standardPieceMap[pieceChar],
					}
					position.pieceLocations[id] = piece
					if piece.pieceType == King {
						if color == ColorBlack {
							position.additionalProps.blackKingPos = id
						} else {
							position.additionalProps.whiteKingPos = id
						}
					}
				}
				col++
				id += 1
			}
		}
		if boardData[1] == "w" {
			position.turn = ColorWhite
		} else {
			position.turn = ColorBlack
		}
		if strings.Contains(boardData[2], "k") {
			position.castlingRights.blackKingSide = true
		}
		if strings.Contains(boardData[2], "K") {
			position.castlingRights.whiteKingSide = true
		}
		if strings.Contains(boardData[2], "q") {
			position.castlingRights.blackQueenSide = true
		}
		if strings.Contains(boardData[2], "Q") {
			position.castlingRights.whiteQueenSide = true
		}
		sqId,err:= strconv.Atoi(boardData[3]); 
		if err!=nil {
			position.epSquare = -1
		} else { position.epSquare = sqId}
	}
	position.attackedSquares = map[int]bool{}
	return position, nil
}

func (p *position) getOpponentColor() Color {
	if p.turn == ColorBlack {
		return ColorWhite
	}
	return ColorBlack
}

func (p *position) isOppKingAtTarget(targetSquareID int) bool {
	//assume opponent color check is already done
	return p.pieceLocations[targetSquareID].pieceType == King
}

func (p *position) isSameColorPiecePresent(sourceSquareID int, targetSquareID int) bool {
	if p1, ok := p.pieceLocations[sourceSquareID]; ok {
		if p2, ok := p.pieceLocations[targetSquareID]; ok {
			return p1.color == p2.color
		}
	}
	return false
}

func (p *position) isOpponentPiecePresent(sourceSquareID int, targetSquareID int) bool {
	if p1, ok := p.pieceLocations[sourceSquareID]; ok {
		if p2, ok := p.pieceLocations[targetSquareID]; ok {
			return p1.color != p2.color
		}
	}
	return false
}

func (p *position) isWall(targetSquareID int) bool {
	_, ok := p.wallSquares[targetSquareID]
	return ok
}

func (p *position) isEmpty(targetSquareID int) bool{
	_, ok := p.pieceLocations[targetSquareID]
	return !ok
}
func (p *position) getclassicMoveType(targetSquareID int) classicMoveType {
	if _, ok := p.pieceLocations[targetSquareID]; ok {
		return CaptureMove
	}
	return QuietMove

}

func (p *position) switchTurn() { p.turn = p.getOpponentColor() }

func (p *position) toRowCol(squareId int) (int, int) {
	return squareId / p.dimensions.Files, squareId % p.dimensions.Files
}

func (p *position) toPos(row int, col int) int {
	return row*p.dimensions.Files + col
}

func isStandardPiece(piece string) bool {
	pieces := []string{"p", "k", "n", "b", "r", "q"}
	for _, p := range pieces {
		if p == piece {
			return true
		}
	}
	return false
}
