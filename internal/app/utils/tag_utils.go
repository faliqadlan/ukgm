package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iman231/go-help"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

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

func GotType(object interface{}, variableName string) (result string) {
	valObject := reflect.ValueOf(object)
	if valObject.Kind() != reflect.Ptr {
		valObject = reflect.New(reflect.TypeOf(object))
	}
	field := valObject.Elem().FieldByName(variableName)
	result = field.Kind().String()
	// fullTag := field.Tag.Get(tagName)
	// tagArrays := strings.Split(fullTag, ",")
	// result = tagArrays[0]
	// if result == "" {
	// 	result = strings.ToLower(variableName)
	// }
	return
}

func CheckNil(object interface{}, variableName string) (result bool) {
	valObject := reflect.ValueOf(object)
	if valObject.Kind() != reflect.Ptr {
		valObject = reflect.New(reflect.TypeOf(object))
		valObject.Elem().Set(reflect.ValueOf(object))
	}
	field := valObject.Elem().FieldByName(variableName)
	result = field.IsZero()
	// fullTag := field.Tag.Get(tagName)
	// tagArrays := strings.Split(fullTag, ",")
	// result = tagArrays[0]
	// if result == "" {
	// 	result = strings.ToLower(variableName)
	// }
	return
}

func GotProperties(object interface{}, variableName string) (result []string) {
	tagName := "properties"
	valObject := reflect.ValueOf(object)
	if valObject.Kind() != reflect.Ptr {
		valObject = reflect.New(reflect.TypeOf(object))
	}

	field, ok := valObject.Type().Elem().FieldByName(variableName)
	if !ok {
		return
	}
	fullTag := field.Tag.Get(tagName)
	result = strings.Split(fullTag, ",")
	return
}

func GotUniqueType(object interface{}, variableName string) (result string) {
	tagName := "unique"
	valObject := reflect.ValueOf(object)
	if valObject.Kind() != reflect.Ptr {
		valObject = reflect.New(reflect.TypeOf(object))
	}

	field, ok := valObject.Type().Elem().FieldByName(variableName)
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

func GetSearchableField(object interface{}) (result []string, err error) {
	valObject := reflect.ValueOf(object)
	if valObject.Kind() != reflect.Ptr {
		valObject = reflect.New(reflect.TypeOf(object))
	}
	result = []string{}

	val := valObject.Elem().Type()
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Field(i).Name
		properties := GotProperties(object, fieldName)
		if help.InArray("searchable", properties) {
			result = append(result, GotTagName(object, fieldName, "bson"))
		}
	}

	return
}

func DuplicateConstantVariable(object, destination interface{}) (err error) {
	var ()
	valObject := reflect.ValueOf(object)
	if valObject.Kind() != reflect.Ptr {
		err = fmt.Errorf("sent object need to be a pointer")
		zap.S().Error(err)
		return
	}
	valSource := reflect.ValueOf(destination)
	if valSource.Kind() != reflect.Ptr {
		err = fmt.Errorf("sent object need to be a pointer")
		zap.S().Error(err)
		return
	}

	if valSource.Elem().Kind() != valObject.Elem().Kind() {
		err = fmt.Errorf("object and source is not the same kind")
		zap.S().Error(err)
		return
	}

	val := valObject.Elem().Type()
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Field(i).Name
		properties := GotProperties(object, fieldName)
		if help.InArray("constant", properties) {
			constanvtValue := valObject.Elem().FieldByName(fieldName)
			valSource.Elem().FieldByName(fieldName).Set(constanvtValue)
		}
	}

	return
}

func ConstructUniqueFilter(object interface{}) (filter bson.M, err error) {
	var (
		primary, secondary string
		primarySuccess     bool
	)
	valObject := reflect.ValueOf(object)
	if valObject.Kind() != reflect.Ptr {
		valObject = reflect.New(reflect.TypeOf(object))
		valObject.Elem().Set(reflect.ValueOf(object))
	}
	filter = primitive.M{}

	val := valObject.Elem().Type()
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Field(i).Name
		partType := GotUniqueType(object, fieldName)
		if partType == "primary" {
			primary = fieldName
			// primaryIndex = i
		} else if partType == "secondary" {
			secondary = fieldName
			// secondaryIndex = i
		}
		if primary != "" && secondary != "" {
			break
		}
	}
	if primary != "" || secondary != "" {
		if GotType(object, primary) == "interface" && !CheckNil(object, primary) {
			field := valObject.Elem().FieldByName(primary)
			oid, err := StringToObjectID(field.Interface())
			if err == nil {
				filter[GotTagName(object, primary, "bson")] = oid
				primarySuccess = true
			}
		}

		if !primarySuccess {
			if !CheckNil(object, secondary) {
				field := valObject.Elem().FieldByName(secondary)
				filter[GotTagName(object, secondary, "bson")] = field.Interface()
				primarySuccess = true

			} else {
				err = fmt.Errorf("object doesn't have primary or secondary unique key")
				zap.S().Error(err)
			}
		}
	} else {
		err = fmt.Errorf("object doesn't have primary or secondary unique key")
		zap.S().Error(err)
	}
	return
}
