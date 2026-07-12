package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"vc-server/chess"
	"vc-server/vc-core/internal/db"
	"vc-server/vc-core/internal/logger"
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
		Dimensions  chess.Dimensions            `json:"dimensions" bson:"dimensions"`
		FEN         string                       `json:"fen,omitempty" bson:"fen,omitempty"`
		PieceProps  map[string]chess.PieceProps `json:"pieceProps,omitempty" bson:"pieceProps,omitempty"`
		CustomData  map[string]interface{}       `json:"customData" bson:"customData"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context()).With().Str("handler", "HandleGetTemplates").Logger()
		page, pageSize := utils.GetPaginationParams(r)
		opts := utils.GetPaginationOptions(page, pageSize)
		userID, err := middleware.GetUserObjIDFromContext(r)
		collection := database.Collection("templates")
		filter := bson.M{ "userId": userID}
		cursor, err := collection.Find(r.Context(), filter, opts)
		if err != nil {
			log.Error().Err(err).Msg("Failed to query templates")
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

		log.Info().Int64("totalCount", totalCount).Int("totalPages", totalPages).Msg("Successfully fetched templates")
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
		log := logger.FromContext(r.Context()).With().Str("handler", "HandleGetTemplate").Logger()
		userID, err := middleware.GetUserObjIDFromContext(r)
		log.Info().Str("templateID", templateID).Msg("Fetching template")
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
			log.Warn().Str("templateID", templateID).Msg("Template not found")
			http.Error(w, "Template not found", http.StatusNotFound)
			return
		}
		log.Info().Str("templateID", templateID).Msg("Successfully fetched template")
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
		Dimensions  chess.Dimensions            `json:"dimensions" bson:"dimensions"`
		FEN         string                       `json:"fen,omitempty" bson:"fen,omitempty"`
		PieceProps  map[string]chess.PieceProps `json:"pieceProps,omitempty" bson:"pieceProps,omitempty"`
		CustomData  map[string]interface{}       `json:"customData" bson:"customData"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r)
		log := logger.FromContext(r.Context()).With().Str("handler", "HandleCreateTemplate").Logger()

		var template request
		if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
			http.Error(w, "Invalid template structure", http.StatusBadRequest)
			return
		}
		log.Info().Str("userID", userID).Msg("Creating new template")
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
			log.Error().Err(err).Msg("Failed to insert template")
			http.Error(w, "Failed to create template", http.StatusInternalServerError)
			return
		}
		log.Info().Str("templateID", template.ID.Hex()).Msg("Template created successfully")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(template)
	}
}

func HandleUpdateTemplate(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context()).With().Str("handler", "HandleUpdateTemplate").Logger()
		templateID := chi.URLParam(r, "id")
		log.Info().Str("template_id", templateID).Msg("Update template request received")
		userID, err := middleware.GetUserObjIDFromContext(r)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get user from context")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		objID, err := primitive.ObjectIDFromHex(templateID)
		if err != nil {
			log.Warn().Str("template_id", templateID).Msg("Invalid template ID")
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
		if err != nil {
			log.Error().Err(err).Msg("Failed to update template")
			http.Error(w, "Failed to update template", http.StatusInternalServerError)
			return
		}

		if result.MatchedCount == 0 {
			log.Warn().Msg("No matching template found for update")
			http.Error(w, "Template not found", http.StatusNotFound)
			return
		}

		log.Info().Int64("matched", result.MatchedCount).Int64("modified", result.ModifiedCount).Msg("Template updated successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

// HandleDeleteTemplate deletes a template
func HandleDeleteTemplate(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context()).With().Str("handler", "HandleDeleteTemplate").Logger()
		templateID := chi.URLParam(r, "id")
		log.Info().Str("template_id", templateID).Msg("Delete template request received")
		userID, err := middleware.GetUserObjIDFromContext(r)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get user from context")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		objID, err := primitive.ObjectIDFromHex(templateID)
		if err != nil {
			log.Warn().Str("template_id", templateID).Msg("Invalid template ID")
			http.Error(w, "Invalid template ID", http.StatusBadRequest)
			return
		}

		collection := database.Collection("templates")

		result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID, "userId": userID})
		if err != nil {
			log.Error().Err(err).Msg("Error during deletion")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if result.DeletedCount == 0 {
			log.Warn().Msg("No template deleted — not found or unauthorized")
			http.Error(w, "Template not found", http.StatusNotFound)
			return
		}

		log.Info().Int64("deleted_count", result.DeletedCount).Msg("Template deleted successfully")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"message": "Template deleted successfully"})
	}
}
