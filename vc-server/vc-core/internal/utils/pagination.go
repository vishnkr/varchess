package utils

import (
	"context"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPaginationParams(r *http.Request) (int, int) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = ps
		}
	}

	return page, pageSize
}

func GetPaginationOptions(page, pageSize int) *options.FindOptions {
	skip := (page - 1) * pageSize
	opts := options.Find()
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(pageSize))
	return opts
}

func GetTotalCount(collection *mongo.Collection, ctx context.Context) (int64, error) {
	return collection.CountDocuments(ctx, bson.M{})
}

func CalculateTotalPages(totalCount, pageSize int64) int {
	totalPages := int(totalCount / pageSize)
	if totalCount%pageSize > 0 {
		totalPages++
	}
	return totalPages
}

func FetchPaginatedData(collection *mongo.Collection, ctx context.Context, opts *options.FindOptions, result interface{}) error {
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)
}
