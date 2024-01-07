package game

import (
	"context"

	"github.com/google/uuid"
)



type Service interface{
	CreateGame(ctx context.Context, input CreateGameInput) (Game, error)
}

type service struct{
	repository Repository
}

type CreateGameInput struct {
	ID uuid.UUID
}

func NewService(repository Repository) *service{
	return &service{repository}
}

func (s *service) CreateGame(ctx context.Context, input CreateGameInput) (Game, error){
	return Game{},nil
}