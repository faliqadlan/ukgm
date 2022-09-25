package service

import "firebase.google.com/go/auth"

type StoreDataSrv interface {
	StoreData(data interface{}, auth *auth.Token) (result interface{}, err error)
}

type GetListsrv interface {
	GetList(input ...string) (result interface{}, err error)
}

type GetDetailSrv interface {
	GetDetail(input ...string) (result interface{}, err error)
}

type FetchIndexSrv interface {
	FetchIndex(input ...interface{}) (result interface{}, count int, err error)
}
