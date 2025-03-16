package chess

import (
	"fmt"
	"math/bits"
	"strings"

	"github.com/holiman/uint256"
)




type IBitboard interface {
    SetBit(index int)
    ClearBit(index int)
    HasBit(index int) bool
    GetSetBits() []int
}

type Bitboard64 struct {
    Bits uint64
}

func (b *Bitboard64) SetBit(index int) {
    b.Bits |= (1 << index)
}

func (b *Bitboard64) ClearBit(index int) {
    b.Bits &= ^(1 << index)
}

func (b *Bitboard64) HasBit(index int) bool {
    return (b.Bits & (1 << index)) != 0
}

func (b *Bitboard64) GetSetBits() []int {
    var positions []int
    for i := 0; i < 64; i++ {
        if b.HasBit(i) {
            positions = append(positions, i)
        }
    }
    return positions
}

func (b *Bitboard) FirstSetBit() int {
    bytes := b.Bits.Bytes32()
    for i, v := range bytes {
        if v != 0 {
            return (31-i)*8 + bits.TrailingZeros8(v)
        }
    }
    return -1
}



type Bitboard struct {
    Bits *uint256.Int
}
func NewBitboard() *Bitboard {
    return &Bitboard{Bits: uint256.NewInt(0)}
}

func PrintBitboard(bb *Bitboard, ranks, files int) {
	for r := ranks - 1; r >= 0; r-- {
		var row strings.Builder
		for f := 0; f < files; f++ {
            ogSquare := FileRankToLargeIndex(f,r,files)
			if bb.HasBit(ogSquare) {
				row.WriteString("1 ")
			} else {
				row.WriteString("0 ")
			}
		}
		fmt.Println("rn",r," - ",row.String())
	}
	fmt.Println()
}