package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"vc-server/vc-core/internal/db"
	"vc-server/vc-core/internal/middleware"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Settings struct {
	UserID            string `bson:"user_id" json:"userId"`
	BoardTheme        string `bson:"board_theme" json:"boardTheme"`
	ShowPossibleMoves bool   `bson:"show_moves" json:"showPossibleMoves"`
	EnablePremove     bool   `bson:"enable_premove" json:"enablePremove"`
	HighlightLastMove bool   `bson:"highlight_last_move" json:"highlightLastMove"`
	ShowCoordinates   bool   `bson:"show_coordinates" json:"showCoordinates"`
	ConfirmResign     bool   `bson:"confirm_resign" json:"confirmResign"`
	HideChat          bool   `bson:"hide_chat" json:"hideChat"`
	SoundEnabled      bool   `bson:"sound_enabled" json:"soundEnabled"`
	AppTheme          string `bson:"app_theme" json:"appTheme"`
}

func DefaultSettings(userID string) Settings {
	return Settings{
		UserID:            userID,
		BoardTheme:        "Default",
		ShowPossibleMoves: true,
		EnablePremove:     false,
		HighlightLastMove: true,
		ShowCoordinates:   false,
		ConfirmResign:     true,
		HideChat:          false,
		SoundEnabled:      false,
		AppTheme:          "dark",
	}
}

// HandleGetSettings fetches user-specific settings (defaults if none stored).
func HandleGetSettings(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r)
		collection := database.Collection("settings")

		var settings Settings
		err := collection.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(&settings)
		if err == mongo.ErrNoDocuments {
			settings = DefaultSettings(userID)
		} else if err != nil {
			http.Error(w, "Failed to fetch settings", http.StatusInternalServerError)
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

// HandleUpdateSettings upserts user settings.
func HandleUpdateSettings(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r)

		var body Settings
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		body.UserID = userID
		collection := database.Collection("settings")

		opts := options.Update().SetUpsert(true)
		_, err := collection.UpdateOne(
			context.TODO(),
			bson.M{"user_id": userID},
			bson.M{"$set": body},
			opts,
		)
		if err != nil {
			http.Error(w, "Failed to update settings", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
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
