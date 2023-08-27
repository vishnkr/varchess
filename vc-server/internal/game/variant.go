package game

import "unicode"

type variantType uint8

const (
	Custom variantType = iota
	DuckChess
	ArcherChess
	Wormhole
)


type Variant interface{
	makeMove(Move)
	unmakeMove(Move)
	getPseudoLegalMoves(color,bool) []Move
	GetLegalMoves() []Move
	IsGameOver() (result,bool)
}

type variant struct{
	Objective Objective `json:"objective"`
	position
	variantType
	recentCapture piece
}

func (v *variant) getTargetSquare(currentSquareID int, offset moveOffset) (int,bool){
	row,col:= v.toRowCol(currentSquareID)
	newRow,newCol := row+offset.yOffset,col+offset.xOffset
	target:= v.toPos(newRow,newCol)
	if newRow<0 || newCol<0 || newRow>=v.Ranks || newCol>=v.Files || v.isDisabled(target) {
		return -1,false
	} 
	return target,true
}

func (v variant) getPseudoLegalMoves(color color,kingCaptureAllowed bool) []Move{
	validMoves := []Move{}
	for currentSquareID,piece:=range v.pieceLocations{
		if color != piece.color { continue }
		if piece.pieceType == Pawn{
			var rowOffset, doubleMoveStartRank int
			if color == ColorWhite {
				doubleMoveStartRank = v.Ranks - 2
				rowOffset = -1
			} else {
				doubleMoveStartRank = 1
				rowOffset = 1
			}
			srcRow,srcCol := v.toRowCol(currentSquareID)
			targetRow := srcRow + rowOffset
			for i := -1; i <= 1; i++ {
				//non-capture moves
				if i == 0 {
					target1, target2 := v.toPos(targetRow,srcCol), v.toPos(targetRow+rowOffset,srcCol)
					if _,ok := v.pieceLocations[target1]; !ok && !v.disabledSquares[target1]{

					}
					if _,ok := v.pieceLocations[target2]; !ok && srcRow == doubleMoveStartRank && !v.disabledSquares[target2]{

					}
				}
				//capture moves
				/*if board.isSquareInBoardRange(targetRow, srcCol+i) && !board.IsEmpty(targetRow, srcCol+i) && !board.isSquareDisabled(targetRow, srcCol+i) && !board.isSameColorPieceAtDest(piece.Color, targetRow, srcCol+i) {
					
				}*/
			}
		} else{
			var props pieceProperties = v.pieceProps[unicode.ToLower(piece.notation)]
			slideMoves := v.genSlideMoves(&piece,currentSquareID,&props,kingCaptureAllowed)
			jumpMoves := v.genJumpMoves(&piece,currentSquareID,&props,kingCaptureAllowed)
			specialMoves:= v.genSpecialMoves()
			validMoves = append(validMoves, slideMoves...)
			validMoves = append(validMoves, jumpMoves...)
			validMoves = append(validMoves, specialMoves...)
		}
			
	}
	return validMoves
}

func (v *variant) genSlideMoves(piece *piece,currentSquareID int,props *pieceProperties,kingCaptureAllowed bool) []Move{
	validMoves := []Move{}
	for _,offset:= range props.SlideOffsets{
		target,valid := v.getTargetSquare(currentSquareID,offset)
		for valid{
			if (!v.isSameColorPiecePresent(currentSquareID,target)){
				if (v.isOppKingAtTarget(target) && !kingCaptureAllowed){
					continue
				}
				v.position.attackedSquares[target] = true
				mType:=v.position.getclassicMoveType(target)
				validMoves = append(validMoves, Move{
					Source: currentSquareID,
					Target: target,
					ClassicMoveType: mType,
					PieceType: piece.pieceType,
					PieceNotation: piece.notation,
					Turn: piece.color,
				})
				if mType!=QuietMove{
					break
				}
			} else {
				break
			}
			target,valid = v.getTargetSquare(target,offset)
		}
	}
	return validMoves
}


