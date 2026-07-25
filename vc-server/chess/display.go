package chess

import (
	"fmt"
	"strings"
	"unicode"
)

// SquareToAlgebraic converts an internal large-board index back to an algebraic
// square name (e.g. 0 on an 8x8 board → "a8").
// Convention: rank 0 is the top of the board (chess rank = Ranks).
func SquareToAlgebraic(sq, files, ranks int) string {
	lbd := max(ranks, files)
	file, internalRank := FileRankFromIndex(sq, lbd)
	if file >= files || internalRank >= ranks {
		return "??"
	}
	// internalRank 0 = chess rank Ranks (top), so chess rank = Ranks - internalRank.
	chessRank := ranks - internalRank
	return fmt.Sprintf("%c%d", rune('a'+file), chessRank)
}

// FormatMove renders a move as "e2e4", "e7e8Q" (promotion), or "O-O" / "O-O-O"
// (castling).
func FormatMove(m Move, pos *Position) string {
	if m.ClassicMoveType == CastleMove {
		// King-side or queen-side based on direction.
		fromFile, _ := FileRankFromIndex(m.From, pos.LargestDimension)
		toFile, _ := FileRankFromIndex(m.To, pos.LargestDimension)
		if toFile > fromFile {
			return "O-O"
		}
		return "O-O-O"
	}
	from := SquareToAlgebraic(m.From, pos.Files, pos.Ranks)
	to := SquareToAlgebraic(m.To, pos.Files, pos.Ranks)
	s := from + to
	if m.Promotion != 0 {
		s += string(unicode.ToUpper(m.Promotion))
	}
	return s
}

// FormatLegalMoves returns a numbered list of legal moves for display.
func FormatLegalMoves(moves []Move, pos *Position) string {
	var sb strings.Builder
	for i, m := range moves {
		if i > 0 && i%5 == 0 {
			sb.WriteRune('\n')
		}
		sb.WriteString(fmt.Sprintf("%3d) %-8s", i+1, FormatMove(m, pos)))
	}
	return sb.String()
}

// PositionToASCII renders the position as a human-readable board string.
// Uppercase = White, lowercase = Black. Rank 8 is at the top.
func PositionToASCII(pos *Position) string {
	// Build a 2D square map from the bitboards.
	board := make([]rune, pos.Ranks*pos.Files)
	for i := range board {
		board[i] = '.'
	}

	for piece, bb := range pos.Pieces {
		for _, idx := range bb.GetSetBits() {
			// idx is in large-board space; convert to small-board (ranks×files).
			file, internalRank := FileRankFromIndex(idx, pos.LargestDimension)
			if internalRank < pos.Ranks && file < pos.Files {
				board[internalRank*pos.Files+file] = piece
			}
		}
	}

	var sb strings.Builder
	// Column labels.
	sb.WriteString("   ")
	for f := 0; f < pos.Files; f++ {
		sb.WriteString(fmt.Sprintf(" %c ", rune('a'+f)))
	}
	sb.WriteRune('\n')
	sb.WriteString("   ")
	sb.WriteString(strings.Repeat("---", pos.Files))
	sb.WriteRune('\n')

	for r := 0; r < pos.Ranks; r++ {
		chessRank := pos.Ranks - r
		sb.WriteString(fmt.Sprintf("%2d |", chessRank))
		for f := 0; f < pos.Files; f++ {
			sq := r*pos.Files + f
			sb.WriteString(fmt.Sprintf(" %c ", board[sq]))
		}
		sb.WriteString(fmt.Sprintf("| %d\n", chessRank))
	}

	sb.WriteString("   ")
	sb.WriteString(strings.Repeat("---", pos.Files))
	sb.WriteRune('\n')
	sb.WriteString("   ")
	for f := 0; f < pos.Files; f++ {
		sb.WriteString(fmt.Sprintf(" %c ", rune('a'+f)))
	}
	sb.WriteRune('\n')
	return sb.String()
}

// ParseAlgebraicMove parses a move string like "e2e4" or "e7e8Q" into square
// indices using the position's coordinate system. Returns the from/to indices
// and an optional promotion rune.
func ParseAlgebraicMove(s string, pos *Position) (from, to int, promotion rune, err error) {
	s = strings.TrimSpace(s)
	if len(s) < 4 {
		return 0, 0, 0, fmt.Errorf("move too short: %q", s)
	}
	fromSq := s[:2]
	toSq := s[2:4]
	from, err = pos.Square(fromSq)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid from-square %q: %w", fromSq, err)
	}
	to, err = pos.Square(toSq)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid to-square %q: %w", toSq, err)
	}
	if len(s) >= 5 {
		r := rune(unicode.ToLower(rune(s[4])))
		promotion = r
	}
	return from, to, promotion, nil
}

// MatchMove finds the legal move that matches from/to/promotion among legalMoves.
// Needed because Move contains ClassicMoveType which the user doesn't know.
func MatchMove(from, to int, promotion rune, legalMoves []Move) (Move, bool) {
	for _, m := range legalMoves {
		if m.From == from && m.To == to {
			if promotion == 0 || m.Promotion == promotion {
				return m, true
			}
		}
	}
	return Move{}, false
}
