package v1

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dontang97/AU0/router/api/oauth/v1/oidc"
	"github.com/dontang97/AU0/router/api/oauth/v1/token"
)

type V1 struct {
}

func (v1 *V1) Serve(router *mux.Router) {
	router.HandleFunc("/token/callback", token.CallbackGetHandle).Methods(http.MethodGet)
	router.HandleFunc("/token/callback", token.CallbackPostHandle).Methods(http.MethodPost)

	router.HandleFunc("/token/authorization_code", token.AuthoRizationCodePostHandle).Methods(http.MethodPost)

	router.HandleFunc("/oidc/authkey", oidc.AuthkeyPostHandle).Methods(http.MethodPost)
}
