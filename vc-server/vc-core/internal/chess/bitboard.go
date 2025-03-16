package chess

import (
	"fmt"
	"math/bits"
	"strings"

	"github.com/holiman/uint256"
)




type Bitboard interface {
    SetBit(index int)
    ClearBit(index int)
    HasBit(index int) bool
    GetSetBits() []int
    FirstSetBit() int

    Or(others ...Bitboard)
    And(others ...Bitboard)
    Xor(others ...Bitboard)
    Not()
    Lsh(n uint)
    Rsh(n uint)
    Clone () Bitboard
    Neg() Bitboard
    Sub(other Bitboard) 
    One() Bitboard
    Equal(other Bitboard) bool
}

func GetAttackTables(lbd int) map[MoveOffset]map[int]Bitboard{
    EnsureTablesInitialized()
    if lbd > 8{
        return SlidingAttackTables256
    }
    return SlidingAttackTables
}

func GetRankFileMasks(ranks, files int) (Bitboard,Bitboard){
    EnsureTablesInitialized()
    lbd:= max(ranks,files)

    if lbd>8{
        return RankMasks256[ranks],FileMasks256[files]
    }
    return RankMasks[ranks],FileMasks[files]

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

func (b *Bitboard64) FirstSetBit() int {
    if b.Bits == 0 {
        return -1
    }
    return bits.TrailingZeros64(b.Bits)
}

func (b *Bitboard64) Or(others ...Bitboard) {
    for _, other := range others {
        if bb, ok := other.(*Bitboard64); ok {
            b.Bits |= bb.Bits
        }
    }
}

func (b *Bitboard64) Xor(others ...Bitboard) {
    for _, other := range others {
        if bb, ok := other.(*Bitboard64); ok {
            b.Bits ^= bb.Bits
        }
    }
}


func (b *Bitboard64) And(others ...Bitboard) {
    for _, other := range others {
        if bb, ok := other.(*Bitboard64); ok {
            b.Bits &= bb.Bits
        }
    }
}

func (b *Bitboard64) One() Bitboard {
    return &Bitboard64{Bits: 1}
}

func (b *Bitboard64) Neg() Bitboard {
    return &Bitboard64{Bits:-b.Bits}
}

func (b *Bitboard64) Sub(other Bitboard) {
    if o, ok := other.(*Bitboard64); ok {
        b.Bits -= o.Bits
    }
}

func (b *Bitboard64) Equal(other Bitboard) bool{
    if o, ok := other.(*Bitboard64); ok {
        return b.Bits == o.Bits
    }
    return false
}

func (b *Bitboard64) Not() {
    b.Bits = ^b.Bits
}

func (b *Bitboard64) Lsh(n uint) {
    b.Bits <<= n
}

func (b *Bitboard64) Rsh(n uint) {
    b.Bits >>= n
}

func (b *Bitboard64) Clone() Bitboard{
    return &Bitboard64{b.Bits} 
}

type Bitboard256 struct {
    Bits *uint256.Int
}

func NewBitboard(largestDim int) Bitboard {
    if largestDim <= 8 {
        return &Bitboard64{Bits: 0}
    }
    return &Bitboard256{Bits: uint256.NewInt(0)}
}

func (b *Bitboard256) SetBit(index int) {
	mask := uint256.NewInt(0).SetUint64(1)
	mask.Lsh(mask, uint(index)) // mask = 1 << index
	b.Bits.Or(b.Bits, mask) // b.Bits |= mask
}


func (b *Bitboard256) ClearBit(index int) {
	mask := uint256.NewInt(0).SetUint64(1)
	mask.Lsh(mask, uint(index)) // mask = 1 << index
	mask.Not(mask)              // mask = ~mask
	b.Bits.And(b.Bits, mask)    // b.Bits &= mask
}

func (b *Bitboard256) HasBit(index int) bool {
	mask := uint256.NewInt(0).SetUint64(1)
	mask.Lsh(mask, uint(index))
	temp := b.Bits.Clone()
	return temp.And(b.Bits, mask).Cmp(uint256.NewInt(0)) != 0
}

func (b *Bitboard256) GetSetBits() []int {
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

func (b *Bitboard256) FirstSetBit() int {
    bytes := b.Bits.Bytes32()
    for i, v := range bytes {
        if v != 0 {
            return (31-i)*8 + bits.TrailingZeros8(v)
        }
    }
    return -1
}

func (b *Bitboard256) And(others ...Bitboard) {
    for _, other := range others {
        if bb, ok := other.(*Bitboard256); ok {
            b.Bits.And(b.Bits, bb.Bits)
        }
    }
}

func (b *Bitboard256) Or(others ...Bitboard) {
    for _, other := range others {
        if bb, ok := other.(*Bitboard256); ok {
            b.Bits.Or(b.Bits, bb.Bits)
        }
    }
}

func (b *Bitboard256) Xor(others ...Bitboard) {
    for _, other := range others {
        if bb, ok := other.(*Bitboard256); ok {
            b.Bits.Xor(b.Bits, bb.Bits)
        }
    }
}


func (b *Bitboard256) Neg() Bitboard{
    return &Bitboard256{Bits: b.Bits.Neg(b.Bits)}
}

func (b *Bitboard256) Sub(other Bitboard) {
    if o, ok := other.(*Bitboard256); ok {
        b.Bits.Sub(b.Bits, o.Bits)
    } 
}

func (b *Bitboard256) Not() {
    b.Bits.Not(b.Bits)
}

func (b *Bitboard256) Lsh(n uint) {
    b.Bits.Lsh(b.Bits, n)
}

func (b *Bitboard256) Rsh(n uint) {
    b.Bits.Rsh(b.Bits, n)
}

func (b *Bitboard256) Clone() Bitboard{
    return &Bitboard256{Bits: b.Bits.Clone()}
}

func (b *Bitboard256) One() Bitboard {
    return &Bitboard256{Bits: uint256.NewInt(1)}
}

func (b *Bitboard256) Equal(other Bitboard) bool{
    if o, ok := other.(*Bitboard256); ok {
        return b.Bits.Eq(o.Bits)
    }
    return false
}


func PrintBitboard(bb Bitboard, ranks, files int) {
	for r := ranks - 1; r >= 0; r-- {
		var row strings.Builder
		for f := 0; f < files; f++ {
            ogSquare := FileRankToLargeIndex(f,r,files,max(ranks,files))
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