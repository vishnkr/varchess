package game

import "unicode"

type variantType uint8

const (
	Custom variantType = iota
	DuckChess
	ArcherChess
	Wormhole
)

type Variant interface {
	makeMove(Move)
	unmakeMove(Move)
	getPseudoLegalMoves(color, bool) []Move
	GetLegalMoves() []Move
	IsGameOver() (result, bool)
}

type variant struct {
	Objective Objective `json:"objective"`
	position
	variantType
	recentCapture piece
}

func (v *variant) getTargetSquare(currentSquareID int, offset moveOffset) (int, bool) {
	row, col := v.toRowCol(currentSquareID)
	newRow, newCol := row+offset.yOffset, col+offset.xOffset
	target := v.toPos(newRow, newCol)
	if newRow < 0 || newCol < 0 || newRow >= v.Ranks || newCol >= v.Files || v.isDisabled(target) {
		return -1, false
	}
	return target, true
}

func (v variant) genPawnMoves(piece piece, currentSquareID int) []Move{
	validMoves := []Move{}
	var rowOffset, doubleMoveStartRank int
	if piece.color == ColorWhite {
		doubleMoveStartRank = v.Ranks - 2
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
			if _, ok := v.pieceLocations[target1]; !ok && !v.disabledSquares[target1] {
				validMoves = append(validMoves, Move{
					Source:          currentSquareID,
					Target:          target1,
					ClassicMoveType: QuietMove,
					PieceType:       Pawn,
					PieceNotation:   piece.notation,
					Turn:            piece.color,
				})
			}
			if _, ok := v.pieceLocations[target2]; !ok && srcRow == doubleMoveStartRank && !v.disabledSquares[target2] {
				validMoves = append(validMoves, Move{
					Source:          currentSquareID,
					Target:          target2,
					ClassicMoveType: QuietMove,
					PieceType:       Pawn,
					PieceNotation:   piece.notation,
					Turn:            piece.color,
				})
			}
		} else {
			//capture moves
			target := v.toPos(targetRow, srcCol+i)
			if v.isOpponentPiecePresent(currentSquareID, target) && (!v.isOppKingAtTarget(target) || (v.isOppKingAtTarget(target) && v.additionalProps.kingCaptureAllowed)) {
				validMoves = append(validMoves, Move{
					Source:          currentSquareID,
					Target:          target,
					ClassicMoveType: CaptureMove,
					PieceType:       Pawn,
					PieceNotation:   piece.notation,
					Turn:            piece.color,
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
			if (row == 0 && col == 0) || (target<0 || target>=v.Ranks*v.Files) {
				continue
			}
			if !v.isSameColorPiecePresent(currentSquareId, target) && !v.isDisabled(target) {
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
	if validKingside && v.isCastleAllowed(piece.color,kingPos,true) {
		validMoves = append(validMoves, Move{Source: kingPos, Target: v.toPos(curRow, curCol-2), Turn:piece.color, ClassicMoveType: CastleMove, PieceType: King})
	}
	if validQueenSide && v.isCastleAllowed(piece.color,kingPos,false) {
		validMoves = append(validMoves, Move{Source: v.toPos(curRow, curCol), Target: v.toPos(curRow, curCol+2), Turn:piece.color, ClassicMoveType: CastleMove, PieceType: King})
	}
	return validMoves
}

func (v variant) isCastleAllowed(color color,kingPos int,isKingside bool) bool{
	curRow,_ := v.toRowCol(kingPos)
	var rookPos, dx,i int
	
	if isKingside{
		rookPos,dx = v.toPos(curRow,v.Files-1),1
		i = kingPos+dx
		if piece,ok := v.pieceLocations[rookPos]; ok && piece.color == color && piece.pieceType == Rook{
			for i<rookPos{
				if _,notEmpty := v.pieceLocations[i]; notEmpty || v.isDisabled(i){ return false}
				i+=dx
			}
			
		}
	} else { 
		rookPos,dx = v.toPos(curRow,0),-1
		i = kingPos+dx
		if piece,ok := v.pieceLocations[rookPos]; ok && piece.color == color && piece.pieceType == Rook{
			for i>rookPos{
				if _,notEmpty := v.pieceLocations[i]; notEmpty || v.isDisabled(i){ return false}
				i+=dx
			}
		}
	}
	return true
}

func (v variant) getPseudoLegalMoves(color color, kingCaptureAllowed bool) []Move {
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
			var props pieceProperties = v.pieceProps[unicode.ToLower(piece.notation)]
			validMoves = append(validMoves, v.genSlideMoves(&piece, currentSquareID, &props)...)
			validMoves = append(validMoves, v.genJumpMoves(&piece, currentSquareID, &props)...)
		}
	}
	return validMoves
}

func (v *variant) genSlideMoves(piece *piece, currentSquareID int, props *pieceProperties) []Move{
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

func (v *variant) genJumpMoves(piece *piece, currentSquareID int, props *pieceProperties) []Move{
	validMoves := []Move{}
	for _, jumpMove := range props.JumpProps {
		target, valid := v.getTargetSquare(currentSquareID, jumpMove.Offset)
		if valid {
			if v.position.isSameColorPiecePresent(currentSquareID, target) {
				continue
			}

			var moveType classicMoveType = NullMove
			if jumpMove.IsCaptureAllowed && v.position.isOpponentPiecePresent(currentSquareID, target) {
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
	case QuietMove, CaptureMove:
		if move.ClassicMoveType == CaptureMove {
			v.recentCapture = v.pieceLocations[move.Target]
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
	}
}

func (v *variant) unmakeMove(move Move) {
	switch move.ClassicMoveType {
	case QuietMove, CaptureMove:
		v.pieceLocations[move.Source] = v.pieceLocations[move.Target]
		if move.PieceType==King{
			if move.Turn == ColorBlack{
				v.additionalProps.blackKingPos = move.Source
			} else {
				v.additionalProps.whiteKingPos = move.Source
			}
		}
		if move.ClassicMoveType == CaptureMove {
			v.pieceLocations[move.Target] = v.recentCapture
		} else { delete(v.pieceLocations,move.Target)}
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

func (cv *CheckmateVariant) isKingUnderCheck(color color) bool {
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

func (cv *CheckmateVariant) IsGameOver() (result, bool) {
	return 0, false
}

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
		//check king in check
		ncv.unmakeMove(mv)
	}
	return []Move{}
}

func (ncv *NCheckVariant) IsGameOver() (result, bool) {
	//check if opponents king received nchecks after making the move
	if ncv.turn == ColorBlack && ncv.whiteKingCheckCount == ncv.targetChecks {
		return BlackWins, true
	} else if ncv.blackKingCheckCount == ncv.targetChecks {
		return WhiteWins, true
	}

	//check for checkmate
	return 0, false
}

type AntichessVariant struct {
	variant
}

func (av *AntichessVariant) GetLegalMoves() []Move {
	pseudoMoves := av.getPseudoLegalMoves(av.turn, true)
	var captureMoves []Move = []Move{}
	for _, move := range pseudoMoves {
		if move.ClassicMoveType == CaptureMove {
			captureMoves = append(captureMoves, move)
		}
	}

	if len(captureMoves) > 0 {
		pseudoMoves = captureMoves
	}
	return pseudoMoves
}

func (av *AntichessVariant) IsGameOver() (result, bool) {

	return 0, false
}
