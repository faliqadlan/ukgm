package user_api

import (
	"context"
	"net/http"
	"ukgm/internal/app/ukgm/model/entities"
	"ukgm/internal/app/ukgm/model/request"
	responses "ukgm/internal/app/ukgm/model/response"

	"firebase.google.com/go/auth"
	"go.uber.org/zap"
)

func (a *UserApi) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// GET DATA CONTEXT MIDDLEWARE
		auth := r.Context().Value(request.GetAuthCtxKey()).(*auth.Token)

		data, err := a.srv.ShowUserByUid(auth.UID)
		if err != nil {
			zap.S().Error(err)
			responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
			return
		}

		ctx := context.WithValue(r.Context(), request.GetUserCtxKey(), data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *UserApi) ShowUserByUid(w http.ResponseWriter, r *http.Request) {
	// get data context middleware
	user := r.Context().Value(request.GetUserCtxKey()).(*entities.User)

	responses.SendJSON(responses.NewResponse(user), w)
}
