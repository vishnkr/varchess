package chess

import (
	"sync"
	"unicode"
)

type Coordinate struct {
	X int `json:"x" bson:"x"`
	Y int `json:"y" bson:"y"`
}

type JumpProps struct {
	Offset           Coordinate `json:"offset" bson:"offset"`
	IsCaptureAllowed bool       `json:"isCaptureAllowed" bson:"isCaptureAllowed"`
}

type PromoProps struct {
	PromotionSquares []int             `json:"promotionSquares" bson:"promotionSquares"`
	CanPromoteTo     []PromotionOption `json:"canPromoteTo" bson:"canPromoteTo"`
}

type PromotionOption struct {
	PieceType string `json:"pieceType" bson:"pieceType"`
	Name      string `json:"name" bson:"name"`
}

type PieceProps struct {
	SlideOffsets []Coordinate `json:"slideOffsets,omitempty" bson:"slideOffsets,omitempty"`
	JumpProps    []Coordinate `json:"jumpOffsets,omitempty" bson:"jumpOffsets,omitempty"`
	PromoProps   *PromoProps  `json:"promoProps,omitempty" bson:"promoProps,omitempty"`
}

type Dimensions struct {
	Ranks int `json:"ranks" bson:"ranks"`
	Files int `json:"files" bson:"files"`
}

type GameConfig struct {
	VariantType    string                      `json:"variantType" bson:"variantType"`
	Name           string                      `json:"name,omitempty"`
	Dimensions     Dimensions                  `json:"dimensions" bson:"dimensions"`
	FEN            string                      `json:"fen,omitempty" bson:"fen,omitempty"`
	PieceProps     map[string]PieceProps       `json:"pieceProps,omitempty" bson:"pieceProps,omitempty"`
	PieceLocations map[string]map[string][]int `json:"pieceLocations,omitempty" bson:"pieceLocations,omitempty"`
	CustomData     map[string]interface{}      `json:"customData" bson:"customData"`
}

type Game struct {
	ID          string
	Position    Position
	MoveHistory []string
}

func SmallToLargeBoardIndex(sq, smallBoardWidth int, largestDim int) int {
	var largeBoardWidth int = 8
	if largestDim > 8 {
		largeBoardWidth = 16
	}
	rank := sq / smallBoardWidth
	file := sq % smallBoardWidth
	return FileRankToIndex(file, rank, largeBoardWidth)
}

func FileRankToIndex(file, rank, width int) int {
	return rank*width + file
}

func FileRankToLargeIndex(file, rank, smallBoardWidth, largestDim int) int {
	return SmallToLargeBoardIndex(FileRankToIndex(file, rank, smallBoardWidth), smallBoardWidth, largestDim)
}

var SlidingAttackTables256 map[MoveOffset]map[int]Bitboard
var SlidingAttackTables map[MoveOffset]map[int]Bitboard
var RankMasks256 map[int]Bitboard
var FileMasks256 map[int]Bitboard
var RankMasks map[int]Bitboard
var FileMasks map[int]Bitboard

func PrecomputeRankFileMasks() {
	RankMasks256 = make(map[int]Bitboard)
	FileMasks256 = make(map[int]Bitboard)
	RankMasks = make(map[int]Bitboard)
	FileMasks = make(map[int]Bitboard)

	for ranks := 1; ranks <= 16; ranks++ {
		rankMask256 := NewBitboard(16)
		rankMask64 := NewBitboard(8)

		for r := 0; r < ranks; r++ {
			for f := 0; f < 16; f++ {
				rankMask256.SetBit(r*16 + f)
				if r < 8 && f < 8 {
					rankMask64.SetBit(r*8 + f)
				}
			}
		}
		RankMasks256[ranks] = rankMask256
		if ranks <= 8 {
			RankMasks[ranks] = rankMask64
		}
	}

	for files := 1; files <= 16; files++ {
		fileMask256 := NewBitboard(16)
		fileMask64 := NewBitboard(8)

		for f := 0; f < files; f++ {
			for r := 0; r < 16; r++ {
				fileMask256.SetBit(r*16 + f)
				if r < 8 && f < 8 {
					fileMask64.SetBit(r*8 + f)
				}
			}
		}
		FileMasks256[files] = fileMask256
		if files <= 8 {
			FileMasks[files] = fileMask64
		}
	}
}


