package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"vc-server/vc-core/internal/db"
	"vc-server/vc-core/internal/middleware"
	"vc-server/vc-core/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	// Identifier is username or email.
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleSignup(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		req.Username = strings.TrimSpace(req.Username)
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		if req.Username == "" || req.Email == "" || len(req.Password) < 8 {
			http.Error(w, "Username, email, and password (min 8 chars) are required", http.StatusBadRequest)
			return
		}

		collection := database.Collection("users")

		var existingUser bson.M
		err := collection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&existingUser)
		if err == nil {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}
		var existingUsername bson.M
		err = collection.FindOne(context.TODO(), bson.M{"username": req.Username}).Decode(&existingUsername)
		if err == nil {
			http.Error(w, "Username already in use", http.StatusConflict)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		user := bson.M{
			"username":      req.Username,
			"email":         req.Email,
			"password":      string(hashedPassword),
			"auth_provider": "password",
			"created_at":    time.Now(),
		}

		result, err := collection.InsertOne(context.TODO(), user)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		userID := result.InsertedID.(primitive.ObjectID).Hex()
		accessToken, err := utils.GenerateToken(userID)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}
		refreshToken, err := utils.GenerateRefreshToken(userID)
		if err != nil {
			http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
			return
		}

		setAuthCookies(w, accessToken, refreshToken)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(bson.M{
			"uid":          userID,
			"username":     req.Username,
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}

func HandleLogin(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		identifier := strings.TrimSpace(req.Username)
		if identifier == "" {
			identifier = strings.TrimSpace(req.Email)
		}
		if identifier == "" || req.Password == "" {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		collection := database.Collection("users")
		var user bson.M

		filter := bson.M{"username": identifier}
		if strings.Contains(identifier, "@") {
			filter = bson.M{"email": strings.ToLower(identifier)}
		}

		err := collection.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		pw, _ := user["password"].(string)
		if pw == "" {
			http.Error(w, "This account uses social sign-in. Continue with Google or GitHub.", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(pw), []byte(req.Password))
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		userID := user["_id"].(primitive.ObjectID).Hex()
		username, _ := user["username"].(string)
		token, err := utils.GenerateToken(userID)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := utils.GenerateRefreshToken(userID)
		if err != nil {
			http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
			return
		}
		setAuthCookies(w, token, refreshToken)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(bson.M{
			"uid":          userID,
			"username":     username,
			"accessToken":  token,
			"refreshToken": refreshToken,
		})
	}
}

func HandleValidateAuth(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r)
		objectId, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid user id", http.StatusUnauthorized)
			return
		}
		var user bson.M
		err = database.Collection("users").FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		response := map[string]string{
			"userId":   userID,
			"username": user["username"].(string),
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func HandleRefresh(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
			// Fall back to cookie if body omitted.
			if cookie, err := r.Cookie("refreshToken"); err == nil {
				req.RefreshToken = cookie.Value
			}
		}
		if req.RefreshToken == "" {
			http.Error(w, "Missing refresh token", http.StatusBadRequest)
			return
		}

		userID, err := utils.ParseRefreshToken(req.RefreshToken)
		if err != nil {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		objectId, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid user id", http.StatusUnauthorized)
			return
		}
		var user bson.M
		err = database.Collection("users").FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&user)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		accessToken, err := utils.GenerateToken(userID)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(bson.M{
			"accessToken": accessToken,
			"uid":         userID,
			"username":    user["username"],
		})
	}
}
