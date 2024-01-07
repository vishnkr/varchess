package chesscore

import (
	"strings"
)

type variantType string

const (
	Custom variantType = "Custom"
	DuckChess variantType = "DuckChess"
	ArcherChess variantType = "ArcherChess"
	Wormhole variantType = "Wormhole"
)

type Variant interface {
	makeMove(Move)
	unmakeMove(Move)
	getPseudoLegalMoves(Color, bool) []Move
	GetLegalMoves() []Move
	PerformMove(Move) (result,bool)
}

type variant struct {
	Objective Objective `json:"objective"`
	position
	variantType
	recentCapture recentCaptureInfo
	possibleLegalMoves []Move
	gameResult result
	isGameOverBool bool
}

type recentCaptureInfo struct{
	piece
	squareId int
	moveType classicMoveType
}

func (v *variant) getTargetSquare(currentSquareID int, offset moveOffset) (int, bool) {
	row, col := v.toRowCol(currentSquareID)
	newRow, newCol := row+offset.y, col+offset.x
	target := v.toPos(newRow, newCol)
	if newRow < 0 || newCol < 0 || newRow >= v.dimensions.Ranks || newCol >= v.dimensions.Files || v.isWall(target) {
		return -1, false
	}
	return target, true
}

func (v variant) genPawnMoves(piece piece, currentSquareID int) []Move{
	validMoves := []Move{}
	var rowOffset, doubleMoveStartRank int
	if piece.color == ColorWhite {
		doubleMoveStartRank = v.dimensions.Ranks - 2
		rowOffset = -1
	} else {
		doubleMoveStartRank = 1
		rowOffset = 1
	}
	srcRow, srcCol := v.toRowCol(currentSquareID)
	targetRow := srcRow + rowOffset
	for i := -1; i <= 1; i++ {
		//non-capture moves
		if i == 0 {
			target1, target2 := v.toPos(targetRow, srcCol), v.toPos(targetRow+rowOffset, srcCol)
			_, ok1 := v.pieceLocations[target1]
			_, ok2 := v.pieceLocations[target2]
			if !ok1 && !v.wallSquares[target1] {
				validMoves = append(validMoves, Move{
					Source:          currentSquareID,
					Target:          target1,
					ClassicMoveType: QuietMove,
					PieceType:       Pawn,
					PieceNotation:   piece.notation,
					Turn:            piece.color,
				})
			}
			if !ok1 && !v.wallSquares[target1] && !ok2 && srcRow == doubleMoveStartRank && !v.wallSquares[target2] {
				validMoves = append(validMoves, Move{
					Source:          currentSquareID,
					Target:          target2,
					ClassicMoveType: DoublePawnPush,
					PieceType:       Pawn,
					PieceNotation:   piece.notation,
					Turn:            piece.color,
				})
			}
		} else {
			//capture moves
			target := v.toPos(targetRow, srcCol+i)
			if v.isOpponentPiecePresent(currentSquareID, target) {
				if (!v.isOppKingAtTarget(target) || (v.isOppKingAtTarget(target) && v.additionalProps.kingCaptureAllowed)) {
					validMoves = append(validMoves, Move{
						Source:          currentSquareID,
						Target:          target,
						ClassicMoveType: CaptureMove,
						PieceType:       Pawn,
						PieceNotation:   piece.notation,
						Turn:            piece.color,
					})
				}
			} else if (v.isEmpty(target) && v.epSquare == target){
				validMoves = append(validMoves,Move{
					Source: currentSquareID,
					Target: target,
					ClassicMoveType: EnPassantMove,
					PieceType: Pawn,
					PieceNotation: piece.notation,
					Turn: piece.color,
				})
			}

		}
	}
	return validMoves
}

