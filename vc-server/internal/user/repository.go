package user

import (
	"context"
	"varchess/internal/db"
)


type Repository interface {
	GetUserIDFromSessionID (ctx context.Context, sessionID string) (string,error)
	GetUserFromUserID (ctx context.Context, userID string) (User,error)
}

type repository struct {
	db *db.Database
	q  *db.Queries
}

func NewRepository(conn *db.Database) *repository{
	return &repository{
		db:conn,
		q:db.New(conn),
	}
}

func (r *repository) GetUserFromUserID(ctx context.Context, userID string) (User,error){
	return User{},nil
}

func (r *repository) GetUserIDFromSessionID(ctx context.Context, sessionID string) (string,error){
	userSess,err := r.q.GetUserFromSession(ctx,sessionID)
	return userSess.UserID,err
}

