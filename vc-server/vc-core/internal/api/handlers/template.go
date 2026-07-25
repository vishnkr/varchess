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

// templateDoc supports both nested `position` (canonical) and legacy flat fields.
type templateDoc struct {
	Name        string                       `json:"name" bson:"name"`
	ID          primitive.ObjectID           `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId      primitive.ObjectID           `json:"userId,omitempty" bson:"userId,omitempty"`
	VariantType string                       `json:"variantType" bson:"variantType"`
	Dimensions  chess.Dimensions             `json:"dimensions,omitempty" bson:"dimensions,omitempty"`
	FEN         string                       `json:"fen,omitempty" bson:"fen,omitempty"`
	PieceProps  map[string]chess.PieceProps  `json:"pieceProps,omitempty" bson:"pieceProps,omitempty"`
	Position    models.Position              `json:"position,omitempty" bson:"position,omitempty"`
	CustomData  map[string]interface{}       `json:"customData" bson:"customData"`
}

// flatTemplateResponse is the client-facing shape (flat + nested for compatibility).
type flatTemplateResponse struct {
	Name        string                      `json:"name"`
	ID          primitive.ObjectID          `json:"_id,omitempty"`
	UserId      primitive.ObjectID          `json:"userId,omitempty"`
	VariantType string                      `json:"variantType"`
	Dimensions  chess.Dimensions            `json:"dimensions"`
	FEN         string                      `json:"fen,omitempty"`
	PieceProps  map[string]chess.PieceProps `json:"pieceProps,omitempty"`
	Position    models.Position             `json:"position"`
	CustomData  map[string]interface{}      `json:"customData"`
}

func normalizeTemplate(doc templateDoc) flatTemplateResponse {
	pos := doc.Position
	if pos.FEN == "" && doc.FEN != "" {
		pos.FEN = doc.FEN
	}
	if pos.Dimensions.Ranks == 0 && doc.Dimensions.Ranks != 0 {
		pos.Dimensions = doc.Dimensions
	}
	if pos.PieceProps == nil && doc.PieceProps != nil {
		pos.PieceProps = doc.PieceProps
	}
	return flatTemplateResponse{
		Name:        doc.Name,
		ID:          doc.ID,
		UserId:      doc.UserId,
		VariantType: doc.VariantType,
		Dimensions:  pos.Dimensions,
		FEN:         pos.FEN,
		PieceProps:  pos.PieceProps,
		Position:    pos,
		CustomData:  doc.CustomData,
	}
}

func HandleGetTemplates(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context()).With().Str("handler", "HandleGetTemplates").Logger()
		page, pageSize := utils.GetPaginationParams(r)
		opts := utils.GetPaginationOptions(page, pageSize)
		userID, err := middleware.GetUserObjIDFromContext(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		collection := database.Collection("templates")
		filter := bson.M{"userId": userID}
		cursor, err := collection.Find(r.Context(), filter, opts)
		if err != nil {
			log.Error().Err(err).Msg("Failed to query templates")
			http.Error(w, "Failed to fetch templates", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		var docs []templateDoc
		if err := cursor.All(r.Context(), &docs); err != nil {
			http.Error(w, "Failed to decode templates", http.StatusInternalServerError)
			return
		}

		items := make([]flatTemplateResponse, 0, len(docs))
		for _, doc := range docs {
			items = append(items, normalizeTemplate(doc))
		}

		totalCount, err := collection.CountDocuments(r.Context(), filter)
		if err != nil {
			http.Error(w, "Failed to count templates", http.StatusInternalServerError)
			return
		}

		totalPages := utils.CalculateTotalPages(totalCount, int64(pageSize))

		log.Info().Int64("totalCount", totalCount).Int("totalPages", totalPages).Msg("Successfully fetched templates")
		response := models.PaginatedResponse[flatTemplateResponse]{
			Items:      items,
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

		var doc templateDoc
		err = collection.FindOne(context.TODO(), bson.M{"_id": objID, "userId": userID}).Decode(&doc)
		if err != nil {
			log.Warn().Str("templateID", templateID).Msg("Template not found")
			http.Error(w, "Template not found", http.StatusNotFound)
			return
		}
		log.Info().Str("templateID", templateID).Msg("Successfully fetched template")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(normalizeTemplate(doc))
	}
}

func HandleCreateTemplate(database *db.DB) http.HandlerFunc {
	type request struct {
		Name        string                      `json:"name"`
		VariantType string                      `json:"variantType"`
		Dimensions  chess.Dimensions            `json:"dimensions"`
		FEN         string                      `json:"fen,omitempty"`
		PieceProps  map[string]chess.PieceProps `json:"pieceProps,omitempty"`
		CustomData  map[string]interface{}      `json:"customData"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r)
		log := logger.FromContext(r.Context()).With().Str("handler", "HandleCreateTemplate").Logger()

		var body request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid template structure", http.StatusBadRequest)
			return
		}
		log.Info().Str("userID", userID).Msg("Creating new template")
		userObjectId, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Invalid user id", http.StatusBadRequest)
			return
		}

		template := models.Template{
			Name:        body.Name,
			ID:          primitive.NewObjectID(),
			UserId:      userObjectId,
			VariantType: body.VariantType,
			Position: models.Position{
				Dimensions: body.Dimensions,
				FEN:        body.FEN,
				PieceProps: body.PieceProps,
			},
			CustomData: body.CustomData,
		}

		collection := database.Collection("templates")
		_, err = collection.InsertOne(context.TODO(), template)
		if err != nil {
			log.Error().Err(err).Msg("Failed to insert template")
			http.Error(w, "Failed to create template", http.StatusInternalServerError)
			return
		}
		log.Info().Str("templateID", template.ID.Hex()).Msg("Template created successfully")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(normalizeTemplate(templateDoc{
			Name:        template.Name,
			ID:          template.ID,
			UserId:      template.UserId,
			VariantType: template.VariantType,
			Position:    template.Position,
			CustomData:  template.CustomData,
		}))
	}
}

func HandleUpdateTemplate(database *db.DB) http.HandlerFunc {
	type request struct {
		Name        string                      `json:"name"`
		VariantType string                      `json:"variantType"`
		Dimensions  chess.Dimensions            `json:"dimensions"`
		FEN         string                      `json:"fen,omitempty"`
		PieceProps  map[string]chess.PieceProps `json:"pieceProps,omitempty"`
		CustomData  map[string]interface{}      `json:"customData"`
	}

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

		var body request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		collection := database.Collection("templates")
		update := bson.M{
			"$set": bson.M{
				"name":                   body.Name,
				"variantType":            body.VariantType,
				"customData":             body.CustomData,
				"position.dimensions":    body.Dimensions,
				"position.fen":           body.FEN,
				"position.pieceProps":    body.PieceProps,
			},
			"$unset": bson.M{
				// Clear legacy flat fields if present.
				"dimensions":  "",
				"fen":         "",
				"pieceProps":  "",
			},
		}

		result, err := collection.UpdateOne(
			context.TODO(),
			bson.M{"_id": objID, "userId": userID},
			update,
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
