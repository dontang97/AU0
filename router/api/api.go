package api

import (
	"github.com/dontang97/AU0/router/api/oauth"
	"github.com/gorilla/mux"
)

type API struct {
}

func (api *API) Serve(router *mux.Router) {
	oauthap := &oauth.OAuth{}

	oauthap.Serve(router.PathPrefix("/oauth").Subrouter())
}
