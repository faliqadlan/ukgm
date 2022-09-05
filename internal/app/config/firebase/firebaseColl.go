package firebase_api

import "firebase.google.com/go/auth"

type FirebaseService interface {
	AuthenticationByID(uid string) *auth.UserRecord
	GenerateToken(uid string) string
	EncodeToken(token string) (auth *auth.Token, err error)
	SignInWithCustomToken(token string) (string, error)
	SignInWithPassword(email, password string) (string, error)
}
