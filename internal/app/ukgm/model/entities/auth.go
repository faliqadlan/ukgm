package entities

import "firebase.google.com/go/auth"

// type AuthUser struct {
// 	ID   string `json:"id" bson:"id,omitempty"`
// 	Name string `json:"name" bson:"name"`
// }

type Auth struct {
	Token    string `json:"token"`
	AuthUser *auth.Token
}

func GenerateAuth(authToken *auth.Token, token string) (*Auth) {
	return &Auth{
		Token:    token,
		AuthUser: authToken,
	}
}
