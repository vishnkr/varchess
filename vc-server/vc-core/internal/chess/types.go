package chess

import (
	"fmt"

	"github.com/holiman/uint256"
)

type Move struct {
    Piece        rune                    
    From         int             
    To           int             
    Capture      bool                 
    Promotion    rune
    ClassicMoveType classicMoveType
    VariantMoveType string
    AdditionalData interface{}
}

type MoveJSON struct {
    Piece           string `json:"p"`
    From           int    `json:"f"`
    To             int    `json:"t"`
    Capture        bool   `json:"capture"`
    Promotion      string `json:"promotion,omitempty"`
    ClassicMoveType classicMoveType `json:"classic_move_type"`
    VariantMoveType string `json:"variant_move_type,omitempty"`
    AdditionalData interface{} `json:"additional_data,omitempty"`
}

func (m Move) ToJSON() MoveJSON {
    return MoveJSON{
        Piece:           string(m.Piece),
        From:            m.From,
        To:              m.To,
        Capture:         m.Capture,
        Promotion:       string(m.Promotion),
        ClassicMoveType: m.ClassicMoveType,
        VariantMoveType: m.VariantMoveType,
        AdditionalData:  m.AdditionalData,
    }
}

func (mj MoveJSON) ToMove() (Move, error) {
    if len(mj.Piece) != 1 {
        return Move{}, fmt.Errorf("invalid piece value: expected single character, got %q", mj.Piece)
    }

    return Move{
        Piece:           rune(mj.Piece[0]),
        From:            mj.From,
        To:              mj.To,
        Capture:         mj.Capture,
        Promotion:       rune(mj.Promotion[0]),
        ClassicMoveType: mj.ClassicMoveType,
        VariantMoveType: mj.VariantMoveType,
        AdditionalData:  mj.AdditionalData,
    }, nil
}


type Piece struct {
    Symbol      rune
    IsCustom    bool
    MovePattern *MovePattern
}


type PieceDefinition struct {
    MoveType MoveType
    MovePatterns []Bitboard
}

type MovePattern struct{
	MoveType MoveType
	MoveOffsets []*MoveOffset
}
type MoveOffset struct {
	x, y int
}



type MoveType int
const (
    Slide MoveType = iota 
    Jump
)

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
	mask.Lsh(mask, uint(index))
	temp := b.Bits.Clone()
	return temp.And(b.Bits, mask).Cmp(uint256.NewInt(0)) != 0
}

func (b *Bitboard) GetSetBits() []int {
    var positions []int
    bytes := b.Bits.Bytes32()
    for byteFileRankToIndex, byteVal := range bytes {
        if byteVal == 0 {
            continue
        }
        for bitPos := 0; bitPos < 8; bitPos++ {
            if (byteVal & (1 << bitPos)) != 0 {
                positions = append(positions, (31-byteFileRankToIndex)*8+bitPos)
            }
        }
    }
    return positions
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
    PositionBitBoard *Bitboard
    ColorBitboards map[Color]*Bitboard
	CustomPieceRules map[rune][]MovePattern
}

type Color uint8
const (
    Black Color = iota
    White
)

type Coords struct{
	X int `json:"x"`
	Y int `json:"y"`
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

type PieceRepr string

type result uint8

type variantMoveType uint8

type allowedJumpMoveType uint

const (
	Pawn        PieceRepr = "pawn"
	Knight      PieceRepr = "knight"
	Bishop      PieceRepr = "bishop"
	Rook        PieceRepr = "rook"
	Queen       PieceRepr = "queen"
	King        PieceRepr = "king"
	Duck        PieceRepr = "duck"
	CustomPiece PieceRepr = "custom"

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


type RecentMoveInfo struct {
    Piece          rune
    squareId       int
    moveType       classicMoveType
    capturedPiece  rune
    prevCastling   uint8
	epSquare int
    prevEnPassant  int
    prevHalfMove   int
    prevFullMove   int
}

type GameState struct{
    Variant Variant
    history []Move
}