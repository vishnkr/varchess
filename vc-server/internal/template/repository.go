package template

import (
	"context"
	"varchess/internal/db"
)

type Template struct {
	ID           int32 	`json:"id"`
	Name string	`json:"name"`
	GameTemplate []byte `json:"game_template"`
	UserID       string `json:"user_id"`
}

type Repository interface {
	CreateTemplate(ctx context.Context, template Template) error
	GetTemplate(ctx context.Context, templateID int32) (Template,error)
	UpdateTemplate(ctx context.Context, template Template) error
	DeleteTemplate(ctx context.Context, templateID int32) error
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

func (r *repository) CreateTemplate(ctx context.Context, template Template) error{
	arg := db.CreateTemplateParams{
		TemplateName: template.Name,
		GameTemplate: template.GameTemplate,
		UserID: template.UserID,
	}
	return r.q.CreateTemplate(ctx,arg)
}

func (r *repository) GetTemplate(ctx context.Context, templateID int32) (Template,error){
	return Template{},nil
}
func (r *repository) UpdateTemplate(ctx context.Context, template Template) error{
	return nil
}

func (r *repository) DeleteTemplate(ctx context.Context, templateID int32) error{
	return nil
}