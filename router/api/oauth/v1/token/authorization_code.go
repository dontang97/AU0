package token

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/dontang97/AU0/util/auth"
	"github.com/dontang97/AU0/util/constant"
	"github.com/sirupsen/logrus"
)

type AuthCodePostReq struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
}

func AuthoRizationCodePostHandle(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	log, _ := ctx.Value(constant.RoutingRequestCtxFieldKeyLog{}).(*logrus.Logger)

	log.Debug("invoke token authorization_code post v1 handle")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithField("error", err.Error()).Error("read request body failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	aReq := &AuthCodePostReq{}
	if err = json.Unmarshal(body, aReq); err != nil {
		log.WithField("error", err.Error()).Error("parse request body failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	auther, _ := ctx.Value(constant.RoutingRequestCtxFieldKeyAuth{}).(*auth.Auth)

	tokenUrl := auther.OIDCProvider.Endpoint().TokenURL
	log.WithField("url", tokenUrl).Debug("auth token url")

	form := url.Values{}
	form.Add("grant_type", aReq.GrantType)
	form.Add("client_id", aReq.ClientID)
	form.Add("client_secret", aReq.ClientSecret)
	form.Add("code", aReq.Code)
	form.Add("redirect_uri", aReq.RedirectURI)
	req, err := http.NewRequest("POST", tokenUrl, strings.NewReader(form.Encode()))
	if err != nil {
		log.WithField("error", err.Error()).Error("build token request failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.WithField("error", err.Error()).Error("fire token request failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithField("error", err.Error()).Error("read token response failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(respBody); err != nil {
		log.WithField("error", err.Error()).Error("write response failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
