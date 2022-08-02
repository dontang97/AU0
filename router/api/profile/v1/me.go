package v1

import (
	"encoding/json"
	"net/http"

	"github.com/dontang97/AU0/util/auth"
	"github.com/dontang97/AU0/util/constant"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type ProfileGetReq struct{}

func (r *ProfileGetReq) Token() (*oauth2.Token, error) {
	return &oauth2.Token{}, nil
}

func ProfileGetHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log, _ := ctx.Value(constant.RoutingRequestCtxFieldKeyLog{}).(*logrus.Logger)

	log.Debug("invoke profile get v1 me handle")

	auther, _ := ctx.Value(constant.RoutingRequestCtxFieldKeyAuth{}).(*auth.Auth)

	pReq := &ProfileGetReq{}
	info, err := auther.UserInfo(ctx, pReq)
	if err != nil {
		log.WithField("error", err.Error()).Error("get user info failed")
		w.WriteHeader(http.StatusFailedDependency)
	}

	body := map[string]interface{}{}

	body["info"] = info

	resp, err := json.Marshal(body)
	if err != nil {
		log.WithField("error", err.Error()).Error("marshal response failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(resp); err != nil {
		log.WithField("error", err.Error()).Error("write response failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
