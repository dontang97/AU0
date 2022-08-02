package oauth

import (
	"context"
	"net/http"

	"github.com/dontang97/AU0/dto"
	v1 "github.com/dontang97/AU0/router/api/oauth/v1"
	"github.com/gorilla/mux"
)

type OAuth struct {
	Cryptor dto.TokenCryptor
}

func (oauth *OAuth) Serve(router *mux.Router) {
	var middleFunc mux.MiddlewareFunc = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(
				context.WithValue(r.Context(), "htc_token_cryptor", oauth.Cryptor),
			)

			next.ServeHTTP(w, r)
		})
	}

	router.Use(middleFunc)

	v1ap := &v1.V1{}

	v1ap.Serve(router.PathPrefix("/v1").Subrouter())
}