func PrecomputeSlidingAttacks() {
	SlidingAttackTables = make(map[MoveOffset]map[int]Bitboard)
	SlidingAttackTables256 = make(map[MoveOffset]map[int]Bitboard)
	directions := []*MoveOffset{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, 1}, {-1, 1}, {1, -1}, {-1, -1},
	}

	for _, dir := range directions {
		SlidingAttackTables256[*dir] = make(map[int]Bitboard)
		SlidingAttackTables[*dir] = make(map[int]Bitboard)
		for square := 0; square < 256; square++ {
			if square < 64 {
				SlidingAttackTables[*dir][square] = ComputeSlideAttacks(square, dir, 8)
			}
			SlidingAttackTables256[*dir][square] = ComputeSlideAttacks(square, dir, 16)
		}
	}
}

func ComputeSlideAttacks(src int, dir *MoveOffset, largestDimension int) Bitboard {
	moves := NewBitboard(largestDimension)
	x, y := src%largestDimension, src/largestDimension
	for {
		x += dir.x
		y += dir.y
		if x < 0 || y < 0 || x >= largestDimension || y >= largestDimension {
			break
		}
		moves.SetBit(y*largestDimension + x)
	}
	return moves
}

var once sync.Once

func EnsureTablesInitialized() {
	once.Do(func() {
		PrecomputeRankFileMasks()
		PrecomputeSlidingAttacks()
	})
}

func init() {
	EnsureTablesInitialized()
}

var knightOffsets = []*MoveOffset{
	{-1, -2}, {-1, 2}, {-2, -1}, {-2, 1}, {1, 2}, {1, -2}, {2, 1}, {2, -1},
}
var bishopOffsets = []*MoveOffset{
	{-1, -1}, {-1, 1}, {1, -1}, {1, 1},
}
var rookOffsets = []*MoveOffset{
	{-1, 0}, {0, 1}, {1, 0}, {0, -1},
}

var queenOffsets = append(bishopOffsets, rookOffsets...)

