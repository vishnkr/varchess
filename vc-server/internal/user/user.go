package user

import "context"

type User struct{
	ID string `json:"id"`
	Username string `json:"username"`
}

type service struct{
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

type Service interface{
	ValidateSession(ctx context.Context, sessionID string, userID string) error
}


