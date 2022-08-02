package oidc

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dontang97/AU0/util/auth"
	"github.com/dontang97/AU0/util/constant"
	"github.com/sirupsen/logrus"
)

type AuthkeyPostReq struct {
	IdToken string `json:"id_token"`
	//AppID   string `json:"app_id"`
}

func AuthkeyPostHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log, _ := ctx.Value(constant.RoutingRequestCtxFieldKeyLog{}).(*logrus.Logger)

	log.Debug("invoke oidc authkey post v1 handle")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithField("error", err.Error()).Error("read request body failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	aReq := &AuthkeyPostReq{}
	if err = json.Unmarshal(body, aReq); err != nil {
		log.WithField("error", err.Error()).Error("parse request body failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/*
		tokenSplits := strings.Split(aReq.IdToken, ".")
		if len(tokenSplits) != 3 {
			log.WithField("error", err.Error()).Error("invalid id token (split by dot should get 3 strings)")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		claimsJson, err := b64.StdEncoding.DecodeString(tokenSplits[1])
		if err != nil {
			log.WithField("error", err.Error()).Error("base64 decode claims failed")
			log.WithField("token", tokenSplits[1]).Error("ggggggggggggggggggggggg")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		claimsMap := map[string]interface{}{}
		if err = json.Unmarshal(claimsJson, &claimsMap); err != nil {
			log.WithField("error", err.Error()).Error("parse claims json failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		clientId, ok := claimsMap["aud"]
		if !ok {
			log.Error("no aud in claims")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		clientIdStr, ok := clientId.(string)
		if !ok {
			log.WithField("error", err.Error()).Error("aud is not string format")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	*/
	auther, _ := ctx.Value(constant.RoutingRequestCtxFieldKeyAuth{}).(*auth.Auth)

	idToken, err := auther.VerifyIDToken(ctx, aReq.IdToken)
	if err != nil {
		log.WithField("error", err.Error()).Error("verify id token failed")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := auther.Manager.User.Read(idToken.Subject)
	if err != nil {
		log.WithField("error", err.Error()).Error("get user failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(user)
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
