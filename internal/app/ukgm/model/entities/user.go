package entities

import (
	"github.com/faliqadlan/ukgm/internal/app/utils"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	UserUrl = "/api/v1/users"
)

type (
	User struct {
		ID     interface{} `json:"id" bson:"_id,omitempty" part:"main" unique:"primary" properties:"constant"`
		Name   string      `json:"name" bson:"name" part:"main" properties:"searchable"`
		Code   string      `json:"code" bson:"code" part:"main" properties:"searchable,constant" unique:"secondary"`
		Gender string      `json:"gender" bson:"gender" part:"detail" properties:"searchable"`
		Timestamp
	}

	UserRef struct {
		ID   string `json:"id" bson:"id"`
		Name string `json:"name" bson:"name"`
	}
)

// FUNCTION

func (user *User) GetSearchableField() (result []string) {
	if user == nil {
		user = &User{}
	}
	result, err := utils.GetSearchableField(user)
	if err != nil {
		result = []string{}
	}
	return
}

func (user *User) GetUniqueFilterMongodb() (result bson.M) {
	result, _ = utils.ConstructUniqueFilter(user)
	return
}

func (user *User) GetMainOrDetail() (main, detail []string, err error) {
	main, detail, err = utils.GetMainOrDetail(*user, user)
	return
}

func (user *User) CreateReference() *UserReference {
	userRef := &UserReference{
		FullName: user.Name,
		ID:       utils.ObjectIDHex(user.ID),
	}
	return userRef
}
