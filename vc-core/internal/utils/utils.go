package utils

import (
	"math/big"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)



func GenerateShortID(objectID primitive.ObjectID) string {
	intVal := new(big.Int)
	intVal.SetBytes(objectID[:])
	shortID := strings.ToUpper(intVal.Text(36))
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}
	return shortID
}