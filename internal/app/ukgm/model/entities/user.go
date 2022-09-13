package entities

import "firebase.google.com/go/auth"

type User struct {
	ID         interface{} `json:"id" bson:"_id,omitempty"`
	UID        string      `json:"uid" bson:"uid,omitempty"`
	DistrictID string      `json:"district_id" bson:"district_id,omitempty"`
	Email      string      `json:"email" bson:"email,omitempty"`

	// include
	Timestamp
}

type SimpleUser struct {
}

func (u *User) CreateNew(userAuth auth.Token, disctrictID string) {
	// user = User{
	// 	UserAuth:   &userAuth,
	// 	UID:        userAuth.UID,
	// 	DistrictID: disctrictID,
	// 	Email:      userAuth.Claims["email"].(string),
	// }

	u.UID = userAuth.UID
	u.DistrictID = disctrictID
	u.Email = userAuth.Claims["email"].(string)
}
