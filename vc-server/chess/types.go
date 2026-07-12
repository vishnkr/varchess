package chess

import (
	"fmt"
)

type Move struct {
	Piece           rune
	From            int
	To              int
	Capture         bool
	Promotion       rune
	ClassicMoveType classicMoveType
	VariantMoveType string
	AdditionalData  interface{}
}

type MoveJSON struct {
	Piece           string          `json:"piece"`
	From            int             `json:"from"`
	To              int             `json:"to"`
	Capture         bool            `json:"capture,omitempty"`
	Promotion       string          `json:"promotion,omitempty"`
	ClassicMoveType classicMoveType `json:"classic_move_type,omitempty"`
	VariantMoveType string          `json:"variant_move_type,omitempty"`
	AdditionalData  interface{}     `json:"additional_data,omitempty"`
}

func (m Move) ToJSON() MoveJSON {
	promo := ""
	if m.Promotion != 0 {
		promo = string(m.Promotion)
	}
	return MoveJSON{
		Piece:           string(m.Piece),
		From:            m.From,
		To:              m.To,
		Capture:         m.Capture,
		Promotion:       promo,
		ClassicMoveType: m.ClassicMoveType,
		VariantMoveType: m.VariantMoveType,
		AdditionalData:  m.AdditionalData,
	}
}

func (mj MoveJSON) ToMove() (Move, error) {
	if len(mj.Piece) != 1 {
		return Move{}, fmt.Errorf("invalid piece value: expected single character, got %q", mj.Piece)
	}
	var promo rune
	if len(mj.Promotion) == 1 {
		promo = rune(mj.Promotion[0])
	}
	return Move{
		Piece:           rune(mj.Piece[0]),
		From:            mj.From,
		To:              mj.To,
		Capture:         mj.Capture,
		Promotion:       promo,
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
	MoveType     MoveType
	MovePatterns []Bitboard
}

type MovePattern struct {
	MoveType    MoveType
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

type Color uint8

const (
	Black Color = iota
	White
)

type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type classicMoveType uint8

type PieceRepr string

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

	NullMove       classicMoveType = iota
	CastleMove     classicMoveType = iota
	CaptureMove    classicMoveType = iota
	QuietMove      classicMoveType = iota
	DoublePawnPush classicMoveType = iota
	EnPassant      classicMoveType = iota
	PromotionMove  classicMoveType = iota

	DuckPlacement variantMoveType = iota
	Teleport      variantMoveType = iota

	CaptureOnlyJump allowedJumpMoveType = iota
	QuietOnlyJump   allowedJumpMoveType = iota
	AllJumps        allowedJumpMoveType = iota
)

// RecentMoveInfo stores enough data to fully undo a move.
type RecentMoveInfo struct {
	Piece          rune
	moveType       classicMoveType
	capturedPiece  rune
	capturedSquare int // actual square of captured piece (differs from To for en passant)
	prevCastling   uint8
	prevEnPassant  int
	prevHalfMove   int
	prevFullMove   int
	castleRookFrom int
	castleRookTo   int
}

// MoveResult is returned by PerformMove.
type MoveResult struct {
	IsCapture  bool
	IsCheck    bool
	IsGameOver bool
	Result     *EngineResult
}
