package game

type VariantType uint8

const (
	Custom VariantType = iota
	DuckChess
	PoisonedPawn
	Wormhole
)

type VariantObjective interface{
	getPseudoLegalMoves() Moves
	GetLegalMoves() Moves
	IsGameOver() bool
}

type Variant interface{
	makeMove(Move)
	unmakeMove(Move)
}

type CustomVariant struct{
	Objective Objective
	Position Position
	VariantType VariantType
	recentCapture Piece
}

type CheckmateVariant struct{
	CustomVariant
}

type Moves struct{
	validMoves []Move
	attackedSquares map[int]bool
}

func (v CheckmateVariant) getPseudoLegalMoves() Moves{
	validMoves := []Move{}
	attackedSquares := make(map[int]bool)
	for currentSquareID,piece:=range v.Position.PieceLocations{
		if(piece.Color == v.Position.Turn){
			var props PieceProps = v.Position.PieceProps[piece]
			slideMoves := v.genSlideMoves(&piece,currentSquareID,&props,&attackedSquares)
			jumpMoves := v.genJumpMoves(&piece,currentSquareID,&props,&attackedSquares)
			if(v.Position.PieceProps[piece].CanDoubleJump){
				if _,ok:= v.Position.PieceProps[piece].DoubleJumpSquares[currentSquareID];ok{
					var doubleJumpMoves []Move = []Move{}
					for _,mv:=range jumpMoves{
						if mv.ClassicMoveType == QuietMove{
							doubleJumpMoves = append(doubleJumpMoves, v.genJumpMoves(&piece,mv.Target,&props,&attackedSquares)...)
						}
					}
					validMoves = append(validMoves, doubleJumpMoves...)
				}
			}
			specialMoves:= v.genSpecialMoves()
			validMoves = append(validMoves, slideMoves...)
			jumpMoves = append(jumpMoves, slideMoves...)
			specialMoves = append(specialMoves, slideMoves...)
		}
		
	}
	return Moves{validMoves,attackedSquares}
}

func (v *CustomVariant) getTargetSquare(currentSquareID int, offset MoveOffset) (int,bool){
	target:= currentSquareID + (v.Position.Dimensions.Ranks * offset.YOffset) + offset.YOffset
	if target<0 || target >= (v.Position.Dimensions.Ranks*v.Position.Dimensions.Files){
		return -1,false
	}
	return target,true
}

func (p *Position) isSameColorPiecePresent(sourceSquareID int, targetSquareID int) bool{
	//assuming both squareIds are vald
	return p.PieceLocations[sourceSquareID].Color == p.PieceLocations[targetSquareID].Color
}

func (p *Position) isOpponentPiecePresent(sourceSquareID int,targetSquareID int) bool{
	//assuming both squareIds are vald
	return p.PieceLocations[sourceSquareID].Color != p.PieceLocations[targetSquareID].Color
}

func (p *Position) isNotDisabled(targetSquareID int) bool{
	_,ok:= p.DisabledSquares[targetSquareID]
	return ok
}
func (p *Position) getClassicMoveType(targetSquareID int) ClassicMoveType{
	if _,ok:= p.PieceLocations[targetSquareID]; ok{
		return CaptureMove
	}
	return QuietMove
	
}

func (v *CustomVariant) genSlideMoves(piece *Piece,currentSquareID int,props *PieceProps,attackedSquares *map[int]bool) []Move{
	validMoves := []Move{}
	for moveOffset:= range props.SlideOffsets{
		target,valid := v.getTargetSquare(currentSquareID,moveOffset)
		for valid{
			if(!v.Position.isSameColorPiecePresent(currentSquareID,target)){
				validMoves = append(validMoves, Move{
					Source: currentSquareID,
					Target: target,
					ClassicMoveType: v.Position.getClassicMoveType(target),
					PieceType: piece.Piecetype,
					PieceName: piece.Name,
					Turn: piece.Color,
				})
			}
			target,valid = v.getTargetSquare(target,moveOffset)
		}
	}
	return validMoves
}

func (v *CustomVariant) genJumpMoves(piece *Piece,currentSquareID int,props *PieceProps,attackedSquares *map[int]bool) []Move{
	validMoves := []Move{}
	for moveOffset, allowedJumpMove:= range props.JumpProps{
		target,valid := v.getTargetSquare(currentSquareID,moveOffset)
		if( valid && 
			!v.Position.isSameColorPiecePresent(currentSquareID,target) &&
			v.Position.isNotDisabled(target)){
			var moveType ClassicMoveType = NullMove
			switch allowedJumpMove{
				case CaptureOnlyJump:
					if(v.Position.isOpponentPiecePresent(currentSquareID,target)) {moveType = CaptureMove}
					
				case QuietOnlyJump:
					moveType = QuietMove
				case AllJumps:
					if(v.Position.isOpponentPiecePresent(currentSquareID,target)){
						moveType = CaptureMove
					} else {
						moveType = QuietMove
					}
			}
			move := Move {
				Source: currentSquareID,
				Target: target,
				ClassicMoveType: moveType,
				PieceType: piece.Piecetype,
				PieceName: piece.Name,
				Turn: piece.Color,
			}
			validMoves = append(validMoves, move)
		}
	}
	return validMoves
}

func (v *CustomVariant) genSpecialMoves() []Move{
	validMoves := []Move{}
	return validMoves
}

func (v CheckmateVariant) GetLegalMoves() Moves{
	pseudoMoves:=v.getPseudoLegalMoves()
	for _,mv:=range pseudoMoves.validMoves{
		v.makeMove(mv)
		//check king in check
		v.unmakeMove(mv)
	} 
 	return Moves{validMoves:[]Move{}}
}


func (v CustomVariant) makeMove(move Move){
	switch move.ClassicMoveType{
	case QuietMove, CaptureMove:
		if move.ClassicMoveType==CaptureMove{
			v.recentCapture = v.Position.PieceLocations[move.Target]
		}
		v.Position.PieceLocations[move.Target] = v.Position.PieceLocations[move.Source]
		delete(v.Position.PieceLocations,move.Source)
	}
}

func (v CustomVariant) unmakeMove(move Move){
	switch move.ClassicMoveType{
	case QuietMove, CaptureMove:
		v.Position.PieceLocations[move.Source] = v.Position.PieceLocations[move.Target]
		if move.ClassicMoveType==CaptureMove{
			v.Position.PieceLocations[move.Target] = v.recentCapture
		}
	}
}

func (v CheckmateVariant) IsGameOver()bool{
	return false
}