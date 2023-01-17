package simple

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"ukgm/internal/app/ukgm/model/request"
	responses "ukgm/internal/app/ukgm/model/response"

	"firebase.google.com/go/auth"
	"github.com/go-co-op/gocron"
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

func MaintainHelloOnLogin() {
	url := "https://ukgm-faliqadlan.cloud.okteto.net/api/test"

	// init request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		zap.S().Error(err)
	}

	// EXECUTE REQUEST
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zap.S().Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.S().Error(err)
	}

	log.Println(string(body))
}

func CronMaintainHelloOnLogin() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(12).Hours().Do(func() {
		MaintainHelloOnLogin()
	})

	s.StartAsync()
}
