package chess

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitboardOperations(t *testing.T) {
	bb := NewBitboard(16)

	bb.SetBit(10)
	require.True(t, bb.HasBit(10), "Bit should be set")

	bb.ClearBit(10)
	require.False(t, bb.HasBit(10), "Bit should be cleared")

	bb.SetBit(134)
	bb.SetBit(77)
	setBits := bb.GetSetBits()
	require.ElementsMatch(t, []int{134, 77}, setBits, "GetSetBits should return correct positions")
}

func TestInferDimensions(t *testing.T) {
	tests := []struct {
		fen    string
		ranks  int
		files  int
		hasErr bool
	}{
		{"8/8/8/8/8/8/8/8", 8, 8, false},
		{"5/5/5", 3, 5, false},
		{"4/4/44/4", 0, 0, true},
	}

	for _, tt := range tests {
		r, f, err := InferDimensions(tt.fen)
		if tt.hasErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.ranks, r)
			require.Equal(t, tt.files, f)
		}
	}
}

func TestParseStandardFEN(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	pos, err := ParseFEN(fen)
	require.NoError(t, err)
	require.Equal(t, 8, pos.Ranks)
	require.Equal(t, 8, pos.Files)
	require.True(t, pos.WhiteToMove)
	require.Equal(t, uint8(15), pos.Castling)
	require.Equal(t, -1, pos.EnPassant)
}

func TestParseCastlingRights(t *testing.T) {
	require.Equal(t, uint8(15), ParseCastlingRights("KQkq"))
	require.Equal(t, uint8(12), ParseCastlingRights("KQ"))
	require.Equal(t, uint8(3), ParseCastlingRights("kq"))
	require.Equal(t, uint8(0), ParseCastlingRights("-"))
}

// TestParseEnPassant: uses rank-from-top convention.
// "e3" = chess rank 3, internal rank = 8-3=5, file e=4 → index 5*8+4=44.
// "d6" = chess rank 6, internal rank = 8-6=2, file d=3 → index 2*8+3=19.
func TestParseEnPassant(t *testing.T) {
	require.Equal(t, FileRankToIndex(4, 5, 8), ParseEnPassant("e3", 8, 8))
	require.Equal(t, FileRankToIndex(3, 2, 8), ParseEnPassant("d6", 8, 8))
	require.Equal(t, -1, ParseEnPassant("-", 8, 8))
	require.Equal(t, -1, ParseEnPassant("z9", 8, 8))
}

