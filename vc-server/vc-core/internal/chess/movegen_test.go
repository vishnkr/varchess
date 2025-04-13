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

    require.Equal(t, 20,len(pseudoMoves), "Expected at least 20 pseudo-legal moves for White")


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
	expectedBitboard := NewBitboard(8)
	expectedMoves := generateExpectedSlidingMovesUtil(src,offsets,pos.Ranks,pos.Files)
	sort.Ints(expectedMoves)
	for _,sq := range expectedMoves{ 
		ogSrc := SmallToLargeBoardIndex(sq,pos.Files,max(pos.Ranks,pos.Files))
		expectedBitboard.SetBit(ogSrc)
	}
    var moves []Move
	generateSlideMoves2(src,'Q', pos, false, offsets,&moves)
	require.Equal(t, len(expectedMoves),len(moves), "Number of expected moves not matching \nFEN: %s", fen)
}


func TestAttackTable(t *testing.T) {
    ranks, files := 8, 8
	largestDim:= max(ranks,files)
    src := SmallToLargeBoardIndex(27, files,largestDim)

    offsets := queenOffsets
    bb := ComputeSlideAttacks(src, offsets[0],largestDim)
    b1 := ComputeSlideAttacks(src, offsets[1],largestDim)
    b2 := ComputeSlideAttacks(src, offsets[2],largestDim)
    b3 := ComputeSlideAttacks(src, offsets[3],largestDim)
    b4 := ComputeSlideAttacks(src, offsets[4],largestDim)
    b5 := ComputeSlideAttacks(src, offsets[5],largestDim)
    b6 := ComputeSlideAttacks(src, offsets[6],largestDim)
    b7 := ComputeSlideAttacks(src, offsets[7],largestDim)
	EnsureTablesInitialized()
    validSquares := GetValidSquares(8, 8)
	PrintBitboard(validSquares,8,8)
    bb.And(validSquares)
    b1.And(validSquares)
    b2.And(validSquares)
    b3.And(validSquares)
    b4.And(validSquares)
    b5.And(validSquares)
    b6.And(validSquares)
    b7.And(validSquares)

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
