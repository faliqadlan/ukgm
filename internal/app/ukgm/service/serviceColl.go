package service

import "firebase.google.com/go/auth"

type StoreDataSrv interface {
	StoreData(data interface{}, auth *auth.Token) (result interface{}, err error)
}

type GetListsrv interface {
	GetList(input ...string) (result []interface{}, err error)
}