func (v variant) genKingMoves(piece piece, currentSquareId int) []Move{
	validMoves := []Move{}
	curRow, curCol := v.toRowCol(currentSquareId)
	for row := -1; row <= 1; row += 1 {
		for col := -1; col <= 1; col += 1 {
			target := v.toPos(curRow+row, curCol+col)
			if (row == 0 && col == 0) || (target<0 || target>=v.dimensions.Ranks*v.dimensions.Files) {
				continue
			}
			if !v.isSameColorPiecePresent(currentSquareId, target) && !v.isWall(target) {
				targetPiece, ok := v.pieceLocations[target]
				if !ok {
					validMoves = append(validMoves, Move{
						Source:          currentSquareId,
						Target:          target,
						ClassicMoveType: QuietMove,
						PieceType:       King,
						PieceNotation:   piece.notation,
						Turn:            piece.color,
					})
				} else if targetPiece.pieceType != King || (targetPiece.pieceType == King && v.additionalProps.kingCaptureAllowed) {
					validMoves = append(validMoves, Move{Source: currentSquareId, Target: target, ClassicMoveType: CaptureMove, PieceType: King, PieceNotation: piece.notation, Turn: piece.color})
				}
			}
		}
	}

	var validKingside, validQueenSide bool
	var kingPos int
	if piece.color == ColorBlack {
		validKingside = v.blackKingSide
		validQueenSide = v.blackQueenSide
		kingPos = v.additionalProps.blackKingPos
	} else {
		validKingside = v.whiteKingSide
		validQueenSide = v.whiteQueenSide
		kingPos = v.additionalProps.whiteKingPos
	}
	kCastleAllowed,kRookPos := v.isCastleAllowed(piece.color,kingPos,false)
	if validKingside && kCastleAllowed {
		validMoves = append(validMoves, Move{
			Source: kingPos, 
			Target: v.toPos(curRow, curCol-2), 
			Turn:piece.color, 
			ClassicMoveType: CastleMove,
			PieceType: King,
			AdditionalData: map[string][]int{"rookPos":kRookPos},
		})
	}
	qCastleAllowed,qRookPos := v.isCastleAllowed(piece.color,kingPos,false)
	if validQueenSide && qCastleAllowed {
		validMoves = append(validMoves, Move{
			Source: v.toPos(curRow, curCol), 
			Target: v.toPos(curRow, curCol+2), 
			Turn:piece.color, 
			ClassicMoveType: CastleMove, 
			PieceType: King,
			AdditionalData: map[string][]int{"rookPos":qRookPos},
		})
	}
	return validMoves
}

// returns if castle is valid(bool) and rook src/target if it is (int)
func (v variant) isCastleAllowed(color Color,kingPos int,isKingside bool) (bool,[]int){
	curRow,_ := v.toRowCol(kingPos)
	var rookSrc, rookTarget, dx,i int
	if v.Objective.Type == Antichess{
		return false,[]int{-1,-1}
	}
	if isKingside{
		rookSrc,dx = v.toPos(curRow,v.dimensions.Files-1),1
		i = kingPos+dx
		rookTarget = i
		if piece,ok := v.pieceLocations[rookSrc]; ok && piece.color == color && piece.pieceType == Rook{
			for i<rookSrc{
				if !v.isEmpty(i) || v.isWall(i){ return false,[]int{-1,-1}}
				i+=dx
			}
			
		}
	} else { 
		rookSrc,dx = v.toPos(curRow,0),-1
		i = kingPos+dx
		rookTarget = i
		if piece,ok := v.pieceLocations[rookSrc]; ok && piece.color == color && piece.pieceType == Rook{
			for i>rookSrc{
				if !v.isEmpty(i) || v.isWall(i){ return false,[]int{-1,-1}}
				i+=dx
			}
		}
	}
	return true,[]int{rookSrc,rookTarget}
}