var kingOffsets = []*MoveOffset{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

// GenerateMovesForPiece generates pseudo-legal moves for a piece at src.
// canCaptureKing allows generating moves that land on the opponent's king
// (used when checking if a square is attacked).
func GenerateMovesForPiece(piece rune, src int, position *Position, canCaptureKing bool) []Move {
	var moves []Move
	switch unicode.ToLower(piece) {
	case 'n':
		generateJumpMoves(src, piece, position, canCaptureKing, knightOffsets, &moves)
	case 'b':
		generateSlideMoves2(src, piece, position, canCaptureKing, bishopOffsets, &moves)
	case 'r':
		generateSlideMoves2(src, piece, position, canCaptureKing, rookOffsets, &moves)
	case 'q':
		generateSlideMoves2(src, piece, position, canCaptureKing, queenOffsets, &moves)
	case 'k':
		generateKingMoves(src, position, canCaptureKing, &moves)
	case 'p':
		generatePawnMoves(src, position, canCaptureKing, &moves)
	default:
		if patterns, exists := position.CustomPieceRules[piece]; exists {
			for _, pattern := range patterns {
				if pattern.MoveType == Jump {
					generateJumpMoves(src, piece, position, canCaptureKing, pattern.MoveOffsets, &moves)
				} else if pattern.MoveType == Slide {
					generateSlideMoves2(src, piece, position, canCaptureKing, pattern.MoveOffsets, &moves)
				}
			}
		}
	}
	return moves
}

// IsAttacked returns true if the square sq is attacked by any piece of byColor.
func IsAttacked(sq int, byColor Color, p *Position) bool {
	for piece, bb := range p.Pieces {
		pieceColor := White
		if unicode.IsLower(piece) {
			pieceColor = Black
		}
		if pieceColor != byColor {
			continue
		}
		for _, pos := range bb.GetSetBits() {
			attackMoves := GenerateMovesForPiece(piece, pos, p, true)
			for _, m := range attackMoves {
				if m.To == sq {
					return true
				}
			}
		}
	}
	return false
}

// generateKingMoves generates one-step king moves plus castling pseudo-legals.
func generateKingMoves(src int, p *Position, canCaptureKing bool, moves *[]Move) Bitboard {
	lbd := p.LargestDimension
	x, y := src%lbd, src/lbd
	piece := p.getPieceAt(src)
	if piece == -1 {
		return NewBitboard(lbd)
	}
	color := White
	if unicode.IsLower(piece) {
		color = Black
	}

	for _, offset := range kingOffsets {
		nx, ny := x+offset.x, y+offset.y
		if nx < 0 || ny < 0 || nx >= p.Files || ny >= p.Ranks {
			continue
		}
		newPos := FileRankToLargeIndex(nx, ny, p.Files, lbd)
		if p.Walls.HasBit(newPos) {
			continue
		}
		if p.ColorBitboards[color].HasBit(newPos) {
			continue
		}
		targetPiece := p.getPieceAt(newPos)
		isOpponentKing := targetPiece == 'k' || targetPiece == 'K'
		if isOpponentKing && !canCaptureKing {
			continue
		}
		isCapture := p.PositionBitBoard.HasBit(newPos)
		moveType := QuietMove
		if isCapture {
			moveType = CaptureMove
		}
		*moves = append(*moves, Move{
			Piece: piece, From: src, To: newPos,
			Capture: isCapture, ClassicMoveType: moveType,
		})
	}

	// Castling — only generate when not in the "can capture king" pass (attack detection).
	if !canCaptureKing {
		generateCastlingMoves(src, p, color, moves)
	}
	return NewBitboard(lbd)
}

// generateCastlingMoves appends castling pseudo-legals. Check/intermediate-square
// validation happens in GetLegalMoves.
func generateCastlingMoves(src int, p *Position, color Color, moves *[]Move) {
	lbd := p.LargestDimension
	rank := src / lbd
	kingFile := src % lbd
	piece := rune('K')
	if color == Black {
		piece = 'k'
	}

	kscBit := uint8(1 << KSCW)
	qscBit := uint8(1 << QSCW)
	if color == Black {
		kscBit = 1 << KSCB
		qscBit = 1 << QSCB
	}

	// King-side castling: squares between king and h-file rook must be empty.
	if p.Castling&kscBit != 0 {
		clear := true
		for f := kingFile + 1; f < p.Files-1; f++ {
			if p.PositionBitBoard.HasBit(FileRankToLargeIndex(f, rank, p.Files, lbd)) {
				clear = false
				break
			}
		}
		if clear && kingFile+2 < p.Files {
			castleTo := FileRankToLargeIndex(kingFile+2, rank, p.Files, lbd)
			*moves = append(*moves, Move{
				Piece: piece, From: src, To: castleTo, ClassicMoveType: CastleMove,
			})
		}
	}

	// Queen-side castling: squares between a-file rook and king must be empty.
	if p.Castling&qscBit != 0 {
		clear := true
		for f := 1; f < kingFile; f++ {
			if p.PositionBitBoard.HasBit(FileRankToLargeIndex(f, rank, p.Files, lbd)) {
				clear = false
				break
			}
		}
		if clear && kingFile-2 >= 0 {
			castleTo := FileRankToLargeIndex(kingFile-2, rank, p.Files, lbd)
			*moves = append(*moves, Move{
				Piece: piece, From: src, To: castleTo, ClassicMoveType: CastleMove,
			})
		}
	}
}

// generatePawnMoves generates pawn moves in the rank-from-top convention.
// White pawns (uppercase) move towards rank 0 (direction = -1).
// Black pawns (lowercase) move towards rank Ranks-1 (direction = +1).
func generatePawnMoves(src int, p *Position, canCaptureKing bool, moves *[]Move) {
	piece := p.getPieceAt(src)
	if piece == -1 {
		return
	}
	lbd := p.LargestDimension
	file := src % lbd
	rank := src / lbd

	color := White
	if unicode.IsLower(piece) {
		color = Black
	}
	oppColor := opponent(color)

	// White: direction=-1 (towards rank 0), starts at rank Ranks-2, promotes at rank 0.
	// Black: direction=+1 (towards rank Ranks-1), starts at rank 1, promotes at rank Ranks-1.
	direction := -1
	startRank := p.Ranks - 2
	promotionRank := 0
	promoOptions := []rune{'Q', 'R', 'B', 'N'}
	if color == Black {
		direction = 1
		startRank = 1
		promotionRank = p.Ranks - 1
		promoOptions = []rune{'q', 'r', 'b', 'n'}
	}

	// One step forward.
	newRank := rank + direction
	if newRank >= 0 && newRank < p.Ranks {
		oneStep := FileRankToLargeIndex(file, newRank, p.Files, lbd)
		if !p.PositionBitBoard.HasBit(oneStep) && !p.Walls.HasBit(oneStep) {
			if newRank == promotionRank {
				for _, promo := range promoOptions {
					*moves = append(*moves, Move{
						Piece: piece, From: src, To: oneStep,
						Promotion: promo, ClassicMoveType: PromotionMove,
					})
				}
			} else {
				*moves = append(*moves, Move{
					Piece: piece, From: src, To: oneStep, ClassicMoveType: QuietMove,
				})
				// Double push from starting rank.
				if rank == startRank {
					twoStepRank := newRank + direction
					if twoStepRank >= 0 && twoStepRank < p.Ranks {
						twoStep := FileRankToLargeIndex(file, twoStepRank, p.Files, lbd)
						if !p.PositionBitBoard.HasBit(twoStep) && !p.Walls.HasBit(twoStep) {
							*moves = append(*moves, Move{
								Piece: piece, From: src, To: twoStep, ClassicMoveType: DoublePawnPush,
							})
						}
					}
				}
			}
		}
	}

	// Diagonal captures.
	captureRank := rank + direction
	if captureRank >= 0 && captureRank < p.Ranks {
		for _, fileOffset := range []int{-1, 1} {
			captureFile := file + fileOffset
			if captureFile < 0 || captureFile >= p.Files {
				continue
			}
			capturePos := FileRankToLargeIndex(captureFile, captureRank, p.Files, lbd)
			if p.Walls.HasBit(capturePos) {
				continue
			}

			hasEnemy := p.ColorBitboards[oppColor].HasBit(capturePos)
			// When canCaptureKing is true we also generate moves that land on the king
			// (used purely for attack detection, not for game play).
			if !hasEnemy {
				if canCaptureKing {
					oppKing := rune('k')
					if oppColor == White {
						oppKing = 'K'
					}
					if kingBB, ok := p.Pieces[oppKing]; ok {
						hasEnemy = kingBB.HasBit(capturePos)
					}
				}
			}

			if hasEnemy {
				if captureRank == promotionRank {
					for _, promo := range promoOptions {
						*moves = append(*moves, Move{
							Piece: piece, From: src, To: capturePos,
							Capture: true, Promotion: promo, ClassicMoveType: PromotionMove,
						})
					}
				} else {
					*moves = append(*moves, Move{
						Piece: piece, From: src, To: capturePos,
						Capture: true, ClassicMoveType: CaptureMove,
					})
				}
			}
		}
	}

	// En passant.
	if p.EnPassant != -1 {
		for _, fileOffset := range []int{-1, 1} {
			captureFile := file + fileOffset
			if captureFile < 0 || captureFile >= p.Files {
				continue
			}
			targetSquare := FileRankToLargeIndex(captureFile, rank+direction, p.Files, lbd)
			if targetSquare == p.EnPassant {
				*moves = append(*moves, Move{
					Piece: piece, From: src, To: targetSquare,
					Capture: true, ClassicMoveType: EnPassant,
				})
			}
		}
	}
}

// getPieceAt returns the piece symbol at the given large-board index, or -1 if empty.
func (p *Position) getPieceAt(src int) rune {
	for piece, bb := range p.Pieces {
		if bb.HasBit(src) {
			return piece
		}
	}
	return -1
}

// generateSlideMoves2 generates sliding moves using direct coordinate iteration.
// Uses p.LargestDimension for correct index extraction on any board size.
func generateSlideMoves2(src int, piece rune, p *Position, canCaptureKing bool, offsets []*MoveOffset, moves *[]Move) {
	lbd := p.LargestDimension
	x, y := src%lbd, src/lbd

	srcColor := White
	if unicode.IsLower(piece) {
		srcColor = Black
	}

	for _, offset := range offsets {
		nx, ny := x, y
		for {
			nx += offset.x
			ny += offset.y

			if nx < 0 || ny < 0 || nx >= p.Files || ny >= p.Ranks {
				break
			}

			newPos := FileRankToLargeIndex(nx, ny, p.Files, lbd)
			if p.Walls.HasBit(newPos) {
				break
			}

			isOccupied := p.PositionBitBoard.HasBit(newPos)
			if !isOccupied {
				*moves = append(*moves, Move{
					Piece: piece, From: src, To: newPos, ClassicMoveType: QuietMove,
				})
				continue
			}

			// Occupied square.
			if p.ColorBitboards[srcColor].HasBit(newPos) {
				break // own piece blocks
			}

			targetPiece := p.getPieceAt(newPos)
			isOpponentKing := targetPiece == 'k' || targetPiece == 'K'
			if isOpponentKing && !canCaptureKing {
				break
			}

			*moves = append(*moves, Move{
				Piece: piece, From: src, To: newPos,
				Capture: true, ClassicMoveType: CaptureMove,
			})
			break
		}
	}
}

// generateJumpMoves generates jump (non-sliding) moves.
func generateJumpMoves(src int, piece rune, p *Position, canCaptureKing bool, offsets []*MoveOffset, moves *[]Move) {
	lbd := p.LargestDimension
	x, y := src%lbd, src/lbd

	srcColor := White
	if unicode.IsLower(piece) {
		srcColor = Black
	}

	for _, offset := range offsets {
		nx, ny := x+offset.x, y+offset.y
		if nx < 0 || ny < 0 || nx >= p.Files || ny >= p.Ranks {
			continue
		}

		newPos := FileRankToLargeIndex(nx, ny, p.Files, lbd)
		if p.Walls.HasBit(newPos) {
			continue
		}
		if p.ColorBitboards[srcColor].HasBit(newPos) {
			continue
		}

		targetPiece := p.getPieceAt(newPos)
		isOpponentKing := targetPiece == 'k' || targetPiece == 'K'
		if isOpponentKing && !canCaptureKing {
			continue
		}

		isCapture := p.PositionBitBoard.HasBit(newPos)
		moveType := QuietMove
		if isCapture {
			moveType = CaptureMove
		}
		*moves = append(*moves, Move{
			Piece: piece, From: src, To: newPos,
			Capture: isCapture, ClassicMoveType: moveType,
		})
	}
}

// flattenBitboard converts a move bitboard to a Move slice (used by legacy tests).
func flattenBitboard(src int, position *Position, moveBitboard Bitboard, opponentBitboard Bitboard, piece rune) []Move {
	var moves []Move
	setBits := moveBitboard.GetSetBits()

	for _, pos := range setBits {
		isCapture := opponentBitboard.HasBit(pos)
		move := Move{
			Piece:           piece,
			From:            src,
			To:              pos,
			Capture:         isCapture,
			ClassicMoveType: QuietMove,
		}
		if isCapture {
			move.ClassicMoveType = CaptureMove
		}
		moves = append(moves, move)
	}
	return moves
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GetValidSquares returns a bitboard with bits set for all squares within
// the given ranks×files board dimensions. Used for display/masking.
func GetValidSquares(ranks, files int) Bitboard {
	lbd := max(ranks, files)
	bb := NewBitboard(lbd)
	for r := 0; r < ranks; r++ {
		for f := 0; f < files; f++ {
			bb.SetBit(FileRankToLargeIndex(f, r, files, lbd))
		}
	}
	return bb
}
