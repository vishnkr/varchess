package chess

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"vc-core/internal/models"

	"github.com/holiman/uint256"
)


type Bitboard struct {
    Bits *uint256.Int
}

type Position struct {
    Files      int
    Ranks      int
    WhiteToMove bool
    Castling   uint8
    EnPassant  int
    HalfMove   int
    FullMove   int
    Pieces     map[rune]*Bitboard
    Walls      *Bitboard
}


type MoveType int
const (
    Slide MoveType = iota 
    Jump
    PawnMove
    CastlingMove
    EnPassantMove
)

type PieceDefinition struct {
    MoveType MoveType
    MovePatterns []Bitboard
}


type Color uint8
const (
    Black Color = iota
    White
)

func NewBitboard() *Bitboard {
    return &Bitboard{Bits: uint256.NewInt(0)}
}

func Index(file, rank, width int) int {
    return rank*width + file
}

func (b *Bitboard) SetBit(index int) {
	mask := uint256.NewInt(0).SetUint64(1)
	mask.Lsh(mask, uint(index)) // mask = 1 << index
	b.Bits.Or(b.Bits, mask) // b.Bits |= mask
}


func (b *Bitboard) ClearBit(index int) {
	mask := uint256.NewInt(0).SetUint64(1)
	mask.Lsh(mask, uint(index)) // mask = 1 << index
	mask.Not(mask)              // mask = ~mask
	b.Bits.And(b.Bits, mask)    // b.Bits &= mask
}

func (b *Bitboard) HasBit(index int) bool {
	mask := uint256.NewInt(0).SetUint64(1)
	mask.Lsh(mask, uint(index)) // mask = 1 << index
	return b.Bits.And(b.Bits, mask).Cmp(uint256.NewInt(0)) != 0
}

type Game struct {
    ID        string
    Position  Position
    MoveHistory []string
}


func InferDimensions(fen string) (int, int, error) {
    ranks := strings.Split(fen, "/")
    numRanks := len(ranks)

    var expectedFiles int = -1
    for i, rank := range ranks {
        files := 0
        numBuffer := ""

        for _, ch := range rank {
            if unicode.IsDigit(ch) {
                numBuffer += string(ch)
            } else {
                if numBuffer != "" {
                    n, _ := strconv.Atoi(numBuffer)
                    files += n
                    numBuffer = ""
                }
                files++
            }
        }

        if numBuffer != "" {
            n, _ := strconv.Atoi(numBuffer)
            files += n
        }

        if expectedFiles == -1 {
            expectedFiles = files
        } else if files != expectedFiles {
            return 0, 0, fmt.Errorf("invalid FEN: inconsistent file count in rank %d", i+1)
        }
    }

    return numRanks, expectedFiles, nil
}


func ParseFEN(fen string) (*Position,error) {

    
    parts := strings.Split(fen, " ")
    board := parts[0]
    sideToMove := parts[1] == "w"
    castlingRights := ParseCastlingRights(parts[2])
    ranks, files,err := InferDimensions(board)
    if err!=nil{
        return nil,err
    }
    enPassant := ParseEnPassant(parts[3],files)
    
    pos := Position{
        Files:      files,
        Ranks:      ranks,
        WhiteToMove: sideToMove,
        Castling:   castlingRights,
        EnPassant:  enPassant,
        Pieces:     make(map[rune]*Bitboard),
        Walls:      NewBitboard(),
    }

    rank := 0
    file := 0

    for _, ch := range board {
        if ch == '/' {
            rank++
            file = 0
            continue
        }

        if ch >= '1' && ch <= '9' {
            file += int(ch - '0')
            continue
        }

        index := Index(file, rank, files)

        if ch == '.' {
            pos.Walls.SetBit(index)
        } else {
            if _, exists := pos.Pieces[ch]; !exists {
                pos.Pieces[ch] = NewBitboard()
            }
            pos.Pieces[ch].SetBit(index)
        }

        file++
    }

    return &pos,nil
}

// Castling Rights Parsing
func ParseCastlingRights(castling string) uint8 {
    var rights uint8
    if strings.Contains(castling, "K") { rights |= 1 << 3 }
    if strings.Contains(castling, "Q") { rights |= 1 << 2 }
    if strings.Contains(castling, "k") { rights |= 1 << 1 }
    if strings.Contains(castling, "q") { rights |= 1 << 0 }
    return rights
}