// TestSquare: uses rank-from-top convention.
// a1 = chess rank 1, internal rank 7, file 0 → 7*8+0 = 56.
// h8 = chess rank 8, internal rank 0, file 7 → 0*8+7 = 7.
// e4 = chess rank 4, internal rank 4, file 4 → 4*8+4 = 36.
func TestSquare(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		square   string
		expected int
		wantErr  error
	}{
		{"Valid A1 (8x8)", Position{Files: 8, Ranks: 8}, "a1", 56, nil},
		{"Valid H8 (8x8)", Position{Files: 8, Ranks: 8}, "h8", 7, nil},
		{"Valid E4 (8x8)", Position{Files: 8, Ranks: 8}, "e4", 36, nil},

		// Custom Board Sizes
		// c3 in 5x5: internal rank=2, file c=2. Large board uses 8-wide grid for dim≤8:
		// FileRankToLargeIndex(2,2,5,5) = SmallToLargeBoardIndex(2*5+2=12, 5, 5)
		//   = FileRankToIndex(12%5=2, 12/5=2, 8) = 2*8+2 = 18
		{"Valid C3 (5x5)", Position{Files: 5, Ranks: 5}, "c3", 18, nil},
		// j10 in 10x12: internal rank = 12-10=2, file j=9 → index via FileRankToLargeIndex(9,2,10,12) = SmallToLargeBoardIndex(2*10+9=29, 10, 12) = FileRankToIndex(9,2,16) = 2*16+9=41
		{"Valid J10 (10x12)", Position{Files: 10, Ranks: 12}, "j10", 41, nil},

		// Invalid Format
		{"Invalid Empty", Position{Files: 8, Ranks: 8}, "", -1, errors.New("invalid square")},
		{"Invalid Too Long", Position{Files: 8, Ranks: 8}, "abc", -1, errors.New("invalid square")},
		{"Invalid Number First", Position{Files: 8, Ranks: 8}, "1a", -1, errors.New("invalid square")},

		// Out of Bounds
		{"Out of Bounds File", Position{Files: 8, Ranks: 8}, "z1", -1, errors.New("invalid square")},
		{"Out of Bounds Rank", Position{Files: 8, Ranks: 8}, "a9", -1, errors.New("invalid square")},
		{"Out of Bounds Both", Position{Files: 8, Ranks: 8}, "z99", -1, errors.New("invalid square")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.position.Square(tt.square)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestGenerateMoveCountForKnight(t *testing.T) {
	tests := []struct {
		fen        string
		file, rank int
		expected   int
		piece      rune
	}{
		{"4rK2/2p1R1n1/5p2/5N1k/2p1pp2/2Q5/8/8 w - - 0 1", 6, 1, 2, 'n'},
		{"8/8/8/3N4/8/8/8/8 w - - 0 1", 3, 3, 8, 'N'},
		{"8/8/8/8/8/N7/8/8 w - - 0 1", 0, 5, 4, 'N'},
		{"N7/8/8/8/8/8/8/8 w - - 0 1", 0, 0, 2, 'N'},
		{"8/2P1P3/1P3P2/3N4/1P2P3/2P2P2/8/8 w - - 0 1", 3, 3, 2, 'N'},
		{"8/8/8/3n4/2p1p3/8/8/8 w - - 0 1", 3, 3, 8, 'n'},
		{"8/2p1R3/1K3p2/3n4/2p1pp2/2Q1N3/8/8 w - - 0 1", 3, 3, 4, 'n'},
		{"5/5/2N2/5/5 w - - 0 1", 2, 2, 8, 'N'},
		{"10/10/10/10/10/5n4/10/10/10/10 w - - 0 1", 5, 5, 8, 'n'},
	}

	for _, tt := range tests {
		pos, err := ParseFEN(tt.fen)
		require.NoError(t, err)
		moves := GenerateMovesForPiece(tt.piece, FileRankToLargeIndex(tt.file, tt.rank, pos.Files, pos.LargestDimension), pos, false)
		require.Len(t, moves, tt.expected, "FEN: %s", tt.fen)
	}
}

func TestGenerateMoveCountForBishop(t *testing.T) {
	tests := []struct {
		fen        string
		file, rank int
		expected   int
		piece      rune
	}{
		{"8/8/8/3B4/8/8/8/8 w - - 0 1", 3, 3, 13, 'B'},
		{"8/8/8/8/8/B7/8/8 w - - 0 1", 0, 5, 7, 'B'},
		{"8/2P1P3/1P3P2/3B4/1P2P3/2P2P2/8/8 w - - 0 1", 3, 3, 9, 'B'},
		{"6k1/1n6/8/3b4/8/1Q6/8/7K w - - 0 1", 3, 3, 8, 'b'},
		{"8/2p1P3/1K3p2/3B4/2p1pp2/2Q1N3/8/8 w - - 0 1", 3, 3, 8, 'B'},
	}

	for _, tt := range tests {
		pos, err := ParseFEN(tt.fen)
		require.NoError(t, err)
		moves := GenerateMovesForPiece(tt.piece, FileRankToLargeIndex(tt.file, tt.rank, pos.Files, pos.LargestDimension), pos, false)
		require.Len(t, moves, tt.expected, "FEN: %s", tt.fen)
	}
}

func TestGenerateMoveCountForRook(t *testing.T) {
	tests := []struct {
		fen        string
		file, rank int
		expected   int
		piece      rune
	}{
		{"8/8/8/3R4/8/8/8/8 w - - 0 1", 3, 3, 14, 'R'},
		{"8/8/8/8/8/R7/8/8 w - - 0 1", 0, 5, 14, 'R'},
		{"8/3PP3/2Pk1P2/2KRB3/3P4/3P1P2/8/3P4 w - - 0 1", 3, 3, 0, 'R'},
		{"3r4/8/8/k2rQ3/8/8/8/3K4 w - - 0 1", 3, 3, 8, 'r'},
		{"8/2p1R3/1K3p2/3R4/2p1pp2/2Q1N3/8/8 w - - 0 1", 4, 1, 9, 'R'},
		{"8/8/3K4/3r4/8/8/8/8 w - - 0 1", 3, 3, 11, 'r'},
	}

	for _, tt := range tests {
		pos, err := ParseFEN(tt.fen)
		require.NoError(t, err)
		moves := GenerateMovesForPiece(tt.piece, FileRankToLargeIndex(tt.file, tt.rank, pos.Files, pos.LargestDimension), pos, false)
		require.Len(t, moves, tt.expected, "FEN: %s", tt.fen)
	}
}

func TestFlattenBitboard(t *testing.T) {
	pos := &Position{
		Files: 8,
		Ranks: 8,
		Pieces: make(map[rune]Bitboard),
		ColorBitboards: map[Color]Bitboard{
			White: NewBitboard(8),
			Black: NewBitboard(8),
		},
		PositionBitBoard: NewBitboard(8),
		Walls:            NewBitboard(8),
	}

	src := FileRankToIndex(3, 3, 8) // d5 (rank-from-top: rank 3 from top)
	moveBitboard := NewBitboard(8)
	moveBitboard.SetBit(FileRankToIndex(5, 4, 8)) // f4
	moveBitboard.SetBit(FileRankToIndex(2, 5, 8)) // c3

	moves := flattenBitboard(src, pos, moveBitboard, pos.ColorBitboards[Black], 'N')
	require.Len(t, moves, 2)
	require.Equal(t, moves[0].From, src)
	require.Equal(t, moves[1].From, src)
}
