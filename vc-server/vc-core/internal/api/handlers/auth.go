package handlers

import (
	"context"
	"encoding/json"
	"net/http"
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
	Username    string `json:"username"`
	Password string `json:"password"`
}

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


func HandleLogin(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		collection := database.Collection("users")

		var user bson.M
		err := collection.FindOne(context.TODO(), bson.M{"username": req.Username}).Decode(&user)
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
		http.SetCookie(w, &http.Cookie{
			Name:     "accessToken",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24),
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bson.M{
			"uid": userID,
			"username": req.Username,
			"accessToken":  token,
			"refreshToken": refreshToken,
		})
	}
}


func HandleValidateAuth(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r)
		objectId,err := primitive.ObjectIDFromHex(userID)
		if err!=nil{
			http.Error(w, "Invalid user id", http.StatusUnauthorized)
			return
		}
		var user bson.M
		err = database.Collection("users").FindOne(context.TODO(),bson.M{"_id":objectId}).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		response := map[string]string{
			"status": "ok",
			"user_id": userID,
			"username": user["username"].(string),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
}
}