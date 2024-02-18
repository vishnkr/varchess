package template

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
)

type Service interface{
	CreateTemplate(context.Context, CreateTemplateInput) error
}

type service struct{
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

type CreateTemplateInput struct{
	ID           int32 	`json:"id"`
	Name 		 string	`json:"name" validate:"required,min=1"`
	GameTemplate []byte `json:"game_template"`
	UserID       string `json:"user_id"`
}

func (s *service) CreateTemplate(ctx context.Context, templateInput CreateTemplateInput)error{
	//logger := logger.FromContext(ctx)
	validate := validator.New()
	err:=validate.Struct(templateInput)
	if err!=nil{
		log.Printf("service: invalid template name")
	}
	return s.repository.CreateTemplate(ctx,Template(templateInput))
}

func (s *service) GetTemplatesByUser(ctx context.Context, userId string) error{
	return nil
}