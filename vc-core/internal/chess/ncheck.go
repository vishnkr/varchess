package chesscore

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