package request

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"ukgm/internal/app/config/env"
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
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		firebaseApp := firebase_api.SetupFirebase()

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

		tokenID, err := firebase_api.SignInWithCustomToken(token, env.Get("WEB_API_KEY"))
		if err != nil {
			zap.S().Error(err)
			responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		}

		authToken, err := firebase_api.EncodeToken(firebaseApp, tokenID)
		if err != nil {
			zap.S().Error(err)
			responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		}

		ctx := context.WithValue(r.Context(), GetAuthCtxKey(), authToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
