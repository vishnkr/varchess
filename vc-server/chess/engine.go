package chess

// EngineResult describes the outcome of a finished game.
type EngineResult struct {
	Winner string `json:"winner"` // "white", "black", or "draw"
	Reason string `json:"reason"` // "checkmate", "stalemate", "no-pieces", "n-checks", etc.
}

// Engine is the primary interface for a chess game instance.
// All variant implementations satisfy this interface.
type Engine interface {
	// GetPosition returns the current board position (read-only).
	GetPosition() *Position

	// GetLegalMoves returns all fully-legal moves for the side to move.
	GetLegalMoves() []Move

	// IsLegalMove reports whether move is legal in the current position.
	IsLegalMove(move Move) bool

	// PerformMove validates and executes a move, updating game state.
	PerformMove(move Move) (MoveResult, error)

	// IsGameOver reports whether the game has ended and the result.
	IsGameOver() (bool, EngineResult)

	// Clone returns a deep copy of the engine (safe for speculative search).
	Clone() Engine

	// VariantState returns serialisable variant-specific state (or nil).
	VariantState() interface{}

	// LoadVariantState restores variant-specific state from a deserialised value.
	LoadVariantState(s interface{}) error
}
