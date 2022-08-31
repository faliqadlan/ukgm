package webapiv1

import (
	"ukgm/internal/app/ukgm/service"
	simple "ukgm/internal/app/ukgm/webApiV1/simpleHandler"

	"github.com/go-chi/chi"
)

type Api struct {
	srv    *service.ServiceNew
	simple *simple.SimpleApi
}

func NewApi(service *service.ServiceNew, simple *simple.SimpleApi) (result *Api) {
	result = &Api{
		srv: service,
		simple:simple ,
	}
	return
}

func (api *Api) Init(r chi.Router) {
	r.Route("/simple", func(r chi.Router) {
		r.Route("/", api.simple.Init)
	})
}
