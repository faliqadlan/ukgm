package request

import (
	"context"
	"errors"
	"net/http"
	"strings"

	firebase_api "ukgm/internal/app/config/firebase"
	responses "ukgm/internal/app/ukgm/model/response"

	"go.uber.org/zap"
)

/*
|--------------------------------------------------------------------------
| Middlewares
|--------------------------------------------------------------------------
|
*/

type Middleware struct {
	firebaseApi *firebase_api.FirebaseApi
}

func NewMiddleware(firebaseApi *firebase_api.FirebaseApi) *Middleware {
	return &Middleware{
		firebaseApi: firebaseApi,
	}
}

func (m *Middleware) CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		splitToken := strings.Split(token, "Bearer ")
		if len(token) <= 0 {
			token = r.URL.Query().Get("token")

			if len(token) <= 0 {
				responses.SendErrorJSON(responses.PrintErrorMessage(errors.New("401")), w)
				return
			}
		} else {
			token = splitToken[1]
		}

		tokenID, err := m.firebaseApi.SignInWithCustomToken(token)
		if err != nil {
			zap.S().Error(err)
			responses.SendErrorJSON(responses.PrintErrorMessage(errors.New("token is invalid")), w)
		}

		authToken, err := m.firebaseApi.EncodeToken(tokenID)
		if err != nil {
			zap.S().Error(err)
			responses.SendErrorJSON(responses.PrintErrorMessage(errors.New("token is invalid")), w)
		}

		ctx := context.WithValue(r.Context(), GetAuthCtxKey(), authToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
