package game

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
	allowedPromotionTypes map[pieceType]bool
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
	pieceProps map[rune]pieceProperties
	turn       color
	dimensions
	disabledSquares map[int]bool
	castlingRights
	attackedSquares map[int]bool
	additionalProps AdditionalProps
}

type JumpMove struct {
	Offset           moveOffset `json:"offset"`
	IsCaptureAllowed bool       `json:"isCaptureAllowed"`
}

type pieceProperties struct {
	Name         string       `json:"name"`
	SlideOffsets []moveOffset `json:"slideOffsets,omitempty"`
	JumpProps    []JumpMove   `json:"jumpProps,omitempty"`
}

type gameConfigPosition struct {
	Dimensions dimensions               `json:"dimensions"`
	Fen        string                   `json:"fen"`
	PieceProps map[rune]pieceProperties `json:"pieceProps,omitempty"`
}

type GameConfig struct {
	VariantType        variantType `json:"variantType"`
	gameConfigPosition `json:"position"`
	Objective          Objective `json:"objective"`
}

func (p *position) addStandardPieceProps() {
	diagonals := []moveOffset{
		{xOffset: -1, yOffset: -1},
		{xOffset: 1, yOffset: -1},
		{xOffset: 1, yOffset: 1},
		{xOffset: -1, yOffset: 1},
	}
	nonDiagonals := []moveOffset{
		{xOffset: -1, yOffset: 0},
		{xOffset: 1, yOffset: 0},
		{xOffset: 0, yOffset: 1},
		{xOffset: 0, yOffset: -1},
	}
	props := map[rune]pieceProperties{
		'q': pieceProperties{
			Name:         "Queen",
			SlideOffsets: append(nonDiagonals, diagonals...),
		},
		'b': pieceProperties{
			Name:         "Bishop",
			SlideOffsets: diagonals,
		},
		'r': pieceProperties{
			Name:         "Rook",
			SlideOffsets: nonDiagonals,
		},
		'n': pieceProperties{
			Name: "Knight",
			JumpProps: []JumpMove{
				JumpMove{moveOffset{xOffset: -2, yOffset: 1}, true},
				JumpMove{moveOffset{xOffset: -2, yOffset: -1}, true},
				JumpMove{moveOffset{xOffset: 2, yOffset: 1}, true},
				JumpMove{moveOffset{xOffset: 2, yOffset: -1}, true},
				JumpMove{moveOffset{xOffset: 1, yOffset: -2}, true},
				JumpMove{moveOffset{xOffset: 1, yOffset: 2}, true},
				JumpMove{moveOffset{xOffset: -1, yOffset: 2}, true},
				JumpMove{moveOffset{xOffset: -1, yOffset: -2}, true},
			},
		}}
	p.pieceProps = props
}
func newPosition(gameConfig GameConfig) (position, error) {
	boardData := strings.Split(gameConfig.Fen, " ")
	rowsData := strings.Split(boardData[0], "/")
	var standardPieceMap = map[rune]pieceType{'p': Pawn, 'k': King, 'n': Knight, 'q': Queen, 'r': Rook, 'b': Bishop}
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
			if unicode.IsNumber(rune(char)) {
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
					position.disabledSquares[id] = true
					col++
					id += 1
					continue
				}
				var color color
				if unicode.IsUpper(char) {
					color = ColorWhite
				} else {
					color = ColorBlack
				}
				pieceRune := unicode.ToLower(char)
				if !isStandardPiece(pieceRune) {
					position.pieceLocations[id] = piece{
						color:     color,
						notation:  pieceRune,
						pieceType: CustomPiece,
					}
					position.pieceProps[pieceRune] = gameConfig.PieceProps[pieceRune]
				} else {
					piece := piece{
						color:     color,
						notation:  pieceRune,
						pieceType: standardPieceMap[pieceRune],
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
	}
	position.attackedSquares = map[int]bool{}
	return position, nil
}

func (p *position) getOpponentColor() color {
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

func (p *position) isDisabled(targetSquareID int) bool {
	_, ok := p.disabledSquares[targetSquareID]
	return ok
}
func (p *position) getclassicMoveType(targetSquareID int) classicMoveType {
	if _, ok := p.pieceLocations[targetSquareID]; ok {
		return CaptureMove
	}
	return QuietMove

}

func (p *position) switchTurn() { p.turn = p.getOpponentColor() }

func (p *position) toRowCol(squareId int) (int, int) {
	return squareId / p.Files, squareId % p.Files
}

func (p *position) toPos(row int, col int) int {
	return row*p.Files + col
}

func isStandardPiece(pieceName rune) bool {
	pieces := []rune{'p', 'k', 'n', 'b', 'r', 'q'}
	for _, p := range pieces {
		if p == pieceName {
			return true
		}
	}
	return false
}
