package router

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/dontang97/AU0/dto"
	"github.com/dontang97/AU0/router/api"
	"github.com/dontang97/AU0/util/auth"
	"github.com/dontang97/AU0/util/constant"
)

func Route(log *logrus.Logger, db *gorm.DB, crypt dto.TokenCryptor, auther *auth.Auth) *mux.Router {
	var rootMiddleFunc mux.MiddlewareFunc = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), constant.RoutingRequestCtxFieldKeyLog{}, log)
			ctx = context.WithValue(ctx, constant.RoutingRequestCtxFieldKeyDB{}, db)
			ctx = context.WithValue(ctx, constant.RoutingRequestCtxFieldKeyHTCTokenCryptor{}, crypt)
			ctx = context.WithValue(ctx, constant.RoutingRequestCtxFieldKeyAuth{}, auther)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}

	root := mux.NewRouter()
	root.Use(rootMiddleFunc)

	apiap := &api.API{}
	apiap.Serve(root.PathPrefix("/api").Subrouter())

	return root
}