func ParseEnPassant(enPassant string, files int) int {
    if enPassant == "-" || len(enPassant) < 2 {
        return -1
    }

    fileChar := enPassant[0]
    file := int(fileChar - 'a')

    rankStr := enPassant[1:]
    rank, err := strconv.Atoi(rankStr)
    if err != nil || rank < 1 || rank > 16 {
        return -1
    }

    return Index(file, rank-1, files)
}


type Coords struct{
	X int `json:"x"`
	Y int `json:"y"`
}

type Move struct {
    Piece        rune            `json:"p"`          // Piece type and color (e.g., 'P' for White Pawn, 'p' for Black Pawn)
    From         int             `json:"f"`           // Source bit index (from position)
    To           int             `json:"t"`             // Target bit index (to position)
    Capture      bool            `json:"capture"`        // Whether the move is a capture
    Promotion    string          `json:"promotion,omitempty"` // Promotion piece, if applicable (e.g., 'Q' for Queen)
    ClassicMoveType string       `json:"classic_move_type"` // Classical move type (normal, castle, en_passant, etc.)
    VariantMoveType string      `json:"variant_move_type,omitempty"` // Special move type (e.g., 'duck_move')
    AdditionalData interface{}  `json:"additional_data,omitempty"` // Any extra data (like promotion choice)
}

type variantType string
const (
	Checkmate variantType = "checkmate"
	Antichess variantType = "antichess"
	NCheck variantType = "ncheck"
	DuckChess variantType = "duckchess"
	ArcherChess variantType = "archerchess"
	Wormhole variantType = "wormhole"
)


type classicMoveType uint8

type Piece string

type result uint8

type variantMoveType uint8

type allowedJumpMoveType uint

const (
	Pawn        Piece = "pawn"
	Knight      Piece = "knight"
	Bishop      Piece = "bishop"
	Rook        Piece = "rook"
	Queen       Piece = "queen"
	King        Piece = "king"
	Duck        Piece = "duck"
	CustomPiece Piece = "custom"

	WhiteWins result = iota
	BlackWins
	Stalemate

	NullMove classicMoveType = iota
	CastleMove
	CaptureMove
	QuietMove
	DoublePawnPush
	EnPassant
	PromotionMove

	DuckPlacement variantMoveType = iota
	Teleport

	CaptureOnlyJump allowedJumpMoveType = iota
	QuietOnlyJump
	AllJumps
)





type recentCaptureInfo struct{
	Piece
	squareId int
	moveType classicMoveType
}

type Variant interface {
    makeMove(Move)
    unmakeMove(Move)
    getPseudoLegalMoves(Color, bool) []Move
    GetLegalMoves() []Move
    PerformMove(Move) (result bool, err error)
}

type variant struct {
	*Position
	variantType
	recentCapture recentCaptureInfo
	possibleLegalMoves []Move
	gameResult result
	isGameOverBool bool
}

type CheckmateVariant struct {
    variant
}

func (cv *CheckmateVariant) GetLegalMoves() []Move {return nil}
func (cv *CheckmateVariant) PerformMove(Move) (result bool, err error){ return true,nil}
func (cv *CheckmateVariant) makeMove(Move) {}
func (cv *CheckmateVariant) unmakeMove(Move){}
func (cv *CheckmateVariant)  getPseudoLegalMoves(Color, bool) []Move {return nil}

type GameState struct{
    Variant Variant
    history []Move
}

func newVariant(gameConfig models.GameConfig) (Variant, error) {
	var variantType = Checkmate
	var newVariant Variant
    fen := gameConfig.Position.FEN
	position, err := ParseFEN(fen)
	if err != nil {
		return nil, err
	}
	var variant = variant{
		variantType: variantType,
		Position:    position,
	}
	switch variantType {
	case Checkmate:
		newVariant = &CheckmateVariant{variant}
	case Antichess:
		//newVariant = &AntichessVariant{variant}
	default:
		//newVariant = &CheckmateVariant{variant}
	}

	return newVariant, nil
}

func CreateGame(gameConfig models.GameConfig) (*GameState,error){
    variant , err := newVariant(gameConfig)
    if err!=nil{
        return nil,err
    }
    game := GameState{
        Variant: variant,
        history: make([]Move,0),
    }
    return &game, nil
}
