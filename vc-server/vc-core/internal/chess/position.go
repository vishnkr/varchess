package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Position struct {
    Files      int
    Ranks      int
    LargestDimension int
    WhiteToMove bool
    Castling   uint8
    EnPassant  int
    HalfMove   int
    FullMove   int
    Pieces     map[rune]Bitboard
    Walls      Bitboard
    PositionBitBoard Bitboard
    ColorBitboards map[Color]Bitboard
	CustomPieceRules map[rune][]MovePattern
}

func ParseFEN(fen string) (*Position,error) {
    parts := strings.Split(fen, " ")
    board := parts[0]
    sideToMove := parts[1] == "w"
    castlingRights := ParseCastlingRights(parts[2])
    ranks, files,err := InferDimensions(board)
    if err!=nil{
        return nil,err
    }
    enPassant := ParseEnPassant(parts[3],files)
    var largestDimension int = max(ranks,files)
    pos := Position{
        Files:      files,
        Ranks:      ranks,
        WhiteToMove: sideToMove,
        Castling:   castlingRights,
        EnPassant:  enPassant,
        Pieces:     make(map[rune]Bitboard),
        Walls:      NewBitboard(largestDimension),
        PositionBitBoard: NewBitboard(largestDimension),
        ColorBitboards: make(map[Color]Bitboard),

    }

    rank := 0
    file := 0
    pos.ColorBitboards[Black] = NewBitboard(largestDimension)
    pos.ColorBitboards[White] = NewBitboard(largestDimension)
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

        //index := FileRankToIndex(file, rank, files)
        //ogIndex := SmallToLargeBoardIndex(index,files)
        ogIndex := FileRankToLargeIndex(file,rank,files,largestDimension)
        if ch == '.' {
            pos.Walls.SetBit(ogIndex)
        } else {
            if _, exists := pos.Pieces[ch]; !exists {
                pos.Pieces[ch] = NewBitboard(largestDimension)
            }
            if unicode.IsLower(ch){
                pos.ColorBitboards[Black].SetBit(ogIndex)
            } else { pos.ColorBitboards[White].SetBit(ogIndex) }
            pos.Pieces[ch].SetBit(ogIndex)
            pos.PositionBitBoard.SetBit(ogIndex)
        }

        file++
    }
    pos.LargestDimension = largestDimension
    return &pos,nil
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
    if strings.Contains(castling, "K") { rights |= 1 << KSCW }
    if strings.Contains(castling, "Q") { rights |= 1 << QSCW }
    if strings.Contains(castling, "k") { rights |= 1 << KSCB }
    if strings.Contains(castling, "q") { rights |= 1 << QSCB }
    return rights
}


func ParseEnPassant(enPassant string, files int) int {
    if enPassant == "-" || len(enPassant) < 2 {
        return -1
    }

    fileChar := enPassant[0]
    file := int(fileChar - 'a')
    if file>=16{
        return -1
    }
    rankStr := enPassant[1:]
    rank, err := strconv.Atoi(rankStr)
    if err != nil || rank < 1 || rank > 16 {
        return -1
    }

    return FileRankToIndex(file, rank-1, files)
}

func (p *Position) Square(square string) (int, error) {
    if len(square) < 2 || len(square) > 3 {
        return -1, errors.New("invalid square")
    }

    file := square[0] - 'a'
    var rank int

    if len(square) == 2 {
        rank = int(square[1] - '1')
    } else {
        rankValue, err := strconv.Atoi(square[1:])
        if err != nil {
            return -1, errors.New("invalid square")
        }
        rank = rankValue - 1
    }

    if file >= byte(p.Files) || rank >= p.Ranks {
        return -1, errors.New("invalid square")
    }

    return rank*p.Files + int(file), nil
}

