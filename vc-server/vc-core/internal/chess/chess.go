package chess

import (
	"sync"
	"unicode"
	"vc-server/vc-core/internal/models"
)



type Game struct {
    ID        string
    Position  Position
    MoveHistory []string
}


func SmallToLargeBoardIndex(sq, smallBoardWidth int,largestDim int) int {
    var largeBoardWidth int = 8
    if largestDim > 8{
        largeBoardWidth = 16
    }
    rank := sq / smallBoardWidth
    file := sq % smallBoardWidth
    return FileRankToIndex(file, rank, largeBoardWidth)
}


func FileRankToIndex(file, rank, width int) int {
    return rank*width + file
}

func FileRankToLargeIndex(file,rank,smallBoardWidth,largestDim int) int{
    return SmallToLargeBoardIndex(FileRankToIndex(file,rank,smallBoardWidth),smallBoardWidth,largestDim)
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

    for ranks := 5; ranks <= 16; ranks++ {
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

    for files := 5; files <= 16; files++ {
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
            if square<64{
                SlidingAttackTables[*dir][square] = ComputeSlideAttacks(square,dir,8)
            }
            SlidingAttackTables256[*dir][square] = ComputeSlideAttacks(square, dir,16)
        }
    }
}

func ComputeSlideAttacks(src int, dir *MoveOffset, largestDimension int) Bitboard {
    moves := NewBitboard(largestDimension)
    x, y := src%largestDimension, src/largestDimension
    for {
        x += dir.x
        y += dir.y

        if x < 0 || y < 0 || x >= 16 || y >= 16 {
            break
        }
        moves.SetBit(y*largestDimension + x)
    }
    return moves
}

func GetValidSquares(ranks, files int) Bitboard {
    bb := NewBitboard(max(ranks,files))
    rankMask, fileMask := GetRankFileMasks(ranks,files)

    bb.Or(rankMask,fileMask)
    PrintBitboard(bb,8,8)
    //PrintBitboard(bb,16,16)
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
    //var largestDimension int = max(position.Ranks,position.Files)
    //moveBitboard := NewBitboard(largestDimension)
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
            generateJumpMoves(src,piece, position, canCaptureKing, knightOffsets,&moves)
        case 'b':
            generateSlideMoves2(src, piece, position, canCaptureKing, bishopOffsets,&moves)
        case 'r':
            generateSlideMoves2(src, piece,position, canCaptureKing, rookOffsets,&moves)
        case 'q':
            generateSlideMoves2(src,piece, position, canCaptureKing, queenOffsets,&moves)
        case 'k':
            generateKingMoves(src, position, canCaptureKing,&moves)
        case 'p':
            generatePawnMoves(src, position, canCaptureKing,&moves)
        }
    } else { 
        if patterns, exists := position.CustomPieceRules[piece]; exists {
            for _,pattern := range patterns{
                if pattern.MoveType==Jump{
                    generateJumpMoves(src,piece, position, canCaptureKing, pattern.MoveOffsets,&moves)
                } else if pattern.MoveType == Slide{
                    generateSlideMoves2(src,piece, position, canCaptureKing, pattern.MoveOffsets,&moves)
                }
            }        
        }
    }
    /*var opponentBitboard Bitboard
    if unicode.IsLower(piece){
        opponentBitboard = position.ColorBitboards[White]
    } else { opponentBitboard = position.ColorBitboards[Black]}

    moves = flattenBitboard(src,position,moveBitboard, opponentBitboard, piece)*/
    return moves
}

