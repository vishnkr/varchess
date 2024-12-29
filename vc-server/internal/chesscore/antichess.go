package chesscore

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