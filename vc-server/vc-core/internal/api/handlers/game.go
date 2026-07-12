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

		collection := database.Collection("games")
		cursor, err := collection.Find(r.Context(), bson.M{}, opts)
		if err != nil {
			http.Error(w, "Failed to fetch templates", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		var games []models.Game
		if err := cursor.All(r.Context(), &games); err != nil {
			http.Error(w, "Failed to decode templates", http.StatusInternalServerError)
			return
		}

		totalCount, err := collection.CountDocuments(r.Context(), bson.M{})
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

		objID, err := primitive.ObjectIDFromHex(gameID)
		if err != nil {
			http.Error(w, "Invalid game ID", http.StatusBadRequest)
			return
		}

		collection := database.Collection("games")

		var game bson.M
		err = collection.FindOne(context.TODO(), bson.M{"_id": objID, "user_id": userID}).Decode(&game)
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
		if request.TemplateID != nil {
			collection := database.Collection("templates")
			objectID, err := primitive.ObjectIDFromHex(*request.TemplateID)
			var template models.Template
			err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&template)
			if err != nil {
				http.Error(w, "Template not found", http.StatusNotFound)
				return
			}

			gameConfig = chess.GameConfig{
				VariantType: template.VariantType,
				Name:    template.Name,
				PieceProps: template.Position.PieceProps,
				PieceLocations: template.Position.PieceLocations,
				FEN: template.Position.FEN,
				CustomData:   template.CustomData,
			}
		} else if request.GameConfig != nil {
			gameConfig = *request.GameConfig
		} else {
			http.Error(w, "Either templateId or gameConfig must be provided", http.StatusBadRequest)
			return
		}
		fmt.Println(userID,"usr")
		shortID := utils.GenerateShortID(primitive.NewObjectID())
		game := models.ActiveGame{
			ID:    	shortID,
			Config: gameConfig,
			State:  models.Waiting,
			Players: make(map[string]models.Player),
			Moves: make([]string, 0),
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