package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"vc-server/vc-core/internal/models"

	"github.com/holiman/uint256"
)



type Game struct {
    ID        string
    Position  Position
    MoveHistory []string
}


func SmallToLargeBoardIndex(sq, smallBoardWidth int) int {
    rank := sq / smallBoardWidth
    file := sq % smallBoardWidth
    return FileRankToIndex(file, rank, 16)
}


func FileRankToIndex(file, rank, width int) int {
    return rank*width + file
}

func FileRankToLargeIndex(file,rank,sbm int) int{
    return SmallToLargeBoardIndex(FileRankToIndex(file,rank,sbm),sbm)
}

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



var SlidingAttackTables map[MoveOffset]map[int]*Bitboard

var RankMasks map[int]*Bitboard
var FileMasks map[int]*Bitboard

func PrecomputeRankFileMasks() {
    RankMasks = make(map[int]*Bitboard)
    FileMasks = make(map[int]*Bitboard)

    for ranks := 1; ranks <= 16; ranks++ {
        mask := NewBitboard()
        for r := 0; r < ranks; r++ {
            for f := 0; f < 16; f++ {
                mask.SetBit(r*16 + f)
            }
        }
        RankMasks[ranks] = mask
    }

    for files := 1; files <= 16; files++ {
        mask := NewBitboard()
        for r := 0; r < 16; r++ {
            for f := 0; f < files; f++ {
                mask.SetBit(r*16 + f)
            }
        }
        FileMasks[files] = mask
    }
}

func PrecomputeSlidingAttacks() {
    SlidingAttackTables = make(map[MoveOffset]map[int]*Bitboard)

    directions := []*MoveOffset{
        {1, 0}, {-1, 0}, {0, 1}, {0, -1}, 
        {1, 1}, {-1, 1}, {1, -1}, {-1, -1},
    }

    for _, dir := range directions {
        SlidingAttackTables[*dir] = make(map[int]*Bitboard)
        for square := 0; square < 256; square++ {
            SlidingAttackTables[*dir][square] = ComputeSlideAttacks(square, dir)
        }
    }
}

func ComputeSlideAttacks(src int, dir *MoveOffset) *Bitboard {
    moves := NewBitboard()
    x, y := src%16, src/16
    for {
        x += dir.x
        y += dir.y

        if x < 0 || y < 0 || x >= 16 || y >= 16 {
            break
        }
        moves.SetBit(y*16 + x)
    }
    return moves
}

func GetValidSquares(ranks, files int) *Bitboard {
    bb := NewBitboard()
    bb.Bits.And(RankMasks[ranks].Bits,FileMasks[files].Bits)
    PrintBitboard(bb,8,8)
    PrintBitboard(bb,16,16)
    return bb
}

var once sync.Once
func EnsureTablesInitialized() {
    once.Do(func (){
        PrecomputeRankFileMasks()
        PrecomputeSlidingAttacks()
    })
}
func init(){
    EnsureTablesInitialized()
}

var knightOffsets = []*MoveOffset{
    {-1, -2}, {-1, 2}, {-2, -1}, {-2, 1}, {1,2}, {1,-2}, {2,1}, {2,-1},
}
var bishopOffsets = []*MoveOffset{
    {-1, -1}, {-1, 1}, {1, -1}, {1, 1},
}
var rookOffsets = []*MoveOffset{
    {-1, 0}, {0, 1}, {1, 0}, {0, -1},
}
var queenOffsets = append(bishopOffsets,rookOffsets...)



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
    
    pos := Position{
        Files:      files,
        Ranks:      ranks,
        WhiteToMove: sideToMove,
        Castling:   castlingRights,
        EnPassant:  enPassant,
        Pieces:     make(map[rune]*Bitboard),
        Walls:      NewBitboard(),
        PositionBitBoard: NewBitboard(),
        ColorBitboards: make(map[Color]*Bitboard),

    }

    rank := 0
    file := 0
    pos.ColorBitboards[Black] = NewBitboard()
    pos.ColorBitboards[White] = NewBitboard()
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
        ogIndex := FileRankToLargeIndex(file,rank,files)
        if ch == '.' {
            pos.Walls.SetBit(ogIndex)
        } else {
            if _, exists := pos.Pieces[ch]; !exists {
                pos.Pieces[ch] = NewBitboard()
            }
            if unicode.IsLower(ch){
                pos.ColorBitboards[Black].SetBit(ogIndex)
            } else { pos.ColorBitboards[White].SetBit(ogIndex) }
            pos.Pieces[ch].SetBit(ogIndex)
            pos.PositionBitBoard.SetBit(ogIndex)
        }

        file++
    }

    return &pos,nil
}

const (
    KSCW uint8 = 3
    QSCW uint8 = 2
    KSCB uint8 = 1
    QSCB uint8 = 0
)
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





