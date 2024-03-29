// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package db

import (
	"context"
)

const createGame = `-- name: CreateGame :exec
INSERT INTO saved_game
(current_state, game_template, player1, player2) 
VALUES ($1,$2,$3,$4)
`

type CreateGameParams struct {
	CurrentState []byte
	GameTemplate []byte
	Player1      string
	Player2      string
}

func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) error {
	_, err := q.db.Exec(ctx, createGame,
		arg.CurrentState,
		arg.GameTemplate,
		arg.Player1,
		arg.Player2,
	)
	return err
}

const createTemplate = `-- name: CreateTemplate :exec
INSERT INTO template
(template_name, game_template, user_id) 
VALUES ($1,$2,$3)
`

type CreateTemplateParams struct {
	TemplateName string
	GameTemplate []byte
	UserID       string
}

func (q *Queries) CreateTemplate(ctx context.Context, arg CreateTemplateParams) error {
	_, err := q.db.Exec(ctx, createTemplate, arg.TemplateName, arg.GameTemplate, arg.UserID)
	return err
}

const deleteGame = `-- name: DeleteGame :exec
DELETE FROM saved_game 
WHERE id = $1
`

func (q *Queries) DeleteGame(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteGame, id)
	return err
}

const deleteTemplate = `-- name: DeleteTemplate :exec
DELETE FROM template 
WHERE id = $1
`

func (q *Queries) DeleteTemplate(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteTemplate, id)
	return err
}

const getTemplate = `-- name: GetTemplate :one
SELECT id, template_name, game_template, user_id, last_updated FROM template 
WHERE template.id = $1
`

func (q *Queries) GetTemplate(ctx context.Context, id int32) (Template, error) {
	row := q.db.QueryRow(ctx, getTemplate, id)
	var i Template
	err := row.Scan(
		&i.ID,
		&i.TemplateName,
		&i.GameTemplate,
		&i.UserID,
		&i.LastUpdated,
	)
	return i, err
}

const listGames = `-- name: ListGames :many
SELECT id FROM saved_game 
WHERE saved_game.player1 = $1 OR saved_game.player2 = $1
ORDER BY saved_game.last_updated DESC
`

func (q *Queries) ListGames(ctx context.Context, player1 string) ([]int32, error) {
	rows, err := q.db.Query(ctx, listGames, player1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTemplates = `-- name: ListTemplates :many
SELECT id,template_name FROM template 
WHERE template.user_id = $1
ORDER BY template.last_updated DESC
`

type ListTemplatesRow struct {
	ID           int32
	TemplateName string
}

func (q *Queries) ListTemplates(ctx context.Context, userID string) ([]ListTemplatesRow, error) {
	rows, err := q.db.Query(ctx, listTemplates, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListTemplatesRow
	for rows.Next() {
		var i ListTemplatesRow
		if err := rows.Scan(&i.ID, &i.TemplateName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateGame = `-- name: UpdateGame :exec
UPDATE saved_game SET (current_state) = ($2)
WHERE id = $1
`

func (q *Queries) UpdateGame(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, updateGame, id)
	return err
}

const updateTemplate = `-- name: UpdateTemplate :exec
UPDATE template SET (id, game_template) = ($1,$2)
WHERE id = $1
`

type UpdateTemplateParams struct {
	ID           int32
	GameTemplate []byte
}

func (q *Queries) UpdateTemplate(ctx context.Context, arg UpdateTemplateParams) error {
	_, err := q.db.Exec(ctx, updateTemplate, arg.ID, arg.GameTemplate)
	return err
}
