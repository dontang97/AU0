package token

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/biter777/countries"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dontang97/AU0/dto"
	"github.com/dontang97/AU0/util/auth"
	"github.com/dontang97/AU0/util/cache"
	"github.com/dontang97/AU0/util/constant"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const html = `
<html>
	<head>
		<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>
		<script>
		$(document).ready(function() {
			const Url = 'http://localhost:9998/api/oauth/v1/token/callback';
			const data = window.location.hash.split("#")[1];
		
			$.post(Url, data, function(_, status){
				console.log(status);
			})
		});
		</script>
	</head>
	<body>
	</body>
</html>
`

func CallbackGetHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log, _ := ctx.Value(constant.RoutingRequestCtxFieldKeyLog{}).(*logrus.Logger)

	log.Debug("invoke token get v1 callback handle")

	w.Header().Set("Content-Type", "text/html")

	if _, err := w.Write([]byte(html)); err != nil {
		log.WithField("error", err.Error()).Error("write response body failed")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// superfluous WriteHeader call
	// w.WriteHeader(http.StatusOK)

	return
}

type State struct {
	AppID       string `json:"app_id"`
	RedirectURI string `json:"redirect_uri"`
	State       string `json:"state"`
}

type CallbackPostReq struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	Scope        []string      `json:"scope"`
	ExpiresIn    time.Duration `json:"expires_in"`
	Expiry       time.Time     `json:"expiry"`
	TokenType    string        `json:"token_type"`
	State        State         `json:"state"`
	IdToken      string        `json:"id_token"`
}

func (r *CallbackPostReq) parse(form url.Values) error {
	r.AccessToken = form.Get("access_token")
	r.RefreshToken = form.Get("refresh_token")
	r.Scope = strings.Split(form.Get("scope"), " ")

	if s := form.Get("expires_in"); s != "" {
		expiresIn, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		r.ExpiresIn = time.Second * time.Duration(expiresIn)
		r.Expiry = time.Now().Add(r.ExpiresIn)
	}

	r.TokenType = form.Get("token_type")

	if state := form.Get("state"); state != "" {
		if err := json.Unmarshal([]byte(state), &r.State); err != nil {
			return err
		}
	}

	r.IdToken = form.Get("id_token")
	return nil
}

func (r *CallbackPostReq) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken:  r.AccessToken,
		TokenType:    r.TokenType,
		RefreshToken: r.RefreshToken,
		Expiry:       r.Expiry,
	}, nil
}

// TODO: database support account id mapping
func toClaims(r CallbackPostReq, userInfo *oidc.UserInfo) dto.Claims {
	var accountID uuid.UUID
	if idStr, ok := cache.UserID2AccountIDMap[userInfo.Subject]; ok {
		accountID = uuid.MustParse(idStr)
	} else {
		accountID = uuid.New()
		cache.UserID2AccountIDMap[userInfo.Subject] = accountID.String()
		cache.AccountID2UserIDMap[accountID.String()] = userInfo.Subject
	}

	var clientID uuid.UUID
	if idStr, ok := cache.AppID2ClientIDMap[r.State.AppID]; ok {
		clientID = uuid.MustParse(idStr)
	} else {
		clientID = uuid.New()
		cache.AppID2ClientIDMap[r.State.AppID] = clientID.String()
		cache.ClientID2AppIDMap[clientID.String()] = r.State.AppID
	}

	claims := dto.Claims{
		IssuedAt:        r.Expiry.Add(-1 * r.ExpiresIn).UnixMilli(),
		AccountID:       accountID,
		VirtualDeviceID: "",
		HandsetDeviceID: "",
		Verified:        userInfo.EmailVerified,
		HandSetVerified: false,
		LegalDocToSign:  false,
		CountryCode:     countries.CountryCode(0),
		DataCenterID:    uuid.MustParse(constant.RAMDataCenterID),
		IDKey:           "",
		IssueTo:         clientID,
		Scopes:          r.Scope,
		ExpiredInterval: int(r.ExpiresIn.Seconds()),
	}
	return claims
}

func CallbackPostHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log, _ := ctx.Value(constant.RoutingRequestCtxFieldKeyLog{}).(*logrus.Logger)

	log.Debug("invoke token post v1 callback handle")

	if err := r.ParseForm(); err != nil {
		log.WithField("error", err.Error()).Error("parse url-encoded form body failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cReq := &CallbackPostReq{}
	if err := cReq.parse(r.PostForm); err != nil {
		log.WithField("error", err.Error()).Error("parse url-encoded form to CallbackPostReq failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	auther, _ := ctx.Value(constant.RoutingRequestCtxFieldKeyAuth{}).(*auth.Auth)

	idToken, err := auther.VerifyIDToken(ctx, cReq.IdToken)
	if err != nil {
		log.WithField("error", err.Error()).Error("verify id token failed")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// TODO: convert to HTC token
	userInfo, err := auther.UserInfo(ctx, cReq)
	if err != nil {
		log.WithField("error", err.Error()).Error("get user info failed")
		w.WriteHeader(http.StatusFailedDependency)
	}

	claims := toClaims(*cReq, userInfo)

	cryptor, _ := ctx.Value(constant.RoutingRequestCtxFieldKeyHTCTokenCryptor{}).(dto.TokenCryptor)
	token, err := claims.Encrypt(cryptor)

	body := map[string]interface{}{}
	body["id_token"] = idToken
	body["htc_token"] = string(token)

	_, err = json.Marshal(body)
	if err != nil {
		log.WithField("error", err.Error()).Error("marshal response failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/*
		if _, err = w.Write(resp); err != nil {
			log.WithField("error", err.Error()).Error("write response failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	*/

	location, err := http.NewRequest("GET", cReq.State.RedirectURI, nil)
	if err != nil {
		log.WithField("error", err.Error()).Error("build location request failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	//w.Header().Set("Location", "https://account-stage.htcvive.com/"+"#access_token="+url.QueryEscape(string(token)))
	w.Header().Set("Location", location.URL.String()+"#access_token="+url.QueryEscape(string(token)))

	w.WriteHeader(http.StatusFound)

	// superfluous WriteHeader call
	// w.WriteHeader(http.StatusOK)

	return
}
