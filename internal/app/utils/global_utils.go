package utils

import (
	"fmt"
	"strconv"

	"github.com/iman231/go-help"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
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
			zap.S().Error(err)
			return
		} else {
			result = id
			return
		}
	default:
		result = ID
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

// ToFixedCustom :
func ToFixedCustom(comma string, v float64) float64 {
	// FLOAT TO STRING FORMATING DIGIT DECIMAL
	res := fmt.Sprintf("%."+comma+"f", v)
	// STRING TO FLOAT64
	v, _ = strconv.ParseFloat(res, 64)

	return v
}

func RateDivider(index, totalAmmount int) (result float64) {
	onePerRate := 100 / float64(totalAmmount)
	onePerRate = ToFixedCustom("4", onePerRate)
	if index >= totalAmmount {
		restOfRate := 0.0
		for i := 0; i < (totalAmmount - 1); i++ {
			restOfRate += onePerRate
			restOfRate = ToFixedCustom("4", restOfRate)
		}
		result = 100 - restOfRate
	} else {
		result = onePerRate
	}
	return
}