func (v variant) getPseudoLegalMoves(color Color, kingCaptureAllowed bool) []Move {
	validMoves := []Move{}
	for currentSquareID, piece := range v.pieceLocations {
		if color != piece.color {
			continue
		}
		if piece.pieceType == Pawn {
			validMoves = append(validMoves,v.genPawnMoves(piece, currentSquareID)...)
		} else if piece.pieceType == King {
			validMoves = append(validMoves,v.genKingMoves(piece, currentSquareID)...)
		} else {
			var props PieceProperties = v.pieceProps[strings.ToLower(piece.notation)]
			validMoves = append(validMoves, v.genSlideMoves(&piece, currentSquareID, &props)...)
			validMoves = append(validMoves, v.genJumpMoves(&piece, currentSquareID, &props)...)
		}
	}
	return validMoves
}

func (v *variant) genSlideMoves(piece *piece, currentSquareID int, props *PieceProperties) []Move{
	validMoves := []Move{}
	for _, offset := range props.SlideOffsets {
		target, valid := v.getTargetSquare(currentSquareID, offset)
		for valid {
			if !v.isSameColorPiecePresent(currentSquareID, target) {
				v.position.attackedSquares[target] = true
				if v.isOppKingAtTarget(target) && !v.additionalProps.kingCaptureAllowed {
					valid = false
					continue
				}
				mType := v.position.getclassicMoveType(target)
				validMoves = append(validMoves, Move{
					Source:          currentSquareID,
					Target:          target,
					ClassicMoveType: mType,
					PieceType:       piece.pieceType,
					PieceNotation:   piece.notation,
					Turn:            piece.color,
				})
				if mType != QuietMove {
					break
				}
			} else {
				break
			}
			target, valid = v.getTargetSquare(target, offset)
		}
	}
	return validMoves
}

func (v *variant) genJumpMoves(piece *piece, currentSquareID int, props *PieceProperties) []Move{
	validMoves := []Move{}
	for _, jumpMove := range props.JumpOffsets {
		target, valid := v.getTargetSquare(currentSquareID, jumpMove)
		if valid {
			if v.position.isSameColorPiecePresent(currentSquareID, target) {
				continue
			}
			var moveType classicMoveType = NullMove
			if  v.position.isOpponentPiecePresent(currentSquareID, target) {
				if !v.position.isOppKingAtTarget(target) || (!v.position.isOppKingAtTarget(target) && v.additionalProps.kingCaptureAllowed) {
					moveType = CaptureMove
				}
			} else {
				moveType = QuietMove
			}
			v.position.attackedSquares[target] = true
			move := Move{
				Source:          currentSquareID,
				Target:          target,
				ClassicMoveType: moveType,
				PieceType:       piece.pieceType,
				PieceNotation:   piece.notation,
				Turn:            piece.color,
			}
			validMoves = append(validMoves, move)
		}
	}
	return validMoves
}

func (v *variant) makeMove(move Move) {
	switch move.ClassicMoveType {
	case QuietMove, CaptureMove, DoublePawnPush:
		if move.ClassicMoveType == CaptureMove {
			v.recentCapture = recentCaptureInfo{
				piece: v.pieceLocations[move.Target],
				squareId: move.Target,
				moveType: move.ClassicMoveType,
			}
		}
		if move.PieceType==King{
			if move.Turn == ColorBlack{
				v.additionalProps.blackKingPos = move.Target
			} else {
				v.additionalProps.whiteKingPos = move.Target
			}
		}
		v.pieceLocations[move.Target] = v.pieceLocations[move.Source]
		delete(v.pieceLocations, move.Source)
	case EnPassantMove:
		v.pieceLocations[move.Target] = v.pieceLocations[move.Source]
		v.recentCapture = recentCaptureInfo{
			piece: v.pieceLocations[v.epSquare],
			squareId: v.epSquare,
			moveType: move.ClassicMoveType,
		}
		delete(v.pieceLocations,v.epSquare)
		v.epSquare = -1
	case CastleMove:
		additionalData, ok := move.AdditionalData.(map[string][]int)
		if ok {
			rookPos := additionalData["rookPos"]
			v.pieceLocations[move.Target] = v.pieceLocations[move.Source]
			v.pieceLocations[rookPos[1]] = v.pieceLocations[rookPos[0]]
			delete(v.pieceLocations,rookPos[0])
			delete(v.pieceLocations,move.Source)
		}
	}
}

