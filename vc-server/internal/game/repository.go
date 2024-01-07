package game

import (
	"context"
	"fmt"
	"varchess/internal/db"

	"github.com/google/uuid"
)

type Game struct {
	ID           uuid.UUID 	`json:"id"`
	CurrentState string	`json:"current_state"`
	GameTemplate string `json:"game_template"`
	Player1       string `json:"player1"`
	Player2       string `json:"playr2"`
}

type Repository interface {
	CreateGame(ctx context.Context, game Game) error
	UpdateGame(ctx context.Context, game Game) error
	DeleteGame(ctx context.Context, gameID int32) error
}

type repository struct {
	db *db.Database
	q  *db.Queries
}

func NewRepository(conn *db.Database) *repository {
	q := db.New(conn)
	return &repository{conn, q}
}

func (r *repository) CreateGame(ctx context.Context, game Game) error{
	arg := db.CreateGameParams{
		CurrentState: []byte(game.CurrentState),
		GameTemplate: []byte(game.GameTemplate),
		Player1: game.Player1,
		Player2: game.Player2,
	}
	err := r.q.CreateGame(ctx,arg)
	if err!=nil{
		return fmt.Errorf("repository: failed to save game: %w", err)
	}
	return nil
}

func (r *repository) UpdateGame(ctx context.Context, game Game) error{ return nil}
func (r *repository) DeleteGame(ctx context.Context, gameID int32) error{ return nil }