func (v *variant) genJumpMoves(piece *piece,currentSquareID int,props *pieceProperties,kingCaptureAllowed bool) []Move{
	validMoves := []Move{}
	for _, jumpMove:= range props.JumpProps{
		target,valid := v.getTargetSquare(currentSquareID,jumpMove.Offset)
		if (valid){
			if (v.position.isSameColorPiecePresent(currentSquareID,target)){ continue}
			
			var moveType classicMoveType = NullMove
			if (jumpMove.IsCaptureAllowed && v.position.isOpponentPiecePresent(currentSquareID,target)){
					if (v.position.isOppKingAtTarget(target) && kingCaptureAllowed){
						moveType = CaptureMove
					} else{ continue }
			} else{
				moveType = QuietMove
			}
			v.position.attackedSquares[target]=true
			move := Move {
				Source: currentSquareID,
				Target: target,
				ClassicMoveType: moveType,
				PieceType: piece.pieceType,
				PieceNotation: piece.notation,
				Turn: piece.color,
			}
			validMoves = append(validMoves, move)
		}
	}
	return validMoves
}

func (v *variant) genSpecialMoves() []Move{
	//castle and en passant
	validMoves := []Move{}
	return validMoves
}



func (v variant) makeMove(move Move){
	switch move.ClassicMoveType{
	case QuietMove, CaptureMove:
		if move.ClassicMoveType==CaptureMove{
			v.recentCapture = v.pieceLocations[move.Target]
		}
		v.pieceLocations[move.Target] = v.pieceLocations[move.Source]
		delete(v.pieceLocations,move.Source)
	}
}

func (v variant) unmakeMove(move Move){
	switch move.ClassicMoveType{
	case QuietMove, CaptureMove:
		v.pieceLocations[move.Source] = v.pieceLocations[move.Target]
		if move.ClassicMoveType==CaptureMove{
			v.pieceLocations[move.Target] = v.recentCapture
		}
	}
}



type CheckmateVariant struct{
	variant
}


func (cv CheckmateVariant) GetLegalMoves() []Move{
	pseudoMoves:=cv.getPseudoLegalMoves(cv.turn,false)
	var validMoves []Move = []Move{}
	for _,mv:=range pseudoMoves{
		cv.makeMove(mv)
		cv.position.switchTurn()
		//check king in check
		if !cv.isKingUnderCheck(cv.turn){
			validMoves = append(validMoves,mv)
		}
		cv.position.switchTurn()
		cv.unmakeMove(mv)
	} 
 	return validMoves
}

func (cv CheckmateVariant) isKingUnderCheck(color color) bool{
	cv.getPseudoLegalMoves(color,false)
	var kingPos int
	if color==ColorWhite{
		kingPos = cv.kingPos.whiteKingPos
	}  else {
		kingPos = cv.kingPos.blackKingPos
	}
	_, ok := cv.attackedSquares[kingPos] 
	return ok
}

func (cv CheckmateVariant) IsGameOver() (result,bool){
	return 0,false
}

type NCheckVariant struct{
	variant
	blackKingCheckCount int
	whiteKingCheckCount int
	targetChecks int
}

func (ncv NCheckVariant) GetLegalMoves() []Move{
	pseudoMoves:=ncv.getPseudoLegalMoves(ncv.turn,false)
	for _,mv:=range pseudoMoves{
		ncv.makeMove(mv)
		//check king in check
		ncv.unmakeMove(mv)
	} 
	return []Move{}
}

func (ncv NCheckVariant) IsGameOver() (result,bool){
	//check if opponents king received nchecks after making the move
	if ncv.turn == ColorBlack && ncv.whiteKingCheckCount == ncv.targetChecks{
		return BlackWins, true
	} else if ncv.blackKingCheckCount == ncv.targetChecks {
		return WhiteWins,true
	}

	//check for checkmate
	return 0, false
}

type AntichessVariant struct{
	variant
}

func (av AntichessVariant) GetLegalMoves() []Move{
	pseudoMoves:= av.getPseudoLegalMoves(av.turn,true)
	var captureMoves []Move = []Move{}
	for _,move := range pseudoMoves{
		if move.ClassicMoveType == CaptureMove{
			captureMoves = append(captureMoves, move)
		} 
	}

	if len(captureMoves)>0{
		pseudoMoves = captureMoves
	}
	return pseudoMoves
}

func (av AntichessVariant) IsGameOver()(result,bool){
	
	return 0,false
}