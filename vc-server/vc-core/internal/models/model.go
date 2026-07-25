package models

import (
	"encoding/json"

	"vc-server/chess"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	UserId primitive.ObjectID `json:"userId" bson:"_id"`
	Name 		string `json:"name" bson:"username"`
}

type Template struct {
	Name 		string `json:"name" bson:"name"`
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId      primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	VariantType string             `json:"variantType" bson:"variantType"`
	Position    Position           `json:"position" bson:"position"`
	CustomData   map[string]interface{}          `json:"customData" bson:"customData"`
}

type Position struct {
	Dimensions     chess.Dimensions                  `json:"dimensions" bson:"dimensions"`
	FEN            string                      `json:"fen,omitempty" bson:"fen,omitempty"`
	PieceProps     map[string]chess.PieceProps       `json:"pieceProps,omitempty" bson:"pieceProps,omitempty"`
	PieceLocations map[string]map[string][]int `json:"pieceLocations,omitempty" bson:"pieceLocations,omitempty"` // Maps color -> piece type -> piece locations
}


type Objective struct {
	Type           string                 `json:"type" bson:"type"`
	ObjectiveProps map[string]interface{} `json:"objectiveProps,omitempty" bson:"objectiveProps,omitempty"`
}

type GameResult struct {
	Winner string `json:"winner" bson:"winner"`
	Reason string `json:"reason" bson:"reason"`
}

type Game struct {
	ID          primitive.ObjectID            `json:"_id,omitempty" bson:"_id,omitempty"`
	ShortID     string                        `json:"shortId,omitempty" bson:"short_id,omitempty"`
	Players     map[string]primitive.ObjectID `json:"players" bson:"players"`
	PlayerNames map[string]string             `json:"playerNames,omitempty" bson:"player_names,omitempty"`
	TemplateID  *primitive.ObjectID           `json:"templateId,omitempty" bson:"template_id,omitempty"`
	Config      chess.GameConfig              `json:"config,omitempty" bson:"config,omitempty"`
	Moves       []string                      `json:"moves" bson:"moves"`
	Result      GameResult                    `json:"result" bson:"result"`
	CreatedAt   primitive.DateTime            `json:"createdAt" bson:"created_at"`
	EndedAt     *primitive.DateTime           `json:"endedAt,omitempty" bson:"ended_at,omitempty"`
}

type ActiveGameState string
const (
	Waiting ActiveGameState = "w"
	InProgress ActiveGameState = "ip"
	Complete ActiveGameState = "c"
)

type ActiveGame struct {
	ID      string            `json:"id"`
	Players map[string]Player `json:"players"`
	Config  chess.GameConfig  `json:"gameConfig"`
	Moves   []string          `json:"moves"`
	State   ActiveGameState   `json:"state"`
	Turn    string            `json:"turn"`
	// VariantState holds live variant-specific state (e.g. wormhole cooldowns)
	// as a JSON object for clients and reconnect snapshots.
	VariantState json.RawMessage `json:"variantState,omitempty"`
}


type PaginatedResponse[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

type User struct{
	UserId string `json:"userId"`
	Name string `json:"name"`
}