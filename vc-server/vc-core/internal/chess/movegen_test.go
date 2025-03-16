package chess

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)



func TestMoveGenerationFromStartPos(t *testing.T) {
    pos, err := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
    require.NoError(t, err)

    v := &StandardVariant{variant{Position: pos}}

    pseudoMoves := v.GetPseudoLegalMoves(White, false)
    legalMoves := v.GetLegalMoves()

    require.GreaterOrEqual(t, len(pseudoMoves), 20, "Expected at least 20 pseudo-legal moves for White")


    require.Equal(t, 20, len(legalMoves), "Expected exactly 20 legal moves from starting position")
    
}

func generateExpectedSlidingMovesUtil(src int, offsets []*MoveOffset, ranks, files int) []int {
	x, y := src%files, src/files
    moves := []int{}

    for _, offset := range offsets {
        nx, ny := x, y
        for {
            nx += offset.x
            ny += offset.y

            if nx < 0 || ny < 0 || nx >= files || ny >= ranks {
                break
            }

            moves = append(moves, ny*files+nx)
        }
    }

    return moves
}

func TestGenerateSlideMoves(t *testing.T) {
	fen := "8/8/8/3Q4/8/8/8/8 w - - 0 1"
	pos, err := ParseFEN(fen)
	require.NoError(t, err)

	src := FileRankToIndex(3, 3, pos.Files)
	offsets := queenOffsets
	expectedBitboard := NewBitboard()
	expectedMoves := generateExpectedSlidingMovesUtil(src,offsets,pos.Ranks,pos.Files)
	sort.Ints(expectedMoves)
	for _,sq := range expectedMoves{ 
		ogSrc := SmallToLargeBoardIndex(sq,pos.Files)
		expectedBitboard.SetBit(ogSrc)
	}

	moveBitboard := generateSlideMoves(src, pos, false, offsets)
	require.Equal(t, expectedBitboard.Bits, moveBitboard.Bits, "FEN: %s", fen)
}


func TestAttackTable(t *testing.T) {
    ranks, files := 8, 8
    src := SmallToLargeBoardIndex(27, files)

    offsets := queenOffsets
    bb := ComputeSlideAttacks(src, offsets[0])
    b1 := ComputeSlideAttacks(src, offsets[1])
    b2 := ComputeSlideAttacks(src, offsets[2])
    b3 := ComputeSlideAttacks(src, offsets[3])
    b4 := ComputeSlideAttacks(src, offsets[4])
    b5 := ComputeSlideAttacks(src, offsets[5])
    b6 := ComputeSlideAttacks(src, offsets[6])
    b7 := ComputeSlideAttacks(src, offsets[7])
	EnsureTablesInitialized()
    validSquares := GetValidSquares(8, 8)
	PrintBitboard(validSquares,8,8)
    bb.Bits.And(bb.Bits,validSquares.Bits)
    b1.Bits.And(b1.Bits,validSquares.Bits)
    b2.Bits.And(b2.Bits,validSquares.Bits)
    b3.Bits.And(b3.Bits,validSquares.Bits)
    b4.Bits.And(b4.Bits,validSquares.Bits)
    b5.Bits.And(b5.Bits,validSquares.Bits)
    b6.Bits.And(b6.Bits,validSquares.Bits)
    b7.Bits.And(b7.Bits,validSquares.Bits)

    fmt.Println("bb")
    PrintBitboard(bb, ranks, files)
    fmt.Println("bb1")
    PrintBitboard(b1, ranks, files)
    fmt.Println("bb2")
    PrintBitboard(b2, ranks, files)
	fmt.Println("bb3")
    PrintBitboard(b3, ranks, files)
	fmt.Println("bb4")
    PrintBitboard(b4, ranks, files)
	fmt.Println("bb5")
    PrintBitboard(b5, ranks, files)
	fmt.Println("bb6")
    PrintBitboard(b6, ranks, files)
	fmt.Println("bb7")
    PrintBitboard(b7, ranks, files)
}