func (v *variant) unmakeMove(move Move) {
	switch move.ClassicMoveType {
	case QuietMove, CaptureMove, DoublePawnPush:
		v.pieceLocations[move.Source] = v.pieceLocations[move.Target]
		if move.PieceType==King{
			if move.Turn == ColorBlack{
				v.additionalProps.blackKingPos = move.Source
			} else {
				v.additionalProps.whiteKingPos = move.Source
			}
		}
		if move.ClassicMoveType == CaptureMove {
			v.pieceLocations[v.recentCapture.squareId] = v.recentCapture.piece
		} else { delete(v.pieceLocations,move.Target)}
	case EnPassantMove:
		v.pieceLocations[v.recentCapture.squareId] = v.recentCapture.piece
		v.pieceLocations[move.Source] = v.pieceLocations[move.Target]
		delete(v.pieceLocations,move.Target)
	case CastleMove:
		additionalData, ok := move.AdditionalData.(map[string][]int)
		if ok {
			rookPos := additionalData["rookPos"]
			v.pieceLocations[move.Source] = v.pieceLocations[move.Target]
			v.pieceLocations[rookPos[0]] = v.pieceLocations[rookPos[1]]
			delete(v.pieceLocations,rookPos[1])
			delete(v.pieceLocations,move.Target)
		}
	}
}

type CheckmateVariant struct {
	variant
}

func (cv *CheckmateVariant) GetLegalMoves() []Move {
	pseudoMoves := cv.getPseudoLegalMoves(cv.turn, false)
	var validMoves []Move = []Move{}
	attackedOriginal := cv.attackedSquares
	for _, mv := range pseudoMoves {
		cv.variant.makeMove(mv)
		//check king in check
		cv.attackedSquares = make(map[int]bool)
		if !cv.isKingUnderCheck(cv.turn) {
			validMoves = append(validMoves, mv)
		}
		cv.variant.unmakeMove(mv)
	}
	cv.attackedSquares = attackedOriginal
	return validMoves
}

func (cv *CheckmateVariant) isKingUnderCheck(color Color) bool {
	cv.getPseudoLegalMoves(cv.getOpponentColor(), false)
	var kingPos int
	if color == ColorWhite {
		kingPos = cv.variant.position.additionalProps.whiteKingPos
	} else {
		kingPos = cv.variant.position.additionalProps.blackKingPos
	}
	_, ok := cv.attackedSquares[kingPos]
	return ok
}

func (cv *CheckmateVariant) checkGameOver() (result, bool) {
	cv.possibleLegalMoves = cv.GetLegalMoves()
	if cv.turn==ColorBlack{
		bc := cv.isKingUnderCheck(ColorBlack)
		if len(cv.possibleLegalMoves)==0{
			if bc { return WhiteWins,true } else {return Stalemate,true}
		}
	} else {
		wc := cv.isKingUnderCheck(ColorBlack)
		if len(cv.possibleLegalMoves)==0{
			if wc { return BlackWins,true } else {return Stalemate,true}
		}
	}
	return 0, false
}

func (v *variant) IsGameOver() (result, bool) {
	return v.gameResult,v.isGameOverBool
}

func (cv *CheckmateVariant) PerformMove(move Move)(result, bool){
	cv.makeMove(move)
	cv.switchTurn()
	res,over:= cv.checkGameOver()
	cv.gameResult = res
	cv.isGameOverBool = over
	return res,over
}

// ------------ N-Check -------------------
type NCheckVariant struct {
	variant
	blackKingCheckCount int
	whiteKingCheckCount int
	targetChecks        int
}

func (ncv *NCheckVariant) GetLegalMoves() []Move {
	pseudoMoves := ncv.getPseudoLegalMoves(ncv.turn, false)
	for _, mv := range pseudoMoves {
		ncv.makeMove(mv)
		//check if king is in check
		ncv.unmakeMove(mv)
	}
	return []Move{}
}

