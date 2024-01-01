package user

import (
	"context"
	"varchess/internal/db"
)


type Repository interface {
	GetUser (ctx context.Context, userID string) (User,error)
}

type repository struct {
	db *db.Database
	q  *db.Queries
}

func (r *repository) GetUser(ctx context.Context, userID string) (User,error){

	return User{},nil
}