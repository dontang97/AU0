package auth

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"gopkg.in/auth0.v5/management"
)

type Auth struct {
	OIDCProvider *oidc.Provider
	Manager      *management.Management
	//ClientID     string
}

// New instantiates the *Authenticator.
func New(domain string, clientID string, secret string) (*Auth, error) {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+domain+"/",
	)
	if err != nil {
		return nil, err
	}

	auth0Mgr, err := management.New(
		domain,
		management.WithClientCredentials(
			clientID,
			secret,
		),
		//management.WithContext(context.Background()),
		//management.WithDebug(true),
	)
	if err != nil {
		return nil, err
	}

	return &Auth{
		OIDCProvider: provider,
		Manager:      auth0Mgr,
		//ClientID:     clientID,
	}, nil
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (auth *Auth) VerifyIDToken(ctx context.Context, token string) (*oidc.IDToken, error) {
	oidcConfig := &oidc.Config{
		SkipClientIDCheck: true,
	}

	return auth.OIDCProvider.Verifier(oidcConfig).Verify(ctx, token)
}

func (auth *Auth) UserInfo(ctx context.Context, token oauth2.TokenSource) (*oidc.UserInfo, error) {
	return auth.OIDCProvider.UserInfo(ctx, token)
}
