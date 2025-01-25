package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"vc-core/internal/db"
	"vc-core/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the request body for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HandleSignup registers a new user
func HandleSignup(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		collection := database.Collection("users")

		var existingUser bson.M
		err := collection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&existingUser)
		if err == nil {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// Create user
		user := bson.M{
			"username": req.Username,
			"email":    req.Email,
			"password": string(hashedPassword),
		}

		result, err := collection.InsertOne(context.TODO(), user)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(bson.M{"user_id": result.InsertedID})
	}
}

// HandleLogin authenticates a user and returns a JWT token
func HandleLogin(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		collection := database.Collection("users")

		// Find user by email
		var user bson.M
		err := collection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user["password"].(string)), []byte(req.Password))
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		userID := user["_id"].(primitive.ObjectID).Hex()
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bson.M{
			"access_token": token,
			"refresh_token": refreshToken,
		})
	}
}
