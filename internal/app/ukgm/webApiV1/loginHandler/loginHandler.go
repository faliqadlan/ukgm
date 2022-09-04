package login_handler

import (
	"errors"
	"net/http"

	"ukgm/internal/app/ukgm/model/entities"
	responses "ukgm/internal/app/ukgm/model/response"

	"go.uber.org/zap"
)

func (l *LoginApi) Login(w http.ResponseWriter, r *http.Request) {

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	tokenId, err := l.firebaseApi.SignInWithPassword(email, password)
	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(errors.New("password or email is invalid")), w)
	}

	authToken, err := l.firebaseApi.EncodeToken(tokenId)
	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(errors.New("invalid account")), w)
	}


	authUser := entities.GenerateAuth(authToken, tokenId)


	responses.SendJSON(responses.NewResponse(authUser), w)
}
