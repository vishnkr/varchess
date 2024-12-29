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
	ValidateSession(ctx context.Context, sessionID string) (string,error)
}

func (s *service) ValidateSession(ctx context.Context, sessionId string) (string,error){
	return s.repository.GetUserIDFromSessionID(ctx,sessionId)
}



