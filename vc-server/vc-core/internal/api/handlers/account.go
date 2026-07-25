package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"vc-server/vc-core/internal/db"
	"vc-server/vc-core/internal/middleware"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteAccountRequest struct {
	// Confirm must match the account username exactly.
	Confirm string `json:"confirm"`
}

// HandleDeleteAccount permanently deletes the authenticated user and related data.
func HandleDeleteAccount(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r)
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid user id", http.StatusUnauthorized)
			return
		}

		var req DeleteAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		req.Confirm = strings.TrimSpace(req.Confirm)
		if req.Confirm == "" {
			http.Error(w, "Type your username to confirm deletion", http.StatusBadRequest)
			return
		}

		users := database.Collection("users")
		var user bson.M
		err = users.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		username, _ := user["username"].(string)
		if username != req.Confirm {
			http.Error(w, "Confirmation does not match username", http.StatusBadRequest)
			return
		}

		ctx := context.TODO()

		if _, err := users.DeleteOne(ctx, bson.M{"_id": objectID}); err != nil {
			http.Error(w, "Failed to delete account", http.StatusInternalServerError)
			return
		}

		// Best-effort cleanup of user-owned data.
		_, _ = database.Collection("settings").DeleteMany(ctx, bson.M{"user_id": userID})
		_, _ = database.Collection("templates").DeleteMany(ctx, bson.M{"userId": objectID})

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(bson.M{"message": "Account deleted"})
	}
}
