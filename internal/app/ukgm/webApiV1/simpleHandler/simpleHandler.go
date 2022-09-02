package simple

import (
	"errors"
	"net/http"
	responses "ukgm/internal/app/ukgm/model/response"

	"go.uber.org/zap"
)

func (SimpeApi *SimpleApi) simpleHello(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	errorTest := query.Get("error")

	if errorTest == "true" {
		err := errors.New("bad request")
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}

	responses.SendJSON(responses.NewResponse("Welcome to the simpleHello"), w)
}
