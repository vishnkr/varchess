package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"vc-core/internal/db"
	"vc-core/internal/middleware"

	"go.mongodb.org/mongo-driver/bson"
)

type Settings struct {
	UserID       string `bson:"user_id" json:"user_id"`
	BoardTheme   string `bson:"board_theme" json:"board_theme"`
	ShowMoves    bool   `bson:"show_moves" json:"show_moves"`
	EnablePremove bool  `bson:"enable_premove" json:"enable_premove"`
}

// HandleGetSettings fetches user-specific settings
func HandleGetSettings(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r)

		collection := database.Collection("settings")

		var settings Settings
		err := collection.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(&settings)
		if err != nil {
			http.Error(w, "Settings not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}
}

// HandleCreateSettings initializes settings for a user
func HandleCreateSettings(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r)

		var settings Settings
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		settings.UserID = userID

		collection := database.Collection("settings")

		_, err := collection.InsertOne(context.TODO(), settings)
		if err != nil {
			http.Error(w, "Failed to create settings", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(settings)
	}
}

// HandleUpdateSettings updates user settings
func HandleUpdateSettings(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r)

		var updates bson.M
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		collection := database.Collection("settings")

		result, err := collection.UpdateOne(
			context.TODO(),
			bson.M{"user_id": userID},
			bson.M{"$set": updates},
		)
		if err != nil || result.MatchedCount == 0 {
			http.Error(w, "Failed to update settings", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"message": "Settings updated successfully"})
	}
}

// HandleDeleteSettings resets settings to default (delete user settings)
func HandleDeleteSettings(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r)

		collection := database.Collection("settings")

		result, err := collection.DeleteOne(context.TODO(), bson.M{"user_id": userID})
		if err != nil || result.DeletedCount == 0 {
			http.Error(w, "Settings not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"message": "Settings deleted successfully"})
	}
}
