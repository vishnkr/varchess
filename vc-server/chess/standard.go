package chess

import (
	"fmt"
	"unicode"
)

// StandardEngine implements Engine for all chess variants.
// Variant-specific behaviour is delegated to the embedded VariantRules.
type StandardEngine struct {
	variant
}

// --- Engine interface ---

func (se *StandardEngine) GetPosition() *Position {
	return se.Position
}

func (se *StandardEngine) VariantState() interface{} {
	if se.rules == nil {
		return nil
	}
	return se.rules.VariantState()
}

func (se *StandardEngine) LoadVariantState(s interface{}) error {
	if se.rules == nil || s == nil {
		return nil
	}
	return se.rules.LoadVariantState(s)
}

func (se *StandardEngine) IsGameOver() (bool, EngineResult) {
	if se.isGameOverBool && se.gameResult != nil {
		return true, *se.gameResult
	}
	return false, EngineResult{}
}

func (se *StandardEngine) Clone() Engine {
	cloned := *se
	if se.rules != nil {
		cloned.rules = se.rules.Clone()
	}
	// Deep-copy position bitboards.
	pos := *se.Position
	pieces := make(map[rune]Bitboard, len(se.Pieces))
	for k, v := range se.Pieces {
		pieces[k] = v.Clone()
	}
	pos.Pieces = pieces
	colorBBs := make(map[Color]Bitboard, 2)
	for k, v := range se.ColorBitboards {
		colorBBs[k] = v.Clone()
	}
	pos.ColorBitboards = colorBBs
	pos.PositionBitBoard = se.PositionBitBoard.Clone()
	pos.Walls = se.Walls.Clone()
	cloned.Position = &pos
	return &cloned
}

// --- Public helpers ---

func (se *StandardEngine) IsCheckmate() bool {
	return se.IsCheck() && len(se.GetLegalMoves()) == 0
}

func (se *StandardEngine) IsStalemate() bool {
	return !se.IsCheck() && len(se.GetLegalMoves()) == 0
}

// IsCheck reports whether the side-to-move's king is currently in check.
func (se *StandardEngine) IsCheck() bool {
	return se.isColorInCheck(se.CurrentTurn())
}

// isColorInCheck checks whether the given color's king is attacked.
func (se *StandardEngine) isColorInCheck(color Color) bool {
	kingPiece := rune('K')
	if color == Black {
		kingPiece = 'k'
	}
	kingBB, exists := se.Position.Pieces[kingPiece]
	if !exists {
		return false
	}
	kings := kingBB.GetSetBits()
	if len(kings) == 0 {
		return false
	}
	return IsAttacked(kings[0], opponent(color), se.Position)
}

// GetLegalMoves returns all fully-legal moves for the side to move,
// including variant-specific extra/filtered moves.
func (se *StandardEngine) GetLegalMoves() []Move {
	// Antichess treats the king as a normal capturable piece.
	canCaptureKing := relaxesKingSafety(se.rules)
	pseudoMoves := se.GetPseudoLegalMoves(se.CurrentTurn(), canCaptureKing)

	// Variant may add extra moves (archer shots, teleports, etc.).
	if se.rules != nil {
		pseudoMoves = append(pseudoMoves, se.rules.ExtraMoves(se.Position, se.CurrentTurn())...)
	}

	// Variant may filter moves (e.g., force captures in Antichess).
	if se.rules != nil {
		pseudoMoves = se.rules.FilterMoves(pseudoMoves, se.Position)
	}

	// Antichess (and similar) ignore king safety — moving into check is allowed.
	if canCaptureKing {
		return pseudoMoves
	}

	// movedColor is the player making each pseudo-legal move.
	movedColor := se.CurrentTurn()

	legalMoves := make([]Move, 0, len(pseudoMoves))
	for _, move := range pseudoMoves {
		// Castling: king must not pass through or start in check.
		if move.ClassicMoveType == CastleMove {
			if se.isColorInCheck(movedColor) {
				continue // can't castle while in check
			}
			lbd := se.Position.LargestDimension
			dir := 1
			if move.To < move.From {
				dir = -1
			}
			fromFile, midRank := FileRankFromIndex(move.From, lbd)
			midFile := fromFile + dir
			midSq := FileRankToLargeIndex(midFile, midRank, se.Position.Files, lbd)
			if IsAttacked(midSq, opponent(movedColor), se.Position) {
				continue
			}
		}
		se.MakeMove(move)
		// After MakeMove the turn has switched; check if the MOVER's king is exposed.
		inCheck := se.isColorInCheck(movedColor)
		se.UnmakeMove(move)
		if !inCheck {
			legalMoves = append(legalMoves, move)
		}
	}
	return legalMoves
}

