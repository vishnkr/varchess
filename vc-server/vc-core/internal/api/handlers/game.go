package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"vc-server/chess"
	"vc-server/vc-core/internal/db"
	"vc-server/vc-core/internal/middleware"
	"vc-server/vc-core/internal/models"
	"vc-server/vc-core/internal/utils"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleGetGames(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, pageSize := utils.GetPaginationParams(r)
		opts := utils.GetPaginationOptions(page, pageSize)
		opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
		userID := middleware.GetUserIDFromContext(r)
		oid, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		}

		// Games where this user played white or black.
		filter := bson.M{
			"$or": []bson.M{
				{"players.w": oid},
				{"players.b": oid},
			},
		}

		collection := database.Collection("games")
		cursor, err := collection.Find(r.Context(), filter, opts)
		if err != nil {
			http.Error(w, "Failed to fetch games", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		var games []models.Game
		if err := cursor.All(r.Context(), &games); err != nil {
			http.Error(w, "Failed to decode games", http.StatusInternalServerError)
			return
		}
		if games == nil {
			games = []models.Game{}
		}

		totalCount, err := collection.CountDocuments(r.Context(), filter)
		if err != nil {
			http.Error(w, "Failed to count games", http.StatusInternalServerError)
			return
		}

		totalPages := utils.CalculateTotalPages(totalCount, int64(pageSize))

		response := models.PaginatedResponse[models.Game]{
			Items:      games,
			TotalCount: totalCount,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode games", http.StatusInternalServerError)
		}
	}
}


func HandleGetGame(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := chi.URLParam(r, "id")
		userID := middleware.GetUserIDFromContext(r)
		userOID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		}

		objID, err := primitive.ObjectIDFromHex(gameID)
		if err != nil {
			http.Error(w, "Invalid game ID", http.StatusBadRequest)
			return
		}

		collection := database.Collection("games")

		var game models.Game
		err = collection.FindOne(context.TODO(), bson.M{
			"_id": objID,
			"$or": []bson.M{
				{"players.w": userOID},
				{"players.b": userOID},
			},
		}).Decode(&game)
		if err != nil {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	}
}

func HandleCreateGame(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			GameConfig *chess.GameConfig `json:"gc,omitempty"`
			TemplateID *string `json:"templateId,omitempty"`
		}
	
		var response struct{
			GameID string `json:"gameId"`
		}

		userID := middleware.GetUserIDFromContext(r)
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		var gameConfig chess.GameConfig
		if request.GameConfig != nil {
			gameConfig = *request.GameConfig
		} else if request.TemplateID != nil {
			collection := database.Collection("templates")
			objectID, err := primitive.ObjectIDFromHex(*request.TemplateID)
			if err != nil {
				http.Error(w, "Invalid template ID", http.StatusBadRequest)
				return
			}
			var doc templateDoc
			err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&doc)
			if err != nil {
				http.Error(w, "Template not found", http.StatusNotFound)
				return
			}
			flat := normalizeTemplate(doc)
			gameConfig = chess.GameConfig{
				VariantType:    flat.VariantType,
				Name:           flat.Name,
				Dimensions:     flat.Dimensions,
				PieceProps:     flat.PieceProps,
				PieceLocations: flat.Position.PieceLocations,
				FEN:            flat.FEN,
				CustomData:     flat.CustomData,
			}
		} else {
			http.Error(w, "Either templateId or gameConfig must be provided", http.StatusBadRequest)
			return
		}
		if errMsg := validateGameConfig(gameConfig); errMsg != "" {
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
		fmt.Println(userID, "usr")
		shortID := utils.GenerateShortID(primitive.NewObjectID())
		game := models.ActiveGame{
			ID:      shortID,
			Config:  gameConfig,
			State:   models.Waiting,
			Players: make(map[string]models.Player),
			Moves:   make([]string, 0),
		}
		gameJSON, _ := json.Marshal(game)
		if _, redisErr := database.RedisClient.SetEx(context.TODO(), "game:"+shortID, gameJSON, 30*time.Minute).Result(); redisErr != nil {
			fmt.Println("err creating key:", redisErr)
		}
		response.GameID = shortID
		json.NewEncoder(w).Encode(response)
	}
}



func HandleUpdateGame(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := chi.URLParam(r, "id")
		userID := middleware.GetUserIDFromContext(r)

		objID, err := primitive.ObjectIDFromHex(gameID)
		if err != nil {
			http.Error(w, "Invalid game ID", http.StatusBadRequest)
			return
		}

		var updates bson.M
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		collection := database.Collection("games")

		result, err := collection.UpdateOne(
			context.TODO(),
			bson.M{"_id": objID, "user_id": userID},
			bson.M{"$set": updates},
		)
		if err != nil || result.MatchedCount == 0 {
			http.Error(w, "Failed to update game", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"message": "Game updated successfully"})
	}
}

func HandleDeleteGame(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := chi.URLParam(r, "id")
		userID := middleware.GetUserIDFromContext(r)

		objID, err := primitive.ObjectIDFromHex(gameID)
		if err != nil {
			http.Error(w, "Invalid game ID", http.StatusBadRequest)
			return
		}

		collection := database.Collection("games")

		result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID, "user_id": userID})
		if err != nil || result.DeletedCount == 0 {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"message": "Game deleted successfully"})
	}
}

func validateGameConfig(cfg chess.GameConfig) string {
	if cfg.FEN == "" {
		return "Game config fen is required"
	}
	vt := cfg.VariantType
	if vt != "wormhole" && vt != "teleport" {
		return ""
	}
	if cfg.CustomData == nil {
		return "Wormhole games require at least one wormhole pair"
	}
	raw, ok := cfg.CustomData["wormholePairs"]
	if !ok || raw == nil {
		return "Wormhole games require at least one wormhole pair"
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return "Wormhole games require at least one wormhole pair"
	}
	var pairs [][]int
	if err := json.Unmarshal(b, &pairs); err != nil || len(pairs) == 0 {
		return "Wormhole games require at least one wormhole pair"
	}
	for _, p := range pairs {
		if len(p) < 2 || p[0] == p[1] {
			return "Each wormhole pair needs two different squares"
		}
	}
	return ""
}