func flattenBitboard(src int, position *Position, moveBitboard Bitboard, opponentBitboard Bitboard, piece rune) []Move {
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



func generateKingMoves(src int, p *Position, canCaptureKing bool,moves *[]Move)Bitboard{
    b := NewBitboard(max(p.Ranks,p.Files))
    return b
}

func generatePawnMoves(src int, p *Position, canCaptureKing bool,moves *[]Move){
    piece := p.getPieceAt(src)
    oppKing := 'k'
    var color Color = White
    var oppColor Color = Black
    if unicode.IsLower(piece){
        color = Black
        oppColor = White
    }
    direction := -1
    startRank, promotionRank := 1, p.Ranks-1

    if color==Black {
        oppKing = 'K'
        direction = 1
        startRank, promotionRank = p.Ranks-2, 0
    }
    file := src % p.Files
    rank := src / p.Files
    oneStep := FileRankToLargeIndex(file, rank+direction, p.Files,p.LargestDimension)
    if !p.PositionBitBoard.HasBit(oneStep) && rank+direction >= 0 && rank+direction < p.Ranks {
        if rank+direction == promotionRank {
            for _, promoPiece := range []rune{'Q', 'R', 'B', 'N'} {
                *moves = append(*moves, Move{
                    Piece: piece, From: src, To: oneStep,
                    Promotion: promoPiece, ClassicMoveType: PromotionMove,
                })
            }
        } else {
            *moves = append(*moves, Move{
                Piece: piece, From: src, To: oneStep, ClassicMoveType: QuietMove,
            })
        }

        twoStep := FileRankToLargeIndex(file, rank+2*direction, p.Files,p.LargestDimension)
        if rank == startRank && !p.PositionBitBoard.HasBit(twoStep) {
            *moves = append(*moves, Move{
                Piece: piece, From: src, To: twoStep, ClassicMoveType: DoublePawnPush,
            })
        }
    }
    
    captureOffsets := []int{-1, 1}
    for _, offset := range captureOffsets {
        newFile := file + offset
        newRank := rank + direction
        if newFile >= 0 && newFile < p.Files && newRank >= 0 && newRank < p.Ranks {
            capturePos := FileRankToLargeIndex(newFile, newRank, p.Files,p.LargestDimension)
            //capturedPiece := p.getPieceAt(capturePos)

            if p.ColorBitboards[oppColor].HasBit(capturePos) || (canCaptureKing && p.Pieces[oppKing].HasBit(capturePos)) {
                if newRank == promotionRank {
                    for _, promoPiece := range []rune{'Q', 'R', 'B', 'N'} {
                        *moves = append(*moves, Move{
                            Piece: piece, From: src, To: capturePos,
                            Capture: true, Promotion: promoPiece,
                            ClassicMoveType: PromotionMove,
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

    if p.EnPassant != -1 {
        for _, offset := range captureOffsets {
            enPassantFile := file + offset
            if enPassantFile >= 0 && enPassantFile < p.Files {
                enPassantPos := FileRankToLargeIndex(enPassantFile, rank, p.Files, p.LargestDimension)
                if enPassantPos == p.EnPassant {
                    targetSquare := FileRankToLargeIndex(enPassantFile, rank+direction, p.Files, p.LargestDimension)
                    *moves = append(*moves, Move{
                        Piece: piece, From: src, To: targetSquare,
                        Capture: true, ClassicMoveType: EnPassant,
                        VariantMoveType: "En Passant",
                    })
                }
            }
        }
    }
}

func (p *Position) getPieceAt(src int) rune{
    for piece, bb := range p.Pieces {
        if bb.HasBit(src) {
            return piece
        }
    }
    return -1
}

func generateSlideMoves2(src int,piece rune, p *Position, canCaptureKing bool, offsets []*MoveOffset,moves *[]Move){
    moveBitboard := NewBitboard(p.LargestDimension)
    x, y := src%p.Files, src/p.Files

    for _, offset := range offsets {
        nx, ny := x, y
        for {
            nx += offset.x
            ny += offset.y

            if nx < 0 || ny < 0 || nx >= p.Files || ny >= p.Ranks {
                break
            }

            newPos := ny*p.Files + nx

            if p.Walls.HasBit(newPos) {
                break
            }

            isOccupied := p.PositionBitBoard.HasBit(newPos)

            if !isOccupied {
                moveBitboard.SetBit(newPos)
                continue
            }

            var targetPiece rune
            for piece, bb := range p.Pieces {
                if bb.HasBit(newPos) {
                    targetPiece = piece
                    break
                }
            }

            isWhite := p.ColorBitboards[White].HasBit(newPos)
            isSameColor := p.ColorBitboards[White].HasBit(src) == isWhite
            isOpponentKing := targetPiece == 'k' || targetPiece == 'K'

            if !isSameColor && (!isOpponentKing || (canCaptureKing && isOpponentKing)) {
                moveBitboard.SetBit(newPos)
            }

            break
        }
    }
    validBB := GetValidSquares(p.Ranks,p.Files)
    moveBitboard.And(validBB,moveBitboard)
    targets := moveBitboard.GetSetBits()
    for _,target := range targets{
    
        isCapture := p.PositionBitBoard.HasBit(target)

        moveType := QuietMove
        if isCapture {
            moveType = CaptureMove
        }
        *moves = append(*moves, Move{
            Piece: piece,
            From: src,
            To: target,
            Capture: isCapture,
            ClassicMoveType: moveType,
        })
    }
}

func generateSlideMoves(src int, piece rune,p *Position, canCaptureKing bool, offsets []*MoveOffset, moves *[]Move) Bitboard{
    largestDimension := max(p.Ranks,p.Files)
    moveBitboard := NewBitboard(largestDimension)
    EnsureTablesInitialized()
    ogSrc := SmallToLargeBoardIndex(src,p.Files,largestDimension)
    validBB := GetValidSquares(p.Ranks,p.Files)
    attackTables := GetAttackTables(largestDimension)
    for _, dir := range offsets {
        dirKey := *dir
        attackMask, exists := attackTables[dirKey][ogSrc]
        if !exists {
            continue
        }
        
        blockers := NewBitboard(largestDimension)
        PrintBitboard(p.PositionBitBoard,p.Ranks,p.Files)
        //PrintBitboard(p.PositionBitBoard,16,16)
        PrintBitboard(attackMask,p.Ranks,p.Files)
        //PrintBitboard(attackMask,16,16)
        blockers.And(attackMask,p.PositionBitBoard)
        PrintBitboard(blockers,p.Ranks,p.Files)
        //PrintBitboard(blockers,16,16)
        blockerBit := blockers.Clone()
        blockerBit.And(blockerBit, blockerBit.Neg())
        PrintBitboard(blockerBit,p.Ranks,p.Files)
        
        limit := blockerBit.Clone()
        limit.Sub(limit.One())
        PrintBitboard(limit,p.Ranks,p.Files)
        legalMoves := NewBitboard(largestDimension)
        legalMoves.And(attackMask, limit)
        PrintBitboard(legalMoves,p.Ranks,p.Files)
        //PrintBitboard(legalMoves,16,16)
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
        moveBitboard.Or(moveBitboard, legalMoves)
    }
    PrintBitboard(moveBitboard,p.Ranks,p.Files)
    PrintBitboard(moveBitboard,16,16)
    moveBitboard.And(validBB,moveBitboard)
    return moveBitboard
}

func generateJumpMoves(src int,piece rune, p *Position,canCaptureKing bool, offsets []*MoveOffset, moves *[]Move) {
    largestDimension := max(p.Ranks,p.Files)
    moveBitboard := NewBitboard(largestDimension)
    //ogSrc := SmallToLargeBoardIndex(src,p.Files,largestDimension)
    x, y := src%p.Files, src/p.Files
    for _, offset := range offsets {
        nx, ny := x+offset.x, y+offset.y
        if nx < 0 || ny < 0 || nx >= p.Files || ny >= p.Ranks {
            continue
        }
        pos := FileRankToIndex(ny,nx,p.Files)
        newPos := SmallToLargeBoardIndex(pos,p.Files,largestDimension)// ny*position.Files + nx
        isOccupied := p.PositionBitBoard.HasBit(newPos) 
        if !isOccupied {
            moveBitboard.SetBit(newPos)
            continue
        }

        var targetPiece rune
        for piece, bb := range p.Pieces {
            if bb.HasBit(newPos) {
                targetPiece = piece
                break
            }
        }

        isWhite := p.ColorBitboards[White].HasBit(newPos)
        isSameColor := p.ColorBitboards[White].HasBit(src) == isWhite
        isOpponentKing := targetPiece == 'k' || targetPiece == 'K'

        if !isSameColor && (!isOpponentKing || (canCaptureKing && isOpponentKing)) {
            moveBitboard.SetBit(newPos)
        }
    }
    validBB := GetValidSquares(p.Ranks,p.Files)
    moveBitboard.And(validBB)
    targets := moveBitboard.GetSetBits()
    for _,target := range targets{
    
        isCapture := p.PositionBitBoard.HasBit(target)

        moveType := QuietMove
        if isCapture {
            moveType = CaptureMove
        }
        *moves = append(*moves, Move{
            Piece: piece,
            From: src,
            To: target,
            Capture: isCapture,
            ClassicMoveType: moveType,
        })
    }

}