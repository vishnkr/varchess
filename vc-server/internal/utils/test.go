package utils

import (
	"varchess/internal/template"
	"varchess/internal/user"
)


type mockTemplatesRepo struct{
	templates map[int32]template.Template
	users map[string]user.User
}

func NewMockTemplatesRepo() *mockTemplatesRepo{
	templates := make(map[int32]template.Template)
	users := make(map[string]user.User)
	return &mockTemplatesRepo{
		templates,
		users,
	}
}