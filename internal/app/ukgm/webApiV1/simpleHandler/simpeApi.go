package simple

import (
	"ukgm/internal/app/ukgm/service"

	"github.com/go-chi/chi"
)

type SimpleApi struct {
	srv *service.ServiceNew
}

func NewSimpleApi(srv *service.ServiceNew) *SimpleApi {
	return &SimpleApi{srv: srv}
}

func (api SimpleApi) Init(r chi.Router) {
	r.Get("/simpleHello", api.simpleHello)
}