// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthUser struct {
	ID          string
	Username    string
	Email       pgtype.Text
	LastUpdated pgtype.Timestamp
}

type SavedGame struct {
	ID           int32
	CurrentState []byte
	GameTemplate []byte
	Player1      string
	Player2      string
	LastUpdated  pgtype.Timestamp
}

type Template struct {
	ID           int32
	TemplateName string
	GameTemplate []byte
	UserID       string
	LastUpdated  pgtype.Timestamp
}

type UserKey struct {
	ID             string
	UserID         string
	HashedPassword pgtype.Text
}

type UserSession struct {
	ID            string
	UserID        string
	ActiveExpires int64
	IdleExpires   int64
}
