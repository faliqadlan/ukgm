package utils

import (
	"fmt"

	"github.com/iman231/go-help"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ObjectIDHex(ID interface{}) string {
	var objectID string
	switch v := ID.(type) {
	case string:
		objectID = v
	default:
		objectID = v.(primitive.ObjectID).Hex()
	}
	return objectID
}

func StringToObjectID(ID interface{}) (result interface{}, err error) {
	switch v := ID.(type) {
	case string:
		id, errTemp := primitive.ObjectIDFromHex(v)
		if errTemp != nil {
			err = errTemp
			return
		} else {
			result = id
			return
		}
	default:
		err = fmt.Errorf("Invalid type")
		return
	}
}

func NotDeleted(filter bson.M) bson.M {
	filter["timestamp.deleted_at"] = nil
	return filter
}

func ValidateInArray(needed bool, array []string, name, input string) error {
	preWarning := fmt.Sprintf("Field %s ", name)

	if needed && input == "" {
		return fmt.Errorf(preWarning + "is needed")
	}

	if !help.InArray(input, array) {
		return fmt.Errorf(preWarning + "is invalid")
	}

	return nil
}
