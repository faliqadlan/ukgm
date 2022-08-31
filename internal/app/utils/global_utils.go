package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iman231/go-help"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GotTagName(object interface{}, variableName, tagName string) (result string) {
	field, ok := reflect.TypeOf(object).Elem().FieldByName(variableName)
	if !ok {
		return
	}
	fullTag := field.Tag.Get(tagName)
	tagArrays := strings.Split(fullTag, ",")
	result = tagArrays[0]
	if result == "" {
		result = strings.ToLower(variableName)
	}
	return
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

func GetMainOrDetail(object, objectPointer interface{}) (main, detail []string, err error) {
	var (
		interfaceVar interface{}
	)
	interfaceVar = object
	main = []string{}
	detail = []string{}

	val := reflect.TypeOf(interfaceVar)
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Field(i).Name
		partType := GotTagName(objectPointer, fieldName, "part")
		if partType == "main" {
			main = append(main, GotTagName(objectPointer, fieldName, "bson"))
		} else if partType == "detail" {
			detail = append(detail, GotTagName(objectPointer, fieldName, "bson"))
		}
	}

	return
}

func BuildProjectionsFindOneOptions(include, exclude []string) (result *options.FindOneOptions, err error) {
	result = &options.FindOneOptions{}
	if len(include) > 0 {
		filter := bson.M{}
		for _, field := range include {
			filter[field] = 1
		}
		result.SetProjection(filter)
	}

	if len(exclude) > 0 {
		filter := bson.M{}
		for _, field := range exclude {
			filter[field] = 0
		}
		result.SetProjection(filter)
	}
	return
}
