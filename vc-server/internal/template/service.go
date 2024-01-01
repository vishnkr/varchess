package template

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
)

type Service interface{

}

type service struct{
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

type CreateTemplateInput struct{
	ID           int32 	`json:"id"`
	TemplateName string	`json:"template_name" validate:"required,min=1"`
	GameTemplate []byte `json:"game_template"`
	UserID       string `json:"user_id"`
}
func (s *service) CreateTemplate(ctx context.Context, template_input CreateTemplateInput)error{
	//logger := logger.FromContext(ctx)
	validate := validator.New()
	err:=validate.Struct(template_input)
	if err!=nil{
		log.Printf("service: invalid template name")
	}
	return s.repository.CreateTemplate(ctx,Template(template_input))
}

