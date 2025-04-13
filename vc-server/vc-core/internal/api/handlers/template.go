package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"vc-server/vc-core/internal/db"
	"vc-server/vc-core/internal/middleware"
	"vc-server/vc-core/internal/models"
	"vc-server/vc-core/internal/utils"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleGetTemplates(database *db.DB) http.HandlerFunc {
	type responseItem struct {
		Name        string                       `json:"name" bson:"name"`
		ID          primitive.ObjectID           `json:"_id,omitempty" bson:"_id,omitempty"`
		UserId      primitive.ObjectID           `json:"userId,omitempty" bson:"userId,omitempty"`
		VariantType string                       `json:"variantType" bson:"variantType"`
		Dimensions  models.Dimensions            `json:"dimensions" bson:"dimensions"`
		FEN         string                       `json:"fen,omitempty" bson:"fen,omitempty"`
		PieceProps  map[string]models.PieceProps `json:"pieceProps,omitempty" bson:"pieceProps,omitempty"`
		CustomData  map[string]interface{}       `json:"customData" bson:"customData"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		page, pageSize := utils.GetPaginationParams(r)
		opts := utils.GetPaginationOptions(page, pageSize)

		collection := database.Collection("templates")
		cursor, err := collection.Find(r.Context(), bson.M{}, opts)
		if err != nil {
			http.Error(w, "Failed to fetch templates", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		var templates []responseItem
		if err := cursor.All(r.Context(), &templates); err != nil {
			http.Error(w, "Failed to decode templates", http.StatusInternalServerError)
			return
		}

		totalCount, err := collection.CountDocuments(r.Context(), bson.M{})
		if err != nil {
			http.Error(w, "Failed to count templates", http.StatusInternalServerError)
			return
		}

		totalPages := utils.CalculateTotalPages(totalCount, int64(pageSize))

		response := models.PaginatedResponse[responseItem]{
			Items:      templates,
			TotalCount: totalCount,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode templates", http.StatusInternalServerError)
		}
	}
}

func HandleGetTemplate(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templateID := chi.URLParam(r, "id")
		userID, err := middleware.GetUserObjIDFromContext(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		objID, err := primitive.ObjectIDFromHex(templateID)
		if err != nil {
			http.Error(w, "Invalid template ID", http.StatusBadRequest)
			return
		}

		collection := database.Collection("templates")

		var template bson.M
		err = collection.FindOne(context.TODO(), bson.M{"_id": objID, "userId": userID}).Decode(&template)
		if err != nil {
			http.Error(w, "Template not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(template)
	}
}

func HandleCreateTemplate(database *db.DB) http.HandlerFunc {
	type request struct {
		Name        string                       `json:"name" bson:"name"`
		ID          primitive.ObjectID           `json:"_id,omitempty" bson:"_id,omitempty"`
		UserId      primitive.ObjectID           `json:"userId,omitempty" bson:"userId,omitempty"`
		VariantType string                       `json:"variantType" bson:"variantType"`
		Dimensions  models.Dimensions            `json:"dimensions" bson:"dimensions"`
		FEN         string                       `json:"fen,omitempty" bson:"fen,omitempty"`
		PieceProps  map[string]models.PieceProps `json:"pieceProps,omitempty" bson:"pieceProps,omitempty"`
		CustomData  map[string]interface{}       `json:"customData" bson:"customData"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r)

		var template request
		if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
			http.Error(w, "Invalid template structure", http.StatusBadRequest)
			return
		}

		userObjectId, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid user id", http.StatusBadRequest)
			return
		}
		template.UserId = userObjectId
		template.ID = primitive.NewObjectID()

		collection := database.Collection("templates")

		_, err = collection.InsertOne(context.TODO(), template)
		if err != nil {
			http.Error(w, "Failed to create template", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(template)
	}
}

func HandleUpdateTemplate(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templateID := chi.URLParam(r, "id")
		userID, err := middleware.GetUserObjIDFromContext(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		objID, err := primitive.ObjectIDFromHex(templateID)
		if err != nil {
			http.Error(w, "Invalid template ID", http.StatusBadRequest)
			return
		}

		var updates bson.M
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		collection := database.Collection("templates")

		result, err := collection.UpdateOne(
			context.TODO(),
			bson.M{"_id": objID, "userId": userID},
			bson.M{"$set": updates},
		)
		if err != nil || result.MatchedCount == 0 {
			http.Error(w, "Failed to update template", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

// HandleDeleteTemplate deletes a template
func HandleDeleteTemplate(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templateID := chi.URLParam(r, "id")
		userID, err := middleware.GetUserObjIDFromContext(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		objID, err := primitive.ObjectIDFromHex(templateID)
		if err != nil {
			http.Error(w, "Invalid template ID", http.StatusBadRequest)
			return
		}

		collection := database.Collection("templates")

		result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID, "userId": userID})
		if err != nil || result.DeletedCount == 0 {
			http.Error(w, "Template not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"message": "Template deleted successfully"})
	}
}