func CreateGame(gameConfig models.GameConfig) (*GameState,error){
    variant , err := NewVariant(gameConfig)
    if err!=nil{
        return nil,err
    }
    game := GameState{
        Variant: variant,
        history: make([]Move,0),
    }
    return &game, nil
}

func GenerateMovesForPiece(piece rune, src int, position *Position, canCaptureKing bool) []Move {
    var moves []Move
    
    moveBitboard := NewBitboard()
    standardPieces := map[rune]bool{
        'p': true,
        'k': true,
        'q': true,
        'r': true,
        'b': true,
        'n': true,
    }
    if _,ok := standardPieces[unicode.ToLower(piece)]; ok {
        switch unicode.ToLower(piece) {
        case 'n':
            moveBitboard = generateJumpMoves(src, position, canCaptureKing, knightOffsets)
        case 'b':
            moveBitboard = generateSlideMoves(src, position, canCaptureKing, bishopOffsets)
        case 'r':
            moveBitboard = generateSlideMoves(src, position, canCaptureKing, rookOffsets)
        case 'q':
            moveBitboard = generateSlideMoves(src, position, canCaptureKing, queenOffsets)
        case 'k':
            moveBitboard = generateKingMoves(src, position, canCaptureKing)
        case 'p':
            moveBitboard = generatePawnMoves(src, position, canCaptureKing)
        }
    } else { 
        if patterns, exists := position.CustomPieceRules[piece]; exists {
            for _,pattern := range patterns{
                if pattern.MoveType==Jump{
                    moveBitboard = generateJumpMoves(src, position, canCaptureKing, pattern.MoveOffsets)
                } else if pattern.MoveType == Slide{
                    moveBitboard = generateSlideMoves(src, position, canCaptureKing, pattern.MoveOffsets)
                }
            }        
        }
    }
    var opponentBitboard *Bitboard
    if unicode.IsLower(piece){
        opponentBitboard = position.ColorBitboards[White]
    } else { opponentBitboard = position.ColorBitboards[Black]}

    moves = flattenBitboard(src,position,moveBitboard, opponentBitboard, piece)
    return moves
}

