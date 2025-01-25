package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
	Color  string             `json:"color" bson:"color"`
}

type Template struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId      primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	VariantType string             `json:"variantType" bson:"variantType"`
	Position    Position           `json:"position" bson:"position"`
	Objective   Objective          `json:"objective" bson:"objective"`
}

type GameConfig struct {
	VariantType string    `json:"variantType" bson:"variantType"`
	Position    Position  `json:"position" bson:"position"`
	Objective   Objective `json:"objective" bson:"objective"`
}

type Position struct {
	Dimensions     Dimensions                  `json:"dimensions" bson:"dimensions"`
	FEN            string                      `json:"fen,omitempty" bson:"fen,omitempty"`
	PieceProps     map[string]PieceProps       `json:"pieceProps,omitempty" bson:"pieceProps,omitempty"`
	PieceLocations map[string]map[string][]int `json:"pieceLocations,omitempty" bson:"pieceLocations,omitempty"` // Maps color -> piece type -> piece locations
}

type Dimensions struct {
	Ranks int `json:"ranks" bson:"ranks"`
	Files int `json:"files" bson:"files"`
}

type PieceProps struct {
	SlideOffsets []Coordinate `json:"slideOffsets,omitempty" bson:"slideOffsets,omitempty"`
	JumpProps    []JumpProps  `json:"jumpProps,omitempty" bson:"jumpProps,omitempty"`
	PromoProps   *PromoProps  `json:"promoProps,omitempty" bson:"promoProps,omitempty"`
}

type Coordinate struct {
	X int `json:"x" bson:"x"`
	Y int `json:"y" bson:"y"`
}

type JumpProps struct {
	Offset           Coordinate `json:"offset" bson:"offset"`
	IsCaptureAllowed bool       `json:"isCaptureAllowed" bson:"isCaptureAllowed"`
}

type PromoProps struct {
	PromotionSquares []int             `json:"promotionSquares" bson:"promotionSquares"`
	CanPromoteTo     []PromotionOption `json:"canPromoteTo" bson:"canPromoteTo"`
}

type PromotionOption struct {
	PieceType string `json:"pieceType" bson:"pieceType"`
	Name      string `json:"name" bson:"name"`
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
	ID         primitive.ObjectID            `json:"_id,omitempty" bson:"_id,omitempty"`
	Players    map[string]primitive.ObjectID `json:"players" bson:"players"`
	TemplateID *primitive.ObjectID           `json:"templateId,omitempty" bson:"template_id,omitempty"`
	Template   interface{}                   `json:"template,omitempty" bson:"template,omitempty"`
	Moves      []string                      `json:"moves" bson:"moves"`
	Result     GameResult                    `json:"result" bson:"result"`
	CreatedAt  primitive.DateTime            `json:"createdAt" bson:"created_at"`
	EndedAt    *primitive.DateTime           `json:"endedAt,omitempty" bson:"ended_at,omitempty"`
}

type ActiveGame struct {
	ID      string                        `json:"id"`
	Players map[string]primitive.ObjectID `json:"players"`
	Config  GameConfig                    `json:"gameConfig"`
	Moves   []string                      `json:"moves"`
	State   string                        `json:"state"`
	Turn    string                        `json:"turn"`
}

type PaginatedResponse[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}
