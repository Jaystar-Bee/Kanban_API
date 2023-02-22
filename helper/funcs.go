package helpers

import (

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToPrimitive(id string) (primitive.ObjectID, error) {
	columnID, err := primitive.ObjectIDFromHex(id)

	return columnID, err
}
