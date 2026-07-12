package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Position struct {
	Files            int
	Ranks            int
	FEN              string
	LargestDimension int
	WhiteToMove      bool
	Castling         uint8
	EnPassant        int
	HalfMove         int
	FullMove         int
	Pieces           map[rune]Bitboard
	Walls            Bitboard
	PositionBitBoard Bitboard
	ColorBitboards   map[Color]Bitboard
	CustomPieceRules map[rune][]MovePattern
}

// ParseFEN parses a FEN string and returns a Position.
// Internal rank representation: rank 0 = top of board (chess rank 8 / black's back rank).
func ParseFEN(fen string) (*Position, error) {
	parts := strings.Split(fen, " ")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid FEN: expected at least 4 parts")
	}
	board := parts[0]
	sideToMove := parts[1] == "w"
	castlingRights := ParseCastlingRights(parts[2])
	ranks, files, err := InferDimensions(board)
	if err != nil {
		return nil, err
	}
	enPassant := ParseEnPassant(parts[3], ranks, files)
	var largestDimension int = max(ranks, files)
	pos := Position{
		Files:            files,
		Ranks:            ranks,
		FEN:              fen,
		WhiteToMove:      sideToMove,
		Castling:         castlingRights,
		EnPassant:        enPassant,
		LargestDimension: largestDimension,
		Pieces:           make(map[rune]Bitboard),
		Walls:            NewBitboard(largestDimension),
		PositionBitBoard: NewBitboard(largestDimension),
		ColorBitboards:   make(map[Color]Bitboard),
		CustomPieceRules: make(map[rune][]MovePattern),
	}

	pos.ColorBitboards[Black] = NewBitboard(largestDimension)
	pos.ColorBitboards[White] = NewBitboard(largestDimension)

	// rank counts from top of FEN (rank=0 is chess rank 8)
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

		// Internal storage: rank 0 = top of board = chess rank Ranks from bottom.
		ogIndex := FileRankToLargeIndex(file, rank, files, largestDimension)
		if ch == '.' {
			pos.Walls.SetBit(ogIndex)
		} else {
			if _, exists := pos.Pieces[ch]; !exists {
				pos.Pieces[ch] = NewBitboard(largestDimension)
			}
			if unicode.IsLower(ch) {
				pos.ColorBitboards[Black].SetBit(ogIndex)
			} else {
				pos.ColorBitboards[White].SetBit(ogIndex)
			}
			pos.Pieces[ch].SetBit(ogIndex)
			pos.PositionBitBoard.SetBit(ogIndex)
		}

		file++
	}
	return &pos, nil
}

const (
	KSCW uint8 = 3
	QSCW uint8 = 2
	KSCB uint8 = 1
	QSCB uint8 = 0
)

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

func ParseCastlingRights(castling string) uint8 {
	var rights uint8
	if strings.Contains(castling, "K") {
		rights |= 1 << KSCW
	}
	if strings.Contains(castling, "Q") {
		rights |= 1 << QSCW
	}
	if strings.Contains(castling, "k") {
		rights |= 1 << KSCB
	}
	if strings.Contains(castling, "q") {
		rights |= 1 << QSCB
	}
	return rights
}

// ParseEnPassant converts an algebraic en-passant square (e.g. "e3") to an internal index.
// Internal convention: rank 0 = top (chess rank = Ranks from bottom).
// chessRank is 1-indexed from bottom; internalRank = ranks - chessRank.
func ParseEnPassant(enPassant string, ranks, files int) int {
	if enPassant == "-" || len(enPassant) < 2 {
		return -1
	}

	fileChar := enPassant[0]
	file := int(fileChar - 'a')
	if file < 0 || file >= 16 {
		return -1
	}
	rankStr := enPassant[1:]
	chessRank, err := strconv.Atoi(rankStr)
	if err != nil || chessRank < 1 || chessRank > 16 {
		return -1
	}
	// Convert from 1-indexed chess rank (from bottom) to 0-indexed internal rank (from top).
	internalRank := ranks - chessRank
	if internalRank < 0 || internalRank >= ranks || file >= files {
		return -1
	}
	return FileRankToLargeIndex(file, internalRank, files, max(ranks, files))
}

// Square converts an algebraic square name (e.g. "e4") to an internal large-board index.
// Uses the rank-from-top convention: rank 0 = top = chess rank Ranks.
func (p *Position) Square(square string) (int, error) {
	if len(square) < 2 || len(square) > 3 {
		return -1, errors.New("invalid square")
	}

	fileChar := square[0]
	if fileChar < 'a' || fileChar > 'z' {
		return -1, errors.New("invalid square")
	}
	file := int(fileChar - 'a')

	rankStr := square[1:]
	chessRank, err := strconv.Atoi(rankStr)
	if err != nil {
		return -1, errors.New("invalid square")
	}

	internalRank := p.Ranks - chessRank

	if file >= p.Files || internalRank < 0 || internalRank >= p.Ranks {
		return -1, errors.New("invalid square")
	}

	lbd := max(p.Ranks, p.Files)
	return FileRankToLargeIndex(file, internalRank, p.Files, lbd), nil
}