// kingSafetyRelaxed is implemented by variants (e.g. antichess) where leaving
// the king in check is legal.
type kingSafetyRelaxed interface {
	RelaxesKingSafety() bool
}

func relaxesKingSafety(rules VariantRules) bool {
	if rules == nil {
		return false
	}
	if r, ok := rules.(kingSafetyRelaxed); ok {
		return r.RelaxesKingSafety()
	}
	return false
}

// IsLegalMove checks if move is among the legal moves for the current position.
// Compares core fields only so variant metadata (VariantMoveType / AdditionalData)
// on generated moves does not reject client payloads that omit them.
func (se *StandardEngine) IsLegalMove(move Move) bool {
	return se.resolveLegalMove(move) != nil
}

func (se *StandardEngine) resolveLegalMove(move Move) *Move {
	for _, lm := range se.GetLegalMoves() {
		if lm.From != move.From || lm.Piece != move.Piece || lm.Promotion != move.Promotion {
			continue
		}
		if lm.To == move.To {
			matched := lm
			return &matched
		}
		// Wormhole: clients click the portal entrance; legal To is the exit.
		if portal, ok := teleportPortalFromAdditional(lm); ok && portal == move.To {
			matched := lm
			return &matched
		}
	}
	return nil
}

func teleportPortalFromAdditional(m Move) (int, bool) {
	if m.VariantMoveType != "teleport" {
		return 0, false
	}
	data, ok := m.AdditionalData.(map[string]interface{})
	if !ok {
		return 0, false
	}
	switch n := data["portal"].(type) {
	case int:
		return n, true
	case float64:
		return int(n), true
	case int32:
		return int(n), true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}

// PerformMove validates and executes a move, updating all game state.
func (se *StandardEngine) PerformMove(move Move) (MoveResult, error) {
	legal := se.resolveLegalMove(move)
	if legal == nil {
		return MoveResult{}, fmt.Errorf("illegal move: %+v", move)
	}
	move = *legal

	if se.rules != nil {
		if err := se.rules.OnBeforeMove(move, se.Position); err != nil {
			return MoveResult{}, err
		}
	}

	se.MakeMove(move)

	result := MoveResult{
		IsCapture: move.Capture,
	}

	// Switch turn before checking (IsCheck checks the NEW side-to-move).
	result.IsCheck = se.IsCheck()

	if se.rules != nil {
		se.rules.OnAfterMove(move, se.Position, &result)
	}

	// OnAfterMove may have already ended the game (e.g. N-Check reaching targetN).
	// Honour that before falling through to CheckTermination / standard rules.
	if result.IsGameOver && result.Result != nil {
		se.isGameOverBool = true
		se.gameResult = result.Result
		return result, nil
	}

	// Check variant termination first, then standard.
	var gameResult *EngineResult
	if se.rules != nil {
		gameResult = se.rules.CheckTermination(se.Position)
	}
	if gameResult == nil {
		gameResult = se.checkStandardTermination()
	}

	if gameResult != nil {
		result.IsGameOver = true
		result.Result = gameResult
		se.isGameOverBool = true
		se.gameResult = gameResult
	}

	return result, nil
}

// checkStandardTermination handles checkmate and stalemate.
func (se *StandardEngine) checkStandardTermination() *EngineResult {
	if len(se.GetLegalMoves()) == 0 {
		if se.IsCheck() {
			// The side that just MOVED wins (opponent of current turn).
			winner := "white"
			if se.CurrentTurn() == White {
				winner = "black"
			}
			return &EngineResult{Winner: winner, Reason: "checkmate"}
		}
		return &EngineResult{Winner: "draw", Reason: "stalemate"}
	}
	return nil
}

// --- MakeMove / UnmakeMove ---

// MakeMove executes move on the position. Does not validate legality.
func (se *StandardEngine) MakeMove(move Move) {
	piece := move.Piece
	from, to := move.From, move.To
	lbd := se.Position.LargestDimension
	color := se.CurrentTurn()

	capturedPiece := rune(0)
	capturedSquare := to

	// En passant: captured pawn is not at the destination.
	if move.ClassicMoveType == EnPassant {
		dir := -1
		if color == Black {
			dir = 1
		}
		// White direction=-1: captured pawn is below target (to + stride).
		// Black direction=+1: captured pawn is above target (to - stride).
		stride := IndexStride(lbd)
		capturedSquare = to - dir*stride
		capturedPiece = se.GetPieceAt(capturedSquare)
	} else {
		capturedPiece = se.GetPieceAt(to)
	}

	se.recentMove = RecentMoveInfo{
		Piece:          piece,
		moveType:       move.ClassicMoveType,
		capturedPiece:  capturedPiece,
		capturedSquare: capturedSquare,
		prevCastling:   se.Position.Castling,
		prevEnPassant:  se.Position.EnPassant,
		prevHalfMove:   se.Position.HalfMove,
		prevFullMove:   se.Position.FullMove,
		castleRookFrom: -1,
		castleRookTo:   -1,
	}

	// Remove piece from source.
	se.Position.Pieces[piece].ClearBit(from)
	se.Position.ColorBitboards[color].ClearBit(from)
	se.Position.PositionBitBoard.ClearBit(from)

	// Remove captured piece.
	if capturedPiece != 0 {
		se.Position.Pieces[capturedPiece].ClearBit(capturedSquare)
		se.Position.ColorBitboards[opponent(color)].ClearBit(capturedSquare)
		se.Position.PositionBitBoard.ClearBit(capturedSquare)
	}

	// Place piece at destination (use promoted piece if applicable).
	landingPiece := piece
	if move.Promotion != 0 {
		landingPiece = move.Promotion
		if _, ok := se.Position.Pieces[landingPiece]; !ok {
			se.Position.Pieces[landingPiece] = NewBitboard(lbd)
		}
	}
	se.Position.Pieces[landingPiece].SetBit(to)
	se.Position.ColorBitboards[color].SetBit(to)
	se.Position.PositionBitBoard.SetBit(to)

	// Castling: also move the rook (skip if no rook bitboard — custom/miniboards).
	if move.ClassicMoveType == CastleMove {
		fromFile, rank := FileRankFromIndex(from, lbd)
		toFile, _ := FileRankFromIndex(to, lbd)
		var rookFrom, rookTo int
		rookPiece := rune('R')
		if color == Black {
			rookPiece = 'r'
		}
		if toFile > fromFile { // king-side
			rookFrom = FileRankToLargeIndex(se.Position.Files-1, rank, se.Position.Files, lbd)
			rookTo = FileRankToLargeIndex(toFile-1, rank, se.Position.Files, lbd)
		} else { // queen-side
			rookFrom = FileRankToLargeIndex(0, rank, se.Position.Files, lbd)
			rookTo = FileRankToLargeIndex(toFile+1, rank, se.Position.Files, lbd)
		}
		if rookBB, ok := se.Position.Pieces[rookPiece]; ok && rookBB != nil {
			rookBB.ClearBit(rookFrom)
			se.Position.ColorBitboards[color].ClearBit(rookFrom)
			se.Position.PositionBitBoard.ClearBit(rookFrom)
			rookBB.SetBit(rookTo)
			se.Position.ColorBitboards[color].SetBit(rookTo)
			se.Position.PositionBitBoard.SetBit(rookTo)
			se.recentMove.castleRookFrom = rookFrom
			se.recentMove.castleRookTo = rookTo
		}
	}

	updateGameState(se, move)
	se.SwitchTurn()
}

// UnmakeMove reverses the last MakeMove call.
// After MakeMove, the turn has been switched, so CurrentTurn() is the opponent of the mover.
func (se *StandardEngine) UnmakeMove(move Move) {
	piece := move.Piece
	from, to := move.From, move.To
	// The side that made the move is the opponent of the current turn.
	movedColor := opponent(se.CurrentTurn())
	lbd := se.Position.LargestDimension

	// Remove piece from destination (handle promotion).
	removePiece := piece
	if move.Promotion != 0 {
		removePiece = move.Promotion
	}
	se.Position.Pieces[removePiece].ClearBit(to)
	se.Position.ColorBitboards[movedColor].ClearBit(to)
	se.Position.PositionBitBoard.ClearBit(to)

	// Restore piece at source.
	se.Position.Pieces[piece].SetBit(from)
	se.Position.ColorBitboards[movedColor].SetBit(from)
	se.Position.PositionBitBoard.SetBit(from)

	// Restore captured piece.
	if se.recentMove.capturedPiece != 0 {
		cp := se.recentMove.capturedPiece
		cs := se.recentMove.capturedSquare
		se.Position.Pieces[cp].SetBit(cs)
		// Captured piece belongs to CurrentTurn() (the side that did NOT move).
		se.Position.ColorBitboards[se.CurrentTurn()].SetBit(cs)
		se.Position.PositionBitBoard.SetBit(cs)
	}

	// Undo castling: move rook back.
	if move.ClassicMoveType == CastleMove && se.recentMove.castleRookFrom != -1 {
		rookPiece := rune('R')
		if movedColor == Black {
			rookPiece = 'r'
		}
		_ = lbd
		if rookBB, ok := se.Position.Pieces[rookPiece]; ok && rookBB != nil {
			rookBB.ClearBit(se.recentMove.castleRookTo)
			se.Position.ColorBitboards[movedColor].ClearBit(se.recentMove.castleRookTo)
			se.Position.PositionBitBoard.ClearBit(se.recentMove.castleRookTo)
			rookBB.SetBit(se.recentMove.castleRookFrom)
			se.Position.ColorBitboards[movedColor].SetBit(se.recentMove.castleRookFrom)
			se.Position.PositionBitBoard.SetBit(se.recentMove.castleRookFrom)
		}
	}

	// Restore position state.
	se.Position.Castling = se.recentMove.prevCastling
	se.Position.EnPassant = se.recentMove.prevEnPassant
	se.Position.HalfMove = se.recentMove.prevHalfMove
	se.Position.FullMove = se.recentMove.prevFullMove

	se.SwitchTurn()
}

// GetPseudoLegalMoves generates all pseudo-legal moves for color (before king-safety filtering).
func (se *StandardEngine) GetPseudoLegalMoves(color Color, canCaptureKing bool) []Move {
	moves := []Move{}
	for piece, bitboard := range se.Pieces {
		isWhite := unicode.IsUpper(piece)
		if (color == White && !isWhite) || (color == Black && isWhite) {
			continue
		}
		for _, pos := range bitboard.GetSetBits() {
			pieceMoves := GenerateMovesForPiece(piece, pos, se.Position, canCaptureKing)
			moves = append(moves, pieceMoves...)
		}
	}
	return moves
}

// updateGameState updates castling rights, en-passant square, and move counters.
func updateGameState(se *StandardEngine, move Move) {
	lbd := se.Position.LargestDimension
	color := se.CurrentTurn()

	// Update castling rights when king or rook moves.
	switch unicode.ToLower(move.Piece) {
	case 'k':
		if color == White {
			se.Position.Castling &^= (1 << KSCW) | (1 << QSCW)
		} else {
			se.Position.Castling &^= (1 << KSCB) | (1 << QSCB)
		}
	case 'r':
		fromFile, fromRank := FileRankFromIndex(move.From, lbd)
		if color == White {
			// White rooks start at rank Ranks-1 (bottom of board in rank-from-top).
			if fromRank == se.Position.Ranks-1 {
				if fromFile == se.Position.Files-1 {
					se.Position.Castling &^= 1 << KSCW
				} else if fromFile == 0 {
					se.Position.Castling &^= 1 << QSCW
				}
			}
		} else {
			// Black rooks start at rank 0 (top).
			if fromRank == 0 {
				if fromFile == se.Position.Files-1 {
					se.Position.Castling &^= 1 << KSCB
				} else if fromFile == 0 {
					se.Position.Castling &^= 1 << QSCB
				}
			}
		}
	}

	// En passant square: set after a double pawn push.
	if move.ClassicMoveType == DoublePawnPush {
		// The en passant target is the square the capturing pawn would land on,
		// which is the midpoint between from and to (works for any board size
		// because the large-board index encodes rank*lbd+file).
		se.Position.EnPassant = (move.From + move.To) / 2
	} else {
		se.Position.EnPassant = -1
	}

	// Half-move clock.
	if unicode.ToLower(move.Piece) == 'p' || move.Capture {
		se.Position.HalfMove = 0
	} else {
		se.Position.HalfMove++
	}

	// Full-move counter increments after Black moves.
	if color == Black {
		se.Position.FullMove++
	}
}

// --- StandardVariant alias for backward compatibility ---

// StandardVariant is an alias kept so existing tests compile unchanged.
type StandardVariant = StandardEngine
