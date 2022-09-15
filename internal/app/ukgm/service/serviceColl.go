package service

import "firebase.google.com/go/auth"

type ServiceCollection interface {
	StoreData(data interface{}, auth *auth.Token) (result interface{}, err error)
}
