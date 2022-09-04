package login_handler

import (
	firebase_api "ukgm/internal/app/config/firebase"

	"github.com/go-chi/chi"
)

type LoginApi struct {
	firebaseApi *firebase_api.FirebaseApi
}

func NewLoginApi(firebaseApi *firebase_api.FirebaseApi) *LoginApi {

	return &LoginApi{
		firebaseApi: firebaseApi,
	}
}

func (api *LoginApi) Init(r chi.Router) {
	r.Post("/", api.Login)
}
