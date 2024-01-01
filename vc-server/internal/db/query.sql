-- name: CreateTemplate :exec
INSERT INTO template
(template_name, game_template, user_id) 
VALUES ($1,$2,$3);

-- name: UpdateTemplate :exec
UPDATE template SET (id, game_template) = ($1,$2)
WHERE id = $1;

-- name: DeleteTemplate :exec
DELETE FROM template 
WHERE id = $1;

-- name: ListTemplates :many
SELECT id,template_name FROM template 
WHERE template.user_id = $1
ORDER BY template.last_updated DESC;

-- name: GetTemplate :one
SELECT * FROM template 
WHERE template.id = $1;

-- name: CreateGame :exec
INSERT INTO saved_game
(current_state, game_template, player1, player2) 
VALUES ($1,$2,$3,$4);

-- name: ListGames :many
SELECT id FROM saved_game 
WHERE saved_game.player1 = $1 OR saved_game.player2 = $1
ORDER BY saved_game.last_updated DESC;

-- name: UpdateGame :exec
UPDATE saved_game SET (current_state) = ($2)
WHERE id = $1;

-- name: DeleteGame :exec
DELETE FROM saved_game 
WHERE id = $1;