func flattenBitboard(src int, position *Position, moveBitboard *Bitboard, opponentBitboard *Bitboard, piece rune) []Move {
    var moves []Move
    setBits := moveBitboard.GetSetBits()

    for _, pos := range setBits {
        move := Move{
            Piece:        piece,
            From:         src,
            To:           pos,
            Capture:      opponentBitboard.HasBit(pos),
            ClassicMoveType: QuietMove,
        }

        if move.Capture {
            move.ClassicMoveType = CaptureMove
        }

        if piece == 'p' || piece == 'P' {
            lastRank := 0
            firstRank := (position.Ranks - 1) * position.Files

            if pos >= firstRank || pos <= lastRank {
                move.ClassicMoveType = PromotionMove
                move.Promotion = 'q'
            }

            startRankWhite := 1
            startRankBlack := position.Ranks - 2

            srcRank := src / position.Files
            dstRank := pos / position.Files
            rankDiff := abs(dstRank - srcRank)

            if rankDiff == 2 && !move.Capture {
                if (piece == 'p' && srcRank == startRankWhite) || (piece == 'P' && srcRank == startRankBlack) {
                    move.ClassicMoveType = DoublePawnPush
                }
            }
        }

        if piece == 'k' || piece == 'K' {
            dx := abs(src%position.Files - pos%position.Files)
            if dx == 2 {
                move.ClassicMoveType = CastleMove
            }
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



func generateKingMoves(src int, position *Position, canCaptureKing bool,)*Bitboard{
    b := NewBitboard()
    return b
}

func generatePawnMoves(src int, position *Position, canCaptureKing bool,)*Bitboard{
    b := NewBitboard()
    return b
}

func (p *Position) getPieceAt(src int) rune{
    for piece, bb := range p.Pieces {
        if bb.HasBit(src) {
            return piece
        }
    }
    return -1
}


func generateSlideMoves2(src int, position *Position, canCaptureKing bool, offsets []*MoveOffset) *Bitboard {
    moveBitboard := NewBitboard()
    x, y := src%position.Files, src/position.Files

    for _, offset := range offsets {
        nx, ny := x, y
        for {
            nx += offset.x
            ny += offset.y

            if nx < 0 || ny < 0 || nx >= position.Files || ny >= position.Ranks {
                break
            }

            newPos := ny*position.Files + nx

            if position.Walls.HasBit(newPos) {
                break
            }

            isOccupied := position.PositionBitBoard.HasBit(newPos)

            if !isOccupied {
                moveBitboard.SetBit(newPos)
                continue
            }

            var targetPiece rune
            for piece, bb := range position.Pieces {
                if bb.HasBit(newPos) {
                    targetPiece = piece
                    break
                }
            }

            isWhite := position.ColorBitboards[White].HasBit(newPos)
            isSameColor := position.ColorBitboards[White].HasBit(src) == isWhite
            isOpponentKing := targetPiece == 'k' || targetPiece == 'K'

            if !isSameColor && (canCaptureKing || !isOpponentKing) {
                moveBitboard.SetBit(newPos)
            }

            break
        }
    }

    return moveBitboard
}

func generateSlideMoves(src int, p *Position, canCaptureKing bool, offsets []*MoveOffset) *Bitboard{
    moveBitboard := NewBitboard()
    EnsureTablesInitialized()
    ogSrc := SmallToLargeBoardIndex(src,p.Files)
    validBB := GetValidSquares(p.Ranks,p.Files)
    for _, dir := range offsets {
        dirKey := *dir
        attackMask, exists := SlidingAttackTables[dirKey][ogSrc]
        if !exists {
            continue
        }
        
        blockers := NewBitboard()
        PrintBitboard(p.PositionBitBoard,p.Ranks,p.Files)
        PrintBitboard(p.PositionBitBoard,16,16)
        PrintBitboard(attackMask,p.Ranks,p.Files)
        PrintBitboard(attackMask,16,16)
        blockers.Bits.And(attackMask.Bits,p.PositionBitBoard.Bits)
        PrintBitboard(blockers,p.Ranks,p.Files)
        PrintBitboard(blockers,16,16)
        blockerBit := blockers.Bits.Clone()
        blockerBit.And(blockerBit, blockerBit.Neg(blockerBit))
        limit := blockerBit.Clone()
        limit.Sub(limit, uint256.NewInt(1))
        legalMoves := NewBitboard()
        legalMoves.Bits.And(attackMask.Bits, limit)
        PrintBitboard(legalMoves,p.Ranks,p.Files)
        PrintBitboard(legalMoves,16,16)
        firstBlockerIdx := blockers.FirstSetBit()
        if firstBlockerIdx != -1 {
            isWhite := p.ColorBitboards[White].HasBit(firstBlockerIdx)
            isSameColor := p.ColorBitboards[White].HasBit(src) == isWhite
            targetPiece := p.getPieceAt(firstBlockerIdx)

            isOpponentKing := targetPiece == 'k' || targetPiece == 'K'
            if !isSameColor && (canCaptureKing || !isOpponentKing) {
                legalMoves.SetBit(firstBlockerIdx)
            }
        }
        moveBitboard.Bits.Or(moveBitboard.Bits, legalMoves.Bits)
    }
    PrintBitboard(moveBitboard,p.Ranks,p.Files)
    PrintBitboard(moveBitboard,16,16)
    moveBitboard.Bits.And(validBB.Bits,moveBitboard.Bits)
    return moveBitboard
}

/*func GeneratePieceMove2s(src int, position *Position, piece rune) *Bitboard {
    moves := NewBitboard()

    // Get allowed movement directions
    patterns, exists := position.CustomPieceRules[piece]
    if !exists {
        return moves // No move pattern found
    }

    // Combine precomputed attack tables for selected directions
    for _, dir := range patterns {
        attackMask, ok := SlidingAttackTables[dir][src]
        if !ok {
            continue
        }

        // Apply occupancy constraints
        blockers := NewBitboard()
        blockers.Bits.And(position.PositionBitBoard.Bits, attackMask.Bits)

        // Stop sliding at first blocker
        blockerBit := blockers.Bits.Clone()
        blockerBit.And(blockerBit, blockerBit.Neg(blockerBit)) // Get lowest blocker bit
        limit := blockerBit.Clone()
        limit.Sub(limit, uint256.NewInt(1)) // Mask up to blocker

        moves.Bits.And(attackMask.Bits, limit)
    }

    return moves
}*/


func generateJumpMoves(src int, position *Position,canCaptureKing bool, offsets []*MoveOffset) *Bitboard {
    moveBitboard := NewBitboard()

    x, y := src%position.Files, src/position.Files
    for _, offset := range offsets {
        nx, ny := x+offset.x, y+offset.y
        if nx < 0 || ny < 0 || nx >= position.Files || ny >= position.Ranks {
            continue
        }

        newPos := ny*position.Files + nx
        //println("bb pos BEFORE",position.PositionBitBoard.Bits.String())
        isOccupied := position.PositionBitBoard.HasBit(newPos) 
        //println("bb pos AFTER",position.PositionBitBoard.Bits.String())
        if !isOccupied {
            moveBitboard.SetBit(newPos)
            continue
        }

        var targetPiece rune
        for piece, bb := range position.Pieces {
            if bb.HasBit(newPos) {
                targetPiece = piece
                break
            }
        }
        //println("bb White",position.ColorBitboards[White].Bits.String(),src,newPos)
        isWhite := position.ColorBitboards[White].HasBit(newPos)
        isSameColor := position.ColorBitboards[White].HasBit(src) == isWhite
        isOpponentKing := targetPiece == 'k' || targetPiece == 'K'

        if !isSameColor && (canCaptureKing || !isOpponentKing) {
            moveBitboard.SetBit(newPos)
        }
    }

    return moveBitboard
}


