package user_api

import (
	"ukgm/internal/app/ukgm/service"

	"github.com/go-chi/chi"
)

type UserApi struct {
	srv *service.ServiceNew
}

func NewUserApi(srv *service.ServiceNew) *UserApi {
	return &UserApi{srv: srv}
}

func (u *UserApi) Init(r chi.Router) {
	r.Use(u.UserCtx)
	r.Get("/", u.ShowUserByUid)
}
