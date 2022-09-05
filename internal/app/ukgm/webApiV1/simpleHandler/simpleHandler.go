package simple

import (
	"errors"
	"log"
	"net/http"
	"ukgm/internal/app/ukgm/model/request"
	responses "ukgm/internal/app/ukgm/model/response"

	"firebase.google.com/go/auth"
	"go.uber.org/zap"
)

func (SimpeApi *SimpleApi) simpleHello(w http.ResponseWriter, r *http.Request) {
	// GET DATA CONTEXT MIDDLEWARE
	auth := r.Context().Value(request.GetAuthCtxKey()).(*auth.Token)

	log.Println(auth)

	query := r.URL.Query()
	errorTest := query.Get("error")

	if errorTest == "true" {
		err := errors.New("bad request")
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}

	// log.Println(authTo)

	responses.SendJSON(responses.NewResponse("Welcome to the simpleHello"), w)
}
