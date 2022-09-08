package entities

import "firebase.google.com/go/auth"

type User struct {
	UserAuth   *auth.Token
	ID         interface{} `json:"id" bson:"_id,omitempty"`
	UID        string      `json:"uid" bson:"_uid,omitempty"`
	DistrictID string      `json:"district_id" bson:"_district_id,omitempty"`
	Email      string      `json:"email" bson:"email"`
}

func NewUser(userAuth auth.Token, disctrictID string) (user User) {
	user = User{
		UserAuth: &userAuth,
		UID: userAuth.UID,
		
	}
	return
}