func (ncv *NCheckVariant) isKingUnderCheck(color Color) bool {
	ncv.getPseudoLegalMoves(ncv.getOpponentColor(), false)
	var kingPos int
	if color == ColorWhite {
		kingPos = ncv.variant.position.additionalProps.whiteKingPos
	} else {
		kingPos = ncv.variant.position.additionalProps.blackKingPos
	}
	_, ok := ncv.attackedSquares[kingPos]
	return ok
}

func (ncv *NCheckVariant) checkGameOver() (result, bool) {
	ncv.possibleLegalMoves = ncv.GetLegalMoves()
	if ncv.turn == ColorBlack && ncv.whiteKingCheckCount == ncv.targetChecks {
		return BlackWins, true
	} else if ncv.blackKingCheckCount == ncv.targetChecks {
		return WhiteWins, true
	}
	if ncv.turn==ColorBlack{
		bc := ncv.isKingUnderCheck(ColorBlack)
		if len(ncv.possibleLegalMoves)==0{
			if bc { return WhiteWins,true } else {return Stalemate,true}
		}
	} else {
		wc := ncv.isKingUnderCheck(ColorBlack)
		if len(ncv.possibleLegalMoves)==0{
			if wc { return BlackWins,true } else {return Stalemate,true}
		}
	}
	return 0, false
}

func (ncv *NCheckVariant) PerformMove(move Move)(result, bool){
	ncv.makeMove(move)
	ncv.switchTurn()
	res,over:= ncv.checkGameOver()
	ncv.gameResult = res
	ncv.isGameOverBool = over
	return res,over
}

// ------------ Antichess -------------------
type AntichessVariant struct {
	variant
}

func (av *AntichessVariant) GetLegalMoves() []Move {
	pseudoMoves := av.getPseudoLegalMoves(av.turn, true)
	var captureMoves []Move = []Move{}
	for _, move := range pseudoMoves {
		if move.ClassicMoveType == CaptureMove || move.ClassicMoveType == EnPassantMove {
			captureMoves = append(captureMoves, move)
		}
	}

	if len(captureMoves) > 0 {
		pseudoMoves = captureMoves
	}
	return pseudoMoves
}


func (av *AntichessVariant) PerformMove(move Move)(result, bool){
	av.makeMove(move)
	av.switchTurn()
	res,over:= av.checkGameOver()
	av.gameResult = res
	av.isGameOverBool = over
	return res,over
}

func (av *AntichessVariant) checkGameOver() (result, bool) {
	var whitePieceCount int = 0
	var blackPieceCount int = 0
	for _,piece := range av.position.pieceLocations{
		if piece.color==ColorBlack{
			blackPieceCount+=1
		} else{
			whitePieceCount+=1
		}
	}
	if whitePieceCount == 0{
		return WhiteWins,true
	} else if blackPieceCount == 0{
		return BlackWins,true
	}
	legalMoves := av.GetLegalMoves()
	if len(legalMoves)==0{
		return Stalemate,true
	}
	return av.checkStalemate()
}

func (av *AntichessVariant) checkStalemate() (result,bool){
	darkBlackBishops,lightBlackBishops,lightWhiteBishops, darkWhiteBishops :=0,0,0,0
	for pos,piece := range av.pieceLocations{
		if piece.pieceType!=Bishop{
			continue
		}
		if piece.color==ColorBlack{
			if pos%2==0{
				darkBlackBishops+=1
			} else { lightBlackBishops+=1}
		} else { 
			if pos%2==0{
				darkWhiteBishops+=1
			} else{ lightWhiteBishops+=1}
		}
	}
	if (darkBlackBishops!=0 && darkWhiteBishops!=0) || (lightBlackBishops!=0 && lightWhiteBishops!=0){
		return 0, false
	} 
	return Stalemate,true
}